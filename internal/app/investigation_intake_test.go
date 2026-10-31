package app

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
)

// Real OTLP/HTTP intake and real five-second reset/snapshot reconciliation are
// deliberately distinct from the high-rate synthetic browser notification test.
func TestInvestigationActualOTLPFixedBaselineResetAndReplay(t *testing.T) {
	registry, _, server, client := managedFixture(t)
	app := mustRegister(t, registry, "Actual intake isolation")
	csrf := loginTest(t, client, server.URL)
	base := "/api/apps/" + app.ID + "/"
	def := Definition{Mappings: []Mapping{{Event: "A", Step: "A", Source: "browser", Kind: "step"}, {Event: "B", Step: "B", Source: "browser", Kind: "step"}}}
	if code, _ := appRequest(t, client, server.URL, "POST", base+"definitions", csrf, def); code != 200 {
		t.Fatalf("definition status %d", code)
	}
	liveStart := time.Now().UTC().Add(-time.Minute)
	logs := func(scope string, at time.Time) plog.Logs {
		ls := logFixture("A", "B")
		rows := ls.ResourceLogs().At(0).ScopeLogs().At(0).LogRecords()
		for i := 0; i < rows.Len(); i++ {
			r := rows.At(i)
			r.SetTimestamp(pcommon.NewTimestampFromTime(at.Add(time.Duration(i) * time.Millisecond)))
			r.SetObservedTimestamp(pcommon.NewTimestampFromTime(time.Now()))
			r.Attributes().PutStr("event.id", scope+"-"+r.EventName())
			r.Attributes().PutStr("session.id", scope)
			r.Attributes().PutInt("source.sequence", int64(i+1))
		}
		return ls
	}
	baseline := logs("baseline-before", liveStart.Add(-time.Hour))
	live := logs("selected-live", liveStart.Add(time.Second))
	for _, ls := range []plog.Logs{baseline, live} {
		if code := requestLogs(t, client, server.URL, app.Browser, ls); code != 200 {
			t.Fatalf("intake %d", code)
		}
	}
	query := "snapshot?mode=live&baseline=1&scope=selected-live&liveStart=" + url.QueryEscape(liveStart.Format(time.RFC3339Nano))
	type snapshot struct {
		Revision string           `json:"revision"`
		Model    Model            `json:"model"`
		Baseline Model            `json:"baselineModel"`
		Bounds   map[string]int64 `json:"baselineBounds"`
		Coverage struct {
			Source struct {
				CollectionCompleteness string `json:"collectionCompleteness"`
			} `json:"source"`
			Current  PopulationCoverage `json:"current"`
			Baseline PopulationCoverage `json:"baseline"`
		} `json:"coverage"`
	}
	read := func() snapshot {
		t.Helper()
		code, raw := appRequest(t, client, server.URL, "GET", base+query, "", nil)
		if code != 200 {
			t.Fatalf("snapshot %d", code)
		}
		var s snapshot
		if err := json.Unmarshal(raw, &s); err != nil {
			t.Fatal(err)
		}
		return s
	}
	first := read()
	if first.Model.Stats.Records != 2 || len(first.Model.Journeys) != 1 || first.Model.Journeys[0].Scope != "selected-live" || first.Baseline.Stats.Records != 2 {
		t.Fatal("initial populations not independently derived")
	}
	if first.Coverage.Current.Scope != "selected-live" || first.Coverage.Baseline.Scope != "" || first.Coverage.Source.CollectionCompleteness != "unknown; exporter sampling/loss not observed" {
		t.Fatal("scope/completeness contract")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, "GET", server.URL+base+"events?generation="+app.Generation, nil)
	response, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	reader := bufio.NewReader(response.Body)
	reset := func() string {
		t.Helper()
		var id string
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					t.Fatal("reset stream ended")
				}
				t.Fatal(err)
			}
			if strings.HasPrefix(line, "id: ") {
				id = strings.TrimSpace(strings.TrimPrefix(line, "id: "))
			}
			if line == "\n" {
				if id == "" {
					t.Fatal("missing reset id")
				}
				return id
			}
		}
	}
	initialID := reset()
	started := time.Now()
	late := logs("baseline-late", liveStart.Add(-30*time.Minute))
	if code := requestLogs(t, client, server.URL, app.Browser, late); code != 200 {
		t.Fatalf("late intake %d", code)
	}
	second := read()
	if second.Model.Stats.Records != 2 || second.Baseline.Stats.Records != 4 || len(second.Baseline.Journeys) != 2 {
		t.Fatal("receipt time moved late baseline evidence into live selection")
	}
	if second.Bounds["from"] != first.Bounds["from"] || second.Bounds["to"] != first.Bounds["to"] {
		t.Fatal("baseline bounds shifted")
	}
	if second.Coverage.Current.Scope != "selected-live" || second.Coverage.Baseline.Scope != "" {
		t.Fatal("late intake changed scope")
	}
	periodicID := reset()
	elapsed := time.Since(started)
	if elapsed < 4500*time.Millisecond || elapsed > 7*time.Second || periodicID == initialID {
		t.Fatalf("actual periodic reset cadence/revision: %s changed=%v", elapsed, periodicID != initialID)
	}
	for _, ls := range []plog.Logs{baseline, live, late} {
		if code := requestLogs(t, client, server.URL, app.Browser, ls); code != 200 {
			t.Fatalf("exact replay %d", code)
		}
	}
	replay := read()
	if replay.Model.Stats.Records != second.Model.Stats.Records || replay.Baseline.Stats.Records != second.Baseline.Stats.Records || len(replay.Model.Edges) != len(second.Model.Edges) || len(replay.Baseline.Edges) != len(second.Baseline.Edges) {
		t.Fatal("replay inflated retained contributions")
	}
	t.Logf("actual OTLP browser logs: live2/baseline2→4; selected live/app-wide baseline; fixed bounds; exact replay stable; periodic reset %.3fs", elapsed.Seconds())
}
