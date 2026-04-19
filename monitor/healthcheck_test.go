package monitor

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestHealthCheck() *HealthCheck {
	return NewHealthCheck(":0")
}

func TestHealthCheckHealthzOK(t *testing.T) {
	h := newTestHealthCheck()
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	h.handleHealth(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &m); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if m["status"] != "ok" {
		t.Errorf("expected status ok, got %v", m["status"])
	}
}

func TestHealthCheckReadyzNotReady(t *testing.T) {
	h := newTestHealthCheck()
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	h.handleReady(rr, req)
	if rr.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503 before first cycle, got %d", rr.Code)
	}
}

func TestHealthCheckReadyzReady(t *testing.T) {
	h := newTestHealthCheck()
	h.RecordCycle()
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	h.handleReady(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 after cycle, got %d", rr.Code)
	}
}

func TestHealthCheckCycleCounter(t *testing.T) {
	h := newTestHealthCheck()
	for i := 0; i < 5; i++ {
		h.RecordCycle()
	}
	if h.cycles.Load() != 5 {
		t.Errorf("expected 5 cycles, got %d", h.cycles.Load())
	}
}
