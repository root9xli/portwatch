package monitor

import (
	"testing"
	"time"
)

func newTestVelocity() *PortVelocity {
	return NewPortVelocity(10*time.Second, 0.5)
}

func TestPortVelocityBelowThresholdReturnsNil(t *testing.T) {
	pv := newTestVelocity()
	now := time.Now()
	// Single observation — rate too low
	ev := pv.Observe(8080, now)
	if ev != nil {
		t.Fatalf("expected nil, got event with rate %f", ev.Rate)
	}
}

func TestPortVelocityAtThresholdReturnsEvent(t *testing.T) {
	pv := NewPortVelocity(10*time.Second, 0.3)
	now := time.Now()
	for i := 0; i < 4; i++ {
		pv.Observe(9090, now.Add(time.Duration(i)*time.Second))
	}
	ev := pv.Observe(9090, now.Add(4*time.Second))
	if ev == nil {
		t.Fatal("expected velocity event, got nil")
	}
	if ev.Port != 9090 {
		t.Errorf("expected port 9090, got %d", ev.Port)
	}
	if ev.Rate < ev.Threshold {
		t.Errorf("rate %f should be >= threshold %f", ev.Rate, ev.Threshold)
	}
}

func TestPortVelocityIndependentPorts(t *testing.T) {
	pv := NewPortVelocity(10*time.Second, 0.3)
	now := time.Now()
	for i := 0; i < 5; i++ {
		pv.Observe(1111, now.Add(time.Duration(i)*time.Second))
	}
	// Port 2222 has only one observation — should not trigger
	ev := pv.Observe(2222, now)
	if ev != nil {
		t.Errorf("expected nil for independent port, got event")
	}
}

func TestPortVelocityEvictsOldSamples(t *testing.T) {
	pv := NewPortVelocity(5*time.Second, 0.3)
	base := time.Now()
	for i := 0; i < 5; i++ {
		pv.Observe(3000, base.Add(time.Duration(i)*time.Second))
	}
	// Evict at base+20s — all samples should be gone
	pv.Evict(base.Add(20 * time.Second))
	pv.mu.Lock()
	_, exists := pv.samples[3000]
	pv.mu.Unlock()
	if exists {
		t.Error("expected samples for port 3000 to be evicted")
	}
}

func TestPortVelocityEventFields(t *testing.T) {
	pv := NewPortVelocity(10*time.Second, 0.2)
	now := time.Now()
	for i := 0; i < 3; i++ {
		pv.Observe(4444, now.Add(time.Duration(i)*time.Second))
	}
	ev := pv.Observe(4444, now.Add(3*time.Second))
	if ev == nil {
		t.Fatal("expected event")
	}
	if ev.Threshold != 0.2 {
		t.Errorf("expected threshold 0.2, got %f", ev.Threshold)
	}
	if ev.At.IsZero() {
		t.Error("expected non-zero At timestamp")
	}
}
