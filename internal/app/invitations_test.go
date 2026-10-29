package app

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

type invitationHarness struct {
	t     *testing.T
	s     *Service
	h     http.Handler
	owner *http.Cookie
	csrf  string
}

func newInvitationHarness(t *testing.T) *invitationHarness {
	t.Helper()
	store, err := Open(filepath.Join(t.TempDir(), "db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { store.DB.Close() })
	s, err := New(store, "synthetic owner passphrase", "browser-test", "server-test", "http://localhost")
	if err != nil {
		t.Fatal(err)
	}
	x := &invitationHarness{t: t, s: s, h: s.Handler("")}
	w := x.call("POST", "/api/login", `{"passphrase":"synthetic owner passphrase"}`, nil, "", true)
	if w.Code != 200 {
		t.Fatal(w.Code, w.Body)
	}
	x.owner = w.Result().Cookies()[0]
	var data map[string]string
	json.Unmarshal(w.Body.Bytes(), &data)
	x.csrf = data["csrf"]
	return x
}
func (x *invitationHarness) call(method, path, body string, cookie *http.Cookie, csrf string, origin bool) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, "http://localhost"+path, strings.NewReader(body))
	if cookie != nil {
		r.AddCookie(cookie)
	}
	if csrf != "" {
		r.Header.Set("X-CSRF-Token", csrf)
	}
	if origin {
		r.Header.Set("Origin", x.s.Origin)
	}
	w := httptest.NewRecorder()
	x.h.ServeHTTP(w, r)
	return w
}
func (x *invitationHarness) create(hours int) (int64, string) {
	x.t.Helper()
	w := x.call("POST", "/api/invitations", fmt.Sprintf(`{"label":"Synthetic review","hours":%d}`, hours), x.owner, x.csrf, true)
	if w.Code != 200 {
		x.t.Fatal(w.Code, w.Body)
	}
	var v struct {
		Invite struct{ ID int64 }
		URL    string
	}
	if err := json.Unmarshal(w.Body.Bytes(), &v); err != nil {
		x.t.Fatal(err)
	}
	token := strings.TrimPrefix(v.URL, x.s.Origin+"/#invite=")
	if len(token) != 64 {
		x.t.Fatal("bounded fragment token missing")
	}
	return v.Invite.ID, token
}
func (x *invitationHarness) redeem(token string) (*http.Cookie, string) {
	x.t.Helper()
	w := x.call("POST", "/api/invitations/redeem", fmt.Sprintf(`{"token":%q}`, token), nil, "", true)
	if w.Code != 200 {
		x.t.Fatal(w.Code, w.Body)
	}
	var d map[string]string
	json.Unmarshal(w.Body.Bytes(), &d)
	if d["role"] != "viewer" || d["expires"] == "" {
		x.t.Fatal("role/expiry absent")
	}
	return w.Result().Cookies()[0], d["csrf"]
}
func TestInvitationsReadOnlyBoundaryAndRevocation(t *testing.T) {
	x := newInvitationHarness(t)
	id, token := x.create(24)
	viewer, csrf := x.redeem(token)
	if viewer.Name != "sentrame_viewer" || !viewer.HttpOnly || viewer.SameSite != http.SameSiteStrictMode {
		t.Fatal("viewer cookie policy", viewer)
	}
	list := x.call("GET", "/api/invitations", "", x.owner, "", false)
	if list.Code != 200 || strings.Contains(list.Body.String(), token) || strings.Contains(list.Body.String(), digest(token)) {
		t.Fatal("list leaked bearer/hash")
	}
	for _, path := range []string{"/api/session", "/api/snapshot", "/api/timeline?window=1h", "/api/errors", "/api/views"} {
		if w := x.call("GET", path, "", viewer, "", false); w.Code != 200 {
			t.Fatalf("read %s got %d %s", path, w.Code, w.Body)
		}
	}
	for _, path := range []string{"/api/definitions", "/api/preview", "/api/views", "/api/views/delete", "/api/delete", "/api/invitations", "/api/invitations/revoke"} {
		if w := x.call("POST", path, `{}`, viewer, csrf, true); w.Code != 403 {
			t.Fatalf("write %s got %d %s", path, w.Code, w.Body)
		}
	}
	if w := x.call("GET", "/api/invitations", "", viewer, "", false); w.Code != 403 {
		t.Fatal("viewer listed invites", w.Code)
	}
	if w := x.call("POST", "/api/invitations/redeem", fmt.Sprintf(`{"token":%q}`, token), nil, "", true); w.Code != 401 {
		t.Fatal("reused token", w.Code)
	}
	if w := x.call("POST", "/api/invitations/revoke", fmt.Sprintf(`{"id":%d}`, id), x.owner, x.csrf, true); w.Code != 200 {
		t.Fatal(w.Code, w.Body)
	}
	if _, kept := x.s.sessions[digest(viewer.Value)]; kept {
		t.Fatal("revoked viewer still consumes session capacity")
	}
	if w := x.call("GET", "/api/snapshot", "", viewer, "", false); w.Code != 401 {
		t.Fatal("revoked read", w.Code)
	}
	if w := x.call("GET", "/api/snapshot", "", x.owner, "", false); w.Code != 200 {
		t.Fatal("owner affected", w.Code)
	}
}
func TestInvitationsOneUseExpiryAndOrigin(t *testing.T) {
	x := newInvitationHarness(t)
	_, token := x.create(1)
	if w := x.call("POST", "/api/invitations/redeem", fmt.Sprintf(`{"token":%q}`, token), nil, "", false); w.Code != 403 {
		t.Fatal("wrong origin", w.Code)
	}
	var wg sync.WaitGroup
	codes := make(chan int, 2)
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			codes <- x.call("POST", "/api/invitations/redeem", fmt.Sprintf(`{"token":%q}`, token), nil, "", true).Code
		}()
	}
	wg.Wait()
	close(codes)
	counts := map[int]int{}
	for c := range codes {
		counts[c]++
	}
	if counts[200] != 1 || counts[401] != 1 {
		t.Fatal("not atomic", counts)
	}
	id, token := x.create(1)
	viewer, _ := x.redeem(token)
	x.s.Store.Now = func() time.Time { return time.Now().Add(2 * time.Hour) }
	if w := x.call("GET", "/api/errors", "", viewer, "", false); w.Code != 401 {
		t.Fatal("expiry", w.Code)
	}
	x.s.Store.Now = time.Now
	_, expired := x.create(1)
	x.s.Store.Now = func() time.Time { return time.Now().Add(2 * time.Hour) }
	if w := x.call("POST", "/api/invitations/redeem", fmt.Sprintf(`{"token":%q}`, expired), nil, "", true); w.Code != 401 {
		t.Fatal("expired invite", w.Code)
	}
	x.s.Store.Now = time.Now
	id, token = x.create(1)
	if w := x.call("POST", "/api/invitations/revoke", fmt.Sprintf(`{"id":%d}`, id), x.owner, x.csrf, true); w.Code != 200 {
		t.Fatal(w.Code)
	}
	if w := x.call("POST", "/api/invitations/redeem", fmt.Sprintf(`{"token":%q}`, token), nil, "", true); w.Code != 401 {
		t.Fatal("revoked invite", w.Code)
	}
}
func TestInvitationsRestartGenerationAndSSE(t *testing.T) {
	x := newInvitationHarness(t)
	_, token := x.create(1)
	s, err := New(x.s.Store, "synthetic owner passphrase", "browser-test", "server-test", "http://localhost")
	if err != nil {
		t.Fatal(err)
	}
	x.s = s
	x.h = s.Handler("")
	viewer, _ := x.redeem(token)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	r := httptest.NewRequest("GET", "http://localhost/api/events", nil).WithContext(ctx)
	r.AddCookie(viewer)
	recorder := &invitationStreamRecorder{ResponseRecorder: httptest.NewRecorder(), flushed: make(chan struct{})}
	done := make(chan struct{})
	go func() { x.h.ServeHTTP(recorder, r); close(done) }()
	// Synchronize on the actual first frame, not an arbitrary scheduling delay.
	select {
	case <-recorder.flushed:
	case <-time.After(time.Second):
		t.Fatal("initial stream frame not flushed")
	}
	x.s.Store.DB.Exec(`UPDATE invitations SET revoked=1`)
	select {
	case <-done:
	case <-time.After(6 * time.Second):
		t.Fatal("revoked stream remained open")
	}
	if !strings.Contains(recorder.Body.String(), "event: reset") {
		t.Fatal("no initial stream frame")
	}
	if w := x.call("GET", "/api/snapshot", "", viewer, "", false); w.Code != 401 {
		t.Fatal("revoked restart session", w.Code)
	}
	_, token = x.createAfterRestart(t)
	x.s.Store.DB.Exec(`UPDATE settings SET value='another-generation' WHERE key='generation'`)
	if w := x.call("POST", "/api/invitations/redeem", fmt.Sprintf(`{"token":%q}`, token), nil, "", true); w.Code != 401 {
		t.Fatal("old app generation redeemed", w.Code)
	}
}
func (x *invitationHarness) createAfterRestart(t *testing.T) (int64, string) {
	// Restore an owner session after restart without opening listeners.
	w := x.call("POST", "/api/login", `{"passphrase":"synthetic owner passphrase"}`, nil, "", true)
	if w.Code != 200 {
		t.Fatal(w.Code)
	}
	x.owner = w.Result().Cookies()[0]
	var d map[string]string
	json.Unmarshal(w.Body.Bytes(), &d)
	x.csrf = d["csrf"]
	return x.create(1)
}

