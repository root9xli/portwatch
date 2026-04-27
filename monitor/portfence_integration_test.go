package monitor

import (
	"strings"
	"testing"
	"time"
)

func TestPortFenceFiltersMessagesFromDiff(t *testing.T) {
	snap1 := makeSnapshot(80, 443)
	snap2 := makeSnapshot(80, 443, 8080)
	msgs := Diff(snap1, snap2)
	if len(msgs) != 1 {
		t.Fatalf("expected 1 diff message, got %d", len(msgs))
	}

	f := NewPortFence()
	f.Fence(8080, 5*time.Minute, "integration-test")

	out := f.Filter(msgs)
	if len(out) != 0 {
		t.Fatalf("expected fenced port to be filtered, got %d messages", len(out))
	}
}

func TestPortFenceAllowsUnfencedPortsFromDiff(t *testing.T) {
	snap1 := makeSnapshot(80)
	snap2 := makeSnapshot(80, 443, 8080)
	msgs := Diff(snap1, snap2)
	if len(msgs) != 2 {
		t.Fatalf("expected 2 diff messages, got %d", len(msgs))
	}

	f := NewPortFence()
	f.Fence(443, 5*time.Minute, "only-443-fenced")

	out := f.Filter(msgs)
	if len(out) != 1 {
		t.Fatalf("expected 1 message after filter, got %d", len(out))
	}
	if out[0].Port != 8080 {
		t.Fatalf("expected port 8080 to pass through, got %d", out[0].Port)
	}
}

func TestPortFenceReporterSummaryAfterFence(t *testing.T) {
	f := NewPortFence()
	f.Fence(9090, 10*time.Minute, "suspicious")
	f.Fence(7070, 2*time.Minute, "noisy")

	r := NewPortFenceReporter(f)
	summary := r.Summary()

	if !strings.Contains(summary, "9090") {
		t.Errorf("expected summary to contain port 9090")
	}
	if !strings.Contains(summary, "7070") {
		t.Errorf("expected summary to contain port 7070")
	}
	if !strings.Contains(summary, "suspicious") {
		t.Errorf("expected summary to contain reason 'suspicious'")
	}
	if r.Count() != 2 {
		t.Fatalf("expected count 2, got %d", r.Count())
	}
}

func TestPortFenceExpiredNotReportedAfterDiff(t *testing.T) {
	now := time.Now()
	f := NewPortFence()
	f.clock = func() time.Time { return now }
	f.Fence(3000, 1*time.Second, "short")

	now = now.Add(2 * time.Second)

	r := NewPortFenceReporter(f)
	if r.Count() != 0 {
		t.Fatalf("expected 0 active fences after expiry, got %d", r.Count())
	}
	if strings.Contains(r.Summary(), "3000") {
		t.Error("expired fence should not appear in summary")
	}
}
