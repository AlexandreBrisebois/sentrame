package app

import (
	"bufio"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"go.opentelemetry.io/collector/pdata/pcommon"
)

// This exercises real HTTP stream cancellation and reconnection. The separate
// browser proxy test proves EventSource error/open and the visible gap lifecycle.
func TestInvestigationActualResetReconnectUsesCurrentSnapshot(t *testing.T) {
	registry, _, server, client := managedFixture(t)
	entry := mustRegister(t, registry, "Reconnect isolation")
	csrf := loginTest(t, client, server.URL)
	base := "/api/apps/" + entry.ID + "/"
	def := Definition{Mappings: []Mapping{{Event: "A", Step: "A", Source: "browser", Kind: "step"}, {Event: "B", Step: "B", Source: "browser", Kind: "step"}}}
	if code, _ := appRequest(t, client, server.URL, "POST", base+"definitions", csrf, def); code != 200 {
		t.Fatalf("definition %d", code)
	}
	started := time.Now().Add(-time.Minute).UTC()
	emit := func(scope string, at time.Time) {
		t.Helper()
		ls := logFixture("A", "B")
		rows := ls.ResourceLogs().At(0).ScopeLogs().At(0).LogRecords()
		for i := 0; i < rows.Len(); i++ {
			r := rows.At(i)
			r.SetTimestamp(pcommon.NewTimestampFromTime(at.Add(time.Duration(i) * time.Millisecond)))
			r.SetObservedTimestamp(pcommon.NewTimestampFromTime(time.Now()))
			r.Attributes().PutStr("session.id", scope)
			r.Attributes().PutStr("event.id", scope+"-"+r.EventName())
			r.Attributes().PutInt("source.sequence", int64(i+1))
		}
		if code := requestLogs(t, client, server.URL, entry.Browser, ls); code != 200 {
			t.Fatalf("intake %d", code)
		}
	}
	query := "snapshot?mode=live&baseline=1&liveStart=" + url.QueryEscape(started.Format(time.RFC3339Nano))
	type state struct {
		Model    Model            `json:"model"`
		Baseline Model            `json:"baselineModel"`
		Bounds   map[string]int64 `json:"baselineBounds"`
	}
	read := func() state {
		t.Helper()
		code, b := appRequest(t, client, server.URL, "GET", base+query, "", nil)
		if code != 200 {
			t.Fatalf("snapshot %d", code)
		}
		var s state
		if err := json.Unmarshal(b, &s); err != nil {
			t.Fatal(err)
		}
		return s
	}
	stream := func(cursor string) (string, func()) {
		t.Helper()
		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		req, _ := http.NewRequestWithContext(ctx, "GET", server.URL+base+"events?generation="+entry.Generation, nil)
		if cursor != "" {
			req.Header.Set("Last-Event-ID", cursor)
		}
		r, err := client.Do(req)
		if err != nil {
			cancel()
			t.Fatal(err)
		}
		if r.StatusCode != 200 {
			r.Body.Close()
			cancel()
			t.Fatalf("events %d", r.StatusCode)
		}
		reader := bufio.NewReader(r.Body)
		id := ""
		event := ""
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				r.Body.Close()
				cancel()
				t.Fatal(err)
			}
			if strings.HasPrefix(line, "event: ") {
				event = strings.TrimSpace(strings.TrimPrefix(line, "event: "))
			}
			if strings.HasPrefix(line, "id: ") {
				id = strings.TrimSpace(strings.TrimPrefix(line, "id: "))
			}
			if line == "\n" {
				break
			}
		}
		if event != "reset" || !strings.HasPrefix(id, entry.ID+":"+entry.Generation+":") {
			t.Fatalf("unexpected reset %q %q", event, id)
		}
		return id, func() { r.Body.Close(); cancel() }
	}
	emit("live-before", started.Add(time.Second))
	emit("baseline-before", started.Add(-time.Hour))
	before := read()
	first, closeFirst := stream("")
	closeFirst() // Actual client-side HTTP disconnect.
	emit("live-during-disconnect", started.Add(2*time.Second))
	emit("baseline-late", started.Add(-30*time.Minute))
	second, closeSecond := stream(first)
	closeSecond()
	if second == first {
		t.Fatal("reconnected reset used obsolete revision")
	}
	after := read()
	if after.Model.Stats.Records != 4 || after.Baseline.Stats.Records != 4 || len(after.Model.Journeys) != 2 || len(after.Baseline.Journeys) != 2 {
		t.Fatal("reconnect did not reconcile both populations")
	}
	if before.Bounds["from"] != after.Bounds["from"] || before.Bounds["to"] != after.Bounds["to"] {
		t.Fatal("reconnect shifted fixed baseline bounds")
	}
	third, closeThird := stream(second)
	closeThird()
	again := read()
	if third != second || again.Model.Stats.Records != 4 || again.Baseline.Stats.Records != 4 {
		t.Fatal("unchanged reset inflated evidence or changed revision")
	}
	req, _ := http.NewRequest("GET", server.URL+base+"events?generation="+entry.Generation, nil)
	req.Header.Set("Last-Event-ID", "another-app:another-generation:1")
	r, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	r.Body.Close()
	if r.StatusCode != 400 {
		t.Fatalf("foreign cursor accepted %d", r.StatusCode)
	}
	t.Log("actual HTTP disconnect/reconnect: current revision reset, no replay/count deltas, fixed baseline, foreign cursor rejected")
}