type invitationStreamRecorder struct {
	*httptest.ResponseRecorder
	flushed chan struct{}
	once    sync.Once
}

func (w *invitationStreamRecorder) Flush() {
	w.ResponseRecorder.Flush()
	w.once.Do(func() { close(w.flushed) })
}

func TestInvitationsIDsNeverReuseAndSessionsDoNotSurviveRestart(t *testing.T) {
	x := newInvitationHarness(t)
	id, token := x.create(1)
	viewer, _ := x.redeem(token)
	x.s.Store.DB.Exec(`DELETE FROM invitations WHERE id=?`, id)
	next, token := x.create(1)
	if next <= id {
		t.Fatal("grant identifier reused", id, next)
	}
	x.redeem(token)
	if w := x.call("GET", "/api/session", "", viewer, "", false); w.Code != 401 {
		t.Fatal("old grant resurrected", w.Code)
	}
	s, err := New(x.s.Store, "synthetic owner passphrase", "browser-test", "server-test", "http://localhost")
	if err != nil {
		t.Fatal(err)
	}
	x.s = s
	x.h = s.Handler("")
	if w := x.call("POST", "/api/invitations/redeem", fmt.Sprintf(`{"token":%q}`, token), nil, "", true); w.Code != 401 {
		t.Fatal("consumed link survived restart as usable", w.Code)
	}
	if w := x.call("GET", "/api/session", "", viewer, "", false); w.Code != 401 {
		t.Fatal("guest session survived restart", w.Code)
	}
}

