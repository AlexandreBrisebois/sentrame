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

func timelineFixture(t *testing.T) (*Store, http.Handler, *http.Cookie) {
	t.Helper()
	store, err := Open(filepath.Join(t.TempDir(), "db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { store.DB.Close() })
	svc, err := New(store, "synthetic owner passphrase", "browser-test", "server-test", "http://localhost")
	if err != nil {
		t.Fatal(err)
	}
	handler := svc.Handler("")
	request := httptest.NewRequest("POST", "http://localhost/api/login", strings.NewReader(`{"passphrase":"synthetic owner passphrase"}`))
	request.Header.Set("Origin", svc.Origin)
	result := httptest.NewRecorder()
	handler.ServeHTTP(result, request)
	if result.Code != 200 {
		t.Fatal(result.Code)
	}
	return store, handler, result.Result().Cookies()[0]
}

func timelineRequest(handler http.Handler, cookie *http.Cookie, q string) *httptest.ResponseRecorder {
	request := httptest.NewRequest("GET", "http://localhost/api/timeline"+q, nil)
	if cookie != nil {
		request.AddCookie(cookie)
	}
	result := httptest.NewRecorder()
	handler.ServeHTTP(result, request)
	return result
}

func TestTimelineProtectedOccurrenceHistogramAndFiltersIgnoreSelectedBounds(t *testing.T) {
	store, handler, cookie := timelineFixture(t)
	now := time.Date(2026, 10, 28, 15, 0, 0, 0, time.UTC)
	store.Now = func() time.Time { return now }
	records := []Evidence{
		{Name: "first", Time: now.Add(-time.Hour).UnixNano(), Scope: "visit", Version: "v1", Definition: 1},
		{Name: "second-bin", Time: now.Add(-59 * time.Minute).UnixNano(), Scope: "visit", Version: "v1", Definition: 1},
		{Name: "last-bin", Time: now.Add(-time.Nanosecond).UnixNano(), Scope: "visit", Version: "v1", Definition: 1},
		{Name: "at-now-excluded", Time: now.UnixNano(), Scope: "visit", Version: "v1", Definition: 1},
		{Name: "before-context", Time: now.Add(-time.Hour - time.Nanosecond).UnixNano(), Scope: "visit", Version: "v1", Definition: 1},
		{Name: "other-visit", Time: now.Add(-time.Minute).UnixNano(), Scope: "other", Version: "v1", Definition: 1},
		{Name: "other-version", Time: now.Add(-time.Minute).UnixNano(), Scope: "visit", Version: "v2", Definition: 1},
		{Name: "other-definition", Time: now.Add(-time.Minute).UnixNano(), Scope: "visit", Version: "v1", Definition: 2},
		{Name: "unknown-version", Time: now.Add(-time.Minute).UnixNano(), Scope: "unknown", Definition: 1},
	}
	if err := store.Insert(records); err != nil {
		t.Fatal(err)
	}
	if result := timelineRequest(handler, nil, ""); result.Code != 401 {
		t.Fatalf("timeline unprotected: %d", result.Code)
	}
	result := timelineRequest(handler, cookie, "?window=1h&scope=visit&version=v1&definition=1&from=invalid&to=invalid")
	if result.Code != 200 {
		t.Fatalf("histogram status%d: %s", result.Code, result.Body)
	}
	var got struct {
		From, To, Unit string
		Total          int
		Truncated      bool
		Bins           []struct {
			From, To string
			Count    int
		}
	}
	if err := json.Unmarshal(result.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.From != now.Add(-time.Hour).Format(time.RFC3339Nano) || got.To != now.Format(time.RFC3339Nano) || got.Total != 3 || got.Truncated || got.Unit != "retained records by occurrence time" || len(got.Bins) != 60 {
		t.Fatalf("wrong context/count/unit: %+v", got)
	}
	if got.Bins[0].Count != 1 || got.Bins[1].Count != 1 || got.Bins[59].Count != 1 {
		t.Fatalf("half-open occurrence binning failed: %+v", got.Bins)
	}
	sum := 0
	for i, bin := range got.Bins {
		sum += bin.Count
		if bin.From != now.Add(-time.Hour+time.Duration(i)*time.Minute).Format(time.RFC3339Nano) || bin.To != now.Add(-time.Hour+time.Duration(i+1)*time.Minute).Format(time.RFC3339Nano) {
			t.Fatal("bins are not equal contiguous intervals")
		}
	}
	if sum != got.Total {
		t.Fatal("histogram total differs from bins")
	}
	if strings.Contains(result.Body.String(), "other-visit") || strings.Contains(result.Body.String(), "first") {
		t.Fatal("histogram disclosed event payload")
	}
	unknown := timelineRequest(handler, cookie, "?window=1h&scope=unknown&version=unknown")
	json.Unmarshal(unknown.Body.Bytes(), &got)
	if got.Total != 1 {
		t.Fatal("unknown version semantics differ from snapshot")
	}
	for _, q := range []string{"?window=bad", "?window=", "?definition=nope", "?definition=0", "?definition=-1", "?definition=1.5"} {
		if result := timelineRequest(handler, cookie, q); result.Code != 400 {
			t.Errorf("invalid timeline %s accepted%d", q, result.Code)
		}
	}
}

func TestTimelineExcludesExpiredDetailImmediatelyAndDisclosesBoundedScan(t *testing.T) {
	store, handler, cookie := timelineFixture(t)
	now := time.Date(2026, 10, 28, 15, 0, 0, 0, time.UTC)
	store.Now = func() time.Time { return now }
	if err := store.Insert([]Evidence{{Name: "expired-exact-cutoff", Time: now.Add(-30 * 24 * time.Hour).UnixNano()}, {Name: "retained", Time: now.Add(-30*24*time.Hour + time.Nanosecond).UnixNano()}}); err != nil {
		t.Fatal(err)
	}
	result := timelineRequest(handler, cookie, "?window=30d")
	if result.Code != 200 {
		t.Fatalf("timeline status%d", result.Code)
	}
	var got struct {
		Total     int
		Truncated bool
		Bins      []struct{ Count int }
	}
	json.Unmarshal(result.Body.Bytes(), &got)
	if got.Total != 1 || got.Bins[0].Count != 1 {
		t.Fatal("expired detail remained in histogram")
	}
	records := make([]Evidence, 10001)
	for i := range records {
		records[i] = Evidence{Name: "bounded", Record: fmt.Sprint(i), Time: now.Add(-time.Second).UnixNano()}
	}
	if err := store.Insert(records); err != nil {
		t.Fatal(err)
	}
	result = timelineRequest(handler, cookie, "?window=1h")
	if result.Code != 200 {
		t.Fatal(result.Code)
	}
	json.Unmarshal(result.Body.Bytes(), &got)
	if !got.Truncated || got.Total != 10000 {
		t.Fatalf("bounded scan silently claimed complete counts: %+v", got)
	}
}
