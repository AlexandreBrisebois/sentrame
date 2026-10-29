package app

import (
	"bytes"
	"compress/gzip"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	statuspb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func random() string         { b := make([]byte, 32); rand.Read(b); return hex.EncodeToString(b) }
func digest(s string) string { x := sha256.Sum256([]byte(s)); return hex.EncodeToString(x[:]) }
func write(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}
func decode(w http.ResponseWriter, r *http.Request, v any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, 65536)
	if e := json.NewDecoder(r.Body).Decode(v); e != nil {
		http.Error(w, "Invalid request", 400)
		return false
	}
	return true
}
func (s *Service) Handler(static string) http.Handler {
	if static != "" {
		s.static = http.FileServer(http.Dir(static))
	}
	return http.HandlerFunc(s.serve)
}
func (s *Service) serve(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Referrer-Policy", "no-referrer")
	w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; connect-src 'self'; img-src 'self' data:; frame-ancestors 'none'")
	if strings.HasPrefix(r.URL.Path, "/v1/") {
		s.ingest(w, r)
		return
	}
	if r.URL.Path == "/ready" {
		s.readiness(w, r)
		return
	}
	if r.URL.Path == "/health" {
		write(w, map[string]string{"status": "local synthetic"})
		return
	}
	if r.URL.Path == "/api/login" {
		s.login(w, r)
		return
	}
	if r.URL.Path == "/api/invitations/redeem" {
		s.redeemInvitation(w, r)
		return
	}
	if strings.HasPrefix(r.URL.Path, "/api/") {
		sessionKey, ss, ok := s.authorize(r)
		if !ok {
			http.Error(w, "Session expired or access unavailable. Sign in or ask the owner for a new invite.", 401)
			return
		}
		if ss.InviteID != 0 && !viewerReadAllowed(r) {
			http.Error(w, "Read-only viewer access", 403)
			return
		}
		if r.Method != "GET" {
			if r.Header.Get("Origin") != s.Origin || r.Header.Get("X-CSRF-Token") != ss.CSRF {
				http.Error(w, "Origin and CSRF proof required", 403)
				return
			}
		}
		switch r.URL.Path {
		case "/api/session":
			write(w, sessionPayload(ss))
		case "/api/logout":
			if r.Method != "POST" {
				http.Error(w, "POST required", 405)
				return
			}
			s.Mu.Lock()
			delete(s.sessions, sessionKey)
			s.Mu.Unlock()
			cookieName := "sentrame_owner"
			if ss.InviteID != 0 {
				cookieName = "sentrame_viewer"
			}
			http.SetCookie(w, &http.Cookie{Name: cookieName, Value: "", Path: "/", MaxAge: -1, HttpOnly: true, Secure: strings.HasPrefix(s.Origin, "https://"), SameSite: http.SameSiteStrictMode})
			write(w, map[string]bool{"ok": true})
		case "/api/invitations":
			s.invitations(w, r)
		case "/api/invitations/revoke":
			s.revokeInvitation(w, r)
		case "/api/errors":
			s.errors(w, r)
		case "/api/timeline":
			s.timeline(w, r)
		case "/api/snapshot":
			s.snapshot(w, r)
		case "/api/definitions":
			s.definitions(w, r)
		case "/api/preview":
			s.preview(w, r)
		case "/api/views":
			s.views(w, r)
		case "/api/views/delete":
			s.deleteView(w, r)
		case "/api/delete":
			s.deletion(w, r)
		case "/api/events":
			s.events(w, r)
		default:
			http.NotFound(w, r)
		}
		return
	}
	if s.static != nil {
		s.static.ServeHTTP(w, r)
	} else {
		http.NotFound(w, r)
	}
}
func (s *Service) login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "POST required", 405)
		return
	}
	if r.Header.Get("Origin") != s.Origin {
		http.Error(w, "Origin required", 403)
		return
	}
	s.Mu.Lock()
	now := time.Now()
	kept := s.failures[:0]
	for _, t := range s.failures {
		if now.Sub(t) < time.Minute {
			kept = append(kept, t)
		}
	}
	s.failures = kept
	if len(kept) >= 5 {
		s.Mu.Unlock()
		http.Error(w, "Wait one minute before retrying", 429)
		return
	}
	s.Mu.Unlock()
	select {
	case s.logins <- struct{}{}:
		defer func() { <-s.logins }()
	default:
		http.Error(w, "Login busy; retry", 429)
		return
	}
	var v struct {
		Passphrase string `json:"passphrase"`
	}
	if !decode(w, r, &v) {
		return
	}
	if bcrypt.CompareHashAndPassword(s.password, []byte(v.Passphrase)) != nil {
		s.Mu.Lock()
		s.failures = append(s.failures, now)
		s.Mu.Unlock()
		http.Error(w, "Passphrase not accepted", 401)
		return
	}
	token, csrf := random(), random()
	s.Mu.Lock()
	s.pruneSessionsLocked()
	if len(s.sessions) >= 128 {
		s.Mu.Unlock()
		http.Error(w, "Viewer session limit; retry later", 503)
		return
	}
	s.sessions[digest(token)] = session{CSRF: csrf, Created: now, Last: now}
	s.Mu.Unlock()
	http.SetCookie(w, &http.Cookie{Name: "sentrame_owner", Value: token, Path: "/", HttpOnly: true, Secure: strings.HasPrefix(s.Origin, "https://"), SameSite: http.SameSiteStrictMode, MaxAge: 28800})
	write(w, sessionPayload(session{CSRF: csrf, Created: now, Last: now}))
}
func (s *Service) authorize(r *http.Request) (string, session, bool) {
	for _, name := range []string{"sentrame_owner", "sentrame_viewer"} {
		c, err := r.Cookie(name)
		if err != nil {
			continue
		}
		key := digest(c.Value)
		s.Mu.Lock()
		ss, ok := s.sessions[key]
		if !ok || (name == "sentrame_owner") != (ss.InviteID == 0) {
			s.Mu.Unlock()
			continue
		}
		if ss.InviteID == 0 {
			now := time.Now()
			if now.Sub(ss.Last) > 30*time.Minute || now.Sub(ss.Created) > 8*time.Hour {
				delete(s.sessions, key)
				s.Mu.Unlock()
				continue
			}
			ss.Last = now
			s.sessions[key] = ss
			s.Mu.Unlock()
			return key, ss, true
		}
		s.Mu.Unlock()
		if !s.Store.Now().Before(ss.Expires) || !s.invitationValid(ss.InviteID) {
			s.Mu.Lock()
			delete(s.sessions, key)
			s.Mu.Unlock()
			continue
		}
		return key, ss, true
	}
	return "", session{}, false
}
func (s *Service) ingest(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		otlpError(w, r, 405, "POST required")
		return
	}
	if _, e := s.producer(r.Header.Get("Authorization")); e != nil {
		otlpError(w, r, 401, "Producer credential required")
		return
	}
	switch r.URL.Path {
	case "/v1/logs", "/v1/traces", "/v1/metrics":
	default:
		otlpError(w, r, 404, "Signal not supported")
		return
	}
	select {
	case s.slots <- struct{}{}:
		defer func() { <-s.slots }()
	default:
		otlpError(w, r, 503, "Intake busy; retry")
		return
	}
	b, e := io.ReadAll(io.LimitReader(r.Body, (1<<20)+1))
	if e != nil || len(b) > 1<<20 {
		otlpError(w, r, 400, "Request limit")
		return
	}
	switch r.Header.Get("Content-Encoding") {
	case "":
	case "gzip":
		gz, err := gzip.NewReader(bytes.NewReader(b))
		if err != nil {
			otlpError(w, r, 400, "Invalid gzip")
			return
		}
		b, e = io.ReadAll(io.LimitReader(gz, (4<<20)+1))
		gz.Close()
		if e != nil || len(b) > 4<<20 {
			otlpError(w, r, 400, "Decompression limit")
			return
		}
		r.Header.Del("Content-Encoding")
	default:
		otlpError(w, r, 415, "Encoding not supported")
		return
	}
	r.Body = io.NopCloser(bytes.NewReader(b))
	r.ContentLength = int64(len(b))
	target, _ := url.Parse(s.ReceiverURL)
	p := httputil.NewSingleHostReverseProxy(target)
	p.ModifyResponse = func(resp *http.Response) error {
		if resp.StatusCode >= 400 {
			resp.Body.Close()
			msg, ct := otlpStatus(r.Header.Get("Content-Type"), resp.StatusCode, "OTLP request rejected; inspect status and sanitized collection health")
			resp.Body = io.NopCloser(bytes.NewReader(msg))
			resp.ContentLength = int64(len(msg))
			resp.Header.Set("Content-Length", strconv.Itoa(len(msg)))
			resp.Header.Set("Content-Type", ct)
		}
		return nil
	}
	p.ErrorHandler = func(w http.ResponseWriter, r *http.Request, e error) { otlpError(w, r, 503, "Receiver unavailable") }
	p.ServeHTTP(w, r)
}
func (s *Service) events(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "GET required", 405)
		return
	}
	select {
	case s.streams <- struct{}{}:
		defer func() { <-s.streams }()
	default:
		http.Error(w, "Stream limit; retry", 503)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("X-Accel-Buffering", "no")
	fl, ok := w.(http.Flusher)
	if !ok {
		return
	}
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		if _, _, ok := s.authorize(r); !ok {
			return
		}
		var rev string
		s.Store.DB.QueryRow(`SELECT value FROM settings WHERE key='revision'`).Scan(&rev)
		http.NewResponseController(w).SetWriteDeadline(time.Now().Add(5 * time.Second))
		if _, e := fmt.Fprintf(w, "event: reset\nid: %s\ndata: {\"revision\":%q}\n\n", rev, rev); e != nil {
			return
		}
		fl.Flush()
		select {
		case <-r.Context().Done():
			return
		case <-ticker.C:
		}
	}
}

func otlpStatus(contentType string, httpCode int, msg string) ([]byte, string) {
	codes := map[int]int32{400: 3, 401: 16, 403: 7, 404: 12, 405: 12, 413: 3, 415: 3, 429: 8, 500: 13, 503: 14}
	code := codes[httpCode]
	if code == 0 {
		code = 13
	}
	v := &statuspb.Status{Code: code, Message: msg}
	if strings.HasPrefix(contentType, "application/x-protobuf") {
		b, _ := proto.Marshal(v)
		return b, "application/x-protobuf"
	}
	b, _ := protojson.Marshal(v)
	return b, "application/json"
}
func otlpError(w http.ResponseWriter, r *http.Request, httpCode int, msg string) {
	b, ct := otlpStatus(r.Header.Get("Content-Type"), httpCode, msg)
	w.Header().Set("Content-Type", ct)
	w.WriteHeader(httpCode)
	w.Write(b)
}