func TestInvitationsCreationGuardsCapacityAndLogout(t *testing.T) {
	x := newInvitationHarness(t)
	for _, tc := range []struct {
		method, path, body string
		cookie             *http.Cookie
		csrf               string
		origin             bool
		want               int
	}{
		{"POST", "/api/invitations", `{"label":"test"}`, nil, "", true, 401},
		{"POST", "/api/invitations", `{"label":"test"}`, x.owner, "", true, 403},
		{"POST", "/api/invitations", `{"label":"test"}`, x.owner, x.csrf, false, 403},
		{"POST", "/api/invitations", `{"hours":2}`, x.owner, x.csrf, true, 400},
		{"GET", "/api/invitations/redeem", "", nil, "", false, 405},
	} {
		if w := x.call(tc.method, tc.path, tc.body, tc.cookie, tc.csrf, tc.origin); w.Code != tc.want {
			t.Fatalf("%+v got %d %s", tc, w.Code, w.Body)
		}
	}
	w := x.call("POST", "/api/invitations", `{}`, x.owner, x.csrf, true)
	if w.Code != 200 {
		t.Fatal("defaults", w.Code, w.Body)
	}
	var made struct {
		Invite Invitation
		URL    string
	}
	json.Unmarshal(w.Body.Bytes(), &made)
	if made.Invite.Label != "Temporary viewer" || made.Invite.Expires.Sub(made.Invite.Created) != 24*time.Hour {
		t.Fatal("defaults wrong", made.Invite)
	}
	viewer, csrf := x.redeem(strings.TrimPrefix(made.URL, x.s.Origin+"/#invite="))
	logout := x.call("POST", "/api/logout", `{}`, viewer, csrf, true)
	if logout.Code != 200 {
		t.Fatal("logout", logout.Code)
	}
	if got := logout.Result().Cookies(); len(got) != 1 || got[0].Name != "sentrame_viewer" || got[0].MaxAge != -1 {
		t.Fatal("logout cookie", got)
	}
	if w := x.call("GET", "/api/session", "", viewer, "", false); w.Code != 401 {
		t.Fatal("logout session survived", w.Code)
	}
	// Expired owner sessions are removed before a fresh owner login reaches its cap.
	x.s.Mu.Lock()
	for i := 0; i < 128; i++ {
		x.s.sessions[fmt.Sprint(i)] = session{Created: time.Now().Add(-9 * time.Hour), Last: time.Now().Add(-9 * time.Hour)}
	}
	x.s.Mu.Unlock()
	x.createAfterRestart(t)
	if len(x.s.sessions) >= 128 {
		t.Fatal("expired owner sessions not pruned")
	}
}

