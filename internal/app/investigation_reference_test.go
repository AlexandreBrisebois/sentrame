package app

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"
)

// These cases use the protected handler and disposable databases, with no
// listener, producer credentials, or dependency on the joint receiver fixture.
func TestInvestigationReferencePresentationRoundTripAndReopen(t *testing.T) {
	x := newInvitationHarness(t)
	metrics := []string{"occurrences", "journeys", "visits", "share"}
	sorts := []string{"traffic", "alphabetic", "journeys", "prevalence", "extras"}
	for i, metric := range metrics {
		for j, sort := range sorts {
			p := &PresentationPreferences{Version: 1, Metric: metric, Stages: []string{"Home", "Vote", "Done"}, Sort: sort, GroupRepeats: (i+j)%2 == 0}
			v := SavedView{Name: "Presentation roundtrip", Mode: "overview", Layout: "matrix", Scope: "selected-visit", Version: "v1", Definition: 1, SelectionPolicy: "captured-start", Presentation: p}
			body, _ := json.Marshal(v)
			w := x.call("POST", "/api/views", string(body), x.owner, x.csrf, true)
			if w.Code != 200 {
				t.Fatalf("%s/%s: %d %s", metric, sort, w.Code, w.Body)
			}
			var saved SavedView
			if err := json.Unmarshal(w.Body.Bytes(), &saved); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(saved.Presentation, p) || saved.Scope != v.Scope || saved.Layout != "matrix" || saved.Version != "v1" || saved.Definition != 1 {
				t.Fatalf("contract lost: %+v", saved)
			}
		}
	}
	path := x.s.Store.Path
	if err := x.s.Store.DB.Close(); err != nil {
		t.Fatal(err)
	}
	reopened, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	x.s.Store = reopened
	t.Cleanup(func() { reopened.DB.Close() })
	w := x.call("GET", "/api/views", "", x.owner, "", false)
	if w.Code != 200 {
		t.Fatal(w.Code, w.Body)
	}
	var views []SavedView
	if err := json.Unmarshal(w.Body.Bytes(), &views); err != nil {
		t.Fatal(err)
	}
	if len(views) != 20 {
		t.Fatalf("reopen lost saved queries: %d", len(views))
	}
	for i, v := range views {
		want := &PresentationPreferences{Version: 1, Metric: metrics[i/5], Stages: []string{"Home", "Vote", "Done"}, Sort: sorts[i%5], GroupRepeats: (i/5+i%5)%2 == 0}
		if !reflect.DeepEqual(v.Presentation, want) {
			t.Fatalf("reopened preferences %d: %+v", i, v.Presentation)
		}
	}
}

func TestInvestigationReferenceRejectsInvalidPreferencesWithoutWriting(t *testing.T) {
	x := newInvitationHarness(t)
	stages := make([]string, 32)
	for i := range stages {
		stages[i] = fmt.Sprintf("Stage %02d", i)
	}
	valid := PresentationPreferences{Version: 1, Metric: "share", Sort: "alphabetic", Stages: stages, GroupRepeats: true}
	for _, tc := range []struct {
		name string
		p    PresentationPreferences
	}{
		{"unsupported-version", PresentationPreferences{Version: 2}},
		{"missing-version", PresentationPreferences{}},
		{"unsupported-metric", PresentationPreferences{Version: 1, Metric: "people"}},
		{"unsupported-sort", PresentationPreferences{Version: 1, Sort: "causal"}},
		{"33-stages", PresentationPreferences{Version: 1, Stages: append(append([]string{}, stages...), "Stage 32")}},
		{"duplicate-stage", PresentationPreferences{Version: 1, Stages: []string{"Vote", "Vote"}}},
		{"empty-stage", PresentationPreferences{Version: 1, Stages: []string{""}}},
		{"overlong-stage", PresentationPreferences{Version: 1, Stages: []string{strings.Repeat("a", 161)}}},
		{"unsafe-stage", PresentationPreferences{Version: 1, Stages: []string{"<script>"}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(SavedView{Name: "Invalid", Mode: "overview", Presentation: &tc.p})
			w := x.call("POST", "/api/views", string(body), x.owner, x.csrf, true)
			if w.Code != 400 || strings.TrimSpace(w.Body.String()) == "" {
				t.Fatalf("invalid accepted or unexplained: %d %s", w.Code, w.Body)
			}
		})
	}
	views, err := x.s.getViews()
	if err != nil || len(views) != 0 {
		t.Fatalf("invalid attempts changed stored queries: %+v %v", views, err)
	}
	body, _ := json.Marshal(SavedView{Name: "Maximum valid", Mode: "overview", Presentation: &valid})
	w := x.call("POST", "/api/views", string(body), x.owner, x.csrf, true)
	if w.Code != 200 {
		t.Fatalf("32-stage boundary rejected: %d %s", w.Code, w.Body)
	}
}

