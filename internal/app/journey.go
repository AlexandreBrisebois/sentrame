package app

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
)

type Model struct {
	Computation string       `json:"computation"`
	Journeys    []Journey    `json:"journeys"`
	Nodes       []Node       `json:"nodes"`
	Edges       []Edge       `json:"edges"`
	Stats       ModelStats   `json:"stats"`
	Diagnostics []Diagnostic `json:"diagnostics"`
	Unassigned  []Evidence   `json:"unassigned"`
	Policy      string       `json:"policy"`
	Warnings    []string     `json:"warnings"`
}
type Journey struct {
	NavigationOnly bool             `json:"navigationOnly"`
	WindowFrom     int64            `json:"windowFrom,string"`
	WindowTo       int64            `json:"windowTo,string"`
	CapturedStatus string           `json:"capturedStatus"`
	RightCensored  bool             `json:"rightCensored"`
	ID             string           `json:"id"`
	Scope          string           `json:"scope"`
	Tab            string           `json:"tab"`
	Definition     int64            `json:"definition"`
	Version        string           `json:"version"`
	MixedVersion   bool             `json:"mixedVersion"`
	Start          int64            `json:"start,string"`
	End            int64            `json:"end,string"`
	LastMapped     int64            `json:"lastMapped,string"`
	Status         string           `json:"status"`
	Reason         string           `json:"reason"`
	Path           []string         `json:"path"`
	Steps          []StepOccurrence `json:"steps"`
	Evidence       []Evidence       `json:"evidence"`
	Timings        []Timing         `json:"timings"`
	Warnings       []string         `json:"warnings"`
	Partial        bool             `json:"partial"`
	Recovered      bool             `json:"recovered"`
}
type StepOccurrence struct {
	GapBefore  bool   `json:"gapBefore"`
	Step       string `json:"step"`
	Event      string `json:"event"`
	EvidenceID int64  `json:"evidenceId"`
	Time       int64  `json:"time,string"`
	Source     string `json:"source"`
}
type Node struct {
	ID          string          `json:"id"`
	Label       string          `json:"label"`
	Occurrences int             `json:"occurrences"`
	Journeys    int             `json:"journeys"`
	Visits      int             `json:"visits"`
	Stops       int             `json:"stops"`
	Pending     int             `json:"pending"`
	Timing      []TimingSummary `json:"timing"`
}
type Edge struct {
	From        string `json:"from"`
	To          string `json:"to"`
	Occurrences int    `json:"occurrences"`
	Journeys    int    `json:"journeys"`
	Visits      int    `json:"visits"`
}
type Timing struct {
	Service      string  `json:"service"`
	DurationNs   *string `json:"durationNs"`
	Computation  string  `json:"computation"`
	Definition   int64   `json:"definition"`
	Stage        string  `json:"stage"`
	Occurrence   string  `json:"occurrence"`
	Source       string  `json:"source"`
	Start        int64   `json:"start,string"`
	End          int64   `json:"end,string"`
	Milliseconds float64 `json:"milliseconds"`
	Basis        string  `json:"basis"`
	Quality      string  `json:"quality"`
	Complete     bool    `json:"complete"`
	Provisional  bool    `json:"provisional"`
	EvidenceIDs  []int64 `json:"evidenceIds"`
}
type TimingSummary struct {
	Computation  string   `json:"computation"`
	Definition   int64    `json:"definition"`
	Basis        string   `json:"basis"`
	Samples      int      `json:"samples"`
	Excluded     int      `json:"excluded"`
	MedianMs     *float64 `json:"medianMs"`
	P90Ms        *float64 `json:"p90Ms"`
	Completeness float64  `json:"completeness"`
}
type Diagnostic struct {
	Source     string `json:"source"`
	EvidenceID int64  `json:"evidenceId"`
	JourneyID  string `json:"journeyId"`
	Scope      string `json:"scope"`
	Tab        string `json:"tab"`
	Type       string `json:"type"`
	Severity   int32  `json:"severity"`
	Linkage    string `json:"linkage"`
	Occurrence string `json:"occurrence"`
	Qualified  bool   `json:"qualified"`
	Continued  bool   `json:"continued"`
	Completed  bool   `json:"completed"`
}
type ModelStats struct {
	EligibleJourneys              int  `json:"eligibleJourneys"`
	EligibleCompleted             int  `json:"eligibleCompleted"`
	PartialJourneys               int  `json:"partialJourneys"`
	Records                       int  `json:"records"`
	Assigned                      int  `json:"assigned"`
	Unassigned                    int  `json:"unassigned"`
	Journeys                      int  `json:"journeys"`
	Visits                        int  `json:"visits"`
	Completed                     int  `json:"completed"`
	Inactive                      int  `json:"inactive"`
	Exited                        int  `json:"exited"`
	Active                        int  `json:"active"`
	DiagnosticRecords             int  `json:"diagnosticRecords"`
	QualifiedExceptionOccurrences int  `json:"qualifiedExceptionOccurrences"`
	AffectedVisits                int  `json:"affectedVisits"`
	ExceptionIdentityUnknown      int  `json:"exceptionIdentityUnknown"`
	LowSupport                    bool `json:"lowSupport"`
}

