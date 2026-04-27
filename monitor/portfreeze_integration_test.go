package monitor

import (
	"strings"
	"testing"
	"time"
)

func TestPortFreezeFiltersMessagesFromDiff(t *testing.T) {
	pf := NewPortFreeze()
	pf.Freeze(8080, 5*time.Minute, "deploy")

	msgs := []Message{
		{Port: 8080, Action: "added", Level: "warn"},
		{Port: 9090, Action: "added", Level: "warn"},
	}

	filtered := pf.Filter(msgs)
	if len(filtered) != 1 {
		t.Fatalf("expected 1 message after filter, got %d", len(filtered))
	}
	if filtered[0].Port != 9090 {
		t.Fatalf("expected port 9090 to pass through, got %d", filtered[0].Port)
	}
}

func TestPortFreezeAllowsUnfrozenPortsFromDiff(t *testing.T) {
	pf := NewPortFreeze()

	msgs := []Message{
		{Port: 80, Action: "added", Level: "info"},
		{Port: 443, Action: "added", Level: "info"},
	}

	filtered := pf.Filter(msgs)
	if len(filtered) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(filtered))
	}
}

func TestPortFreezeReporterSummaryAfterFreeze(t *testing.T) {
	pf := NewPortFreeze()
	pf.Freeze(8080, 10*time.Minute, "maintenance")
	pf.Freeze(443, 30*time.Minute, "cert-renewal")

	rep := NewPortFreezeReporter(pf)
	summary := rep.Summary()

	if !strings.Contains(summary, "8080") {
		t.Errorf("expected summary to contain port 8080")
	}
	if !strings.Contains(summary, "maintenance") {
		t.Errorf("expected summary to contain reason 'maintenance'")
	}
	if !strings.Contains(summary, "443") {
		t.Errorf("expected summary to contain port 443")
	}
	if rep.Count() != 2 {
		t.Errorf("expected count 2, got %d", rep.Count())
	}
}

func TestPortFreezeExpiredNotReportedAfterEvict(t *testing.T) {
	pf := NewPortFreeze()
	now := time.Now()
	pf.clock = func() time.Time { return now }
	pf.Freeze(8080, 1*time.Second, "short")
	pf.clock = func() time.Time { return now.Add(2 * time.Second) }
	pf.Evict()

	rep := NewPortFreezeReporter(pf)
	if rep.Count() != 0 {
		t.Errorf("expected count 0 after evict, got %d", rep.Count())
	}
	if strings.Contains(rep.Summary(), "8080") {
		t.Errorf("expected expired port to be absent from summary")
	}
}
