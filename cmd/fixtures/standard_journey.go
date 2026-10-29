package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/plog/plogotlp"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.opentelemetry.io/collector/pdata/ptrace/ptraceotlp"
)

// A progressive wire story. These synthetic states do not establish real
// browser collection, or either language's exporter behavior.
func buildStandardJourney(now time.Time) (map[string][]byte, string, string, error) {
	bodies := map[string][]byte{}
	raw, err := os.ReadFile("rest/fixtures/standard-journey.json")
	if err != nil {
		return nil, "", "", err
	}
	var spec struct {
		BrowserService, ServerService string
		BrowserEvents, ServerEvents   []Event
		Goal                          Event
		Operations                    []struct {
			Name, Start, End, Status, Event, EventOffset string
			HTTPStatus                                   int
		}
	}
	if err = json.Unmarshal(raw, &spec); err != nil {
		return nil, "", "", err
	}
	scope, tab, page := randomID(), randomID(), randomID()
	offset := func(value string) (pcommon.Timestamp, error) {
		d, err := time.ParseDuration(value)
		return pcommon.Timestamp(now.Add(d).UnixNano()), err
	}
	for name, events := range map[string][]Event{"standard-journey-browser": spec.BrowserEvents, "standard-journey-server": spec.ServerEvents, "standard-journey-goal": {spec.Goal}} {
		q := plogotlp.NewExportRequest()
		r := q.Logs().ResourceLogs().AppendEmpty()
		service := spec.BrowserService
		if name == "standard-journey-server" {
			service = spec.ServerService
		}
		r.Resource().Attributes().PutStr("service.name", service)
		ls := r.ScopeLogs().AppendEmpty().LogRecords()
		for _, event := range events {
			l := ls.AppendEmpty()
			l.SetEventName(event.Name)
			stamp, err := offset(event.Offset)
			if err != nil {
				return nil, "", "", err
			}
			l.SetTimestamp(stamp)
			l.SetSeverityNumber(plog.SeverityNumber(event.Severity))
			a := l.Attributes()
			a.PutStr("session.id", scope)
			a.PutStr("tab.id", tab)
			a.PutStr("event.id", randomID())
			if event.Occurrence != "" {
				a.PutStr("stage.occurrence.id", page)
			}
			for key, value := range event.Attributes {
				a.PutStr(key, value)
			}
		}
		b, err := q.MarshalProto()
		if err != nil {
			return nil, "", "", err
		}
		bodies[name+".bin"] = b
	}
	q := ptraceotlp.NewExportRequest()
	r := q.Traces().ResourceSpans().AppendEmpty()
	r.Resource().Attributes().PutStr("service.name", spec.ServerService)
	spans := r.ScopeSpans().AppendEmpty().Spans()
	firstTrace := ""
	for i, operation := range spec.Operations {
		p := spans.AppendEmpty()
		p.SetName(operation.Name)
		p.SetKind(ptrace.SpanKindServer)
		var trace pcommon.TraceID
		copy(trace[:], diagnosticBytesID(16))
		var span pcommon.SpanID
		copy(span[:], diagnosticBytesID(8))
		p.SetTraceID(trace)
		p.SetSpanID(span)
		if i == 0 {
			firstTrace = hex.EncodeToString(trace[:])
		}
		start, err := offset(operation.Start)
		if err != nil {
			return nil, "", "", err
		}
		end, err := offset(operation.End)
		if err != nil {
			return nil, "", "", err
		}
		p.SetStartTimestamp(start)
		p.SetEndTimestamp(end)
		p.Attributes().PutStr("session.id", scope)
		p.Attributes().PutStr("tab.id", tab)
		p.Attributes().PutStr("http.request.method", "POST")
		p.Attributes().PutStr("http.route", "/votes")
		p.Attributes().PutInt("http.response.status_code", int64(operation.HTTPStatus))
		if operation.Status == "error" {
			p.Status().SetCode(ptrace.StatusCodeError)
			p.Status().SetMessage("Vote persistence unavailable; retry permitted")
		} else {
			p.Status().SetCode(ptrace.StatusCodeOk)
		}
		event := p.Events().AppendEmpty()
		event.SetName(operation.Event)
		stamp, err := offset(operation.EventOffset)
		if err != nil {
			return nil, "", "", err
		}
		event.SetTimestamp(stamp)
		if operation.Event == "exception" {
			event.Attributes().PutStr("exception.type", "SyntheticPersistenceUnavailable")
		}
	}
	b, err := q.MarshalProto()
	if err != nil {
		return nil, "", "", err
	}
	bodies["standard-journey-traces.bin"] = b
	if firstTrace == "" {
		return nil, "", "", fmt.Errorf("journey fixture requires native operations")
	}
	return bodies, scope, firstTrace, nil
}

