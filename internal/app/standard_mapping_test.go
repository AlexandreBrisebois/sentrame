package app

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"testing"
)

func postMeaning(t *testing.T, c *http.Client, base, csrf, path string, d any, want int) map[string]any {
	t.Helper()
	b, _ := json.Marshal(d)
	r, _ := http.NewRequest("POST", base+path, bytes.NewReader(b))
	r.Header.Set("Origin", base)
	r.Header.Set("X-CSRF-Token", csrf)
	r.Header.Set("Content-Type", "application/json")
	resp, err := c.Do(r)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != want {
		t.Fatalf("%s status%d want%d: %s", path, resp.StatusCode, want, raw)
	}
	out := map[string]any{}
	json.Unmarshal(raw, &out)
	return out
}
func ordinaryDefinition() map[string]any {
	return map[string]any{"goal": "vote.committed", "mappings": []map[string]any{
		{"event": "POST /votes", "step": "Vote processing", "kind": "step", "source": "server", "service": "synthetic-voting", "carrier": "span", "key": "http.request.method", "value": "POST"},
		{"event": "vote.committed", "step": "Vote saved", "kind": "step", "source": "server", "service": "synthetic-voting", "carrier": "span-event"},
		{"event": "vote.persist", "step": "Persistence operation", "kind": "step", "source": "server", "service": "synthetic-voting", "carrier": "span"},
	}}
}
func TestStandardMappingContextualPreviewAndTraceCapture(t *testing.T) {
	_, server, c := fixtureService(t)
	tb, lb, _ := standardBodies(t)
	for signal, body := range map[string][]byte{"traces": tb, "logs": lb} {
		code, _ := protocolRequest(t, c, server.URL, signal, "application/x-protobuf", "", "server-test", body)
		if code != 200 {
			t.Fatal(code)
		}
	}
	csrf := loginTest(t, c, server.URL)
	d := ordinaryDefinition()
	preview := postMeaning(t, c, server.URL, csrf, "/api/preview", d, 200)
	if preview["conflict"] != "" {
		t.Fatal(preview["conflict"])
	}
	counts := preview["matches"].(map[string]any)
	if counts["vote.committed"] != float64(1) || counts["vote.persist"] != float64(2) {
		t.Fatalf("carrier selector ignored: %v", counts)
	}
	changed := ordinaryDefinition()
	changed["mappings"].([]map[string]any)[0]["service"] = "unrelated-service"
	p := postMeaning(t, c, server.URL, csrf, "/api/preview", changed, 200)
	if p["matches"].(map[string]any)["POST /votes"] != nil {
		t.Fatal("service selector ignored")
	}
	activation := postMeaning(t, c, server.URL, csrf, "/api/definitions", d, 200)
	version := int(activation["version"].(float64))
	path := "/api/traces/11111111111111111111111111111111"
	old := getTraceJSON(t, c, server.URL, path, 200)
	for _, v := range old["business"].([]any) {
		b := v.(map[string]any)
		if b["definition"] != float64(1) || b["stage"] == "Vote saved" {
			t.Fatal("activation silently remapped historical capture")
		}
	}
	prospective := getTraceJSON(t, c, server.URL, path+"?previewDefinition="+strconv.Itoa(version), 200)
	business := prospective["business"].([]any)
	if len(business) != 4 {
		t.Fatalf("want native span labels plus trusted event: %v", business)
	}
	samples := 0
	for _, v := range business {
		m := v.(map[string]any)
		if m["carrier"] == "span" {
			if m["basis"] != "native-inclusive" || m["occurrence"] == "" || m["durationNs"] == nil {
				t.Fatal("span requires custom timing")
			}
			samples++
		}
		if m["goal"] == true && m["stage"] != "Vote saved" {
			t.Fatal("untrusted/generic evidence established persistence")
		}
	}
	if samples != 3 {
		t.Fatal("repeated operations collapsed")
	}
	if prospective["interpretation"] != "explicit non-persisting preview" {
		t.Fatal("preview context missing")
	}
}
func TestStandardMappingDisjointAndConflictingSelectors(t *testing.T) {
	_, server, c := fixtureService(t)
	csrf := loginTest(t, c, server.URL)
	d := ordinaryDefinition()
	rules := d["mappings"].([]map[string]any)
	r := map[string]any{"event": "vote.persist", "step": "Other service", "kind": "step", "source": "server", "service": "other-service", "carrier": "span"}
	d["mappings"] = append(rules, r)
	postMeaning(t, c, server.URL, csrf, "/api/definitions", d, 200)
	r["service"] = ""
	postMeaning(t, c, server.URL, csrf, "/api/definitions", d, 400)
}

func TestStandardMappingNativeGoalEvaluatesAtEndAndRejectsAmbiguity(t *testing.T) {
	_, server, c := fixtureService(t)
	tb, _, _ := standardBodies(t)
	code, _ := protocolRequest(t, c, server.URL, "traces", "application/x-protobuf", "", "server-test", tb)
	if code != 200 {
		t.Fatal(code)
	}
	csrf := loginTest(t, c, server.URL)
	d := ordinaryDefinition()
	d["goal"] = "vote.persist"
	activated := postMeaning(t, c, server.URL, csrf, "/api/definitions", d, 200)
	version := int(activated["version"].(float64))
	detail := getTraceJSON(t, c, server.URL, "/api/traces/11111111111111111111111111111111?previewDefinition="+strconv.Itoa(version), 200)
	for _, v := range detail["business"].([]any) {
		m := v.(map[string]any)
		if m["stage"] == "Persistence operation" {
			if m["goal"] != true || m["boundary"] != "span-end" || m["boundaryNs"] == nil {
				t.Fatal("native goal not evaluated at completion boundary")
			}
		}
	}
	rules := d["mappings"].([]map[string]any)
	d["mappings"] = append(rules, map[string]any{"event": "vote.persist", "step": "Browser intent", "kind": "step", "source": "browser", "carrier": "log"})
	postMeaning(t, c, server.URL, csrf, "/api/definitions", d, 400)
	d = ordinaryDefinition()
	d["mappings"].([]map[string]any)[0]["key"] = "url.full"
	d["mappings"].([]map[string]any)[0]["value"] = "unsafe"
	postMeaning(t, c, server.URL, csrf, "/api/definitions", d, 400)
}

func TestStandardMappingSafeHTTPRouteTemplateCondition(t *testing.T) {
	_, server, c := fixtureService(t)
	csrf := loginTest(t, c, server.URL)
	d := ordinaryDefinition()
	m := d["mappings"].([]map[string]any)[0]
	m["key"] = "http.route"
	m["value"] = "/votes/{voteId}"
	postMeaning(t, c, server.URL, csrf, "/api/definitions", d, 200)
}
