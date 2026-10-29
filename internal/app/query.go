package app

import (
	"encoding/json"
	"fmt"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

func (s *Service) definitions(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		d, e := s.Store.Definitions()
		if e != nil {
			http.Error(w, "Store unavailable", 503)
			return
		}
		write(w, d)
		return
	}
	if r.Method != "POST" {
		http.Error(w, "POST required", 405)
		return
	}
	var d Definition
	if !decode(w, r, &d) {
		return
	}
	if msg := validateDefinition(d); msg != "" {
		http.Error(w, msg, 400)
		return
	}
	s.Store.Mu.Lock()
	locked := true
	defer func() {
		if locked {
			s.Store.Mu.Unlock()
		}
	}()
	tx, e := s.Store.DB.Begin()
	if e != nil {
		http.Error(w, "Writer unavailable", 503)
		return
	}
	defer tx.Rollback()
	var count int
	if e = tx.QueryRow(`SELECT count(*) FROM definitions`).Scan(&count); e != nil || count >= 1000 {
		http.Error(w, "Definition limit or unavailable store", 503)
		return
	}
	if e = tx.QueryRow(`SELECT coalesce(max(version),0)+1 FROM definitions`).Scan(&d.Version); e != nil {
		http.Error(w, "Writer unavailable", 503)
		return
	}
	b, _ := json.Marshal(d)
	if _, e = tx.Exec(`INSERT INTO definitions VALUES(?,?)`, d.Version, string(b)); e == nil {
		_, e = tx.Exec(`UPDATE settings SET value=CAST(value AS INTEGER)+1 WHERE key='revision'`)
	}
	if e != nil || tx.Commit() != nil {
		http.Error(w, "Writer unavailable", 503)
		return
	}
	s.Store.Mu.Unlock()
	locked = false
	write(w, d)
}
func validateDefinition(d Definition) string {
	if d.Goal == "" && len(d.Mappings) == 0 && d.GoalSelector == nil {
		return ""
	}
	if d.Goal == "" || len(d.Mappings) == 0 || len(d.Mappings) > 100 {
		return "Choose a mapped goal and 1–100 mappings"
	}
	seen := []Mapping{}
	goal := false
	for _, m := range d.Mappings {
		if safeEvidenceName(m.Event) == "" || safe(m.Step) == "" || m.Source != "" && m.Source != "browser" && m.Source != "server" {
			return "Invalid event, step or producer"
		}
		switch m.Kind {
		case "step", "exit", "stage-start", "stage-end", "activity":
		default:
			return "Unknown mapping kind"
		}
		if m.Event == d.Goal {
			goal = true
		}
		for _, prev := range seen {
			if prev.Event != m.Event {
				continue
			}
			disjointContext := prev.Service != "" && m.Service != "" && prev.Service != m.Service || prev.Signal != "" && m.Signal != "" && prev.Signal != m.Signal || mappingCarrier(prev) != "" && mappingCarrier(m) != "" && mappingCarrier(prev) != mappingCarrier(m)
			disjointSource := prev.Source != "" && m.Source != "" && prev.Source != m.Source
			disjointCondition := prev.Key != "" && prev.Key == m.Key && prev.Value != m.Value
			if !disjointContext && !disjointSource && !disjointCondition {
				return "Conflicting event mappings require disjoint producer or attribute conditions"
			}
		}
		if m.Service != "" && safe(m.Service) == "" {
			return "Invalid service selector"
		}
		if m.Signal != "" && m.Signal != "traces" && m.Signal != "span-events" && m.Signal != "logs" {
			return "Choose a supported signal; metrics remain receipt-only"
		}
		if m.Carrier != "" && m.Carrier != "span" && m.Carrier != "span-event" && m.Carrier != "log" {
			return "Choose span, span-event or log carrier"
		}
		if m.Signal != "" && m.Carrier != "" && evidenceCarrier(Evidence{Signal: m.Signal}) != m.Carrier {
			return "Signal and carrier must agree"
		}
		switch m.Timing {
		case "", "native-span", "paired", "estimated", "none":
		default:
			return "Choose an explicit timing basis"
		}
		if m.Timing == "native-span" && m.Carrier != "span" {
			return "Native timing requires a span selector"
		}
		if m.Timing == "paired" && (m.Kind != "stage-start" && m.Kind != "stage-end") {
			return "Paired timing requires explicit boundary kinds"
		}
		if (m.Key == "") != (m.Value == "") {
			return "Condition key and value must both be present or empty"
		}
		seen = append(seen, m)
		if m.Key != "" && !allowed[m.Key] {
			return "Condition field is outside the safe allowlist"
		}
		if m.Value != "" && !validConditionValue(m.Key, m.Value) {
			return "Invalid condition value"
		}
	}
	if _, ok := goalMapping(d); goal && !ok {
		return "Choose exactly one goal selector; disjoint mappings with the same name need an explicit goalSelector"
	}
	if !goal {
		return "Goal must be a mapped event"
	}
	return ""
}
func (s *Service) preview(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "POST required", 405)
		return
	}
	var d Definition
	if !decode(w, r, &d) {
		return
	}
	es, e := s.Store.EvidenceLimit(10001)
	if e != nil {
		http.Error(w, "Store unavailable", 503)
		return
	}
	counts := map[string]int{}
	missingFields := map[string]int{}
	missingByEvent := map[string]int{}
	cut := s.Store.Now().Add(-30 * 24 * time.Hour).UnixNano()
	missing := 0
	for _, x := range es {
		if x.Time <= cut {
			continue
		}
		for _, m := range d.Mappings {
			if m.Event == x.Name && m.Key != "" && x.Attrs[m.Key] == "" {
				missingFields[m.Key]++
			}
			if mappingMatches(x, m) {
				counts[m.Event]++
				if x.Scope == "" || x.Source == "browser" && x.Tab == "" {
					missing++
					missingByEvent[m.Event]++
				}
			}
		}
	}
	write(w, map[string]any{"conflict": validateDefinition(d), "matches": counts, "missingContext": missing, "missingContextByEvent": missingByEvent, "unavailableFields": missingFields, "missingContextCapability": "multi-request journey only; native trace exploration remains available", "policy": "Activation is prospective. Existing captured journeys keep their definition."})
}
func (s *Service) getViews() ([]SavedView, error) {
	rows, e := s.Store.DB.Query(`SELECT id,data FROM views ORDER BY id`)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []SavedView{}
	for rows.Next() {
		var v SavedView
		var b string
		rows.Scan(&v.ID, &b)
		id := v.ID
		json.Unmarshal([]byte(b), &v)
		v.ID = id
		out = append(out, v)
	}
	return out, rows.Err()
}
func (s *Service) views(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		v, e := s.getViews()
		if e != nil {
			http.Error(w, "Store unavailable", 503)
			return
		}
		write(w, v)
		return
	}
	if r.Method != "POST" {
		http.Error(w, "POST required", 405)
		return
	}
	var v SavedView
	if !decode(w, r, &v) {
		return
	}
	if safe(v.Name) == "" {
		http.Error(w, "Name required", 400)
		return
	}
	if message := validateReferenceQuery(v, true); message != "" {
		http.Error(w, message, 400)
		return
	}
	if v.Mode == "trace" {
		if v.Layout != "" && v.Layout != "tree" && v.Layout != "timeline" && v.Layout != "table" {
			http.Error(w, "Choose a trace tree, timeline or table", 400)
			return
		}
	} else if v.Layout != "" && v.Layout != "hierarchy" && v.Layout != "force" && v.Layout != "sankey" {
		http.Error(w, "Choose a supported graph layout", 400)
		return
	}
	if v.Lens != "" && v.Lens != "none" && v.Lens != "rare" && v.Lens != "endings" && v.Lens != "exceptions" && v.Lens != "recovery" {
		http.Error(w, "Choose a supported graph highlight", 400)
		return
	}
	s.Store.Mu.Lock()
	locked := true
	defer func() {
		if locked {
			s.Store.Mu.Unlock()
		}
	}()
	if v.Mode == "trace" && v.Definition > 0 {
		defs, err := s.Store.Definitions()
		found := false
		for _, d := range defs {
			if d.Version == v.Definition {
				found = true
			}
		}
		if err != nil {
			http.Error(w, "Store unavailable", 503)
			return
		}
		if !found {
			http.Error(w, "Reference definition unavailable", 400)
			return
		}
	}
	var n int
	s.Store.DB.QueryRow(`SELECT count(*) FROM views`).Scan(&n)
	if n >= 100 {
		http.Error(w, "Saved-view limit 100", 400)
		return
	}
	b, _ := json.Marshal(v)
	res, e := s.Store.DB.Exec(`INSERT INTO views(data) VALUES(?)`, string(b))
	if e != nil {
		http.Error(w, "Writer unavailable", 503)
		return
	}
	v.ID, _ = res.LastInsertId()
	s.Store.Mu.Unlock()
	locked = false
	write(w, v)
}

