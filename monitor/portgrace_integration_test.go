package monitor

import (
	"testing"
	"time"
)

// TestPortGraceFiltersNewPortsDuringWindow verifies that messages for newly
// added ports are suppressed while the grace period is active.
func TestPortGraceFiltersNewPortsDuringWindow(t *testing.T) {
	base := time.Now()
	g := NewPortGrace(time.Minute)
	g.now = func() time.Time { return base }

	old := makeSnapshot([]int{80})
	new_ := makeSnapshot([]int{80, 9090})
	msgs := Diff(old, new_)

	if len(msgs) != 1 {
		t.Fatalf("expected 1 diff message, got %d", len(msgs))
	}

	for _, m := range msgs {
		if m.Action == "added" {
			g.Enter(m.Port)
		}
	}

	var filtered []Message
	for _, m := range msgs {
		if !g.InGrace(m.Port) {
			filtered = append(filtered, m)
		}
	}

	if len(filtered) != 0 {
		t.Fatalf("expected all added ports suppressed during grace, got %d messages", len(filtered))
	}
}

// TestPortGraceAllowsPortAfterExpiry verifies that once the grace window
// expires the port is no longer suppressed.
func TestPortGraceAllowsPortAfterExpiry(t *testing.T) {
	base := time.Now()
	g := NewPortGrace(50 * time.Millisecond)
	g.now = func() time.Time { return base }

	g.Enter(7777)

	// advance past grace window
	g.now = func() time.Time { return base.Add(200 * time.Millisecond) }

	if g.InGrace(7777) {
		t.Fatal("expected port 7777 to be out of grace after window expiry")
	}
}
