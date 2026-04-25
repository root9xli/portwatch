package monitor

import (
	"testing"
	"time"
)

func newTestPulse(maxAge time.Duration) *PortPulse {
	p := NewPortPulse(maxAge)
	base := time.Now()
	p.now = func() time.Time { return base }
	return p
}

func TestPortPulseSilentBeforeObserve(t *testing.T) {
	p := newTestPulse(5 * time.Second)
	if !p.Silent(8080) {
		t.Fatal("expected port to be silent before any observation")
	}
}

func TestPortPulseNotSilentAfterObserve(t *testing.T) {
	p := newTestPulse(5 * time.Second)
	p.Observe(8080)
	if p.Silent(8080) {
		t.Fatal("expected port to not be silent immediately after observation")
	}
}

func TestPortPulseSilentAfterMaxAge(t *testing.T) {
	base := time.Now()
	p := NewPortPulse(5 * time.Second)
	p.now = func() time.Time { return base }
	p.Observe(8080)
	// advance time beyond maxAge
	p.now = func() time.Time { return base.Add(6 * time.Second) }
	if !p.Silent(8080) {
		t.Fatal("expected port to be silent after maxAge elapsed")
	}
}

func TestPortPulseIndependentPorts(t *testing.T) {
	p := newTestPulse(5 * time.Second)
	p.Observe(443)
	if p.Silent(443) {
		t.Fatal("expected 443 to be active")
	}
	if !p.Silent(8080) {
		t.Fatal("expected 8080 to be silent")
	}
}

func TestPortPulseEvictRemovesStale(t *testing.T) {
	base := time.Now()
	p := NewPortPulse(5 * time.Second)
	p.now = func() time.Time { return base }
	p.Observe(9000)
	p.Observe(9001)
	p.now = func() time.Time { return base.Add(6 * time.Second) }
	p.Evict()
	if p.Len() != 0 {
		t.Fatalf("expected 0 entries after eviction, got %d", p.Len())
	}
}

func TestPortPulseLastSeen(t *testing.T) {
	base := time.Now()
	p := NewPortPulse(5 * time.Second)
	p.now = func() time.Time { return base }
	p.Observe(3000)
	ts, ok := p.LastSeen(3000)
	if !ok {
		t.Fatal("expected LastSeen to return true for observed port")
	}
	if !ts.Equal(base) {
		t.Fatalf("expected timestamp %v, got %v", base, ts)
	}
}

func TestPortPulseLastSeenMissing(t *testing.T) {
	p := newTestPulse(5 * time.Second)
	_, ok := p.LastSeen(1234)
	if ok {
		t.Fatal("expected LastSeen to return false for unobserved port")
	}
}
