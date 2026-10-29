package app

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace/ptraceotlp"
	"io"
	"net/http"
	"testing"
	"time"
)

func getTraceJSON(t *testing.T, c *http.Client, base, path string, want int) map[string]any {
	t.Helper()
	r, err := c.Get(base + path)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Body.Close()
	b, _ := io.ReadAll(r.Body)
	if r.StatusCode != want {
		t.Fatalf("%s status %d want %d: %s", path, r.StatusCode, want, b)
	}
	out := map[string]any{}
	if want == 200 {
		if err = json.Unmarshal(b, &out); err != nil {
			t.Fatal(err)
		}
	}
	return out
}
func TestTraceExplorerStandardProtectedStructure(t *testing.T) {
	_, server, c := fixtureService(t)
	tb, lb, _ := standardBodies(t)
	for signal, body := range map[string][]byte{"traces": tb, "logs": lb} {
		code, _ := protocolRequest(t, c, server.URL, signal, "application/x-protobuf", "", "server-test", body)
		if code != 200 {
			t.Fatal(code)
		}
	}
	path := "/api/traces/11111111111111111111111111111111"
	getTraceJSON(t, c, server.URL, path, 401)
	loginTest(t, c, server.URL)
	list := getTraceJSON(t, c, server.URL, "/api/traces", 200)
	if list["mode"] != "trace" {
		t.Fatal("trace list lacks explicit mode")
	}
	detail := getTraceJSON(t, c, server.URL, path, 200)
	if detail["mode"] != "trace" || detail["computation"] != "native-inclusive-v1" || detail["unit"] != "observed spans" {
		t.Fatal("trace semantics missing")
	}
	spans, ok := detail["operations"].([]any)
	if !ok || len(spans) != 5 {
		t.Fatalf("want five native operations: %v", detail)
	}
	for _, v := range spans {
		op := v.(map[string]any)
		e := op["evidence"].(map[string]any)
		switch e["name"] {
		case "POST /votes":
			if op["durationNs"] != "1000000001" || op["timing"] != "complete" {
				t.Fatal("native inclusive precision lost")
			}
		case "missing.parent":
			if op["timing"] != "right-censored" || op["durationNs"] != nil {
				t.Fatal("missing end became duration")
			}
		case "invalid.clock":
			if op["timing"] != "unavailable" || op["durationNs"] != nil {
				t.Fatal("bad clock became sample")
			}
		}
	}
	rels := detail["relationships"].([]any)
	counts := map[string]int{}
	for _, v := range rels {
		r := v.(map[string]any)
		counts[r["type"].(string)]++
		if r["type"] == "parent" && r["state"] == "resolved" {
			if r["targetSpan"] != "0000000000000001" {
				t.Fatal("invented sibling adjacency")
			}
		}
	}
	if counts["parent"] != 3 || counts["span-event"] != 2 || counts["log-context"] != 2 || counts["link"] != 32 {
		t.Fatalf("typed relationships missing: %v", counts)
	}
	events := detail["events"].([]any)
	if events[0].(map[string]any)["name"] != "vote.committed" || events[1].(map[string]any)["name"] != "retry.observed" {
		t.Fatal("events reordered by clock instead of export index")
	}
	if detail["completeness"] == "complete" || detail["revision"] == nil {
		t.Fatal("partial sampling/revision disclosure missing")
	}
	// Browser claims stay separate; duplicate parent candidates become ambiguous.
	code, _ := protocolRequest(t, c, server.URL, "traces", "application/x-protobuf", "", "browser-test", tb)
	if code != 200 {
		t.Fatal(code)
	}
	changed := getTraceJSON(t, c, server.URL, path, 200)
	if changed["revision"] == detail["revision"] {
		t.Fatal("new evidence did not invalidate revision")
	}
	ambiguous := 0
	for _, v := range changed["relationships"].([]any) {
		r := v.(map[string]any)
		if r["type"] == "parent" && r["targetSpan"] == "0000000000000001" {
			if r["state"] != "ambiguous" {
				t.Fatal("cross-source parent authority promoted")
			}
			ambiguous++
		}
	}
	if ambiguous != 4 {
		t.Fatalf("ambiguous parents %d", ambiguous)
	}
}

func TestTraceViewerAllowlistRevocationAndNoMutation(t *testing.T) {
	x := newInvitationHarness(t)
	id, token := x.create(24)
	viewer, csrf := x.redeem(token)
	// Empty list is useful; unknown detail stays unavailable under authorized read.
	if w := x.call("GET", "/api/traces", "", viewer, "", false); w.Code != 200 {
		t.Fatalf("viewer trace read %d", w.Code)
	}
	if w := x.call("GET", "/api/traces/11111111111111111111111111111111", "", viewer, "", false); w.Code != 404 {
		t.Fatalf("unknown trace %d", w.Code)
	}
	if w := x.call("POST", "/api/traces", "{}", viewer, csrf, true); w.Code != 403 {
		t.Fatal("viewer trace mutation allowed")
	}
	if w := x.call("POST", "/api/invitations/revoke", fmt.Sprintf(`{"id":%d}`, id), x.owner, x.csrf, true); w.Code != 200 {
		t.Fatal(w.Code)
	}
	if w := x.call("GET", "/api/traces", "", viewer, "", false); w.Code != 401 {
		t.Fatal("revoked trace read allowed")
	}
}

func TestTraceDetailBoundUsesActualIngestion(t *testing.T) {
	_, server, c := fixtureService(t)
	for batch := 0; batch < 2; batch++ {
		q := ptraceotlp.NewExportRequest()
		sp := q.Traces().ResourceSpans().AppendEmpty().ScopeSpans().AppendEmpty().Spans()
		n := 1001
		if batch == 1 {
			n = 1000
		}
		for i := 0; i < n; i++ {
			p := sp.AppendEmpty()
			p.SetName("bounded.operation")
			p.SetTraceID(pcommon.TraceID{6})
			var id pcommon.SpanID
			binary.BigEndian.PutUint64(id[:], uint64(batch*1001+i+1))
			p.SetSpanID(id)
			p.SetStartTimestamp(pcommon.Timestamp(time.Now().Add(-time.Minute).UnixNano()))
			p.SetEndTimestamp(p.StartTimestamp() + 1)
		}
		body, _ := q.MarshalProto()
		code, _ := protocolRequest(t, c, server.URL, "traces", "application/x-protobuf", "", "server-test", body)
		if code != 200 {
			t.Fatal(code)
		}
	}
	loginTest(t, c, server.URL)
	d := getTraceJSON(t, c, server.URL, "/api/traces/06000000000000000000000000000000", 200)
	if d["truncated"] != true || len(d["operations"].([]any)) != 2000 {
		t.Fatal("trace detail did not disclose bounded omission")
	}
}
