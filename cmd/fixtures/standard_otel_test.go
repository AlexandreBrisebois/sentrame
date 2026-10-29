package main

import (
	"bytes"
	tracesv1 "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	"google.golang.org/protobuf/proto"
	"os"
	"testing"
	"time"
)

func TestStandardFixtureUsesNativeStructureWithoutJourneyFields(t *testing.T) {
	raw, err := os.ReadFile("../../rest/fixtures/standard-otel.json")
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now().Add(-time.Minute)
	g, err := buildStandard(raw, now)
	if err != nil {
		t.Fatal(err)
	}
	var tr tracesv1.ExportTraceServiceRequest
	if err = proto.Unmarshal(g.Bodies["standard-traces.bin"], &tr); err != nil {
		t.Fatal(err)
	}
	spans := tr.ResourceSpans[0].ScopeSpans[0].Spans
	if len(spans[0].TraceId) != 16 || len(spans[0].SpanId) != 8 {
		t.Fatal("OTLP hex IDs decoded as generic Protobuf bytes")
	}
	if len(spans) != 5 || len(spans[0].Links) != 34 || spans[0].EndTimeUnixNano-spans[0].StartTimeUnixNano != 1000000001 || spans[3].EndTimeUnixNano != 0 {
		t.Fatal("fixture lost native structure, exact duration or absent end")
	}
	for _, key := range []string{"session.id", "stage.name", "stage.occurrence.id", "event.id", "activity.state"} {
		if bytes.Contains(raw, []byte(key)) {
			t.Fatalf("ordinary fixture requires custom context %s", key)
		}
	}
	if len(g.Hashes) != 6 {
		t.Fatal("exact wire hashes missing")
	}
}
