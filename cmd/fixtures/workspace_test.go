package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sentrame/internal/app"
	"strings"
	"testing"
	"time"

	logsv1 "go.opentelemetry.io/proto/otlp/collector/logs/v1"
	"google.golang.org/protobuf/proto"
)

func TestWorkspaceFixtureCreatesComparableClosedCohortAndUnmappedDiscovery(t *testing.T) {
	raw, err := os.ReadFile("../../rest/fixtures/workspace.json")
	if err != nil {
		t.Fatal(err)
	}
	var specs []Scenario
	if err = json.Unmarshal(raw, &specs); err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 10, 28, 15, 0, 0, 0, time.UTC)
	g, err := build(specs, now)
	if err != nil {
		t.Fatal(err)
	}
	var req logsv1.ExportLogsServiceRequest
	if err = proto.Unmarshal(g.Bodies["workspace"], &req); err != nil {
		t.Fatal(err)
	}
	records := req.ResourceLogs[0].ScopeLogs[0].LogRecords
	counts := map[string]int{}
	identities := map[string]bool{}
	path := ""
	unmapped := 0
	for _, r := range records {
		if r.TimeUnixNano >= uint64(now.UnixNano()) {
			t.Fatal("workspace journeys must be closed before generation time")
		}
		record, scope := "", ""
		for _, a := range r.Attributes {
			if a.Key == "event.id" {
				record = a.Value.GetStringValue()
			}
			if a.Key == "session.id" {
				scope = a.Value.GetStringValue()
			}
		}
		if record == "" || identities[record] || scope != g.Scopes["workspace"] {
			t.Fatal("workspace identities must be fresh, unique and scoped")
		}
		identities[record] = true
		if r.EventName == "workspace.hint.shown" {
			unmapped++
			continue
		}
		if path != "" {
			path += " → "
		}
		path += r.EventName
		if r.EventName == "completion.entered" {
			counts[path]++
			path = ""
		}
	}
	if counts["discovery.entered → completion.entered"] != 20 || counts["discovery.entered → vote.requested → completion.entered"] != 1 || len(counts) != 2 || path != "" || unmapped != 1 {
		t.Fatalf("need 20 common,1 rare,1 unmapped; got%v remaining=%s unmapped=%d", counts, path, unmapped)
	}
	if req.ResourceLogs[0].Resource.Attributes[1].Value.GetStringValue() != "synthetic-a1" {
		t.Fatal("cohort version must stay comparable")
	}
}

func TestWorkspaceGenerationPreservesExistingDemoReplayAndRESTVariables(t *testing.T) {
	t.Setenv("SENTRAME_APP_ID", strings.Repeat("a", 64))
	raw, err := os.ReadFile("../../rest/fixtures/workspace.json")
	if err != nil {
		t.Fatal(err)
	}
	httpFixture, err := os.ReadFile("../../rest/workspace.http")
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	if err = os.MkdirAll(filepath.Join(dir, "rest", "fixtures"), 0700); err != nil {
		t.Fatal(err)
	}
	if err = os.MkdirAll(filepath.Join(dir, "rest", "generated"), 0700); err != nil {
		t.Fatal(err)
	}
	os.WriteFile(filepath.Join(dir, "rest", "fixtures", "workspace.json"), raw, 0600)
	os.WriteFile(filepath.Join(dir, "rest", "workspace.http"), httpFixture, 0600)
	os.WriteFile(filepath.Join(dir, "rest", "generated", "manifest.json"), []byte("existing-replay-state"), 0600)
	os.WriteFile(filepath.Join(dir, "rest", ".env"), []byte("healthyScope=preserved\nbaseUrl=http://localhost:1\n"), 0600)
	t.Chdir(dir)
	t.Setenv("SENTRAME_PASSPHRASE", "synthetic test passphrase")
	t.Setenv("SENTRAME_BROWSER_TOKEN", "synthetic-test-token")
	t.Setenv("SENTRAME_URL", "http://127.0.0.1:8096")
	if err = generateWorkspace(); err != nil {
		t.Fatal(err)
	}
	prior, err := os.ReadFile("rest/generated/manifest.json")
	if err != nil || string(prior) != "existing-replay-state" {
		t.Fatal("workspace generation replaced existing replay state")
	}
	vars, err := os.ReadFile("rest/.env")
	if err != nil || !bytes.Contains(vars, []byte("healthyScope=preserved")) || !bytes.Contains(vars, []byte("workspaceScope=")) {
		t.Fatal("workspace generation lost prior REST scope")
	}
	var g Generated
	manifest, err := os.ReadFile("rest/generated/workspace-manifest.json")
	if err != nil {
		t.Fatal(err)
	}
	if err = json.Unmarshal(manifest, &g); err != nil {
		t.Fatal(err)
	}
	requests, err := parseREST("rest/workspace.http", map[string]string{"baseUrl": baseURL(), "browserToken": "synthetic-test-token", "workspaceScope": g.Scopes["workspace"]})
	if err != nil {
		t.Fatal(err)
	}
	body, err := os.ReadFile("rest/generated/workspace.bin")
	if err != nil {
		t.Fatal(err)
	}
	if len(requests) != 5 || !bytes.Equal(requests[1].Body, body) || !bytes.Equal(requests[2].Body, body) || requests[3].URL != baseURL()+appAPI("snapshot")+"?scope="+g.Scopes["workspace"] || requests[4].URL != baseURL()+appAPI("timeline")+"?window=1h&scope="+g.Scopes["workspace"] {
		t.Fatal("versioncontrolled workspace REST requests did not preserve ingestion/replay/scoped read")
	}
	if g.Scopes["workspace"] == "" {
		t.Fatal("workspace manifest omitted generated scope")
	}
}