const inactivityNanos = int64(30 * time.Minute)

// Derive reads an already sanitized, app-isolated retained evidence snapshot. It never
// infers a person, cross-tab join, or a causal relation from temporal proximity.
func Derive(evidence []Evidence, definitions []Definition, nowNanos int64) Model {
	m := Model{Computation: "journey-navigation-v4", Journeys: []Journey{}, Nodes: []Node{}, Edges: []Edge{}, Diagnostics: []Diagnostic{}, Unassigned: []Evidence{}, Warnings: []string{}, Policy: "Correlated browser navigation or first mapped occurrence / optional goal / 30-minute inactivity; event-time display is not causal order. Transition occurrences, journeys and anonymous scopes are separate units. Timings never measure attention."}
	defs := map[int64]Definition{}
	for _, d := range definitions {
		defs[d.Version] = d
	}
	groups := map[string][]Evidence{}
	seen := map[string]bool{}
	for _, e := range evidence {
		if e.Record != "" {
			key := e.Source + "\x00" + e.Signal + "\x00" + e.Record
			if seen[key] {
				continue
			}
			seen[key] = true
		}
		m.Stats.Records++
		if e.Signal == "metrics" || e.Scope == "" || (e.Source == "browser" && e.Tab == "") {
			m.Unassigned = append(m.Unassigned, e)
			continue
		}
		groups[e.Scope+"\x00"+e.Tab] = append(groups[e.Scope+"\x00"+e.Tab], e)
	}
	keys := make([]string, 0, len(groups))
	for k := range groups {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		es := groups[k]
		sequenceOrder := len(es) > 0
		source := ""
		sequences := map[int64]bool{}
		for _, e := range es {
			if source == "" {
				source = e.Source + "\x00" + e.Service
			}
			if e.Source+"\x00"+e.Service != source || e.Sequence <= 0 || sequences[e.Sequence] {
				sequenceOrder = false
			}
			sequences[e.Sequence] = true
		}
		sort.SliceStable(es, func(i, j int) bool {
			if sequenceOrder && es[i].Sequence != es[j].Sequence {
				return es[i].Sequence < es[j].Sequence
			}
			if es[i].Time != es[j].Time {
				return es[i].Time < es[j].Time
			}
			if es[i].Receipt != es[j].Receipt {
				return es[i].Receipt < es[j].Receipt
			}
			return es[i].ID < es[j].ID
		})
		var current *Journey
		var currentDef Definition
		var terminalAt int64
		terminalStatus := ""
		closeTerminal := func() {
			if current != nil && terminalAt > 0 {
				current.End = terminalAt
				current.Status = terminalStatus
				terminalAt = 0
				terminalStatus = ""
			}
		}
		scheduleTerminal := func(at int64, state string) {
			if terminalAt == 0 || at < terminalAt {
				terminalAt = at
				terminalStatus = state
			}
		}
		for _, e := range es {
			if terminalAt > 0 && e.Time >= terminalAt {
				closeTerminal()
			}
			candidateVersion := e.Definition
			if e.CaptureDefinition > 0 {
				candidateVersion = e.CaptureDefinition
			}
			candidate := defs[candidateVersion]
			d := candidate
			if current != nil && current.Status == "active" && e.Time-current.LastMapped < inactivityNanos {
				d = currentDef
			}
			mapping, matched, conflict := journeyMappingFor(e, d)
			if conflict {
				m.Unassigned = append(m.Unassigned, e)
				m.Warnings = appendUnique(m.Warnings, "Conflicting mappings left evidence unassigned; activate a non-conflicting definition.")
				continue
			}
			if matched {
				if current != nil && current.Status == "active" && e.Time-current.LastMapped >= inactivityNanos {
					current.Status = "inactive"
					current.End = current.LastMapped + inactivityNanos
				}
				if current == nil || current.Status != "active" {
					// Select the version captured on the new occurrence, rather than changing
					// an existing active journey's definition when configuration is activated.
					d = candidate
					mapping, matched, conflict = journeyMappingFor(e, d)
					if !matched || conflict {
						m.Unassigned = append(m.Unassigned, e)
						continue
					}
					m.Journeys = append(m.Journeys, Journey{ID: fmt.Sprintf("%s/%s/d%d/e%d", e.Scope, e.Tab, d.Version, e.ID), Scope: e.Scope, Tab: e.Tab, Definition: d.Version, Version: e.Version, Start: e.Time, LastMapped: e.Time, Status: "active", Path: []string{}, Steps: []StepOccurrence{}, Evidence: []Evidence{}, Timings: []Timing{}, Warnings: []string{}})
					current = &m.Journeys[len(m.Journeys)-1]
					currentDef = d
					current.NavigationOnly = d.Goal == ""
					if current.NavigationOnly {
						current.Warnings = appendUnique(current.Warnings, "Correlated navigation; no completion goal configured. Excluded from goal conversion.")
					}
				}
				if e.Time > current.LastMapped {
					current.LastMapped = e.Time
				}
				if e.Signal == "traces" && !e.Fallback && e.End >= e.Time && e.End > current.LastMapped {
					current.LastMapped = e.End
				}
				if mapping.Kind != "stage-end" && mapping.Kind != "activity" {
					label := mapping.Step
					if label == "" {
						label = e.Name
					}
					current.Path = append(current.Path, label)
					current.Steps = append(current.Steps, StepOccurrence{Step: label, Event: e.Name, EvidenceID: e.ID, Time: e.Time, Source: e.Source})
				}
				if isGoal(e, mapping, currentDef) {
					if e.Signal == "traces" {
						if _, quality := nativeDuration(e); quality == "complete" && e.Status != 2 {
							scheduleTerminal(e.End, "completed")
						} else {
							current.Warnings = appendUnique(current.Warnings, "Native goal completion unavailable: failed operation or missing/incompatible end.")
						}
					} else {
						current.Status = "completed"
						current.End = e.Time
						terminalAt = 0
					}
				}
				if mapping.Kind == "exit" && e.Signal == "traces" {
					if _, quality := nativeDuration(e); quality == "complete" {
						state := "exited"
						if isGoal(e, mapping, currentDef) && e.Status != 2 {
							state = "conflict"
							current.Warnings = appendUnique(current.Warnings, "Goal and explicit exit match the same native operation.")
						}
						scheduleTerminal(e.End, state)
						current.Reason = e.Attrs["reason"]
					} else {
						current.Warnings = appendUnique(current.Warnings, "Native exit boundary unavailable: missing/incompatible end.")
					}
				}
				if mapping.Kind == "exit" && e.Signal != "traces" {
					if current.Status == "completed" {
						current.Warnings = appendUnique(current.Warnings, "Goal and explicit exit match the same evidence; configured meanings conflict.")
						current.Status = "conflict"
					}
					if current.Status != "conflict" {
						current.Status = "exited"
					}
					current.End = e.Time
					current.Reason = e.Attrs["reason"]
					if current.Reason == "" {
						current.Reason = "unspecified explicit exit"
					}
				}
			}
			if current != nil && e.Time >= current.Start && e.Time-current.LastMapped <= inactivityNanos && (!current.NavigationOnly || len(currentDef.Mappings) > 0 || e.Source == "browser" && e.Signal == "logs") {
				current.Evidence = append(current.Evidence, e)
				if e.Attrs["_partial_start"] == "true" {
					current.Partial = true
					current.Warnings = appendUnique(current.Warnings, "Original captured beginning is unavailable; partial evidence is excluded from started-journey conversion.")
				}
				if current.Version != e.Version {
					current.MixedVersion = true
				}
				if e.Fallback {
					current.Warnings = appendUnique(current.Warnings, "Missing source time uses receipt time; timing/order unavailable.")
				}
				if e.Record == "" {
					current.Warnings = appendUnique(current.Warnings, "Record identity unavailable; retry counts are uncertain.")
				}
				if !sequenceOrder {
					current.Warnings = appendUnique(current.Warnings, "Display uses occurrence time and stable receipt tie-break; causal order is unproven.")
				}
			} else {
				m.Unassigned = append(m.Unassigned, e)
			}
		}
		if terminalAt > 0 {
			closeTerminal()
		}
		if current != nil && current.Status == "active" && nowNanos-current.LastMapped >= inactivityNanos {
			current.Status = "inactive"
			current.End = current.LastMapped + inactivityNanos
		}
	}
	// Pointer to a slice item is not retained across appends; all remaining
	// calculations use the finalized slice and cannot mutate caller evidence.
	for i := range m.Journeys {
		j := &m.Journeys[i]
		j.Timings = deriveTimings(*j, defs[j.Definition], nowNanos)
		for ti := range j.Timings {
			j.Timings[ti].Definition = j.Definition
		}
		if j.Status == "active" {
			j.Warnings = appendUnique(j.Warnings, "Pending follow-up / right-censored; no stopping cause is known.")
		}
		hasFriction := false
		for _, e := range j.Evidence {
			if isDiagnostic(e) {
				hasFriction = true
			}
		}
		j.Recovered = hasFriction && j.Status == "completed"
	}
	sort.SliceStable(m.Journeys, func(i, j int) bool {
		if m.Journeys[i].Start != m.Journeys[j].Start {
			return m.Journeys[i].Start < m.Journeys[j].Start
		}
		return m.Journeys[i].ID < m.Journeys[j].ID
	})
	summarizeModel(&m)
	return m
}
func mappingFor(e Evidence, d Definition) (Mapping, bool, bool) {
	var found Mapping
	count := 0
	for _, r := range d.Mappings {
		if !mappingMatches(e, r) {
			continue
		}
		found = r
		count++
	}
	return found, count == 1, count > 1
}

