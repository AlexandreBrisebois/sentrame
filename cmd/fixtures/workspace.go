package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// generateWorkspace uses the same protobuf generator as the release demonstrations.
// Its files and scope are separate so it cannot replace existing replay/followup state.
func generateWorkspace() error {
	raw, err := os.ReadFile("rest/fixtures/workspace.json")
	if err != nil {
		return err
	}
	var specs []Scenario
	if err = json.Unmarshal(raw, &specs); err != nil {
		return err
	}
	g, err := build(specs, time.Now())
	if err != nil {
		return err
	}
	if err = os.MkdirAll("rest/generated", 0700); err != nil {
		return err
	}
	for name, body := range g.Bodies {
		if err = os.WriteFile(filepath.Join("rest/generated", name+".bin"), body, 0600); err != nil {
			return err
		}
		if err = os.WriteFile(filepath.Join("rest/generated", name+".json"), g.JSON[name], 0600); err != nil {
			return err
		}
	}
	manifest, _ := json.MarshalIndent(g, "", "  ")
	if err = os.WriteFile("rest/generated/workspace-manifest.json", manifest, 0600); err != nil {
		return err
	}
	pass, err := envRequired("SENTRAME_PASSPHRASE")
	if err != nil {
		return err
	}
	token, err := envRequired("SENTRAME_BROWSER_TOKEN")
	if err != nil {
		return err
	}
	login, _ := json.Marshal(map[string]string{"passphrase": pass})
	if err = os.WriteFile("rest/generated/workspace-login.json", login, 0600); err != nil {
		return err
	}
	// Keep the existing demonstration's variables and overwrite only shared local
	// endpoint/token values plus this new scope. Never print credentials.
	values := map[string]string{"baseUrl": baseURL(), "browserToken": token, "workspaceScope": g.Scopes["workspace"]}
	existing, err := os.ReadFile("rest/.env")
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	lines := []string{}
	for _, line := range strings.Split(strings.TrimSpace(string(existing)), "\n") {
		key, _, ok := strings.Cut(line, "=")
		if _, replace := values[key]; ok && replace {
			continue
		}
		if line != "" {
			lines = append(lines, line)
		}
	}
	keys := []string{}
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		lines = append(lines, key+"="+values[key])
	}
	if err = os.WriteFile("rest/.env", []byte(strings.Join(lines, "\n")+"\n"), 0600); err != nil {
		return err
	}
	fmt.Println("Generated workspace OTLP protobuf and readable projection; 21 closed journeys plus one unmapped event. Existing demo replay/followup files preserved.")
	return nil
}

func verifyWorkspace() error {
	raw, err := os.ReadFile("rest/generated/workspace-manifest.json")
	if err != nil {
		return err
	}
	var g Generated
	if err = json.Unmarshal(raw, &g); err != nil {
		return err
	}
	token, err := envRequired("SENTRAME_BROWSER_TOKEN")
	if err != nil {
		return err
	}
	requests, err := parseREST("rest/workspace.http", map[string]string{"baseUrl": baseURL(), "browserToken": token, "workspaceScope": g.Scopes["workspace"]})
	if err != nil {
		return err
	}
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar, Timeout: 20 * time.Second}
	var snapshot, timeline []byte
	for i, r := range requests {
		req, err := http.NewRequest(r.Method, r.URL, bytes.NewReader(r.Body))
		if err != nil {
			return err
		}
		req.Header = r.Headers.Clone()
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("workspace request%d failed: %w", i+1, err)
		}
		body, readErr := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
		resp.Body.Close()
		if readErr != nil {
			return readErr
		}
		if resp.StatusCode != r.Expected {
			return fmt.Errorf("workspace request%d %s %s status%d, expected%d", i+1, r.Method, req.URL.Path, resp.StatusCode, r.Expected)
		}
		if r.Method == "GET" && req.URL.Path == "/api/snapshot" {
			snapshot = body
		}
		if r.Method == "GET" && req.URL.Path == "/api/timeline" {
			timeline = body
		}
		fmt.Printf("PASS workspace %02d %s %s → %d\n", i+1, r.Method, req.URL.Path, resp.StatusCode)
	}
	if err = checkWorkspaceSnapshot(snapshot, g); err != nil {
		return err
	}
	var histogram struct {
		Total     int
		Truncated bool
		Unit      string
		Bins      []struct{ Count int }
	}
	if err = json.Unmarshal(timeline, &histogram); err != nil {
		return err
	}
	sum := 0
	for _, bin := range histogram.Bins {
		sum += bin.Count
	}
	if histogram.Total != 44 || sum != 44 || len(histogram.Bins) != 60 || histogram.Truncated || histogram.Unit != "retained records by occurrence time" {
		return errors.New("workspace histogram must contain44 occurrence-time record counts across60 bins without truncation")
	}
	fmt.Println("PASS workspace timeline:44 records across60 occurrence-time bins in independent1h context.")
	fmt.Println("PASS workspace cohort:20 common paths,1 rare path (1/21=4.8%),1 unmapped event; exact replay preserved44 records. Verified literal rest/workspace.http through actual OTLP/storage/derivation/protected read paths.")
	return nil
}

func checkWorkspaceSnapshot(raw []byte, g Generated) error {
	var snap struct {
		Evidence   []struct{ Name, Scope string } `json:"evidence"`
		Discovered []struct{ Name string }        `json:"discovered"`
		Model      struct {
			Journeys []struct {
				ID, Scope, Status, Version string
				Definition                 int64
				Partial, RightCensored     bool
				Path                       []string
			} `json:"journeys"`
			Unassigned []struct{ Name, Scope string } `json:"unassigned"`
		} `json:"model"`
	}
	if err := json.Unmarshal(raw, &snap); err != nil {
		return err
	}
	if len(snap.Evidence) != 44 || len(snap.Model.Journeys) != 21 {
		return fmt.Errorf("workspace requires44 retained records and21 journeys; got%d/%d. Check active mappings/goal without resetting user configuration", len(snap.Evidence), len(snap.Model.Journeys))
	}
	counts := map[string]int{}
	definition := int64(0)
	for _, j := range snap.Model.Journeys {
		if j.Scope != g.Scopes["workspace"] || j.Version != "synthetic-a1" || j.Status != "completed" || j.Partial || j.RightCensored {
			return errors.New("workspace cohort must be closed,uncensored,comparable synthetic-a1 journeys; review active mappings/goal")
		}
		if definition == 0 {
			definition = j.Definition
		}
		if definition != j.Definition {
			return errors.New("workspace cohort crossed definition versions")
		}
		counts[strings.Join(j.Path, " → ")]++
	}
	if len(counts) != 2 || counts["Discovery → Completion"] != 20 || counts["Discovery → Intent → Completion"] != 1 {
		return fmt.Errorf("workspace paths differ from approved fixture mapping: %v; current configuration was not changed", counts)
	}
	found := false
	for _, event := range snap.Model.Unassigned {
		if event.Name == "workspace.hint.shown" && event.Scope == g.Scopes["workspace"] {
			found = true
		}
	}
	discovered := false
	for _, event := range snap.Discovered {
		if event.Name == "workspace.hint.shown" {
			discovered = true
		}
	}
	if !found || !discovered {
		return errors.New("workspace unmapped event must remain accepted-unassigned and discovered")
	}
	return nil
}
