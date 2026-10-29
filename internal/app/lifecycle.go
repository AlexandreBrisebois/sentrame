package app

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math"
	"net/url"
	"time"
)

func (s *Store) Cleanup() error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	now := s.Now()
	cut := now.Add(-30 * 24 * time.Hour).UnixNano()
	es, e := s.Evidence()
	if e != nil {
		return e
	}
	defs, e := s.Definitions()
	if e != nil {
		return e
	}
	model := Derive(es, defs, now.UnixNano())
	days := map[string]*DaySummary{}
	barriers := []Journey{}
	load := func(day string) *DaySummary {
		if d := days[day]; d != nil {
			return d
		}
		d := &DaySummary{Day: day, Final: true, Groups: map[string]*SummaryGroup{}}
		var b string
		if s.DB.QueryRow(`SELECT data FROM summaries WHERE day=?`, day).Scan(&b) == nil {
			json.Unmarshal([]byte(b), d)
		}
		days[day] = d
		return d
	}
	group := func(d *DaySummary, key string) *SummaryGroup {
		g := d.Groups[key]
		if g == nil {
			g = &SummaryGroup{Histograms: map[string]map[int]int{}, Excluded: map[string]int{}, Policy: "UTC captured-start cohort; correction ends start+30 days; log histogram ratio 1.1 / 1ms-24h approximate; no individual membership"}
			d.Groups[key] = g
		}
		return g
	}
	assigned := map[int64]bool{}
	updates := map[int64]Evidence{}
	for _, j := range model.Journeys {
		for _, x := range j.Evidence {
			assigned[x.ID] = true
		}
		if j.Start > cut {
			continue
		}
		already := false
		for _, x := range j.Evidence {
			if x.Attrs["_summary_finalized"] == "true" {
				already = true
			}
		}
		if !already {
			barriers = append(barriers, j)
			d := load(time.Unix(0, j.Start).UTC().Format("2006-01-02"))
			g := group(d, fmt.Sprintf("definition:%d/version:%s", j.Definition, url.QueryEscape(j.Version)))
			g.Records += len(j.Evidence)
			if !j.Partial {
				g.Journeys++
				switch j.Status {
				case "completed":
					g.Completed++
				case "inactive":
					g.Inactive++
				case "exited":
					g.Exits++
				}
			}
			for _, t := range j.Timings {
				if !t.Complete || t.Milliseconds < 0 {
					if g.Excluded == nil {
						g.Excluded = map[string]int{}
					}
					g.Excluded[t.Stage+"/"+t.Basis]++
					continue
				}
				key := t.Stage + "/" + t.Basis
				h := g.Histograms[key]
				if h == nil {
					h = map[int]int{}
					g.Histograms[key] = h
				}
				bucket := -1
				if t.Milliseconds >= 1 {
					bucket = int(math.Floor(math.Log(t.Milliseconds) / math.Log(1.1)))
					if t.Milliseconds > 86400000 {
						bucket = 999
					}
				}
				h[bucket]++
			}
		}
		for _, x := range j.Evidence {
			if x.Time > cut {
				if x.Attrs == nil {
					x.Attrs = map[string]string{}
				}
				x.Attrs["_partial_start"] = "true"
				x.Attrs["_summary_finalized"] = "true"
				updates[x.ID] = x
			}
		}
	}
	for _, x := range es {
		if !assigned[x.ID] && x.Time <= cut {
			d := load(time.Unix(0, x.Time).UTC().Format("2006-01-02"))
			g := group(d, fmt.Sprintf("definition:%d/version:%s/unassigned", x.Definition, url.QueryEscape(x.Version)))
			g.Records++
		}
	}
	tx, e := s.DB.Begin()
	if e != nil {
		return e
	}
	defer tx.Rollback()
	for _, j := range barriers {
		end := j.End
		if end == 0 {
			end = j.LastMapped + int64(30*time.Minute)
		}
		until := end + int64(30*24*time.Hour)
		if _, e = tx.Exec(`INSERT OR REPLACE INTO boundaries VALUES(?,?,?,?,?)`, j.Scope, j.Tab, j.Start, end, until); e != nil {
			return e
		}
	}
	for day, d := range days {
		b, _ := json.Marshal(d)
		if _, e = tx.Exec(`INSERT OR REPLACE INTO summaries VALUES(?,?)`, day, string(b)); e != nil {
			return e
		}
	}
	for id, x := range updates {
		b, _ := json.Marshal(x)
		if _, e = tx.Exec(`UPDATE evidence SET data=? WHERE id=?`, string(b), id); e != nil {
			return e
		}
	}
	res, e := tx.Exec(`DELETE FROM evidence WHERE at<=?`, cut)
	if e != nil {
		return e
	}
	n, _ := res.RowsAffected()
	if _, e = tx.Exec(`DELETE FROM summaries WHERE day<?`, now.AddDate(-1, 0, 0).UTC().Format("2006-01-02")); e != nil {
		return e
	}
	if _, e = tx.Exec(`DELETE FROM tombstones WHERE until<=?`, now.UnixNano()); e != nil {
		return e
	}
	if _, e = tx.Exec(`DELETE FROM boundaries WHERE until<=?`, now.UnixNano()); e != nil {
		return e
	}
	state, _ := json.Marshal(model.Stats)
	signature := fmt.Sprintf("%x", sha256.Sum256(state))
	var oldSignature string
	if err := tx.QueryRow(`SELECT value FROM settings WHERE key='lifecycle_signature'`).Scan(&oldSignature); err != nil && err.Error() != "sql: no rows in result set" {
		return err
	}
	changed := oldSignature != signature
	if changed {
		if _, e = tx.Exec(`INSERT OR REPLACE INTO settings VALUES('lifecycle_signature',?)`, signature); e != nil {
			return e
		}
	}
	if n > 0 || changed {
		if _, e = tx.Exec(`UPDATE settings SET value=CAST(value AS INTEGER)+1 WHERE key='revision'`); e != nil {
			return e
		}
	}
	return tx.Commit()
}