// Browser logs with explicit visit/tab correlation can describe navigation
// without assigning business meaning to traces or claiming completion.
// An explicit definition remains authoritative, including excluded events.
func journeyMappingFor(e Evidence, d Definition) (Mapping, bool, bool) {
	if len(d.Mappings) > 0 || d.Goal != "" || d.GoalSelector != nil {
		return mappingFor(e, d)
	}
	if e.Source != "browser" || e.Signal != "logs" || e.Scope == "" || e.Tab == "" || isDiagnostic(e) || safeEvidenceName(e.Name) == "" || e.Name == "unnamed" {
		return Mapping{}, false, false
	}
	label := e.Attrs["stage.name"]
	if label == "" {
		label = e.Name
	}
	return Mapping{Event: e.Name, Step: label, Kind: "step", Source: "browser", Signal: "logs", Timing: "none"}, true, false
}

func appendUnique(xs []string, x string) []string {
	for _, y := range xs {
		if x == y {
			return xs
		}
	}
	return append(xs, x)
}
func isDiagnostic(e Evidence) bool {
	return e.Signal != "metrics" && (e.Severity >= 13 || isErrorDiagnostic(e))
}

func deriveTimings(j Journey, d Definition, now int64) []Timing {
	out := []Timing{}
	starts := map[string]Evidence{}
	ambiguous := map[string]bool{}
	keyFor := func(e Evidence, m Mapping) string {
		return e.Source + "\x00" + e.Service + "\x00" + m.Step + "\x00" + e.Occurrence
	}
	for _, e := range j.Evidence {
		r, ok, _ := mappingFor(e, d)
		if !ok {
			continue
		}
		if e.Signal == "traces" {
			if r.Timing == "" || r.Timing == "native-span" {
				ns, quality := nativeDuration(e)
				t := Timing{Stage: r.Step, Source: e.Source, Service: e.Service, Occurrence: e.Source + "/" + e.Trace + "/" + e.Span, Start: e.Time, End: e.End, Basis: "native-inclusive", Computation: nativeTimingPolicy, Quality: quality, EvidenceIDs: []int64{e.ID}}
				if ns != nil {
					t.Complete = true
					t.Milliseconds = float64(e.End-e.Time) / 1e6
					t.Quality = "inclusive native operation; nested/parallel intervals are not unique time"
				}
				out = append(out, t)
				if t.Complete && e.Source == "browser" && (j.End == 0 || e.End <= j.End) {
					end := e
					end.Time = e.End
					if active, has := activeDuration(j.Evidence, e, end); has {
						duration := fmt.Sprint(active)
						out = append(out, Timing{Stage: r.Step, Occurrence: t.Occurrence, Source: e.Source, Service: e.Service, Start: e.Time, End: e.End, Milliseconds: float64(active) / 1e6, DurationNs: &duration, Basis: "observed-active", Quality: "explicit browser activity intervals; unobserved gaps unknown; 60-second heartbeat cap; not attention", Complete: true, EvidenceIDs: []int64{e.ID}})
					}
				}
			}
			continue
		}
		if r.Timing == "none" {
			continue
		}
		if r.Kind != "stage-start" && r.Kind != "stage-end" {
			continue
		}
		key := keyFor(e, r)
		stage := r.Step
		if stage == "" {
			stage = e.Attrs["stage.name"]
		}
		if e.Occurrence == "" {
			out = append(out, Timing{Stage: stage, Source: e.Source, Service: e.Service, Start: e.Time, Basis: "measured-elapsed", Quality: "unavailable: missing stage occurrence identity", EvidenceIDs: []int64{e.ID}})
			continue
		}
		if r.Kind == "stage-start" {
			if _, exists := starts[key]; exists {
				ambiguous[key] = true
			}
			starts[key] = e
			continue
		}
		start, exists := starts[key]
		t := Timing{Stage: stage, Occurrence: e.Occurrence, Source: e.Source, Service: e.Service, End: e.Time, Basis: "measured-elapsed", Quality: "unavailable: missing start", EvidenceIDs: []int64{e.ID}}
		if exists {
			t.Start = start.Time
			t.EvidenceIDs = []int64{start.ID, e.ID}
			if ambiguous[key] {
				t.Quality = "unavailable: ambiguous repeated start"
			} else if start.Fallback || e.Fallback || e.Time < start.Time {
				t.Quality = "unavailable: incompatible or missing clock"
			} else {
				t.Complete = true
				t.Quality = "same-producer/service wall-clock elapsed; clock stability assumed; may include background/idle"
				t.Milliseconds = float64(e.Time-start.Time) / 1e6
			}
			delete(starts, key)
			delete(ambiguous, key)
			if t.Complete {
				if active, has := activeDuration(j.Evidence, start, e); has {
					out = append(out, Timing{Stage: stage, Occurrence: e.Occurrence, Source: e.Source, Service: e.Service, Start: start.Time, End: e.Time, Milliseconds: float64(active) / 1e6, Basis: "observed-active", Quality: "explicit activity intervals; 60-second heartbeat cap; not attention", Complete: true, EvidenceIDs: []int64{start.ID, e.ID}})
				}
			}
		}
		out = append(out, t)
	}
	keys := make([]string, 0, len(starts))
	for k := range starts {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		e := starts[k]
		r, _, _ := mappingFor(e, d)
		quality := "ongoing / right-censored"
		provisional := j.Status == "active"
		end := now
		if j.End > 0 && j.End < end {
			end = j.End
		}
		if end < e.Time {
			end = e.Time
		}
		if e.Fallback || ambiguous[k] {
			quality = "unavailable: missing clock or ambiguous start"
		}
		out = append(out, Timing{Stage: r.Step, Occurrence: e.Occurrence, Source: e.Source, Service: e.Service, Start: e.Time, Milliseconds: float64(end-e.Time) / 1e6, Basis: "measured-elapsed", Quality: quality, Provisional: provisional, EvidenceIDs: []int64{e.ID}})
	}
	if d.FallbackTiming || hasEstimatedMapping(d) {
		for i := 0; i+1 < len(j.Evidence); i++ {
			a, b := j.Evidence[i], j.Evidence[i+1]
			r, ok, _ := mappingFor(a, d)
			if !ok || !d.FallbackTiming && r.Timing != "estimated" || a.Signal == "traces" || r.Timing == "none" || r.Timing == "paired" || r.Kind != "step" || a.Source != b.Source || a.Service != b.Service || a.Fallback || b.Fallback || b.Time < a.Time || b.Time-a.Time > inactivityNanos {
				continue
			}
			out = append(out, Timing{Stage: r.Step, Occurrence: fmt.Sprintf("estimate-%d", a.ID), Source: a.Source, Service: a.Service, Start: a.Time, End: b.Time, Milliseconds: float64(b.Time-a.Time) / 1e6, Basis: "next-event-estimated", Quality: "same-source adjacent retained evidence; not measured stage end", Complete: true, EvidenceIDs: []int64{a.ID, b.ID}})
		}
	}
	groups := map[string][]int{}
	for i := range out {
		if out[i].Computation == "" {
			out[i].Computation = timingComputation(out[i].Basis)
		}
		if out[i].Complete && out[i].DurationNs == nil {
			duration := out[i].End - out[i].Start
			if out[i].Basis == "observed-active" {
				duration = int64(math.Round(out[i].Milliseconds * 1e6))
			}
			ns := fmt.Sprint(duration)
			out[i].DurationNs = &ns
		}
		if out[i].Occurrence != "" && out[i].Basis != "observed-active" {
			key := out[i].Stage + "\x00" + out[i].Source + "\x00" + out[i].Service + "\x00" + out[i].Occurrence
			groups[key] = append(groups[key], i)
		}
	}
	// Competing elapsed bases cannot silently double-pool one semantic occurrence.
	for _, indices := range groups {
		bases := map[string]bool{}
		for _, i := range indices {
			bases[out[i].Basis] = true
		}
		if len(bases) > 1 {
			for _, i := range indices {
				out[i].Complete = false
				out[i].DurationNs = nil
				out[i].Quality = "unavailable: conflicting occurrence timing policies"
			}
		}
	}
	return out
}
func hasEstimatedMapping(d Definition) bool {
	for _, m := range d.Mappings {
		if m.Timing == "estimated" {
			return true
		}
	}
	return false
}
func timingComputation(basis string) string {
	switch basis {
	case "native-inclusive":
		return nativeTimingPolicy
	case "measured-elapsed":
		return "paired-wall-clock-v1"
	case "observed-active":
		return "presence-cap60-browser-v2"
	case "next-event-estimated":
		return "adjacent-window30m-v1"
	}
	return "unavailable"
}
func activeDuration(es []Evidence, start, end Evidence) (int64, bool) {
	if start.Source != "browser" || end.Source != "browser" {
		return 0, false
	}
	var activeAt int64
	var sum int64
	has := false
	addInterval := func(until int64) {
		if activeAt == 0 {
			return
		}
		stop := until
		if stop > activeAt+int64(time.Minute) {
			stop = activeAt + int64(time.Minute)
		}
		if stop > end.Time {
			stop = end.Time
		}
		begin := activeAt
		if begin < start.Time {
			begin = start.Time
		}
		if stop > begin {
			sum += stop - begin
		}
	}
	for _, e := range es {
		if e.Source != start.Source || e.Service != start.Service || e.Time > end.Time || e.Fallback {
			continue
		}
		state := e.Attrs["activity.state"]
		if state != "active" && state != "visible" && state != "idle" && state != "hidden" {
			continue
		}
		if e.Time >= start.Time || (state == "active" || state == "visible") && start.Time-e.Time <= int64(time.Minute) {
			has = true
		}
		addInterval(e.Time)
		activeAt = 0
		if state == "active" || state == "visible" {
			activeAt = e.Time
		}
	}
	addInterval(end.Time)
	return sum, has
}

