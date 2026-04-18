package monitor

import (
	"fmt"
	"strings"
	"time"
)

// MetricsReporter formats Metrics into human-readable summaries.
type MetricsReporter struct {
	metrics *Metrics
}

// NewMetricsReporter creates a reporter backed by the given Metrics.
func NewMetricsReporter(m *Metrics) *MetricsReporter {
	return &MetricsReporter{metrics: m}
}

// Summary returns a multi-line string with all key metrics.
func (r *MetricsReporter) Summary() string {
	snap := r.metrics.Snapshot()
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("uptime:       %s\n", fmtDuration(time.Since(snap.StartTime))))
	sb.WriteString(fmt.Sprintf("diffs:        %d\n", snap.DiffsTotal))
	sb.WriteString(fmt.Sprintf("alerts:       %d\n", snap.AlertsTotal))
	sb.WriteString(fmt.Sprintf("suppressed:   %d\n", snap.Suppressed))
	sb.WriteString(fmt.Sprintf("rate_limited: %d\n", snap.RateLimited))
	if !snap.lastDiff.IsZero() {
		sb.WriteString(fmt.Sprintf("last_diff:    %s\n", snap.lastDiff.Format(time.RFC3339)))
	} else {
		sb.WriteString("last_diff:    never\n")
	}
	return sb.String()
}

// fmtDuration formats a duration as a concise string (e.g. "2h3m4s").
func fmtDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	if h > 0 {
		return fmt.Sprintf("%dh%dm%ds", h, m, s)
	}
	if m > 0 {
		return fmt.Sprintf("%dm%ds", m, s)
	}
	return fmt.Sprintf("%ds", s)
}
