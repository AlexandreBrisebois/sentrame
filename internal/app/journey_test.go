package app

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func def() Definition {
	return Definition{Version: 1, Goal: "completion.entered", Mappings: []Mapping{{Event: "discovery.entered", Step: "Discovery", Kind: "step"}, {Event: "vote.requested", Step: "Intent", Kind: "step", Source: "browser"}, {Event: "vote.committed", Step: "Vote", Kind: "step", Source: "server"}, {Event: "completion.entered", Step: "Completion", Kind: "step", Source: "browser"}, {Event: "journey.exited", Step: "Exit", Kind: "exit"}, {Event: "stage.entered", Step: "Viewing", Kind: "stage-start"}, {Event: "stage.exited", Step: "Viewing", Kind: "stage-end"}, {Event: "activity.changed", Kind: "activity"}}}
}

var epoch = time.Date(2026, 10, 28, 12, 0, 0, 0, time.UTC).UnixNano()

func ev(id int, name, source string, seconds int) Evidence {
	return Evidence{ID: int64(id), Name: name, Signal: "logs", Source: source, Time: epoch + int64(seconds)*int64(time.Second), Receipt: epoch + int64(seconds)*int64(time.Second), Scope: "visit-a", Tab: "tab-a", Record: fmt.Sprint(id), Definition: 1, Attrs: map[string]string{}}
}
func TestJourneyHealthyLoopRecoveryAndBoundaries(t *testing.T) {
	es := []Evidence{ev(1, "discovery.entered", "browser", 0)}
	for i := 0; i < 10; i++ {
		es = append(es, ev(2+i, "vote.committed", "server", i+1))
	}
	bad := ev(12, "publication.failed", "server", 12)
	bad.Severity = 17
	bad.Attrs["exception.type"] = "PublicationError"
	es = append(es, bad, ev(13, "completion.entered", "browser", 14), ev(14, "discovery.entered", "browser", 15))
	m := Derive(es, []Definition{def()}, epoch+int64(20*time.Second))
	if len(m.Journeys) != 2 {
		t.Fatalf("goal must close first journey and new mapped event restart: %#v", m.Journeys)
	}
	j := m.Journeys[0]
	if j.Status != "completed" || !j.Recovered || len(j.Evidence) != 13 {
		t.Fatalf("error must remain without ending completion: %#v", j)
	}
	for _, edge := range m.Edges {
		if edge.From == "Vote" && edge.To == "Vote" {
			if edge.Occurrences != 9 || edge.Journeys != 1 || edge.Visits != 1 {
				t.Fatalf("healthy loop counts wrong: %#v", edge)
			}
			return
		}
	}
	t.Fatal("nine healthy Vote→Vote transitions missing")
}

