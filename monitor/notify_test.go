package monitor

import (
	"log"
	"os"
	"strings"
	"testing"
)

func newTestNotifier(cmd string, args []string) *Notifier {
	return NewNotifier(cmd, args, log.New(os.Stderr, "test ", 0))
}

func TestNotifyEmptyCmdNoError(t *testing.T) {
	n := newTestNotifier("", nil)
	if err := n.Notify("hello"); err != nil {
		t.Fatalf("expected no error for empty cmd, got %v", err)
	}
}

func TestNotifyEchoCommand(t *testing.T) {
	n := newTestNotifier("echo", nil)
	if err := n.Notify("port 8080 opened"); err != nil {
		t.Fatalf("echo should succeed: %v", err)
	}
}

func TestNotifyBadCommandReturnsError(t *testing.T) {
	n := newTestNotifier("false", nil)
	err := n.Notify("msg")
	if err == nil {
		t.Fatal("expected error from failing command")
	}
	if !strings.Contains(err.Error(), "command failed") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestNotifyMultipleReturnsErrorsOnly(t *testing.T) {
	n := newTestNotifier("echo", nil)
	msgs := []string{"a", "b", "c"}
	errs := n.NotifyMultiple(msgs)
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
}

func TestNotifyMultipleCollectsErrors(t *testing.T) {
	n := newTestNotifier("false", nil)
	msgs := []string{"a", "b"}
	errs := n.NotifyMultiple(msgs)
	if len(errs) != 2 {
		t.Fatalf("expected 2 errors, got %d", len(errs))
	}
}
