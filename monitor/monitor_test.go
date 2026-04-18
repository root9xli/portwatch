package monitor

import (
	"testing"
)

func TestIsAllowed(t *testing.T) {
	cfg := &Config{AllowedPorts: []int{80, 443, 8080}}
	if !cfg.IsAllowed(80) {
		t.Error("expected 80 to be allowed")
	}
	if !cfg.IsAllowed(443) {
		t.Error("expected 443 to be allowed")
	}
	if cfg.IsAllowed(9000) {
		t.Error("expected 9000 to not be allowed")
	}
}

func TestMonitorSkipsAllowed(t *testing.T) {
	cfg := &Config{AllowedPorts: []int{22, 80}}
	m := New(cfg)
	// Simulate already-seen unexpected port to ensure no double alert.
	m.seen[9999] = true

	listeners := []Listener{
		{Port: 80, PID: 1, Process: "nginx", Proto: "tcp"},
		{Port: 9999, PID: 2, Process: "unknown", Proto: "tcp"},
	}

	for _, l := range listeners {
		if cfg.IsAllowed(l.Port) {
			continue
		}
		if !m.seen[l.Port] {
			t.Errorf("expected port %d to already be seen", l.Port)
		}
	}
}

func TestNewMonitor(t *testing.T) {
	cfg := &Config{}
	m := New(cfg)
	if m == nil {
		t.Fatal("expected non-nil monitor")
	}
	if m.seen == nil {
		t.Fatal("expected initialized seen map")
	}
}
