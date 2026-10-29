package app

import (
	"context"
	"encoding/hex"
	"encoding/json"
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
var allowed = map[string]bool{"http.request.method": true, "http.route": true, "http.response.status_code": true, "stage.name": true, "stage.state": true, "activity.state": true, "reason": true, "exception.type": true, "exception.occurrence.id": true, "telemetry.sdk.language": true, "process.runtime.name": true, "deployment.environment.name": true}

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
	if k != "source.sequence" && v.Type() != pcommon.ValueTypeStr {
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
	if s.Registry != nil && s.selected == nil {
		s.Registry.mu.RLock()
		defer s.Registry.mu.RUnlock()
		_, role, ok := s.Registry.credential(h)
		if ok {
			return role, nil
		}
		return "", status.Error(codes.Unauthenticated, "invalid producer credential")
	}

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
	x := Evidence{Projection: projectionPolicy, Source: source, Signal: signal, Name: safeEvidenceName(name), Time: int64(t), Observed: int64(observed), Receipt: now, Attrs: map[string]string{}, Service: attr(r, "service.name"), Version: attr(r, "service.version"), Scope: attr(a, "session.id"), Tab: attr(a, "tab.id"), Record: attr(a, "event.id"), Action: attr(a, "action.id"), Occurrence: attr(a, "stage.occurrence.id")}
	x.NameState = "retained"
	if x.Name == "" {
		x.Name = "unnamed"
		x.NameState = "missing"
		if name != "" {
			x.NameState = "omitted"
		}
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
			if val := projectedAttribute(k, v); val != "" {
				x.Attrs[k] = val
			} else {
				x.OmittedAttributes++
			}
		} else if contextAttributes[k] {
			if attr(a, k) == "" {
				x.OmittedAttributes++
			}
		} else {
			x.OmittedAttributes++
		}
		return true
	})
	r.Range(func(k string, v pcommon.Value) bool {
		if k == "service.name" || k == "service.version" {
			if attr(r, k) == "" {
				x.OmittedAttributes++
			}
		} else if k == "telemetry.sdk.language" || k == "process.runtime.name" || k == "deployment.environment.name" {
			if val := projectedAttribute(k, v); val != "" {
				x.Attrs[k] = val
			} else {
				x.OmittedAttributes++
			}
		} else {
			x.OmittedAttributes++
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
	if len(s.credentials) > 0 {
		raw, _ := json.Marshal(es)
		for _, secret := range s.credentials {
			if secret != "" && strings.Contains(string(raw), secret) {
				return status.Error(codes.InvalidArgument, "ingestion credentials must not appear in telemetry")
			}
		}
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
	if s.Registry != nil && s.selected == nil {
		s.Registry.mu.RLock()
		defer s.Registry.mu.RUnlock()
		hs := client.FromContext(ctx).Metadata.Get("authorization")
		if len(hs) != 1 {
			return status.Error(codes.Unauthenticated, "producer credential required")
		}
		a, _, ok := s.Registry.credential(hs[0])
		if !ok {
			return status.Error(codes.Unauthenticated, "invalid producer credential")
		}
		return s.forApp(a).consumeLogs(ctx, ls)
	}

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
					if value, ok := l.Attributes().Get("event.name"); ok && value.Type() == pcommon.ValueTypeStr {
						name = value.Str()
					}
				}
				x := s.project(source, "logs", name, l.Timestamp(), l.ObservedTimestamp(), l.Attributes(), r.Resource().Attributes())
				structuralResourceCounts(&x, r.Resource(), ss.At(j).Scope())
				x.Severity = int32(l.SeverityNumber())
				x.Flags = uint32(l.Flags())
				x.DroppedAttributes = l.DroppedAttributesCount()
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
	if s.Registry != nil && s.selected == nil {
		s.Registry.mu.RLock()
		defer s.Registry.mu.RUnlock()
		hs := client.FromContext(ctx).Metadata.Get("authorization")
		if len(hs) != 1 {
			return status.Error(codes.Unauthenticated, "producer credential required")
		}
		a, _, ok := s.Registry.credential(hs[0])
		if !ok {
			return status.Error(codes.Unauthenticated, "invalid producer credential")
		}
		return s.forApp(a).consumeTraces(ctx, ts)
	}

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
				structuralResourceCounts(&x, r.Resource(), ss.At(j).Scope())
				x.Trace, x.Span = sourceIDs(p.TraceID(), p.SpanID())
				x.Record = x.Trace + "/" + x.Span
				x.End = int64(p.EndTimestamp())
				if uint64(p.EndTimestamp()) > uint64(1<<63-1) {
					return status.Error(codes.InvalidArgument, "unsupported span end timestamp")
				}
				_, x.Parent = sourceIDs(pcommon.TraceID{}, p.ParentSpanID())
				x.Kind = int32(p.Kind())
				x.Status = int32(p.Status().Code())
				x.Flags = p.Flags()
				x.EventCount = p.Events().Len()
				x.DroppedAttributes = p.DroppedAttributesCount()
				x.DroppedEvents = p.DroppedEventsCount()
				x.DroppedLinks = p.DroppedLinksCount()
				for n := 0; n < p.Links().Len(); n++ {
					if n >= 32 {
						x.OmittedLinks = p.Links().Len() - 32
						break
					}
					link := p.Links().At(n)
					trace, span := sourceIDs(link.TraceID(), link.SpanID())
					x.Links = append(x.Links, SpanLink{Trace: trace, Span: span})
				}

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
					structuralResourceCounts(&z, r.Resource(), ss.At(j).Scope())
					z.EventIndex = q
					z.DroppedAttributes = v.DroppedAttributesCount()
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
	if s.Registry != nil && s.selected == nil {
		s.Registry.mu.RLock()
		defer s.Registry.mu.RUnlock()
		hs := client.FromContext(ctx).Metadata.Get("authorization")
		if len(hs) != 1 {
			return status.Error(codes.Unauthenticated, "producer credential required")
		}
		a, _, ok := s.Registry.credential(hs[0])
		if !ok {
			return status.Error(codes.Unauthenticated, "invalid producer credential")
		}
		return s.forApp(a).consumeMetrics(ctx, ms)
	}

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
					structuralResourceCounts(&x, r.Resource(), ss.At(j).Scope())
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

// HTTP selectors are deliberately narrower than arbitrary safe tokens. A route
// must be a minimized template supplied by the producer, never a URL or query.
var routeLiteral = regexp.MustCompile(`^[a-z][a-z-]*$`)
var routeParameter = regexp.MustCompile(`^(\{[a-zA-Z][a-zA-Z0-9_]*\}|:[a-zA-Z][a-zA-Z0-9_]*|\[[a-zA-Z][a-zA-Z0-9_]*\])$`)

func projectedAttribute(key string, v pcommon.Value) string {
	switch key {
	case "stage.name":
		if v.Type() != pcommon.ValueTypeStr || strings.TrimSpace(v.Str()) != v.Str() || strings.TrimSpace(v.Str()) == "" {
			return ""
		}
		return safe(v.Str())
	case "http.request.method":
		switch v.Str() {
		case "GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS", "CONNECT", "TRACE":
			return v.Str()
		}
		return ""
	case "http.response.status_code":
		if v.Type() != pcommon.ValueTypeInt || v.Int() < 100 || v.Int() > 599 {
			return ""
		}
		return strconv.FormatInt(v.Int(), 10)
	case "http.route":
		if v.Type() != pcommon.ValueTypeStr {
			return ""
		}
		route := v.Str()
		if len(route) > 160 || !strings.HasPrefix(route, "/") {
			return ""
		}
		if route == "/" {
			return route
		}
		for _, part := range strings.Split(strings.TrimPrefix(route, "/"), "/") {
			if !routeLiteral.MatchString(part) && !routeParameter.MatchString(part) {
				return ""
			}
		}
		return route
	default:
		return safe(v.AsString())
	}
}

var contextAttributes = map[string]bool{"event.name": true, "session.id": true, "tab.id": true, "event.id": true, "action.id": true, "stage.occurrence.id": true, "source.sequence": true, "exception.message": true, "exception.stacktrace": true}

func structuralResourceCounts(e *Evidence, r pcommon.Resource, scope pcommon.InstrumentationScope) {
	e.ResourceDroppedAttributes = r.DroppedAttributesCount()
	e.ScopeDroppedAttributes = scope.DroppedAttributesCount()
	e.OmittedAttributes += scope.Attributes().Len()
}
func safeEvidenceName(name string) string {
	if len(name) == 0 || len(name) > 160 || strings.Contains(name, "://") {
		return ""
	}
	if strings.HasPrefix(name, "/") {
		return projectedAttribute("http.route", pcommon.NewValueStr(name))
	}
	if method, route, found := strings.Cut(name, " "); found && projectedAttribute("http.request.method", pcommon.NewValueStr(strings.ToUpper(method))) != "" {
		if projectedAttribute("http.route", pcommon.NewValueStr(route)) == route {
			return name
		}
		return ""
	}
	return safe(name)
}
