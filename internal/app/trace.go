package app

import (
	"encoding/hex"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

const nativeTimingPolicy = "native-inclusive-v1"

type TraceOperation struct {
	Evidence   Evidence `json:"evidence"`
	DurationNs *string  `json:"durationNs"`
	Timing     string   `json:"timing"`
}
type TraceRelationship struct {
	Type        string `json:"type"`
	From        int64  `json:"from"`
	To          int64  `json:"to,omitempty"`
	TargetTrace string `json:"targetTrace"`
	TargetSpan  string `json:"targetSpan"`
	State       string `json:"state"`
}
type TraceDetail struct {
	Suggestions    []StageSuggestion    `json:"suggestions"`
	Business       []BusinessOccurrence `json:"business"`
	Interpretation string               `json:"interpretation"`
	Mode           string               `json:"mode"`
	Trace          string               `json:"trace"`
	Revision       string               `json:"revision"`
	Computation    string               `json:"computation"`
	Unit           string               `json:"unit"`
	Completeness   string               `json:"completeness"`
	Truncated      bool                 `json:"truncated"`
	Operations     []TraceOperation     `json:"operations"`
	Events         []Evidence           `json:"events"`
	Logs           []Evidence           `json:"logs"`
	Relationships  []TraceRelationship  `json:"relationships"`
}

func nativeDuration(e Evidence) (*string, string) {
	if e.Fallback || e.Time <= 0 {
		return nil, "unavailable"
	}
	if e.End == 0 {
		return nil, "right-censored"
	}
	if e.End < e.Time {
		return nil, "unavailable"
	}
	ns := strconv.FormatInt(e.End-e.Time, 10)
	return &ns, "complete"
}
func validTraceID(value string) bool {
	if len(value) != 32 {
		return false
	}
	b, e := hex.DecodeString(value)
	if e != nil {
		return false
	}
	for _, v := range b {
		if v != 0 {
			return true
		}
	}
	return false
}

// All relationships describe observed context. They confer no source authority.
// The candidate index lives only for this bounded read, never beyond retention.
func traceDetail(es []Evidence, trace, revision string, truncated bool) TraceDetail {
	d := TraceDetail{Mode: "trace", Trace: trace, Revision: revision, Computation: nativeTimingPolicy, Unit: "observed spans", Completeness: "unknown: sampling and unexported work cannot establish completeness", Truncated: truncated, Operations: []TraceOperation{}, Events: []Evidence{}, Logs: []Evidence{}, Relationships: []TraceRelationship{}}
	candidates := map[string][]Evidence{}
	for _, e := range es {
		if e.Signal == "traces" {
			key := e.Trace + "/" + e.Span
			candidates[key] = append(candidates[key], e)
		}
	}
	relationship := func(kind string, e Evidence, targetTrace, targetSpan, source string) {
		r := TraceRelationship{Type: kind, From: e.ID, TargetTrace: targetTrace, TargetSpan: targetSpan, State: "missing"}
		matches := []Evidence{}
		for _, c := range candidates[targetTrace+"/"+targetSpan] {
			if source == "" || source == c.Source {
				matches = append(matches, c)
			}
		}
		if kind == "log-context" && targetSpan == "" {
			r.State = "trace-only"
		} else if len(matches) == 1 {
			r.State = "resolved"
			r.To = matches[0].ID
		} else if len(matches) > 1 {
			r.State = "ambiguous"
		}
		d.Relationships = append(d.Relationships, r)
	}
	count := 0
	for _, e := range es {
		if e.Trace != trace {
			continue
		}
		count++
		if count > 2000 {
			d.Truncated = true
			continue
		}
		switch e.Signal {
		case "traces":
			ns, timing := nativeDuration(e)
			d.Operations = append(d.Operations, TraceOperation{Evidence: e, DurationNs: ns, Timing: timing})
			if e.Parent != "" {
				relationship("parent", e, e.Trace, e.Parent, "")
			}
			for _, link := range e.Links {
				relationship("link", e, link.Trace, link.Span, "")
			}
		case "span-events":
			d.Events = append(d.Events, e)
			relationship("span-event", e, e.Trace, e.Span, e.Source)
		case "logs":
			d.Logs = append(d.Logs, e)
			relationship("log-context", e, e.Trace, e.Span, "")
		}
	}
	// Event order is defined inside its exported span, not by event timestamps.
	sort.SliceStable(d.Events, func(i, j int) bool {
		a, b := d.Events[i], d.Events[j]
		if a.Source != b.Source {
			return a.Source < b.Source
		}
		if a.Span != b.Span {
			return a.Span < b.Span
		}
		return a.EventIndex < b.EventIndex
	})
	return d
}
func (s *Service) traces(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "GET required", 405)
		return
	}
	trace := ""
	if r.URL.Path != "/api/traces" {
		trace = strings.TrimPrefix(r.URL.Path, "/api/traces/")
		if !validTraceID(trace) {
			http.Error(w, "Trace unavailable", 404)
			return
		}
	}
	s.Store.Mu.Lock()
	defer s.Store.Mu.Unlock()
	q := r.URL.Query()
	var reference *SavedView
	if q.Get("reference") != "" {
		if q.Has("previewDefinition") {
			http.Error(w, "Choose a reference or an explicit preview", 400)
			return
		}
		views, err := s.getViews()
		if err != nil {
			http.Error(w, "Store unavailable", 503)
			return
		}
		for i := range views {
			if strconv.FormatInt(views[i].ID, 10) == q.Get("reference") && views[i].Mode == "trace" {
				copy := views[i]
				reference = &copy
				break
			}
		}
		if reference == nil || trace != "" && trace != reference.Trace {
			http.Error(w, "Trace reference unavailable", 404)
			return
		}
		if message := validateReferenceQuery(*reference, true); message != "" {
			http.Error(w, message, 400)
			return
		}
		trace = reference.Trace
		for k, v := range map[string]string{"from": reference.From, "to": reference.To, "version": reference.Version} {
			if !q.Has(k) && v != "" {
				q.Set(k, v)
			}
		}
	}
	from, err := parseTime(q.Get("from"), s.Store.Now().Add(-30*24*time.Hour).UnixNano())
	if err != nil {
		http.Error(w, "Use an ISO trace start", 400)
		return
	}
	to, err := parseTime(q.Get("to"), s.Store.Now().Add(5*time.Minute).UnixNano())
	if err != nil || from >= to {
		http.Error(w, "Choose an increasing trace interval", 400)
		return
	}
	es, err := s.Store.EvidenceLimit(10001)
	if err != nil {
		http.Error(w, "Store unavailable", 503)
		return
	}
	var rev string
	if err = s.Store.DB.QueryRow(`SELECT value FROM settings WHERE key='revision'`).Scan(&rev); err != nil {
		http.Error(w, "Store unavailable", 503)
		return
	}
	truncated := len(es) > 10000
	if truncated {
		es = es[len(es)-10000:]
	}
	cut := s.Store.Now().Add(-30 * 24 * time.Hour).UnixNano()
	retained := []Evidence{}
	for _, e := range es {
		if e.Time > cut && e.Time >= from && e.Time < to && matchesVersion(e.Version, q.Get("version")) {
			retained = append(retained, e)
		}
	}
	if trace != "" {
		d := traceDetail(retained, trace, rev, truncated)
		defs, err := s.Store.Definitions()
		if err != nil {
			http.Error(w, "Store unavailable", 503)
			return
		}
		var preview int64
		if q.Has("previewDefinition") {
			preview, err = strconv.ParseInt(q.Get("previewDefinition"), 10, 64)
			found := false
			for _, def := range defs {
				if def.Version == preview {
					found = true
				}
			}
			if err != nil || !found {
				http.Error(w, "Preview definition unavailable", 400)
				return
			}
			d.Interpretation = "explicit non-persisting preview"
		} else {
			d.Interpretation = "capture-time definitions"
		}
		if reference != nil {
			preview = reference.Definition
			d.Interpretation = "reference capture-time definitions"
			if preview > 0 {
				found := false
				for _, def := range defs {
					if def.Version == preview {
						found = true
					}
				}
				if !found {
					http.Error(w, "Reference definition unavailable", 400)
					return
				}
				d.Interpretation = "reference pinned definition preview"
			}
		}
		projected := []Evidence{}
		for _, op := range d.Operations {
			projected = append(projected, op.Evidence)
		}
		projected = append(projected, d.Events...)
		projected = append(projected, d.Logs...)
		d.Business = traceBusiness(projected, defs, preview)
		d.Suggestions = []StageSuggestion{}
		active := defs[len(defs)-1]
		for _, e := range projected {
			_, matched, conflict := mappingFor(e, active)
			if !matched && !conflict && e.Attrs["stage.name"] != "" {
				d.Suggestions = append(d.Suggestions, StageSuggestion{EvidenceID: e.ID, Label: e.Attrs["stage.name"]})
			}
		}
		if len(d.Operations)+len(d.Events)+len(d.Logs) == 0 {
			http.Error(w, "Trace unavailable or outside retained scan", 404)
			return
		}
		write(w, d)
		return
	}
	type item struct {
		Trace    string   `json:"trace"`
		Spans    int      `json:"spans"`
		Records  int      `json:"records"`
		Services []string `json:"services"`
	}
	items := []item{}
	index := map[string]int{}
	for i := len(retained) - 1; i >= 0; i-- {
		e := retained[i]
		if !validTraceID(e.Trace) {
			continue
		}
		n, ok := index[e.Trace]
		if !ok {
			if len(items) >= 200 {
				truncated = true
				continue
			}
			n = len(items)
			index[e.Trace] = n
			items = append(items, item{Trace: e.Trace, Services: []string{}})
		}
		x := &items[n]
		x.Records++
		if e.Signal == "traces" {
			x.Spans++
		}
		service := e.Service
		if service == "" {
			service = "unknown"
		}
		found := false
		for _, v := range x.Services {
			if v == service {
				found = true
			}
		}
		if !found {
			x.Services = append(x.Services, service)
		}
	}
	write(w, map[string]any{"mode": "trace", "revision": rev, "traces": items, "unit": "observed traces", "truncated": truncated, "limits": map[string]int{"scan": 10000, "traces": 200, "detail": 2000}, "completeness": "unknown; observed evidence is not a population or journey count"})
}

