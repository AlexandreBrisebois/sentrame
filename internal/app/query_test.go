package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"io"
	"maps"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func snapshotTest(t *testing.T, c *http.Client, base, q string) map[string]json.RawMessage {
	t.Helper()
	r, e := c.Get(base + "/api/snapshot" + q)
	if e != nil {
		t.Fatal(e)
	}
	defer r.Body.Close()
	if r.StatusCode != 200 {
		t.Fatalf("snapshot %d", r.StatusCode)
	}
	v := map[string]json.RawMessage{}
	json.NewDecoder(r.Body).Decode(&v)
	return v
}
func TestLiveWindowIncludesContinuationWithoutBaselineDoubleCount(t *testing.T) {
	svc, server, c := fixtureService(t)
	loginTest(t, c, server.URL)
	now := time.Now()
	es := []Evidence{{Source: "browser", Signal: "logs", Name: "discovery.entered", Scope: "ongoing", Tab: "tab", Record: "entry", Time: now.Add(-time.Minute).UnixNano(), Attrs: map[string]string{}}, {Source: "browser", Signal: "logs", Name: "completion.entered", Scope: "ongoing", Tab: "tab", Record: "goal", Time: now.Add(-time.Second).UnixNano(), Attrs: map[string]string{}}}
	svc.Store.Insert(es)
	v := snapshotTest(t, c, server.URL, "?mode=live&baseline=1&liveStart="+url.QueryEscape(now.Add(-30*time.Second).Format(time.RFC3339Nano)))
	var live, base Model
	json.Unmarshal(v["model"], &live)
	json.Unmarshal(v["baselineModel"], &base)
	if len(live.Journeys) != 1 || len(live.Journeys[0].Evidence) != 1 || !live.Journeys[0].Partial {
		t.Fatalf("continuation missing or baseline leaked into live: %+v", live)
	}
	if base.Stats.Completed != 0 || base.Stats.Records != 1 {
		t.Fatalf("baseline future goal leaked: %+v", base.Stats)
	}
}
func TestSummaryScopeCannotSilentlyWidenAndFiltersApply(t *testing.T) {
	svc, server, c := fixtureService(t)
	loginTest(t, c, server.URL)
	day := time.Now().Add(-40 * 24 * time.Hour).UTC().Format("2006-01-02")
	d := DaySummary{Day: day, Groups: map[string]*SummaryGroup{"definition:1/version:a": {Journeys: 3}, "definition:2/version:b": {Journeys: 7}}}
	b, _ := json.Marshal(d)
	svc.Store.DB.Exec(`INSERT INTO summaries VALUES(?,?)`, day, string(b))
	from := time.Now().Add(-60 * 24 * time.Hour).Format(time.RFC3339)
	v := snapshotTest(t, c, server.URL, "?from="+url.QueryEscape(from)+"&scope=one-visit")
	var got []DaySummary
	json.Unmarshal(v["summary"], &got)
	if len(got) != 0 {
		t.Fatal("isolated visit exposed app-wide summary")
	}
	v = snapshotTest(t, c, server.URL, "?from="+url.QueryEscape(from)+"&version=a&definition=1")
	json.Unmarshal(v["summary"], &got)
	if len(got) != 1 || len(got[0].Groups) != 1 || got[0].Groups["definition:1/version:a"] == nil {
		t.Fatalf("summary filters ignored %+v", got)
	}
}

