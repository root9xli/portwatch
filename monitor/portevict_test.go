package monitor

import (
	"testing"
	"time"
)

func newTestEvict(maxAge time.Duration) *PortEvict {
	pe := NewPortEvict(maxAge)
	return pe
}

func TestPortEvictNotPresentInitially(t *testing.T) {
	pe := newTestEvict(time.Minute)
	_, ok := pe.Get(8080)
	if ok {
		t.Fatal("expected port to be absent initially")
	}
}

func TestPortEvictRecordStored(t *testing.T) {
	pe := newTestEvict(time.Minute)
	now := time.Now()
	pe.Evict(8080, now)
	r, ok := pe.Get(8080)
	if !ok {
		t.Fatal("expected eviction record to be present")
	}
	if r.Port != 8080 {
		t.Errorf("expected port 8080, got %d", r.Port)
	}
	if r.Count != 1 {
		t.Errorf("expected count 1, got %d", r.Count)
	}
}

func TestPortEvictCountIncrements(t *testing.T) {
	pe := newTestEvict(time.Minute)
	now := time.Now()
	pe.Evict(9090, now)
	pe.Evict(9090, now.Add(time.Second))
	r, ok := pe.Get(9090)
	if !ok {
		t.Fatal("expected record present")
	}
	if r.Count != 2 {
		t.Errorf("expected count 2, got %d", r.Count)
	}
}

func TestPortEvictExpiresAfterMaxAge(t *testing.T) {
	pe := newTestEvict(50 * time.Millisecond)
	fixed := time.Now()
	pe.now = func() time.Time { return fixed }
	pe.Evict(1234, fixed)

	// advance clock past maxAge
	pe.now = func() time.Time { return fixed.Add(100 * time.Millisecond) }
	_, ok := pe.Get(1234)
	if ok {
		t.Fatal("expected record to have expired")
	}
}

func TestPortEvictAllReturnsActive(t *testing.T) {
	pe := newTestEvict(time.Minute)
	now := time.Now()
	pe.Evict(80, now)
	pe.Evict(443, now)
	all := pe.All()
	if len(all) != 2 {
		t.Errorf("expected 2 records, got %d", len(all))
	}
}

func TestPortEvictLenMatchesActive(t *testing.T) {
	pe := newTestEvict(time.Minute)
	now := time.Now()
	pe.Evict(22, now)
	pe.Evict(25, now)
	pe.Evict(53, now)
	if pe.Len() != 3 {
		t.Errorf("expected Len 3, got %d", pe.Len())
	}
}

func TestPortEvictIndependentPorts(t *testing.T) {
	pe := newTestEvict(time.Minute)
	now := time.Now()
	pe.Evict(3000, now)
	_, ok := pe.Get(4000)
	if ok {
		t.Fatal("port 4000 should not be present")
	}
}
