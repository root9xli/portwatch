package monitor

import (
	"testing"
	"time"
)

func newTestWindow() *PortWindow {
	pw := NewPortWindow()
	return pw
}

func TestPortWindowNotInWindowInitially(t *testing.T) {
	pw := newTestWindow()
	if pw.InWindow(8080) {
		t.Fatal("expected port 8080 not to be in window initially")
	}
}

func TestPortWindowInWindowAfterRegister(t *testing.T) {
	pw := newTestWindow()
	now := time.Now()
	pw.Register(8080, now.Add(-time.Second), now.Add(time.Hour))
	if !pw.InWindow(8080) {
		t.Fatal("expected port 8080 to be in window after register")
	}
}

func TestPortWindowNotInWindowBeforeStart(t *testing.T) {
	pw := newTestWindow()
	now := time.Now()
	pw.Register(9090, now.Add(time.Hour), now.Add(2*time.Hour))
	if pw.InWindow(9090) {
		t.Fatal("expected port 9090 not to be in window before start")
	}
}

func TestPortWindowNotInWindowAfterEnd(t *testing.T) {
	pw := newTestWindow()
	base := time.Now()
	pw.Register(7070, base.Add(-2*time.Hour), base.Add(-time.Hour))
	if pw.InWindow(7070) {
		t.Fatal("expected port 7070 not to be in window after end")
	}
}

func TestPortWindowRemoveClearsEntry(t *testing.T) {
	pw := newTestWindow()
	now := time.Now()
	pw.Register(3000, now.Add(-time.Second), now.Add(time.Hour))
	pw.Remove(3000)
	if pw.InWindow(3000) {
		t.Fatal("expected port 3000 not to be in window after remove")
	}
	if pw.Len() != 0 {
		t.Fatalf("expected len 0, got %d", pw.Len())
	}
}

func TestPortWindowEvictRemovesExpired(t *testing.T) {
	pw := newTestWindow()
	fixed := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	pw.now = func() time.Time { return fixed }

	pw.Register(1111, fixed.Add(-2*time.Hour), fixed.Add(-time.Hour)) // expired
	pw.Register(2222, fixed.Add(-time.Hour), fixed.Add(time.Hour))    // active

	pw.Evict()

	if pw.Len() != 1 {
		t.Fatalf("expected 1 entry after evict, got %d", pw.Len())
	}
	if _, ok := pw.Get(2222); !ok {
		t.Fatal("expected port 2222 to remain after evict")
	}
}

func TestPortWindowGetReturnsEntry(t *testing.T) {
	pw := newTestWindow()
	now := time.Now()
	start := now.Add(-time.Minute)
	end := now.Add(time.Minute)
	pw.Register(4444, start, end)

	e, ok := pw.Get(4444)
	if !ok {
		t.Fatal("expected to find entry for port 4444")
	}
	if e.Port != 4444 {
		t.Fatalf("expected port 4444, got %d", e.Port)
	}
	if !e.Start.Equal(start) {
		t.Fatalf("start mismatch: got %v", e.Start)
	}
}

func TestPortWindowIndependentPorts(t *testing.T) {
	pw := newTestWindow()
	now := time.Now()
	pw.Register(5000, now.Add(-time.Second), now.Add(time.Hour))

	if pw.InWindow(6000) {
		t.Fatal("port 6000 should not be in window")
	}
	if !pw.InWindow(5000) {
		t.Fatal("port 5000 should be in window")
	}
}
