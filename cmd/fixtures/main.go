// Command fixtures generates real standard OTLP protobuf requests and executes
// the version-controlled REST Client demonstration against the actual service.
package main

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	logsv1 "go.opentelemetry.io/proto/otlp/collector/logs/v1"
	commonv1 "go.opentelemetry.io/proto/otlp/common/v1"
	logdata "go.opentelemetry.io/proto/otlp/logs/v1"
	resourcev1 "go.opentelemetry.io/proto/otlp/resource/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type Event struct {
	Name       string            `json:"name"`
	Record     string            `json:"record,omitempty"`
	Offset     string            `json:"offset"`
	Action     string            `json:"action,omitempty"`
	Occurrence string            `json:"occurrence,omitempty"`
	Severity   int32             `json:"severity,omitempty"`
	Attributes map[string]string `json:"attributes,omitempty"`
}
type Scenario struct {
	Name       string  `json:"name"`
	Scope      string  `json:"scope,omitempty"`
	Producer   string  `json:"producer"`
	NoTab      bool    `json:"noTab,omitempty"`
	NoIdentity bool    `json:"noIdentity,omitempty"`
	Events     []Event `json:"events"`
}
type Generated struct {
	Bodies        map[string][]byte `json:"-"`
	JSON          map[string][]byte `json:"-"`
	Scopes        map[string]string `json:"scopes"`
	Counts        map[string]int    `json:"counts"`
	Hashes        map[string]string `json:"hashes"`
	LiveStart     string            `json:"liveStart"`
	GeneratedAt   string            `json:"generatedAt"`
	FollowupState string            `json:"followupState"`
}

