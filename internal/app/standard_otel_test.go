package app

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog/plogotlp"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.opentelemetry.io/collector/pdata/ptrace/ptraceotlp"
)

// Readable shared wire source: times are offsets from a fresh synthetic epoch.
func standardBodies(t *testing.T) ([]byte, []byte, int64) {
	t.Helper()
	raw, err := os.ReadFile("../../rest/fixtures/standard-otel.json")
	if err != nil {
		t.Fatal(err)
	}
	var source map[string]json.RawMessage
	if err = json.Unmarshal(raw, &source); err != nil {
		t.Fatal(err)
	}
	tr := ptraceotlp.NewExportRequest()
	b, _ := json.Marshal(map[string]json.RawMessage{"resourceSpans": source["resourceSpans"]})
	if err = tr.UnmarshalJSON(b); err != nil {
		t.Fatal(err)
	}
	epoch := time.Now().Add(-time.Minute).UnixNano()
	spans := tr.Traces().ResourceSpans().At(0).ScopeSpans().At(0).Spans()
	for i := 0; i < spans.Len(); i++ {
		s := spans.At(i)
		s.SetStartTimestamp(pcommon.Timestamp(epoch + int64(s.StartTimestamp())))
		if s.EndTimestamp() != 0 {
			s.SetEndTimestamp(pcommon.Timestamp(epoch + int64(s.EndTimestamp())))
		}
		for j := 0; j < s.Events().Len(); j++ {
			e := s.Events().At(j)
			e.SetTimestamp(pcommon.Timestamp(epoch + int64(e.Timestamp())))
		}
	}
	logs := plogotlp.NewExportRequest()
	b, _ = json.Marshal(map[string]json.RawMessage{"resourceLogs": source["resourceLogs"]})
	if err = logs.UnmarshalJSON(b); err != nil {
		t.Fatal(err)
	}
	ls := logs.Logs().ResourceLogs().At(0).ScopeLogs().At(0).LogRecords()
	for i := 0; i < ls.Len(); i++ {
		l := ls.At(i)
		l.SetTimestamp(pcommon.Timestamp(epoch + int64(l.Timestamp())))
	}
	tb, err := tr.MarshalProto()
	if err != nil {
		t.Fatal(err)
	}
	lb, err := logs.MarshalProto()
	if err != nil {
		t.Fatal(err)
	}
	return tb, lb, epoch
}
func protectedEvidence(t *testing.T, c *http.Client, base string) ([]map[string]any, []byte) {
	t.Helper()
	r, err := c.Get(base + "/api/snapshot")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Body.Close()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatal(err)
	}
	if r.StatusCode != 200 {
		t.Fatalf("protected read %d", r.StatusCode)
	}
	var snap struct {
		Evidence []map[string]any `json:"evidence"`
		Model    struct {
			Journeys []any `json:"journeys"`
		} `json:"model"`
	}
	if err = json.Unmarshal(b, &snap); err != nil {
		t.Fatal(err)
	}
	if len(snap.Model.Journeys) != 0 {
		t.Fatal("ordinary OTel invented journey")
	}
	return snap.Evidence, b
}
func TestStandardOTelStructuralProjection(t *testing.T) {
	svc, server, c := fixtureService(t)
	tb, lb, epoch := standardBodies(t)
	for signal, body := range map[string][]byte{"traces": tb, "logs": lb} {
		code, _ := protocolRequest(t, c, server.URL, signal, "application/x-protobuf", "", "server-test", body)
		if code != 200 {
			t.Fatalf("ingest %s: %d", signal, code)
		}
	}
	loginTest(t, c, server.URL)
	es, b := protectedEvidence(t, c, server.URL)
	if len(es) != 9 {
		t.Fatalf("want 5 spans,2 events,2 logs; got %d", len(es))
	}
	if bytes.Contains(b, []byte("STRUCTURAL-PRIVATE-SENTINEL")) {
		t.Fatal("private structural payload leaked")
	}
	var root, child, event, log map[string]any
	for _, x := range es {
		if x["scope"] != "" || x["occurrence"] != "" {
			t.Fatal("custom context invented")
		}
		if x["projection"] != "structural-v2" {
			t.Errorf("missing projection policy: %v", x["projection"])
		}
		switch {
		case x["signal"] == "traces" && x["span"] == "0000000000000001":
			root = x
		case x["signal"] == "traces" && x["span"] == "0000000000000002":
			child = x
		case x["signal"] == "span-events" && x["name"] == "retry.observed":
			event = x
		case x["signal"] == "logs" && x["name"] == "vote.committed":
			log = x
		}
	}
	if root["end"] != strconv.FormatInt(epoch+1000000001, 10) || root["kind"] != float64(2) || root["status"] != float64(1) || root["flags"] != float64(1) {
		t.Errorf("native fields discarded: %v", root)
	}
	if child["parent"] != "0000000000000001" {
		t.Error("native parent missing")
	}
	if event["eventIndex"] != float64(1) {
		t.Error("exported event index missing")
	}
	if log["trace"] != "11111111111111111111111111111111" || log["span"] != "0000000000000001" || log["flags"] != float64(1) {
		t.Error("log context lost")
	}
	links, _ := root["links"].([]any)
	if len(links) != 32 || root["omittedLinks"] != float64(2) || root["droppedLinks"] != float64(4) || root["droppedEvents"] != float64(2) || root["droppedAttributes"] != float64(3) {
		t.Errorf("bounded links/drop disclosure missing: retained %d", len(links))
	}
	attrs := root["attrs"].(map[string]any)
	if attrs["http.route"] != "/votes/{voteId}" || attrs["http.request.method"] != "POST" || attrs["http.response.status_code"] != "200" {
		t.Error("safe HTTP selectors missing")
	}
	rows, err := svc.Store.DB.Query(`SELECT data FROM evidence`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var s string
		rows.Scan(&s)
		if bytes.Contains([]byte(s), []byte("STRUCTURAL-PRIVATE-SENTINEL")) {
			t.Fatal("private payload persisted")
		}
	}
}

