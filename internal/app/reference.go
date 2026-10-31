package app

// referenceJourney projects a query over recorded boundaries. It never merges or derives journeys.
func referenceJourney(j Journey, v SavedView) (Journey, bool) {
	if v.Mode == "overview" {
		return j, true
	}
	if len(j.Evidence) == 0 {
		return j, false
	}
	labels := map[int64]string{}
	if j.NavigationOnly {
		for _, step := range j.Steps {
			labels[step.EvidenceID] = step.Step
		}
	}
	find := func(name, choice string) int {
		found := -1
		for i, e := range j.Evidence {
			if e.Name == name || labels[e.ID] == name {
				found = i
				if choice != "last" {
					break
				}
			}
		}
		return found
	}
	startChoice := v.StartOccurrence
	if startChoice == "" {
		startChoice = "first"
	}
	endChoice := v.EndOccurrence
	if endChoice == "" {
		endChoice = "last"
	}
	a, b := 0, len(j.Evidence)-1
	missing := false
	switch v.Mode {
	case "goal-backward":
		b = find(v.End, endChoice)
		if b < 0 {
			return j, false
		}
	case "origin-forward":
		a = find(v.Start, startChoice)
		if a < 0 {
			return j, false
		}
	case "segment":
		x, y := find(v.Start, startChoice), find(v.End, endChoice)
		if x < 0 && y < 0 {
			return j, false
		}
		if x >= 0 && y >= x {
			a, b = x, y
		} else {
			missing = true
			if x >= 0 {
				a = x
			}
			if y >= 0 && x < 0 {
				b = y
			}
		}
	default:
		return j, false
	}
	original := j
	if j.CapturedStatus == "" {
		j.CapturedStatus = j.Status
	}
	j.Evidence = append([]Evidence{}, original.Evidence[a:b+1]...)
	j.Path = []string{}
	j.Steps = []StepOccurrence{}
	j.Timings = []Timing{}
	j.Warnings = append([]string{}, original.Warnings...)
	selected := map[int64]bool{}
	for _, e := range j.Evidence {
		selected[e.ID] = true
	}
	previous := -1
	for i, step := range original.Steps {
		if !selected[step.EvidenceID] {
			continue
		}
		if previous >= 0 && i != previous+1 {
			step.GapBefore = true
		}
		j.Steps = append(j.Steps, step)
		j.Path = append(j.Path, step.Step)
		previous = i
	}
	if a > 0 || missing {
		j.Partial = true
	}
	if missing {
		j.Warnings = append(j.Warnings, "Requested endpoint is unavailable or crosses a captured boundary; showing only this recorded segment, without joining journeys.")
	}
	if a > 0 {
		j.Warnings = append(j.Warnings, "Reference begins inside a captured journey; earlier evidence remains outside this projection.")
	}
	if b < len(original.Evidence)-1 {
		j.RightCensored = true
		j.Warnings = append(j.Warnings, "Reference ends before later captured evidence; operational goal and captured boundary are unchanged.")
		if j.Status == "completed" || j.Status == "exited" {
			j.Status = "active"
		}
	}
	j.WindowFrom = j.Evidence[0].Time
	j.WindowTo = j.Evidence[len(j.Evidence)-1].Time + 1
	for _, t := range original.Timings {
		overlap := false
		all := len(t.EvidenceIDs) > 0
		for _, id := range t.EvidenceIDs {
			if selected[id] {
				overlap = true
			} else {
				all = false
			}
		}
		if !overlap {
			continue
		}
		if !all {
			t.Complete = false
			t.Provisional = false
			t.Quality = "unavailable: reference excludes stage boundary"
		}
		j.Timings = append(j.Timings, t)
	}
	return j, true
}

// validateReferenceQuery shares endpoint/occurrence rules between saved and
// unsaved projections. Only saved workspace views can use overview mode.
func validateReferenceQuery(v SavedView, allowOverview bool) string {
	if v.SelectionPolicy != "" && v.SelectionPolicy != "captured-start" && v.SelectionPolicy != "evidence-window" {
		return "Unsupported reference selection policy"
	}
	if v.SelectionPolicy == "evidence-window" && (v.Mode == "trace" || v.From == "" || v.To == "") {
		return "Evidence-window references require fixed historical bounds"
	}
	if v.SelectionPolicy == "evidence-window" {
		a, ae := parseTime(v.From, 0)
		b, be := parseTime(v.To, 0)
		if ae != nil || be != nil || a >= b {
			return "Choose increasing ISO historical evidence-window bounds"
		}
	}
	if message := validatePresentation(v.Presentation); message != "" {
		return message
	}
	if v.Mode == "trace" {
		if !allowOverview || !validTraceID(v.Trace) || v.Scope != "" || v.Start != "" || v.End != "" || v.Definition < 0 {
			return "Choose one retained trace and an explicit definition context"
		}
		if v.Version != "" && safe(v.Version) == "" {
			return "Invalid version filter"
		}
		if _, err := parseTime(v.From, 0); err != nil {
			return "Use an ISO trace start"
		}
		if _, err := parseTime(v.To, 0); err != nil {
			return "Use an ISO trace end"
		}
		return ""
	}

	if (v.StartOccurrence != "" && v.StartOccurrence != "first" && v.StartOccurrence != "last") || (v.EndOccurrence != "" && v.EndOccurrence != "first" && v.EndOccurrence != "last") {
		return "Choose first or last endpoint occurrence"
	}
	if v.Mode != "goal-backward" && v.Mode != "origin-forward" && v.Mode != "segment" && !(allowOverview && v.Mode == "overview") {
		return "Choose a supported reference mode"
	}
	if v.Mode == "goal-backward" && safe(v.End) == "" || v.Mode == "origin-forward" && safe(v.Start) == "" || v.Mode == "segment" && (safe(v.Start) == "" || safe(v.End) == "") {
		return "Select reference endpoints"
	}
	return ""
}
