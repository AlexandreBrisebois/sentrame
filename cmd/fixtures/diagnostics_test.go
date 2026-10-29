package main

import (
	"bytes"
	"encoding/json"
	logsv1 "go.opentelemetry.io/proto/otlp/collector/logs/v1"
	tracesv1 "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	tracev1 "go.opentelemetry.io/proto/otlp/trace/v1"
	"google.golang.org/protobuf/proto"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func diagnosticSpec(t *testing.T) DiagnosticSpec {
	t.Helper()
	raw, err := os.ReadFile("../../rest/fixtures/diagnostics.json")
	if err != nil {
		t.Fatal(err)
	}
	var spec DiagnosticSpec
	if err = json.Unmarshal(raw, &spec); err != nil {
		t.Fatal(err)
	}
	return spec
}
func TestDiagnosticFixturesUseStandardWireIdentityStacksAndUnscopedBody(t *testing.T) {
	spec := diagnosticSpec(t)
	now := time.Date(2026, 10, 28, 20, 0, 0, 0, time.UTC)
	a, err := buildDiagnostics(spec, now)
	if err != nil {
		t.Fatal(err)
	}
	b, err := buildDiagnostics(spec, now)
	if err != nil {
		t.Fatal(err)
	}
	if a.Tag == b.Tag || bytes.Equal(a.Traces, b.Traces) || bytes.Equal(a.Logs, b.Logs) {
		t.Fatal("new generation must have fresh identities and safe tag")
	}
	var traces tracesv1.ExportTraceServiceRequest
	if err = proto.Unmarshal(a.Traces, &traces); err != nil {
		t.Fatal(err)
	}
	span := traces.ResourceSpans[0].ScopeSpans[0].Spans[0]
	if len(span.TraceId) != 16 || len(span.SpanId) != 8 || span.Status.Code != tracev1.Status_STATUS_CODE_ERROR || len(span.Events) != 1 || !strings.Contains(span.Name, a.Tag) {
		t.Fatal("expected standard identified Error span with exception event")
	}
	attrs := span.Events[0].Attributes
	foundStack := false
	for _, kv := range attrs {
		if kv.Key == "exception.stacktrace" && strings.Contains(kv.Value.GetStringValue(), "OrderService.CommitAsync") {
			foundStack = true
		}
	}
	if !foundStack {
		t.Fatal(".NET-like standard exception stack missing")
	}
	var logs logsv1.ExportLogsServiceRequest
	if err = proto.Unmarshal(a.Logs, &logs); err != nil {
		t.Fatal(err)
	}
	records := logs.ResourceLogs[0].ScopeLogs[0].LogRecords
	if len(records) != 2 || records[1].Body.GetStringValue() != spec.CustomBody {
		t.Fatal("custom error must use ordinary standard log body")
	}
	for _, kv := range records[1].Attributes {
		if kv.Key == "session.id" {
			t.Fatal("custom error must demonstrate appwide discovery without session")
		}
	}
	for _, r := range records {
		if r.SeverityNumber != 17 || r.TimeUnixNano >= uint64(now.UnixNano()) || !strings.Contains(r.EventName, a.Tag) {
			t.Fatal("fresh error records need bounded past occurrence time and searchable tag")
		}
	}
}
func TestDiagnosticGenerationPreservesReplayStateAndLiteralRESTBodies(t *testing.T) {
	specRaw, err := os.ReadFile("../../rest/fixtures/diagnostics.json")
	if err != nil {
		t.Fatal(err)
	}
	httpRaw, err := os.ReadFile("../../rest/diagnostics.http")
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "rest", "fixtures"), 0700)
	os.MkdirAll(filepath.Join(dir, "rest", "generated"), 0700)
	os.WriteFile(filepath.Join(dir, "rest", "fixtures", "diagnostics.json"), specRaw, 0600)
	os.WriteFile(filepath.Join(dir, "rest", "diagnostics.http"), httpRaw, 0600)
	os.WriteFile(filepath.Join(dir, "rest", "generated", "manifest.json"), []byte("preserved"), 0600)
	os.WriteFile(filepath.Join(dir, "rest", ".env"), []byte("workspaceScope=preserved\n"), 0600)
	t.Chdir(dir)
	t.Setenv("SENTRAME_PASSPHRASE", "synthetic test passphrase")
	t.Setenv("SENTRAME_SERVER_TOKEN", "synthetic-server-token")
	if err = generateDiagnostics(); err != nil {
		t.Fatal(err)
	}
	prior, _ := os.ReadFile("rest/generated/manifest.json")
	env, _ := os.ReadFile("rest/.env")
	if string(prior) != "preserved" || !bytes.Contains(env, []byte("workspaceScope=preserved")) {
		t.Fatal("diagnostics generation changed unrelated replay state")
	}
	var manifest DiagnosticGenerated
	raw, _ := os.ReadFile("rest/generated/diagnostics-manifest.json")
	json.Unmarshal(raw, &manifest)
	reqs, err := parseREST("rest/diagnostics.http", map[string]string{"baseUrl": baseURL(), "serverToken": "synthetic-server-token", "diagnosticsTag": manifest.Tag})
	if err != nil {
		t.Fatal(err)
	}
	if len(reqs) != 7 || !bytes.Equal(reqs[1].Body, reqs[4].Body) || !bytes.Equal(reqs[2].Body, reqs[5].Body) || !strings.Contains(reqs[3].URL, "/api/errors?q="+manifest.Tag) {
		t.Fatal("literal requests must use exact prepared protobuf and appwide error discovery")
	}
}
func TestDiagnosticVerifierRejectsMissingStacksAndRawSentinels(t *testing.T) {
	tag := "diagnostics-safe-tag"
	g := DiagnosticGenerated{Tag: tag, Sentinels: []string{"private-sentinel"}}
	records := []map[string]any{{"id": 1, "name": tag + ".dotnet.request", "severity": 17, "attrs": map[string]string{"exception.message": "request failed", "diagnostic.policy": "scrubbed-v1"}}, {"id": 2, "name": tag + ".dotnet.exception", "attrs": map[string]string{"exception.type": "System.InvalidOperationException", "exception.message": "Cannot persist synthetic checkout", "exception.stacktrace": "at OrderService.CommitAsync()", "diagnostic.policy": "scrubbed-v1"}}, {"id": 3, "name": tag + ".next.render", "severity": 17, "attrs": map[string]string{"exception.type": "TypeError", "exception.message": "Cannot read properties of undefined", "exception.stacktrace": "at CheckoutPage()", "diagnostic.policy": "scrubbed-v1"}}, {"id": 4, "name": tag + ".custom.error", "severity": 17, "scope": "", "attrs": map[string]string{"log.message": "Synthetic custom error: pending state was preserved", "diagnostic.policy": "scrubbed-v1"}}}
	response := map[string]any{"records": records, "total": 4, "unit": "diagnostic records", "policy": "scrubbed-v1"}
	raw, _ := json.Marshal(response)
	if _, err := checkDiagnostics(raw, g); err != nil {
		t.Fatal(err)
	}
	records[1]["attrs"].(map[string]string)["exception.stacktrace"] = ""
	raw, _ = json.Marshal(response)
	if _, err := checkDiagnostics(raw, g); err == nil {
		t.Fatal("missing stack accepted")
	}
	records[1]["attrs"].(map[string]string)["exception.stacktrace"] = "private-sentinel"
	raw, _ = json.Marshal(response)
	if _, err := checkDiagnostics(raw, g); err == nil {
		t.Fatal("raw private content accepted")
	}
}
