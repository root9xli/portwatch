package monitor

import (
	"testing"
	"time"
)

// TestPortTrendRecordedFromDiff verifies that recording diff counts drives trend detection.
func TestPortTrendRecordedFromDiff(t *testing.T) {
	pt, now := newTestTrend(time.Minute)

	// Simulate three diff cycles with increasing new-port counts.
	for cycle, count := range []int{1, 1, 4, 5} {
		_ = cycle
		pt.Record(8080, count)
		*now = now.Add(5 * time.Second)
	}

	if got := pt.Trend(8080); got != TrendUp {
		t.Errorf("expected TrendUp after increasing counts, got %s", got)
	}
}

// TestPortTrendStablePortNotFlagged verifies a port with constant counts stays stable.
func TestPortTrendStablePortNotFlagged(t *testing.T) {
	pt, now := newTestTrend(time.Minute)

	for i := 0; i < 6; i++ {
		pt.Record(443, 2)
		*now = now.Add(4 * time.Second)
	}

	if got := pt.Trend(443); got != TrendStable {
		t.Errorf("expected TrendStable for constant counts, got %s", got)
	}
}

// TestPortTrendReporterWithMultiplePorts exercises the reporter over a realistic set of ports.
func TestPortTrendReporterWithMultiplePorts(t *testing.T) {
	pt, now := newTestTrend(time.Minute)
	labels := NewPortLabeler(nil)
	rep := NewPortTrendReporter(pt, labels)

	ports := []uint16{22, 80, 443}
	for _, p := range ports {
		pt.Record(p, 1)
	}
	*now = now.Add(10 * time.Second)
	for _, p := range ports {
		pt.Record(p, 1)
	}

	entries := rep.Entries(ports)
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	for _, e := range entries {
		if e.Direction != TrendStable {
			t.Errorf("port %d: expected stable, got %s", e.Port, e.Direction)
		}
	}
}
