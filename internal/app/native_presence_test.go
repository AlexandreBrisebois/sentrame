package app

import (
	"fmt"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.opentelemetry.io/collector/pdata/ptrace/ptraceotlp"
	"testing"
	"time"
)

func TestNativeBrowserIntervalUsesOnlyBrowserPresence(t *testing.T) {
	_, server, c := fixtureService(t)
	csrf := loginTest(t, c, server.URL)
	d := Definition{Goal: "finished", Mappings: []Mapping{{Event: "voting.page", Step: "Voting page", Source: "browser", Service: "synthetic-presence", Carrier: "span", Kind: "step", Timing: "native-span"}, {Event: "server.operation", Step: "Server operation", Source: "server", Service: "synthetic-presence", Carrier: "span", Kind: "step", Timing: "native-span"}, {Event: "activity.changed", Step: "Presence", Source: "browser", Carrier: "log", Kind: "activity"}, {Event: "finished", Step: "Completion", Source: "browser", Carrier: "log", Kind: "step"}}}
	postMeaning(t, c, server.URL, csrf, "/api/definitions", d, 200)
	base := time.Now().Add(-time.Minute)
	logs := plog.NewLogs()
	r := logs.ResourceLogs().AppendEmpty()
	r.Resource().Attributes().PutStr("service.name", "synthetic-presence")
	ls := r.ScopeLogs().AppendEmpty().LogRecords()
	for i, observation := range []struct {
		offset time.Duration
		state  string
	}{{100 * time.Millisecond, "active"}, {400 * time.Millisecond, "hidden"}, {600 * time.Millisecond, "active"}, {900 * time.Millisecond, "idle"}} {
		l := ls.AppendEmpty()
		l.SetEventName("activity.changed")
		l.SetTimestamp(pcommon.NewTimestampFromTime(base.Add(observation.offset)))
		l.Attributes().PutStr("session.id", "native-presence-scope")
		l.Attributes().PutStr("tab.id", "native-presence-tab")
		l.Attributes().PutStr("event.id", fmt.Sprint(i))
		l.Attributes().PutStr("activity.state", observation.state)
	}
	if code := requestLogs(t, c, server.URL, "browser-test", logs); code != 200 {
		t.Fatal(code)
	}
	// The real SDK-style page span arrives only when ended, after the states.
	q := ptraceotlp.NewExportRequest()
	sr := q.Traces().ResourceSpans().AppendEmpty()
	sr.Resource().Attributes().PutStr("service.name", "synthetic-presence")
	p := sr.ScopeSpans().AppendEmpty().Spans().AppendEmpty()
	p.SetName("voting.page")
	p.SetTraceID(pcommon.TraceID{7})
	p.SetSpanID(pcommon.SpanID{8})
	p.SetStartTimestamp(pcommon.NewTimestampFromTime(base))
	p.SetEndTimestamp(pcommon.NewTimestampFromTime(base.Add(time.Second)))
	p.Attributes().PutStr("session.id", "native-presence-scope")
	p.Attributes().PutStr("tab.id", "native-presence-tab")
	p.Status().SetCode(ptrace.StatusCodeOk)
	body, _ := q.MarshalProto()
	code, _ := protocolRequest(t, c, server.URL, "traces", "application/x-protobuf", "", "browser-test", body)
	if code != 200 {
		t.Fatal(code)
	}
	p.SetName("server.operation")
	body, _ = q.MarshalProto()
	code, _ = protocolRequest(t, c, server.URL, "traces", "application/x-protobuf", "", "server-test", body)
	if code != 200 {
		t.Fatal(code)
	}
	goal := plog.NewLogs()
	gr := goal.ResourceLogs().AppendEmpty()
	gr.Resource().Attributes().PutStr("service.name", "synthetic-presence")
	g := gr.ScopeLogs().AppendEmpty().LogRecords().AppendEmpty()
	g.SetEventName("finished")
	g.SetTimestamp(pcommon.NewTimestampFromTime(base.Add(1100 * time.Millisecond)))
	g.Attributes().PutStr("session.id", "native-presence-scope")
	g.Attributes().PutStr("tab.id", "native-presence-tab")
	if code := requestLogs(t, c, server.URL, "browser-test", goal); code != 200 {
		t.Fatal(code)
	}
	model := readNativeModel(t, c, server.URL)
	if len(model.Journeys) != 1 {
		t.Fatal("native scope split")
	}
	found := false
	for _, timing := range model.Journeys[0].Timings {
		if timing.Stage == "Voting page" && timing.Basis == "observed-active" {
			if !timing.Complete || timing.Milliseconds != 600 || timing.DurationNs == nil || *timing.DurationNs != "600000000" {
				t.Fatalf("native browser active measurement %v", timing)
			}
			found = true
		}
		if timing.Stage == "Server operation" && timing.Basis == "observed-active" {
			t.Fatal("server borrowed browser presence")
		}
	}
	if !found {
		t.Fatal("native browser interval has no observed-active basis despite scoped actual wire states")
	}
}