// Summarize recomputes map/statistic projections for a selected cohort without
// touching captured segmentation or manufacturing adjacency across hidden events.
func Summarize(journeys []Journey, unassigned []Evidence, nowNanos int64) Model {
	m := Model{Computation: "journey-navigation-v4", Journeys: journeys, Nodes: []Node{}, Edges: []Edge{}, Diagnostics: []Diagnostic{}, Unassigned: unassigned, Warnings: []string{}, Policy: "Captured journey boundaries are preserved; cohort selection does not resegment. Counts describe observed eligible journeys; timing kinds and censoring remain separate."}
	if m.Journeys == nil {
		m.Journeys = []Journey{}
	}
	if m.Unassigned == nil {
		m.Unassigned = []Evidence{}
	}
	for _, j := range journeys {
		m.Stats.Records += len(j.Evidence)
	}
	m.Stats.Records += len(unassigned)
	summarizeModel(&m)
	return m
}

func timingSummaries(ts []Timing) []TimingSummary {
	by := map[string][]Timing{}
	for _, t := range ts {
		key := fmt.Sprintf("%d/%s/%s", t.Definition, t.Basis, t.Computation)
		by[key] = append(by[key], t)
	}
	keys := make([]string, 0, len(by))
	for k := range by {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	out := []TimingSummary{}
	for _, key := range keys {
		s := TimingSummary{Computation: by[key][0].Computation, Definition: by[key][0].Definition, Basis: by[key][0].Basis}
		values := []float64{}
		for _, t := range by[key] {
			if t.Complete && !t.Provisional && t.Milliseconds >= 0 {
				values = append(values, t.Milliseconds)
			} else {
				s.Excluded++
			}
		}
		s.Samples = len(values)
		if len(by[key]) > 0 {
			s.Completeness = float64(s.Samples) / float64(len(by[key]))
		}
		sort.Float64s(values)
		if len(values) > 0 {
			median := values[len(values)/2]
			if len(values)%2 == 0 {
				median = (values[len(values)/2-1] + median) / 2
			}
			s.MedianMs = &median
		}
		if len(values) >= 20 {
			p := values[int(math.Ceil(.9*float64(len(values))))-1]
			s.P90Ms = &p
		}
		out = append(out, s)
	}
	return out
}
func summarizeModel(m *Model) {
	nodes := map[string]*Node{}
	edges := map[string]*Edge{}
	nodeJ := map[string]map[string]bool{}
	nodeV := map[string]map[string]bool{}
	edgeJ := map[string]map[string]bool{}
	edgeV := map[string]map[string]bool{}
	timings := map[string][]Timing{}
	visits := map[string]bool{}
	affected := map[string]bool{}
	exceptionOccurrences := map[string]bool{}
	for _, j := range m.Journeys {
		m.Stats.Journeys++
		if j.Partial {
			m.Stats.PartialJourneys++
		} else if !j.NavigationOnly {
			m.Stats.EligibleJourneys++
			if j.Status == "completed" {
				m.Stats.EligibleCompleted++
			}
		}
		visit := j.Scope + "\x00" + j.Tab
		visits[visit] = true
		m.Stats.Assigned += len(j.Evidence)
		switch j.Status {
		case "completed":
			m.Stats.Completed++
		case "inactive":
			m.Stats.Inactive++
		case "exited":
			m.Stats.Exited++
		case "active":
			m.Stats.Active++
		}
		for i, step := range j.Path {
			n := nodes[step]
			if n == nil {
				n = &Node{ID: step, Label: step, Timing: []TimingSummary{}}
				nodes[step] = n
				nodeJ[step] = map[string]bool{}
				nodeV[step] = map[string]bool{}
			}
			n.Occurrences++
			nodeJ[step][j.ID] = true
			nodeV[step][visit] = true
			if i == len(j.Path)-1 {
				if j.Status == "inactive" || j.Status == "exited" {
					n.Stops++
				}
				if j.Status == "active" {
					n.Pending++
				}
			}
			if i > 0 && (i >= len(j.Steps) || !j.Steps[i].GapBefore) {
				from := j.Path[i-1]
				key := from + "\x00" + step
				edge := edges[key]
				if edge == nil {
					edge = &Edge{From: from, To: step}
					edges[key] = edge
					edgeJ[key] = map[string]bool{}
					edgeV[key] = map[string]bool{}
				}
				edge.Occurrences++
				edgeJ[key][j.ID] = true
				edgeV[key][visit] = true
			}
		}
		for _, t := range j.Timings {
			timings[t.Stage] = append(timings[t.Stage], t)
		}
		for _, e := range j.Evidence {
			if !isDiagnostic(e) {
				continue
			}
			diag := diagnostic(e)
			diag.JourneyID = j.ID
			diag.Linkage = "explicit app-scoped logical correlation and tab; source authority retained"
			diag.Completed = j.Status == "completed"
			for _, next := range j.Steps {
				if next.Time > e.Time {
					diag.Continued = true
					break
				}
			}
			m.Diagnostics = append(m.Diagnostics, diag)
			affected[visit] = true
			if diag.Qualified {
				exceptionOccurrences[visit+"\x00"+diag.Source+"\x00"+diag.Occurrence] = true
			}
		}
	}
	for _, e := range m.Unassigned {
		if isDiagnostic(e) {
			m.Diagnostics = append(m.Diagnostics, diagnostic(e))
		}
	}
	for _, d := range m.Diagnostics {
		m.Stats.DiagnosticRecords++
		if !d.Qualified {
			m.Stats.ExceptionIdentityUnknown++
		}
	}
	m.Stats.QualifiedExceptionOccurrences = len(exceptionOccurrences)
	m.Stats.AffectedVisits = len(affected)
	m.Stats.Visits = len(visits)
	m.Stats.Unassigned = len(m.Unassigned)
	m.Stats.LowSupport = m.Stats.EligibleJourneys < 20
	for key, n := range nodes {
		n.Journeys = len(nodeJ[key])
		n.Visits = len(nodeV[key])
		n.Timing = timingSummaries(timings[key])
		m.Nodes = append(m.Nodes, *n)
	}
	for key, e := range edges {
		e.Journeys = len(edgeJ[key])
		e.Visits = len(edgeV[key])
		m.Edges = append(m.Edges, *e)
	}
	sort.Slice(m.Nodes, func(i, j int) bool { return m.Nodes[i].ID < m.Nodes[j].ID })
	sort.Slice(m.Edges, func(i, j int) bool {
		if m.Edges[i].From != m.Edges[j].From {
			return m.Edges[i].From < m.Edges[j].From
		}
		return m.Edges[i].To < m.Edges[j].To
	})
	if m.Stats.LowSupport {
		m.Warnings = appendUnique(m.Warnings, "Low support: fewer than 20 eligible journeys. Rarity or a large timing node is not a defect.")
	}
	m.Warnings = appendUnique(m.Warnings, "Raw detail can be incomplete due to instrumentation, sampling, retention or filters; conversion describes only observed eligible journeys.")
}
func diagnostic(e Evidence) Diagnostic {
	occ := e.Attrs["exception.occurrence.id"]
	return Diagnostic{EvidenceID: e.ID, Source: e.Source, Scope: e.Scope, Tab: e.Tab, Type: e.Attrs["exception.type"], Severity: e.Severity, Linkage: "unlinked; time proximity does not establish a journey cause", Occurrence: occ, Qualified: occ != "" && strings.TrimSpace(e.Scope) != ""}
}

// ClipJourney projects occurrence-time evidence into a half-open display window.
// Captured boundary metadata remains intact; the query cannot create a journey.
func ClipJourney(j Journey, from, to int64) Journey {
	original := j
	j.WindowFrom = from
	j.WindowTo = to
	j.CapturedStatus = original.Status
	j.Evidence = []Evidence{}
	j.Steps = []StepOccurrence{}
	j.Path = []string{}
	j.Timings = []Timing{}
	j.Warnings = append([]string{}, original.Warnings...)
	ids := map[int64]bool{}
	for _, e := range original.Evidence {
		if e.Time >= from && e.Time < to {
			j.Evidence = append(j.Evidence, e)
			ids[e.ID] = true
		}
	}
	prev := -1
	for i, step := range original.Steps {
		if step.Time < from || step.Time >= to {
			continue
		}
		if i > 0 && (prev < 0 || i != prev+1) {
			step.GapBefore = true
		}
		j.Steps = append(j.Steps, step)
		j.Path = append(j.Path, step.Step)
		prev = i
	}
	if original.Start < from {
		j.Partial = true
		j.Warnings = appendUnique(j.Warnings, "Window omits the captured beginning; original boundary remains recorded.")
	}
	for _, step := range j.Steps {
		if step.GapBefore {
			j.Partial = true
			j.Warnings = appendUnique(j.Warnings, "Window omits intermediate steps; no direct transition is invented across the gap.")
		}
	}
	if original.End == 0 || original.End >= to {
		j.Status = "active"
		j.RightCensored = true
		j.Reason = ""
		j.Recovered = false
		j.Warnings = appendUnique(j.Warnings, "Window ends before a captured closure; pending/right-censored, not an observed dead end.")
	}
	for _, t := range original.Timings {
		if t.Start >= to || t.End > 0 && t.End < from {
			continue
		}
		if t.Start >= from && t.End > 0 && t.End < to {
			j.Timings = append(j.Timings, t)
			continue
		}
		t.Complete = false
		t.Provisional = j.RightCensored
		t.Quality = "window clips captured stage evidence; excluded from completed-duration statistics"
		start := t.Start
		if start < from {
			start = from
		}
		end := t.End
		if end == 0 || end > to {
			end = to
		}
		if end < start {
			end = start
		}
		t.Milliseconds = float64(end-start) / 1e6
		originalIDs := t.EvidenceIDs
		t.EvidenceIDs = []int64{}
		for _, id := range originalIDs {
			if ids[id] {
				t.EvidenceIDs = append(t.EvidenceIDs, id)
			}
		}
		j.Timings = append(j.Timings, t)
	}
	return j
}

func evidenceCarrier(e Evidence) string {
	switch e.Signal {
	case "traces":
		return "span"
	case "span-events":
		return "span-event"
	case "logs":
		return "log"
	case "metrics":
		return "metric"
	}
	return "unknown"
}
func mappingMatches(e Evidence, r Mapping) bool {
	return e.Signal != "metrics" && r.Event == e.Name && (r.Source == "" || r.Source == e.Source) && (r.Service == "" || r.Service == e.Service) && (r.Signal == "" || r.Signal == e.Signal) && (r.Carrier == "" || r.Carrier == evidenceCarrier(e)) && (r.Key == "" || e.Attrs[r.Key] == r.Value)
}
func sameSelector(a, b Mapping) bool {
	return a.Event == b.Event && a.Source == b.Source && a.Service == b.Service && a.Signal == b.Signal && a.Carrier == b.Carrier && a.Key == b.Key && a.Value == b.Value
}
func goalMapping(d Definition) (Mapping, bool) {
	var found Mapping
	count := 0
	for _, m := range d.Mappings {
		if m.Event == d.Goal && (d.GoalSelector == nil || sameSelector(m, *d.GoalSelector)) {
			found = m
			count++
		}
	}
	return found, count == 1
}
func isGoal(e Evidence, m Mapping, d Definition) bool {
	goal, ok := goalMapping(d)
	return ok && sameSelector(m, goal) && mappingMatches(e, goal)
}

func mappingCarrier(m Mapping) string {
	if m.Carrier != "" {
		return m.Carrier
	}
	if m.Signal != "" {
		return evidenceCarrier(Evidence{Signal: m.Signal})
	}
	return ""
}
