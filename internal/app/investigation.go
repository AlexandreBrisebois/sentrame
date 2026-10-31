package app

import "sort"

func validatePresentation(p *PresentationPreferences) string {
	if p == nil {
		return ""
	}
	if p.Version != 1 {
		return "Unsupported presentation preference version"
	}
	switch p.Metric {
	case "", "occurrences", "journeys", "visits", "share":
	default:
		return "Choose a supported transition metric"
	}
	switch p.Sort {
	case "", "traffic", "alphabetic", "journeys", "prevalence", "extras":
	default:
		return "Choose a supported presentation sort"
	}
	if len(p.Stages) > 32 {
		return "Choose at most 32 displayed stages"
	}
	seen := map[string]bool{}
	for _, stage := range p.Stages {
		if safe(stage) == "" || seen[stage] {
			return "Choose distinct supported stage names"
		}
		seen[stage] = true
	}
	return ""
}

type PopulationCoverage struct {
	From              int64   `json:"from,string"`
	To                int64   `json:"to,string"`
	Scope             string  `json:"scope"`
	SelectionPolicy   string  `json:"selectionPolicy"`
	Journeys          int     `json:"journeys"`
	SourcePartial     bool    `json:"sourcePartial"`
	PartialJourneys   int     `json:"partialJourneys"`
	CensoredJourneys  int     `json:"censoredJourneys"`
	Definitions       []int64 `json:"definitions"`
	GraphNodesOmitted int     `json:"graphNodesOmitted"`
	GraphEdgesOmitted int     `json:"graphEdgesOmitted"`
}

func populationCoverage(model Model, from, to int64, scope string, clip, sourcePartial bool, cut int64) PopulationCoverage {
	c := PopulationCoverage{From: from, To: to, Scope: scope, SelectionPolicy: "captured-start", Journeys: len(model.Journeys), SourcePartial: sourcePartial || from < cut, Definitions: []int64{}}
	if clip {
		c.SelectionPolicy = "evidence-window"
	}
	defs := map[int64]bool{}
	for _, j := range model.Journeys {
		defs[j.Definition] = true
		if j.Partial {
			c.PartialJourneys++
		}
		if j.RightCensored {
			c.CensoredJourneys++
		}
		for _, step := range j.Steps {
			if step.GapBefore {
				c.SourcePartial = true
			}
		}
	}
	for d := range defs {
		c.Definitions = append(c.Definitions, d)
	}
	sort.Slice(c.Definitions, func(i, j int) bool { return c.Definitions[i] < c.Definitions[j] })
	return c
}
