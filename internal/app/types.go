package app

// Evidence is the only retained projection. Attribute values outside this contract are never stored.
const projectionPolicy = "structural-v2"
const schemaVersion = 3

type SpanLink struct {
	Trace string `json:"trace"`
	Span  string `json:"span"`
}

type Evidence struct {
	NameState                 string `json:"nameState"`
	OmittedAttributes         int    `json:"omittedAttributes"`
	ResourceDroppedAttributes uint32 `json:"resourceDroppedAttributes"`
	ScopeDroppedAttributes    uint32 `json:"scopeDroppedAttributes"`

	CaptureDefinition int64      `json:"captureDefinition"`
	Projection        string     `json:"projection"`
	End               int64      `json:"end"`
	Parent            string     `json:"parent"`
	Kind              int32      `json:"kind"`
	Status            int32      `json:"status"`
	Flags             uint32     `json:"flags"`
	EventIndex        int        `json:"eventIndex"`
	EventCount        int        `json:"eventCount"`
	Links             []SpanLink `json:"links"`
	OmittedLinks      int        `json:"omittedLinks"`
	DroppedLinks      uint32     `json:"droppedLinks"`
	DroppedEvents     uint32     `json:"droppedEvents"`
	DroppedAttributes uint32     `json:"droppedAttributes"`

	ID         int64             `json:"id"`
	Source     string            `json:"source"`
	Signal     string            `json:"signal"`
	Name       string            `json:"name"`
	Time       int64             `json:"time"`
	Receipt    int64             `json:"receipt"`
	Observed   int64             `json:"observed"`
	Version    string            `json:"version"`
	Service    string            `json:"service"`
	Scope      string            `json:"scope"`
	Tab        string            `json:"tab"`
	Record     string            `json:"record"`
	Action     string            `json:"action"`
	Occurrence string            `json:"occurrence"`
	Trace      string            `json:"trace"`
	Span       string            `json:"span"`
	Sequence   int64             `json:"sequence"`
	Severity   int32             `json:"severity"`
	Attrs      map[string]string `json:"attrs"`
	Fallback   bool              `json:"fallback"`
	Definition int64             `json:"definition"`
}
type Mapping struct {
	Service string `json:"service,omitempty"`
	Signal  string `json:"signal,omitempty"`
	Carrier string `json:"carrier,omitempty"`
	Timing  string `json:"timing,omitempty"`

	Event  string `json:"event"`
	Step   string `json:"step"`
	Source string `json:"source"`
	Kind   string `json:"kind"`
	Key    string `json:"key"`
	Value  string `json:"value"`
}
type Definition struct {
	GoalSelector   *Mapping  `json:"goalSelector,omitempty"`
	Version        int64     `json:"version"`
	Goal           string    `json:"goal"`
	Mappings       []Mapping `json:"mappings"`
	FallbackTiming bool      `json:"fallbackTiming"`
}
type PresentationPreferences struct {
	Version      int      `json:"version"`
	Metric       string   `json:"metric,omitempty"`
	Stages       []string `json:"stages,omitempty"`
	Sort         string   `json:"sort,omitempty"`
	GroupRepeats bool     `json:"groupRepeats"`
}
type SavedView struct {
	SelectionPolicy     string                   `json:"selectionPolicy,omitempty"`
	Presentation        *PresentationPreferences `json:"presentation,omitempty"`
	PresentationWarning string                   `json:"presentationWarning,omitempty"`
	Trace               string                   `json:"trace,omitempty"`
	Scope               string                   `json:"scope,omitempty"`
	Layout              string                   `json:"layout,omitempty"`
	Lens                string                   `json:"lens,omitempty"`
	StartOccurrence     string                   `json:"startOccurrence"`
	EndOccurrence       string                   `json:"endOccurrence"`
	ID                  int64                    `json:"id"`
	Name                string                   `json:"name"`
	Mode                string                   `json:"mode"`
	Start               string                   `json:"start"`
	End                 string                   `json:"end"`
	Version             string                   `json:"version"`
	Definition          int64                    `json:"definition"`
	From                string                   `json:"from"`
	To                  string                   `json:"to"`
}
