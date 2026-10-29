package app

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

func TestDeleteReferenceIsProtectedAndPreservesEvidenceAndOtherViews(t *testing.T) {
	path := filepath.Join(t.TempDir(), "db")
	store, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { store.DB.Close() }()
	svc, err := New(store, "synthetic owner passphrase", "browser-test", "server-test", "http://localhost")
	if err != nil {
		t.Fatal(err)
	}
	h := svc.Handler("")
	login := httptest.NewRequest("POST", "http://localhost/api/login", strings.NewReader(`{"passphrase":"synthetic owner passphrase"}`))
	login.Header.Set("Origin", svc.Origin)
	logged := httptest.NewRecorder()
	h.ServeHTTP(logged, login)
	if logged.Code != 200 {
		t.Fatal(logged.Body)
	}
	cookie := logged.Result().Cookies()[0]
	var session map[string]string
	json.Unmarshal(logged.Body.Bytes(), &session)
	request := func(method, path, body string, auth, proof bool) *httptest.ResponseRecorder {
		r := httptest.NewRequest(method, "http://localhost"+path, strings.NewReader(body))
		if auth {
			r.AddCookie(cookie)
		}
		if proof {
			r.Header.Set("Origin", svc.Origin)
			r.Header.Set("X-CSRF-Token", session["csrf"])
		}
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		return w
	}
	var ids []int64
	for _, name := range []string{"remove me", "keep me"} {
		w := request("POST", "/api/views", fmt.Sprintf(`{"name":%q,"mode":"overview"}`, name), true, true)
		if w.Code != 200 {
			t.Fatal(w.Code, w.Body)
		}
		var v SavedView
		json.Unmarshal(w.Body.Bytes(), &v)
		ids = append(ids, v.ID)
	}
	if err := store.Insert([]Evidence{{Source: "server", Signal: "logs", Name: "synthetic.retained", Time: 1, Attrs: map[string]string{}}}); err != nil {
		t.Fatal(err)
	}
	body := fmt.Sprintf(`{"id":%d}`, ids[0])
	for _, tc := range []struct {
		method, body string
		auth, proof  bool
		want         int
	}{
		{"POST", body, false, false, 401}, {"POST", body, true, false, 403}, {"GET", "", true, true, 405}, {"POST", `{"id":0}`, true, true, 400}, {"POST", `{"id":99999}`, true, true, 404},
	} {
		if w := request(tc.method, "/api/views/delete", tc.body, tc.auth, tc.proof); w.Code != tc.want {
			t.Fatalf("%+v: got %d %s", tc, w.Code, w.Body)
		}
	}
	if w := request("POST", "/api/views/delete", body, true, true); w.Code != 200 {
		t.Fatal(w.Code, w.Body)
	}
	views, err := svc.getViews()
	if err != nil || len(views) != 1 || views[0].ID != ids[1] {
		t.Fatalf("views %+v error %v", views, err)
	}
	var count int
	if err := store.DB.QueryRow(`SELECT count(*) FROM evidence`).Scan(&count); err != nil || count != 1 {
		t.Fatalf("evidence changed: %d %v", count, err)
	}
	if w := request("POST", "/api/views/delete", body, true, true); w.Code != 404 {
		t.Fatal("repeat delete", w.Code)
	}
	store.DB.Close()
	store, err = Open(path)
	if err != nil {
		t.Fatal(err)
	}
	svc.Store = store
	views, err = svc.getViews()
	if err != nil || len(views) != 1 || views[0].ID != ids[1] {
		t.Fatalf("restart views %+v error %v", views, err)
	}
}
