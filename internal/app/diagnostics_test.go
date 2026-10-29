package app

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	"go.opentelemetry.io/collector/client"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

func diagnosticFixture(t *testing.T) (*Service, http.Handler, *http.Cookie) {
	t.Helper()
	store, err := Open(filepath.Join(t.TempDir(), "db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { store.DB.Close() })
	svc, err := New(store, "synthetic owner passphrase", "browser-test", "server-test", "http://localhost")
	if err != nil {
		t.Fatal(err)
	}
	handler := svc.Handler("")
	request := httptest.NewRequest("POST", "http://localhost/api/login", strings.NewReader(`{"passphrase":"synthetic owner passphrase"}`))
	request.Header.Set("Origin", svc.Origin)
	result := httptest.NewRecorder()
	handler.ServeHTTP(result, request)
	if result.Code != 200 {
		t.Fatal(result.Code)
	}
	return svc, handler, result.Result().Cookies()[0]
}
func diagnosticContext() context.Context {
	return client.NewContext(context.Background(), client.Info{Metadata: client.NewMetadata(map[string][]string{"authorization": {"Bearer server-test"}})})
}
func errorRequest(handler http.Handler, cookie *http.Cookie, q string) *httptest.ResponseRecorder {
	request := httptest.NewRequest("GET", "http://localhost/api/errors"+q, nil)
	if cookie != nil {
		request.AddCookie(cookie)
	}
	result := httptest.NewRecorder()
	handler.ServeHTTP(result, request)
	return result
}

func TestDiagnosticLogsScrubPersistRestartAndAppWideReadWithoutJourney(t *testing.T) {
	svc, handler, cookie := diagnosticFixture(t)
	now := time.Now().UTC().Add(-time.Second)
	logs := plog.NewLogs()
	records := logs.ResourceLogs().AppendEmpty().ScopeLogs().AppendEmpty().LogRecords()
	info := records.AppendEmpty()
	info.SetEventName("info.body")
	info.SetTimestamp(pcommon.NewTimestampFromTime(now))
	info.Body().SetStr("ORDINARY-INFO-BODY-MUST-BE-OMITTED")
	failure := records.AppendEmpty()
	failure.SetEventName("synthetic.failure")
	failure.SetTimestamp(pcommon.NewTimestampFromTime(now))
	failure.SetSeverityNumber(plog.SeverityNumberError)
	failure.Attributes().PutStr("event.id", "diagnostic-log")
	failure.Attributes().PutStr("exception.type", "SyntheticWriteError")
	failure.Attributes().PutStr("exception.message", "write failed password=PASSWORD_SENTINEL user=alex@example.test at https://example.test/users/PATH_SENTINEL?token=QUERY_SENTINEL <script>alert(1)</script>")
	failure.Attributes().PutStr("exception.stacktrace", "SyntheticWriteError: bad write\n    at SaveOrder (/Users/PRIVATE_USER/project/src/orders.ts:42:9)\n    at HandleRequest (https://example.test/users/PRIVATE_ROUTE/app.js?token=QUERY_SENTINEL:12:3)\nCookie: session=COOKIE_SENTINEL\nAuthorization: Bearer AUTH_SENTINEL")
	failure.Attributes().PutStr("private.attribute", "UNKNOWN-ATTR-MUST-BE-OMITTED")
	failure.Body().SetStr("write operation failed token=BODY_SECRET_SENTINEL")
	if err := svc.consumeLogs(diagnosticContext(), logs); err != nil {
		t.Fatal(err)
	}
	retained, err := svc.Store.Evidence()
	if err != nil {
		t.Fatal(err)
	}
	if len(retained) != 2 {
		t.Fatal(len(retained))
	}
	attrs := retained[1].Attrs
	if attrs["exception.message"] == "" || attrs["exception.stacktrace"] == "" || attrs["log.message"] == "" || attrs["diagnostic.policy"] != "scrubbed-v1" {
		t.Fatalf("diagnostic projection missing: %+v", attrs)
	}
	if attrs["private.attribute"] != "" || retained[0].Attrs["log.message"] != "" {
		t.Fatal("ordinary body or unknown attr retained")
	}
	if !strings.Contains(attrs["exception.stacktrace"], "SaveOrder") || !strings.Contains(attrs["exception.stacktrace"], "orders.ts:42:9") {
		t.Fatal("stack lost useful symbols/source basenames")
	}
	if retained[1].Scope != "" {
		t.Fatal("unscoped error acquired fabricated journey")
	}
	path := svc.Store.Path
	svc.Store.DB.Close()
	reopened, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer reopened.DB.Close()
	svc.Store = reopened
	if response := errorRequest(handler, nil, ""); response.Code != 401 {
		t.Fatal("errors read not protected")
	}
	response := errorRequest(handler, cookie, "")
	if response.Code != 200 {
		t.Fatalf("errors status%d %s", response.Code, response.Body)
	}
	raw := response.Body.String()
	for _, sentinel := range []string{"PASSWORD_SENTINEL", "PATH_SENTINEL", "QUERY_SENTINEL", "PRIVATE_USER", "PRIVATE_ROUTE", "COOKIE_SENTINEL", "AUTH_SENTINEL", "BODY_SECRET_SENTINEL", "alex@example.test", "UNKNOWN-ATTR-MUST-BE-OMITTED", "ORDINARY-INFO-BODY-MUST-BE-OMITTED"} {
		if strings.Contains(raw, sentinel) {
			t.Errorf("diagnostic privacy leak: %s", sentinel)
		}
	}
	var got struct {
		Records      []Evidence
		Total        int
		Truncated    bool
		Unit, Policy string
	}
	if err = json.Unmarshal(response.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.Total != 1 || len(got.Records) != 1 || got.Records[0].Record != "diagnostic-log" || got.Unit != "diagnostic records" || got.Policy == "" {
		t.Fatalf("wrong appwide diagnostic result: %+v", got)
	}
	if strings.Contains(raw, "<script>") || !strings.Contains(got.Records[0].Attrs["exception.message"], "<script>") {
		t.Fatal("diagnostic text must remain inert JSON text, not markup")
	}
}

func TestDiagnosticFailedSpansAndBoundedErrorSearchExpiry(t *testing.T) {
	svc, handler, cookie := diagnosticFixture(t)
	now := time.Now().UTC()
	svc.Store.Now = func() time.Time { return now }
	traces := ptrace.NewTraces()
	rs := traces.ResourceSpans().AppendEmpty()
	rs.Resource().Attributes().PutStr("service.version", "v1")
	spans := rs.ScopeSpans().AppendEmpty().Spans()
	span := spans.AppendEmpty()
	span.SetName("database.write")
	span.SetStartTimestamp(pcommon.NewTimestampFromTime(now.Add(-time.Second)))
	span.SetTraceID(pcommon.TraceID{1})
	span.SetSpanID(pcommon.SpanID{2})
	span.Status().SetCode(ptrace.StatusCodeError)
	span.Status().SetMessage("write timeout client_secret=SPAN_SECRET_SENTINEL")
	span.Attributes().PutStr("session.id", "operation")
	if err := svc.consumeTraces(diagnosticContext(), traces); err != nil {
		t.Fatal(err)
	}
	response := errorRequest(handler, cookie, "?scope=operation&version=v1&definition=1&q=timeout")
	if response.Code != 200 {
		t.Fatalf("errors status%d %s", response.Code, response.Body)
	}
	var got struct {
		Records   []Evidence
		Total     int
		Truncated bool
	}
	json.Unmarshal(response.Body.Bytes(), &got)
	if got.Total != 1 || len(got.Records) != 1 || got.Records[0].Severity != 17 || got.Records[0].Trace == "" || got.Records[0].Span == "" || got.Records[0].Attrs["exception.message"] == "" || strings.Contains(response.Body.String(), "SPAN_SECRET_SENTINEL") {
		t.Fatalf("failed span omitted or unsanitized: %+v", got)
	}
	for _, q := range []string{"?scope=other", "?version=v2", "?definition=2", "?q=notfound"} {
		r := errorRequest(handler, cookie, q)
		json.Unmarshal(r.Body.Bytes(), &got)
		if r.Code != 200 || got.Total != 0 {
			t.Errorf("error filter ignored %s", q)
		}
	}
	for _, q := range []string{"?q=" + strings.Repeat("a", 161), "?definition=bad", "?definition=0", "?from=invalid", "?to=invalid", "?from=2026-10-29T00:00:00Z&to=2026-10-28T00:00:00Z"} {
		if r := errorRequest(handler, cookie, q); r.Code != 400 {
			t.Errorf("invalid errors query accepted: %s status%d", q, r.Code)
		}
	}
	svc.Store.Now = func() time.Time { return now.Add(30 * 24 * time.Hour) }
	response = errorRequest(handler, cookie, "")
	json.Unmarshal(response.Body.Bytes(), &got)
	if response.Code != 200 || got.Total != 0 {
		t.Fatal("expired diagnostic survived protected read")
	}
}

func TestDiagnosticTextLimitsPreserveUTF8AndNormalizeControls(t *testing.T) {
	svc, _, _ := diagnosticFixture(t)
	logs := plog.NewLogs()
	entry := logs.ResourceLogs().AppendEmpty().ScopeLogs().AppendEmpty().LogRecords().AppendEmpty()
	entry.SetEventName("limits.error")
	entry.SetTimestamp(pcommon.NewTimestampFromTime(time.Now()))
	entry.SetSeverityNumber(plog.SeverityNumberError)
	entry.Attributes().PutStr("exception.message", strings.Repeat("é", 40000)+"\x00\x1b[31m")
	entry.Attributes().PutStr("exception.stacktrace", strings.Repeat("at Function (/home/SECRET_USER/src/main.go:42)\r\n", 100))
	entry.Body().SetStr("body\x00\x1b[31m\r\npassword=LIMIT_SECRET")
	if err := svc.consumeLogs(diagnosticContext(), logs); err != nil {
		t.Fatal(err)
	}
	records, err := svc.Store.Evidence()
	if err != nil {
		t.Fatal(err)
	}
	a := records[0].Attrs
	if len(a["exception.message"]) > 2048 || len(a["exception.stacktrace"]) > 8192 || len(strings.Split(a["exception.stacktrace"], "\n")) > 40 || !utf8.ValidString(a["exception.message"]) || a["diagnostic.truncated"] != "true" {
		t.Fatalf("diagnostic limits absent: %+v", a)
	}
	if a["exception.message"] == "" || a["exception.stacktrace"] == "" || strings.Contains(a["log.message"], "\x00") || strings.Contains(a["log.message"], "\x1b") || strings.Contains(a["log.message"], "LIMIT_SECRET") || strings.Contains(a["exception.stacktrace"], "SECRET_USER") {
		t.Fatal("diagnostic normalization/scrubbing failed")
	}
}

func TestDiagnosticScrubKeepsDotNetSymbolsAndRedactsQuotedCredentials(t *testing.T) {
	stack := `SyntheticException: failed
   at Example.Repository.SaveAsync() in C:\Users\PRIVATE_WINDOWS_USER\project\Repository.cs:line 42
   at Example.API.Handle() in /home/PRIVATE_LINUX_USER/src/Handler.cs:line 19
{"cookie":"jwt=QUOTED_COOKIE_SENTINEL; opaque=OTHER_COOKIE_SENTINEL", "api_key":"QUOTED_KEY_SENTINEL", "password":"TWO WORD SENTINEL"}`
	got, _ := scrubDiagnostic(stack, true)
	for _, sentinel := range []string{"PRIVATE_WINDOWS_USER", "PRIVATE_LINUX_USER", "QUOTED_COOKIE_SENTINEL", "OTHER_COOKIE_SENTINEL", "QUOTED_KEY_SENTINEL", "TWO WORD SENTINEL"} {
		if strings.Contains(got, sentinel) {
			t.Errorf("scrub leaked %s", sentinel)
		}
	}
	for _, symbol := range []string{"Example.Repository.SaveAsync()", "Repository.cs:line 42", "Example.API.Handle()", "Handler.cs:line 19"} {
		if !strings.Contains(got, symbol) {
			t.Errorf("lost safe stack symbol %s: %s", symbol, got)
		}
	}
}

func TestDiagnosticLegacyReplayKeepsOriginalProjectionAndConflicts(t *testing.T) {
	svc, _, _ := diagnosticFixture(t)
	now := time.Now().UTC()
	svc.Store.Now = func() time.Time { return now }
	legacy := Evidence{Source: "server", Signal: "logs", Name: "legacy.error", Record: "legacy-log", Time: now.UnixNano(), Severity: 17, Attrs: map[string]string{"exception.type": "LegacyError"}}
	if err := svc.Store.Insert([]Evidence{legacy}); err != nil {
		t.Fatal(err)
	}
	logs := plog.NewLogs()
	entry := logs.ResourceLogs().AppendEmpty().ScopeLogs().AppendEmpty().LogRecords().AppendEmpty()
	entry.SetEventName(legacy.Name)
	entry.SetTimestamp(pcommon.NewTimestampFromTime(now))
	entry.SetSeverityNumber(plog.SeverityNumberError)
	entry.Attributes().PutStr("event.id", legacy.Record)
	entry.Attributes().PutStr("exception.type", "LegacyError")
	entry.Attributes().PutStr("exception.message", "newly projected text")
	entry.Body().SetStr("newly projected body")
	if err := svc.consumeLogs(diagnosticContext(), logs); err != nil {
		t.Fatalf("upgrade broke qualified legacy replay: %v", err)
	}
	got, err := svc.Store.Evidence()
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].Attrs["exception.message"] != "" || got[0].Attrs["diagnostic.policy"] != "" {
		t.Fatal("legacy replay backfilled or duplicated old projection")
	}
	entry.SetSeverityNumber(plog.SeverityNumberFatal)
	if err := svc.consumeLogs(diagnosticContext(), logs); err == nil {
		t.Fatal("legacy replay lost log-severity conflict detection")
	}
	entry.SetSeverityNumber(plog.SeverityNumberError)
	entry.Attributes().PutStr("exception.type", "DifferentError")
	if err := svc.consumeLogs(diagnosticContext(), logs); err == nil {
		t.Fatal("legacy replay lost existing attribute conflict detection")
	}
	entry.Attributes().PutStr("event.id", "new-policy-log")
	entry.Attributes().PutStr("exception.type", "LegacyError")
	if err := svc.consumeLogs(diagnosticContext(), logs); err != nil {
		t.Fatal(err)
	}
	entry.Attributes().PutStr("exception.message", "conflicting new-policy text")
	if err := svc.consumeLogs(diagnosticContext(), logs); err == nil {
		t.Fatal("new privacy-policy replay ignored changed diagnostic text")
	}
	trace, spanID := pcommon.TraceID{3}, pcommon.SpanID{4}
	traceText, spanText := sourceIDs(trace, spanID)
	oldSpan := Evidence{Source: "server", Signal: "traces", Name: "legacy.span", Record: traceText + "/" + spanText, Trace: traceText, Span: spanText, Time: now.UnixNano(), Attrs: map[string]string{}}
	if err := svc.Store.Insert([]Evidence{oldSpan}); err != nil {
		t.Fatal(err)
	}
	traces := ptrace.NewTraces()
	span := traces.ResourceSpans().AppendEmpty().ScopeSpans().AppendEmpty().Spans().AppendEmpty()
	span.SetName(oldSpan.Name)
	span.SetStartTimestamp(pcommon.NewTimestampFromTime(now))
	span.SetTraceID(trace)
	span.SetSpanID(spanID)
	span.Status().SetCode(ptrace.StatusCodeError)
	span.Status().SetMessage("legacy span status was not previously projected")
	if err := svc.consumeTraces(diagnosticContext(), traces); err != nil {
		t.Fatalf("upgrade broke legacy span replay: %v", err)
	}
	span.SetName("changed-original-name")
	if err := svc.consumeTraces(diagnosticContext(), traces); err == nil {
		t.Fatal("legacy trace replay ignored original field conflict")
	}
	got, err = svc.Store.Evidence()
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 3 || got[2].Severity != 0 || got[2].Attrs["diagnostic.policy"] != "" {
		t.Fatal("legacy span replay overwrote original projection")
	}
}

func TestDiagnosticReadLimitsAndExplicitVisitDeletion(t *testing.T) {
	svc, handler, cookie := diagnosticFixture(t)
	now := time.Now().UTC()
	svc.Store.Now = func() time.Time { return now }
	logs := plog.NewLogs()
	records := logs.ResourceLogs().AppendEmpty().ScopeLogs().AppendEmpty().LogRecords()
	for i := 0; i < 101; i++ {
		entry := records.AppendEmpty()
		entry.SetEventName("bounded.error")
		entry.SetTimestamp(pcommon.NewTimestampFromTime(now.Add(-time.Second)))
		entry.SetSeverityNumber(plog.SeverityNumberError)
		entry.Attributes().PutStr("event.id", fmt.Sprint("real-", i))
		entry.Attributes().PutStr("session.id", "bounded-operation")
		entry.Body().SetStr("bounded error details")
	}
	if err := svc.consumeLogs(diagnosticContext(), logs); err != nil {
		t.Fatal(err)
	}
	response := errorRequest(handler, cookie, "?q=bounded")
	var got struct {
		Records   []Evidence
		Total     int
		Truncated bool
	}
	json.Unmarshal(response.Body.Bytes(), &got)
	if response.Code != 200 || got.Total != 101 || len(got.Records) != 100 || !got.Truncated {
		t.Fatalf("returned-record bound not disclosed: %+v", got)
	}
	more := make([]Evidence, 9900)
	for i := range more {
		more[i] = Evidence{Source: "server", Signal: "logs", Name: "bounded.error", Record: fmt.Sprint("bound-", i), Scope: "bounded-operation", Severity: 17, Time: now.Add(-time.Second).UnixNano(), Attrs: map[string]string{}}
	}
	if err := svc.Store.Insert(more); err != nil {
		t.Fatal(err)
	}
	response = errorRequest(handler, cookie, "?q=bounded")
	json.Unmarshal(response.Body.Bytes(), &got)
	if response.Code != 200 || got.Total != 10000 || len(got.Records) != 100 || !got.Truncated {
		t.Fatalf("scan bound silently treated counts as complete: total%d rows%d truncated%v", got.Total, len(got.Records), got.Truncated)
	}
	sessionRequest := httptest.NewRequest("GET", "http://localhost/api/session", nil)
	sessionRequest.AddCookie(cookie)
	sessionResponse := httptest.NewRecorder()
	handler.ServeHTTP(sessionResponse, sessionRequest)
	var session map[string]string
	json.Unmarshal(sessionResponse.Body.Bytes(), &session)
	request := httptest.NewRequest("POST", "http://localhost/api/delete", strings.NewReader(`{"scope":"bounded-operation"}`))
	request.Header.Set("Origin", svc.Origin)
	request.Header.Set("X-CSRF-Token", session["csrf"])
	request.AddCookie(cookie)
	deleted := httptest.NewRecorder()
	handler.ServeHTTP(deleted, request)
	if deleted.Code != 200 {
		t.Fatalf("deletion status%d: %s", deleted.Code, deleted.Body)
	}
	response = errorRequest(handler, cookie, "")
	json.Unmarshal(response.Body.Bytes(), &got)
	if response.Code != 200 || got.Total != 0 || len(got.Records) != 0 {
		t.Fatal("deleted diagnostics survived errors read")
	}
}

func TestDiagnosticExplicitSpanMessageWinsAndExceptionBodyEligibility(t *testing.T) {
	svc, _, _ := diagnosticFixture(t)
	now := time.Now().UTC()
	traces := ptrace.NewTraces()
	span := traces.ResourceSpans().AppendEmpty().ScopeSpans().AppendEmpty().Spans().AppendEmpty()
	span.SetName("explicit.message")
	span.SetStartTimestamp(pcommon.NewTimestampFromTime(now))
	span.SetTraceID(pcommon.TraceID{8})
	span.SetSpanID(pcommon.SpanID{9})
	span.Status().SetCode(ptrace.StatusCodeError)
	span.Status().SetMessage("fallback status message")
	span.Attributes().PutStr("exception.message", "explicit diagnostic message")
	if err := svc.consumeTraces(diagnosticContext(), traces); err != nil {
		t.Fatal(err)
	}
	logs := plog.NewLogs()
	entry := logs.ResourceLogs().AppendEmpty().ScopeLogs().AppendEmpty().LogRecords().AppendEmpty()
	entry.SetEventName("info.exception")
	entry.SetTimestamp(pcommon.NewTimestampFromTime(now))
	entry.Attributes().PutStr("exception.type", "SyntheticError")
	entry.Body().SetStr("exception details at info severity")
	if err := svc.consumeLogs(diagnosticContext(), logs); err != nil {
		t.Fatal(err)
	}
	got, err := svc.Store.Evidence()
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[0].Attrs["exception.message"] != "explicit diagnostic message" || got[1].Attrs["log.message"] != "exception details at info severity" {
		t.Fatalf("explicit diagnostic text precedence/body eligibility incorrect: %+v", got)
	}
}

func TestDiagnosticMessageOnlyExceptionIsDiscoverableButMetricsStayReceiptOnly(t *testing.T) {
	svc, handler, cookie := diagnosticFixture(t)
	now := time.Now().UTC().Add(-time.Second)
	logs := plog.NewLogs()
	entry := logs.ResourceLogs().AppendEmpty().ScopeLogs().AppendEmpty().LogRecords().AppendEmpty()
	entry.SetEventName("exception.message.only")
	entry.SetTimestamp(pcommon.NewTimestampFromTime(now))
	entry.Attributes().PutStr("exception.message", "message-only exception details")
	entry.Body().SetStr("INFO_BODY_STILL_OMITTED")
	if err := svc.consumeLogs(diagnosticContext(), logs); err != nil {
		t.Fatal(err)
	}
	metrics := pmetric.NewMetrics()
	metric := metrics.ResourceMetrics().AppendEmpty().ScopeMetrics().AppendEmpty().Metrics().AppendEmpty()
	metric.SetName("metric.with.exception.metadata")
	point := metric.SetEmptyGauge().DataPoints().AppendEmpty()
	point.SetTimestamp(pcommon.NewTimestampFromTime(now))
	point.SetIntValue(1)
	point.Attributes().PutStr("exception.type", "MetricMetadata")
	point.Attributes().PutStr("exception.message", "METRIC_TEXT_OMITTED")
	point.Attributes().PutStr("exception.stacktrace", "METRIC_STACK_OMITTED")
	if err := svc.consumeMetrics(diagnosticContext(), metrics); err != nil {
		t.Fatal(err)
	}
	response := errorRequest(handler, cookie, "")
	var got struct {
		Records []Evidence
		Total   int
	}
	json.Unmarshal(response.Body.Bytes(), &got)
	if response.Code != 200 || got.Total != 1 || len(got.Records) != 1 || got.Records[0].Name != "exception.message.only" {
		t.Fatalf("message-only exception absent or metric became error: %+v", got)
	}
	records, err := svc.Store.Evidence()
	if err != nil {
		t.Fatal(err)
	}
	raw, _ := json.Marshal(records)
	for _, sentinel := range []string{"INFO_BODY_STILL_OMITTED", "METRIC_TEXT_OMITTED", "METRIC_STACK_OMITTED"} {
		if strings.Contains(string(raw), sentinel) {
			t.Fatal("unexpected retained field: " + sentinel)
		}
	}
}
