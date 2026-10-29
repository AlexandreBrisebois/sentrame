package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestUnsavedReferencePreviewUsesCapturedReadPathWithoutSaving(t *testing.T) {
	store, err := Open(filepath.Join(t.TempDir(), "db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.DB.Close()
	svc, err := New(store, "synthetic owner passphrase", "browser-test", "server-test", "http://localhost")
	if err != nil {
		t.Fatal(err)
	}
	handler := svc.Handler("")
	login := httptest.NewRequest("POST", "http://localhost/api/login", strings.NewReader(`{"passphrase":"synthetic owner passphrase"}`))
	login.Header.Set("Origin", svc.Origin)
	logged := httptest.NewRecorder()
	handler.ServeHTTP(logged, login)
	if logged.Code != 200 {
		t.Fatalf("login %d", logged.Code)
	}
	cookie := logged.Result().Cookies()[0]
	now := time.Now().Add(-time.Minute)
	for i, name := range []string{"discovery.entered", "vote.requested", "completion.entered"} {
		if err := store.Insert([]Evidence{{Source: "browser", Signal: "logs", Name: name, Scope: "visit", Tab: "tab", Record: fmt.Sprint(i), Time: now.Add(time.Duration(i) * time.Second).UnixNano(), Attrs: map[string]string{}}}); err != nil {
			t.Fatal(err)
		}
	}
	query := func(q string) (int, Model) {
		r := httptest.NewRequest("GET", "http://localhost/api/snapshot"+q, nil)
		r.AddCookie(cookie)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)
		var snapshot struct {
			Model Model `json:"model"`
		}
		if w.Code == http.StatusOK {
			if err := json.Unmarshal(w.Body.Bytes(), &snapshot); err != nil {
				t.Fatal(err)
			}
		}
		return w.Code, snapshot.Model
	}
	_, full := query("")
	if len(full.Journeys) != 1 {
		t.Fatalf("fixture %+v", full)
	}
	captured := full.Journeys[0]
	for _, tc := range []struct {
		q        string
		first    string
		last     string
		censored bool
	}{
		{"?referenceMode=goal-backward&referenceEnd=vote.requested", "discovery.entered", "vote.requested", true},
		{"?referenceMode=origin-forward&referenceStart=vote.requested", "vote.requested", "completion.entered", false},
		{"?referenceMode=segment&referenceStart=vote.requested&referenceEnd=completion.entered&startOccurrence=last&endOccurrence=first", "vote.requested", "completion.entered", false},
	} {
		code, m := query(tc.q)
		if code != 200 || len(m.Journeys) != 1 {
			t.Fatalf("preview %s status=%d model=%+v", tc.q, code, m)
		}
		j := m.Journeys[0]
		if len(j.Evidence) != 2 || j.Evidence[0].Name != tc.first || j.Evidence[1].Name != tc.last || j.ID != captured.ID || j.Start != captured.Start || j.End != captured.End || j.CapturedStatus != captured.Status || j.RightCensored != tc.censored {
			t.Fatalf("preview changed capture or failed projection: %+v", j)
		}
	}
	for _, q := range []string{
		"?referenceMode=unsupported", "?referenceMode=overview", "?referenceMode=goal-backward", "?referenceMode=origin-forward", "?referenceMode=segment&referenceStart=vote.requested", "?referenceMode=goal-backward&referenceEnd=vote.requested&startOccurrence=middle", "?referenceMode=goal-backward&referenceEnd=vote.requested&endOccurrence=middle", "?referenceMode=goal-backward&referenceEnd=vote.requested&reference=1", "?referenceStart=vote.requested", "?referenceMode=",
	} {
		if code, _ := query(q); code != 400 {
			t.Errorf("invalid preview accepted: %s status=%d", q, code)
		}
	}
	var count int
	if err := store.DB.QueryRow(`SELECT count(*) FROM views`).Scan(&count); err != nil || count != 0 {
		t.Fatalf("preview persisted a saved view: %d %v", count, err)
	}
}
