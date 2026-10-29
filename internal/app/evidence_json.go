package app

import (
	"encoding/json"
	"strconv"
)

// Evidence timestamps use decimal strings on JSON boundaries, preserving exact OTLP nanoseconds.
func (e Evidence) MarshalJSON() ([]byte, error) {
	type alias Evidence
	return json.Marshal(struct {
		alias
		Time     string `json:"time"`
		Receipt  string `json:"receipt"`
		Observed string `json:"observed"`
	}{alias: alias(e), Time: strconv.FormatInt(e.Time, 10), Receipt: strconv.FormatInt(e.Receipt, 10), Observed: strconv.FormatInt(e.Observed, 10)})
}

// Accept both the initial local numeric storage format and the exact decimal-string format.
func (e *Evidence) UnmarshalJSON(b []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	for _, key := range []string{"time", "receipt", "observed"} {
		v := raw[key]
		if len(v) > 0 && v[0] == '"' {
			var text string
			if err := json.Unmarshal(v, &text); err != nil {
				return err
			}
			if _, err := strconv.ParseInt(text, 10, 64); err != nil {
				return err
			}
			raw[key] = json.RawMessage(text)
		}
	}
	normalized, err := json.Marshal(raw)
	if err != nil {
		return err
	}
	type alias Evidence
	var out alias
	if err = json.Unmarshal(normalized, &out); err != nil {
		return err
	}
	*e = Evidence(out)
	return nil
}