func TestJourneyInactivityLateCorrectionAndReplay(t *testing.T) {
	a := ev(1, "discovery.entered", "browser", 0)
	m := Derive([]Evidence{a}, []Definition{def()}, epoch+int64(31*time.Minute))
	if len(m.Journeys) != 1 || m.Journeys[0].Status != "inactive" || m.Journeys[0].Reason != "" {
		t.Fatalf("silence must become inactive without invented reason: %#v", m.Journeys)
	}
	goal := ev(2, "completion.entered", "browser", 20)
	goal.Receipt = epoch + int64(time.Hour)
	m = Derive([]Evidence{a, goal, goal}, []Definition{def()}, epoch+int64(time.Hour))
	if len(m.Journeys) != 1 || m.Stats.Completed != 1 || len(m.Journeys[0].Evidence) != 2 {
		t.Fatalf("late actual goal and supported replay should revise one journey: %#v", m)
	}
}
func TestJourneyAuthorityTabIsolationAndServerOnly(t *testing.T) {
	browserCommit := ev(1, "vote.committed", "browser", 0)
	b := ev(2, "discovery.entered", "browser", 1)
	b.Tab = "tab-b"
	c := ev(3, "completion.entered", "browser", 2)
	c.Tab = "tab-c"
	noTab := ev(4, "discovery.entered", "browser", 3)
	noTab.Tab = ""
	server := ev(5, "vote.committed", "server", 4)
	server.Tab = ""
	server.Scope = "operation-a"
	m := Derive([]Evidence{browserCommit, b, c, noTab, server}, []Definition{def()}, epoch+int64(10*time.Second))
	if len(m.Journeys) != 3 || m.Stats.Completed != 1 || len(m.Unassigned) != 2 {
		t.Fatalf("authority/tab/server-only scopes invalid: %#v", m)
	}
	for _, j := range m.Journeys {
		if j.Scope == "operation-a" && j.Tab != "" {
			t.Fatal("fabricated server tab")
		}
	}
}
func TestJourneyOngoingDefinitionPinnedAndNextUsesNewVersion(t *testing.T) {
	d2 := def()
	d2.Version = 2
	d2.Goal = "vote.committed"
	a := ev(1, "discovery.entered", "browser", 0)
	b := ev(2, "vote.committed", "server", 1)
	b.Definition = 2
	c := ev(3, "completion.entered", "browser", 2)
	c.Definition = 2
	next := ev(4, "vote.committed", "server", 3)
	next.Definition = 2
	m := Derive([]Evidence{a, b, c, next}, []Definition{def(), d2}, epoch+int64(time.Minute))
	if len(m.Journeys) != 2 || m.Journeys[0].Definition != 1 || m.Journeys[0].Status != "completed" || len(m.Journeys[0].Evidence) != 3 || m.Journeys[1].Definition != 2 || m.Journeys[1].Status != "completed" {
		t.Fatalf("version activation changed captured boundary: %#v", m.Journeys)
	}
}
func TestJourneyDiagnosticsIdentityNoProximityAndSourceAuthority(t *testing.T) {
	a := ev(1, "discovery.entered", "browser", 0)
	bad := ev(2, "error", "server", 1)
	bad.Severity = 17
	bad.Attrs["exception.type"] = "SameError"
	bad.Attrs["exception.occurrence.id"] = "failure-a"
	secondCarrier := bad
	secondCarrier.ID = 3
	secondCarrier.Record = "3"
	secondCarrier.Signal = "span-events"
	secondFailure := bad
	secondFailure.ID = 4
	secondFailure.Record = "4"
	secondFailure.Attrs = map[string]string{"exception.type": "SameError", "exception.occurrence.id": "failure-b"}
	browserFailure := bad
	browserFailure.ID = 5
	browserFailure.Source = "browser"
	browserFailure.Record = "5"
	nearby := ev(6, "error", "server", 2)
	nearby.Scope = ""
	nearby.Severity = 17
	goal := ev(7, "completion.entered", "browser", 3)
	m := Derive([]Evidence{a, bad, secondCarrier, secondFailure, browserFailure, nearby, goal}, []Definition{def()}, epoch+int64(time.Minute))
	if m.Stats.DiagnosticRecords != 5 || m.Stats.QualifiedExceptionOccurrences != 3 || m.Stats.AffectedVisits != 1 {
		t.Fatalf("records/occurrences/source authority counts wrong: %#v", m.Stats)
	}
	for _, d := range m.Diagnostics {
		if d.EvidenceID == 6 && d.JourneyID != "" {
			t.Fatal("nearby unlinked error was attached")
		}
	}
}
func TestJourneyExplicitStageTimingActivityAndMissingEnd(t *testing.T) {
	es := []Evidence{ev(1, "discovery.entered", "browser", 0)}
	active := ev(2, "activity.changed", "browser", 1)
	active.Attrs["activity.state"] = "active"
	start := ev(3, "stage.entered", "browser", 2)
	start.Occurrence = "stage-a"
	hidden := ev(4, "activity.changed", "browser", 7)
	hidden.Attrs["activity.state"] = "hidden"
	end := ev(5, "stage.exited", "browser", 12)
	end.Occurrence = "stage-a"
	missing := ev(6, "stage.entered", "browser", 13)
	missing.Occurrence = "stage-b"
	es = append(es, active, start, hidden, end, missing)
	m := Derive(es, []Definition{def()}, epoch+int64(20*time.Second))
	var measured, observed, ongoing bool
	for _, tm := range m.Journeys[0].Timings {
		switch {
		case tm.Occurrence == "stage-a" && tm.Basis == "measured-elapsed":
			measured = tm.Complete && tm.Milliseconds == 10000
		case tm.Occurrence == "stage-a" && tm.Basis == "observed-active":
			observed = tm.Complete && tm.Milliseconds == 5000
		case tm.Occurrence == "stage-b":
			ongoing = !tm.Complete && tm.Provisional
		}
	}
	if !measured || !observed || !ongoing {
		t.Fatalf("elapsed/active/missing evidence mixed: %#v", m.Journeys[0].Timings)
	}
}
func TestJourneyTimingKindsSupportAndNoNegativeClock(t *testing.T) {
	ts := []Timing{}
	for i := 1; i <= 20; i++ {
		ts = append(ts, Timing{Basis: "measured-elapsed", Complete: true, Milliseconds: float64(i)})
	}
	ts = append(ts, Timing{Basis: "measured-elapsed", Provisional: true}, Timing{Basis: "next-event-estimated", Complete: true, Milliseconds: 999})
	sums := timingSummaries(ts)
	for _, s := range sums {
		if s.Basis == "measured-elapsed" && (s.Samples != 20 || s.Excluded != 1 || s.MedianMs == nil || *s.MedianMs != 10.5 || s.P90Ms == nil || *s.P90Ms != 18) {
			t.Fatalf("wrong supported distribution: %#v", s)
		}
		if s.Basis == "next-event-estimated" && s.P90Ms != nil {
			t.Fatal("low support p90 shown")
		}
	}
	a := ev(1, "stage.entered", "browser", 10)
	a.Occurrence = "a"
	a.Sequence = 1
	b := ev(2, "stage.exited", "browser", 5)
	b.Occurrence = "a"
	b.Sequence = 2
	m := Derive([]Evidence{a, b}, []Definition{def()}, epoch+int64(time.Minute))
	for _, tm := range m.Journeys[0].Timings {
		if tm.Complete || tm.Milliseconds < 0 {
			t.Fatalf("inconsistent clock became a measured duration: %#v", tm)
		}
	}
}

