package monitor

import (
	"testing"
	"time"
)

// TestDeadPortWithDiff simulates a port disappearing and reappearing via Diff.
func TestDeadPortWithDiff(t *testing.T) {
	now := time.Now()
	dp := NewDeadPort(10 * time.Second)
	dp.clock = func() time.Time { return now }

	snap1 := makeSnapshot([]uint16{80, 443, 8080})
	snap2 := makeSnapshot([]uint16{80, 443})
	snap3 := makeSnapshot([]uint16{80, 443, 8080})

	diff1 := Diff(snap1, snap2)
	for _, msg := range diff1 {
		if msg.Action == "removed" {
			dp.MarkGone(msg.Port)
		}
	}

	diff2 := Diff(snap2, snap3)
	recycled := []uint16{}
	for _, msg := range diff2 {
		if msg.Action == "added" && dp.IsRecycled(msg.Port) {
			recycled = append(recycled, msg.Port)
		}
	}

	if len(recycled) != 1 || recycled[0] != 8080 {
		t.Fatalf("expected port 8080 flagged as recycled, got %v", recycled)
	}
}

// TestDeadPortNoFalsePositiveAfterWindow ensures recycled detection expires.
func TestDeadPortNoFalsePositiveAfterWindow(t *testing.T) {
	now := time.Now()
	dp := NewDeadPort(5 * time.Second)
	dp.clock = func() time.Time { return now }

	snap1 := makeSnapshot([]uint16{8080})
	snap2 := makeSnapshot([]uint16{})
	snap3 := makeSnapshot([]uint16{8080})

	diff1 := Diff(snap1, snap2)
	for _, msg := range diff1 {
		if msg.Action == "removed" {
			dp.MarkGone(msg.Port)
		}
	}

	dp.clock = func() time.Time { return now.Add(20 * time.Second) }

	diff2 := Diff(snap2, snap3)
	for _, msg := range diff2 {
		if msg.Action == "added" && dp.IsRecycled(msg.Port) {
			t.Fatalf("expected no recycled detection after window expiry")
		}
	}
}
