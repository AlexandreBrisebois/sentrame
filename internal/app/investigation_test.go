package app

import (
	"encoding/json"
	"fmt"
	"net/url"
	"testing"
	"time"
)

func TestInvestigationReferenceKeepsLiveContinuationAndPreferences(t *testing.T) {
	svc, server, c := fixtureService(t)
	csrf := loginTest(t, c, server.URL)
	now := time.Now().UTC()
	svc.Store.Insert([]Evidence{
		{Source: "browser", Signal: "logs", Name: "discovery.entered", Scope: "continuation", Tab: "tab", Record: "before", Time: now.Add(-time.Minute).UnixNano(), Attrs: map[string]string{}},
		{Source: "browser", Signal: "logs", Name: "completion.entered", Scope: "continuation", Tab: "tab", Record: "inside", Time: now.Add(-time.Second).UnixNano(), Attrs: map[string]string{}},
	})
	input := map[string]any{"name": "Frozen continuation", "mode": "overview", "layout": "matrix", "from": now.Add(-30 * time.Second).Format(time.RFC3339Nano), "to": now.Format(time.RFC3339Nano), "selectionPolicy": "evidence-window", "presentation": map[string]any{"version": 1, "metric": "journeys", "stages": []string{"Discovery", "Completion"}, "sort": "traffic", "groupRepeats": true}}
	code, body := appRequest(t, c, server.URL, "POST", "/api/views", csrf, input)
	if code != 200 {
		t.Fatalf("save %d %s", code, body)
	}
	var saved map[string]any
	json.Unmarshal(body, &saved)
	if saved["selectionPolicy"] != "evidence-window" || saved["presentation"] == nil {
		t.Fatalf("reference contract dropped: %s", body)
	}
	snap := snapshotTest(t, c, server.URL, "?reference="+fmt.Sprint(saved["id"]))
	var model Model
	json.Unmarshal(snap["model"], &model)
	if len(model.Journeys) != 1 || len(model.Journeys[0].Evidence) != 1 || !model.Journeys[0].Partial {
		t.Fatalf("lost or broadened continuation: %+v", model.Journeys)
	}
	// The old start-cohort contract remains the default.
	old := snapshotTest(t, c, server.URL, "?from="+url.QueryEscape(input["from"].(string))+"&to="+url.QueryEscape(input["to"].(string)))
	json.Unmarshal(old["model"], &model)
	if len(model.Journeys) != 0 {
		t.Fatal("old default eligibility changed")
	}
	for _, bad := range []map[string]any{
		{"selectionPolicy": "unknown"},
		{"presentation": map[string]any{"version": 2}},
		{"presentation": map[string]any{"version": 1, "metric": "people"}},
		{"presentation": map[string]any{"version": 1, "stages": make([]string, 33)}},
	} {
		candidate := map[string]any{}
		for k, v := range input {
			candidate[k] = v
		}
		for k, v := range bad {
			candidate[k] = v
		}
		if code, _ := appRequest(t, c, server.URL, "POST", "/api/views", csrf, candidate); code != 400 {
			t.Fatalf("invalid preferences accepted: %+v (%d)", bad, code)
		}
	}
}

func TestInvestigationCoverageSeparatesEvidencePageFromPopulation(t *testing.T) {
	svc, server, c := fixtureService(t)
	loginTest(t, c, server.URL)
	now := time.Now().UTC()
	es := make([]Evidence, 2002)
	for i := range es {
		es[i] = Evidence{Source: "server", Signal: "logs", Name: "unmapped", Record: fmt.Sprint(i), Time: now.Add(-time.Minute).UnixNano() + int64(i), Attrs: map[string]string{}}
	}
	if e := svc.Store.Insert(es); e != nil {
		t.Fatal(e)
	}
	snap := snapshotTest(t, c, server.URL, "")
	var coverage struct {
		Source  struct{ Truncated bool }
		Current struct {
			SourcePartial     bool
			GraphNodesOmitted int
			GraphEdgesOmitted int
		}
		Evidence struct{ Omitted int }
	}
	if err := json.Unmarshal(snap["coverage"], &coverage); err != nil {
		t.Fatalf("missing independent coverage: %v", err)
	}
	if coverage.Source.Truncated || coverage.Current.SourcePartial || coverage.Evidence.Omitted != 2 {
		t.Fatalf("list cap contaminated source/population: %+v", coverage)
	}
}

