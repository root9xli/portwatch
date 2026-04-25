package monitor

import (
	"testing"
	"time"
)

func newTestExpiry() *PortExpiry {
	pe := NewPortExpiry()
	return pe
}

func TestPortExpiryNotOverdueIfNotRegistered(t *testing.T) {
	pe := newTestExpiry()
	if pe.IsOverdue(8080) {
		t.Fatal("expected port not registered to not be overdue")
	}
}

func TestPortExpiryNotOverdueBeforeDeadline(t *testing.T) {
	pe := newTestExpiry()
	pe.Register(8080, time.Now().Add(10*time.Minute))
	if pe.IsOverdue(8080) {
		t.Fatal("expected port to not be overdue before deadline")
	}
}

func TestPortExpiryOverdueAfterDeadline(t *testing.T) {
	pe := newTestExpiry()
	past := time.Now().Add(-1 * time.Second)
	pe.Register(8080, past)
	if !pe.IsOverdue(8080) {
		t.Fatal("expected port to be overdue after deadline")
	}
}

func TestPortExpiryForgetRemovesEntry(t *testing.T) {
	pe := newTestExpiry()
	past := time.Now().Add(-1 * time.Second)
	pe.Register(8080, past)
	pe.Forget(8080)
	if pe.IsOverdue(8080) {
		t.Fatal("expected forgotten port to not be overdue")
	}
	if pe.Len() != 0 {
		t.Fatalf("expected len 0, got %d", pe.Len())
	}
}

func TestPortExpiryOverdueReturnsMultiple(t *testing.T) {
	pe := newTestExpiry()
	past := time.Now().Add(-1 * time.Second)
	future := time.Now().Add(10 * time.Minute)
	pe.Register(8080, past)
	pe.Register(9090, past)
	pe.Register(443, future)
	overdue := pe.Overdue()
	if len(overdue) != 2 {
		t.Fatalf("expected 2 overdue ports, got %d", len(overdue))
	}
}

func TestPortExpiryEvictRemovesOverdue(t *testing.T) {
	pe := newTestExpiry()
	past := time.Now().Add(-1 * time.Second)
	future := time.Now().Add(10 * time.Minute)
	pe.Register(8080, past)
	pe.Register(443, future)
	evicted := pe.Evict()
	if len(evicted) != 1 || evicted[0] != 8080 {
		t.Fatalf("expected [8080] evicted, got %v", evicted)
	}
	if pe.Len() != 1 {
		t.Fatalf("expected 1 remaining, got %d", pe.Len())
	}
}

func TestPortExpiryLenReflectsRegistrations(t *testing.T) {
	pe := newTestExpiry()
	if pe.Len() != 0 {
		t.Fatalf("expected 0, got %d", pe.Len())
	}
	pe.Register(8080, time.Now().Add(time.Minute))
	pe.Register(9090, time.Now().Add(time.Minute))
	if pe.Len() != 2 {
		t.Fatalf("expected 2, got %d", pe.Len())
	}
}
