package monitor_test

import (
	"testing"
	"time"
)

// TestPortPriorityWithDiffMessages verifies that PortPriority correctly
// assigns and retrieves priorities for ports observed through a diff cycle.
func TestPortPriorityWithDiffMessages(t *testing.T) {
	pri := NewPortPriority()

	// Simulate a diff producing added ports
	added := makeSnapshot([]int{8080, 9090, 3000})
	removed := makeSnapshot([]int{})
	msgs := Diff(makeSnapshot([]int{}), added)
	_ = removed

	if len(msgs) == 0 {
		t.Fatal("expected diff messages for added ports")
	}

	// Assign priorities based on observed ports
	for _, msg := range msgs {
		switch msg.Port {
		case 8080:
			pri.Set(msg.Port, PriorityHigh)
		case 9090:
			pri.Set(msg.Port, PriorityMedium)
		case 3000:
			pri.Set(msg.Port, PriorityLow)
		}
	}

	if p, ok := pri.Get(8080); !ok || p != PriorityHigh {
		t.Errorf("expected port 8080 to have high priority, got %v (ok=%v)", p, ok)
	}
	if p, ok := pri.Get(9090); !ok || p != PriorityMedium {
		t.Errorf("expected port 9090 to have medium priority, got %v (ok=%v)", p, ok)
	}
	if p, ok := pri.Get(3000); !ok || p != PriorityLow {
		t.Errorf("expected port 3000 to have low priority, got %v (ok=%v)", p, ok)
	}
}

// TestPortPriorityRemovedPortClearedAfterDiff verifies that when a port
// disappears from a diff, its priority entry can be explicitly removed.
func TestPortPriorityRemovedPortClearedAfterDiff(t *testing.T) {
	pri := NewPortPriority()
	pri.Set(4444, PriorityHigh)

	// Port 4444 is now removed in the diff
	prev := makeSnapshot([]int{4444})
	curr := makeSnapshot([]int{})
	msgs := Diff(prev, curr)

	for _, msg := range msgs {
		if msg.Action == "removed" {
			pri.Remove(msg.Port)
		}
	}

	if _, ok := pri.Get(4444); ok {
		t.Error("expected port 4444 priority to be removed after diff removal")
	}
}

// TestPortPriorityHighPriorityPortsReportedFirst verifies that Top returns
// ports ordered by descending priority after a diff cycle populates them.
func TestPortPriorityHighPriorityPortsReportedFirst(t *testing.T) {
	pri := NewPortPriority()

	ports := []int{1111, 2222, 3333}
	msgs := Diff(makeSnapshot([]int{}), makeSnapshot(ports))
	if len(msgs) == 0 {
		t.Fatal("expected diff messages")
	}

	pri.Set(1111, PriorityLow)
	pri.Set(2222, PriorityHigh)
	pri.Set(3333, PriorityMedium)

	top := pri.Top(3)
	if len(top) != 3 {
		t.Fatalf("expected 3 results from Top, got %d", len(top))
	}

	if top[0].Port != 2222 {
		t.Errorf("expected first port to be 2222 (high), got %d", top[0].Port)
	}
	if top[1].Port != 3333 {
		t.Errorf("expected second port to be 3333 (medium), got %d", top[1].Port)
	}
	if top[2].Port != 1111 {
		t.Errorf("expected third port to be 1111 (low), got %d", top[2].Port)
	}
}

// TestPortPriorityLenAfterDiff verifies that Len reflects the number of
// ports assigned a priority following a diff cycle.
func TestPortPriorityLenAfterDiff(t *testing.T) {
	pri := NewPortPriority()

	msgs := Diff(makeSnapshot([]int{}), makeSnapshot([]int{5000, 6000, 7000}))
	for _, msg := range msgs {
		pri.Set(msg.Port, PriorityMedium)
	}

	if pri.Len() != 3 {
		t.Errorf("expected Len() == 3, got %d", pri.Len())
	}

	// Ensure no stale state from time-based eviction within a short window
	time.Sleep(1 * time.Millisecond)
	if pri.Len() != 3 {
		t.Errorf("expected Len() still 3 after brief sleep, got %d", pri.Len())
	}
}
