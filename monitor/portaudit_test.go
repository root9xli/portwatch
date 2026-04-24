package monitor

import (
	"strings"
	"testing"
)

func newTestAudit() *PortAudit {
	return NewPortAudit(10)
}

func TestPortAuditEmptyLen(t *testing.T) {
	a := newTestAudit()
	if a.Len() != 0 {
		t.Fatalf("expected 0, got %d", a.Len())
	}
}

func TestPortAuditRecordIncreasesLen(t *testing.T) {
	a := newTestAudit()
	a.Record(8080, "added", "warn")
	a.Record(443, "removed", "info")
	if a.Len() != 2 {
		t.Fatalf("expected 2, got %d", a.Len())
	}
}

func TestPortAuditEntriesForPort(t *testing.T) {
	a := newTestAudit()
	a.Record(8080, "added", "warn")
	a.Record(443, "removed", "info")
	a.Record(8080, "removed", "info")
	entries := a.EntriesForPort(8080)
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries for port 8080, got %d", len(entries))
	}
	for _, e := range entries {
		if e.Port != 8080 {
			t.Errorf("unexpected port %d in entries", e.Port)
		}
	}
}

func TestPortAuditEvictsOldestWhenFull(t *testing.T) {
	a := NewPortAudit(3)
	a.Record(1, "added", "info")
	a.Record(2, "added", "info")
	a.Record(3, "added", "info")
	a.Record(4, "added", "warn")
	if a.Len() != 3 {
		t.Fatalf("expected 3 after eviction, got %d", a.Len())
	}
	// port 1 should have been evicted
	if len(a.EntriesForPort(1)) != 0 {
		t.Error("expected port 1 to be evicted")
	}
	if len(a.EntriesForPort(4)) != 1 {
		t.Error("expected port 4 to be present")
	}
}

func TestPortAuditSummaryEmpty(t *testing.T) {
	a := newTestAudit()
	s := a.Summary()
	if !strings.Contains(s, "no entries") {
		t.Errorf("expected 'no entries' in summary, got: %s", s)
	}
}

func TestPortAuditSummaryContainsPort(t *testing.T) {
	a := newTestAudit()
	a.Record(9090, "added", "critical")
	s := a.Summary()
	if !strings.Contains(s, "9090") {
		t.Errorf("expected port 9090 in summary, got: %s", s)
	}
	if !strings.Contains(s, "critical") {
		t.Errorf("expected level 'critical' in summary, got: %s", s)
	}
}
