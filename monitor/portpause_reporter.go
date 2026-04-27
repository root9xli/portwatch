package monitor

import (
	"fmt"
	"sort"
	"strings"
)

// PortPauseReporter produces human-readable summaries of paused ports.
type PortPauseReporter struct {
	pause *PortPause
}

// NewPortPauseReporter creates a reporter backed by the given PortPause.
func NewPortPauseReporter(p *PortPause) *PortPauseReporter {
	return &PortPauseReporter{pause: p}
}

// Summary returns a formatted string listing all active paused ports.
func (r *PortPauseReporter) Summary() string {
	r.pause.mu.Lock()
	defer r.pause.mu.Unlock()

	now := r.pause.now()
	type row struct {
		port   int
		reason string
		ttl    string
	}

	var rows []row
	for port, e := range r.pause.entries {
		if now.After(e.ExpiresAt) {
			continue
		}
		remaining := e.ExpiresAt.Sub(now).Round(1e9)
		rows = append(rows, row{
			port:   port,
			reason: e.Reason,
			ttl:    remaining.String(),
		})
	}

	if len(rows) == 0 {
		return "paused ports: none"
	}

	sort.Slice(rows, func(i, j int) bool { return rows[i].port < rows[j].port })

	var sb strings.Builder
	sb.WriteString("paused ports:\n")
	for _, row := range rows {
		sb.WriteString(fmt.Sprintf("  port=%d reason=%q ttl=%s\n", row.port, row.reason, row.ttl))
	}
	return strings.TrimRight(sb.String(), "\n")
}

// Count returns the number of currently active (non-expired) paused ports.
func (r *PortPauseReporter) Count() int {
	r.pause.mu.Lock()
	defer r.pause.mu.Unlock()
	now := r.pause.now()
	count := 0
	for _, e := range r.pause.entries {
		if !now.After(e.ExpiresAt) {
			count++
		}
	}
	return count
}
