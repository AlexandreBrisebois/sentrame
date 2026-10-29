package app

import (
	"database/sql"
	"encoding/hex"
	"net/http"
	"strings"
	"time"
)

type Invitation struct {
	ID       int64     `json:"id"`
	Label    string    `json:"label"`
	Created  time.Time `json:"created"`
	Expires  time.Time `json:"expires"`
	Redeemed bool      `json:"redeemed"`
	Revoked  bool      `json:"revoked"`
}

func sessionRole(ss session) string {
	if ss.InviteID != 0 {
		return "viewer"
	}
	return "owner"
}
func sessionExpiry(ss session) time.Time {
	if ss.InviteID != 0 {
		return ss.Expires
	}
	return ss.Created.Add(8 * time.Hour)
}
func sessionPayload(ss session) map[string]string {
	return map[string]string{"csrf": ss.CSRF, "role": sessionRole(ss), "expires": sessionExpiry(ss).UTC().Format(time.RFC3339Nano)}
}
func viewerReadAllowed(r *http.Request) bool {
	if r.Method == "POST" && r.URL.Path == "/api/logout" {
		return true
	}
	if r.Method != "GET" {
		return false
	}
	switch r.URL.Path {
	case "/api/session", "/api/snapshot", "/api/timeline", "/api/errors", "/api/views", "/api/events", "/api/definitions":
		return true
	}
	return false
}

