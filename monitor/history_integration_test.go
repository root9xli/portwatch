package monitor

import (
	"strings"
	"testing"
)

// TestHistoryRecordedOnDiff verifies that History.Record is called correctly
// when processing diff results, simulating monitor integration.
func TestHistoryRecordedOnDiff(t *testing.T) {
	h := NewHistory(50)

	old := makeSnapshot([]int{80, 443})
	new_ := makeSnapshot([]int{80, 443, 8080})
	result := Diff(old, new_)

	for _, port := range result.Added {
		h.Record(port, "added", "warn")
	}
	for _, port := range result.Removed {
		h.Record(port, "removed", "info")
	}

	if h.Len() != 1 {
		t.Fatalf("expected 1 history entry, got %d", h.Len())
	}
	entries := h.Last(1)
	if entries[0].Port != 8080 {
		t.Errorf("expected port 8080, got %d", entries[0].Port)
	}
	if entries[0].Action != "added" {
		t.Errorf("expected action added, got %s", entries[0].Action)
	}
}

func TestReporterSummaryAfterMultipleDiffs(t *testing.T) {
	h := NewHistory(50)
	r := NewReporter(h)

	h.Record(22, "added", "critical")
	h.Record(8080, "added", "warn")
	h.Record(22, "removed", "info")

	summary := r.Summary(10)
	if !strings.Contains(summary, "3 recent event") {
		t.Errorf("expected 3 events in summary, got: %s", summary)
	}
	if !r.HasCritical() {
		t.Error("expected HasCritical true")
	}
	counts := r.CountByLevel()
	if counts["info"] != 1 || counts["warn"] != 1 || counts["critical"] != 1 {
		t.Errorf("unexpected counts: %v", counts)
	}
}

// TestHistoryRecordedOnRemoval verifies that removed ports are correctly
// recorded in history when a diff detects a port disappearing.
func TestHistoryRecordedOnRemoval(t *testing.T) {
	h := NewHistory(50)

	old := makeSnapshot([]int{80, 443, 9090})
	new_ := makeSnapshot([]int{80, 443})
	result := Diff(old, new_)

	for _, port := range result.Added {
		h.Record(port, "added", "warn")
	}
	for _, port := range result.Removed {
		h.Record(port, "removed", "info")
	}

	if h.Len() != 1 {
		t.Fatalf("expected 1 history entry, got %d", h.Len())
	}
	entries := h.Last(1)
	if entries[0].Port != 9090 {
		t.Errorf("expected port 9090, got %d", entries[0].Port)
	}
	if entries[0].Action != "removed" {
		t.Errorf("expected action removed, got %s", entries[0].Action)
	}
	if entries[0].Level != "info" {
		t.Errorf("expected level info, got %s", entries[0].Level)
	}
}
