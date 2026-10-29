package app

import (
	"net/http"
	"strconv"
	"time"
)

type timelineBin struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Count int    `json:"count"`
}
type timelineResponse struct {
	From      string        `json:"from"`
	To        string        `json:"to"`
	Bins      []timelineBin `json:"bins"`
	Total     int           `json:"total"`
	Truncated bool          `json:"truncated"`
	Unit      string        `json:"unit"`
}

// timeline counts retained records, never metric values, people or reconstructed
// paths. Its context is independent of the current graph's selected from/to.
func (s *Service) timeline(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "GET required", 405)
		return
	}
	q := r.URL.Query()
	window := q.Get("window")
	if !q.Has("window") {
		window = "7d"
	}
	windows := map[string]time.Duration{"1h": time.Hour, "24h": 24 * time.Hour, "7d": 7 * 24 * time.Hour, "30d": 30 * 24 * time.Hour}
	span, ok := windows[window]
	if !ok {
		http.Error(w, "Choose timeline window 1h, 24h, 7d or 30d", 400)
		return
	}
	var definition int64
	if value := q.Get("definition"); value != "" {
		parsed, err := strconv.ParseInt(value, 10, 64)
		if err != nil || parsed <= 0 {
			http.Error(w, "Choose a positive configuration definition version", 400)
			return
		}
		definition = parsed
	}
	s.Store.Mu.Lock()
	defer s.Store.Mu.Unlock()
	now := s.Store.Now().UTC()
	from := now.Add(-span)
	expiry := now.Add(-30 * 24 * time.Hour).UnixNano()
	const limit = 10000
	records, err := s.Store.EvidenceLimit(limit + 1)
	if err != nil {
		http.Error(w, "Store unavailable", 503)
		return
	}
	response := timelineResponse{From: from.Format(time.RFC3339Nano), To: now.Format(time.RFC3339Nano), Bins: make([]timelineBin, 60), Truncated: len(records) > limit, Unit: "retained records by occurrence time"}
	if response.Truncated {
		records = records[len(records)-limit:]
	}
	width := span / 60
	for i := range response.Bins {
		response.Bins[i] = timelineBin{From: from.Add(time.Duration(i) * width).Format(time.RFC3339Nano), To: from.Add(time.Duration(i+1) * width).Format(time.RFC3339Nano)}
	}
	for _, record := range records {
		if record.Time <= expiry || record.Time < from.UnixNano() || record.Time >= now.UnixNano() || q.Get("scope") != "" && record.Scope != q.Get("scope") || !matchesVersion(record.Version, q.Get("version")) || definition != 0 && record.Definition != definition {
			continue
		}
		index := int((record.Time - from.UnixNano()) / int64(width))
		if index < 0 || index >= len(response.Bins) {
			http.Error(w, "Timeline bin unavailable", 503)
			return
		}
		response.Bins[index].Count++
		response.Total++
	}
	write(w, response)
}
