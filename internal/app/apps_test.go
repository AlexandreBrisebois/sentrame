package app

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"go.opentelemetry.io/collector/client"
)

func managedFixture(t *testing.T) (*Registry, *Service, *httptest.Server, *http.Client) {
	t.Helper()
	m, _ := registryFixture(t)
	s, e := NewManaged(m, "synthetic owner passphrase", "")
	if e != nil {
		t.Fatal(e)
	}
	if e = s.StartReceiver(); e != nil {
		t.Fatal(e)
	}
	server := httptest.NewServer(s.Handler(""))
	s.Origin = server.URL
	jar, _ := cookiejar.New(nil)
	c := &http.Client{Jar: jar, Timeout: 10 * time.Second}
	t.Cleanup(func() { server.Close(); s.Close() })
	return m, s, server, c
}
func appRequest(t *testing.T, c *http.Client, base, method, path, csrf string, body any) (int, []byte) {
	t.Helper()
	var b []byte
	if body != nil {
		b, _ = json.Marshal(body)
	}
	r, _ := http.NewRequest(method, base+path, bytes.NewReader(b))
	r.Header.Set("Origin", base)
	r.Header.Set("X-CSRF-Token", csrf)
	r.Header.Set("Content-Type", "application/json")
	resp, e := c.Do(r)
	if e != nil {
		t.Fatal(e)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.Header.Get("Cache-Control") != "no-store" {
		t.Fatal("missing no-store")
	}
	return resp.StatusCode, raw
}
func TestManagedAppsOwnerAPIIsolationAndPrivacy(t *testing.T) {
	m, _, server, c := managedFixture(t)
	if code, _ := appRequest(t, c, server.URL, "GET", "/api/apps", "", nil); code != 401 {
		t.Fatal("unauthenticated list")
	}
	csrf := loginTest(t, c, server.URL)
	if code, _ := appRequest(t, c, server.URL, "POST", "/api/apps", "", map[string]string{"name": "Alpha"}); code != 403 {
		t.Fatal("CSRF omitted")
	}
	code, raw := appRequest(t, c, server.URL, "POST", "/api/apps", csrf, map[string]string{"name": "Alpha"})
	if code != 200 {
		t.Fatal(code, string(raw))
	}
	var a AppEntry
	json.Unmarshal(raw, &a)
	a = m.current.Apps[0]
	b := mustRegister(t, m, "Beta")
	for _, path := range []string{"/api/apps", "/api/session"} {
		code, raw = appRequest(t, c, server.URL, "GET", path, "", nil)
		if code != 200 || bytes.Contains(raw, []byte(a.Server)) || bytes.Contains(raw, []byte(a.Browser)) {
			t.Fatal("ordinary payload disclosed secret")
		}
	}
	if code, _ = appRequest(t, c, server.URL, "GET", "/api/snapshot", "", nil); code != 400 {
		t.Fatal("ambiguous route accepted")
	}
	// Deliberately equal sender-controlled identities and forged app claims.
	ls := logFixture("colliding")
	ls.ResourceLogs().At(0).Resource().Attributes().PutStr("app.id", b.ID)
	for _, a := range []AppEntry{a, b} {
		if code := requestLogs(t, c, server.URL, a.Browser, ls); code != 200 {
			t.Fatal(code)
		}
	}
	for _, a := range []AppEntry{a, b} {
		code, raw = appRequest(t, c, server.URL, "GET", "/api/apps/"+a.ID+"/snapshot", "", nil)
		if code != 200 {
			t.Fatal(code, string(raw))
		}
		var snap struct{ Evidence []Evidence }
		json.Unmarshal(raw, &snap)
		if len(snap.Evidence) != 1 || snap.Evidence[0].Source != "browser" {
			t.Fatal("attribution/isolation")
		}
	}
	base := "/api/apps/" + a.ID
	if code, raw = appRequest(t, c, server.URL, "GET", base+"/connection", "", nil); code != 200 || !bytes.Contains(raw, []byte(a.Server)) {
		t.Fatal("explicit settings unavailable")
	}
	if code, _ = appRequest(t, c, server.URL, "POST", base+"/rename", csrf, map[string]string{"name": "Renamed"}); code != 200 {
		t.Fatal(code)
	}
	if code, _ = appRequest(t, c, server.URL, "POST", base+"/rotate", csrf, map[string]string{"role": "server"}); code != 200 {
		t.Fatal(code)
	}
	if _, _, ok := m.credential("Bearer " + a.Server); ok {
		t.Fatal("rotation kept old authority")
	}
	if _, _, ok := m.credential("Bearer " + a.Browser); !ok {
		t.Fatal("rotation changed other role")
	}
	if code, _ = appRequest(t, c, server.URL, "POST", base+"/delete", csrf, map[string]string{"scope": "random-scope"}); code != 200 {
		t.Fatal(code)
	}
	ea, _ := m.stores[a.ID].Evidence()
	eb, _ := m.stores[b.ID].Evidence()
	if len(ea) != 0 || len(eb) != 1 {
		t.Fatal("visit deletion crossed apps")
	}
	if requestLogs(t, c, server.URL, a.Browser, ls) != 200 {
		t.Fatal("replay rejected")
	}
	ea, _ = m.stores[a.ID].Evidence()
	if len(ea) != 0 {
		t.Fatal("replay resurrected deleted visit")
	}
	if code, _ = appRequest(t, c, server.URL, "POST", base+"/remove", csrf, map[string]string{"confirmation": "Renamed"}); code != 200 {
		t.Fatal(code)
	}
	if code := requestLogs(t, c, server.URL, a.Browser, ls); code != 401 {
		t.Fatal("removed credential accepted", code)
	}
	if code := requestLogs(t, c, server.URL, b.Browser, ls); code != 200 {
		t.Fatal("other app affected", code)
	}
}
func TestInvitationsManagedAppCollisionAndRoutes(t *testing.T) {
	m, _, server, owner := managedFixture(t)
	a := mustRegister(t, m, "Alpha")
	b := mustRegister(t, m, "Beta")
	csrf := loginTest(t, owner, server.URL)
	create := func(a AppEntry) (int64, string) {
		code, raw := appRequest(t, owner, server.URL, "POST", "/api/apps/"+a.ID+"/invitations", csrf, map[string]any{"hours": 1})
		if code != 200 {
			t.Fatal(code, string(raw))
		}
		var v struct {
			Invite Invitation
			URL    string
		}
		json.Unmarshal(raw, &v)
		return v.Invite.ID, strings.Split(v.URL, "#invite=")[1]
	}
	ai, at := create(a)
	bi, bt := create(b)
	if ai != bi {
		t.Fatal("fixture did not collide")
	}
	redeem := func(token string) (*http.Client, string) {
		jar, _ := cookiejar.New(nil)
		c := &http.Client{Jar: jar, Timeout: 10 * time.Second}
		code, raw := appRequest(t, c, server.URL, "POST", "/api/invitations/redeem", "", map[string]string{"token": token})
		if code != 200 {
			t.Fatal(code, string(raw))
		}
		var v map[string]string
		json.Unmarshal(raw, &v)
		return c, v["csrf"]
	}
	av, ac := redeem(at)
	bv, _ := redeem(bt)
	if code, _ := appRequest(t, av, server.URL, "POST", "/api/invitations/redeem", "", map[string]string{"token": at}); code != 401 {
		t.Fatal("one-use replay succeeded")
	}
	for _, route := range []string{"snapshot", "timeline?window=24h", "traces", "traces/0123456789abcdef0123456789abcdef", "errors", "definitions", "views", "events", "context", "invitations", "connection", "rename", "rotate", "remove"} {
		method := "GET"
		if route == "rotate" || route == "rename" || route == "remove" {
			method = "POST"
		}
		if code, _ := appRequest(t, av, server.URL, method, "/api/apps/"+b.ID+"/"+route, ac, map[string]string{}); code != 403 {
			t.Fatalf("forged route %s=%d", route, code)
		}
	}
	for _, path := range []string{"/api/apps", "/api/apps/" + a.ID + "/connection", "/api/apps/" + a.ID + "/invitations"} {
		if code, _ := appRequest(t, av, server.URL, "GET", path, "", nil); code != 403 {
			t.Fatal("viewer administration", path, code)
		}
	}
	stream, e := av.Get(server.URL + "/api/apps/" + a.ID + "/events")
	if e != nil {
		t.Fatal(e)
	}
	reader := bufio.NewReader(stream.Body)
	for i := 0; i < 4; i++ {
		line, e := reader.ReadString('\n')
		if e != nil {
			t.Fatal(e)
		}
		if i == 1 && !strings.HasPrefix(line, "id: "+a.ID+":"+a.Generation+":") {
			t.Fatal("unqualified cursor")
		}
	}
	code, _ := appRequest(t, owner, server.URL, "POST", "/api/apps/"+a.ID+"/invitations/revoke", csrf, map[string]int64{"id": ai})
	if code != 200 {
		t.Fatal(code)
	}
	done := make(chan struct{})
	go func() { io.Copy(io.Discard, reader); close(done) }()
	select {
	case <-done:
	case <-time.After(6 * time.Second):
		t.Fatal("revoked stream remained open")
	}
	stream.Body.Close()
	if code, _ := appRequest(t, av, server.URL, "GET", "/api/session", "", nil); code != 401 {
		t.Fatal("revoked session remained valid")
	}
	if code, _ := appRequest(t, bv, server.URL, "GET", "/api/apps/"+b.ID+"/snapshot", "", nil); code != 200 {
		t.Fatal("collision revoked other app", code)
	}
	r, _ := http.NewRequest("GET", server.URL+"/api/apps/"+b.ID+"/events", nil)
	r.Header.Set("Last-Event-ID", a.ID+":"+a.Generation+":1")
	resp, e := bv.Do(r)
	if e != nil {
		t.Fatal(e)
	}
	resp.Body.Close()
	if resp.StatusCode != 400 {
		t.Fatal("cross-app cursor accepted")
	}
	if e = m.Delete(b.ID, b.Name); e != nil {
		t.Fatal(e)
	}
	if code, _ := appRequest(t, bv, server.URL, "GET", "/api/session", "", nil); code != 401 {
		t.Fatal("deleted app viewer active")
	}
}
func TestManagedIntakeFinalAuthorityAndSharedBudgets(t *testing.T) {
	m, s, server, c := managedFixture(t)
	a := mustRegister(t, m, "Alpha")
	b := mustRegister(t, m, "Beta")
	// A stale forwarded request is independently rejected by the loopback consumer.
	ctx := client.NewContext(context.Background(), client.Info{Metadata: client.NewMetadata(map[string][]string{"authorization": {"Bearer " + a.Server}, "x-sentrame-app": {b.ID}})})
	if e := m.Update(a.ID, "", "server"); e != nil {
		t.Fatal(e)
	}
	if e := s.consumeLogs(ctx, logFixture("stale")); e == nil {
		t.Fatal("consumer trusted stale credentials/header")
	}
	csrf := loginTest(t, c, server.URL)
	for i := 0; i < 8; i++ {
		s.queries <- struct{}{}
	}
	if code, _ := appRequest(t, c, server.URL, "GET", "/api/apps/"+b.ID+"/snapshot", csrf, nil); code != 503 {
		t.Fatal("global query limit missing")
	}
	for i := 0; i < 8; i++ {
		<-s.queries
	}
	for i := 0; i < 8; i++ {
		s.slots <- struct{}{}
	}
	if code := requestLogs(t, c, server.URL, b.Browser, logFixture("budget")); code != 503 {
		t.Fatal("global intake limit missing")
	}
	for i := 0; i < 8; i++ {
		<-s.slots
	}
	// Hold the store writer: intake owns registry authority until commit; rotation
	// cannot report success before that request has completed its commit.
	store := m.stores[b.ID]
	store.Mu.Lock()
	ingested := make(chan int, 1)
	go func() { ingested <- requestLogs(t, c, server.URL, b.Browser, logFixture("inflight")) }()
	// Wait for the authority lock rather than assuming scheduling order.
	deadline := time.Now().Add(time.Second)
	for m.mu.TryLock() {
		m.mu.Unlock()
		if time.Now().After(deadline) {
			store.Mu.Unlock()
			t.Fatal("intake never acquired authority")
		}
		time.Sleep(time.Millisecond)
	}
	rotated := make(chan error, 1)
	go func() { rotated <- m.Update(b.ID, "", "browser") }()
	select {
	case <-rotated:
		store.Mu.Unlock()
		t.Fatal("rotation overtook commit")
	case <-time.After(20 * time.Millisecond):
	}
	store.Mu.Unlock()
	if code := <-ingested; code != 200 {
		t.Fatal(code)
	}
	if e := <-rotated; e != nil {
		t.Fatal(e)
	}
	if code := requestLogs(t, c, server.URL, b.Browser, logFixture("after")); code != 401 {
		t.Fatal("old credential committed after successful rotation")
	}
}
func TestManagedWireMatrixRestartAndTwoApps(t *testing.T) {
	path := filepath.Join(t.TempDir(), "unused.db")
	m, e := OpenRegistry(path)
	if e != nil {
		t.Fatal(e)
	}
	a := mustRegister(t, m, "Alpha")
	b := mustRegister(t, m, "Beta")
	type request struct {
		signal, ct, compression, token string
		body                           []byte
	}
	requests := []request{}
	for _, role := range []string{"browser", "server"} {
		token := a.Browser
		if role == "server" {
			token = a.Server
		}
		for _, signal := range []string{"logs", "traces", "metrics"} {
			for _, encoding := range []string{"json", "protobuf"} {
				for _, compression := range []string{"", "gzip"} {
					body := protocolBody(t, signal, encoding, fmt.Sprintf("matrix.%s.%s.%s.%s", role, signal, encoding, compression))
					ct := "application/x-protobuf"
					if encoding == "json" {
						ct = "application/json"
					}
					if compression == "gzip" {
						var zipped bytes.Buffer
						w := gzip.NewWriter(&zipped)
						w.Write(body)
						w.Close()
						body = zipped.Bytes()
					}
					requests = append(requests, request{signal, ct, compression, token, body})
				}
			}
		}
	}
	counts := map[string]int{}
	for restart := 0; restart < 2; restart++ {
		s, e := NewManaged(m, "synthetic owner passphrase", "")
		if e != nil {
			t.Fatal(e)
		}
		if e = s.StartReceiver(); e != nil {
			t.Fatal(e)
		}
		server := httptest.NewServer(s.Handler(""))
		s.Origin = server.URL
		for repeat := 0; repeat < 2; repeat++ {
			for _, r := range requests {
				code, _ := protocolRequest(t, server.Client(), server.URL, r.signal, r.ct, r.compression, r.token, r.body)
				if code != 200 {
					t.Fatalf("wire %s/%s/%s=%d", r.signal, r.ct, r.compression, code)
				}
			}
		}
		// Identical payload IDs also belong independently to B, attributed by token.
		for _, r := range requests {
			token := b.Browser
			if r.token == a.Server {
				token = b.Server
			}
			code, _ := protocolRequest(t, server.Client(), server.URL, r.signal, r.ct, r.compression, token, r.body)
			if code != 200 {
				t.Fatal(code)
			}
		}
		for _, app := range []AppEntry{a, b} {
			es, e := m.stores[app.ID].Evidence()
			if e != nil || len(es) == 0 {
				t.Fatal("no durable evidence")
			}
			stable := 0
			for _, record := range es {
				if record.Signal != "metrics" {
					stable++
				}
			}
			if restart == 0 {
				counts[app.ID] = stable
			} else if counts[app.ID] != stable {
				t.Fatal("restart replay changed stable identities")
			}
		}
		server.Close()
		s.Close()
		m.Close()
		if restart == 0 {
			m, e = OpenRegistry(path)
			if e != nil {
				t.Fatal(e)
			}
		}
	}
}

func TestManagedRejectsStaleRegistrationAndCredentialProjection(t *testing.T) {
	m, _, server, c := managedFixture(t)
	a := mustRegister(t, m, "Alpha")
	b := mustRegister(t, m, "Beta")
	loginTest(t, c, server.URL)
	for _, target := range []string{"snapshot", "traces", "views", "events", "context"} {
		r, _ := http.NewRequest("GET", server.URL+"/api/apps/"+a.ID+"/"+target+"?generation="+b.Generation, nil)
		resp, e := c.Do(r)
		if e != nil {
			t.Fatal(e)
		}
		resp.Body.Close()
		if resp.StatusCode != 403 {
			t.Fatalf("stale generation %s=%d", target, resp.StatusCode)
		}
	}
	ls := logFixture(b.Server)
	if code := requestLogs(t, c, server.URL, a.Browser, ls); code != 400 {
		t.Fatal("known app credential entered telemetry", code)
	}
	es, _ := m.stores[a.ID].Evidence()
	if len(es) != 0 {
		t.Fatal("secret-bearing batch committed")
	}
}
