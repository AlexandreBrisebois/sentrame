package main

import (
	"go.opentelemetry.io/collector/pdata/plog/plogotlp"
	"go.opentelemetry.io/collector/pdata/ptrace/ptraceotlp"
	"testing"
	"time"
)

func TestStandardJourneyFixtureUsesExplicitBoundedScopeAndNativeAttempts(t *testing.T) {
	t.Chdir("../..")
	bodies, scope, trace, err := buildStandardJourney(time.Now().Add(-time.Minute))
	if err != nil {
		t.Fatal(err)
	}
	if len(scope) != 32 || len(trace) != 32 || len(bodies) != 4 {
		t.Fatal("fixture lost fresh scope/native trace identity")
	}
	q := ptraceotlp.NewExportRequest()
	if err = q.UnmarshalProto(bodies["standard-journey-traces.bin"]); err != nil {
		t.Fatal(err)
	}
	sp := q.Traces().ResourceSpans().At(0).ScopeSpans().At(0).Spans()
	if sp.Len() != 2 || sp.At(0).TraceID() == sp.At(1).TraceID() {
		t.Fatal("separate attempts collapsed into one request trace")
	}
	for i := 0; i < sp.Len(); i++ {
		a := sp.At(i).Attributes()
		v, ok := a.Get("session.id")
		if !ok || v.Str() != scope {
			t.Fatal("native operation lacks explicit journey scope")
		}
		if _, ok = a.Get("stage.occurrence.id"); ok {
			t.Fatal("native operation needlessly requires custom occurrence identity")
		}
	}
	logs := plogotlp.NewExportRequest()
	if err = logs.UnmarshalProto(bodies["standard-journey-goal.bin"]); err != nil {
		t.Fatal(err)
	}
	l := logs.Logs().ResourceLogs().At(0).ScopeLogs().At(0).LogRecords().At(0)
	if l.EventName() != "completion.entered" {
		t.Fatal("server commit invented a browser-screen goal")
	}
}