// Caller holds Mu. Session storage is bounded without evicting active viewers.
func (s *Service) pruneSessionsLocked() {
	now := s.Store.Now()
	for key, ss := range s.sessions {
		if ss.InviteID != 0 && !now.Before(ss.Expires) {
			delete(s.sessions, key)
		} else if ss.InviteID == 0 && (time.Since(ss.Last) > 30*time.Minute || time.Since(ss.Created) > 8*time.Hour) {
			delete(s.sessions, key)
		}
	}
}
func (s *Service) invitationValid(id int64) bool {
	var n int
	err := s.Store.DB.QueryRow(`SELECT count(*) FROM invitations WHERE id=? AND redeemed=1 AND revoked=0 AND expires>? AND app_id=(SELECT value FROM settings WHERE key='app_id') AND generation=(SELECT value FROM settings WHERE key='generation') AND NOT EXISTS(SELECT 1 FROM settings WHERE key='disabled' AND value='true')`, id, s.Store.Now().UnixNano()).Scan(&n)
	return err == nil && n == 1
}
func (s *Service) invitations(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		rows, err := s.Store.DB.Query(`SELECT id,label,created,expires,redeemed,revoked FROM invitations ORDER BY id DESC LIMIT 200`)
		if err != nil {
			http.Error(w, "Invitation store unavailable", 503)
			return
		}
		defer rows.Close()
		out := []Invitation{}
		for rows.Next() {
			var v Invitation
			var created, expires int64
			if err := rows.Scan(&v.ID, &v.Label, &created, &expires, &v.Redeemed, &v.Revoked); err != nil {
				http.Error(w, "Invitation store unavailable", 503)
				return
			}
			v.Created = time.Unix(0, created).UTC()
			v.Expires = time.Unix(0, expires).UTC()
			out = append(out, v)
		}
		if rows.Err() != nil {
			http.Error(w, "Invitation store unavailable", 503)
			return
		}
		write(w, out)
	case "POST":
		var input struct {
			Label string `json:"label"`
			Hours int    `json:"hours"`
		}
		if !decode(w, r, &input) {
			return
		}
		if input.Hours == 0 {
			input.Hours = 24
		}
		if input.Hours != 1 && input.Hours != 24 && input.Hours != 168 {
			http.Error(w, "Choose 1, 24 or 168 hours", 400)
			return
		}
		if strings.TrimSpace(input.Label) == "" {
			input.Label = "Temporary viewer"
		}
		label := safe(input.Label)
		if label == "" || len(input.Label) > 100 || label != input.Label {
			http.Error(w, "Use a label of 1 to 100 safe characters", 400)
			return
		}
		now := s.Store.Now()
		expires := now.Add(time.Duration(input.Hours) * time.Hour)
		token := random()
		s.Store.Mu.Lock()
		defer s.Store.Mu.Unlock()
		tx, err := s.Store.DB.Begin()
		if err != nil {
			http.Error(w, "Invitation writer unavailable", 503)
			return
		}
		defer tx.Rollback()
		var disabled string
		err = tx.QueryRow(`SELECT value FROM settings WHERE key='disabled'`).Scan(&disabled)
		if err != nil && err != sql.ErrNoRows {
			http.Error(w, "Invitation writer unavailable", 503)
			return
		}
		if disabled == "true" {
			http.Error(w, "Application is disabled", 409)
			return
		}
		// Retain at most 100 inactive audit rows alongside at most 100 active grants.
		if _, err = tx.Exec(`DELETE FROM invitations WHERE (revoked=1 OR expires<=?) AND id NOT IN (SELECT id FROM invitations WHERE revoked=1 OR expires<=? ORDER BY id DESC LIMIT 99)`, now.UnixNano(), now.UnixNano()); err != nil {
			http.Error(w, "Invitation writer unavailable", 503)
			return
		}
		var count int
		if err = tx.QueryRow(`SELECT count(*) FROM invitations WHERE revoked=0 AND expires>?`, now.UnixNano()).Scan(&count); err != nil {
			http.Error(w, "Invitation writer unavailable", 503)
			return
		}
		if count >= 100 {
			http.Error(w, "Active invitation limit 100", 400)
			return
		}
		result, err := tx.Exec(`INSERT INTO invitations(label,created,expires,token_hash,app_id,generation) SELECT ?,?,?,?,(SELECT value FROM settings WHERE key='app_id'),(SELECT value FROM settings WHERE key='generation')`, label, now.UnixNano(), expires.UnixNano(), digest(token))
		if err != nil {
			http.Error(w, "Invitation writer unavailable", 503)
			return
		}
		id, err := result.LastInsertId()
		if err != nil || tx.Commit() != nil {
			http.Error(w, "Invitation writer unavailable", 503)
			return
		}
		v := Invitation{ID: id, Label: label, Created: now.UTC(), Expires: expires.UTC()}
		write(w, map[string]any{"invite": v, "url": s.Origin + "/#invite=" + token})
	default:
		http.Error(w, "GET or POST required", 405)
	}
}
func (s *Service) revokeInvitation(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "POST required", 405)
		return
	}
	var input struct {
		ID int64 `json:"id"`
	}
	if !decode(w, r, &input) {
		return
	}
	if input.ID <= 0 {
		http.Error(w, "Invitation ID required", 400)
		return
	}
	result, err := s.Store.DB.Exec(`UPDATE invitations SET revoked=1 WHERE id=?`, input.ID)
	if err != nil {
		http.Error(w, "Invitation writer unavailable", 503)
		return
	}
	n, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Invitation writer unavailable", 503)
		return
	}
	if n == 0 {
		http.Error(w, "Invitation not found", 404)
		return
	}
	// Reclaim capacity immediately; disconnected revoked viewers may never read again.
	s.Mu.Lock()
	for key, ss := range s.sessions {
		if ss.InviteID == input.ID {
			delete(s.sessions, key)
		}
	}
	s.Mu.Unlock()
	write(w, map[string]bool{"ok": true})
}
func (s *Service) redeemInvitation(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "POST required", 405)
		return
	}
	if r.Header.Get("Origin") != s.Origin {
		http.Error(w, "Origin required", 403)
		return
	}
	select {
	case s.logins <- struct{}{}:
		defer func() { <-s.logins }()
	default:
		http.Error(w, "Invitation busy; retry", 429)
		return
	}
	var input struct {
		Token string `json:"token"`
	}
	if !decode(w, r, &input) {
		return
	}
	b, err := hex.DecodeString(input.Token)
	if err != nil || len(b) != 32 {
		http.Error(w, "Invitation unavailable; ask the owner for a new link", 401)
		return
	}
	now := s.Store.Now()
	var id, expires int64
	// Session capacity, atomic one-use consumption and registration are one critical
	// section. authorize never holds Mu while querying SQLite.
	s.Mu.Lock()
	defer s.Mu.Unlock()
	s.pruneSessionsLocked()
	if len(s.sessions) >= 128 {
		http.Error(w, "Viewer session limit; retry later", 503)
		return
	}
	tx, err := s.Store.DB.Begin()
	if err != nil {
		http.Error(w, "Invitation writer unavailable", 503)
		return
	}
	defer tx.Rollback()
	err = tx.QueryRow(`UPDATE invitations SET redeemed=1 WHERE token_hash=? AND redeemed=0 AND revoked=0 AND expires>? AND app_id=(SELECT value FROM settings WHERE key='app_id') AND generation=(SELECT value FROM settings WHERE key='generation') AND NOT EXISTS(SELECT 1 FROM settings WHERE key='disabled' AND value='true') RETURNING id,expires`, digest(input.Token), now.UnixNano()).Scan(&id, &expires)
	if err == sql.ErrNoRows {
		http.Error(w, "Invitation unavailable; ask the owner for a new link", 401)
		return
	}
	if err != nil {
		http.Error(w, "Invitation writer unavailable", 503)
		return
	}
	if err = tx.Commit(); err != nil {
		http.Error(w, "Invitation writer unavailable", 503)
		return
	}
	token, csrf := random(), random()
	ss := session{CSRF: csrf, Created: now, Last: now, InviteID: id, Expires: time.Unix(0, expires)}
	s.sessions[digest(token)] = ss
	http.SetCookie(w, &http.Cookie{Name: "sentrame_viewer", Value: token, Path: "/", HttpOnly: true, Secure: strings.HasPrefix(s.Origin, "https://"), SameSite: http.SameSiteStrictMode, MaxAge: int(ss.Expires.Sub(now).Seconds()), Expires: ss.Expires})
	write(w, sessionPayload(ss))
}
