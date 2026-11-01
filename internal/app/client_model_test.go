package app

import (
	"testing"
	"time"
)

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
