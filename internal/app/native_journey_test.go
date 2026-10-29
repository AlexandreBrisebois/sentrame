package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.opentelemetry.io/collector/pdata/ptrace/ptraceotlp"
)

func scopedNative(t *testing.T) []byte {
	t.Helper()
	b, _, _ := standardBodies(t)
	q := ptraceotlp.NewExportRequest()
	q.UnmarshalProto(b)
	sp := q.Traces().ResourceSpans().At(0).ScopeSpans().At(0).Spans()
	for i := 0; i < sp.Len(); i++ {
		sp.At(i).Attributes().PutStr("session.id", "synthetic-native-scope")
	}
	b, _ = q.MarshalProto()
	return b
}
func readNativeModel(t *testing.T, c *http.Client, base string) Model {
	t.Helper()
	r, err := c.Get(base + "/api/snapshot")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Body.Close()
	var snap struct {
		Model Model `json:"model"`
	}
	if err = json.NewDecoder(r.Body).Decode(&snap); err != nil {
		t.Fatal(err)
	}
	return snap.Model
}
func TestNativeJourneyInclusiveTimingAndGoalEnd(t *testing.T) {
	_, server, c := fixtureService(t)
	csrf := loginTest(t, c, server.URL)
	d := ordinaryDefinition()
	d["goal"] = "vote.persist"
	d["fallbackTiming"] = true
	postMeaning(t, c, server.URL, csrf, "/api/definitions", d, 200)
	b := scopedNative(t)
	q := ptraceotlp.NewExportRequest()
	q.UnmarshalProto(b)
	sp := q.Traces().ResourceSpans().At(0).ScopeSpans().At(0).Spans()
	goalEnd := int64(sp.At(1).EndTimestamp())
	code, _ := protocolRequest(t, c, server.URL, "traces", "application/x-protobuf", "", "server-test", b)
	if code != 200 {
		t.Fatal(code)
	}
	model := readNativeModel(t, c, server.URL)
	if len(model.Journeys) != 2 || model.Journeys[0].End != goalEnd {
		t.Fatalf("native goal must close at end after parallel starts; got %d journeys", len(model.Journeys))
	}
	j := model.Journeys[0]
	native := 0
	for _, timing := range j.Timings {
		if timing.Basis == "native-inclusive" {
			native++
			if timing.Occurrence == "" || !timing.Complete {
				t.Fatal("native identity or completed duration missing")
			}
		}
		if timing.Basis == "next-event-estimated" {
			t.Fatal("native span double-pooled as next-event estimate")
		}
	}
	if native != 3 {
		t.Fatalf("native samples %d want3", native)
	}
	if j.Tab != "" {
		t.Fatal("server-only scope invented browser tab")
	}
}
func TestNativeJourneyLateEarlierParentPreservesDefinition(t *testing.T) {
	_, server, c := fixtureService(t)
	csrf := loginTest(t, c, server.URL)
	d := ordinaryDefinition()
	postMeaning(t, c, server.URL, csrf, "/api/definitions", d, 200)
	b := scopedNative(t)
	q := ptraceotlp.NewExportRequest()
	q.UnmarshalProto(b)
	sp := q.Traces().ResourceSpans().At(0).ScopeSpans().At(0).Spans()
	late := ptraceotlp.NewExportRequest()
	r := late.Traces().ResourceSpans().AppendEmpty()
	q.Traces().ResourceSpans().At(0).Resource().CopyTo(r.Resource())
	sp.At(0).CopyTo(r.ScopeSpans().AppendEmpty().Spans().AppendEmpty())
	sp.RemoveIf(func(p ptrace.Span) bool { return p.Name() == "POST /votes" })
	first, _ := q.MarshalProto()
	code, _ := protocolRequest(t, c, server.URL, "traces", "application/x-protobuf", "", "server-test", first)
	if code != 200 {
		t.Fatal(code)
	}
	d = ordinaryDefinition()
	d["mappings"].([]map[string]any)[2]["step"] = "Changed persistence"
	postMeaning(t, c, server.URL, csrf, "/api/definitions", d, 200)
	body, _ := late.MarshalProto()
	code, _ = protocolRequest(t, c, server.URL, "traces", "application/x-protobuf", "", "server-test", body)
	if code != 200 {
		t.Fatal(code)
	}
	model := readNativeModel(t, c, server.URL)
	if len(model.Journeys) != 1 || model.Journeys[0].Definition != 2 {
		t.Fatal("late earlier parent silently changed captured definition")
	}
	for _, stage := range model.Journeys[0].Path {
		if stage == "Changed persistence" {
			t.Fatal("historical labels remapped")
		}
	}
}
func TestNativeJourneyMissingEndAndSummaryExpiry(t *testing.T) {
	svc, server, c := fixtureService(t)
	csrf := loginTest(t, c, server.URL)
	d := ordinaryDefinition()
	d["goal"] = "missing.parent"
	rules := d["mappings"].([]map[string]any)
	d["mappings"] = append(rules, map[string]any{"event": "missing.parent", "step": "Incomplete operation", "source": "server", "carrier": "span", "kind": "step"})
	postMeaning(t, c, server.URL, csrf, "/api/definitions", d, 200)
	b := scopedNative(t)
	code, _ := protocolRequest(t, c, server.URL, "traces", "application/x-protobuf", "", "server-test", b)
	if code != 200 {
		t.Fatal(code)
	}
	model := readNativeModel(t, c, server.URL)
	if len(model.Journeys) != 1 || model.Journeys[0].Status == "completed" {
		t.Fatal("missing end completed native goal")
	}
	for _, timing := range model.Journeys[0].Timings {
		if timing.Stage == "Incomplete operation" && timing.Complete {
			t.Fatal("missing end became summary sample")
		}
	}
	now := svc.Store.Now()
	svc.Store.Now = func() time.Time { return now.Add(31 * 24 * time.Hour) }
	if err := svc.Store.Cleanup(); err != nil {
		t.Fatal(err)
	}
	getTraceJSON(t, c, server.URL, "/api/traces/11111111111111111111111111111111", 404)
	code, _ = protocolRequest(t, c, server.URL, "traces", "application/x-protobuf", "", "server-test", b)
	if code != 400 {
		t.Fatal("expired replay resurrected detail")
	}
	day := now.UTC().Truncate(24 * time.Hour)
	summary := getTraceJSON(t, c, server.URL, "/api/snapshot?from="+url.QueryEscape(day.Format(time.RFC3339Nano))+"&to="+url.QueryEscape(day.Add(24*time.Hour).Format(time.RFC3339Nano)), 200)
	summaryJSON, _ := json.Marshal(summary["summary"])
	if !strings.Contains(string(summaryJSON), "native-inclusive") || strings.Contains(string(summaryJSON), "synthetic-native-scope") {
		t.Fatal("protected identity-free native summary unavailable")
	}
	rows, err := svc.Store.DB.Query(`SELECT data FROM summaries`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	found := false
	for rows.Next() {
		var raw string
		rows.Scan(&raw)
		if strings.Contains(raw, "synthetic-native-scope") || strings.Contains(raw, "11111111111111111111111111111111") || strings.Contains(raw, "0000000000000001") {
			t.Fatal("identity retained in long-lived summary")
		}
		if strings.Contains(raw, "native-inclusive") {
			found = true
		}
	}
	if !found {
		t.Fatal("native summary basis missing")
	}
}
func TestTraceReferenceModeRetainsDefinitionAndExpires(t *testing.T) {
	svc, server, c := fixtureService(t)
	csrf := loginTest(t, c, server.URL)
	b := scopedNative(t)
	code, _ := protocolRequest(t, c, server.URL, "traces", "application/x-protobuf", "", "server-test", b)
	if code != 200 {
		t.Fatal(code)
	}
	v := postMeaning(t, c, server.URL, csrf, "/api/views", map[string]any{"name": "Native trace reference", "mode": "trace", "trace": "11111111111111111111111111111111", "definition": 1}, 200)
	id := fmt.Sprint(v["id"])
	d := getTraceJSON(t, c, server.URL, "/api/traces/11111111111111111111111111111111?reference="+id, 200)
	if d["mode"] != "trace" {
		t.Fatal("reference mode lost")
	}
	postMeaning(t, c, server.URL, csrf, "/api/definitions", ordinaryDefinition(), 200)
	d = getTraceJSON(t, c, server.URL, "/api/traces/11111111111111111111111111111111?reference="+id, 200)
	for _, v := range d["business"].([]any) {
		if v.(map[string]any)["definition"] != float64(1) {
			t.Fatal("trace reference silently reinterpreted")
		}
	}
	now := svc.Store.Now()
	svc.Store.Now = func() time.Time { return now.Add(31 * 24 * time.Hour) }
	getTraceJSON(t, c, server.URL, "/api/traces/11111111111111111111111111111111?reference="+id, 404)
}

func TestNativeAndPagePresenceBasesStaySeparate(t *testing.T) {
	now := time.Now().UnixNano()
	base := now - int64(time.Minute)
	d := Definition{Version: 1, Goal: "done", Mappings: []Mapping{{Event: "operation", Step: "Operation", Kind: "step", Carrier: "span", Timing: "native-span"}, {Event: "page.start", Step: "Voting", Kind: "stage-start", Timing: "paired"}, {Event: "page.end", Step: "Voting", Kind: "stage-end", Timing: "paired"}, {Event: "presence", Step: "Presence", Kind: "activity"}, {Event: "done", Step: "Done", Kind: "step"}}}
	es := []Evidence{}
	add := func(name, signal string, offset time.Duration, attrs map[string]string) {
		es = append(es, Evidence{ID: int64(len(es) + 1), Name: name, Signal: signal, Source: "browser", Service: "voting-ui", Scope: "scope", Tab: "tab", Time: base + int64(offset), Record: fmt.Sprint(len(es) + 1), Definition: 1, Attrs: attrs})
	}
	add("page.start", "logs", 0, map[string]string{})
	es[0].Occurrence = "page-pair"
	add("presence", "logs", 0, map[string]string{"activity.state": "visible"})
	add("operation", "traces", time.Second, map[string]string{})
	es[2].Trace = "native-trace"
	es[2].Span = "native-span"
	es[2].End = base + int64(2*time.Second)
	add("presence", "logs", time.Second, map[string]string{"activity.state": "hidden"})
	add("presence", "logs", 2*time.Second, map[string]string{"activity.state": "active"})
	add("presence", "logs", 3*time.Second, map[string]string{"activity.state": "idle"})
	add("page.end", "logs", 4*time.Second, map[string]string{})
	es[6].Occurrence = "page-pair"
	add("done", "logs", 5*time.Second, map[string]string{})
	m := Derive(es, []Definition{d}, now)
	if len(m.Journeys) != 1 {
		t.Fatal("explicit tab scope fragmented")
	}
	bases := map[string]float64{}
	for _, timing := range m.Journeys[0].Timings {
		if timing.Complete {
			bases[timing.Basis] = timing.Milliseconds
		}
	}
	if bases["native-inclusive"] != 1000 || bases["measured-elapsed"] != 4000 || bases["observed-active"] != 2000 {
		t.Fatalf("native/pair/presence missing or pooled: %v", bases)
	}
	// Removing explicit scope/tab, even with valid technical trace IDs, must not
	// infer a person or a journey from time proximity or service identity.
	for i := range es {
		es[i].Scope = ""
		es[i].Tab = ""
		es[i].Attrs["http.client_ip"] = "192.0.2.1"
	}
	if m := Derive(es, []Definition{d}, now); len(m.Journeys) != 0 {
		t.Fatal("trace/time/IP fabricated journey")
	}
}

func TestNativeScopeDeletionAndBrowserSeparation(t *testing.T) {
	_, server, c := fixtureService(t)
	csrf := loginTest(t, c, server.URL)
	postMeaning(t, c, server.URL, csrf, "/api/definitions", ordinaryDefinition(), 200)
	b := scopedNative(t)
	code, _ := protocolRequest(t, c, server.URL, "traces", "application/x-protobuf", "", "browser-test", b)
	if code != 200 {
		t.Fatal(code)
	}
	if len(readNativeModel(t, c, server.URL).Journeys) != 0 {
		t.Fatal("browser without tab or server authority fabricated journey")
	}
	code, _ = protocolRequest(t, c, server.URL, "traces", "application/x-protobuf", "", "server-test", b)
	if code != 200 {
		t.Fatal(code)
	}
	if len(readNativeModel(t, c, server.URL).Journeys) != 1 {
		t.Fatal("server-only explicit scope unavailable")
	}
	postMeaning(t, c, server.URL, csrf, "/api/delete", map[string]string{"scope": "synthetic-native-scope"}, 200)
	getTraceJSON(t, c, server.URL, "/api/traces/11111111111111111111111111111111", 404)
	code, _ = protocolRequest(t, c, server.URL, "traces", "application/x-protobuf", "", "server-test", b)
	if code != 200 {
		t.Fatal(code)
	}
	getTraceJSON(t, c, server.URL, "/api/traces/11111111111111111111111111111111", 404)
}

func TestNativeOccurrenceConflictingTimingBasesDoNotDoublePool(t *testing.T) {
	base := time.Now().Add(-time.Minute).UnixNano()
	occurrence := "server/t/s"
	d := Definition{Version: 1, Goal: "done", Mappings: []Mapping{{Event: "op", Step: "Voting", Kind: "step", Carrier: "span", Timing: "native-span"}, {Event: "begin", Step: "Voting", Kind: "stage-start", Timing: "paired"}, {Event: "end", Step: "Voting", Kind: "stage-end", Timing: "paired"}, {Event: "done", Step: "Done", Kind: "step"}}}
	es := []Evidence{{ID: 1, Name: "op", Signal: "traces", Source: "server", Scope: "scope", Trace: "t", Span: "s", Time: base, End: base + 1000, Definition: 1, Record: "t/s"}, {ID: 2, Name: "begin", Signal: "logs", Source: "server", Scope: "scope", Time: base, Definition: 1, Occurrence: occurrence, Record: "begin"}, {ID: 3, Name: "end", Signal: "logs", Source: "server", Scope: "scope", Time: base + 1000, Definition: 1, Occurrence: occurrence, Record: "end"}, {ID: 4, Name: "done", Signal: "logs", Source: "server", Scope: "scope", Time: base + 2000, Definition: 1, Record: "done"}}
	m := Derive(es, []Definition{d}, base+3000)
	for _, timing := range m.Journeys[0].Timings {
		if timing.Complete {
			t.Fatal("competing timing policies silently double-pooled occurrence")
		}
	}
	d.Mappings[0].Timing = "none"
	m = Derive(es, []Definition{d}, base+3000)
	if len(m.Journeys[0].Timings) != 1 || !m.Journeys[0].Timings[0].Complete || m.Journeys[0].Timings[0].Basis != "measured-elapsed" {
		t.Fatal("explicitly choosing paired basis did not yield exactly one sample")
	}
}
