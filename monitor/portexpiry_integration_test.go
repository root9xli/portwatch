package monitor

import (
	"testing"
	"time"
)

func TestPortExpiryWithDiffMessages(t *testing.T) {
	pe := NewPortExpiry()

	added := []Message{
		{Port: 8080, Action: "added", Level: "info"},
		{Port: 9090, Action: "added", Level: "info"},
	}

	// Register ports with short deadlines
	for _, msg := range added {
		pe.Register(msg.Port, time.Now().Add(50*time.Millisecond))
	}

	if pe.Len() != 2 {
		t.Fatalf("expected 2 registered, got %d", pe.Len())
	}

	// Before expiry, none should be overdue
	overdue := pe.Overdue()
	if len(overdue) != 0 {
		t.Fatalf("expected 0 overdue before deadline, got %d", len(overdue))
	}

	// Advance past deadline
	time.Sleep(60 * time.Millisecond)

	overdue = pe.Overdue()
	if len(overdue) != 2 {
		t.Fatalf("expected 2 overdue after deadline, got %d", len(overdue))
	}
}

func TestPortExpiryEvictClearsAfterDiff(t *testing.T) {
	pe := NewPortExpiry()

	pe.Register(8080, time.Now().Add(-1*time.Second))
	pe.Register(443, time.Now().Add(10*time.Minute))

	evicted := pe.Evict()
	if len(evicted) != 1 {
		t.Fatalf("expected 1 evicted, got %d", len(evicted))
	}
	if evicted[0] != 8080 {
		t.Fatalf("expected port 8080 evicted, got %d", evicted[0])
	}

	// Remaining port should still be tracked and not overdue
	if pe.IsOverdue(443) {
		t.Fatal("expected 443 to not be overdue")
	}
	if pe.Len() != 1 {
		t.Fatalf("expected 1 remaining after evict, got %d", pe.Len())
	}
}
