package monitor

import (
	"testing"
	"time"
)

func newTestSpike(threshold int) *PortSpike {
	return NewPortSpike(200*time.Millisecond, threshold)
}

func TestPortSpikeBelowThresholdReturnsNil(t *testing.T) {
	ps := newTestSpike(3)
	if ev := ps.Observe(8080); ev != nil {
		t.Fatalf("expected nil, got event for count=1")
	}
	if ev := ps.Observe(8080); ev != nil {
		t.Fatalf("expected nil, got event for count=2")
	}
}

func TestPortSpikeAtThresholdReturnsEvent(t *testing.T) {
	ps := newTestSpike(3)
	ps.Observe(9000)
	ps.Observe(9000)
	ev := ps.Observe(9000)
	if ev == nil {
		t.Fatal("expected spike event at threshold")
	}
	if ev.Port != 9000 {
		t.Errorf("expected port 9000, got %d", ev.Port)
	}
	if ev.Count < 3 {
		t.Errorf("expected count >= 3, got %d", ev.Count)
	}
}

func TestPortSpikeIndependentPorts(t *testing.T) {
	ps := newTestSpike(2)
	ps.Observe(1111)
	ps.Observe(1111)
	if ev := ps.Observe(2222); ev != nil {
		t.Fatal("port 2222 should not spike from port 1111 observations")
	}
}

func TestPortSpikeEvictsOldSamples(t *testing.T) {
	ps := NewPortSpike(50*time.Millisecond, 3)
	ps.Observe(7070)
	ps.Observe(7070)
	time.Sleep(80 * time.Millisecond)
	// old samples evicted; only 1 new observation — should not spike
	if ev := ps.Observe(7070); ev != nil {
		t.Fatal("expected nil after window expiry")
	}
}

func TestPortSpikeResetClearsPort(t *testing.T) {
	ps := newTestSpike(2)
	ps.Observe(5050)
	ps.Observe(5050)
	ps.Reset(5050)
	if ps.Len() != 0 {
		t.Errorf("expected 0 tracked ports after reset, got %d", ps.Len())
	}
}

func TestPortSpikeLenTracksMultiplePorts(t *testing.T) {
	ps := newTestSpike(5)
	ps.Observe(1)
	ps.Observe(2)
	ps.Observe(3)
	if ps.Len() != 3 {
		t.Errorf("expected 3, got %d", ps.Len())
	}
}
