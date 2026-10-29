package app

import (
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

type errorRecordsResponse struct {
	Records   []Evidence `json:"records"`
	Total     int        `json:"total"`
	Truncated bool       `json:"truncated"`
	Unit      string     `json:"unit"`
	Policy    string     `json:"policy"`
}

// errors reads stored sanitized evidence independently of journey membership.
// It does not infer incidents, causes or visits from time proximity.
func (s *Service) errors(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "GET required", 405)
		return
	}
	q := r.URL.Query()
	search := q.Get("q")
	if !utf8.ValidString(search) || utf8.RuneCountInString(search) > 160 || strings.ContainsFunc(search, unicode.IsControl) {
		http.Error(w, "Use a diagnostic search of at most 160 characters without controls", 400)
		return
	}
	search = strings.ToLower(strings.TrimSpace(search))
	var definition int64
	if value := q.Get("definition"); value != "" {
		number, err := strconv.ParseInt(value, 10, 64)
		if err != nil || number <= 0 {
			http.Error(w, "Choose a positive configuration definition version", 400)
			return
		}
		definition = number
	}
	s.Store.Mu.Lock()
	defer s.Store.Mu.Unlock()
	now := s.Store.Now()
	cut := now.Add(-30 * 24 * time.Hour).UnixNano()
	from, err := parseTime(q.Get("from"), cut)
	if err != nil {
		http.Error(w, "Use an ISO start date", 400)
		return
	}
	to, err := parseTime(q.Get("to"), now.UnixNano()+1)
	if err != nil || from >= to {
		http.Error(w, "Choose an increasing diagnostic date range", 400)
		return
	}
	const scanLimit = 10000
	const recordLimit = 100
	records, err := s.Store.EvidenceLimit(scanLimit + 1)
	if err != nil {
		http.Error(w, "Store unavailable", 503)
		return
	}
	result := errorRecordsResponse{Records: []Evidence{}, Truncated: len(records) > scanLimit, Unit: "diagnostic records", Policy: diagnosticPolicyDescription}
	if len(records) > scanLimit {
		records = records[len(records)-scanLimit:]
	}
	for i := len(records) - 1; i >= 0; i-- {
		record := records[i]
		if record.Time <= cut || record.Time < from || record.Time >= to || record.Time > now.UnixNano() || !isErrorDiagnostic(record) || q.Get("scope") != "" && record.Scope != q.Get("scope") || !matchesVersion(record.Version, q.Get("version")) || definition != 0 && record.Definition != definition {
			continue
		}
		if search != "" {
			searchable := strings.Join([]string{record.Name, record.Source, record.Signal, record.Service, record.Version, record.Scope, record.Trace, record.Span, record.Attrs["exception.type"], record.Attrs["exception.message"], record.Attrs["exception.stacktrace"], record.Attrs["log.message"]}, "\n")
			if !strings.Contains(strings.ToLower(searchable), search) {
				continue
			}
		}
		result.Total++
		if len(result.Records) < recordLimit {
			result.Records = append(result.Records, record)
		} else {
			result.Truncated = true
		}
	}
	write(w, result)
}
