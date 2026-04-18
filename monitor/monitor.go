package monitor

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// Listener represents a process listening on a port.
type Listener struct {
	Port    int
	PID     int
	Process string
	Proto   string
}

// Monitor polls port usage and alerts on unexpected listeners.
type Monitor struct {
	cfg  *Config
	seen map[int]bool
}

// New creates a new Monitor.
func New(cfg *Config) *Monitor {
	return &Monitor{cfg: cfg, seen: make(map[int]bool)}
}

// Run performs a single scan cycle.
func (m *Monitor) Run() error {
	listeners, err := scanListeners()
	if err != nil {
		return fmt.Errorf("scan: %w", err)
	}

	for _, l := range listeners {
		if m.cfg.IsAllowed(l.Port) {
			continue
		}
		if !m.seen[l.Port] {
			m.seen[l.Port] = true
			log.Printf("ALERT: unexpected listener port=%d pid=%d process=%s proto=%s",
				l.Port, l.PID, l.Process, l.Proto)
			m.runAlert(l)
		}
	}
	return nil
}

func (m *Monitor) runAlert(l Listener) {
	if m.cfg.AlertCommand == "" {
		return
	}
	cmd := strings.ReplaceAll(m.cfg.AlertCommand, "{port}", fmt.Sprintf("%d", l.Port))
	cmd = strings.ReplaceAll(cmd, "{pid}", fmt.Sprintf("%d", l.PID))
	parts := strings.Fields(cmd)
	if err := exec.Command(parts[0], parts[1:]...).Run(); err != nil {
		log.Printf("alert command failed: %v", err)
	}
}
