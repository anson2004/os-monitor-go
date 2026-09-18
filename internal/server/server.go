// Package server exposes collected metrics over HTTP.
package server

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"osmonitor/internal/collector"
)

// Server serves the latest Snapshot. A background loop refreshes it on a
// fixed interval so HTTP requests never block on CPU sampling.
type Server struct {
	collector *collector.Collector
	interval  time.Duration
	log       *slog.Logger

	mu     sync.RWMutex
	latest collector.Snapshot
}

// New creates a Server that refreshes metrics every interval.
func New(c *collector.Collector, interval time.Duration, log *slog.Logger) *Server {
	return &Server{collector: c, interval: interval, log: log}
}

// Run refreshes the snapshot until ctx is cancelled.
func (s *Server) Run(ctx context.Context) {
	s.refresh(ctx)
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.refresh(ctx)
		}
	}
}

func (s *Server) refresh(ctx context.Context) {
	snap := s.collector.Collect(ctx)
	s.mu.Lock()
	s.latest = snap
	s.mu.Unlock()

	s.log.Info("metrics refreshed",
		"cpu_percent", snap.CPU.UsagePercent,
		"mem_percent", snap.Memory.UsedPercent,
		"sensors", len(snap.Temperature),
		"errors", len(snap.Errors),
	)
}

// Handler returns the HTTP routes.
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", s.handleHealth)
	mux.HandleFunc("GET /api/metrics", s.handleMetrics)
	mux.HandleFunc("GET /api/metrics/live", s.handleLive)
	return mux
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// handleMetrics returns the cached snapshot.
func (s *Server) handleMetrics(w http.ResponseWriter, _ *http.Request) {
	s.mu.RLock()
	snap := s.latest
	s.mu.RUnlock()
	writeJSON(w, http.StatusOK, snap)
}

// handleLive collects a fresh snapshot on demand (slower, blocks on sampling).
func (s *Server) handleLive(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.collector.Collect(r.Context()))
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(v)
}