func randomID() string {
	b := make([]byte, 16)
	if _, e := rand.Read(b); e != nil {
		panic(e)
	}
	return hex.EncodeToString(b)
}
func kv(k, v string) *commonv1.KeyValue {
	return &commonv1.KeyValue{Key: k, Value: &commonv1.AnyValue{Value: &commonv1.AnyValue_StringValue{StringValue: v}}}
}
func build(specs []Scenario, now time.Time) (Generated, error) {
	g := Generated{FollowupState: "unsent", Bodies: map[string][]byte{}, JSON: map[string][]byte{}, Scopes: map[string]string{}, Counts: map[string]int{}, Hashes: map[string]string{}, GeneratedAt: now.UTC().Format(time.RFC3339Nano), LiveStart: now.Add(-2 * time.Second).UTC().Format(time.RFC3339Nano)}
	tabs := map[string]string{}
	actions := map[string]string{}
	occurrences := map[string]string{}
	records := map[string]string{}
	for _, s := range specs {
		if s.Name == "" || s.Producer != "browser" && s.Producer != "server" {
			return g, errors.New("fixture needs a name and configured producer")
		}
		if s.Scope != "" && g.Scopes[s.Scope] == "" {
			g.Scopes[s.Scope] = randomID()
			tabs[s.Scope] = randomID()
		}
	}
	for _, s := range specs {
		req := &logsv1.ExportLogsServiceRequest{ResourceLogs: []*logdata.ResourceLogs{{Resource: &resourcev1.Resource{Attributes: []*commonv1.KeyValue{kv("service.name", "sentrame-synthetic-"+s.Producer), kv("service.version", "synthetic-a1"), kv("telemetry.sdk.language", "go")}}, ScopeLogs: []*logdata.ScopeLogs{{Scope: &commonv1.InstrumentationScope{Name: "sentrame-standard-fixtures", Version: "1"}}}}}}
		for _, e := range s.Events {
			off, err := time.ParseDuration(e.Offset)
			if err != nil {
				return g, fmt.Errorf("fixture %s offset: %w", s.Name, err)
			}
			attrs := []*commonv1.KeyValue{}
			if s.Scope != "" {
				attrs = append(attrs, kv("session.id", g.Scopes[s.Scope]))
				if !s.NoTab {
					attrs = append(attrs, kv("tab.id", tabs[s.Scope]))
				}
			}
			if !s.NoIdentity {
				id := randomID()
				if e.Record != "" {
					key := s.Scope + "/" + e.Record
					if records[key] == "" {
						records[key] = randomID()
					}
					id = records[key]
				}
				attrs = append(attrs, kv("event.id", id))
			}
			if e.Action != "" {
				key := s.Scope + "/" + e.Action
				if actions[key] == "" {
					actions[key] = randomID()
				}
				attrs = append(attrs, kv("action.id", actions[key]))
			}
			if e.Occurrence != "" {
				key := s.Scope + "/" + s.Producer + "/" + e.Occurrence
				if occurrences[key] == "" {
					occurrences[key] = randomID()
				}
				attrs = append(attrs, kv("stage.occurrence.id", occurrences[key]))
			}
			keys := []string{}
			for k := range e.Attributes {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				v := e.Attributes[k]
				if k == "exception.occurrence.id" {
					key := s.Scope + "/exception/" + v
					if occurrences[key] == "" {
						occurrences[key] = randomID()
					}
					v = occurrences[key]
				}
				attrs = append(attrs, kv(k, v))
			}
			// No body or proprietary schema is needed by standard OTLP ingestion.
			record := &logdata.LogRecord{TimeUnixNano: uint64(now.Add(off).UnixNano()), ObservedTimeUnixNano: uint64(now.UnixNano()), EventName: e.Name, SeverityNumber: logdata.SeverityNumber(e.Severity), Attributes: attrs}
			req.ResourceLogs[0].ScopeLogs[0].LogRecords = append(req.ResourceLogs[0].ScopeLogs[0].LogRecords, record)
		}
		b, e := proto.Marshal(req)
		if e != nil {
			return g, e
		}
		j, e := protojson.MarshalOptions{Indent: "  "}.Marshal(req)
		if e != nil {
			return g, e
		}
		g.Bodies[s.Name] = b
		g.JSON[s.Name] = j
		h := sha256.Sum256(b)
		g.Hashes[s.Name] = hex.EncodeToString(h[:])
		g.Counts[s.Name] = len(s.Events)
	}
	return g, nil
}
func envRequired(k string) (string, error) {
	v := os.Getenv(k)
	if v == "" {
		return "", fmt.Errorf("set %s before generating/executing local synthetic requests", k)
	}
	return v, nil
}
func baseURL() string {
	v := os.Getenv("SENTRAME_URL")
	if v == "" {
		v = "http://127.0.0.1:8096"
	}
	return strings.TrimRight(v, "/")
}
func generate() error {
	b, e := os.ReadFile("rest/fixtures/scenarios.json")
	if e != nil {
		return e
	}
	var specs []Scenario
	if e = json.Unmarshal(b, &specs); e != nil {
		return e
	}
	g, e := build(specs, time.Now())
	if e != nil {
		return e
	}
	if e = os.MkdirAll("rest/generated", 0700); e != nil {
		return e
	}
	for name, b := range g.Bodies {
		if e = os.WriteFile(filepath.Join("rest/generated", name+".bin"), b, 0600); e != nil {
			return e
		}
		if e = os.WriteFile(filepath.Join("rest/generated", name+".json"), g.JSON[name], 0600); e != nil {
			return e
		}
	}
	manifest, _ := json.MarshalIndent(g, "", "  ")
	if e = os.WriteFile("rest/generated/manifest.json", manifest, 0600); e != nil {
		return e
	}
	// Generate local-only REST credentials and login body with correct JSON escaping.
	pass, e := envRequired("SENTRAME_PASSPHRASE")
	if e != nil {
		return e
	}
	browser, e := envRequired("SENTRAME_BROWSER_TOKEN")
	if e != nil {
		return e
	}
	server, e := envRequired("SENTRAME_SERVER_TOKEN")
	if e != nil {
		return e
	}
	login, _ := json.Marshal(map[string]string{"passphrase": pass})
	if e = os.WriteFile("rest/generated/login.json", login, 0600); e != nil {
		return e
	}
	values := map[string]string{"appId": os.Getenv("SENTRAME_APP_ID"), "baseUrl": baseURL(), "browserToken": browser, "serverToken": server, "healthyScope": g.Scopes["healthy"], "timingScope": g.Scopes["timing-samples"], "liveScope": g.Scopes["live"], "liveStart": url.QueryEscape(g.LiveStart)}
	var env strings.Builder
	env.WriteString("# Generated local synthetic credentials; do not commit.\n")
	keys := []string{}
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Fprintf(&env, "%s=%s\n", k, values[k])
	}
	if e = os.WriteFile("rest/.env", []byte(env.String()), 0600); e != nil {
		return e
	}
	fmt.Printf("Generated %d real OTLP protobuf bodies at %s; replay uses these exact files.\n", len(g.Bodies), g.GeneratedAt)
	return nil
}