func TestInvitationsExpiryClosesSSEAndDisabledAppRejectsGrant(t *testing.T) {
	x := newInvitationHarness(t)
	id, token := x.create(1)
	viewer, _ := x.redeem(token)
	r := httptest.NewRequest("GET", "http://localhost/api/events", nil)
	r.AddCookie(viewer)
	recorder := &invitationStreamRecorder{ResponseRecorder: httptest.NewRecorder(), flushed: make(chan struct{})}
	done := make(chan struct{})
	go func() { x.h.ServeHTTP(recorder, r); close(done) }()
	select {
	case <-recorder.flushed:
	case <-time.After(time.Second):
		t.Fatal("no frame")
	}
	x.s.Store.DB.Exec(`UPDATE invitations SET expires=? WHERE id=?`, time.Now().Add(-time.Second).UnixNano(), id)
	select {
	case <-done:
	case <-time.After(6 * time.Second):
		t.Fatal("expired stream remained open")
	}
	_, token = x.create(1)
	x.s.Store.DB.Exec(`INSERT OR REPLACE INTO settings(key,value) VALUES('disabled','true')`)
	if w := x.call("POST", "/api/invitations/redeem", fmt.Sprintf(`{"token":%q}`, token), nil, "", true); w.Code != 401 {
		t.Fatal("disabled app grant redeemed", w.Code)
	}
}

func TestInvitationsCookieRolesAndCapacityAreFailClosed(t *testing.T) {
	x := newInvitationHarness(t)
	id, token := x.create(1)
	firstToken := token
	viewer, csrf := x.redeem(token)
	r := httptest.NewRequest("GET", "http://localhost/api/session", nil)
	r.AddCookie(viewer)
	r.AddCookie(x.owner)
	w := httptest.NewRecorder()
	x.h.ServeHTTP(w, r)
	var sessionData map[string]string
	json.Unmarshal(w.Body.Bytes(), &sessionData)
	if w.Code != 200 || sessionData["role"] != "owner" || sessionData["csrf"] != x.csrf {
		t.Fatal("valid owner not preferred", w.Code)
	}
	if w := x.call("GET", "/api/session", "", &http.Cookie{Name: "sentrame_owner", Value: viewer.Value}, "", false); w.Code != 401 {
		t.Fatal("guest token impersonated owner", w.Code)
	}
	if w := x.call("POST", "/api/logout", `{}`, viewer, csrf, false); w.Code != 403 {
		t.Fatal("guest logout skipped origin", w.Code)
	}
	_, token = x.create(1)
	x.s.Mu.Lock()
	for i := 0; len(x.s.sessions) < 128; i++ {
		x.s.sessions[fmt.Sprint(i)] = session{CSRF: "synthetic", Created: time.Now(), Last: time.Now(), InviteID: id, Expires: time.Now().Add(time.Hour)}
	}
	x.s.Mu.Unlock()
	if w := x.call("POST", "/api/invitations/redeem", fmt.Sprintf(`{"token":%q}`, token), nil, "", true); w.Code != 503 {
		t.Fatal("session capacity not bounded", w.Code)
	}
	x.s.Mu.Lock()
	delete(x.s.sessions, "0")
	x.s.Mu.Unlock()
	x.redeem(token) // Capacity rejection must not consume the one-use link.
	if w := x.call("GET", "/api/session", "", x.owner, "", false); w.Code != 200 {
		t.Fatal("capacity evicted owner", w.Code)
	}
	for i := 2; i < 100; i++ {
		x.create(1)
	}
	if w := x.call("POST", "/api/invitations", `{"hours":1}`, x.owner, x.csrf, true); w.Code != 400 {
		t.Fatal("active invitation cap missing", w.Code)
	}
	var hash string
	if err := x.s.Store.DB.QueryRow(`SELECT token_hash FROM invitations WHERE id=?`, id).Scan(&hash); err != nil || hash != digest(firstToken) || hash == firstToken {
		t.Fatal("stored bearer rather than digest", err)
	}
	if w := x.call("POST", "/api/invitations/revoke", fmt.Sprintf(`{"id":%d}`, id), x.owner, x.csrf, true); w.Code != 200 {
		t.Fatal(w.Code)
	}
	if w := x.call("POST", "/api/invitations/revoke", fmt.Sprintf(`{"id":%d}`, id), x.owner, x.csrf, true); w.Code != 200 {
		t.Fatal("revoke not idempotent", w.Code)
	}
}
