package monitor

import (
	"strings"
	"testing"
	"time"
)

func newTestPingReporter(t *testing.T) (*PortPingReporter, int, func()) {
	t.Helper()
	port, stop := startTestListener(t)
	pp := NewPortPing(time.Second)
	return NewPortPingReporter(pp), port, stop
}

func TestPortPingReporterEmptySummary(t *testing.T) {
	pp := NewPortPing(time.Second)
	r := NewPortPingReporter(pp)
	if !strings.Contains(r.Summary(), "no probes") {
		t.Error("expected 'no probes' in empty summary")
	}
}

func TestPortPingReporterSummaryContainsPort(t *testing.T) {
	r, port, stop := newTestPingReporter(t)
	defer stop()
	r.Run([]int{port})
	summary := r.Summary()
	if !strings.Contains(summary, "UP") {
		t.Errorf("expected UP in summary, got: %s", summary)
	}
}

func TestPortPingReporterUnreachable(t *testing.T) {
	pp := NewPortPing(100 * time.Millisecond)
	r := NewPortPingReporter(pp)
	r.Run([]int{1})
	u := r.Unreachable()
	if len(u) == 0 {
		t.Skip("port 1 reachable; skipping unreachable test")
	}
	if u[0] != 1 {
		t.Errorf("expected port 1 in unreachable list")
	}
}

func TestPortPingReporterAvgLatency(t *testing.T) {
	r, port, stop := newTestPingReporter(t)
	defer stop()
	r.Run([]int{port})
	if r.AvgLatency() <= 0 {
		t.Error("expected positive avg latency after probing reachable port")
	}
}

func TestPortPingReporterAvgLatencyNoReachable(t *testing.T) {
	pp := NewPortPing(100 * time.Millisecond)
	r := NewPortPingReporter(pp)
	r.Run([]int{1})
	if r.AvgLatency() != 0 && len(r.Unreachable()) > 0 {
		// only assert zero when we know port 1 is unreachable
		if r.AvgLatency() != 0 {
			t.Error("expected zero avg latency when no reachable ports")
		}
	}
}