func TestInvestigationRejectsInvalidHistoricalBounds(t *testing.T) {
	_, server, c := fixtureService(t)
	csrf := loginTest(t, c, server.URL)
	now := time.Now().UTC()
	for _, bounds := range [][2]string{{"nonsense", "nonsense"}, {now.Format(time.RFC3339Nano), now.Format(time.RFC3339Nano)}, {now.Format(time.RFC3339Nano), now.Add(-time.Hour).Format(time.RFC3339Nano)}, {now.Add(-time.Hour).Format(time.RFC3339Nano), now.Add(time.Hour).Format(time.RFC3339Nano)}} {
		code, _ := appRequest(t, c, server.URL, "POST", "/api/views", csrf, map[string]any{"name": "Invalid frozen window", "mode": "overview", "selectionPolicy": "evidence-window", "from": bounds[0], "to": bounds[1]})
		if code != 400 {
			t.Fatalf("invalid historical bounds accepted %v (%d)", bounds, code)
		}
	}
}

func TestInvestigationEvidenceListCapPreservesMappedDenominators(t *testing.T) {
	svc, server, c := fixtureService(t)
	loginTest(t, c, server.URL)
	now := time.Now().UTC()
	es := make([]Evidence, 2002)
	for i := range es {
		es[i] = Evidence{Source: "browser", Signal: "logs", Name: "discovery.entered", Scope: "capped-list", Tab: "one-tab", Record: fmt.Sprint(i), Time: now.Add(-time.Minute).UnixNano() + int64(i), Sequence: int64(i + 1), Attrs: map[string]string{}}
	}
	if err := svc.Store.Insert(es); err != nil {
		t.Fatal(err)
	}
	snap := snapshotTest(t, c, server.URL, "")
	var model Model
	var records []Evidence
	var coverage struct {
		Source   struct{ Truncated bool }
		Current  PopulationCoverage
		Evidence struct{ Omitted int }
	}
	for _, target := range []struct {
		key   string
		value any
	}{{"model", &model}, {"evidence", &records}, {"coverage", &coverage}} {
		if err := json.Unmarshal(snap[target.key], target.value); err != nil {
			t.Fatal(err)
		}
	}
	if len(records) != 2000 || coverage.Evidence.Omitted != 2 || coverage.Source.Truncated || coverage.Current.SourcePartial {
		t.Fatalf("list cap contaminated source/population: records=%d coverage=%+v", len(records), coverage)
	}
	if len(model.Journeys) != 1 || len(model.Journeys[0].Steps) != 2002 || len(model.Journeys[0].Evidence) != 2002 || len(model.Edges) != 1 || model.Edges[0].Occurrences != 2001 {
		t.Fatalf("mapped ordered denominator lost: journeys=%d edges=%+v", len(model.Journeys), model.Edges)
	}
}
func TestInvestigationEndpointKeepsCapturedOutcomeAfterClipping(t *testing.T) {
	j := Journey{ID: "captured", Start: 1, End: 4, LastMapped: 4, Status: "completed", Evidence: []Evidence{{ID: 1, Name: "A", Time: 1}, {ID: 2, Name: "X", Time: 2}, {ID: 3, Name: "B", Time: 4}}, Steps: []StepOccurrence{{Step: "A", Event: "A", EvidenceID: 1, Time: 1}, {Step: "X", Event: "X", EvidenceID: 2, Time: 2}, {Step: "B", Event: "B", EvidenceID: 3, Time: 4}}, Path: []string{"A", "X", "B"}}
	clipped := ClipJourney(j, 1, 3)
	out, ok := referenceJourney(clipped, SavedView{Mode: "goal-backward", End: "X"})
	if !ok || out.CapturedStatus != "completed" {
		t.Fatalf("captured outcome overwritten: %+v", out)
	}
}
func TestInvestigationTraceReferenceFailsClosedOnUnsupportedPolicy(t *testing.T) {
	svc, server, c := fixtureService(t)
	loginTest(t, c, server.URL)
	for _, policy := range []string{"unknown", "evidence-window"} {
		raw, _ := json.Marshal(map[string]any{"mode": "trace", "trace": "11111111111111111111111111111111", "selectionPolicy": policy})
		result, err := svc.Store.DB.Exec(`INSERT INTO views(data) VALUES(?)`, string(raw))
		if err != nil {
			t.Fatal(err)
		}
		id, _ := result.LastInsertId()
		code, _ := appRequest(t, c, server.URL, "GET", "/api/traces?reference="+fmt.Sprint(id), "", nil)
		if code != 400 {
			t.Fatalf("unsupported trace policy %s not rejected (%d)", policy, code)
		}
	}
}
