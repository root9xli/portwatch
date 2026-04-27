package monitor

import (
	"fmt"
	"sort"
	"strings"
)

// PortFreezeReporter produces human-readable summaries of the current
// freeze state tracked by a PortFreeze instance.
type PortFreezeReporter struct {
	freeze *PortFreeze
}

// NewPortFreezeReporter creates a reporter backed by the given PortFreeze.
func NewPortFreezeReporter(pf *PortFreeze) *PortFreezeReporter {
	return &PortFreezeReporter{freeze: pf}
}

// Summary returns a multi-line string listing all currently frozen ports.
func (r *PortFreezeReporter) Summary() string {
	r.freeze.mu.RLock()
	defer r.freeze.mu.RUnlock()

	now := r.freeze.clock()
	type row struct {
		port   int
		entry  *FreezeEntry
	}
	var rows []row
	for port, e := range r.freeze.entries {
		if now.Before(e.ExpiresAt) {
			rows = append(rows, row{port, e})
		}
	}
	if len(rows) == 0 {
		return "frozen ports: none"
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].port < rows[j].port })

	var sb strings.Builder
	sb.WriteString("frozen ports:\n")
	for _, row := range rows {
		remaining := row.entry.ExpiresAt.Sub(now).Round(1e9)
		sb.WriteString(fmt.Sprintf("  port %-6d reason=%-20s expires_in=%s\n",
			row.port, row.entry.Reason, remaining))
	}
	return strings.TrimRight(sb.String(), "\n")
}

// Count returns the number of currently active (non-expired) freezes.
func (r *PortFreezeReporter) Count() int {
	r.freeze.mu.RLock()
	defer r.freeze.mu.RUnlock()
	now := r.freeze.clock()
	count := 0
	for _, e := range r.freeze.entries {
		if now.Before(e.ExpiresAt) {
			count++
		}
	}
	return count
}
