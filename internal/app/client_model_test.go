package app

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestClientModelDoesNotBridgeLegacyMixedEvidence(t *testing.T) {
	a := ev(1, "discovery.entered", "client", 0)
	b := ev(2, "vote.committed", "server", 1)
	c := ev(3, "completion.entered", "client", 2)
	model := Derive([]Evidence{a, b, c}, []Definition{def()}, epoch+int64(time.Minute))
	if len(model.Journeys) != 0 || len(model.Edges) != 0 || model.Stats.Records != 0 {
		t.Fatalf("legacy B removal fabricated a client path/count: %+v", model)
	}
	if !strings.Contains(strings.Join(model.Warnings, " "), "non-comparable") {
		t.Fatal("incompatible retained history must explain its absence")
	}
}

func TestClientModelCannotUseRoleRulesOrFallBackToNavigation(t *testing.T) {
	d := Definition{Version: 1, Mappings: []Mapping{{Event: "screen", Step: "Saved", Kind: "step", Source: "server"}}}
	if msg := validateDefinition(d); !strings.Contains(msg, "unsupported") {
		t.Fatalf("retired role was not explicitly rejected: %q", msg)
	}
	a := ev(1, "screen", "client", 0)
	model := Derive([]Evidence{a}, []Definition{d}, epoch+int64(time.Minute))
	if len(model.Journeys) != 0 || model.Stats.Records != 0 {
		t.Fatal("unsafe historic definition fell back to navigation")
	}
	d.Goal = "screen"
	d.GoalSelector = &Mapping{Event: "screen", Step: "Saved", Kind: "step", Source: "browser"}
	d.Mappings[0].Source = ""
	if msg := validateDefinition(d); !strings.Contains(msg, "unsupported") {
		t.Fatal("retired goal selector was accepted")
	}
}

func TestCompatibleClientDefinitionPreservesOriginalAndMeaning(t *testing.T) {
	d := Definition{Version: 7, Goal: "screen", Mappings: []Mapping{{Event: "screen", Step: "Screen shown", Kind: "step", Source: "browser"}}}
	d.GoalSelector = &d.Mappings[0]
	before := d.Mappings[0]
	converted, ok := compatibleClientDefinition(d)
	if !ok || converted.Mappings[0].Source != "" || converted.GoalSelector.Source != "" {
		t.Fatal("compatible browser-only definition was not projected to client semantics")
	}
	if !reflect.DeepEqual(d.Mappings[0], before) || d.GoalSelector.Source != "browser" {
		t.Fatal("analysis projection mutated persisted definition content")
	}
	d.Mappings[0].Source = "server"
	if _, ok := compatibleClientDefinition(d); ok {
		t.Fatal("server goal acquired client meaning")
	}
}

func TestClientJourneyRequiresTabButTraceReceiptDoesNot(t *testing.T) {
	a := ev(1, "discovery.entered", "client", 0)
	a.Tab = ""
	a.Signal = "traces"
	a.Trace = "11111111111111111111111111111111"
	a.Span = "0000000000000001"
	a.End = a.Time + int64(time.Second)
	a.Kind = 2 // SERVER is protocol metadata, not runtime assurance.
	model := Derive([]Evidence{a}, []Definition{def()}, epoch+int64(time.Minute))
	if len(model.Journeys) != 0 || len(model.Unassigned) != 1 {
		t.Fatal("missing client tab created a server-workflow fallback")
	}
	detail := traceDetail([]Evidence{a}, a.Trace, "1", false)
	if len(detail.Operations) != 1 || detail.Operations[0].DurationNs == nil {
		t.Fatal("otherwise valid standard span lost diagnostic structure")
	}
}
