package monitor

import (
	"testing"
	"time"
)

// TestPortEvictPopulatedFromDiff verifies that removed ports from a Diff are
// recorded in PortEvict.
func TestPortEvictPopulatedFromDiff(t *testing.T) {
	pe := NewPortEvict(time.Minute)

	before := makeSnapshot([]int{80, 443, 8080})
	after := makeSnapshot([]int{80, 443})

	result := Diff(before, after)
	if len(result.Removed) == 0 {
		t.Fatal("expected at least one removed port")
	}

	now := time.Now()
	for _, msg := range result.Removed {
		pe.Evict(msg.Port, now)
	}

	r, ok := pe.Get(8080)
	if !ok {
		t.Fatal("expected port 8080 to be recorded as evicted")
	}
	if r.Port != 8080 {
		t.Errorf("expected port 8080, got %d", r.Port)
	}
}

// TestPortEvictNoFalsePositiveForAdded verifies that added ports are not
// recorded as evicted.
func TestPortEvictNoFalsePositiveForAdded(t *testing.T) {
	pe := NewPortEvict(time.Minute)

	before := makeSnapshot([]int{80})
	after := makeSnapshot([]int{80, 9000})

	result := Diff(before, after)
	now := time.Now()
	for _, msg := range result.Removed {
		pe.Evict(msg.Port, now)
	}

	_, ok := pe.Get(9000)
	if ok {
		t.Fatal("added port 9000 should not appear in eviction records")
	}
	if pe.Len() != 0 {
		t.Errorf("expected 0 eviction records, got %d", pe.Len())
	}
}
