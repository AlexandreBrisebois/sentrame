package main

import (
	"bytes"
	"encoding/json"
	logsv1 "go.opentelemetry.io/proto/otlp/collector/logs/v1"
	"google.golang.org/protobuf/proto"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func TestFixtureFreshIdentityAndReplay(t *testing.T) {
	specs := []Scenario{{Name: "test", Scope: "test", Producer: "browser", Events: []Event{{Name: "discovery.entered", Offset: "-1m"}}}}
	now := time.Date(2026, 10, 28, 15, 0, 0, 0, time.UTC)
	a, err := build(specs, now)
	if err != nil {
		t.Fatal(err)
	}
	if len(a.Bodies) != 1 {
		t.Fatalf("one readable fixture must create one actual protobuf body; got %d", len(a.Bodies))
	}
	b, err := build(specs, now)
	if err != nil {
		t.Fatal(err)
	}
	var req logsv1.ExportLogsServiceRequest
	if err = proto.Unmarshal(a.Bodies["test"], &req); err != nil {
		t.Fatal(err)
	}
	logs := req.ResourceLogs[0].ScopeLogs[0].LogRecords
	if len(logs) != 1 || logs[0].TimeUnixNano != uint64(now.Add(-time.Minute).UnixNano()) {
		t.Fatalf("fresh timestamp not encoded in real protobuf: %v", logs)
	}
	if a.Scopes["test"] == b.Scopes["test"] {
		t.Fatal("new generation reused logical scope")
	}
	first := append([]byte(nil), a.Bodies["test"]...)
	if string(a.Bodies["test"]) != string(first) {
		t.Fatal("replay mutated generated body")
	}
}

func TestReadableFixtureCanPreserveCrossProducerRecordIdentity(t *testing.T) {
	var specs []Scenario
	raw := `[{"name":"browser","scope":"combined","producer":"browser","events":[{"name":"vote.requested","offset":"-2s","record":"shared-occurrence","action":"one-action"}]},{"name":"server","scope":"combined","producer":"server","events":[{"name":"vote.committed","offset":"-1s","record":"shared-occurrence","action":"one-action"}]}]`
	if err := json.Unmarshal([]byte(raw), &specs); err != nil {
		t.Fatal(err)
	}
	g, err := build(specs, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	ids := map[string]string{}
	for _, name := range []string{"browser", "server"} {
		var req logsv1.ExportLogsServiceRequest
		if err = proto.Unmarshal(g.Bodies[name], &req); err != nil {
			t.Fatal(err)
		}
		for _, a := range req.ResourceLogs[0].ScopeLogs[0].LogRecords[0].Attributes {
			if a.Key == "event.id" {
				ids[name] = a.Value.GetStringValue()
			}
		}
	}
	if ids["browser"] == "" || ids["browser"] != ids["server"] {
		t.Fatalf("readable record identity must survive across producers to test namespaced replay; got browser=%s server=%s", ids["browser"], ids["server"])
	}
}

func TestRESTParserReadsExactSuppliedBodyAndHeaders(t *testing.T) {
	dir := t.TempDir()
	body := []byte{0, 1, 2, 255, 10}
	if err := os.WriteFile(filepath.Join(dir, "body.bin"), body, 0600); err != nil {
		t.Fatal(err)
	}
	text := `# @expect 400
POST {{baseUrl}}/v1/logs
Content-Type: application/x-protobuf
Authorization: Bearer {{token}}

< ./body.bin

### Independent protected read
# @expect 401
# @no-cookie-jar
GET {{baseUrl}}/api/snapshot
`
	path := filepath.Join(dir, "demo.http")
	if err := os.WriteFile(path, []byte(text), 0600); err != nil {
		t.Fatal(err)
	}
	requests, err := parseREST(path, map[string]string{"baseUrl": "http://127.0.0.1:8096", "token": "synthetic-token"})
	if err != nil {
		t.Fatal(err)
	}
	if len(requests) != 2 || requests[0].Expected != 400 || !bytes.Equal(requests[0].Body, body) || requests[0].Headers.Get("Authorization") != "Bearer synthetic-token" || !requests[1].NoCookies || requests[1].Expected != 401 {
		t.Fatalf("literal supplied REST semantics were changed: %#v", requests)
	}
}

func TestManualFollowupRefreshesTimeOnceWithoutChangingIdentity(t *testing.T) {
	created := time.Date(2026, 10, 28, 12, 0, 0, 0, time.UTC)
	specs := []Scenario{{Name: "live-followup", Scope: "live", Producer: "browser", Events: []Event{{Name: "completion.entered", Offset: "0s", Action: "persisted-action"}}}}
	generated, err := build(specs, created)
	if err != nil {
		t.Fatal(err)
	}
	original := generated.Bodies["live-followup"]
	later := created.Add(time.Hour)
	refreshed, err := refreshFollowup(original, later, false)
	if err != nil {
		t.Fatal(err)
	}
	var before, after logsv1.ExportLogsServiceRequest
	if err = proto.Unmarshal(original, &before); err != nil {
		t.Fatal(err)
	}
	if err = proto.Unmarshal(refreshed, &after); err != nil {
		t.Fatal(err)
	}
	b := before.ResourceLogs[0].ScopeLogs[0].LogRecords[0]
	a := after.ResourceLogs[0].ScopeLogs[0].LogRecords[0]
	if a.TimeUnixNano != uint64(later.UnixNano()) || a.ObservedTimeUnixNano != uint64(later.UnixNano()) {
		t.Fatalf("first manual followup must occur after current live T0, got%d want%d", a.TimeUnixNano, later.UnixNano())
	}
	if !reflect.DeepEqual(a.Attributes, b.Attributes) || a.EventName != b.EventName {
		t.Fatal("fresh manual timestamp changed scope/action/event identity")
	}
	repeated, err := refreshFollowup(refreshed, later.Add(time.Hour), true)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(repeated, refreshed) {
		t.Fatal("retry/replay must preserve exact prepared protobuf bytes")
	}
}

func TestTimingFixtureProducesTwentyDistinctMeasuredPairs(t *testing.T) {
	raw, err := os.ReadFile("../../rest/fixtures/scenarios.json")
	if err != nil {
		t.Fatal(err)
	}
	var specs []Scenario
	if err = json.Unmarshal(raw, &specs); err != nil {
		t.Fatal(err)
	}
	found := false
	for _, s := range specs {
		if s.Name == "timing-samples" {
			found = true
		}
	}
	if !found {
		t.Fatal("readable demonstration must include a timing-samples scenario")
	}
	now := time.Date(2026, 10, 28, 15, 0, 0, 0, time.UTC)
	g, err := build(specs, now)
	if err != nil {
		t.Fatal(err)
	}
	var req logsv1.ExportLogsServiceRequest
	if err = proto.Unmarshal(g.Bodies["timing-samples"], &req); err != nil {
		t.Fatal(err)
	}
	records := req.ResourceLogs[0].ScopeLogs[0].LogRecords
	if len(records) != 42 {
		t.Fatalf("timing demonstration must contain discovery,20 measured pairs,completion; got%d", len(records))
	}
	starts := map[string]uint64{}
	durations := []int{}
	for _, r := range records {
		if r.TimeUnixNano >= uint64(now.UnixNano()) {
			t.Fatal("timing fixture contains an ongoing/future sample")
		}
		occurrence := ""
		for _, a := range r.Attributes {
			if a.Key == "stage.occurrence.id" {
				occurrence = a.Value.GetStringValue()
			}
		}
		if r.EventName == "stage.entered" {
			if occurrence == "" {
				t.Fatal("stage pairing identity missing")
			}
			if _, exists := starts[occurrence]; exists {
				t.Fatal("two real stage entries reused occurrence identity")
			}
			starts[occurrence] = r.TimeUnixNano
		}
		if r.EventName == "stage.exited" {
			start, ok := starts[occurrence]
			if !ok {
				t.Fatal("stage exit lost its entry identity")
			}
			durations = append(durations, int((r.TimeUnixNano-start)/uint64(time.Second)))
		}
	}
	if len(starts) != 20 || len(durations) != 20 {
		t.Fatalf("need20 unique paired occurrences; got%d/%d", len(starts), len(durations))
	}
	for i, d := range durations {
		if d != i+1 {
			t.Fatalf("duration%d=%ds, expected%ds", i, d, i+1)
		}
	}
}
