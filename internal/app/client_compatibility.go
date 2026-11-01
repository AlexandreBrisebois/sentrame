package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

const clientHistoryMessage = "Legacy history is preserved. Mixed journeys and incompatible definitions or summaries are unavailable in client analysis."

// identitySource preserves the existing client replay domain. This is an internal
// compatibility alias, never a producer role or a source of authorization.
func identitySource(source string) string {
	if source == "client" {
		return "browser"
	}
	return source
}

// stageClientHistory runs only on the offline converter's private staged copy.
// Records, definitions and references are preserved, not rewritten to imply a
// client path. Exclusions persist independently of the excluded record's expiry.
func stageClientHistory(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err = tx.Exec(`CREATE TABLE IF NOT EXISTS legacy_excluded(evidence_id INTEGER PRIMARY KEY); CREATE TABLE IF NOT EXISTS legacy_scope_barriers(scope_hash TEXT PRIMARY KEY,until INTEGER NOT NULL)`); err != nil {
		return err
	}
	defs := map[int64]Definition{}
	rows, err := tx.Query(`SELECT data FROM definitions ORDER BY version`)
	if err != nil {
		return err
	}
	var latest Definition
	for rows.Next() {
		var raw string
		if err = rows.Scan(&raw); err != nil {
			rows.Close()
			return err
		}
		var d Definition
		if err = json.Unmarshal([]byte(raw), &d); err != nil {
			rows.Close()
			return err
		}
		defs[d.Version] = d
		latest = d
	}
	err = rows.Err()
	rows.Close()
	if err != nil {
		return err
	}
	rows, err = tx.Query(`SELECT id,data FROM evidence ORDER BY id LIMIT 100001`)
	if err != nil {
		return err
	}
	evidence := []Evidence{}
	badScopes := map[string]bool{}
	for rows.Next() {
		var id int64
		var raw string
		if err = rows.Scan(&id, &raw); err != nil {
			rows.Close()
			return err
		}
		var x Evidence
		if err = json.Unmarshal([]byte(raw), &x); err != nil {
			rows.Close()
			return err
		}
		x.ID = id
		evidence = append(evidence, x)
		if len(evidence) > 100000 {
			rows.Close()
			return fmt.Errorf("legacy evidence exceeds supported migration bound")
		}
		if x.Source != "browser" && x.Scope != "" {
			badScopes[x.Scope] = true
		}
	}
	err = rows.Err()
	rows.Close()
	if err != nil {
		return err
	}
	until := time.Now().Add(30*24*time.Hour + 5*time.Minute).UnixNano()
	for scope := range badScopes {
		if _, err = tx.Exec(`INSERT OR REPLACE INTO legacy_scope_barriers VALUES(?,?)`, digest(scope), until); err != nil {
			return err
		}
	}
	for _, x := range evidence {
		version := x.CaptureDefinition
		if version == 0 {
			version = x.Definition
		}
		d, ok := defs[version]
		_, compatible := compatibleClientDefinition(d)
		if x.Source != "browser" || badScopes[x.Scope] || !ok || !compatible {
			if _, err = tx.Exec(`INSERT OR IGNORE INTO legacy_excluded VALUES(?)`, x.ID); err != nil {
				return err
			}
		}
	}
	// The next capture has a new prospective definition. Do not strip selectors
	// from an active rule and thereby widen it or relabel a persistence goal.
	hasRoles := false
	for _, m := range latest.Mappings {
		hasRoles = hasRoles || m.Source != ""
	}
	if latest.GoalSelector != nil {
		hasRoles = hasRoles || latest.GoalSelector.Source != ""
	}
	if hasRoles {
		next := Definition{Version: latest.Version + 1, Mappings: []Mapping{}}
		raw, _ := json.Marshal(next)
		if _, err = tx.Exec(`INSERT INTO definitions VALUES(?,?)`, next.Version, string(raw)); err != nil {
			return err
		}
		if _, err = tx.Exec(`INSERT OR REPLACE INTO settings VALUES('client_definition_review','true')`); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *Store) ClientCompatibility() map[string]any {
	if s.ClientCutover == 0 {
		return nil
	}
	var excluded int
	if err := s.DB.QueryRow(`SELECT count(*) FROM legacy_excluded`).Scan(&excluded); err != nil {
		return map[string]any{"message": clientHistoryMessage, "excludedRecords": 0, "cutover": time.Unix(0, s.ClientCutover).UTC().Format(time.RFC3339Nano), "unavailable": true}
	}
	return map[string]any{"message": clientHistoryMessage, "excludedRecords": excluded, "cutover": time.Unix(0, s.ClientCutover).UTC().Format(time.RFC3339Nano)}
}

func (s *Store) Evidence() ([]Evidence, error) {
	// Physical evidence is bounded to 100k by insertion and migration preflight.
	// Eligibility is applied before each caller's smaller source/display cap.
	es, err := s.rawEvidence()
	if err != nil {
		return nil, err
	}
	excluded := map[int64]bool{}
	rows, err := s.DB.Query(`SELECT evidence_id FROM legacy_excluded`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var id int64
		if err = rows.Scan(&id); err != nil {
			rows.Close()
			return nil, err
		}
		excluded[id] = true
	}
	err = rows.Err()
	rows.Close()
	if err != nil {
		return nil, err
	}
	out := make([]Evidence, 0, len(es))
	for _, x := range es {
		if excluded[x.ID] || (x.Source != "client" && x.Source != "browser") {
			continue
		}
		if x.Source == "browser" {
			x.Source = "client"
		}
		out = append(out, x)
	}
	return out, nil
}

// ClientReferenceCompatible rejects only affected historical queries, preserving
// equivalent client-only references and never silently widening their bounds.
func (s *Store) ClientReferenceCompatible(v SavedView) (bool, error) {
	if s.ClientCutover == 0 {
		return true, nil
	}
	defs, err := s.Definitions()
	if err != nil {
		return false, err
	}
	if v.Definition > 0 {
		found := false
		for _, d := range defs {
			if d.Version == v.Definition {
				_, found = compatibleClientDefinition(d)
				break
			}
		}
		if !found {
			return false, nil
		}
	}
	from, to := int64(0), int64(1<<63-1)
	if t, err := time.Parse(time.RFC3339Nano, v.From); err == nil {
		from = t.UnixNano()
	}
	if t, err := time.Parse(time.RFC3339Nano, v.To); err == nil {
		to = t.UnixNano()
	}
	if from >= s.ClientCutover {
		return true, nil
	}
	rows, err := s.DB.Query(`SELECT e.data FROM evidence e JOIN legacy_excluded l ON l.evidence_id=e.id WHERE e.at>=? AND e.at<? LIMIT 100001`, from, to)
	if err != nil {
		return false, err
	}
	affected := false
	for rows.Next() {
		var raw string
		if err = rows.Scan(&raw); err != nil {
			break
		}
		var x Evidence
		if err = json.Unmarshal([]byte(raw), &x); err != nil {
			break
		}
		definition := x.Definition
		if x.CaptureDefinition > 0 {
			definition = x.CaptureDefinition
		}
		if (v.Scope == "" || v.Scope == x.Scope) && (v.Trace == "" || v.Trace == x.Trace) && (v.Version == "" || matchesVersion(x.Version, v.Version)) && (v.Definition == 0 || v.Definition == definition) {
			affected = true
			break
		}
	}
	rowErr := rows.Err()
	rows.Close()
	if err != nil {
		return false, err
	}
	if rowErr != nil {
		return false, rowErr
	}
	if affected {
		return false, nil
	}
	if v.Scope != "" || v.Trace != "" {
		return true, nil
	}
	rows, err = s.DB.Query(`SELECT data FROM summaries`)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	for rows.Next() {
		var raw string
		if err = rows.Scan(&raw); err != nil {
			return false, err
		}
		var d DaySummary
		if err = json.Unmarshal([]byte(raw), &d); err != nil {
			return false, err
		}
		day, err := time.Parse("2006-01-02", d.Day)
		if err != nil {
			return false, err
		}
		if day.UnixNano() >= to || day.Add(24*time.Hour).UnixNano() <= from {
			continue
		}
		for _, g := range d.Groups {
			if g.Computation != "journey-client-v5" {
				return false, nil
			}
		}
	}
	return true, rows.Err()
}
