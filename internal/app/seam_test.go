package app

import (
	"bytes"
	"context"
	"encoding/json"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/plog/plogotlp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	otellog "go.opentelemetry.io/otel/log"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func fixtureService(t *testing.T) (*Service, *httptest.Server, *http.Client) {
	t.Helper()
	s, e := Open(filepath.Join(t.TempDir(), "db"))
	if e != nil {
		t.Fatal(e)
	}
	svc, e := New(s, "synthetic owner passphrase", "browser-test", "server-test", "")
	if e != nil {
		t.Fatal(e)
	}
	if e = svc.StartReceiver(); e != nil {
		t.Fatal(e)
	}
	server := httptest.NewServer(svc.Handler(""))
	svc.Origin = server.URL
	jar, _ := cookiejar.New(nil)
	c := &http.Client{Jar: jar, Timeout: 10 * time.Second}
	t.Cleanup(func() { server.Close(); svc.Close(); s.DB.Close() })
	return svc, server, c
}
func loginTest(t *testing.T, c *http.Client, url string) string {
	t.Helper()
	r, _ := http.NewRequest("POST", url+"/api/login", strings.NewReader(`{"passphrase":"synthetic owner passphrase"}`))
	r.Header.Set("Origin", url)
	resp, e := c.Do(r)
	if e != nil {
		t.Fatal(e)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("login %d", resp.StatusCode)
	}
	var v map[string]string
	json.NewDecoder(resp.Body).Decode(&v)
	return v["csrf"]
}
func requestLogs(t *testing.T, c *http.Client, url, token string, ls plog.Logs) int {
	t.Helper()
	body, e := plogotlp.NewExportRequestFromLogs(ls).MarshalProto()
	if e != nil {
		t.Fatal(e)
	}
	r, _ := http.NewRequest("POST", url+"/v1/logs", bytes.NewReader(body))
	r.Header.Set("Authorization", "Bearer "+token)
	r.Header.Set("Content-Type", "application/x-protobuf")
	resp, e := c.Do(r)
	if e != nil {
		t.Fatal(e)
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)
	return resp.StatusCode
}
func logFixture(names ...string) plog.Logs {
	ls := plog.NewLogs()
	rs := ls.ResourceLogs().AppendEmpty()
	rs.Resource().Attributes().PutStr("service.name", "synthetic")
	ss := rs.ScopeLogs().AppendEmpty()
	for i, n := range names {
		l := ss.LogRecords().AppendEmpty()
		l.SetEventName(n)
		l.SetTimestamp(pcommon.NewTimestampFromTime(time.Now().Add(time.Duration(i) * time.Millisecond)))
		l.Attributes().PutStr("event.id", n)
		l.Attributes().PutStr("session.id", "random-scope")
		l.Attributes().PutStr("tab.id", "random-tab")
	}
	return ls
}
func TestStandardExporterDurableProtectedRead(t *testing.T) {
	svc, server, c := fixtureService(t)
	ctx := context.Background()
	exp, e := otlploghttp.New(ctx, otlploghttp.WithEndpointURL(server.URL+"/v1/logs"), otlploghttp.WithHeaders(map[string]string{"Authorization": "Bearer browser-test"}), otlploghttp.WithInsecure())
	if e != nil {
		t.Fatal(e)
	}
	provider := sdklog.NewLoggerProvider(sdklog.WithProcessor(sdklog.NewSimpleProcessor(exp)))
	var rec otellog.Record
	rec.SetEventName("standard.received")
	rec.SetTimestamp(time.Now())
	rec.SetBody(attribute.StringValue("PRIVATE BODY must never persist"))
	rec.AddAttributes(attribute.String("secret.password", "PRIVATE ATTRIBUTE"))
	provider.Logger("seam").Emit(ctx, rec)
	if e = provider.Shutdown(ctx); e != nil {
		t.Fatal(e)
	}
	resp, e := c.Get(server.URL + "/api/snapshot")
	if e != nil {
		t.Fatal(e)
	}
	resp.Body.Close()
	if resp.StatusCode != 401 {
		t.Fatalf("unauthorized %d", resp.StatusCode)
	}
	loginTest(t, c, server.URL)
	resp, e = c.Get(server.URL + "/api/snapshot")
	if e != nil {
		t.Fatal(e)
	}
	b, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	if !bytes.Contains(b, []byte("standard.received")) || bytes.Contains(b, []byte("PRIVATE")) {
		t.Fatalf("protected sanitized read: %s", b)
	}

	actualPath := svc.Store.Path
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
	again := httptest.NewServer(restarted.Handler(""))
	defer again.Close()
	restarted.Origin = again.URL
	jar, _ := cookiejar.New(nil)
	after := &http.Client{Jar: jar, Timeout: 10 * time.Second}
	loginTest(t, after, again.URL)
	rr, e := after.Get(again.URL + "/api/snapshot")
	if e != nil {
		t.Fatal(e)
	}
	defer rr.Body.Close()
	persisted, _ := io.ReadAll(rr.Body)
	if !bytes.Contains(persisted, []byte("standard.received")) || bytes.Contains(persisted, []byte("PRIVATE")) {
		t.Fatalf("same actual intake database did not survive restart: %s", persisted)
	}

}
func TestAtomicReplayAuthorityAndFailedCommit(t *testing.T) {
	svc, server, c := fixtureService(t)
	ls := logFixture("discovery.entered", "completion.entered")
	if status := requestLogs(t, c, server.URL, "browser-test", ls); status != 200 {
		t.Fatalf("intake %d", status)
	}
	if status := requestLogs(t, c, server.URL, "browser-test", ls); status != 200 {
		t.Fatalf("replay %d", status)
	}
	es, _ := svc.Store.Evidence()
	if len(es) != 2 {
		t.Fatalf("replay grew %d", len(es))
	}
	if status := requestLogs(t, c, server.URL, "server-test", ls); status != 200 {
		t.Fatal(status)
	}
	es, _ = svc.Store.Evidence()
	if len(es) != 4 {
		t.Fatal("namespaces collapsed")
	}
	bad := logFixture("new.valid", "old.invalid")
	bad.ResourceLogs().At(0).ScopeLogs().At(0).LogRecords().At(1).SetTimestamp(pcommon.NewTimestampFromTime(time.Now().Add(-31 * 24 * time.Hour)))
	if status := requestLogs(t, c, server.URL, "browser-test", bad); status != 400 {
		t.Fatalf("atomic status %d", status)
	}
	es, _ = svc.Store.Evidence()
	if len(es) != 4 {
		t.Fatal("mixed batch partially committed")
	}
	ls.ResourceLogs().At(0).ScopeLogs().At(0).LogRecords().At(0).SetEventName("conflict")
	if status := requestLogs(t, c, server.URL, "browser-test", ls); status != 400 {
		t.Fatalf("conflict %d", status)
	}
	svc.Store.DB.Close()
	if status := requestLogs(t, c, server.URL, "browser-test", logFixture("after.close")); status != 503 {
		t.Fatalf("failed commit %d", status)
	}
}
func TestOwnerCSRFAndSafeMutations(t *testing.T) {
	_, server, c := fixtureService(t)
	csrf := loginTest(t, c, server.URL)
	r, _ := http.NewRequest("POST", server.URL+"/api/delete", strings.NewReader(`{"scope":"random-scope"}`))
	resp, e := c.Do(r)
	if e != nil {
		t.Fatal(e)
	}
	resp.Body.Close()
	if resp.StatusCode != 403 {
		t.Fatal("mutation accepted without origin/csrf")
	}
	r, _ = http.NewRequest("POST", server.URL+"/api/delete", strings.NewReader(`{"scope":"random-scope"}`))
	r.Header.Set("Origin", server.URL)
	r.Header.Set("X-CSRF-Token", csrf)
	resp, e = c.Do(r)
	if e != nil {
		t.Fatal(e)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("valid mutation %d", resp.StatusCode)
	}
}
