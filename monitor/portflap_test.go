package monitor

import (
	"testing"
	"time"
)

func newTestFlap(threshold int, window time.Duration) *PortFlap {
	return NewPortFlap(threshold, window)
}

func TestPortFlapBelowThresholdReturnsNil(t *testing.T) {
	f := newTestFlap(3, time.Minute)
	now := time.Now()
	if ev := f.Observe(8080, now); ev != nil {
		t.Fatalf("expected nil, got event")
	}
	if ev := f.Observe(8080, now.Add(time.Second)); ev != nil {
		t.Fatalf("expected nil, got event")
	}
}

func TestPortFlapAtThresholdReturnsEvent(t *testing.T) {
	f := newTestFlap(3, time.Minute)
	now := time.Now()
	f.Observe(9000, now)
	f.Observe(9000, now.Add(time.Second))
	ev := f.Observe(9000, now.Add(2*time.Second))
	if ev == nil {
		t.Fatal("expected flap event, got nil")
	}
	if ev.Port != 9000 {
		t.Errorf("expected port 9000, got %d", ev.Port)
	}
	if ev.Flaps != 3 {
		t.Errorf("expected 3 flaps, got %d", ev.Flaps)
	}
}

func TestPortFlapIndependentPorts(t *testing.T) {
	f := newTestFlap(2, time.Minute)
	now := time.Now()
	f.Observe(80, now)
	ev := f.Observe(443, now.Add(time.Second))
	if ev != nil {
		t.Fatal("port 443 should not flap independently")
	}
}

func TestPortFlapEvictsOldSamples(t *testing.T) {
	f := newTestFlap(2, 5*time.Second)
	now := time.Now()
	f.Observe(3000, now)
	// second observation is outside the window
	ev := f.Observe(3000, now.Add(10*time.Second))
	if ev != nil {
		t.Fatal("old sample should have been evicted, no flap expected")
	}
}

func TestPortFlapForgetResetsPort(t *testing.T) {
	f := newTestFlap(2, time.Minute)
	now := time.Now()
	f.Observe(7000, now)
	f.Forget(7000)
	ev := f.Observe(7000, now.Add(time.Second))
	if ev != nil {
		t.Fatal("expected nil after forget, got event")
	}
}

func TestPortFlapEvictRemovesStale(t *testing.T) {
	f := newTestFlap(2, 5*time.Second)
	now := time.Now()
	f.Observe(5000, now)
	f.Evict(now.Add(10 * time.Second))
	// after evict the counter is gone; a fresh observation should not flap
	ev := f.Observe(5000, now.Add(11*time.Second))
	if ev != nil {
		t.Fatal("expected nil after eviction, got event")
	}
}
