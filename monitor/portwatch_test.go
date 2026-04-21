package monitor

import (
	"testing"
	"time"
)

func newTestWatcher() *PortWatcher {
	w := NewPortWatcher()
	return w
}

func TestPortWatcherWatchAndIsWatched(t *testing.T) {
	w := newTestWatcher()
	w.Watch(8080, "http", 0)
	if !w.IsWatched(8080) {
		t.Fatal("expected port 8080 to be watched")
	}
}

func TestPortWatcherUnwatch(t *testing.T) {
	w := newTestWatcher()
	w.Watch(9090, "custom", 0)
	w.Unwatch(9090)
	if w.IsWatched(9090) {
		t.Fatal("expected port 9090 to be unwatched after removal")
	}
}

func TestPortWatcherNotWatchedInitially(t *testing.T) {
	w := newTestWatcher()
	if w.IsWatched(443) {
		t.Fatal("expected port 443 not to be watched initially")
	}
}

func TestPortWatcherExpiresAfterTTL(t *testing.T) {
	w := newTestWatcher()
	now := time.Now()
	w.now = func() time.Time { return now }
	w.Watch(3000, "temp", 10*time.Millisecond)

	if !w.IsWatched(3000) {
		t.Fatal("expected port 3000 to be watched before expiry")
	}

	w.now = func() time.Time { return now.Add(20 * time.Millisecond) }
	if w.IsWatched(3000) {
		t.Fatal("expected port 3000 to be expired")
	}
}

func TestPortWatcherLenExcludesExpired(t *testing.T) {
	w := newTestWatcher()
	now := time.Now()
	w.now = func() time.Time { return now }
	w.Watch(1111, "a", 5*time.Millisecond)
	w.Watch(2222, "b", 0)

	w.now = func() time.Time { return now.Add(10 * time.Millisecond) }
	if w.Len() != 1 {
		t.Fatalf("expected len 1 after expiry, got %d", w.Len())
	}
}

func TestPortWatcherSummaryEmpty(t *testing.T) {
	w := newTestWatcher()
	s := w.Summary()
	if s != "watched ports: none" {
		t.Fatalf("unexpected summary: %q", s)
	}
}

func TestPortWatcherSummaryContainsPort(t *testing.T) {
	w := newTestWatcher()
	w.Watch(8443, "https", 0)
	s := w.Summary()
	if !contains(s, "8443") {
		t.Fatalf("expected summary to contain port 8443, got: %s", s)
	}
	if !contains(s, "https") {
		t.Fatalf("expected summary to contain note 'https', got: %s", s)
	}
}

func TestPortWatcherAllReturnsEntries(t *testing.T) {
	w := newTestWatcher()
	w.Watch(80, "http", 0)
	w.Watch(443, "https", 0)
	if len(w.All()) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(w.All()))
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		(func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		})())
}
