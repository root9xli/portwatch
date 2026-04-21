package monitor

import (
	"fmt"
	"sort"
	"strings"
)

// PortTrendReporter summarises trend directions for a set of ports.
type PortTrendReporter struct {
	trend  *PortTrend
	labels *PortLabeler
}

// NewPortTrendReporter creates a reporter backed by the given PortTrend and PortLabeler.
func NewPortTrendReporter(trend *PortTrend, labels *PortLabeler) *PortTrendReporter {
	return &PortTrendReporter{trend: trend, labels: labels}
}

// TrendEntry holds a single port's trend information.
type TrendEntry struct {
	Port      uint16
	Label     string
	Direction TrendDirection
}

// Entries returns TrendEntry values for each port in ports, sorted by port number.
func (r *PortTrendReporter) Entries(ports []uint16) []TrendEntry {
	sorted := make([]uint16, len(ports))
	copy(sorted, ports)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })

	out := make([]TrendEntry, 0, len(sorted))
	for _, p := range sorted {
		out = append(out, TrendEntry{
			Port:      p,
			Label:     r.labels.Label(p),
			Direction: r.trend.Trend(p),
		})
	}
	return out
}

// Summary returns a human-readable multi-line summary of trends for ports.
func (r *PortTrendReporter) Summary(ports []uint16) string {
	entries := r.Entries(ports)
	if len(entries) == 0 {
		return "no trend data"
	}
	var sb strings.Builder
	for _, e := range entries {
		sb.WriteString(fmt.Sprintf("port %d (%s): %s\n", e.Port, e.Label, e.Direction))
	}
	return strings.TrimRight(sb.String(), "\n")
}