type BusinessOccurrence struct {
	EvidenceID int64   `json:"evidenceId"`
	Stage      string  `json:"stage"`
	Kind       string  `json:"kind"`
	Carrier    string  `json:"carrier"`
	Occurrence string  `json:"occurrence"`
	Definition int64   `json:"definition"`
	Basis      string  `json:"basis"`
	DurationNs *string `json:"durationNs"`
	Timing     string  `json:"timing"`
	Goal       bool    `json:"goal"`
	Boundary   string  `json:"boundary"`
	BoundaryNs string  `json:"boundaryNs"`
}

func traceBusiness(es []Evidence, defs []Definition, preview int64) []BusinessOccurrence {
	by := map[int64]Definition{}
	for _, d := range defs {
		by[d.Version] = d
	}
	out := []BusinessOccurrence{}
	for _, e := range es {
		version := e.Definition
		if e.CaptureDefinition > 0 {
			version = e.CaptureDefinition
		}
		if preview != 0 {
			version = preview
		}
		d := by[version]
		m, matched, conflict := mappingFor(e, d)
		if !matched || conflict {
			continue
		}
		b := BusinessOccurrence{EvidenceID: e.ID, Stage: m.Step, Kind: m.Kind, Carrier: evidenceCarrier(e), Definition: version, Basis: "unavailable", Timing: "unavailable", Boundary: "event-time", BoundaryNs: strconv.FormatInt(e.Time, 10)}
		if e.Signal == "traces" {
			b.Occurrence = e.Source + "/" + e.Trace + "/" + e.Span
			if m.Timing == "" || m.Timing == "native-span" {
				b.Basis = "native-inclusive"
				b.DurationNs, b.Timing = nativeDuration(e)
			}
			b.Boundary = "span-end"
			b.BoundaryNs = strconv.FormatInt(e.End, 10)
		}
		if m.Timing == "paired" {
			b.Basis = "paired-stage-elapsed"
		}
		if m.Timing == "estimated" {
			b.Basis = "next-event-estimated"
		}
		b.Goal = isGoal(e, m, d) && (e.Signal != "traces" || e.End >= e.Time && e.End > 0 && !e.Fallback && e.Status != 2)
		out = append(out, b)
	}
	return out
}

type StageSuggestion struct {
	EvidenceID int64  `json:"evidenceId"`
	Label      string `json:"label"`
}
