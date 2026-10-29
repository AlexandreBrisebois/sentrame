package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	logsv1 "go.opentelemetry.io/proto/otlp/collector/logs/v1"
	tracesv1 "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	commonv1 "go.opentelemetry.io/proto/otlp/common/v1"
	logdata "go.opentelemetry.io/proto/otlp/logs/v1"
	resourcev1 "go.opentelemetry.io/proto/otlp/resource/v1"
	tracedata "go.opentelemetry.io/proto/otlp/trace/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"time"
	"unicode/utf8"
)

type DiagnosticException struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Stack   string `json:"stack"`
}
type DiagnosticSpec struct {
	DotNet     DiagnosticException `json:"dotnet"`
	Next       DiagnosticException `json:"next"`
	CustomBody string              `json:"customBody"`
	Sentinels  []string            `json:"sentinels"`
}
type DiagnosticGenerated struct {
	Tag         string            `json:"tag"`
	GeneratedAt string            `json:"generatedAt"`
	Sentinels   []string          `json:"sentinels"`
	Hashes      map[string]string `json:"hashes"`
	Traces      []byte            `json:"-"`
	Logs        []byte            `json:"-"`
	TracesJSON  []byte            `json:"-"`
	LogsJSON    []byte            `json:"-"`
}

func diagnosticBytesID(size int) []byte { b, _ := hex.DecodeString(randomID()); return b[:size] }
func buildDiagnostics(spec DiagnosticSpec, now time.Time) (DiagnosticGenerated, error) {
	g := DiagnosticGenerated{Tag: "diagnostics-" + randomID()[:12], GeneratedAt: now.UTC().Format(time.RFC3339Nano), Sentinels: append([]string(nil), spec.Sentinels...), Hashes: map[string]string{}}
	if spec.DotNet.Type == "" || spec.DotNet.Stack == "" || spec.Next.Type == "" || spec.Next.Stack == "" || spec.CustomBody == "" {
		return g, errors.New("diagnostic fixture needs both typed stacks and custom body")
	}
	exceptionAttrs := func(e DiagnosticException) []*commonv1.KeyValue {
		return []*commonv1.KeyValue{kv("exception.type", e.Type), kv("exception.message", e.Message), kv("exception.stacktrace", e.Stack)}
	}
	trace := &tracesv1.ExportTraceServiceRequest{ResourceSpans: []*tracedata.ResourceSpans{{Resource: &resourcev1.Resource{Attributes: []*commonv1.KeyValue{kv("service.name", "synthetic-dotnet-wire"), kv("service.version", "synthetic-a1"), kv("process.runtime.name", ".NET-like synthetic wire example")}}, ScopeSpans: []*tracedata.ScopeSpans{{Scope: &commonv1.InstrumentationScope{Name: "sentrame-synthetic-diagnostic-wire", Version: "1"}, Spans: []*tracedata.Span{{TraceId: diagnosticBytesID(16), SpanId: diagnosticBytesID(8), Name: g.Tag + ".dotnet.request", Kind: tracedata.Span_SPAN_KIND_SERVER, StartTimeUnixNano: uint64(now.Add(-4 * time.Second).UnixNano()), EndTimeUnixNano: uint64(now.Add(-2 * time.Second).UnixNano()), Attributes: []*commonv1.KeyValue{kv("session.id", randomID())}, Status: &tracedata.Status{Code: tracedata.Status_STATUS_CODE_ERROR, Message: spec.DotNet.Message}, Events: []*tracedata.Span_Event{{Name: g.Tag + ".dotnet.exception", TimeUnixNano: uint64(now.Add(-3 * time.Second).UnixNano()), Attributes: exceptionAttrs(spec.DotNet)}}}}}}}}}
	nextAttrs := append(exceptionAttrs(spec.Next), kv("event.id", randomID()), kv("session.id", randomID()))
	logs := &logsv1.ExportLogsServiceRequest{ResourceLogs: []*logdata.ResourceLogs{{Resource: &resourcev1.Resource{Attributes: []*commonv1.KeyValue{kv("service.name", "synthetic-next-wire"), kv("service.version", "synthetic-a1"), kv("process.runtime.name", "Node.js-like synthetic wire example")}}, ScopeLogs: []*logdata.ScopeLogs{{Scope: &commonv1.InstrumentationScope{Name: "sentrame-synthetic-diagnostic-wire", Version: "1"}, LogRecords: []*logdata.LogRecord{{TimeUnixNano: uint64(now.Add(-time.Second).UnixNano()), ObservedTimeUnixNano: uint64(now.UnixNano()), EventName: g.Tag + ".next.render", SeverityNumber: logdata.SeverityNumber_SEVERITY_NUMBER_ERROR, SeverityText: "ERROR", Attributes: nextAttrs}, {TimeUnixNano: uint64(now.Add(-500 * time.Millisecond).UnixNano()), ObservedTimeUnixNano: uint64(now.UnixNano()), EventName: g.Tag + ".custom.error", SeverityNumber: logdata.SeverityNumber_SEVERITY_NUMBER_ERROR, SeverityText: "ERROR", Attributes: []*commonv1.KeyValue{kv("event.id", randomID())}, Body: &commonv1.AnyValue{Value: &commonv1.AnyValue_StringValue{StringValue: spec.CustomBody}}}}}}}}}
	var err error
	if g.Traces, err = proto.Marshal(trace); err != nil {
		return g, err
	}
	if g.Logs, err = proto.Marshal(logs); err != nil {
		return g, err
	}
	if g.TracesJSON, err = (protojson.MarshalOptions{Indent: "  "}).Marshal(trace); err != nil {
		return g, err
	}
	if g.LogsJSON, err = (protojson.MarshalOptions{Indent: "  "}).Marshal(logs); err != nil {
		return g, err
	}
	for name, body := range map[string][]byte{"diagnostics-traces": g.Traces, "diagnostics-logs": g.Logs} {
		h := sha256.Sum256(body)
		g.Hashes[name] = hex.EncodeToString(h[:])
	}
	return g, nil
}
func generateDiagnostics() error {
	raw, err := os.ReadFile("rest/fixtures/diagnostics.json")
	if err != nil {
		return err
	}
	var spec DiagnosticSpec
	if err = json.Unmarshal(raw, &spec); err != nil {
		return err
	}
	g, err := buildDiagnostics(spec, time.Now())
	if err != nil {
		return err
	}
	pass, err := envRequired("SENTRAME_PASSPHRASE")
	if err != nil {
		return err
	}
	token, err := envRequired("SENTRAME_SERVER_TOKEN")
	if err != nil {
		return err
	}
	login, _ := json.Marshal(map[string]string{"passphrase": pass})
	manifest, _ := json.MarshalIndent(g, "", "  ")
	if err = os.MkdirAll("rest/generated", 0700); err != nil {
		return err
	}
	files := map[string][]byte{"diagnostics-traces.bin": g.Traces, "diagnostics-logs.bin": g.Logs, "diagnostics-traces.json": g.TracesJSON, "diagnostics-logs.json": g.LogsJSON, "diagnostics-login.json": login, "diagnostics-manifest.json": manifest}
	for name, body := range files {
		if err = os.WriteFile(filepath.Join("rest/generated", name), body, 0600); err != nil {
			return err
		}
	}
	values := map[string]string{"baseUrl": baseURL(), "serverToken": token, "diagnosticsTag": g.Tag}
	existing, err := os.ReadFile("rest/.env")
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	lines := []string{}
	for _, line := range strings.Split(strings.TrimSpace(string(existing)), "\n") {
		key, _, _ := strings.Cut(line, "=")
		if _, replace := values[key]; !replace && line != "" {
			lines = append(lines, line)
		}
	}
	keys := []string{}
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		lines = append(lines, key+"="+values[key])
	}
	if err = os.WriteFile("rest/.env", []byte(strings.Join(lines, "\n")+"\n"), 0600); err != nil {
		return err
	}
	fmt.Printf("Generated synthetic diagnostic wire fixtures with tag%s; existing replay manifests preserved. These are not real .NET/Next.js exporter qualifications.\n", g.Tag)
	return nil
}
func checkDiagnostics(raw []byte, g DiagnosticGenerated) ([]int64, error) {
	for _, sentinel := range g.Sentinels {
		if sentinel != "" && bytes.Contains(raw, []byte(sentinel)) {
			return nil, errors.New("diagnostic privacy sentinel leaked through protected read")
		}
	}
	var response struct {
		Records []struct {
			ID    int64             `json:"id"`
			Name  string            `json:"name"`
			Scope string            `json:"scope"`
			Attrs map[string]string `json:"attrs"`
		} `json:"records"`
		Total     int    `json:"total"`
		Truncated bool   `json:"truncated"`
		Unit      string `json:"unit"`
		Policy    string `json:"policy"`
	}
	if err := json.Unmarshal(raw, &response); err != nil {
		return nil, err
	}
	if response.Total != 4 || len(response.Records) != 4 || response.Truncated || response.Unit != "diagnostic records" {
		return nil, fmt.Errorf("diagnostic query needs four untruncated records, got%d/%d", response.Total, len(response.Records))
	}
	seen := map[string]bool{}
	seenIDs := map[int64]bool{}
	ids := []int64{}
	for _, r := range response.Records {
		if r.ID == 0 || seenIDs[r.ID] || seen[r.Name] || !strings.HasPrefix(r.Name, g.Tag+".") {
			return nil, errors.New("diagnostic record identity/name mismatch")
		}
		seen[r.Name] = true
		seenIDs[r.ID] = true
		ids = append(ids, r.ID)
		if r.Attrs["diagnostic.policy"] != "scrubbed-v1" {
			return nil, errors.New("diagnostic projection policy marker missing")
		}
		for _, field := range []string{"exception.message", "exception.stacktrace", "log.message"} {
			limit := 2048
			if field == "exception.stacktrace" {
				limit = 8192
			}
			value := r.Attrs[field]
			if len(value) > limit || !utf8.ValidString(value) || field == "exception.stacktrace" && strings.Count(value, "\n") >= 40 {
				return nil, errors.New("diagnostic field exceeded safe UTF8 byte bound")
			}
		}
		switch strings.TrimPrefix(r.Name, g.Tag+".") {
		case "dotnet.request":
			if r.Attrs["exception.message"] == "" {
				return nil, errors.New("Error-status span description missing")
			}
		case "dotnet.exception":
			if r.Attrs["exception.type"] != "System.InvalidOperationException" || !strings.Contains(r.Attrs["exception.message"], "Cannot persist synthetic checkout") || !strings.Contains(r.Attrs["exception.stacktrace"], "OrderService.CommitAsync") {
				return nil, errors.New(".NET-like message/stack missing")
			}
		case "next.render":
			if r.Attrs["exception.type"] != "TypeError" || !strings.Contains(r.Attrs["exception.message"], "Cannot read properties") || !strings.Contains(r.Attrs["exception.stacktrace"], "CheckoutPage") {
				return nil, errors.New("Next-like message/stack missing")
			}
		case "custom.error":
			if r.Scope != "" || !strings.Contains(r.Attrs["log.message"], "Synthetic custom error") || !strings.Contains(r.Attrs["log.message"], "pending state was preserved") {
				return nil, errors.New("unscoped standard custom log body missing")
			}
		default:
			return nil, errors.New("unexpected diagnostic fixture record")
		}
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	return ids, nil
}
func verifyDiagnostics() error {
	raw, err := os.ReadFile("rest/generated/diagnostics-manifest.json")
	if err != nil {
		return err
	}
	var g DiagnosticGenerated
	if err = json.Unmarshal(raw, &g); err != nil {
		return err
	}
	token, err := envRequired("SENTRAME_SERVER_TOKEN")
	if err != nil {
		return err
	}
	for name, expected := range g.Hashes {
		body, err := os.ReadFile(filepath.Join("rest/generated", name+".bin"))
		if err != nil {
			return err
		}
		h := sha256.Sum256(body)
		if hex.EncodeToString(h[:]) != expected {
			return errors.New("prepared diagnostic body changed; regenerate rather than changing replay identity")
		}
	}
	requests, err := parseREST("rest/diagnostics.http", map[string]string{"baseUrl": baseURL(), "serverToken": token, "diagnosticsTag": g.Tag})
	if err != nil {
		return err
	}
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar, Timeout: 20 * time.Second}
	var first []int64
	reads := 0
	for i, r := range requests {
		req, err := http.NewRequest(r.Method, r.URL, bytes.NewReader(r.Body))
		if err != nil {
			return err
		}
		req.Header = r.Headers.Clone()
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		body, readErr := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
		resp.Body.Close()
		if readErr != nil {
			return readErr
		}
		if resp.StatusCode != r.Expected {
			return fmt.Errorf("diagnostic request%d %s %s status%d,expected%d", i+1, r.Method, req.URL.Path, resp.StatusCode, r.Expected)
		}
		if r.Method == "GET" && req.URL.Path == "/api/errors" {
			ids, err := checkDiagnostics(body, g)
			if err != nil {
				return err
			}
			if reads == 0 {
				first = ids
			} else if !reflect.DeepEqual(first, ids) {
				return errors.New("diagnostic records changed after exact replay")
			}
			reads++
		}
		fmt.Printf("PASS diagnostics %02d %s %s → %d\n", i+1, r.Method, req.URL.Path, resp.StatusCode)
	}
	if reads != 2 {
		return errors.New("literal diagnostic requests omitted replay comparison")
	}
	fmt.Println("PASS four diagnostic records, preserved scrubbed .NET-like/Next-like messages and stacks, standard custom body without visit, bounded fields, privacy sentinels absent, stable replay identities. Verified actual OTLP→storage→protected error reads. Wire examples do not qualify .NET/Next.js exporters.")
	return nil
}
