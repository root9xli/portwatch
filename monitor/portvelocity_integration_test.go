package monitor

import (
	"testing"
	"time"
)

func TestPortVelocityDetectedFromRepeatedDiff(t *testing.T) {
	pv := NewPortVelocity(10*time.Second, 0.3)

	msgs := []Message{
		{Port: 7070, Action: "added"},
		{Port: 7070, Action: "added"},
		{Port: 7070, Action: "added"},
		{Port: 7070, Action: "added"},
	}

	now := time.Now()
	var lastEvent *PortVelocityEvent
	for i, m := range msgs {
		if m.Action == "added" {
			ev := pv.Observe(m.Port, now.Add(time.Duration(i)*time.Second))
			if ev != nil {
				lastEvent = ev
			}
		}
	}

	if lastEvent == nil {
		t.Fatal("expected velocity event from repeated diff, got nil")
	}
	if lastEvent.Port != 7070 {
		t.Errorf("expected port 7070, got %d", lastEvent.Port)
	}
}

func TestPortVelocityStablePortNotFlagged(t *testing.T) {
	pv := NewPortVelocity(10*time.Second, 2.0) // high threshold

	now := time.Now()
	// Only two observations spread far apart
	pv.Observe(8181, now)
	ev := pv.Observe(8181, now.Add(9*time.Second))

	if ev != nil {
		t.Errorf("stable port should not trigger velocity event, got rate %f", ev.Rate)
	}
}

func TestPortVelocityEvictClearsAfterDiff(t *testing.T) {
	pv := NewPortVelocity(3*time.Second, 0.5)
	now := time.Now()

	for i := 0; i < 5; i++ {
		pv.Observe(5555, now.Add(time.Duration(i)*100*time.Millisecond))
	}

	// Advance past window and evict
	pv.Evict(now.Add(10 * time.Second))

	// After eviction a fresh observation should not trigger
	ev := pv.Observe(5555, now.Add(10*time.Second))
	if ev != nil {
		t.Errorf("expected nil after eviction, got event with rate %f", ev.Rate)
	}
}
