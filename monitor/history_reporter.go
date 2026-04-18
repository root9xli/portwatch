package monitor

import (
	"fmt"
	"strings"
	"time"
)

// Reporter formats History entries for display or logging.
type Reporter struct {
	history *History
}

// NewReporter wraps a History for reporting.
func NewReporter(h *History) *Reporter {
	return &Reporter{history: h}
}

// Summary returns a human-readable summary of the last n events.
func (r *Reporter) Summary(n int) string {
	entries := r.history.Last(n)
	if len(entries) == 0 {
		return "no recent events"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d recent event(s):\n", len(entries)))
	for _, e := range entries {
		sb.WriteString(fmt.Sprintf(
			"  [%s] port=%d action=%s level=%s\n",
			e.Timestamp.Format(time.RFC3339),
			e.Port,
			e.Action,
			e.Level,
		))
	}
	return sb.String()
}

// CountByLevel returns a map of level -> count across all stored entries.
func (r *Reporter) CountByLevel() map[string]int {
	counts := make(map[string]int)
	for _, e := range r.history.Last(r.history.Len()) {
		counts[e.Level]++
	}
	return counts
}

// HasCritical returns true if any stored entry has level "critical".
func (r *Reporter) HasCritical() bool {
	for _, e := range r.history.Last(r.history.Len()) {
		if e.Level == "critical" {
			return true
		}
	}
	return false
}
