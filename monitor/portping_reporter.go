package monitor

import (
	"fmt"
	"strings"
	"time"
)

// PortPingReporter summarises recent probe results.
type PortPingReporter struct {
	ping    *PortPing
	results []PingResult
}

// NewPortPingReporter creates a reporter backed by the given PortPing.
func NewPortPingReporter(pp *PortPing) *PortPingReporter {
	return &PortPingReporter{ping: pp}
}

// Run probes the given ports and stores the results.
func (r *PortPingReporter) Run(ports []int) {
	r.results = r.ping.ProbeAll(ports)
}

// Summary returns a human-readable report of the last probe run.
func (r *PortPingReporter) Summary() string {
	if len(r.results) == 0 {
		return "portping: no probes run"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("portping: %d ports probed\n", len(r.results)))
	for _, res := range r.results {
		status := "UP"
		detail := fmt.Sprintf("latency=%s", fmtDuration(res.Latency))
		if !res.Reachable {
			status = "DOWN"
			detail = "unreachable"
		}
		sb.WriteString(fmt.Sprintf("  port=%-6d status=%-4s %s\n", res.Port, status, detail))
	}
	return strings.TrimRight(sb.String(), "\n")
}

// Unreachable returns the ports from the last run that were not reachable.
func (r *PortPingReporter) Unreachable() []int {
	var out []int
	for _, res := range r.results {
		if !res.Reachable {
			out = append(out, res.Port)
		}
	}
	return out
}

// AvgLatency returns the mean latency across reachable ports.
// Returns 0 if no reachable ports were probed.
func (r *PortPingReporter) AvgLatency() time.Duration {
	var total time.Duration
	var count int
	for _, res := range r.results {
		if res.Reachable {
			total += res.Latency
			count++
		}
	}
	if count == 0 {
		return 0
	}
	return total / time.Duration(count)
}