// Delete only the saved query; evidence and captured boundaries are independent.
func (s *Service) deleteView(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "POST required", 405)
		return
	}
	var input struct {
		ID int64 `json:"id"`
	}
	if !decode(w, r, &input) {
		return
	}
	if input.ID <= 0 {
		http.Error(w, "Reference ID required", 400)
		return
	}
	s.Store.Mu.Lock()
	defer s.Store.Mu.Unlock()
	result, err := s.Store.DB.Exec(`DELETE FROM views WHERE id=?`, input.ID)
	if err != nil {
		http.Error(w, "Writer unavailable", 503)
		return
	}
	count, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Writer unavailable", 503)
		return
	}
	if count == 0 {
		http.Error(w, "Reference not found", 404)
		return
	}
	write(w, map[string]bool{"ok": true})
}

func parseTime(v string, fallback int64) (int64, error) {
	if v == "" {
		return fallback, nil
	}
	t, e := time.Parse(time.RFC3339Nano, v)
	return t.UnixNano(), e
}
func (s *Service) snapshot(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "GET required", 405)
		return
	}
	s.Store.Mu.Lock()
	locked := true
	defer func() {
		if locked {
			s.Store.Mu.Unlock()
		}
	}()
	now := s.Store.Now().UnixNano()
	cut := s.Store.Now().Add(-30 * 24 * time.Hour).UnixNano()
	q := r.URL.Query()
	es, e := s.Store.EvidenceLimit(10001)
	if e != nil {
		http.Error(w, "Store unavailable", 503)
		return
	}
	defs, e := s.Store.Definitions()
	if e != nil {
		http.Error(w, "Store unavailable", 503)
		return
	}
	views, viewErr := s.getViews()
	if viewErr != nil {
		http.Error(w, "Store unavailable", 503)
		return
	}
	var reference *SavedView
	occurrenceCheck := SavedView{Mode: "overview", StartOccurrence: q.Get("startOccurrence"), EndOccurrence: q.Get("endOccurrence")}
	if message := validateReferenceQuery(occurrenceCheck, true); message != "" {
		http.Error(w, message, 400)
		return
	}
	previewRequested := q.Has("referenceMode") || q.Has("referenceStart") || q.Has("referenceEnd")
	if previewRequested {
		if q.Get("reference") != "" {
			http.Error(w, "Choose a saved reference or an unsaved preview, not both", 400)
			return
		}
		preview := SavedView{Mode: q.Get("referenceMode"), Start: q.Get("referenceStart"), End: q.Get("referenceEnd"), StartOccurrence: q.Get("startOccurrence"), EndOccurrence: q.Get("endOccurrence")}
		if message := validateReferenceQuery(preview, false); message != "" {
			http.Error(w, message, 400)
			return
		}
		reference = &preview
	}
	if q.Get("reference") != "" {
		for i := range views {
			if fmt.Sprint(views[i].ID) == q.Get("reference") {
				copy := views[i]
				if v := q.Get("startOccurrence"); v == "first" || v == "last" {
					copy.StartOccurrence = v
				}
				if v := q.Get("endOccurrence"); v == "first" || v == "last" {
					copy.EndOccurrence = v
				}
				reference = &copy
				break
			}
		}
		if reference == nil {
			http.Error(w, "Saved reference unavailable", 400)
			return
		}
	}
	if reference != nil && reference.Mode == "trace" {
		http.Error(w, "Open this reference in trace exploration", 400)
		return
	}
	// Saved filters are defaults, so opening a reference by ID restores its scope.
	// Explicit query filters may refine the saved investigation without mutating it.
	if reference != nil {
		for key, value := range map[string]string{"from": reference.From, "to": reference.To, "scope": reference.Scope, "version": reference.Version} {
			if !q.Has(key) && value != "" {
				q.Set(key, value)
			}
		}
		if !q.Has("definition") && reference.Definition > 0 {
			q.Set("definition", fmt.Sprint(reference.Definition))
		}
	}
	from, e := parseTime(q.Get("from"), cut)
	if e != nil {
		http.Error(w, "Use an ISO start date", 400)
		return
	}
	to, e := parseTime(q.Get("to"), now)
	if e != nil || from >= to || from < s.Store.Now().AddDate(-1, 0, 0).UnixNano() {
		http.Error(w, "Choose an increasing date range within 12 months", 400)
		return
	}
	var rev string
	s.Store.DB.QueryRow(`SELECT value FROM settings WHERE key='revision'`).Scan(&rev)
	retained := []Evidence{}
	for _, x := range es {
		if x.Time > cut {
			retained = append(retained, x)
		}
	}
	truncated := false
	if len(retained) > 10000 {
		retained = retained[len(retained)-10000:]
		seenScope := map[string]bool{}
		for i := range retained {
			k := retained[i].Scope + "/" + retained[i].Tab
			if !seenScope[k] {
				if retained[i].Attrs == nil {
					retained[i].Attrs = map[string]string{}
				}
				retained[i].Attrs["_partial_start"] = "true"
				seenScope[k] = true
			}
		}
		truncated = true
	}

	type Discovery struct {
		Name     string   `json:"name"`
		Count    int      `json:"count"`
		Sources  []string `json:"sources"`
		Assigned int      `json:"assigned"`
	}
	discoveries := map[string]*Discovery{}
	versionSet := map[string]bool{}
	scopeSet := map[string]bool{}
	for _, x := range retained {
		v := x.Version
		if v == "" {
			v = "unknown"
		}
		versionSet[v] = true
		if x.Scope != "" {
			scopeSet[x.Scope] = true
		}
		d := discoveries[x.Name]
		if d == nil {
			d = &Discovery{Name: x.Name, Sources: []string{}}
			discoveries[x.Name] = d
		}
		d.Count++
		present := false
		for _, v := range d.Sources {
			if v == x.Source {
				present = true
			}
		}
		if !present {
			d.Sources = append(d.Sources, x.Source)
		}
		if x.Scope != "" && (x.Source != "browser" || x.Tab != "") {
			d.Assigned++
		}
	}
	discovered := []Discovery{}
	for _, d := range discoveries {
		sort.Strings(d.Sources)
		discovered = append(discovered, *d)
	}
	sort.Slice(discovered, func(i, j int) bool { return discovered[i].Name < discovered[j].Name })
	appVersions := []string{}
	for v := range versionSet {
		appVersions = append(appVersions, v)
	}
	sort.Strings(appVersions)
	appScopes := []string{}
	for v := range scopeSet {
		appScopes = append(appScopes, v)
	}
	sort.Strings(appScopes)
	full := Derive(retained, defs, now)
	selectModel := func(a, b int64, scope string, clip bool) Model {
		js := []Journey{}
		for _, j := range full.Journeys {
			if (!clip && (j.Start < a || j.Start >= b)) || scope != "" && j.Scope != scope || !matchesVersion(j.Version, q.Get("version")) {
				continue
			}
			if q.Get("definition") != "" && fmt.Sprint(j.Definition) != q.Get("definition") {
				continue
			}
			if clip {
				j = ClipJourney(j, a, b)
				if len(j.Evidence) == 0 {
					continue
				}
			}
			if reference != nil {
				var ok bool
				j, ok = referenceJourney(j, *reference)
				if !ok {
					continue
				}
			}
			js = append(js, j)
		}
		us := []Evidence{}
		for _, x := range full.Unassigned {
			if x.Time >= a && x.Time < b && (scope == "" || x.Scope == scope) && matchesVersion(x.Version, q.Get("version")) {
				us = append(us, x)
			}
		}
		out := Summarize(js, us, now)
		if len(out.Nodes) > 100 {
			out.Nodes = out.Nodes[:100]
			truncated = true
		}
		allowedNodes := map[string]bool{}
		for _, node := range out.Nodes {
			allowedNodes[node.ID] = true
		}
		edges := []Edge{}
		for _, edge := range out.Edges {
			if allowedNodes[edge.From] && allowedNodes[edge.To] {
				edges = append(edges, edge)
			}
		}
		if len(edges) > 300 {
			edges = edges[:300]
			truncated = true
		}
		out.Edges = edges
		return out
	}
	liveStart, err := parseTime(q.Get("liveStart"), now)
	if err != nil || liveStart > now {
		http.Error(w, "Invalid live start", 400)
		return
	}
	model := selectModel(from, to, q.Get("scope"), false)
	var base any = nil
	var bounds any = nil
	if q.Get("mode") == "live" {
		model = selectModel(liveStart, now+1, q.Get("scope"), true)
		if q.Get("baseline") == "1" {
			base = selectModel(liveStart-int64(7*24*time.Hour), liveStart, "", true)
			bounds = map[string]int64{"from": liveStart - int64(7*24*time.Hour), "to": liveStart}
		}
	}
	if q.Get("mode") == "live" {
		from = liveStart
		to = now + 1
	}
	selected := []Evidence{}
	for _, x := range retained {
		if x.Time >= from && x.Time < to && (q.Get("scope") == "" || x.Scope == q.Get("scope")) && matchesVersion(x.Version, q.Get("version")) {
			selected = append(selected, x)
		}
	}
	if len(selected) > 2000 {
		selected = selected[len(selected)-2000:]
		truncated = true
	}
	summary := []DaySummary{}
	rows, e := s.Store.DB.Query(`SELECT data FROM summaries ORDER BY day`)
	if e == nil {
		for rows.Next() {
			var d DaySummary
			var b string
			rows.Scan(&b)
			json.Unmarshal([]byte(b), &d)
			dayStart, _ := time.Parse("2006-01-02", d.Day)
			if q.Get("scope") == "" && q.Get("mode") != "live" && !dayStart.Before(time.Unix(0, from)) && !dayStart.Add(24*time.Hour).After(time.Unix(0, to)) {
				for key := range d.Groups {
					parts := strings.SplitN(key, "/version:", 2)
					versionValue := ""
					if len(parts) == 2 {
						versionValue, _ = url.QueryUnescape(strings.TrimSuffix(parts[1], "/unassigned"))
					}
					if !matchesVersion(versionValue, q.Get("version")) || (q.Get("definition") != "" && !strings.HasPrefix(key, "definition:"+q.Get("definition")+"/")) {
						delete(d.Groups, key)
					}
				}
				if len(d.Groups) > 0 {
					summary = append(summary, d)
				}
			}
		}
		rows.Close()
	}
	s.Store.Mu.Unlock()
	locked = false
	write(w, map[string]any{"mode": "journey", "revision": rev, "now": now, "contextualDiscovery": contextualDiscovery(retained), "discovered": discovered, "appVersions": appVersions, "appScopes": appScopes, "evidence": selected, "definitions": defs, "activeDefinition": defs[len(defs)-1], "views": views, "model": model, "baselineModel": base, "baselineBounds": bounds, "summary": summary, "health": map[string]any{"mode": "local synthetic", "collectionCompleteness": "unknown; exporter sampling/loss not observed", "transport": "origin HTTP only; Tunnel and gRPC unqualified", "privacy": diagnosticPolicyDescription, "detailFrom": cut, "summaryOnlyBefore": cut, "freshnessSeconds": 5, "summaryPolicy": "App-wide complete UTC days only, version/definition dimensions. Isolated-visit and partial-day summary queries are unsupported; no individual drilldown. Correction horizon ends captured start+30 days."}, "limits": map[string]any{"truncated": truncated, "modelRecordLimit": 10000, "evidenceLimit": 2000, "policy": "Captured start cohorts; follow-up to snapshot. Partial starts separate; supported summaries only after detail expires."}})
}

