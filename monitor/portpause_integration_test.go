package monitor

import (
	"strings"
	"testing"
	"time"
)

func TestPortPauseFiltersMessagesFromDiff(t *testing.T) {
	p := NewPortPause()
	now := time.Unix(2000, 0)
	p.now = func() time.Time { return now }
	p.Pause(8080, 30*time.Second, "deploy")

	snap1 := makeSnapshot()
	snap2 := makeSnapshot(8080, 9090)
	msgs := Diff(snap1, snap2)

	filtered := p.Filter(msgs)
	for _, m := range filtered {
		if m.Port == 8080 {
			t.Fatalf("port 8080 should have been filtered out (paused)")
		}
	}

	var found9090 bool
	for _, m := range filtered {
		if m.Port == 9090 {
			found9090 = true
		}
	}
	if !found9090 {
		t.Fatal("expected port 9090 to pass through filter")
	}
}

func TestPortPauseAllowsPortAfterExpiry(t *testing.T) {
	base := time.Unix(2000, 0)
	p := NewPortPause()
	p.now = func() time.Time { return base }
	p.Pause(8080, 5*time.Second, "short pause")

	// advance past expiry
	p.now = func() time.Time { return base.Add(10 * time.Second) }

	snap1 := makeSnapshot()
	snap2 := makeSnapshot(8080)
	msgs := Diff(snap1, snap2)
	filtered := p.Filter(msgs)

	if len(filtered) == 0 {
		t.Fatal("expected port 8080 to pass through after pause expiry")
	}
}

func TestPortPauseReporterSummaryAfterPause(t *testing.T) {
	base := time.Unix(3000, 0)
	p := NewPortPause()
	p.now = func() time.Time { return base }
	p.Pause(443, 60*time.Second, "cert-renewal")
	p.Pause(80, 120*time.Second, "maintenance")

	r := NewPortPauseReporter(p)
	summary := r.Summary()

	if !strings.Contains(summary, "443") {
		t.Errorf("expected summary to contain port 443, got: %s", summary)
	}
	if !strings.Contains(summary, "cert-renewal") {
		t.Errorf("expected summary to contain reason 'cert-renewal', got: %s", summary)
	}
	if !strings.Contains(summary, "80") {
		t.Errorf("expected summary to contain port 80, got: %s", summary)
	}
	if r.Count() != 2 {
		t.Fatalf("expected count 2, got %d", r.Count())
	}
}
