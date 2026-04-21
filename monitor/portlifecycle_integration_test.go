package monitor

import (
	"testing"
	"time"
)

// TestPortLifecycleTracksAddedPorts verifies that ports appearing in a Diff
// are observed and tracked by the lifecycle manager.
func TestPortLifecycleTracksAddedPorts(t *testing.T) {
	pl := NewPortLifecycle()

	before := makeSnapshot([]int{80})
	after := makeSnapshot([]int{80, 443, 8080})
	result := Diff(before, after)

	for _, msg := range result.Added {
		pl.Observe(msg.Port)
	}

	for _, port := range []int{443, 8080} {
		entry, ok := pl.Get(port)
		if !ok {
			t.Errorf("expected lifecycle entry for port %d", port)
			continue
		}
		if entry.SeenCount != 1 {
			t.Errorf("port %d: expected SeenCount=1, got %d", port, entry.SeenCount)
		}
	}

	// Port 80 was not added, so it should not be tracked
	_, ok := pl.Get(80)
	if ok {
		t.Error("port 80 should not be in lifecycle (it was not added)")
	}
}

// TestPortLifecycleForgetsRemovedPorts verifies that removed ports are
// cleaned up from the lifecycle tracker.
func TestPortLifecycleForgetsRemovedPorts(t *testing.T) {
	pl := NewPortLifecycle()

	// Seed lifecycle with some ports
	pl.Observe(22)
	pl.Observe(80)
	time.Sleep(5 * time.Millisecond)

	before := makeSnapshot([]int{22, 80})
	after := makeSnapshot([]int{80})
	result := Diff(before, after)

	for _, msg := range result.Removed {
		pl.Forget(msg.Port)
	}

	_, ok := pl.Get(22)
	if ok {
		t.Error("expected port 22 to be forgotten after removal")
	}

	_, ok = pl.Get(80)
	if !ok {
		t.Error("expected port 80 to still be tracked")
	}
}
