package app

import (
	"context"
	"net/http"
	"time"
)

// Readiness covers initialization and the local SQLite connection. It is not
// a statement about external exporter delivery or routed transport support.
func (s *Service) readiness(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "GET required", http.StatusMethodNotAllowed)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	ready := s.ReceiverURL != ""
	if s.Registry != nil {
		s.Registry.mu.RLock()
		ready = ready && !s.Registry.uncertain
		for _, store := range s.Registry.stores {
			ready = ready && store.DB.PingContext(ctx) == nil
		}
		s.Registry.mu.RUnlock()
	} else {
		ready = ready && s.Store.DB.PingContext(ctx) == nil
	}
	if !ready {
		http.Error(w, "Service not ready", http.StatusServiceUnavailable)
		return
	}
	write(w, map[string]string{"status": "ready"})
}