func checkStandardJourney(raw []byte, definition string, completed bool) error {
	var snap struct {
		Mode  string `json:"mode"`
		Model struct {
			Computation string `json:"computation"`
			Journeys    []struct {
				Definition int64    `json:"definition"`
				Status     string   `json:"status"`
				Recovered  bool     `json:"recovered"`
				Path       []string `json:"path"`
				Timings    []struct {
					Basis, Computation string
					Complete           bool
					Milliseconds       float64
					Start, End         string
					DurationNs         *string `json:"durationNs"`
				}
				Evidence []struct{ Name, Signal, Source string }
			} `json:"journeys"`
		} `json:"model"`
	}
	if err := json.Unmarshal(raw, &snap); err != nil {
		return err
	}
	if snap.Mode != "journey" || snap.Model.Computation != "journey-native-v3" || len(snap.Model.Journeys) != 1 {
		return fmt.Errorf("scoped voting requires one explicit journey under native computation, got%d", len(snap.Model.Journeys))
	}
	j := snap.Model.Journeys[0]
	if fmt.Sprint(j.Definition) != definition {
		return fmt.Errorf("goal/new definition remapped active capture: %d expected%s", j.Definition, definition)
	}
	if completed {
		if j.Status != "completed" || !j.Recovered {
			return fmt.Errorf("actual goal did not preserve completed recovery")
		}
	} else if j.Status != "active" {
		return fmt.Errorf("HTTP success or commitment incorrectly completed browser-screen goal")
	}
	votes, requests, processing := 0, 0, 0
	for _, e := range j.Evidence {
		if e.Name == "vote.committed" && e.Source == "server" && e.Signal == "span-events" {
			votes++
		}
		if e.Name == "vote.requested" && e.Source == "browser" {
			requests++
		}
	}
	for _, step := range j.Path {
		if step == "Vote processing" {
			processing++
		}
		if step == "Next saved vote" {
			return fmt.Errorf("new definition relabeled existing capture")
		}
	}
	if votes != 1 || requests != 2 || processing != 2 {
		return fmt.Errorf("voting story lost distinct attempts/commit/native occurrences: %d/%d/%d", requests, votes, processing)
	}
	native, elapsed, active := 0, 0, 0
	for _, timing := range j.Timings {
		if !timing.Complete {
			continue
		}
		if timing.Computation == "" || timing.Start == "" || timing.End == "" || timing.DurationNs == nil {
			return fmt.Errorf("timing lost exact time or policy")
		}
		switch timing.Basis {
		case "native-inclusive":
			if timing.Milliseconds != 200 && timing.Milliseconds != 300 {
				return fmt.Errorf("native attempt duration changed")
			}
			native++
		case "measured-elapsed":
			if timing.Milliseconds != 1000 {
				return fmt.Errorf("page pair elapsed changed")
			}
			elapsed++
		case "observed-active":
			if timing.Milliseconds != 800 {
				return fmt.Errorf("synthetic presence interval changed")
			}
			active++
		}
	}
	if native != 2 || elapsed != 1 || active != 1 {
		return fmt.Errorf("timing bases missing or double-pooled: native%d pair%d presence%d", native, elapsed, active)
	}
	return nil
}
