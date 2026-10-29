package app

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestConcurrentIntakeReadCleanupAndSSERevocation(t *testing.T) {
	svc, server, c := fixtureService(t)
	csrf := loginTest(t, c, server.URL)
	stream, e := c.Get(server.URL + "/api/events")
	if e != nil {
		t.Fatal(e)
	}
	reader := bufio.NewReader(stream.Body)
	line, e := reader.ReadString('\n')
	if e != nil || line != "event: reset\n" {
		t.Fatalf("initial reset %q %v", line, e)
	}
	var wg sync.WaitGroup
	fail := make(chan string, 100)
	for worker := 0; worker < 4; worker++ {
		wg.Add(1)
		go func(worker int) {
			defer wg.Done()
			for i := 0; i < 10; i++ {
				ls := logFixture(fmt.Sprintf("concurrent.%d.%d", worker, i))
				if status := requestLogs(t, c, server.URL, "server-test", ls); status != 200 {
					fail <- fmt.Sprintf("intake %d", status)
				}
			}
		}(worker)
	}
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			resp, e := c.Get(server.URL + "/api/snapshot")
			if e != nil {
				fail <- e.Error()
				return
			}
			var snap map[string]json.RawMessage
			if e = json.NewDecoder(resp.Body).Decode(&snap); e != nil {
				fail <- e.Error()
			}
			resp.Body.Close()
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			if e := svc.Store.Cleanup(); e != nil {
				fail <- e.Error()
			}
		}
	}()
	wg.Wait()
	close(fail)
	for msg := range fail {
		t.Error(msg)
	}
	es, e := svc.Store.Evidence()
	if e != nil || len(es) != 40 {
		t.Fatalf("lost writes %d %v", len(es), e)
	}
	r, _ := http.NewRequest("POST", server.URL+"/api/logout", strings.NewReader(`{}`))
	r.Header.Set("Origin", server.URL)
	r.Header.Set("X-CSRF-Token", csrf)
	resp, e := c.Do(r)
	if e != nil {
		t.Fatal(e)
	}
	resp.Body.Close()
	done := make(chan error, 1)
	go func() { _, e := io.ReadAll(reader); done <- e }()
	select {
	case <-done:
	case <-time.After(7 * time.Second):
		t.Fatal("revoked stream remained active")
	}
	stream.Body.Close()
	resp, e = c.Get(server.URL + "/api/snapshot")
	if e != nil {
		t.Fatal(e)
	}
	resp.Body.Close()
	if resp.StatusCode != 401 {
		t.Fatal("revoked cookie retained read access")
	}
}
func TestDeleteReplayBarrierSurvivesRestartAndAppDisable(t *testing.T) {
	svc, server, c := fixtureService(t)
	csrf := loginTest(t, c, server.URL)
	ls := logFixture("discovery.entered")
	if code := requestLogs(t, c, server.URL, "browser-test", ls); code != 200 {
		t.Fatal(code)
	}
	mutate := func(body string) {
		r, _ := http.NewRequest("POST", server.URL+"/api/delete", strings.NewReader(body))
		r.Header.Set("Origin", server.URL)
		r.Header.Set("X-CSRF-Token", csrf)
		resp, e := c.Do(r)
		if e != nil {
			t.Fatal(e)
		}
		resp.Body.Close()
		if resp.StatusCode != 200 {
			t.Fatalf("delete %d", resp.StatusCode)
		}
	}
	mutate(`{"scope":"random-scope"}`)
	actualPath := svc.Store.Path
	server.Close()
	svc.Close()
	svc.Store.DB.Close()
	reopened, e := Open(actualPath)
	if e != nil {
		t.Fatal(e)
	}
	defer reopened.DB.Close()
	restarted, e := New(reopened, "synthetic owner passphrase", "browser-test", "server-test", "")
	if e != nil {
		t.Fatal(e)
	}
	if e = restarted.StartReceiver(); e != nil {
		t.Fatal(e)
	}
	defer restarted.Close()
	server = httptest.NewServer(restarted.Handler(""))
	defer server.Close()
	restarted.Origin = server.URL
	svc = restarted
	csrf = loginTest(t, c, server.URL)

	if code := requestLogs(t, c, server.URL, "browser-test", ls); code != 200 {
		t.Fatal(code)
	}
	es, _ := svc.Store.Evidence()
	if len(es) != 0 {
		t.Fatal("deleted visit resurrected")
	}
	mutate(`{"app":true}`)
	if code := requestLogs(t, c, server.URL, "browser-test", ls); code != 403 {
		t.Fatalf("disabled generation accepted %d", code)
	}
	var n int
	svc.Store.DB.QueryRow(`SELECT count(*) FROM tombstones`).Scan(&n)
	if n != 1 {
		t.Fatal("app deletion erased visit barrier")
	}
}
