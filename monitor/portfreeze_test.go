package monitor

import (
	"testing"
	"time"
)

func newTestFreeze() *PortFreeze {
	pf := NewPortFreeze()
	return pf
}

func TestPortFreezeNotFrozenInitially(t *testing.T) {
	pf := newTestFreeze()
	if pf.IsFrozen(8080) {
		t.Fatal("expected port not to be frozen initially")
	}
}

func TestPortFreezeFrozenAfterFreeze(t *testing.T) {
	pf := newTestFreeze()
	pf.Freeze(8080, 5*time.Minute, "maintenance")
	if !pf.IsFrozen(8080) {
		t.Fatal("expected port to be frozen")
	}
}

func TestPortFreezeLiftRemovesFreeze(t *testing.T) {
	pf := newTestFreeze()
	pf.Freeze(8080, 5*time.Minute, "test")
	pf.Lift(8080)
	if pf.IsFrozen(8080) {
		t.Fatal("expected freeze to be lifted")
	}
}

func TestPortFreezeExpiresAfterDuration(t *testing.T) {
	pf := newTestFreeze()
	now := time.Now()
	pf.clock = func() time.Time { return now }
	pf.Freeze(8080, 1*time.Second, "expiry test")
	pf.clock = func() time.Time { return now.Add(2 * time.Second) }
	if pf.IsFrozen(8080) {
		t.Fatal("expected freeze to have expired")
	}
}

func TestPortFreezeGetReturnsEntry(t *testing.T) {
	pf := newTestFreeze()
	pf.Freeze(9090, 10*time.Minute, "deploy")
	e := pf.Get(9090)
	if e == nil {
		t.Fatal("expected non-nil entry")
	}
	if e.Reason != "deploy" {
		t.Fatalf("expected reason 'deploy', got %q", e.Reason)
	}
}

func TestPortFreezeGetMissingReturnsNil(t *testing.T) {
	pf := newTestFreeze()
	if pf.Get(1234) != nil {
		t.Fatal("expected nil for unknown port")
	}
}

func TestPortFreezeEvictRemovesExpired(t *testing.T) {
	pf := newTestFreeze()
	now := time.Now()
	pf.clock = func() time.Time { return now }
	pf.Freeze(80, 1*time.Second, "old")
	pf.Freeze(443, 1*time.Hour, "active")
	pf.clock = func() time.Time { return now.Add(2 * time.Second) }
	pf.Evict()
	if pf.IsFrozen(80) {
		t.Fatal("expected expired entry to be evicted")
	}
	if !pf.IsFrozen(443) {
		t.Fatal("expected active entry to remain")
	}
}

func TestPortFreezeIndependentPorts(t *testing.T) {
	pf := newTestFreeze()
	pf.Freeze(8080, 5*time.Minute, "a")
	if pf.IsFrozen(9090) {
		t.Fatal("expected different port to be unaffected")
	}
}
