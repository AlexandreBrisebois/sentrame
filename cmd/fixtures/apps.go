package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"
	"sentrame/internal/app"
	"strings"
	"time"
)

// Synthetic tools register/select one explicit demonstration app through the same
// owner APIs as the UI. Generated credentials configure the fixture sender only.
func prepareDemo() error {
	jar, _ := cookiejar.New(nil)
	c := &http.Client{Jar: jar, Timeout: 20 * time.Second}
	csrf := ""
	call := func(method, path string, body any, out any) error {
		b, _ := json.Marshal(body)
		r, _ := http.NewRequest(method, baseURL()+path, bytes.NewReader(b))
		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("Origin", baseURL())
		r.Header.Set("X-CSRF-Token", csrf)
		resp, e := c.Do(r)
		if e != nil {
			return errors.New("demonstration requires a running loopback service")
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			io.Copy(io.Discard, resp.Body)
			return errors.New("demonstration app setup rejected; verify owner access and private registry")
		}
		if out != nil {
			return json.NewDecoder(resp.Body).Decode(out)
		}
		io.Copy(io.Discard, resp.Body)
		return nil
	}
	var session map[string]string
	if e := call("POST", "/api/login", map[string]string{"passphrase": os.Getenv("SENTRAME_PASSPHRASE")}, &session); e != nil {
		return e
	}
	csrf = session["csrf"]
	var apps []app.AppEntry
	if e := call("GET", "/api/apps", nil, &apps); e != nil {
		return e
	}
	id := ""
	for _, a := range apps {
		if a.Name == "Sentrame demonstration" && a.State == "active" {
			id = a.ID
		}
	}
	if id == "" {
		var a app.AppEntry
		if e := call("POST", "/api/apps", map[string]string{"name": "Sentrame demonstration"}, &a); e != nil {
			return e
		}
		id = a.ID
	}
	var settings map[string]string
	if e := call("GET", "/api/apps/"+id+"/connection", nil, &settings); e != nil {
		return e
	}
	os.Setenv("SENTRAME_APP_ID", id)
	os.Setenv("SENTRAME_BROWSER_TOKEN", settings["browser"])
	os.Setenv("SENTRAME_SERVER_TOKEN", settings["server"])
	if os.Args[1] == "generate" || os.Args[1] == "demo" || os.Args[1] == "workspace" || os.Args[1] == "workspace-generate" || os.Args[1] == "live" {
		var defs []app.Definition
		if e := call("GET", "/api/apps/"+id+"/definitions", nil, &defs); e != nil {
			return e
		}
		seed := app.InitialDefinition()
		seed.Version = 0
		var latest app.Definition
		if len(defs) > 0 {
			latest = defs[len(defs)-1]
			latest.Version = 0
		}
		old, _ := json.Marshal(latest)
		wanted, _ := json.Marshal(seed)
		if !bytes.Equal(old, wanted) {
			if e := call("POST", "/api/apps/"+id+"/definitions", app.InitialDefinition(), nil); e != nil {
				return e
			}
		}
	}
	return nil
}
func appAPI(path string) string { return "/api/apps/" + os.Getenv("SENTRAME_APP_ID") + "/" + path }

func fixtureAPIPath(path string) string {
	prefix := "/api/apps/" + os.Getenv("SENTRAME_APP_ID") + "/"
	if strings.HasPrefix(path, prefix) {
		return "/api/" + strings.TrimPrefix(path, prefix)
	}
	return path
}
