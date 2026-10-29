package app

import (
	"encoding/json"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLifecycleSummaryAndNoDetailAfterExpiry(t *testing.T) {
	s, e := Open(filepath.Join(t.TempDir(), "lifecycle.db"))
	if e != nil {
		t.Fatal(e)
	}
	defer s.DB.Close()
	now := time.Date(2026, 10, 28, 12, 0, 0, 0, time.UTC)
	s.Now = func() time.Time { return now }
	es := []Evidence{{Source: "browser", Signal: "logs", Name: "discovery.entered", Scope: "synthetic", Tab: "tab", Record: "first", Time: now.UnixNano(), Attrs: map[string]string{}}, {Source: "browser", Signal: "logs", Name: "completion.entered", Scope: "synthetic", Tab: "tab", Record: "last", Time: now.Add(time.Minute).UnixNano(), Attrs: map[string]string{}}}
	if e = s.Insert(es); e != nil {
		t.Fatal(e)
	}
	if e = s.Cleanup(); e != nil {
		t.Fatal(e)
	}
	now = now.Add(31 * 24 * time.Hour)
	if e = s.Cleanup(); e != nil {
		t.Fatal(e)
	}
	got, _ := s.Evidence()
	if len(got) != 0 {
		t.Fatalf("expired detail still persisted %d", len(got))
	}
	var n int
	s.DB.QueryRow(`SELECT count(*) FROM summaries`).Scan(&n)
	if n != 1 {
		t.Fatalf("supported summary missing: %d", n)
	}
}

func TestFinalizedBoundaryLateGoalCannotCreateEligibleCompletionOrSecondSummary(t *testing.T) {
	s, e := Open(filepath.Join(t.TempDir(), "tail.db"))
	if e != nil {
		t.Fatal(e)
	}
	defer s.DB.Close()
	start := time.Date(2026, 8, 1, 12, 0, 0, 0, time.UTC)
	now := start
	s.Now = func() time.Time { return now }
	makeEvent := func(name, id string, offset time.Duration) Evidence {
		return Evidence{Source: "browser", Signal: "logs", Name: name, Scope: "scope", Tab: "tab", Record: id, Time: start.Add(offset).UnixNano(), Attrs: map[string]string{}}
	}
	if e = s.Insert([]Evidence{makeEvent("discovery.entered", "start", 0), makeEvent("completion.entered", "goal", 2*time.Minute)}); e != nil {
		t.Fatal(e)
	}
	now = start.Add(30*24*time.Hour + time.Minute)
	if e = s.Cleanup(); e != nil {
		t.Fatal(e)
	}
	if e = s.Insert([]Evidence{makeEvent("completion.entered", "late-goal", 90*time.Second)}); e != nil {
		t.Fatal(e)
	}
	es, _ := s.Evidence()
	defs, _ := s.Definitions()
	m := Derive(es, defs, now.UnixNano())
	if m.Stats.EligibleCompleted != 0 {
		t.Fatalf("finalized boundary resurrected eligibility: %+v", m.Stats)
	}
	now = now.Add(3 * time.Minute)
	if e = s.Cleanup(); e != nil {
		t.Fatal(e)
	}
	var raw string
	s.DB.QueryRow(`SELECT data FROM summaries`).Scan(&raw)
	var d DaySummary
	json.Unmarshal([]byte(raw), &d)
	completed := 0
	for _, g := range d.Groups {
		completed += g.Completed
	}
	if completed != 1 {
		t.Fatalf("summary multiplied after late goal: %d", completed)
	}
}
func TestInactivityAdvancesCommittedLifecycleRevision(t *testing.T) {
	s, e := Open(filepath.Join(t.TempDir(), "timeout.db"))
	if e != nil {
		t.Fatal(e)
	}
	defer s.DB.Close()
	now := time.Now()
	s.Now = func() time.Time { return now }
	s.Insert([]Evidence{{Source: "server", Name: "discovery.entered", Scope: "scope", Signal: "logs", Time: now.UnixNano(), Attrs: map[string]string{}}})
	s.Cleanup()
	var before, after string
	s.DB.QueryRow(`SELECT value FROM settings WHERE key='revision'`).Scan(&before)
	now = now.Add(31 * time.Minute)
	if e = s.Cleanup(); e != nil {
		t.Fatal(e)
	}
	s.DB.QueryRow(`SELECT value FROM settings WHERE key='revision'`).Scan(&after)
	if before == after {
		t.Fatal("timeout did not advance lifecycle revision")
	}
}

func TestPublicReadsThirtyDayDetailTwelveMonthSummaryAndFreshGeneration(t *testing.T) {
	svc, server, c := fixtureService(t)
	loginTest(t, c, server.URL)
	base := time.Now()
	clock := base
	svc.Store.Now = func() time.Time { return clock }
	ls := logFixture("discovery.entered", "completion.entered")
	logs := ls.ResourceLogs().At(0).ScopeLogs().At(0).LogRecords()
	for i := 0; i < logs.Len(); i++ {
		logs.At(i).SetTimestamp(pcommon.NewTimestampFromTime(base.Add(time.Duration(i) * time.Millisecond)))
	}
	if code := requestLogs(t, c, server.URL, "browser-test", ls); code != 200 {
		t.Fatal(code)
	}
	clock = base.Add(31 * 24 * time.Hour)
	if err := svc.Store.Cleanup(); err != nil {
		t.Fatal(err)
	}
	from := time.Date(base.Year(), base.Month(), base.Day(), 0, 0, 0, 0, time.UTC).Format(time.RFC3339)
	v := snapshotTest(t, c, server.URL, "?from="+url.QueryEscape(from))
	var detailed []Evidence
	json.Unmarshal(v["evidence"], &detailed)
	var days []DaySummary
	json.Unmarshal(v["summary"], &days)
	if len(detailed) != 0 || len(days) != 1 {
		t.Fatalf("public expiry/detail-summary boundary %d/%d", len(detailed), len(days))
	}
	if strings.Contains(string(v["summary"]), "random-scope") || strings.Contains(string(v["summary"]), "random-tab") {
		t.Fatal("summary retained individual membership")
	}
	clock = base.AddDate(0, 13, 0)
	if err := svc.Store.Cleanup(); err != nil {
		t.Fatal(err)
	}
	v = snapshotTest(t, c, server.URL, "")
	json.Unmarshal(v["summary"], &days)
	if len(days) != 0 {
		t.Fatal("summary resurrected beyond12months")
	}
	fresh, err := Open(filepath.Join(t.TempDir(), "fresh-generation.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer fresh.DB.Close()
	var oldGeneration, newGeneration string
	svc.Store.DB.QueryRow(`SELECT value FROM settings WHERE key='generation'`).Scan(&oldGeneration)
	fresh.DB.QueryRow(`SELECT value FROM settings WHERE key='generation'`).Scan(&newGeneration)
	if oldGeneration == newGeneration {
		t.Fatal("reregistered dataset reusedgeneration")
	}
	next, err := New(fresh, "synthetic owner passphrase", "browser-next", "server-next", "")
	if err != nil {
		t.Fatal(err)
	}
	if err = next.StartReceiver(); err != nil {
		t.Fatal(err)
	}
	defer next.Close()
	origin := httptest.NewServer(next.Handler(""))
	defer origin.Close()
	next.Origin = origin.URL
	if code := requestLogs(t, c, origin.URL, "browser-test", ls); code != 401 {
		t.Fatalf("old queued credential entered newgeneration %d", code)
	}
}
