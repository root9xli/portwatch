package monitor

import (
	"testing"
	"time"
)

func newTestPortState() (*PortState, *time.Time) {
	now := time.Now()
	ps := NewPortState(5 * time.Minute)
	ps.now = func() time.Time { return now }
	return ps, &now
}

func TestPortStateTransitionFirstTime(t *testing.T) {
	ps, _ := newTestPortState()
	changed := ps.Transition(8080, "open")
	if !changed {
		t.Fatal("expected changed=true on first transition")
	}
	e := ps.Get(8080)
	if e == nil {
		t.Fatal("expected entry")
	}
	if e.Current != "open" {
		t.Errorf("expected current=open, got %s", e.Current)
	}
	if e.Previous != "" {
		t.Errorf("expected empty previous, got %s", e.Previous)
	}
}

func TestPortStateTransitionSameStateNoChange(t *testing.T) {
	ps, _ := newTestPortState()
	ps.Transition(8080, "open")
	changed := ps.Transition(8080, "open")
	if changed {
		t.Fatal("expected changed=false for same state")
	}
	e := ps.Get(8080)
	if e.ChangeCount != 1 {
		t.Errorf("expected ChangeCount=1, got %d", e.ChangeCount)
	}
}

func TestPortStateTransitionUpdatesFields(t *testing.T) {
	ps, _ := newTestPortState()
	ps.Transition(443, "open")
	ps.Transition(443, "closed")
	e := ps.Get(443)
	if e.Current != "closed" {
		t.Errorf("expected current=closed, got %s", e.Current)
	}
	if e.Previous != "open" {
		t.Errorf("expected previous=open, got %s", e.Previous)
	}
	if e.ChangeCount != 2 {
		t.Errorf("expected ChangeCount=2, got %d", e.ChangeCount)
	}
}

func TestPortStateGetMissingReturnsNil(t *testing.T) {
	ps, _ := newTestPortState()
	if ps.Get(9999) != nil {
		t.Fatal("expected nil for unknown port")
	}
}

func TestPortStateIndependentPorts(t *testing.T) {
	ps, _ := newTestPortState()
	ps.Transition(80, "open")
	ps.Transition(443, "closed")
	if ps.Len() != 2 {
		t.Errorf("expected 2 entries, got %d", ps.Len())
	}
	if ps.Get(80).Current != "open" {
		t.Error("port 80 state mismatch")
	}
	if ps.Get(443).Current != "closed" {
		t.Error("port 443 state mismatch")
	}
}

func TestPortStateEvictsOldEntries(t *testing.T) {
	now := time.Now()
	ps := NewPortState(1 * time.Minute)
	ps.now = func() time.Time { return now }
	ps.Transition(8080, "open")
	if ps.Len() != 1 {
		t.Fatal("expected 1 entry before eviction")
	}
	now = now.Add(2 * time.Minute)
	// trigger eviction via Len
	ps.mu.Lock()
	ps.evict()
	ps.mu.Unlock()
	if ps.Len() != 0 {
		t.Errorf("expected 0 entries after eviction, got %d", ps.Len())
	}
}