func TestStandardOTelNativeReplayConflictsAreAtomic(t *testing.T) {
	svc, server, c := fixtureService(t)
	tb, _, _ := standardBodies(t)
	send := func(b []byte, want int) {
		t.Helper()
		code, _ := protocolRequest(t, c, server.URL, "traces", "application/x-protobuf", "", "server-test", b)
		if code != want {
			t.Fatalf("status %d want %d", code, want)
		}
	}
	send(tb, 200)
	send(tb, 200)
	before, _ := svc.Store.Evidence()
	if len(before) != 7 {
		t.Fatal("exact native replay duplicated spans/events")
	}
	for _, field := range []string{"end", "parent", "event", "event-removal", "link", "status"} {
		t.Run(field, func(t *testing.T) {
			q := ptraceotlp.NewExportRequest()
			if err := q.UnmarshalProto(tb); err != nil {
				t.Fatal(err)
			}
			sp := q.Traces().ResourceSpans().At(0).ScopeSpans().At(0).Spans()
			root := sp.At(0)
			// A new row before a conflict must roll back too.
			newSpan := sp.AppendEmpty()
			root.CopyTo(newSpan)
			newSpan.SetSpanID(pcommon.SpanID{9})
			switch field {
			case "end":
				root.SetEndTimestamp(root.EndTimestamp() + 1)
			case "parent":
				root.SetParentSpanID(pcommon.SpanID{8})
			case "event":
				root.Events().At(0).SetName("changed.event")
			case "event-removal":
				root.Events().RemoveIf(func(ptrace.SpanEvent) bool { return true })
			case "link":
				root.Links().At(0).SetSpanID(pcommon.SpanID{7})
			case "status":
				root.Status().SetCode(ptrace.StatusCodeError)
			}
			b, _ := q.MarshalProto()
			send(b, 400)
			after, _ := svc.Store.Evidence()
			if len(after) != len(before) {
				t.Fatal("conflicting direct request partially committed")
			}
		})
	}
}
func TestStandardOTelRejectsOldSchemaWithoutReset(t *testing.T) {
	path := filepath.Join(t.TempDir(), "unsupported.db")
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = db.Exec("CREATE TABLE evidence(data TEXT); INSERT INTO evidence VALUES('preserve-synthetic-sentinel')"); err != nil {
		t.Fatal(err)
	}
	db.Close()
	st, err := Open(path)
	if err == nil {
		st.DB.Close()
		t.Fatal("unsupported old schema silently opened")
	}
	if !strings.Contains(err.Error(), "schema") || !strings.Contains(err.Error(), "fresh") {
		t.Fatalf("no clear fresh-directory guidance: %v", err)
	}
	db, err = sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	var got string
	if err = db.QueryRow("SELECT data FROM evidence").Scan(&got); err != nil || got != "preserve-synthetic-sentinel" {
		t.Fatal("old data changed during startup")
	}
}
func TestStandardOTelHTTPProjectionPrivacy(t *testing.T) {
	_, server, c := fixtureService(t)
	tb, _, _ := standardBodies(t)
	q := ptraceotlp.NewExportRequest()
	q.UnmarshalProto(tb)
	sp := q.Traces().ResourceSpans().At(0).ScopeSpans().At(0).Spans()
	routes := []string{"https://example.test/votes?token=STRUCTURAL-PRIVATE-SENTINEL", "/votes/12345", "/votes/STRUCTURAL-PRIVATE-SENTINEL", "/votes?token=STRUCTURAL-PRIVATE-SENTINEL", "/votes/{voteId}"}
	for i, route := range routes {
		p := sp.At(i)
		p.Attributes().PutStr("http.route", route)
		p.Attributes().PutStr("url.full", "STRUCTURAL-PRIVATE-SENTINEL")
		p.Attributes().PutStr("authorization", "STRUCTURAL-PRIVATE-SENTINEL")
	}
	b, _ := q.MarshalProto()
	code, _ := protocolRequest(t, c, server.URL, "traces", "application/x-protobuf", "", "server-test", b)
	if code != 200 {
		t.Fatal(code)
	}
	loginTest(t, c, server.URL)
	es, raw := protectedEvidence(t, c, server.URL)
	if bytes.Contains(raw, []byte("STRUCTURAL-PRIVATE-SENTINEL")) {
		t.Fatal("unsafe route persisted")
	}
	for _, x := range es {
		if x["signal"] != "traces" {
			continue
		}
		a := x["attrs"].(map[string]any)
		if x["name"] != "invalid.clock" && a["http.route"] != nil {
			t.Fatal("unsafe route became selectable")
		}
	}
}

