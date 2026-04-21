package monitor

import (
	"strings"
	"testing"
	"time"
)

func TestPortPingIntegrationReachablePort(t *testing.T) {
	port, stop := startTestListener(t)
	defer stop()

	pp := NewPortPing(time.Second)
	r := NewPortPingReporter(pp)
	r.Run([]int{port})

	if len(r.Unreachable()) != 0 {
		t.Errorf("expected no unreachable ports, got: %v", r.Unreachable())
	}
	if r.AvgLatency() <= 0 {
		t.Error("expected positive latency")
	}
	summary := r.Summary()
	if !strings.Contains(summary, "1 ports probed") {
		t.Errorf("unexpected summary: %s", summary)
	}
}

func TestPortPingIntegrationMixedPorts(t *testing.T) {
	port, stop := startTestListener(t)
	defer stop()

	pp := NewPortPing(100 * time.Millisecond)
	r := NewPortPingReporter(pp)
	r.Run([]int{port, 1})

	summary := r.Summary()
	if !strings.Contains(summary, "2 ports probed") {
		t.Errorf("unexpected summary: %s", summary)
	}
	// The live port should show UP
	if !strings.Contains(summary, "UP") {
		t.Errorf("expected UP entry in summary: %s", summary)
	}
}