func TestJourneySummarizePreservesBoundariesAndExcludesPartialConversion(t *testing.T) {
	full := Derive([]Evidence{ev(1, "discovery.entered", "browser", 0), ev(2, "vote.committed", "server", 1), ev(3, "completion.entered", "browser", 2), ev(4, "discovery.entered", "browser", 3)}, []Definition{def()}, epoch+int64(time.Minute))
	chosen := full.Journeys[0]
	chosen.Partial = true
	m := Summarize([]Journey{chosen}, nil, epoch+int64(time.Minute))
	if len(m.Journeys) != 1 || m.Journeys[0].ID != chosen.ID || len(m.Journeys[0].Path) != 3 {
		t.Fatal("cohort selection changed captured path")
	}
	if m.Stats.Completed != 1 || m.Stats.EligibleCompleted != 0 || m.Stats.EligibleJourneys != 0 || m.Stats.PartialJourneys != 1 {
		t.Fatalf("partial goal entered started conversion denominator: %#v", m.Stats)
	}
}
func TestJourneyTimingPolicyVersionsStaySeparate(t *testing.T) {
	ss := timingSummaries([]Timing{{Definition: 1, Basis: "measured-elapsed", Complete: true, Milliseconds: 10}, {Definition: 2, Basis: "measured-elapsed", Complete: true, Milliseconds: 1000}})
	if len(ss) != 2 {
		t.Fatalf("unlike definition policies were combined: %#v", ss)
	}
}
func TestJourneyExpiredStartMarkerAndInitialGoalDiffer(t *testing.T) {
	goal := ev(1, "completion.entered", "browser", 0)
	initial := Derive([]Evidence{goal}, []Definition{def()}, epoch+int64(time.Minute))
	if initial.Stats.EligibleCompleted != 1 || initial.Journeys[0].Partial {
		t.Fatal("initial convergence was incorrectly marked partial")
	}
	goal.Attrs["_partial_start"] = "true"
	retained := Derive([]Evidence{goal}, []Definition{def()}, epoch+int64(time.Minute))
	if retained.Stats.Completed != 1 || retained.Stats.EligibleCompleted != 0 || !retained.Journeys[0].Partial {
		t.Fatal("expired original beginning was counted as a full started journey")
	}
}
func TestJourneyDuplicateSequenceKeepsOrderingUncertainty(t *testing.T) {
	a := ev(1, "discovery.entered", "browser", 0)
	a.Sequence = 1
	b := ev(2, "completion.entered", "browser", 1)
	b.Sequence = 1
	m := Derive([]Evidence{a, b}, []Definition{def()}, epoch+int64(time.Minute))
	if len(m.Journeys[0].Warnings) == 0 {
		t.Fatal("duplicate source sequence was advertised as reliable order")
	}
}
func TestJourneyStageCannotPairDifferentServiceClocks(t *testing.T) {
	a := ev(1, "stage.entered", "server", 0)
	a.Occurrence = "stage-a"
	a.Service = "service-a"
	b := ev(2, "stage.exited", "server", 10)
	b.Occurrence = "stage-a"
	b.Service = "service-b"
	m := Derive([]Evidence{a, b}, []Definition{def()}, epoch+int64(time.Minute))
	for _, tm := range m.Journeys[0].Timings {
		if tm.Complete {
			t.Fatalf("different service clocks fabricated stage duration: %#v", tm)
		}
	}
}
func TestJourneyExplicitExitReasonAndNoScopeJoin(t *testing.T) {
	a := ev(1, "discovery.entered", "browser", 0)
	exit := ev(2, "journey.exited", "browser", 1)
	exit.Attrs["reason"] = "exhausted"
	other := ev(3, "completion.entered", "browser", 2)
	other.Scope = "different-visit"
	m := Derive([]Evidence{a, exit, other}, []Definition{def()}, epoch+int64(time.Minute))
	if m.Stats.Exited != 1 || m.Stats.Completed != 1 {
		t.Fatalf("scope or exit boundary invalid: %#v", m.Stats)
	}
	for _, j := range m.Journeys {
		if j.Scope == "visit-a" && (j.Status != "exited" || j.Reason != "exhausted") {
			t.Fatalf("invented convergence for departed scope: %#v", j)
		}
	}
}
func TestJourneyWindowPreservesBoundaryAndNeverInventsEdge(t *testing.T) {
	d := def()
	a := ev(1, "discovery.entered", "browser", 0)
	b := ev(2, "vote.committed", "server", 1)
	c := ev(3, "completion.entered", "browser", 2)
	whole := Derive([]Evidence{a, b, c}, []Definition{d}, epoch+int64(time.Minute)).Journeys[0]
	clipped := ClipJourney(whole, epoch+int64(time.Second), epoch+int64(3*time.Second))
	if clipped.Start != whole.Start || clipped.ID != whole.ID || !clipped.Partial || len(clipped.Evidence) != 2 || len(clipped.Path) != 2 {
		t.Fatalf("window resegmented or included prior evidence: %#v", clipped)
	}
	// A clock-order conflict can clip an actual middle step; never count direct A→C.
	whole.Steps = []StepOccurrence{{Step: "A", Time: epoch + int64(2*time.Second)}, {Step: "X", Time: epoch}, {Step: "B", Time: epoch + int64(3*time.Second)}}
	whole.Path = []string{"A", "X", "B"}
	clipped = ClipJourney(whole, epoch+int64(time.Second), epoch+int64(4*time.Second))
	m := Summarize([]Journey{clipped}, nil, epoch+int64(time.Minute))
	if len(m.Edges) != 0 {
		t.Fatalf("clipping hidden X manufactured A→B edge: %#v", m.Edges)
	}
	prior := ClipJourney(Derive([]Evidence{a, b, c}, []Definition{d}, epoch+int64(time.Minute)).Journeys[0], epoch, epoch+int64(2*time.Second))
	if prior.Status == "completed" {
		t.Fatal("future completion polluted fixed baseline")
	}
}

