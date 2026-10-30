package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"
)

// Standalone OTLP logs remain inspectable even when neither a native trace nor
// a correlated business journey can be derived from the producer's evidence.
func TestManagedStandaloneWFSLogsRemainVisibleWithoutJourneysOrTraces(t *testing.T) {
	m, _, server, c := managedFixture(t)
	a := mustRegister(t, m, "WFS")
	other := mustRegister(t, m, "Other app")
	csrf := loginTest(t, c, server.URL)
	payload := fmt.Sprintf(`{"resourceLogs":[{"resource":{"attributes":[{"key":"service.name","value":{"stringValue":"wfs-pwa-browser"}},{"key":"service.version","value":{"stringValue":"test-beta"}}]},"scopeLogs":[{"scope":{"name":"wfs.anonymous"},"logRecords":[{"timeUnixNano":"%d","eventName":"pwa.app.opened","severityNumber":9,"body":{"stringValue":"WFS opening observed"},"attributes":[{"key":"stage.name","value":{"stringValue":"search"}}]}]}]}]}`, time.Now().Add(-time.Second).UnixNano())
	code, response := protocolRequest(t, c, server.URL, "logs", "application/json", "", a.Browser, []byte(payload))
	if code != http.StatusOK || !json.Valid([]byte(response)) {
		t.Fatalf("OTLP log intake: status=%d response=%s", code, response)
	}
	base := "/api/apps/" + a.ID
	code, raw := appRequest(t, c, server.URL, "GET", base+"/snapshot", "", nil)
	var snap struct {
		Evidence   []Evidence
		Model      Model
		Discovered []struct {
			Name     string
			Count    int
			Assigned int
		}
	}
	if code != http.StatusOK || json.Unmarshal(raw, &snap) != nil || len(snap.Evidence) != 1 || len(snap.Model.Unassigned) != 1 || len(snap.Model.Journeys) != 0 {
		t.Fatalf("standalone log missing from protected evidence/unassigned read: %d %s", code, raw)
	}
	e := snap.Evidence[0]
	if e.Name != "pwa.app.opened" || e.Service != "wfs-pwa-browser" || e.Version != "test-beta" || e.Source != "browser" || e.Signal != "logs" || e.Attrs["stage.name"] != "search" || e.Scope != "" || e.Trace != "" {
		t.Fatalf("unexpected standalone projection: %+v", e)
	}
	if len(snap.Discovered) != 1 || snap.Discovered[0].Name != e.Name || snap.Discovered[0].Count != 1 || snap.Discovered[0].Assigned != 0 {
		t.Fatalf("standalone log missing from discovery: %+v", snap.Discovered)
	}
	code, raw = appRequest(t, c, server.URL, "GET", base+"/timeline?window=1h", "", nil)
	var timeline timelineResponse
	if code != http.StatusOK || json.Unmarshal(raw, &timeline) != nil || timeline.Total != 1 {
		t.Fatalf("standalone log missing from volume timeline: %d %s", code, raw)
	}
	code, raw = appRequest(t, c, server.URL, "GET", base+"/traces", "", nil)
	var traces struct{ Traces []any }
	if code != http.StatusOK || json.Unmarshal(raw, &traces) != nil || len(traces.Traces) != 0 {
		t.Fatalf("standalone log fabricated a trace: %d %s", code, raw)
	}
	code, raw = appRequest(t, c, server.URL, "GET", "/api/apps/"+other.ID+"/snapshot", "", nil)
	if code != http.StatusOK || strings.Contains(string(raw), "pwa.app.opened") {
		t.Fatalf("standalone evidence crossed app boundary: status=%d", code)
	}
	t.Run("first mapping requires an explicit goal but not journey correlation", func(t *testing.T) {
		definition := Definition{Mappings: []Mapping{{Event: "pwa.app.opened", Step: "App opened", Source: "browser", Service: "wfs-pwa-browser", Signal: "logs", Kind: "step"}}}
		code, raw := appRequest(t, c, server.URL, "POST", base+"/preview", csrf, definition)
		var preview struct {
			Conflict       string
			MissingContext int
		}
		if code != http.StatusOK || json.Unmarshal(raw, &preview) != nil || preview.Conflict == "" || preview.MissingContext != 1 {
			t.Fatalf("first mapping must explain missing goal separately from correlation: %d %s", code, raw)
		}
		definition.Goal = "pwa.app.opened"
		code, raw = appRequest(t, c, server.URL, "POST", base+"/preview", csrf, definition)
		if code != http.StatusOK || json.Unmarshal(raw, &preview) != nil || preview.Conflict != "" || preview.MissingContext != 1 {
			t.Fatalf("explicit opening goal must permit preview without inventing correlation: %d %s", code, raw)
		}
		code, raw = appRequest(t, c, server.URL, "POST", base+"/definitions", csrf, definition)
		if code != http.StatusOK {
			t.Fatalf("missing correlation must not block valid definition activation: %d %s", code, raw)
		}
		code, raw = appRequest(t, c, server.URL, "GET", base+"/snapshot", "", nil)
		if code != http.StatusOK || json.Unmarshal(raw, &snap) != nil || len(snap.Model.Journeys) != 0 || len(snap.Model.Unassigned) != 1 {
			t.Fatalf("mapping fabricated a correlated journey: %d %s", code, raw)
		}
	})

}
