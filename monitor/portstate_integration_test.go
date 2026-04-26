package monitor

import (
	"testing"
)

func TestPortStateTransitionsFromDiffAdded(t *testing.T) {
	ps := NewPortState(defaultCooldown)

	added := []Message{
		{Port: 8080, Action: "added", Level: "warn"},
		{Port: 443, Action: "added", Level: "info"},
	}

	for _, m := range added {
		changed := ps.Transition(m.Port, m.Action)
		if !changed {
			t.Errorf("expected transition for port %d", m.Port)
		}
	}

	if ps.Len() != 2 {
		t.Errorf("expected 2 tracked ports, got %d", ps.Len())
	}

	e := ps.Get(8080)
	if e == nil || e.Current != "added" {
		t.Errorf("expected port 8080 state=added")
	}
}

func TestPortStateTransitionsFromDiffRemovedAfterAdded(t *testing.T) {
	ps := NewPortState(defaultCooldown)

	ps.Transition(8080, "added")
	changed := ps.Transition(8080, "removed")
	if !changed {
		t.Fatal("expected state change from added to removed")
	}

	e := ps.Get(8080)
	if e.Current != "removed" {
		t.Errorf("expected current=removed, got %s", e.Current)
	}
	if e.Previous != "added" {
		t.Errorf("expected previous=added, got %s", e.Previous)
	}
	if e.ChangeCount != 2 {
		t.Errorf("expected ChangeCount=2, got %d", e.ChangeCount)
	}
}

func TestPortStateNoChangeOnSameActionFromDiff(t *testing.T) {
	ps := NewPortState(defaultCooldown)

	messages := []Message{
		{Port: 9090, Action: "added"},
		{Port: 9090, Action: "added"},
		{Port: 9090, Action: "added"},
	}

	changeCount := 0
	for _, m := range messages {
		if ps.Transition(m.Port, m.Action) {
			changeCount++
		}
	}

	if changeCount != 1 {
		t.Errorf("expected 1 real transition, got %d", changeCount)
	}

	e := ps.Get(9090)
	if e.ChangeCount != 1 {
		t.Errorf("expected ChangeCount=1, got %d", e.ChangeCount)
	}
}