type restRequest struct {
	BodyFile          string
	Method, URL       string
	Headers           http.Header
	Body              []byte
	Expected          int
	Stream, NoCookies bool
	Name              string
}

var variableRE = regexp.MustCompile(`\{\{([a-zA-Z0-9_]+)\}\}`)

func substitute(s string, v map[string]string) (string, error) {
	missing := ""
	out := variableRE.ReplaceAllStringFunc(s, func(x string) string {
		k := x[2 : len(x)-2]
		r, ok := v[k]
		if !ok {
			missing = k
		}
		return r
	})
	if missing != "" {
		return "", fmt.Errorf("missing REST variable %s", missing)
	}
	return out, nil
}
func parseREST(path string, v map[string]string) ([]restRequest, error) {
	b, e := os.ReadFile(path)
	if e != nil {
		return nil, e
	}
	if v["appId"] == "" {
		v["appId"] = os.Getenv("SENTRAME_APP_ID")
	}
	text, e := substitute(string(b), v)
	if e != nil {
		return nil, e
	}
	lines := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
	out := []restRequest{}
	r := restRequest{Expected: 200, Headers: http.Header{}}
	inBody := false
	var body []string
	flush := func() error {
		if r.Method == "" {
			return nil
		}
		raw := strings.TrimSpace(strings.Join(body, "\n"))
		if strings.HasPrefix(raw, "< ") {
			file := strings.TrimSpace(strings.TrimPrefix(raw, "< "))
			r.BodyFile = filepath.Base(file)
			var err error
			r.Body, err = os.ReadFile(filepath.Join(filepath.Dir(path), file))
			if err != nil {
				return err
			}
		} else {
			r.Body = []byte(raw)
		}
		out = append(out, r)
		r = restRequest{Expected: 200, Headers: http.Header{}}
		inBody = false
		body = nil
		return nil
	}
	for _, line := range lines {
		if strings.HasPrefix(line, "###") {
			if e = flush(); e != nil {
				return nil, e
			}
			r.Name = strings.TrimSpace(strings.TrimPrefix(line, "###"))
			continue
		}
		if strings.HasPrefix(line, "#") {
			if strings.HasPrefix(line, "# @expect ") {
				fmt.Sscanf(line, "# @expect %d", &r.Expected)
			}
			if line == "# @stream true" {
				r.Stream = true
			}
			if line == "# @no-cookie-jar" {
				r.NoCookies = true
			}
			continue
		}
		if r.Method == "" {
			if strings.TrimSpace(line) == "" {
				continue
			}
			parts := strings.SplitN(line, " ", 2)
			if len(parts) != 2 {
				return nil, fmt.Errorf("invalid REST request line")
			}
			r.Method = parts[0]
			r.URL = parts[1]
			continue
		}
		if !inBody {
			if line == "" {
				inBody = true
				continue
			}
			p := strings.SplitN(line, ":", 2)
			if len(p) != 2 {
				return nil, errors.New("invalid REST header")
			}
			r.Headers.Set(strings.TrimSpace(p[0]), strings.TrimSpace(p[1]))
			continue
		}
		body = append(body, line)
	}
	if e = flush(); e != nil {
		return nil, e
	}
	return out, nil
}
func variables(g Generated) (map[string]string, error) {
	browser, e := envRequired("SENTRAME_BROWSER_TOKEN")
	if e != nil {
		return nil, e
	}
	server, e := envRequired("SENTRAME_SERVER_TOKEN")
	if e != nil {
		return nil, e
	}
	return map[string]string{"appId": os.Getenv("SENTRAME_APP_ID"), "baseUrl": baseURL(), "browserToken": browser, "serverToken": server, "healthyScope": g.Scopes["healthy"], "timingScope": g.Scopes["timing-samples"], "liveScope": g.Scopes["live"], "liveStart": url.QueryEscape(g.LiveStart)}, nil
}
func manifest() (Generated, error) {
	var g Generated
	b, e := os.ReadFile("rest/generated/manifest.json")
	if e != nil {
		return g, e
	}
	e = json.Unmarshal(b, &g)
	return g, e
}
func verify() error {
	g, e := manifest()
	if e != nil {
		return e
	}
	vars, e := variables(g)
	if e != nil {
		return e
	}
	requests, e := parseREST("rest/demo.http", vars)
	if e != nil {
		return e
	}
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar, Timeout: 20 * time.Second}
	for i, r := range requests {
		req, e := http.NewRequest(r.Method, r.URL, bytes.NewReader(r.Body))
		if e != nil {
			return e
		}
		req.Header = r.Headers.Clone()
		c := client
		if r.NoCookies {
			c = &http.Client{Timeout: 20 * time.Second}
		}
		resp, e := c.Do(req)
		if e != nil {
			return fmt.Errorf("request %d failed: %w", i+1, e)
		}
		if resp.StatusCode != r.Expected {
			resp.Body.Close()
			return fmt.Errorf("request %d (%s %s) status %d, expected %d", i+1, r.Method, req.URL.Path, resp.StatusCode, r.Expected)
		}
		if r.Stream {
			if !strings.HasPrefix(resp.Header.Get("Content-Type"), "text/event-stream") {
				resp.Body.Close()
				return errors.New("SSE content type missing")
			}
			reader := bufio.NewReader(io.LimitReader(resp.Body, 8192))
			seen := false
			for {
				line, err := reader.ReadString('\n')
				if strings.HasPrefix(line, "data:") {
					seen = true
					break
				}
				if err != nil {
					break
				}
			}
			resp.Body.Close()
			if !seen {
				return errors.New("SSE did not deliver an initial reset/invalidation")
			}
		} else {
			raw, e := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
			resp.Body.Close()
			if e != nil {
				return e
			}
			if strings.Contains(string(raw), "SHOULD-NOT-PERSIST") {
				return errors.New("privacy sentinel leaked through public response")
			}
		}
		if r.Method == "POST" && r.BodyFile == "live-followup.bin" && resp.StatusCode == 200 {
			g.FollowupState = "sent"
			if e = saveManifest(g); e != nil {
				return e
			}
		}
		fmt.Printf("PASS %02d %s %s → %d\n", i+1, r.Method, req.URL.Path, resp.StatusCode)
	}
	return checkSnapshots(client, g)
}
func checkSnapshots(client *http.Client, g Generated) error {
	type evidence struct {
		Name   string `json:"name"`
		Source string `json:"source"`
		Scope  string `json:"scope"`
	}
	type journey struct {
		Status string   `json:"status"`
		Path   []string `json:"path"`
	}
	type timingSummary struct {
		Basis        string   `json:"basis"`
		Samples      int      `json:"samples"`
		MedianMs     *float64 `json:"medianMs"`
		P90Ms        *float64 `json:"p90Ms"`
		Completeness float64  `json:"completeness"`
	}
	type node struct {
		Label  string          `json:"label"`
		Timing []timingSummary `json:"timing"`
	}
	type snapshot struct {
		Evidence []evidence `json:"evidence"`
		Model    struct {
			Journeys []journey `json:"journeys"`
			Nodes    []node    `json:"nodes"`
		} `json:"model"`
	}
	expected := map[string]int{"healthy": 16, "recovery": 6, "initial": 1, "server-only": 3, "forged": 1, "inactive": 1, "baseline": 2, "live": 2, "atomic": 0, "timing-samples": 42}
	for scope, n := range expected {
		resp, e := client.Get(baseURL() + appAPI("snapshot") + "?scope=" + url.QueryEscape(g.Scopes[scope]))
		if e != nil {
			return e
		}
		raw, e := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
		resp.Body.Close()
		if e != nil {
			return e
		}
		var snap snapshot
		if e = json.Unmarshal(raw, &snap); e != nil {
			return e
		}
		if len(snap.Evidence) != n {
			return fmt.Errorf("scope %s retained %d records, expected %d (including replay/atomicity)", scope, len(snap.Evidence), n)
		}
		if scope == "healthy" {
			votes := 0
			for _, x := range snap.Evidence {
				if x.Name == "vote.committed" && x.Source == "server" {
					votes++
				}
			}
			if votes != 10 {
				return fmt.Errorf("healthy votes %d, expected ten trusted commitments", votes)
			}
			if len(snap.Model.Journeys) != 1 || snap.Model.Journeys[0].Status != "completed" {
				return errors.New("healthy scope did not produce one completed journey")
			}
		}
		if scope == "forged" && len(snap.Model.Journeys) != 0 {
			return errors.New("browser self-asserted server authority became a mapped journey")
		}
		if scope == "inactive" && (len(snap.Model.Journeys) != 1 || snap.Model.Journeys[0].Status != "inactive") {
			return errors.New("silence did not produce inactive journey")
		}
		if scope == "recovery" && (len(snap.Model.Journeys) != 1 || snap.Model.Journeys[0].Status != "completed") {
			return errors.New("post-commit error prevented eventual completion")
		}
		if scope == "timing-samples" {
			if len(snap.Model.Journeys) != 1 || snap.Model.Journeys[0].Status != "completed" {
				return errors.New("timing sample scope must be one completed captured journey")
			}
			found := false
			for _, node := range snap.Model.Nodes {
				if node.Label == "Stage" {
					for _, timing := range node.Timing {
						if timing.Basis == "measured-elapsed" {
							if timing.Samples != 20 || timing.MedianMs == nil || *timing.MedianMs != 10500 || timing.P90Ms == nil || *timing.P90Ms != 18000 || timing.Completeness != 1 {
								return errors.New("twenty measured stage samples must yield median10.5s,p90=18s,completeness100percent")
							}
							found = true
						}
					}
				}
			}
			if !found {
				return errors.New("timing sample scope omitted measured Stage summary")
			}
			fmt.Println("PASS measured Stage:20 occurrences,median10.5s,p90=18s,completeness100percent")
		}
		fmt.Printf("PASS retained scope %-11s = %d records\n", scope, n)
	}
	fmt.Println("Verified literal rest/demo.http plus protected per-scope read assertions through the actual service. No DB preload.")
	return nil
}
func sendFile(name, producer string) error {
	b, e := os.ReadFile("rest/generated/" + name + ".bin")
	if e != nil {
		return e
	}
	var g Generated
	if name == "live-followup" {
		g, e = manifest()
		if e != nil {
			return e
		}
		if g.FollowupState == "unsent" {
			b, e = refreshFollowup(b, time.Now(), false)
			if e != nil {
				return e
			}
			var q logsv1.ExportLogsServiceRequest
			if e = proto.Unmarshal(b, &q); e != nil {
				return e
			}
			j, err := protojson.MarshalOptions{Indent: "  "}.Marshal(&q)
			if err != nil {
				return err
			}
			if e = os.WriteFile("rest/generated/live-followup.bin", b, 0600); e != nil {
				return e
			}
			if e = os.WriteFile("rest/generated/live-followup.json", j, 0600); e != nil {
				return e
			}
			h := sha256.Sum256(b)
			g.Hashes["live-followup"] = hex.EncodeToString(h[:])
			g.FollowupState = "prepared"
			if e = saveManifest(g); e != nil {
				return e
			}
		}
	}
	key := "SENTRAME_BROWSER_TOKEN"
	if producer == "server" {
		key = "SENTRAME_SERVER_TOKEN"
	}
	token, e := envRequired(key)
	if e != nil {
		return e
	}
	req, _ := http.NewRequest("POST", baseURL()+"/v1/logs", bytes.NewReader(b))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/x-protobuf")
	client := &http.Client{Timeout: 15 * time.Second}
	resp, e := client.Do(req)
	if e != nil {
		return e
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("ingestion status %d", resp.StatusCode)
	}
	if name == "live-followup" {
		g.FollowupState = "sent"
		if e = saveManifest(g); e != nil {
			return e
		}
	}
	fmt.Printf("Sent prepared %s body via actual OTLP ingestion; retries preserve exact bytes.\n", name)
	return nil
}
func run() error {
	u, err := url.Parse(baseURL())
	if err != nil || u.Scheme != "http" || u.Hostname() != "127.0.0.1" && u.Hostname() != "localhost" && u.Hostname() != "::1" || u.User != nil || u.RawQuery != "" || u.Fragment != "" {
		return errors.New("SENTRAME_URL must be a plain local loopback HTTP origin for this synthetic demonstration")
	}
	if len(os.Args) < 2 {
		return errors.New("usage: fixtures generate|verify|demo|live|replay|followup|workspace|workspace-generate|workspace-verify|diagnostics|diagnostics-generate|diagnostics-verify")
	}
	if e := prepareDemo(); e != nil {
		return e
	}
	switch os.Args[1] {
	case "standard-verify":
		return verifyStandard()
	case "standard":
		if err := generateStandard(); err != nil {
			return err
		}
		return verifyStandard()
	case "standard-generate":
		return generateStandard()
	case "diagnostics":
		if err := generateDiagnostics(); err != nil {
			return err
		}
		return verifyDiagnostics()
	case "diagnostics-generate":
		return generateDiagnostics()
	case "diagnostics-verify":
		return verifyDiagnostics()

	case "workspace":
		if err := generateWorkspace(); err != nil {
			return err
		}
		return verifyWorkspace()
	case "workspace-generate":
		return generateWorkspace()
	case "workspace-verify":
		return verifyWorkspace()

	case "generate":
		return generate()
	case "verify":
		return verify()
	case "demo":
		if e := generate(); e != nil {
			return e
		}
		return verify()
	case "live":
		if e := generate(); e != nil {
			return e
		}
		return sendFile("live-start", "browser")
	case "replay":
		return sendFile("healthy-server", "server")
	case "followup":
		return sendFile("live-followup", "browser")
	default:
		return errors.New("unknown fixture command")
	}
}
func main() {
	if e := run(); e != nil {
		fmt.Fprintln(os.Stderr, "fixtures:", e)
		os.Exit(1)
	}
}

func refreshFollowup(body []byte, now time.Time, alreadyPrepared bool) ([]byte, error) {
	if alreadyPrepared {
		return append([]byte(nil), body...), nil
	}
	var req logsv1.ExportLogsServiceRequest
	if e := proto.Unmarshal(body, &req); e != nil {
		return nil, e
	}
	count := 0
	for _, r := range req.ResourceLogs {
		for _, s := range r.ScopeLogs {
			for _, l := range s.LogRecords {
				l.TimeUnixNano = uint64(now.UnixNano())
				l.ObservedTimeUnixNano = uint64(now.UnixNano())
				count++
			}
		}
	}
	if count == 0 {
		return nil, errors.New("followup has no log records")
	}
	return proto.Marshal(&req)
}
func saveManifest(g Generated) error {
	b, e := json.MarshalIndent(g, "", "  ")
	if e != nil {
		return e
	}
	return os.WriteFile("rest/generated/manifest.json", b, 0600)
}
