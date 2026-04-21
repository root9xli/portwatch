package monitor

import (
	"testing"
	"time"
)

// TestPortWatcherWithDiffMessages verifies that watched ports appear in diff results.
func TestPortWatcherWithDiffMessages(t *testing.T) {
	watcher := NewPortWatcher()
	watcher.Watch(9999, "test-service", 0)

	before := makeSnapshot([]int{80, 443})
	after := makeSnapshot([]int{80, 443, 9999})
	msgs := Diff(before, after)

	if len(msgs) != 1 {
		t.Fatalf("expected 1 diff message, got %d", len(msgs))
	}

	port := msgs[0].Port
	if !watcher.IsWatched(port) {
		t.Fatalf("expected port %d to be watched", port)
	}
}

// TestPortWatcherExpiredNotFlagged verifies expired watches are no longer active.
func TestPortWatcherExpiredNotFlagged(t *testing.T) {
	watcher := NewPortWatcher()
	now := time.Now()
	watcher.now = func() time.Time { return now }
	watcher.Watch(7777, "ephemeral", 1*time.Millisecond)

	// advance time past TTL
	watcher.now = func() time.Time { return now.Add(10 * time.Millisecond) }

	if watcher.IsWatched(7777) {
		t.Fatal("expected watch to have expired")
	}
	if watcher.Len() != 0 {
		t.Fatalf("expected 0 active watches, got %d", watcher.Len())
	}
}

// TestPortWatcherSummaryAfterDiff verifies summary reflects current watch state.
func TestPortWatcherSummaryAfterDiff(t *testing.T) {
	watcher := NewPortWatcher()
	watcher.Watch(5432, "postgres", 0)
	watcher.Watch(6379, "redis", 0)

	before := makeSnapshot([]int{})
	after := makeSnapshot([]int{5432, 6379})
	msgs := Diff(before, after)

	if len(msgs) != 2 {
		t.Fatalf("expected 2 diff messages, got %d", len(msgs))
	}

	summary := watcher.Summary()
	for _, m := range msgs {
		portStr := itoa(m.Port)
		if !contains(summary, portStr) {
			t.Errorf("expected summary to contain port %s", portStr)
		}
	}
}
