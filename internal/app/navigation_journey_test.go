package app

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
)

func TestCorrelatedNavigationWithoutMappingOrGoal(t *testing.T) {
	d := Definition{Version: 1}
	home := ev(1, "pwa.app.opened", "browser", 0)
	home.Attrs["stage.name"] = "home"
	home.Sequence = 1
	search := ev(2, "pwa.screen.presented", "browser", 1)
	search.Attrs["stage.name"] = "search"
	search.Sequence = 2
	bad := ev(3, "pwa.error", "browser", 2)
	bad.Severity = 17
	bad.Sequence = 3
	planner := ev(4, "pwa.screen.presented", "browser", 3)
	planner.Attrs["stage.name"] = "planner"
	planner.Sequence = 4
	m := Derive([]Evidence{planner, bad, home, search}, []Definition{d}, epoch+int64(10*time.Second))
	if len(m.Journeys) != 1 {
		t.Fatalf("correlated navigation must form one journey: %+v", m)
	}
	j := m.Journeys[0]
	if !reflect.DeepEqual(j.Path, []string{"home", "search", "planner"}) || len(j.Evidence) != 4 || j.Status != "active" {
		t.Fatalf("navigation/error evidence or ordering lost: %+v", j)
	}
	if len(m.Edges) != 2 || m.Stats.Completed != 0 || m.Stats.EligibleJourneys != 0 || len(j.Timings) != 0 {
		t.Fatalf("navigation must connect without fabricated completion, conversion or duration: %+v", m)
	}
	if len(m.Diagnostics) != 1 || m.Diagnostics[0].JourneyID != j.ID {
		t.Fatal("correlated error must attach to navigation")
	}
	m = Derive([]Evidence{home, search, bad, planner}, []Definition{d}, epoch+int64(31*time.Minute))
	if m.Journeys[0].Status != "inactive" || m.Stats.Completed != 0 {
		t.Fatal("silence must close inactive with no completion")
	}
}

func TestCorrelatedNavigationIsolationAndConfiguredPrecedence(t *testing.T) {
	a := ev(1, "screen", "browser", 0)
	a.Attrs["stage.name"] = "home"
	b := ev(2, "screen", "browser", 1)
	b.Tab = "another-tab"
	c := ev(3, "screen", "browser", 2)
	c.Scope = "another-session"
	noTab := ev(4, "screen", "browser", 3)
	noTab.Tab = ""
	server := ev(5, "screen", "server", 4)
	metric := ev(6, "screen", "browser", 5)
	metric.Signal = "metrics"
	trace := ev(7, "screen", "browser", 6)
	trace.Signal = "traces"
	m := Derive([]Evidence{a, b, c, noTab, server, metric, trace}, []Definition{{Version: 1}}, epoch+int64(10*time.Second))
	if len(m.Journeys) != 3 || len(m.Unassigned) != 4 {
		t.Fatalf("navigation must isolate tabs/scopes and not infer trace/server/metric paths: %+v", m)
	}
	d := Definition{Version: 1, Goal: "screen", Mappings: []Mapping{{Event: "screen", Step: "Owner label", Kind: "step", Source: "browser"}}}
	m = Derive([]Evidence{a}, []Definition{d}, epoch+int64(time.Second))
	if len(m.Journeys) != 1 || m.Journeys[0].Path[0] != "Owner label" || m.Stats.Completed != 1 {
		t.Fatal("configured meaning and completion must take precedence")
	}
}

func TestMappingWithoutCompletionGoal(t *testing.T) {
	d := Definition{Version: 1, Mappings: []Mapping{{Event: "screen", Step: "Screen", Kind: "step", Source: "browser"}}}
	if msg := validateDefinition(d); msg != "" {
		t.Fatalf("goal-free mapping must be valid: %s", msg)
	}
	m := Derive([]Evidence{ev(1, "screen", "browser", 0)}, []Definition{d}, epoch+int64(time.Second))
	if len(m.Journeys) != 1 || m.Stats.Completed != 0 || m.Stats.EligibleJourneys != 0 {
		t.Fatal("goal-free mappings must not imply conversion")
	}
}

