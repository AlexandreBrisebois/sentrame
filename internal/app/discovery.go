package app

import "sort"

type DiscoveryContext struct {
	Service           string         `json:"service"`
	Signal            string         `json:"signal"`
	Carrier           string         `json:"carrier"`
	Source            string         `json:"source"`
	Name              string         `json:"name"`
	Projection        string         `json:"projection"`
	Count             int            `json:"count"`
	StageHints        []string       `json:"stageHints"`
	OmittedStageHints int            `json:"omittedStageHints"`
	Fields            map[string]int `json:"fields"`
}

func contextualDiscovery(es []Evidence) []DiscoveryContext {
	groups := map[string]*DiscoveryContext{}
	labels := map[string]map[string]bool{}
	for _, e := range es {
		key := e.Service + "\x00" + e.Signal + "\x00" + e.Source + "\x00" + e.Name + "\x00" + e.Projection
		d := groups[key]
		if d == nil {
			d = &DiscoveryContext{Service: e.Service, Signal: e.Signal, Carrier: evidenceCarrier(e), Source: e.Source, Name: e.Name, Projection: e.Projection, Fields: map[string]int{}}
			groups[key] = d
			labels[key] = map[string]bool{}
		}
		d.Count++
		if label := e.Attrs["stage.name"]; label != "" {
			labels[key][label] = true
		}
		for k, v := range e.Attrs {
			if allowed[k] && v != "" {
				d.Fields[k]++
			}
		}
	}
	keys := []string{}
	for k := range groups {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	out := []DiscoveryContext{}
	for _, k := range keys {
		d := groups[k]
		d.StageHints = []string{}
		all := []string{}
		for label := range labels[k] {
			all = append(all, label)
		}
		sort.Strings(all)
		for i, label := range all {
			if i >= 32 {
				d.OmittedStageHints = len(all) - 32
				break
			}
			d.StageHints = append(d.StageHints, label)
		}
		out = append(out, *d)
	}
	return out
}