// DaySummary contains no individual membership. Quantiles use mergeable histogram counts.
type DaySummary struct {
	Day    string                   `json:"day"`
	Final  bool                     `json:"final"`
	Groups map[string]*SummaryGroup `json:"groups"`
}
type SummaryGroup struct {
	Computation string                 `json:"computation"`
	Records     int                    `json:"records"`
	Journeys    int                    `json:"journeys"`
	Completed   int                    `json:"completed"`
	Inactive    int                    `json:"inactive"`
	Exits       int                    `json:"exits"`
	Histograms  map[string]map[int]int `json:"histograms"`
	Excluded    map[string]int         `json:"excluded"`
	Policy      string                 `json:"policy"`
}

func (s *Service) deletion(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "POST required", 405)
		return
	}
	var v struct {
		Scope  string `json:"scope"`
		App    bool   `json:"app"`
		Bucket string `json:"bucket"`
	}
	if !decode(w, r, &v) {
		return
	}
	if !v.App && safe(v.Scope) == "" {
		http.Error(w, "Choose an explicit app or identifiable visit scope; journey-only deletion unsupported", 400)
		return
	}
	s.Store.Mu.Lock()
	locked := true
	defer func() {
		if locked {
			s.Store.Mu.Unlock()
		}
	}()
	tx, e := s.Store.DB.Begin()
	if e != nil {
		http.Error(w, "Writer unavailable", 503)
		return
	}
	defer tx.Rollback()
	if v.App {
		if v.Bucket != "" {
			if _, e = time.Parse("2006-01-02", v.Bucket); e != nil {
				http.Error(w, "Invalid UTC bucket", 400)
				return
			}
			_, e = tx.Exec(`DELETE FROM summaries WHERE day=?`, v.Bucket)
		} else {
			_, e = tx.Exec(`DELETE FROM evidence; DELETE FROM summaries; INSERT OR REPLACE INTO settings VALUES('disabled','true')`)
		}
	} else {
		_, e = tx.Exec(`INSERT OR REPLACE INTO tombstones VALUES(?,?)`, v.Scope, s.Store.Now().Add(31*24*time.Hour).UnixNano())
		if e == nil {
			_, e = tx.Exec(`DELETE FROM evidence WHERE scope=?`, v.Scope)
		}
	}
	if e == nil {
		_, e = tx.Exec(`UPDATE settings SET value=CAST(value AS INTEGER)+1 WHERE key='revision'`)
	}
	if e != nil || tx.Commit() != nil {
		http.Error(w, "Writer unavailable", 503)
		return
	}
	s.Store.Mu.Unlock()
	locked = false
	write(w, map[string]any{"deleted": true, "limitation": "Individual summary subtraction unavailable after detail expiry; retained client copies/backups are outside server erasure."})
}

func matchesVersion(actual, wanted string) bool {
	return wanted == "" || wanted == actual || (wanted == "unknown" && actual == "")
}

func validConditionValue(key, value string) bool {
	if key == "http.response.status_code" {
		n, err := strconv.ParseInt(value, 10, 64)
		return err == nil && projectedAttribute(key, pcommon.NewValueInt(n)) == value
	}
	if key == "http.route" || key == "http.request.method" {
		return projectedAttribute(key, pcommon.NewValueStr(value)) == value
	}
	return safe(value) != ""
}
