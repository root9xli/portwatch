package monitor

import (
	"fmt"
	"strings"
	"time"
)

// BaselineReporter produces human-readable summaries of a Baseline.
type BaselineReporter struct {
	baseline *Baseline
}

// NewBaselineReporter creates a reporter for the given baseline.
func NewBaselineReporter(b *Baseline) *BaselineReporter {
	return &BaselineReporter{baseline: b}
}

// Summary returns a formatted string listing all baselined ports.
func (r *BaselineReporter) Summary() string {
	r.baseline.mu.RLock()
	defer r.baseline.mu.RUnlock()

	if len(r.baseline.entries) == 0 {
		return "baseline: no ports recorded"
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "baseline: %d port(s)\n", len(r.baseline.entries))
	for _, e := range r.baseline.entries {
		fmt.Fprintf(&sb, "  port=%-6d first_seen=%s\n", e.Port, e.FirstSeen.Format(time.RFC3339))
	}
	return strings.TrimRight(sb.String(), "\n")
}

// ContainsPort returns a one-line string indicating whether a port is baselined.
func (r *BaselineReporter) ContainsPort(port int) string {
	if r.baseline.Contains(port) {
		return fmt.Sprintf("port %d is in baseline", port)
	}
	return fmt.Sprintf("port %d is NOT in baseline", port)
}
