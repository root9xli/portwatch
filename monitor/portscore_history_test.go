package monitor

import (
	"testing"
	"time"
)

func newTestScoreHistory() *PortScoreHistory {
	return NewPortScoreHistory(5, 10*time.Minute)
}

func TestPortScoreHistoryEmptyGet(t *testing.T) {
	h := newTestScoreHistory()
	entries := h.Get(8080)
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

func TestPortScoreHistoryRecordAndGet(t *testing.T) {
	h := newTestScoreHistory()
	h.Record(8080, 3.5, "warn")
	entries := h.Get(8080)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Score != 3.5 {
		t.Errorf("expected score 3.5, got %.2f", entries[0].Score)
	}
	if entries[0].Level != "warn" {
		t.Errorf("expected level warn, got %s", entries[0].Level)
	}
	if entries[0].Port != 8080 {
		t.Errorf("expected port 8080, got %d", entries[0].Port)
	}
}

func TestPortScoreHistoryIndependentPorts(t *testing.T) {
	h := newTestScoreHistory()
	h.Record(80, 1.0, "info")
	h.Record(443, 7.0, "critical")

	if len(h.Get(80)) != 1 {
		t.Error("expected 1 entry for port 80")
	}
	if len(h.Get(443)) != 1 {
		t.Error("expected 1 entry for port 443")
	}
}

func TestPortScoreHistoryCapEnforced(t *testing.T) {
	h := NewPortScoreHistory(3, 10*time.Minute)
	for i := 0; i < 6; i++ {
		h.Record(9000, float64(i), "info")
	}
	entries := h.Get(9000)
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries (cap), got %d", len(entries))
	}
	// Should retain the most recent
	if entries[0].Score != 3.0 {
		t.Errorf("expected first retained score=3, got %.2f", entries[0].Score)
	}
}

func TestPortScoreHistoryEvictsOldEntries(t *testing.T) {
	h := NewPortScoreHistory(10, 50*time.Millisecond)
	h.Record(7000, 5.0, "warn")
	time.Sleep(80 * time.Millisecond)
	entries := h.Get(7000)
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries after expiry, got %d", len(entries))
	}
}

func TestPortScoreHistorySummaryEmpty(t *testing.T) {
	h := newTestScoreHistory()
	s := h.Summary(1234)
	if s == "" {
		t.Error("expected non-empty summary")
	}
}

func TestPortScoreHistorySummaryContainsPort(t *testing.T) {
	h := newTestScoreHistory()
	h.Record(8443, 6.1, "critical")
	s := h.Summary(8443)
	if !contains(s, "8443") {
		t.Errorf("expected summary to contain port 8443, got: %s", s)
	}
	if !contains(s, "critical") {
		t.Errorf("expected summary to contain level critical, got: %s", s)
	}
}
