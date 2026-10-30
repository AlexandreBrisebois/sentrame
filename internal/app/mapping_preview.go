package app

import "sort"

const mappingPreviewExampleLimit = 5
const mappingPreviewConditionLimit = 64

type previewCondition struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Count int    `json:"count"`
}

type mappingPreviewEvidence struct {
	MissingContext      int                `json:"missingContext"`
	Candidates          int                `json:"candidates"`
	Matches             int                `json:"matches"`
	Examples            []Evidence         `json:"examples"`
	MatchingExamples    []Evidence         `json:"matchingExamples"`
	Conditions          []previewCondition `json:"conditions"`
	ConditionsTruncated bool               `json:"conditionsTruncated"`
}

// Candidates share the event's exact context without its optional condition.
// This keeps alternative field values available even when a condition matches
// nothing. Only retained, allowlisted, usable condition values are offered.
func previewMappingEvidence(es []Evidence, mapping Mapping) mappingPreviewEvidence {
	out := mappingPreviewEvidence{Examples: []Evidence{}, MatchingExamples: []Evidence{}, Conditions: []previewCondition{}}
	context := mapping
	context.Key = ""
	context.Value = ""
	options := map[string]*previewCondition{}
	for i := len(es) - 1; i >= 0; i-- {
		e := es[i]
		if !mappingMatches(e, context) {
			continue
		}
		out.Candidates++
		if len(out.Examples) < mappingPreviewExampleLimit {
			out.Examples = append(out.Examples, e)
		}
		if mappingMatches(e, mapping) {
			out.Matches++
			if e.Scope == "" || e.Source == "browser" && e.Tab == "" {
				out.MissingContext++
			}
			if len(out.MatchingExamples) < mappingPreviewExampleLimit {
				out.MatchingExamples = append(out.MatchingExamples, e)
			}
		}
		// Stable order also makes the bounded set of offered fields deterministic.
		keys := []string{}
		for key := range e.Attrs {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			value := e.Attrs[key]
			if !allowed[key] || value == "" || !validConditionValue(key, value) {
				continue
			}
			id := key + "\x00" + value
			if option := options[id]; option != nil {
				option.Count++
				continue
			}
			if len(options) >= mappingPreviewConditionLimit {
				out.ConditionsTruncated = true
				continue
			}
			options[id] = &previewCondition{Key: key, Value: value, Count: 1}
		}
	}
	for _, option := range options {
		out.Conditions = append(out.Conditions, *option)
	}
	sort.Slice(out.Conditions, func(i, j int) bool {
		a, b := out.Conditions[i], out.Conditions[j]
		if a.Key != b.Key {
			return a.Key < b.Key
		}
		return a.Value < b.Value
	})
	return out
}