func TestStandardOTelRestartRetainsExactNativeProjection(t *testing.T) {
	svc, server, c := fixtureService(t)
	tb, _, _ := standardBodies(t)
	code, _ := protocolRequest(t, c, server.URL, "traces", "application/x-protobuf", "", "server-test", tb)
	if code != 200 {
		t.Fatal(code)
	}
	loginTest(t, c, server.URL)
	before, _ := protectedEvidence(t, c, server.URL)
	svc.Close()
	svc.Store.DB.Close()
	reopened, err := Open(svc.Store.Path)
	if err != nil {
		t.Fatal(err)
	}
	defer reopened.DB.Close()
	restarted, err := New(reopened, "synthetic owner passphrase", "browser-test", "server-test", "")
	if err != nil {
		t.Fatal(err)
	}
	defer restarted.Close()
	if err = restarted.StartReceiver(); err != nil {
		t.Fatal(err)
	}
	again := httptest.NewServer(restarted.Handler(""))
	defer again.Close()
	restarted.Origin = again.URL
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar, Timeout: 10 * time.Second}
	loginTest(t, client, again.URL)
	after, _ := protectedEvidence(t, client, again.URL)
	if !reflect.DeepEqual(before, after) {
		t.Fatal("native projection changed across SQLite reopen")
	}
	code, _ = protocolRequest(t, client, again.URL, "traces", "application/x-protobuf", "", "server-test", tb)
	if code != 200 {
		t.Fatal(code)
	}
	replay, _ := protectedEvidence(t, client, again.URL)
	if !reflect.DeepEqual(after, replay) {
		t.Fatal("exact bytes changed durable identity after restart")
	}
}

func TestStandardOTelRouteTemplateNamesAndOmissionDisclosure(t *testing.T) {
	_, server, c := fixtureService(t)
	tb, _, _ := standardBodies(t)
	q := ptraceotlp.NewExportRequest()
	q.UnmarshalProto(tb)
	sp := q.Traces().ResourceSpans().At(0).ScopeSpans().At(0).Spans()
	q.Traces().ResourceSpans().At(0).Resource().SetDroppedAttributesCount(5)
	q.Traces().ResourceSpans().At(0).ScopeSpans().At(0).Scope().SetDroppedAttributesCount(7)
	sp.At(0).SetName("POST /votes/{voteId}")
	sp.At(1).SetName("GET https://example.test/STRUCTURAL-PRIVATE-SENTINEL")
	sp.At(2).SetName("/votes/12345")
	body, _ := q.MarshalProto()
	code, _ := protocolRequest(t, c, server.URL, "traces", "application/x-protobuf", "", "server-test", body)
	if code != 200 {
		t.Fatal(code)
	}
	csrf := loginTest(t, c, server.URL)
	es, raw := protectedEvidence(t, c, server.URL)
	if bytes.Contains(raw, []byte("STRUCTURAL-PRIVATE-SENTINEL")) {
		t.Error("raw URL in span name survived privacy projection")
	}
	found := false
	for _, e := range es {
		if e["span"] == "0000000000000001" && e["signal"] == "traces" {
			if e["name"] != "POST /votes/{voteId}" || e["nameState"] != "retained" || e["resourceDroppedAttributes"] != float64(5) || e["scopeDroppedAttributes"] != float64(7) {
				t.Error("ordinary route-template operation became unnamed")
			}
			found = true
		}
		if e["span"] == "0000000000000002" || e["span"] == "0000000000000003" {
			if e["name"] != "unnamed" || e["nameState"] != "omitted" {
				t.Error("unsafe/dynamic name not explicitly omitted")
			}
		}
		if n, ok := e["omittedAttributes"].(float64); e["span"] == "0000000000000001" && e["signal"] == "traces" && (!ok || n < 1) {
			t.Error("projection omissions not distinct from SDK drop counts")
		}
	}
	if !found {
		t.Fatal("native operation missing")
	}
	d := ordinaryDefinition()
	d["mappings"].([]map[string]any)[0]["event"] = "POST /votes/{voteId}"
	postMeaning(t, c, server.URL, csrf, "/api/definitions", d, 200)
}
