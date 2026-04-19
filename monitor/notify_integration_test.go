package monitor

import (
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

// TestNotifierWithDiffMessages verifies that a Notifier can process real
// formatted diff messages produced by FormatMultiple.
func TestNotifierWithDiffMessages(t *testing.T) {
	snap1 := makeSnapshot([]int{80, 443})
	snap2 := makeSnapshot([]int{80, 443, 9090})
	changes := Diff(snap1, snap2)

	if len(changes) == 0 {
		t.Fatal("expected diff changes")
	}

	msgs := FormatMultiple(changes, "testhost", time.Now())
	if len(msgs) == 0 {
		t.Fatal("expected formatted messages")
	}

	n := NewNotifier("echo", nil, log.New(os.Stderr, "", 0))
	errs := n.NotifyMultiple(msgs)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
}

// TestNotifierMessageContainsPort ensures the message passed to the command
// actually contains the port string (echo output not captured, but no panic).
func TestNotifierMessageContainsPort(t *testing.T) {
	snap1 := makeSnapshot([]int{3000})
	snap2 := makeSnapshot([]int{3000, 4000})
	changes := Diff(snap1, snap2)
	msgs := FormatMultiple(changes, "host", time.Now())

	for _, m := range msgs {
		if !strings.Contains(m, "4000") {
			t.Errorf("expected port 4000 in message: %q", m)
		}
	}
}
