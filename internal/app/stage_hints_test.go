package app

import (
	"fmt"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"testing"
	"time"

	"go.opentelemetry.io/collector/pdata/ptrace/ptraceotlp"
)

func TestStageHintsSuggestionsNeverActivateOrInherit(t *testing.T) {
	_, server, c := fixtureService(t)
	tb, _, _ := standardBodies(t)
	q := ptraceotlp.NewExportRequest()
	q.UnmarshalProto(tb)
	sp := q.Traces().ResourceSpans().At(0).ScopeSpans().At(0).Spans()
	root := sp.At(0)
	root.Attributes().PutStr("stage.name", "voting")
	root.Events().At(0).Attributes().PutStr("stage.name", "voting.persist")
	sp.At(1).Attributes().PutStr("stage.name", "")
	sp.At(2).Attributes().PutStr("stage.name", " ")
	sp.At(3).Attributes().PutInt("stage.name", 42)
	sp.At(4).Attributes().PutStr("stage.name", "unsafe?label")
	b, _ := q.MarshalProto()
	code, _ := protocolRequest(t, c, server.URL, "traces", "application/x-protobuf", "", "server-test", b)
	if code != 200 {
		t.Fatal(code)
	}
	loginTest(t, c, server.URL)
	snap := getTraceJSON(t, c, server.URL, "/api/snapshot", 200)
	contexts := snap["contextualDiscovery"].([]any)
	seen := false
	for _, v := range contexts {
		x := v.(map[string]any)
		hints, ok := x["stageHints"].([]any)
		if !ok {
			t.Fatal("bounded stage hints absent")
		}
		if x["name"] == "POST /votes" {
			if len(hints) != 1 || hints[0] != "voting" {
				t.Fatal("valid emitted hint missing")
			}
			seen = true
		}
		if x["name"] == "vote.persist" || x["name"] == "missing.parent" || x["name"] == "invalid.clock" || x["name"] == "retry.observed" {
			if len(hints) != 0 {
				t.Fatalf("invalid/inherited hint: %v", x)
			}
		}
	}
	if !seen {
		t.Fatal("root suggestion missing")
	}
	detail := getTraceJSON(t, c, server.URL, "/api/traces/11111111111111111111111111111111", 200)
	for _, v := range detail["business"].([]any) {
		x := v.(map[string]any)
		if x["stage"] == "voting" || x["stage"] == "voting.persist" {
			t.Fatal("hint activated without owner definition")
		}
	}
}
func TestStageHintsExplicitOwnerWinsAndHistoryIsLiteral(t *testing.T) {
	_, server, c := fixtureService(t)
	csrf := loginTest(t, c, server.URL)
	d := ordinaryDefinition()
	activation := postMeaning(t, c, server.URL, csrf, "/api/definitions", d, 200)
	version := activation["version"]
	tb, _, _ := standardBodies(t)
	q := ptraceotlp.NewExportRequest()
	q.UnmarshalProto(tb)
	root := q.Traces().ResourceSpans().At(0).ScopeSpans().At(0).Spans().At(0)
	root.Events().At(0).Attributes().PutStr("stage.name", "voting")
	b, _ := q.MarshalProto()
	code, _ := protocolRequest(t, c, server.URL, "traces", "application/x-protobuf", "", "server-test", b)
	if code != 200 {
		t.Fatal(code)
	}
	path := "/api/traces/11111111111111111111111111111111"
	check := func() {
		t.Helper()
		detail := getTraceJSON(t, c, server.URL, path, 200)
		found := false
		for _, v := range detail["business"].([]any) {
			m := v.(map[string]any)
			if m["goal"] == true {
				if m["stage"] != "Vote saved" || m["definition"] != version {
					t.Fatal("emitted tag or new definition renamed capture")
				}
				found = true
			}
		}
		if !found {
			t.Fatal("owner goal missing")
		}
	}
	check()
	d = ordinaryDefinition()
	d["mappings"].([]map[string]any)[1]["step"] = "voting"
	postMeaning(t, c, server.URL, csrf, "/api/definitions", d, 200)
	check()
	// No attribute-implied authority: browser export has a separate namespace and
	// cannot match the server selector, even when parent/tag/IDs are forged.
	code, _ = protocolRequest(t, c, server.URL, "traces", "application/x-protobuf", "", "browser-test", b)
	if code != 200 {
		t.Fatal(code)
	}
	detail := getTraceJSON(t, c, server.URL, path, 200)
	for _, v := range detail["business"].([]any) {
		m := v.(map[string]any)
		if m["stage"] == "voting" {
			t.Fatal("browser hint promoted to trusted server mapping")
		}
	}
	// The stored event retains emitted hint for mismatch inspection, independently
	// of the literal owner label and definition.
	for _, v := range detail["events"].([]any) {
		e := v.(map[string]any)
		if e["source"] == "server" && e["name"] == "vote.committed" {
			a := e["attrs"].(map[string]any)
			if a["stage.name"] != "voting" {
				t.Fatal("emitted evidence lost")
			}
		}
	}
}
func TestStageHintsDistinctLabelsRemainSeparateSuggestions(t *testing.T) {
	_, server, c := fixtureService(t)
	tb, _, _ := standardBodies(t)
	q := ptraceotlp.NewExportRequest()
	q.UnmarshalProto(tb)
	sp := q.Traces().ResourceSpans().At(0).ScopeSpans().At(0).Spans()
	sp.At(1).Attributes().PutStr("stage.name", "voting")
	sp.At(2).Attributes().PutStr("stage.name", "voting.persist")
	b, _ := q.MarshalProto()
	code, _ := protocolRequest(t, c, server.URL, "traces", "application/x-protobuf", "", "server-test", b)
	if code != 200 {
		t.Fatal(code)
	}
	loginTest(t, c, server.URL)
	snap := getTraceJSON(t, c, server.URL, "/api/snapshot", 200)
	for _, v := range snap["contextualDiscovery"].([]any) {
		x := v.(map[string]any)
		if x["name"] == "vote.persist" {
			labels := x["stageHints"].([]any)
			if len(labels) != 2 {
				t.Fatal("conflicting labels collapsed or popularity-selected")
			}
		}
	}

}

func TestStageHintsBoundedContextualLabels(t *testing.T) {
	_, server, c := fixtureService(t)
	logs := plog.NewLogs()
	r := logs.ResourceLogs().AppendEmpty()
	r.Resource().Attributes().PutStr("service.name", "synthetic-voting")
	ls := r.ScopeLogs().AppendEmpty().LogRecords()
	for i := 0; i < 34; i++ {
		l := ls.AppendEmpty()
		l.SetEventName("hint.observed")
		l.SetTimestamp(pcommon.NewTimestampFromTime(time.Now()))
		l.Attributes().PutStr("stage.name", fmt.Sprintf("voting.%02d", i))
	}
	if code := requestLogs(t, c, server.URL, "server-test", logs); code != 200 {
		t.Fatal(code)
	}
	loginTest(t, c, server.URL)
	snap := getTraceJSON(t, c, server.URL, "/api/snapshot", 200)
	contexts := snap["contextualDiscovery"].([]any)
	if len(contexts) != 1 {
		t.Fatal("context grouping changed")
	}
	context := contexts[0].(map[string]any)
	if len(context["stageHints"].([]any)) != 32 || context["omittedStageHints"] != float64(2) || context["count"] != float64(34) {
		t.Fatal("suggestions were unbounded, dropped receipt, or lost omissions")
	}
}