func TestMappingActivationAndSavedViewUseRealAPIAndSurviveRestart(t *testing.T) {
	svc, server, c := fixtureService(t)
	csrf := loginTest(t, c, server.URL)
	base := time.Now().Add(-time.Minute)
	send := func(name, record, producer string, offset time.Duration) {
		ls := logFixture(name)
		l := ls.ResourceLogs().At(0).ScopeLogs().At(0).LogRecords().At(0)
		l.SetTimestamp(pcommon.NewTimestampFromTime(base.Add(offset)))
		l.Attributes().PutStr("event.id", record)
		if status := requestLogs(t, c, server.URL, producer, ls); status != 200 {
			t.Fatal(status)
		}
	}
	send("discovery.entered", "entry", "browser-test", 0)
	post := func(path string, v any) []byte {
		b, _ := json.Marshal(v)
		r, _ := http.NewRequest("POST", server.URL+"/api/"+path, bytes.NewReader(b))
		r.Header.Set("Origin", server.URL)
		r.Header.Set("X-CSRF-Token", csrf)
		resp, e := c.Do(r)
		if e != nil {
			t.Fatal(e)
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode != 200 {
			t.Fatalf("%s %d %s", path, resp.StatusCode, body)
		}
		return body
	}
	d := InitialDefinition()
	d.Goal = "vote.committed"
	preview := post("preview", d)
	if !bytes.Contains(preview, []byte(`"discovery.entered":1`)) {
		t.Fatal("received evidence not previewed")
	}
	var activated Definition
	json.Unmarshal(post("definitions", d), &activated)
	if activated.Version != 2 {
		t.Fatal("definition was not versioned")
	}
	send("vote.committed", "vote1", "server-test", time.Second)
	send("completion.entered", "completion", "browser-test", 2*time.Second)
	send("vote.committed", "vote2", "server-test", 3*time.Second)
	v := snapshotTest(t, c, server.URL, "")
	var m Model
	json.Unmarshal(v["model"], &m)
	if len(m.Journeys) != 2 || m.Journeys[0].Definition != 1 || m.Journeys[1].Definition != 2 || m.Stats.Completed != 2 {
		t.Fatalf("definition change resegmented captured journey %+v", m.Journeys)
	}
	var saved SavedView
	json.Unmarshal(post("views", SavedView{Name: "Captured completion", Mode: "goal-backward", End: "completion.entered", Definition: 1}), &saved)
	if saved.ID == 0 {
		t.Fatal("view not saved")
	}
	path := svc.Store.Path
	svc.Store.DB.Close()
	reopened, e := Open(path)
	if e != nil {
		t.Fatal(e)
	}
	defer reopened.DB.Close()
	defs, e := reopened.Definitions()
	if e != nil || len(defs) != 2 {
		t.Fatal("definitions lost on restart")
	}
	var count int
	reopened.DB.QueryRow(`SELECT count(*) FROM views`).Scan(&count)
	if count != 1 {
		t.Fatal("saved view lost on restart")
	}
}

func TestOverviewSavedWorkspacePreservesScopeFiltersAndPresentation(t *testing.T) {
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
	var session map[string]string
	json.Unmarshal(logged.Body.Bytes(), &session)
	now := time.Now().UTC()
	for i, scope := range []string{"selected-visit", "other-visit"} {
		err := svc.Store.Insert([]Evidence{{Source: "browser", Tab: "tab", Signal: "logs", Name: "discovery.entered", Scope: scope, Record: fmt.Sprint(i), Version: "v1", Time: now.Add(-time.Minute).UnixNano(), Attrs: map[string]string{}}, {Source: "browser", Tab: "tab", Signal: "logs", Name: "completion.entered", Scope: scope, Record: fmt.Sprint(i) + "goal", Version: "v1", Time: now.Add(-time.Second).UnixNano(), Attrs: map[string]string{}}})
		if err != nil {
			t.Fatal(err)
		}
	}
	payload := map[string]any{"name": "My investigation", "mode": "overview", "scope": "selected-visit", "version": "v1", "definition": 1, "from": now.Add(-time.Hour).Format(time.RFC3339Nano), "to": now.Format(time.RFC3339Nano), "layout": "sankey", "lens": "recovery"}
	post := func(v any) (int, []byte) {
		b, _ := json.Marshal(v)
		r := httptest.NewRequest("POST", "http://localhost/api/views", bytes.NewReader(b))
		r.Header.Set("Origin", svc.Origin)
		r.Header.Set("X-CSRF-Token", session["csrf"])
		r.AddCookie(cookie)
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, r)
		return response.Code, response.Body.Bytes()
	}
	code, body := post(payload)
	if code != 200 {
		t.Fatalf("overview save status %d: %s", code, body)
	}
	var saved map[string]any
	json.Unmarshal(body, &saved)
	if saved["scope"] != "selected-visit" || saved["layout"] != "sankey" || saved["lens"] != "recovery" {
		t.Fatalf("saved workspace metadata lost: %s", body)
	}
	request := httptest.NewRequest("GET", "http://localhost/api/snapshot?reference="+fmt.Sprint(saved["id"]), nil)
	request.AddCookie(cookie)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != 200 {
		t.Fatalf("snapshot %d: %s", response.Code, response.Body)
	}
	var snapshot map[string]json.RawMessage
	json.Unmarshal(response.Body.Bytes(), &snapshot)
	var m Model
	json.Unmarshal(snapshot["model"], &m)
	if len(m.Journeys) != 1 || m.Journeys[0].Scope != "selected-visit" || m.Journeys[0].Status != "completed" || len(m.Journeys[0].Evidence) != 2 || len(m.Journeys[0].Path) != 2 {
		t.Fatalf("saved overview lost filters/boundary: %+v", m)
	}
	for _, field := range []string{"layout", "lens"} {
		invalid := maps.Clone(payload)
		invalid[field] = "unsupported"
		if code, _ := post(invalid); code != 400 {
			t.Fatalf("invalid %s accepted: %d", field, code)
		}
	}
}
