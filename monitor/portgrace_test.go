package monitor

import (
	"testing"
	"time"
)

func newTestGrace(window time.Duration) *PortGrace {
	g := NewPortGrace(window)
	return g
}

func TestPortGraceNotInGraceInitially(t *testing.T) {
	g := newTestGrace(time.Second)
	if g.InGrace(8080) {
		t.Fatal("expected port 8080 not to be in grace initially")
	}
}

func TestPortGraceInGraceAfterEnter(t *testing.T) {
	g := newTestGrace(time.Second)
	g.Enter(8080)
	if !g.InGrace(8080) {
		t.Fatal("expected port 8080 to be in grace after Enter")
	}
}

func TestPortGraceExpiresAfterWindow(t *testing.T) {
	g := newTestGrace(50 * time.Millisecond)
	base := time.Now()
	g.now = func() time.Time { return base }
	g.Enter(9000)
	g.now = func() time.Time { return base.Add(100 * time.Millisecond) }
	if g.InGrace(9000) {
		t.Fatal("expected grace to have expired")
	}
}

func TestPortGraceClearRemovesPort(t *testing.T) {
	g := newTestGrace(time.Minute)
	g.Enter(443)
	g.Clear(443)
	if g.InGrace(443) {
		t.Fatal("expected port 443 to be cleared")
	}
}

func TestPortGraceIndependentPorts(t *testing.T) {
	g := newTestGrace(time.Minute)
	g.Enter(80)
	if g.InGrace(443) {
		t.Fatal("port 443 should not be in grace")
	}
	if !g.InGrace(80) {
		t.Fatal("port 80 should be in grace")
	}
}

func TestPortGraceEvictRemovesExpired(t *testing.T) {
	g := newTestGrace(50 * time.Millisecond)
	base := time.Now()
	g.now = func() time.Time { return base }
	g.Enter(1234)
	g.Enter(5678)
	g.now = func() time.Time { return base.Add(100 * time.Millisecond) }
	g.Evict()
	if g.Len() != 0 {
		t.Fatalf("expected 0 entries after evict, got %d", g.Len())
	}
}

func TestPortGraceLenReflectsEntries(t *testing.T) {
	g := newTestGrace(time.Minute)
	if g.Len() != 0 {
		t.Fatal("expected empty initially")
	}
	g.Enter(80)
	g.Enter(443)
	if g.Len() != 2 {
		t.Fatalf("expected 2, got %d", g.Len())
	}
}
