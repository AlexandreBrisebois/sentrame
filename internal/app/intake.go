package app

import (
	"context"
	"encoding/hex"
	"fmt"
	"go.opentelemetry.io/collector/client"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var tokenPattern = regexp.MustCompile(`^[a-zA-Z0-9_.:/ -]{1,160}$`)
var allowed = map[string]bool{"stage.name": true, "stage.state": true, "activity.state": true, "reason": true, "exception.type": true, "exception.occurrence.id": true, "telemetry.sdk.language": true, "process.runtime.name": true, "deployment.environment.name": true}

func safe(v string) string {
	if tokenPattern.MatchString(v) {
		return v
	}
	return ""
}
func attr(m pcommon.Map, k string) string {
	v, ok := m.Get(k)
	if !ok {
		return ""
	}
	return safe(v.AsString())
}
func (s *Service) source(ctx context.Context) (string, error) {
	hs := client.FromContext(ctx).Metadata.Get("authorization")
	if len(hs) != 1 {
		return "", status.Error(codes.Unauthenticated, "producer credential required")
	}
	return s.producer(hs[0])
}
func (s *Service) producer(h string) (string, error) {
	if s.BrowserToken != "" && h == "Bearer "+s.BrowserToken {
		return "browser", nil
	}
	if s.ServerToken != "" && h == "Bearer "+s.ServerToken {
		return "server", nil
	}
	return "", status.Error(codes.Unauthenticated, "invalid producer credential")
}
func (s *Service) project(source, signal, name string, t pcommon.Timestamp, observed pcommon.Timestamp, a, r pcommon.Map) Evidence {
	now := s.Store.Now().UnixNano()
	x := Evidence{Source: source, Signal: signal, Name: safe(name), Time: int64(t), Observed: int64(observed), Receipt: now, Attrs: map[string]string{}, Service: attr(r, "service.name"), Version: attr(r, "service.version"), Scope: attr(a, "session.id"), Tab: attr(a, "tab.id"), Record: attr(a, "event.id"), Action: attr(a, "action.id"), Occurrence: attr(a, "stage.occurrence.id")}
	if x.Name == "" {
		x.Name = "unnamed"
	}
	if x.Time == 0 {
		x.Time = now
		x.Fallback = true
	}
	if seq, err := strconv.ParseInt(attr(a, "source.sequence"), 10, 64); err == nil {
		x.Sequence = seq
	}
	a.Range(func(k string, v pcommon.Value) bool {
		if allowed[k] {
			if val := safe(v.AsString()); val != "" {
				x.Attrs[k] = val
			}
		}
		return true
	})
	r.Range(func(k string, v pcommon.Value) bool {
		if allowed[k] {
			x.Attrs[k] = safe(v.AsString())
		}
		return true
	})
	if signal != "metrics" {
		projectDiagnosticAttributes(&x, a)
	}
	return x
}
func (s *Service) accept(es []Evidence) error {
	if len(es) > 2000 {
		return status.Error(codes.InvalidArgument, "record limit")
	}
	now := s.Store.Now()
	for _, x := range es {
		if x.Time < now.Add(-30*24*time.Hour).UnixNano() || x.Time > now.Add(5*time.Minute).UnixNano() {
			return status.Error(codes.InvalidArgument, "timestamp outside supported horizon")
		}
	}
	if e := s.Store.Insert(es); e != nil {
		if strings.Contains(e.Error(), "disabled") {
			return status.Error(codes.PermissionDenied, "application generation disabled")
		}
		if strings.Contains(e.Error(), "conflict") {
			return status.Error(codes.InvalidArgument, "identity conflict")
		}
		return status.Error(codes.Unavailable, "storage unavailable or bounded capacity reached")
	}
	return nil
}
func (s *Service) consumeLogs(ctx context.Context, ls plog.Logs) error {
	source, e := s.source(ctx)
	if e != nil {
		return e
	}
	es := []Evidence{}
	rs := ls.ResourceLogs()
	for i := 0; i < rs.Len(); i++ {
		r := rs.At(i)
		if r.Resource().Attributes().Len() > 64 {
			return status.Error(codes.InvalidArgument, "resource attribute limit")
		}
		if r.Resource().Attributes().Len() > 64 {
			return status.Error(codes.InvalidArgument, "attribute limit")
		}
		ss := r.ScopeLogs()
		for j := 0; j < ss.Len(); j++ {
			if ss.At(j).Scope().Attributes().Len() > 64 {
				return status.Error(codes.InvalidArgument, "scope attribute limit")
			}
			logs := ss.At(j).LogRecords()
			for k := 0; k < logs.Len(); k++ {
				l := logs.At(k)
				if l.Attributes().Len() > 64 {
					return status.Error(codes.InvalidArgument, "attribute limit")
				}
				name := l.EventName()
				if name == "" {
					name = attr(l.Attributes(), "event.name")
				}
				x := s.project(source, "logs", name, l.Timestamp(), l.ObservedTimestamp(), l.Attributes(), r.Resource().Attributes())
				x.Severity = int32(l.SeverityNumber())
				if x.Severity >= 17 || x.Attrs["exception.type"] != "" {
					x.Attrs["diagnostic.policy"] = diagnosticPolicy
					if l.Body().Type() == pcommon.ValueTypeStr {
						retainDiagnosticText(&x, "log.message", l.Body().Str(), false)
					}
				}
				x.Trace, x.Span = sourceIDs(l.TraceID(), l.SpanID())
				es = append(es, x)
			}
		}
	}
	return s.accept(es)
}
func (s *Service) consumeTraces(ctx context.Context, ts ptrace.Traces) error {
	source, e := s.source(ctx)
	if e != nil {
		return e
	}
	es := []Evidence{}
	rs := ts.ResourceSpans()
	for i := 0; i < rs.Len(); i++ {
		r := rs.At(i)
		if r.Resource().Attributes().Len() > 64 {
			return status.Error(codes.InvalidArgument, "resource attribute limit")
		}
		ss := r.ScopeSpans()
		for j := 0; j < ss.Len(); j++ {
			if ss.At(j).Scope().Attributes().Len() > 64 {
				return status.Error(codes.InvalidArgument, "scope attribute limit")
			}
			sp := ss.At(j).Spans()
			for k := 0; k < sp.Len(); k++ {
				p := sp.At(k)
				if p.TraceID().IsEmpty() || p.SpanID().IsEmpty() {
					return status.Error(codes.InvalidArgument, "standard span requires valid trace and span IDs")
				}
				if p.Attributes().Len() > 64 {
					return status.Error(codes.InvalidArgument, "attribute limit")
				}
				x := s.project(source, "traces", p.Name(), p.StartTimestamp(), 0, p.Attributes(), r.Resource().Attributes())
				x.Trace, x.Span = sourceIDs(p.TraceID(), p.SpanID())
				x.Record = x.Trace + "/" + x.Span
				if p.Status().Code() == ptrace.StatusCodeError {
					x.Severity = 17
					x.Attrs["diagnostic.policy"] = diagnosticPolicy
					if x.Attrs["exception.message"] == "" {
						retainDiagnosticText(&x, "exception.message", p.Status().Message(), false)
					}
				}
				es = append(es, x)
				ev := p.Events()
				for q := 0; q < ev.Len(); q++ {
					v := ev.At(q)
					if v.Attributes().Len() > 64 {
						return status.Error(codes.InvalidArgument, "event attribute limit")
					}
					z := s.project(source, "span-events", v.Name(), v.Timestamp(), 0, v.Attributes(), r.Resource().Attributes())
					z.Trace = x.Trace
					z.Span = x.Span
					if z.Scope == "" {
						z.Scope = x.Scope
						if z.Tab == "" {
							z.Tab = x.Tab
						}
					}
					if z.Scope == x.Scope && z.Tab == x.Tab && z.Action == "" {
						z.Action = x.Action
					}
					z.Record = fmt.Sprintf("%s/%d", x.Record, q)
					es = append(es, z)
				}
			}
		}
	}
	return s.accept(es)
}
func (s *Service) consumeMetrics(ctx context.Context, ms pmetric.Metrics) error {
	source, e := s.source(ctx)
	if e != nil {
		return e
	}
	es := []Evidence{}
	rs := ms.ResourceMetrics()
	for i := 0; i < rs.Len(); i++ {
		r := rs.At(i)
		if r.Resource().Attributes().Len() > 64 {
			return status.Error(codes.InvalidArgument, "resource attribute limit")
		}
		ss := r.ScopeMetrics()
		for j := 0; j < ss.Len(); j++ {
			if ss.At(j).Scope().Attributes().Len() > 64 {
				return status.Error(codes.InvalidArgument, "scope attribute limit")
			}
			mm := ss.At(j).Metrics()
			for k := 0; k < mm.Len(); k++ {
				m := mm.At(k)
				pointOverflow := false
				add := func(t pcommon.Timestamp, a pcommon.Map) {
					if a.Len() > 64 {
						pointOverflow = true
						return
					}
					x := s.project(source, "metrics", m.Name(), t, 0, a, r.Resource().Attributes())
					x.Scope = ""
					x.Tab = ""
					x.Record = ""
					es = append(es, x)
				}
				switch m.Type() {
				case pmetric.MetricTypeGauge:
					ps := m.Gauge().DataPoints()
					for n := 0; n < ps.Len(); n++ {
						p := ps.At(n)
						add(p.Timestamp(), p.Attributes())
					}
				case pmetric.MetricTypeSum:
					ps := m.Sum().DataPoints()
					for n := 0; n < ps.Len(); n++ {
						p := ps.At(n)
						add(p.Timestamp(), p.Attributes())
					}
				case pmetric.MetricTypeHistogram:
					ps := m.Histogram().DataPoints()
					for n := 0; n < ps.Len(); n++ {
						p := ps.At(n)
						if p.BucketCounts().Len() > 2048 {
							return status.Error(codes.InvalidArgument, "bucket limit")
						}
						add(p.Timestamp(), p.Attributes())
					}
				case pmetric.MetricTypeExponentialHistogram:
					ps := m.ExponentialHistogram().DataPoints()
					for n := 0; n < ps.Len(); n++ {
						p := ps.At(n)
						if p.Positive().BucketCounts().Len()+p.Negative().BucketCounts().Len() > 2048 {
							return status.Error(codes.InvalidArgument, "bucket limit")
						}
						add(p.Timestamp(), p.Attributes())
					}
				case pmetric.MetricTypeSummary:
					ps := m.Summary().DataPoints()
					for n := 0; n < ps.Len(); n++ {
						p := ps.At(n)
						add(p.Timestamp(), p.Attributes())
					}
				}
				if pointOverflow {
					return status.Error(codes.InvalidArgument, "point attribute limit")
				}
			}
		}
	}
	return s.accept(es)
}

func sourceIDs(t pcommon.TraceID, s pcommon.SpanID) (string, string) {
	trace, span := "", ""
	if !t.IsEmpty() {
		trace = hex.EncodeToString(t[:])
	}
	if !s.IsEmpty() {
		span = hex.EncodeToString(s[:])
	}
	return trace, span
}
