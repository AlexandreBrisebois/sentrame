package app

import (
	"reflect"
	"testing"
)

func TestReferenceQueryPreservesCapturedBoundariesAndRepeatedEndpoints(t *testing.T) {
	j := Journey{ID: "captured", Start: 1, End: 4, Status: "completed", Evidence: []Evidence{{ID: 1, Name: "A", Time: 1}, {ID: 2, Name: "X", Time: 2}, {ID: 3, Name: "A", Time: 3}, {ID: 4, Name: "B", Time: 4}}, Steps: []StepOccurrence{{Step: "A", Event: "A", EvidenceID: 1}, {Step: "X", Event: "X", EvidenceID: 2}, {Step: "A", Event: "A", EvidenceID: 3}, {Step: "B", Event: "B", EvidenceID: 4}}, Path: []string{"A", "X", "A", "B"}}
	v := SavedView{Mode: "goal-backward", End: "X"}
	out, ok := referenceJourney(j, v)
	if !ok || len(out.Evidence) != 2 || len(out.Path) != 2 || out.ID != j.ID || out.Start != j.Start || out.End != j.End {
		t.Fatalf("reference did not project within captured boundary: %+v", out)
	}
	out, ok = referenceJourney(j, SavedView{Mode: "segment", Start: "A", End: "B", StartOccurrence: "last", EndOccurrence: "last"})
	if !ok || len(out.Evidence) != 2 || out.Evidence[0].ID != 3 || !out.Partial {
		t.Fatal("repeated start selection not preserved")
	}
	out, ok = referenceJourney(j, SavedView{Mode: "segment", Start: "A", End: "outside-this-boundary"})
	if !ok || !out.Partial || len(out.Warnings) == 0 || out.ID != j.ID {
		t.Fatal("cross-boundary endpoint silently invented")
	}
}

func TestOverviewReferenceKeepsCompleteCapturedProjection(t *testing.T) {
	j := Journey{ID: "captured", Start: 1, End: 4, Status: "completed", Partial: true, RightCensored: true, CapturedStatus: "completed", WindowFrom: 2, WindowTo: 5, Evidence: []Evidence{{ID: 2, Name: "A", Time: 2}, {ID: 3, Name: "B", Time: 3}}, Path: []string{"A", "B"}, Steps: []StepOccurrence{{Step: "A", EvidenceID: 2}, {Step: "B", EvidenceID: 3, GapBefore: true}}, Timings: []Timing{{Quality: "measured"}}, Warnings: []string{"existing boundary"}}
	out, ok := referenceJourney(j, SavedView{Mode: "overview"})
	if !ok || !reflect.DeepEqual(out, j) {
		t.Fatalf("overview altered recorded projection: ok=%v got=%+v", ok, out)
	}
}