func TestInvestigationReferenceOldAndUnsupportedStoredPreferences(t *testing.T) {
	x := newInvitationHarness(t)
	old := `{"name":"Old reference","mode":"overview","layout":"matrix","scope":"old-scope","version":"v1","definition":1}`
	future := `{"name":"Future reference","mode":"overview","layout":"matrix","scope":"safe-scope","version":"v1","definition":1,"presentation":{"version":2,"metric":"people","stages":["Home"]}}`
	for _, body := range []string{old, future} {
		if _, err := x.s.Store.DB.Exec(`INSERT INTO views(data) VALUES(?)`, body); err != nil {
			t.Fatal(err)
		}
	}
	w := x.call("GET", "/api/views", "", x.owner, "", false)
	if w.Code != 200 {
		t.Fatal(w.Code, w.Body)
	}
	var views []SavedView
	if err := json.Unmarshal(w.Body.Bytes(), &views); err != nil {
		t.Fatal(err)
	}
	if len(views) != 2 || views[0].Presentation != nil || views[0].PresentationWarning != "" {
		t.Fatalf("old defaults contract: %+v", views)
	}
	v := views[1]
	if v.Presentation != nil || v.PresentationWarning != "Unsupported presentation preferences; defaults applied." || v.Scope != "safe-scope" || v.Version != "v1" || v.Definition != 1 || v.Layout != "matrix" {
		t.Fatalf("unsafe fallback or context lost: %+v", v)
	}
	for i, original := range []string{old, future} {
		var raw string
		if err := x.s.Store.DB.QueryRow(`SELECT data FROM views WHERE id=?`, i+1).Scan(&raw); err != nil || raw != original {
			t.Fatalf("read rewrote stored reference %d: %v", i, err)
		}
		w := x.call("GET", fmt.Sprintf("/api/snapshot?reference=%d", i+1), "", x.owner, "", false)
		if w.Code != 200 {
			t.Fatalf("fallback query unavailable: %d %s", w.Code, w.Body)
		}
	}
}