func TestWorkspaceSnapshotVerificationRejectsFabricatedRarityAndMissingDiscovery(t *testing.T) {
	journeys := []map[string]any{}
	for i := 0; i < 21; i++ {
		path := []string{"Discovery", "Completion"}
		if i == 20 {
			path = []string{"Discovery", "Intent", "Completion"}
		}
		journeys = append(journeys, map[string]any{"id": fmt.Sprint(i), "scope": "scope", "version": "synthetic-a1", "status": "completed", "definition": 1, "path": path})
	}
	snapshot := map[string]any{"evidence": make([]map[string]any, 44), "discovered": []map[string]string{{"name": "workspace.hint.shown"}}, "model": map[string]any{"journeys": journeys, "unassigned": []map[string]string{{"name": "workspace.hint.shown", "scope": "scope"}}}}
	g := Generated{Scopes: map[string]string{"workspace": "scope"}}
	raw, _ := json.Marshal(snapshot)
	if err := checkWorkspaceSnapshot(raw, g); err != nil {
		t.Fatal(err)
	}
	journeys[20]["path"] = []string{"Discovery", "Completion"}
	raw, _ = json.Marshal(snapshot)
	if err := checkWorkspaceSnapshot(raw, g); err == nil {
		t.Fatal("missing uncommon sequence passed verification")
	}
	journeys[20]["path"] = []string{"Discovery", "Intent", "Completion"}
	snapshot["discovered"] = []map[string]string{}
	raw, _ = json.Marshal(snapshot)
	if err := checkWorkspaceSnapshot(raw, g); err == nil {
		t.Fatal("missing discovered unmapped event passed verification")
	}
}

func TestWorkspaceUnmappedFixturePrecedesJourneyCapture(t *testing.T) {
	raw, err := os.ReadFile("../../rest/fixtures/workspace.json")
	if err != nil {
		t.Fatal(err)
	}
	var specs []Scenario
	if err = json.Unmarshal(raw, &specs); err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 10, 28, 15, 0, 0, 0, time.UTC)
	g, err := build(specs, now)
	if err != nil {
		t.Fatal(err)
	}
	var request logsv1.ExportLogsServiceRequest
	if err = proto.Unmarshal(g.Bodies["workspace"], &request); err != nil {
		t.Fatal(err)
	}
	evidence := []app.Evidence{}
	for i, r := range request.ResourceLogs[0].ScopeLogs[0].LogRecords {
		attrs := map[string]string{}
		for _, kv := range r.Attributes {
			attrs[kv.Key] = kv.Value.GetStringValue()
		}
		evidence = append(evidence, app.Evidence{ID: int64(i + 1), Name: r.EventName, Source: "browser", Signal: "logs", Scope: attrs["session.id"], Tab: attrs["tab.id"], Record: attrs["event.id"], Time: int64(r.TimeUnixNano), Receipt: now.UnixNano(), Definition: 1, Version: "synthetic-a1", Attrs: attrs})
	}
	model := app.Derive(evidence, []app.Definition{app.InitialDefinition()}, now.UnixNano())
	if len(model.Journeys) != 21 || len(model.Unassigned) != 1 || model.Unassigned[0].Name != "workspace.hint.shown" {
		t.Fatalf("actual journey derivation must keep 21 completed journeys and one unassigned example: got %d/%d", len(model.Journeys), len(model.Unassigned))
	}
}
