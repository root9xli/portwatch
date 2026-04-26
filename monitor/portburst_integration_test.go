package monitor

import (
	"strings"
	"testing"
	"time"
)

func TestPortBurstDetectedFromRepeatedDiff(t *testing.T) {
	pb := NewPortBurst(30*time.Second, 3)
	reporter := NewPortBurstReporter(pb)

	msgs := []Message{
		{Port: 8080, Level: "warn", Action: "added"},
	}

	// Simulate three consecutive diff cycles adding the same port.
	reporter.Record(msgs)
	reporter.Record(msgs)
	reporter.Record(msgs)

	events := reporter.Events()
	if len(events) == 0 {
		t.Fatal("expected at least one burst event after threshold reached")
	}

	found := false
	for _, ev := range events {
		if ev.Port == 8080 && ev.Count >= 3 {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected burst event for port 8080 with count >= 3, events: %+v", events)
	}
}

func TestPortBurstSummaryAfterDiff(t *testing.T) {
	pb := NewPortBurst(30*time.Second, 2)
	reporter := NewPortBurstReporter(pb)

	msgs := []Message{
		{Port: 9090, Level: "warn", Action: "added"},
	}
	reporter.Record(msgs)
	reporter.Record(msgs)

	summary := reporter.Summary()
	if !strings.Contains(summary, "9090") {
		t.Errorf("expected summary to contain port 9090, got: %s", summary)
	}
	if !strings.Contains(summary, "burst") {
		t.Errorf("expected summary to mention burst, got: %s", summary)
	}
}

func TestPortBurstStablePortNotFlagged(t *testing.T) {
	pb := NewPortBurst(30*time.Second, 5)
	reporter := NewPortBurstReporter(pb)

	msgs := []Message{
		{Port: 443, Level: "info", Action: "added"},
	}
	// Only two observations — below threshold of 5.
	reporter.Record(msgs)
	reporter.Record(msgs)

	if events := reporter.Events(); len(events) != 0 {
		t.Errorf("expected no burst events for stable port, got %d", len(events))
	}

	summary := reporter.Summary()
	if strings.Contains(summary, "443") {
		t.Errorf("stable port 443 should not appear in burst summary")
	}
}
