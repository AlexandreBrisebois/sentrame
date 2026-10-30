package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestMappingPreviewExamplesUseExactContextAndRetainedConditionValues(t *testing.T) {
	m, _, server, c := managedFixture(t)
	a := mustRegister(t, m, "Preview app")
	b := mustRegister(t, m, "Other app")
	csrf := loginTest(t, c, server.URL)
	send := func(app AppEntry, service, stage string) {
		t.Helper()
		payload := fmt.Sprintf(`{"resourceLogs":[{"resource":{"attributes":[{"key":"service.name","value":{"stringValue":%q}}]},"scopeLogs":[{"logRecords":[{"timeUnixNano":"%d","eventName":"pwa.app.opened","attributes":[{"key":"stage.name","value":{"stringValue":%q}},{"key":"private.user.text","value":{"stringValue":"must not be offered"}}]}]}]}]}`, service, time.Now().Add(-time.Second).UnixNano(), stage)
		if code, _ := protocolRequest(t, c, server.URL, "logs", "application/json", "", app.Browser, []byte(payload)); code != 200 {
			t.Fatal(code)
		}
	}
	send(a, "wfs-pwa-browser", "search")
	send(a, "wfs-pwa-browser", "home")
	send(a, "another-service", "capture")
	send(b, "wfs-pwa-browser", "profile")
	definition := Definition{Mappings: []Mapping{{Event: "pwa.app.opened", Step: "Open", Kind: "step", Service: "wfs-pwa-browser", Source: "browser", Carrier: "log"}}}
	read := func() map[string]any {
		t.Helper()
		code, raw := appRequest(t, c, server.URL, "POST", "/api/apps/"+a.ID+"/preview", csrf, definition)
		var response map[string]any
		if code != http.StatusOK || json.Unmarshal(raw, &response) != nil {
			t.Fatalf("preview: %d %s", code, raw)
		}
		previews, ok := response["mappingPreviews"].([]any)
		if !ok || len(previews) != len(definition.Mappings) {
			t.Fatalf("missing exact-context examples: %s", raw)
		}
		return previews[0].(map[string]any)
	}
	preview := read()
	if preview["candidates"] != float64(2) || preview["matches"] != float64(2) {
		t.Fatalf("wrong contextual counts: %v", preview)
	}
	options := preview["conditions"].([]any)
	if len(options) != 2 {
		t.Fatalf("expected two retained stage choices: %v", options)
	}
	for _, v := range options {
		o := v.(map[string]any)
		if o["key"] != "stage.name" || (o["value"] != "home" && o["value"] != "search") || o["count"] != float64(1) {
			t.Fatalf("unusable or cross-context condition: %v", o)
		}
	}
	examples := preview["matchingExamples"].([]any)
	if len(examples) != 2 {
		t.Fatal("matching examples missing")
	}
	for _, v := range examples {
		e := v.(map[string]any)
		if e["service"] != "wfs-pwa-browser" {
			t.Fatal("wrong service example")
		}
	}
	definition.Mappings = append(definition.Mappings, Mapping{Event: "pwa.app.opened", Step: "Other opening", Kind: "step", Service: "another-service", Source: "browser", Carrier: "log"})
	preview = read()
	if preview["matches"] != float64(2) || preview["missingContext"] != float64(2) {
		t.Fatal("per-mapping preview mixed another service with the same event name")
	}
	definition.Mappings = definition.Mappings[:1]
	definition.Mappings[0].Key = "stage.name"
	definition.Mappings[0].Value = "search"
	preview = read()
	if preview["matches"] != float64(1) || len(preview["matchingExamples"].([]any)) != 1 || len(preview["conditions"].([]any)) != 2 {
		t.Fatal("condition must narrow matches while preserving alternative received values")
	}
	definition.Mappings[0].Value = "not-emitted"
	preview = read()
	if preview["matches"] != float64(0) || len(preview["matchingExamples"].([]any)) != 0 || len(preview["examples"].([]any)) != 2 {
		t.Fatal("zero-match condition must still offer contextual examples")
	}
	definition.Mappings[0].Service = "not-a-service"
	preview = read()
	if preview["candidates"] != float64(0) || len(preview["conditions"].([]any)) != 0 {
		t.Fatal("wrong service must not offer unrelated fields")
	}
}

func TestMappingPreviewBoundsAndUsableConditions(t *testing.T) {
	records := []Evidence{}
	for i := 0; i < 80; i++ {
		records = append(records, Evidence{ID: int64(i), Name: "opened", Source: "browser", Signal: "logs", Service: "pwa", Attrs: map[string]string{"stage.name": fmt.Sprintf("screen-%d", i), "private.user.text": "not selectable", "http.request.method": "GET?secret"}})
	}
	preview := previewMappingEvidence(records, Mapping{Event: "opened", Service: "pwa"})
	if preview.Candidates != 80 || preview.Matches != 80 || len(preview.Examples) != 5 || len(preview.MatchingExamples) != 5 || len(preview.Conditions) != 64 || !preview.ConditionsTruncated {
		t.Fatalf("preview limits not disclosed: %+v", preview)
	}
	if preview.Examples[0].ID != 79 {
		t.Fatal("examples must prefer recent receipts")
	}
	for _, option := range preview.Conditions {
		if option.Key != "stage.name" {
			t.Fatalf("unusable condition offered: %+v", option)
		}
	}
}