func TestManagedNavigationIntakeSnapshotReplayAndGoalFreeActivation(t *testing.T) {
	m, _, server, c := managedFixture(t)
	a := mustRegister(t, m, "Browser navigation")
	other := mustRegister(t, m, "Other navigation app")
	csrf := loginTest(t, c, server.URL)
	ls := plog.NewLogs()
	r := ls.ResourceLogs().AppendEmpty()
	r.Resource().Attributes().PutStr("service.name", "wfs-pwa-browser")
	logs := r.ScopeLogs().AppendEmpty().LogRecords()
	start := time.Now().Add(-10 * time.Second)
	for i, screen := range []string{"home", "search", "planner"} {
		l := logs.AppendEmpty()
		l.SetEventName("pwa.screen.presented")
		l.SetTimestamp(pcommon.NewTimestampFromTime(start.Add(time.Duration(i) * time.Second)))
		l.SetSeverityNumber(9)
		l.Attributes().PutStr("stage.name", screen)
		l.Attributes().PutStr("session.id", "random-visit")
		l.Attributes().PutStr("tab.id", "random-tab")
		l.Attributes().PutStr("event.id", fmt.Sprintf("record-%d", i))
		l.Attributes().PutInt("source.sequence", int64(i+1))
	}
	for attempt := 0; attempt < 2; attempt++ {
		if status := requestLogs(t, c, server.URL, a.Browser, ls); status != 200 {
			t.Fatalf("intake: %d", status)
		}
	}
	base := "/api/apps/" + a.ID
	code, raw := appRequest(t, c, server.URL, "GET", base+"/snapshot", "", nil)
	var snapshot struct{ Model Model }
	if code != 200 || json.Unmarshal(raw, &snapshot) != nil || len(snapshot.Model.Journeys) != 1 {
		t.Fatalf("protected navigation read: %d %s", code, raw)
	}
	j := snapshot.Model.Journeys[0]
	if !j.NavigationOnly || !reflect.DeepEqual(j.Path, []string{"home", "search", "planner"}) || len(j.Evidence) != 3 || len(snapshot.Model.Edges) != 2 {
		t.Fatalf("intake/replay lost navigation: %+v", j)
	}
	projected, ok := referenceJourney(j, SavedView{Mode: "segment", Start: "search", End: "planner"})
	if !ok || !reflect.DeepEqual(projected.Path, []string{"search", "planner"}) {
		t.Fatalf("navigation stage reference lost endpoints: %+v", projected)
	}
	code, raw = appRequest(t, c, server.URL, "GET", "/api/apps/"+other.ID+"/snapshot", "", nil)
	if code != 200 || json.Unmarshal(raw, &snapshot) != nil || len(snapshot.Model.Journeys) != 0 {
		t.Fatal("navigation crossed app boundary")
	}
	d := Definition{Mappings: []Mapping{{Event: "pwa.screen.presented", Step: "Screen", Source: "browser", Kind: "step"}}}
	code, raw = appRequest(t, c, server.URL, "POST", base+"/definitions", csrf, d)
	if code != 200 {
		t.Fatalf("goal-free definition activation: %d %s", code, raw)
	}
	code, raw = appRequest(t, c, server.URL, "GET", base+"/snapshot", "", nil)
	if code != 200 || json.Unmarshal(raw, &snapshot) != nil || len(snapshot.Model.Journeys) != 1 || !reflect.DeepEqual(snapshot.Model.Journeys[0].Path, j.Path) {
		t.Fatal("new mapping must preserve captured navigation interpretation")
	}
}

func TestCorrelatedNavigationMissingSequenceAndInactivityBoundary(t *testing.T) {
	a := ev(1, "opened", "browser", 0)
	b := ev(2, "screen", "browser", 1)
	c := ev(3, "screen", "browser", 1900)
	missingScope := ev(4, "screen", "browser", 2)
	missingScope.Scope = ""
	bad := ev(5, "error", "browser", -1)
	bad.Severity = 17
	m := Derive([]Evidence{c, b, a, missingScope, bad}, []Definition{{Version: 1}}, epoch+int64(1901*time.Second))
	if len(m.Journeys) != 2 || m.Journeys[0].Status != "inactive" || m.Journeys[1].Status != "active" || len(m.Unassigned) != 2 {
		t.Fatalf("inactivity, missing context or pre-navigation diagnostic boundary lost: %+v", m)
	}
	if !reflect.DeepEqual(m.Journeys[0].Path, []string{"opened", "screen"}) {
		t.Fatal("missing sequence must fall back to occurrence order")
	}
	warning := "Display uses occurrence time and stable receipt tie-break; causal order is unproven."
	found := false
	for _, w := range m.Journeys[0].Warnings {
		found = found || w == warning
	}
	if !found {
		t.Fatal("missing sequence order uncertainty must be disclosed")
	}
}
