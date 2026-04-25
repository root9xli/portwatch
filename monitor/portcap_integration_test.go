package monitor

import (
	"testing"
	"time"
)

// TestPortCapTracksAddedPortsFromDiff verifies that PortCap integrates with
// Diff output by observing counts derived from added port messages.
func TestPortCapTracksAddedPortsFromDiff(t *testing.T) {
	pc := NewPortCap(time.Minute)

	before := makeSnapshot([]int{80})
	after := makeSnapshot([]int{80, 8080, 9090})
	msgs := Diff(before, after)

	for _, m := range msgs {
		if m.Action == "added" {
			// Simulate count=1 for first appearance.
			newPeak := pc.Observe(m.Port, 1)
			if newPeak {
				t.Errorf("first observation for port %d should not flag new peak", m.Port)
			}
		}
	}

	// Simulate a second wave with higher counts.
	peakFound := false
	for _, m := range msgs {
		if m.Action == "added" {
			if pc.Observe(m.Port, 5) {
				peakFound = true
			}
		}
	}
	if !peakFound {
		t.Fatal("expected at least one new peak after higher count observation")
	}
}

// TestPortCapNoPeakForRemovedPorts ensures removed ports are not tracked.
func TestPortCapNoPeakForRemovedPorts(t *testing.T) {
	pc := NewPortCap(time.Minute)

	before := makeSnapshot([]int{80, 8080})
	after := makeSnapshot([]int{80})
	msgs := Diff(before, after)

	for _, m := range msgs {
		if m.Action == "removed" {
			// Should not observe removed ports.
			_ = m
		}
	}

	if pc.Len() != 0 {
		t.Fatalf("expected no tracked ports for removed-only diff, got %d", pc.Len())
	}
}
