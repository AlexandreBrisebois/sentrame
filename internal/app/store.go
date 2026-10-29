package app

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "modernc.org/sqlite"
	"os"
	"strings"
	"sync"
	"time"
)

type Store struct {
	Path string
	DB   *sql.DB
	Mu   sync.Mutex
	Now  func() time.Time
}

func Open(path string) (*Store, error) {
	db, e := sql.Open("sqlite", path)
	if e != nil {
		return nil, e
	}
	db.SetMaxOpenConns(1)
	var version, tables int
	if e = db.QueryRow(`PRAGMA user_version`).Scan(&version); e != nil {
		db.Close()
		return nil, e
	}
	if e = db.QueryRow(`SELECT count(*) FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'`).Scan(&tables); e != nil {
		db.Close()
		return nil, e
	}
	if version != schemaVersion && (version != 0 || tables != 0) {
		db.Close()
		return nil, fmt.Errorf("unsupported database schema %d; select a fresh synthetic data directory (existing data was not reset)", version)
	}

	_, e = db.Exec(`PRAGMA journal_mode=WAL; PRAGMA synchronous=FULL; PRAGMA foreign_keys=ON; PRAGMA busy_timeout=3000; PRAGMA max_page_count=200000; PRAGMA journal_size_limit=16777216; CREATE TABLE IF NOT EXISTS evidence(id INTEGER PRIMARY KEY, data TEXT NOT NULL, identity TEXT UNIQUE, at INTEGER NOT NULL, scope TEXT NOT NULL); CREATE INDEX IF NOT EXISTS evidence_time ON evidence(at); CREATE INDEX IF NOT EXISTS evidence_scope ON evidence(scope); CREATE TABLE IF NOT EXISTS settings(key TEXT PRIMARY KEY,value TEXT NOT NULL); INSERT OR IGNORE INTO settings VALUES('revision','0'); CREATE TABLE IF NOT EXISTS definitions(version INTEGER PRIMARY KEY,data TEXT NOT NULL); CREATE TABLE IF NOT EXISTS views(id INTEGER PRIMARY KEY,data TEXT NOT NULL); CREATE TABLE IF NOT EXISTS invitations(id INTEGER PRIMARY KEY AUTOINCREMENT,label TEXT NOT NULL,created INTEGER NOT NULL,expires INTEGER NOT NULL,token_hash TEXT NOT NULL UNIQUE,app_id TEXT NOT NULL,generation TEXT NOT NULL,redeemed INTEGER NOT NULL DEFAULT 0,revoked INTEGER NOT NULL DEFAULT 0); CREATE TABLE IF NOT EXISTS tombstones(scope TEXT PRIMARY KEY,until INTEGER NOT NULL); CREATE TABLE IF NOT EXISTS summaries(day TEXT PRIMARY KEY,data TEXT NOT NULL); CREATE TABLE IF NOT EXISTS boundaries(scope TEXT,tab TEXT,start INTEGER, end INTEGER,until INTEGER,PRIMARY KEY(scope,tab,start));`)
	if e != nil {
		db.Close()
		return nil, e
	}
	if _, e = db.Exec(fmt.Sprintf("PRAGMA user_version=%d", schemaVersion)); e != nil {
		db.Close()
		return nil, e
	}
	s := &Store{DB: db, Now: time.Now, Path: path}
	if _, e = db.Exec(`INSERT OR IGNORE INTO settings VALUES('app_id',?),('generation',?)`, random(), random()); e != nil {
		db.Close()
		return nil, e
	}
	ds, err := s.Definitions()
	if err != nil {
		db.Close()
		return nil, err
	}
	if len(ds) == 0 {
		d := InitialDefinition()
		b, _ := json.Marshal(d)
		_, e = db.Exec(`INSERT INTO definitions VALUES(1,?)`, string(b))
	}
	return s, e
}
func (s *Store) Insert(es []Evidence) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	if !strings.Contains(s.Path, ":memory:") {
		var used int64
		for _, suffix := range []string{"", "-wal", "-shm"} {
			info, err := os.Stat(s.Path + suffix)
			if err == nil {
				used += info.Size()
			} else if !os.IsNotExist(err) {
				return fmt.Errorf("storage accounting unavailable")
			}
		}
		if used > 896*1024*1024 {
			return fmt.Errorf("storage byte limit")
		}
	}
	// Capture a late earlier prefix under a known retained journey's policy.
	// This is detail metadata, not a durable membership or identity index.
	captured := Model{}
	needsCapture := false
	for _, x := range es {
		if x.Scope != "" {
			needsCapture = true
			break
		}
	}
	if needsCapture {
		existing, err := s.Evidence()
		if err != nil {
			return fmt.Errorf("writer unavailable")
		}
		defs, err := s.Definitions()
		if err != nil {
			return fmt.Errorf("writer unavailable")
		}
		captured = Derive(existing, defs, s.Now().UnixNano())
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	tx, e := s.DB.BeginTx(ctx, nil)
	if e != nil {
		return fmt.Errorf("writer unavailable")
	}
	defer tx.Rollback()
	var n int
	if e = tx.QueryRow(`SELECT count(*) FROM evidence`).Scan(&n); e != nil {
		return fmt.Errorf("writer unavailable")
	}

	var disabled string
	if e = tx.QueryRow(`SELECT value FROM settings WHERE key='disabled'`).Scan(&disabled); e != nil && e != sql.ErrNoRows {
		return fmt.Errorf("writer unavailable")
	}
	if disabled == "true" {
		return fmt.Errorf("application disabled")
	}
	var ver int64
	if e = tx.QueryRow(`SELECT max(version) FROM definitions`).Scan(&ver); e != nil {
		return fmt.Errorf("writer unavailable")
	}
	for _, x := range es {
		var until int64
		if e = tx.QueryRow(`SELECT until FROM tombstones WHERE scope=?`, x.Scope).Scan(&until); e != nil && e != sql.ErrNoRows {
			return fmt.Errorf("writer unavailable")
		}
		if until > s.Now().UnixNano() {
			continue
		}
		var barrier int
		if e = tx.QueryRow(`SELECT count(*) FROM boundaries WHERE scope=? AND tab=? AND start<=? AND end>=? AND until>?`, x.Scope, x.Tab, x.Time, x.Time, s.Now().UnixNano()).Scan(&barrier); e != nil {
			return fmt.Errorf("writer unavailable")
		}
		if barrier > 0 {
			if x.Attrs == nil {
				x.Attrs = map[string]string{}
			}
			x.Attrs["_partial_start"] = "true"
			x.Attrs["_summary_finalized"] = "true"
		}
		if x.Definition == 0 {
			x.Definition = ver
		}
		if x.CaptureDefinition == 0 && x.Scope != "" && (x.Source != "browser" || x.Tab != "") {
			containing := false
			var nearest *Journey
			for i := range captured.Journeys {
				j := &captured.Journeys[i]
				if j.Scope != x.Scope || j.Tab != x.Tab {
					continue
				}
				end := j.End
				if end == 0 {
					end = j.LastMapped + inactivityNanos
				}
				if x.Time >= j.Start && x.Time < end {
					containing = true
				}
				if x.Time < j.Start && j.Start-x.Time < inactivityNanos && (nearest == nil || j.Start < nearest.Start) {
					nearest = j
				}
			}
			if !containing && nearest != nil {
				x.CaptureDefinition = nearest.Definition
			}
		}
		key := any(nil)
		if x.Record != "" {
			key = x.Source + "/" + x.Signal + "/" + x.Record
		}
		raw, _ := json.Marshal(x)
		if key != nil {
			var old string
			err := tx.QueryRow(`SELECT data FROM evidence WHERE identity=?`, key).Scan(&old)
			if err != nil && err != sql.ErrNoRows {
				return fmt.Errorf("writer unavailable")
			}
			if err == nil {
				var prev Evidence
				json.Unmarshal([]byte(old), &prev)
				x.ID = prev.ID
				x.Receipt = prev.Receipt
				x.Definition = prev.Definition
				x.CaptureDefinition = prev.CaptureDefinition
				if x.Attrs == nil {
					x.Attrs = map[string]string{}
				}
				for k, v := range prev.Attrs {
					if k == "_partial_start" || k == "_summary_finalized" {
						x.Attrs[k] = v
					}
				}
				if x.Fallback {
					x.Time = prev.Time
				}
				raw, _ = json.Marshal(x)
				canonicalOld, _ := json.Marshal(prev)
				if string(canonicalOld) != string(raw) {
					return fmt.Errorf("identity conflict")
				}
				continue
			}
		}
		if n >= 100000 {
			return fmt.Errorf("storage record limit")
		}
		n++
		if _, e = tx.Exec(`INSERT INTO evidence(data,identity,at,scope) VALUES(?,?,?,?)`, string(raw), key, x.Time, x.Scope); e != nil {
			return fmt.Errorf("writer unavailable")
		}
	}
	_, e = tx.Exec(`UPDATE settings SET value=CAST(value AS INTEGER)+1 WHERE key='revision'`)
	if e != nil {
		return e
	}
	return tx.Commit()
}
func (s *Store) Evidence() ([]Evidence, error) {
	rows, e := s.DB.Query(`SELECT id,data FROM evidence ORDER BY at,id`)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []Evidence{}
	for rows.Next() {
		var x Evidence
		var raw string
		if e = rows.Scan(&x.ID, &raw); e != nil {
			return nil, e
		}
		id := x.ID
		if e = json.Unmarshal([]byte(raw), &x); e != nil {
			return nil, e
		}
		x.ID = id
		out = append(out, x)
	}
	return out, rows.Err()
}
func (s *Store) Definitions() ([]Definition, error) {
	rows, e := s.DB.Query(`SELECT data FROM definitions ORDER BY version`)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []Definition{}
	for rows.Next() {
		var d Definition
		var b string
		if err := rows.Scan(&b); err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(b), &d); err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, rows.Err()
}
func InitialDefinition() Definition {
	return Definition{Version: 1, Goal: "completion.entered", Mappings: []Mapping{{Event: "discovery.entered", Step: "Discovery", Kind: "step"}, {Event: "vote.requested", Step: "Intent", Kind: "step", Source: "browser"}, {Event: "vote.committed", Step: "Vote", Kind: "step", Source: "server"}, {Event: "publication.failed", Step: "Publication error", Kind: "step"}, {Event: "publication.recovered", Step: "Recovered", Kind: "step"}, {Event: "completion.entered", Step: "Completion", Kind: "step", Source: "browser"}, {Event: "journey.exited", Step: "Exit", Kind: "exit"}, {Event: "stage.entered", Step: "Stage", Kind: "stage-start"}, {Event: "stage.exited", Step: "Stage", Kind: "stage-end"}, {Event: "activity.changed", Step: "Activity", Kind: "activity"}}}
}

func (s *Store) EvidenceLimit(limit int) ([]Evidence, error) {
	rows, e := s.DB.Query(`SELECT id,data FROM evidence ORDER BY at DESC,id DESC LIMIT ?`, limit)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []Evidence{}
	for rows.Next() {
		var x Evidence
		var raw string
		if e = rows.Scan(&x.ID, &raw); e != nil {
			return nil, e
		}
		id := x.ID
		if e = json.Unmarshal([]byte(raw), &x); e != nil {
			return nil, e
		}
		x.ID = id
		out = append(out, x)
	}
	for a, b := 0, len(out)-1; a < b; a, b = a+1, b-1 {
		out[a], out[b] = out[b], out[a]
	}
	return out, rows.Err()
}
