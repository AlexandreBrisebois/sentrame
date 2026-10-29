package app

import (
	"net/http/httptest"
	"path/filepath"
	"testing"
)

func TestReadinessRequiresStartedReceiverAndUsableStorage(t *testing.T) {
	store, err := Open(filepath.Join(t.TempDir(), "db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.DB.Close()
	svc, err := New(store, "synthetic owner passphrase", "browser-test", "server-test", "http://localhost")
	if err != nil {
		t.Fatal(err)
	}
	status := func(path string) int {
		w := httptest.NewRecorder()
		svc.Handler("").ServeHTTP(w, httptest.NewRequest("GET", path, nil))
		return w.Code
	}
	if got := status("/ready"); got != 503 {
		t.Fatalf("unstarted ready=%d", got)
	}
	svc.ReceiverURL = "http://127.0.0.1:12345"
	if got := status("/ready"); got != 200 {
		t.Fatalf("initialized ready=%d", got)
	}
	store.DB.Close()
	if got := status("/ready"); got != 503 {
		t.Fatalf("closed store ready=%d", got)
	}
	if got := status("/health"); got != 200 {
		t.Fatalf("liveness=%d", got)
	}
}
