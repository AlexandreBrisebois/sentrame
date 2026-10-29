package app

// Evidence is the only retained projection. Attribute values outside this contract are never stored.
type Evidence struct {
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
	Event  string `json:"event"`
	Step   string `json:"step"`
	Source string `json:"source"`
	Kind   string `json:"kind"`
	Key    string `json:"key"`
	Value  string `json:"value"`
}
type Definition struct {
	Version        int64     `json:"version"`
	Goal           string    `json:"goal"`
	Mappings       []Mapping `json:"mappings"`
	FallbackTiming bool      `json:"fallbackTiming"`
}
type SavedView struct {
	Scope           string `json:"scope,omitempty"`
	Layout          string `json:"layout,omitempty"`
	Lens            string `json:"lens,omitempty"`
	StartOccurrence string `json:"startOccurrence"`
	EndOccurrence   string `json:"endOccurrence"`
	ID              int64  `json:"id"`
	Name            string `json:"name"`
	Mode            string `json:"mode"`
	Start           string `json:"start"`
	End             string `json:"end"`
	Version         string `json:"version"`
	Definition      int64  `json:"definition"`
	From            string `json:"from"`
	To              string `json:"to"`
}
