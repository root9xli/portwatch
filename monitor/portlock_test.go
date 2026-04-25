package monitor

import (
	"testing"
	"time"
)

func newTestPortLock() (*PortLock, *time.Time) {
	now := time.Now()
	pl := NewPortLock()
	pl.clock = func() time.Time { return now }
	return pl, &now
}

func TestPortLockNotLockedInitially(t *testing.T) {
	pl, _ := newTestPortLock()
	if pl.IsLocked(8080) {
		t.Fatal("expected port to not be locked initially")
	}
}

func TestPortLockLockedAfterLock(t *testing.T) {
	pl, _ := newTestPortLock()
	pl.Lock(8080, "admin", time.Minute)
	if !pl.IsLocked(8080) {
		t.Fatal("expected port to be locked after Lock()")
	}
}

func TestPortLockUnlockRemovesLock(t *testing.T) {
	pl, _ := newTestPortLock()
	pl.Lock(8080, "admin", time.Minute)
	pl.Unlock(8080)
	if pl.IsLocked(8080) {
		t.Fatal("expected port to be unlocked after Unlock()")
	}
}

func TestPortLockExpiresAfterDuration(t *testing.T) {
	pl, now := newTestPortLock()
	pl.Lock(9000, "ci", 5*time.Second)
	*now = now.Add(6 * time.Second)
	if pl.IsLocked(9000) {
		t.Fatal("expected lock to have expired")
	}
}

func TestPortLockGetReturnsEntry(t *testing.T) {
	pl, _ := newTestPortLock()
	pl.Lock(443, "security", time.Hour)
	e, ok := pl.Get(443)
	if !ok {
		t.Fatal("expected Get to return entry")
	}
	if e.LockedBy != "security" {
		t.Fatalf("expected LockedBy=security, got %s", e.LockedBy)
	}
}

func TestPortLockGetMissingReturnsFalse(t *testing.T) {
	pl, _ := newTestPortLock()
	_, ok := pl.Get(1234)
	if ok {
		t.Fatal("expected Get to return false for missing port")
	}
}

func TestPortLockIndependentPorts(t *testing.T) {
	pl, _ := newTestPortLock()
	pl.Lock(80, "a", time.Minute)
	if pl.IsLocked(443) {
		t.Fatal("locking port 80 should not affect port 443")
	}
}

func TestPortLockEvictRemovesExpired(t *testing.T) {
	pl, now := newTestPortLock()
	pl.Lock(8080, "a", 2*time.Second)
	pl.Lock(9090, "b", time.Minute)
	*now = now.Add(5 * time.Second)
	pl.Evict()
	if pl.Len() != 1 {
		t.Fatalf("expected 1 active lock after eviction, got %d", pl.Len())
	}
}

func TestPortLockLen(t *testing.T) {
	pl, _ := newTestPortLock()
	pl.Lock(80, "a", time.Minute)
	pl.Lock(443, "b", time.Minute)
	pl.Lock(8080, "c", time.Minute)
	if pl.Len() != 3 {
		t.Fatalf("expected Len=3, got %d", pl.Len())
	}
}
