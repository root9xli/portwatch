package monitor

import (
	"testing"
	"time"
)

func newTestBurst(threshold int) *PortBurst {
	return NewPortBurst(10*time.Second, threshold)
}

func TestPortBurstBelowThresholdReturnsNil(t *testing.T) {
	pb := newTestBurst(3)
	now := time.Now()
	if ev := pb.Observe(8080, now); ev != nil {
		t.Fatalf("expected nil, got event with count %d", ev.Count)
	}
	if ev := pb.Observe(8080, now.Add(time.Second)); ev != nil {
		t.Fatalf("expected nil on second observe, got event")
	}
}

func TestPortBurstAtThresholdReturnsEvent(t *testing.T) {
	pb := newTestBurst(3)
	now := time.Now()
	pb.Observe(8080, now)
	pb.Observe(8080, now.Add(time.Second))
	ev := pb.Observe(8080, now.Add(2*time.Second))
	if ev == nil {
		t.Fatal("expected burst event at threshold")
	}
	if ev.Port != 8080 {
		t.Errorf("expected port 8080, got %d", ev.Port)
	}
	if ev.Count != 3 {
		t.Errorf("expected count 3, got %d", ev.Count)
	}
}

func TestPortBurstIndependentPorts(t *testing.T) {
	pb := newTestBurst(2)
	now := time.Now()
	pb.Observe(9090, now)
	ev := pb.Observe(8080, now.Add(time.Second))
	if ev != nil {
		t.Errorf("port 8080 should not burst from port 9090 observations")
	}
}

func TestPortBurstEvictsOldSamples(t *testing.T) {
	pb := NewPortBurst(5*time.Second, 2)
	now := time.Now()
	// first observation is old
	pb.Observe(8080, now.Add(-10*time.Second))
	// second observation is fresh — should not burst because old one evicted
	ev := pb.Observe(8080, now)
	if ev != nil {
		t.Errorf("expected no burst after stale entry evicted, got count %d", ev.Count)
	}
}

func TestPortBurstCountReflectsWindow(t *testing.T) {
	pb := NewPortBurst(5*time.Second, 10)
	now := time.Now()
	pb.Observe(443, now.Add(-6*time.Second)) // outside window
	pb.Observe(443, now.Add(-2*time.Second))
	pb.Observe(443, now)
	count := pb.Count(443, now)
	if count != 2 {
		t.Errorf("expected count 2 (one stale evicted), got %d", count)
	}
}

func TestPortBurstEvictRemovesEmptyBuckets(t *testing.T) {
	pb := NewPortBurst(5*time.Second, 10)
	now := time.Now()
	pb.Observe(1234, now.Add(-10*time.Second))
	pb.Evict(now)
	if count := pb.Count(1234, now); count != 0 {
		t.Errorf("expected 0 after evict, got %d", count)
	}
}
