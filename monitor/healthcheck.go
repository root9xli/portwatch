package monitor

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"time"
)

// HealthCheck exposes a simple HTTP endpoint reporting daemon liveness.
type HealthCheck struct {
	server  *http.Server
	started time.Time
	cycles  atomic.Int64
}

// NewHealthCheck creates a HealthCheck listening on the given address (e.g. ":9090").
func NewHealthCheck(addr string) *HealthCheck {
	h := &HealthCheck{started: time.Now()}
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", h.handleHealth)
	mux.HandleFunc("/readyz", h.handleReady)
	h.server = &http.Server{Addr: addr, Handler: mux}
	return h
}

// RecordCycle increments the monitor cycle counter.
func (h *HealthCheck) RecordCycle() {
	h.cycles.Add(1)
}

// Start begins serving in a background goroutine.
func (h *HealthCheck) Start() {
	go func() { _ = h.server.ListenAndServe() }()
}

// Close shuts down the HTTP server.
func (h *HealthCheck) Close() error {
	return h.server.Close()
}

func (h *HealthCheck) handleHealth(w http.ResponseWriter, _ *http.Request) {
	uptime := time.Since(h.started).Truncate(time.Second)
	fmt.Fprintf(w, `{"status":"ok","uptime":%q,"cycles":%d}`, uptime, h.cycles.Load())
}

func (h *HealthCheck) handleReady(w http.ResponseWriter, _ *http.Request) {
	if h.cycles.Load() == 0 {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprint(w, `{"ready":false}`)
		return
	}
	fmt.Fprint(w, `{"ready":true}`)
}