func TestInvestigationReferenceHistoricalWindowAndViewerReadOnly(t *testing.T) {
	x := newInvitationHarness(t)
	now := time.Now().UTC()
	from, to := now.Add(-30*time.Second), now.Add(-time.Second)
	if err := x.s.Store.Insert([]Evidence{
		{Source: "browser", Signal: "logs", Name: "discovery.entered", Scope: "continuation", Tab: "tab", Record: "prior", Time: now.Add(-time.Minute).UnixNano(), Attrs: map[string]string{}},
		{Source: "browser", Signal: "logs", Name: "completion.entered", Scope: "continuation", Tab: "tab", Record: "inside", Time: now.Add(-2 * time.Second).UnixNano(), Attrs: map[string]string{}},
		{Source: "browser", Signal: "logs", Name: "discovery.entered", Scope: "outside", Tab: "tab", Record: "later", Time: now.UnixNano(), Attrs: map[string]string{}},
	}); err != nil {
		t.Fatal(err)
	}
	p := &PresentationPreferences{Version: 1, Metric: "journeys", Stages: []string{"Completion"}, Sort: "alphabetic", GroupRepeats: true}
	v := SavedView{Name: "Frozen historical", Mode: "overview", Layout: "matrix", From: from.Format(time.RFC3339Nano), To: to.Format(time.RFC3339Nano), SelectionPolicy: "evidence-window", Presentation: p}
	body, _ := json.Marshal(v)
	w := x.call("POST", "/api/views", string(body), x.owner, x.csrf, true)
	if w.Code != 200 {
		t.Fatal(w.Code, w.Body)
	}
	var saved SavedView
	json.Unmarshal(w.Body.Bytes(), &saved)
	_, token := x.create(1)
	viewer, csrf := x.redeem(token)
	for _, auth := range []struct {
		name   string
		viewer bool
	}{{"owner", false}, {"viewer", true}} {
		cookie := x.owner
		if auth.viewer {
			cookie = viewer
		}
		w := x.call("GET", fmt.Sprintf("/api/snapshot?reference=%d", saved.ID), "", cookie, "", false)
		if w.Code != 200 {
			t.Fatal(auth.name, w.Code, w.Body)
		}
		var snap struct {
			Model         Model       `json:"model"`
			BaselineModel *Model      `json:"baselineModel"`
			Views         []SavedView `json:"views"`
			Coverage      struct {
				Current PopulationCoverage `json:"current"`
			} `json:"coverage"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &snap); err != nil {
			t.Fatal(err)
		}
		if len(snap.Model.Journeys) != 1 || snap.Model.Journeys[0].Scope != "continuation" || !snap.Model.Journeys[0].Partial || len(snap.Model.Journeys[0].Evidence) != 1 || snap.Model.Journeys[0].Evidence[0].Record != "inside" || snap.BaselineModel != nil || snap.Coverage.Current.From != from.UnixNano() || snap.Coverage.Current.To != to.UnixNano() || snap.Coverage.Current.SelectionPolicy != "evidence-window" {
			t.Fatalf("%s did not reopen exact historical window: %+v", auth.name, snap)
		}
		if len(snap.Views) != 1 || !reflect.DeepEqual(snap.Views[0].Presentation, p) {
			t.Fatalf("%s lost preferences", auth.name)
		}
	}
	for _, path := range []string{"/api/views", "/api/views/delete"} {
		w := x.call("POST", path, string(body), viewer, csrf, true)
		if w.Code != 403 {
			t.Fatalf("viewer wrote %s: %d %s", path, w.Code, w.Body)
		}
	}
	for _, proof := range []struct {
		origin bool
		csrf   string
	}{{false, ""}, {true, ""}} {
		w := x.call("POST", "/api/views", string(body), x.owner, proof.csrf, proof.origin)
		if w.Code != 403 {
			t.Fatalf("owner write missing proof accepted: %d", w.Code)
		}
	}
	views, err := x.s.getViews()
	if err != nil || len(views) != 1 {
		t.Fatalf("denied writes mutated references: %d %v", len(views), err)
	}
}

func TestInvestigationReferenceNeverStoresTransientOrRawExampleFields(t *testing.T) {
	x := newInvitationHarness(t)
	body := `{"name":"Safe reference","mode":"overview","layout":"matrix","presentation":{"version":1,"metric":"share","stages":["Home"],"sort":"traffic","groupRepeats":false,"rawExamples":[{"secret":"synthetic-unpersisted"}]},"evidence":[{"raw":"synthetic-unpersisted"}],"pausedSnapshot":{"credential":"synthetic-unpersisted"},"notificationCursor":"synthetic-unpersisted","returnStack":["synthetic-unpersisted"]}`
	w := x.call("POST", "/api/views", body, x.owner, x.csrf, true)
	if w.Code != 200 {
		t.Fatal(w.Code, w.Body)
	}
	var saved SavedView
	if err := json.Unmarshal(w.Body.Bytes(), &saved); err != nil {
		t.Fatal(err)
	}
	var raw string
	if err := x.s.Store.DB.QueryRow(`SELECT data FROM views WHERE id=?`, saved.ID).Scan(&raw); err != nil {
		t.Fatal(err)
	}
	for _, forbidden := range []string{"synthetic-unpersisted", "rawExamples", "evidence", "pausedSnapshot", "notificationCursor", "returnStack"} {
		if strings.Contains(raw, forbidden) || strings.Contains(w.Body.String(), forbidden) {
			t.Fatalf("transient/raw field persisted: %s", forbidden)
		}
	}
	if saved.Presentation == nil || saved.Presentation.Metric != "share" || len(saved.Presentation.Stages) != 1 {
		t.Fatal("safe preferences were discarded")
	}
}
