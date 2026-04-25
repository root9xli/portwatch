package monitor

import (
	"testing"
	"time"
)

// TestPortPulseObservedFromDiff verifies that Observe is called for each added
// port produced by a Diff, and that those ports are no longer silent.
func TestPortPulseObservedFromDiff(t *testing.T) {
	pulse := NewPortPulse(10 * time.Second)

	old := makeSnapshot([]int{80})
	new_ := makeSnapshot([]int{80, 443, 8080})
	result := Diff(old, new_)

	for _, msg := range result.Added {
		pulse.Observe(msg.Port)
	}

	for _, port := range []int{443, 8080} {
		if pulse.Silent(port) {
			t.Errorf("expected port %d to be active after Observe", port)
		}
	}
	// Port 80 was not in Added, so pulse has no record of it.
	if !pulse.Silent(80) {
		t.Error("expected port 80 to be silent (never observed via diff)")
	}
}

// TestPortPulseSilentAfterExpiry verifies that a port observed during a diff
// becomes silent once the maxAge window passes.
func TestPortPulseSilentAfterExpiry(t *testing.T) {
	base := time.Now()
	pulse := NewPortPulse(5 * time.Second)
	pulse.now = func() time.Time { return base }

	old := makeSnapshot([]int{})
	new_ := makeSnapshot([]int{9090})
	result := Diff(old, new_)

	for _, msg := range result.Added {
		pulse.Observe(msg.Port)
	}

	if pulse.Silent(9090) {
		t.Fatal("expected 9090 to be active immediately after observation")
	}

	// Advance time past maxAge
	pulse.now = func() time.Time { return base.Add(6 * time.Second) }

	if !pulse.Silent(9090) {
		t.Fatal("expected 9090 to be silent after maxAge elapsed")
	}
}
