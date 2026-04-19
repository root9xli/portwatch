package monitor

import (
	"path/filepath"
	"testing"
)

func TestBaselinePopulatedFromDiff(t *testing.T) {
	b := NewBaseline("")

	prev := makeSnapshot(80, 443)
	curr := makeSnapshot(80, 443, 8080)
	msgs := Diff(prev, curr)

	for _, m := range msgs {
		if m.Action == "added" {
			b.Add(m.Port)
		}
	}

	if !b.Contains(8080) {
		t.Fatal("expected 8080 in baseline after diff")
	}
	if b.Contains(80) || b.Contains(443) {
		t.Fatal("existing ports should not be in baseline")
	}
}

func TestBaselinePersistenceRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bl.json")

	b1 := NewBaseline(path)
	prev := makeSnapshot()
	curr := makeSnapshot(22, 3306)
	msgs := Diff(prev, curr)
	for _, m := range msgs {
		b1.Add(m.Port)
	}
	if err := b1.Save(); err != nil {
		t.Fatalf("save: %v", err)
	}

	b2 := NewBaseline(path)
	if err := b2.Load(); err != nil {
		t.Fatalf("load: %v", err)
	}
	if !b2.Contains(22) || !b2.Contains(3306) {
		t.Fatal("ports missing after round-trip")
	}

	r := NewBaselineReporter(b2)
	if b2.Len() != 2 {
		t.Fatalf("expected 2 ports, got %d", b2.Len())
	}
	_ = r.Summary()
}