func TestJourneyEqualEdgesRetainDifferentOrderedPaths(t *testing.T) {
	now := time.Now()
	d := Definition{Version: 1, Goal: "Goal", Mappings: []Mapping{{Event: "A", Step: "A", Kind: "step"}, {Event: "B", Step: "B", Kind: "step"}, {Event: "C", Step: "C", Kind: "step"}, {Event: "Goal", Step: "Goal", Kind: "step"}}}
	derivePath := func(path []string) Model {
		es := []Evidence{}
		for i, name := range path {
			es = append(es, Evidence{ID: int64(i + 1), Source: "server", Signal: "logs", Name: name, Scope: "logical", Record: fmt.Sprint(i), Time: now.Add(time.Duration(i) * time.Second).UnixNano(), Definition: 1, Attrs: map[string]string{}})
		}
		return Derive(es, []Definition{d}, now.Add(time.Minute).UnixNano())
	}
	a := derivePath([]string{"A", "B", "A", "C", "A"})
	b := derivePath([]string{"A", "C", "A", "B", "A"})
	if !reflect.DeepEqual(a.Edges, b.Edges) {
		t.Fatal("fixture does not have equal pairwise edges")
	}
	if reflect.DeepEqual(a.Journeys[0].Path, b.Journeys[0].Path) {
		t.Fatal("pairwise counters replaced ordered paths")
	}
}
