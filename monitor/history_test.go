package monitor

import (
	"testing"
	"time"
)

func TestHistoryRecordAndLen(t *testing.T) {
	h := NewHistory(10)
	h.Record(8080, "added", "warn")
	h.Record(9090, "removed", "info")
	if h.Len() != 2 {
		t.Fatalf("expected 2 entries, got %d", h.Len())
	}
}

func TestHistoryLastReturnsRecent(t *testing.T) {
	h := NewHistory(10)
	h.Record(1111, "added", "warn")
	h.Record(2222, "added", "critical")
	h.Record(3333, "removed", "info")
	entries := h.Last(2)
	if len(entries) != 2 {
		t.Fatalf("expected 2, got %d", len(entries))
	}
	if entries[0].Port != 2222 || entries[1].Port != 3333 {
		t.Errorf("unexpected ports: %v", entries)
	}
}

func TestHistoryLastAllWhenFewerThanN(t *testing.T) {
	h := NewHistory(10)
	h.Record(80, "added", "info")
	result := h.Last(5)
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
}

func TestHistoryEvictsOldest(t *testing.T) {
	h := NewHistory(3)
	h.Record(1, "added", "info")
	h.Record(2, "added", "info")
	h.Record(3, "added", "info")
	h.Record(4, "added", "info")
	if h.Len() != 3 {
		t.Fatalf("expected 3, got %d", h.Len())
	}
	all := h.Last(3)
	if all[0].Port != 2 {
		t.Errorf("expected oldest to be evicted, got port %d", all[0].Port)
	}
}

func TestHistoryTimestampSet(t *testing.T) {
	before := time.Now()
	h := NewHistory(5)
	h.Record(443, "added", "warn")
	after := time.Now()
	entries := h.Last(1)
	if entries[0].Timestamp.Before(before) || entries[0].Timestamp.After(after) {
		t.Errorf("timestamp out of expected range")
	}
}

func TestHistoryClear(t *testing.T) {
	h := NewHistory(5)
	h.Record(80, "added", "info")
	h.Record(443, "added", "info")
	h.Clear()
	if h.Len() != 0 {
		t.Fatalf("expected 0 after clear, got %d", h.Len())
	}
}
