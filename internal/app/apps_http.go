package app

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

func (s *Service) managedAPI(w http.ResponseWriter, r *http.Request) {
	_, ss, ok := s.authorize(r)
	if !ok {
		http.Error(w, "Session expired or access unavailable", 401)
		return
	}
	if r.Method != "GET" && (r.Header.Get("Origin") != s.Origin || r.Header.Get("X-CSRF-Token") != ss.CSRF) {
		http.Error(w, "Origin and CSRF proof required", 403)
		return
	}
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/"), "/")
	if parts[0] != "apps" {
		http.Error(w, "Select an app and use /api/apps/{appId}/…", 400)
		return
	}
	if len(parts) == 1 {
		if ss.InviteID != 0 {
			http.Error(w, "Owner required", 403)
			return
		}
		switch r.Method {
		case "GET":
			s.Registry.mu.RLock()
			defer s.Registry.mu.RUnlock()
			if s.Registry.uncertain {
				http.Error(w, "Registry recovery required", 503)
				return
			}
			out := []map[string]string{}
			for _, a := range s.Registry.current.Apps {
				out = append(out, a.public())
			}
			write(w, out)
		case "POST":
			var in struct {
				Name string `json:"name"`
			}
			if !decode(w, r, &in) {
				return
			}
			a, e := s.Registry.Register(in.Name)
			if e != nil {
				http.Error(w, e.Error(), 400)
				return
			}
			write(w, a.public())
		default:
			http.Error(w, "GET or POST required", 405)
		}
		return
	}
	id := parts[1]
	action := ""
	if len(parts) > 2 {
		action = strings.Join(parts[2:], "/")
	}
	if generation := r.Header.Get("X-Sentrame-App-Generation"); generation != "" || r.URL.Query().Get("generation") != "" {
		s.Registry.mu.RLock()
		valid := false
		for _, a := range s.Registry.current.Apps {
			if a.ID == id && (generation == "" || a.Generation == generation) && (r.URL.Query().Get("generation") == "" || a.Generation == r.URL.Query().Get("generation")) {
				valid = true
			}
		}
		s.Registry.mu.RUnlock()
		if !valid {
			http.Error(w, "App registration unavailable", 403)
			return
		}
	}
	admin := action == "" || action == "connection" || action == "rename" || action == "rotate" || action == "remove"
	if ss.InviteID != 0 && (admin || ss.AppID != id) {
		http.Error(w, "Access unavailable", 403)
		return
	}
	if admin {
		if action == "connection" && r.Method == "GET" {
			s.Registry.mu.RLock()
			defer s.Registry.mu.RUnlock()
			a, ok := s.Registry.entry(id)
			if !ok {
				http.Error(w, "App unavailable", 404)
				return
			}
			write(w, map[string]string{"browser": a.Browser, "server": a.Server, "endpoint": s.Origin, "name": a.Name})
			return
		}
		if r.Method != "POST" {
			http.Error(w, "POST required", 405)
			return
		}
		var in struct {
			Name         string `json:"name"`
			Role         string `json:"role"`
			Confirmation string `json:"confirmation"`
		}
		if !decode(w, r, &in) {
			return
		}
		var e error
		switch action {
		case "rename":
			e = s.Registry.Update(id, in.Name, "")
		case "rotate":
			if in.Role != "browser" && in.Role != "server" {
				http.Error(w, "Choose browser or server", 400)
				return
			}
			e = s.Registry.Update(id, "", in.Role)
		case "remove":
			e = s.Registry.Delete(id, in.Confirmation)
		default:
			http.NotFound(w, r)
			return
		}
		if e != nil {
			http.Error(w, e.Error(), 409)
			return
		}
		write(w, map[string]bool{"ok": true})
		return
	}
	if action != "context" && action != "snapshot" && action != "timeline" && action != "traces" && !strings.HasPrefix(action, "traces/") && action != "errors" && action != "definitions" && action != "views" && action != "views/delete" && action != "preview" && action != "delete" && action != "events" && action != "invitations" && action != "invitations/revoke" {
		http.NotFound(w, r)
		return
	}
	if action != "events" {
		select {
		case s.queries <- struct{}{}:
			defer func() { <-s.queries }()
		default:
			http.Error(w, "Query budget busy; retry", 503)
			return
		}
	}
	s.Registry.mu.RLock()
	a, ok := s.Registry.entry(id)
	if !ok || (ss.InviteID != 0 && ss.Generation != a.Generation) {
		s.Registry.mu.RUnlock()
		http.Error(w, "Access unavailable", 403)
		return
	}
	if action == "context" {
		if r.Method != "GET" {
			s.Registry.mu.RUnlock()
			http.Error(w, "GET required", 405)
			return
		}
		write(w, a.public())
		s.Registry.mu.RUnlock()
		return
	}
	child := s.forApp(a)
	rr := r.Clone(r.Context())
	u := *r.URL
	rr.URL = &u
	rr.URL.Path = "/api/" + action
	// Streaming checks authority and store availability on every bounded heartbeat;
	// it must not retain the registry read lock while waiting for the next tick.
	if action == "events" {
		s.Registry.mu.RUnlock()
		child.serve(w, rr)
		return
	}
	defer s.Registry.mu.RUnlock()
	child.serve(w, rr)
}
func (s *Service) managedRedemption(w http.ResponseWriter, r *http.Request) {
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
	body, e := io.ReadAll(io.LimitReader(r.Body, 65537))
	if e != nil || len(body) > 65536 {
		http.Error(w, "Invalid request", 400)
		return
	}
	var in struct {
		Token string `json:"token"`
	}
	if json.Unmarshal(body, &in) != nil {
		http.Error(w, "Invalid request", 400)
		return
	}
	s.Registry.mu.RLock()
	defer s.Registry.mu.RUnlock()
	if !s.Registry.uncertain {
		for _, a := range s.Registry.current.Apps {
			if a.State != "active" {
				continue
			}
			store := s.Registry.stores[a.ID]
			var n int
			if store.DB.QueryRow(`SELECT count(*) FROM invitations WHERE token_hash=?`, digest(in.Token)).Scan(&n) == nil && n == 1 {
				r.Body = io.NopCloser(bytes.NewReader(body))
				child := s.forApp(a)
				child.redemptionReserved = true
				child.redeemInvitation(w, r)
				return
			}
		}
	}
	http.Error(w, "Invitation unavailable; ask the owner for a new link", 401)
}
