package app

import (
	"encoding/json"
	"net/url"
	"testing"
	"time"
)

// A paused browser still performs authorized lifecycle reads. Verify that those
// reads can remove one population without reconstructing it from the other.
func TestInvestigationLifecycleReadsRemoveDetailWithoutPopulationFallback(t *testing.T) {
	for _, operation := range []string{"delete-live", "expire-both"} {
		t.Run(operation, func(t *testing.T) {
			svc, server, c := fixtureService(t)
			csrf := loginTest(t, c, server.URL)
			now := time.Now().UTC()
			clock := now
			svc.Store.Now = func() time.Time { return clock }
			var events []Evidence
			for _, population := range []string{"baseline", "live"} {
				start := now.Add(-time.Hour)
				if population == "live" {
					start = now.Add(-time.Minute)
				}
				for i, name := range []string{"discovery.entered", "completion.entered"} {
					events = append(events, Evidence{Source: "browser", Signal: "logs", Name: name, Scope: population, Tab: "same-tab", Record: population + name, Time: start.Add(time.Duration(i) * time.Second).UnixNano(), Attrs: map[string]string{}})
				}
			}
			if err := svc.Store.Insert(events); err != nil {
				t.Fatal(err)
			}
			query := "?mode=live&baseline=1&liveStart=" + url.QueryEscape(now.Add(-30*time.Minute).Format(time.RFC3339Nano))
			read := func() (Model, Model) {
				t.Helper()
				snapshot := snapshotTest(t, c, server.URL, query)
				var current, baseline Model
				if err := json.Unmarshal(snapshot["model"], &current); err != nil {
					t.Fatal(err)
				}
				if err := json.Unmarshal(snapshot["baselineModel"], &baseline); err != nil {
					t.Fatal(err)
				}
				return current, baseline
			}
			current, baseline := read()
			if len(current.Journeys) != 1 || len(baseline.Journeys) != 1 || current.Journeys[0].Scope == baseline.Journeys[0].Scope {
				t.Fatal("fixture must expose independently scoped same-name paths")
			}
			if operation == "delete-live" {
				code, raw := appRequest(t, c, server.URL, "POST", "/api/delete", csrf, map[string]string{"scope": "live"})
				if code != 200 {
					t.Fatalf("delete %d %s", code, raw)
				}
			} else {
				// No Cleanup call: protected reads must exclude expired raw rows
				// immediately, rather than wait for a maintenance deletion.
				clock = now.Add(31 * 24 * time.Hour)
			}
			current, baseline = read()
			if len(current.Journeys) != 0 || current.Stats.Records != 0 {
				t.Fatal("removed live evidence survived or borrowed baseline paths")
			}
			if operation == "delete-live" && (len(baseline.Journeys) != 1 || baseline.Journeys[0].Scope != "baseline") {
				t.Fatal("live deletion cleared unrelated baseline evidence")
			}
			if operation == "expire-both" && (len(baseline.Journeys) != 0 || baseline.Stats.Records != 0) {
				t.Fatal("expired baseline detail survived protected read")
			}
		})
	}
}
