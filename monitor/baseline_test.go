package monitor

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBaselineContainsAfterAdd(t *testing.T) {
	b := NewBaseline("")
	if b.Contains(8080) {
		t.Fatal("expected port not present before add")
	}
	b.Add(8080)
	if !b.Contains(8080) {
		t.Fatal("expected port present after add")
	}
}

func TestBaselineLen(t *testing.T) {
	b := NewBaseline("")
	b.Add(80)
	b.Add(443)
	b.Add(80) // duplicate
	if b.Len() != 2 {
		t.Fatalf("expected 2, got %d", b.Len())
	}
}

func TestBaselineSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")

	b1 := NewBaseline(path)
	b1.Add(22)
	b1.Add(3306)
	if err := b1.Save(); err != nil {
		t.Fatalf("save error: %v", err)
	}

	b2 := NewBaseline(path)
	if err := b2.Load(); err != nil {
		t.Fatalf("load error: %v", err)
	}
	if !b2.Contains(22) || !b2.Contains(3306) {
		t.Fatal("loaded baseline missing ports")
	}
}

func TestBaselineLoadMissingFile(t *testing.T) {
	b := NewBaseline("/nonexistent/path/baseline.json")
	if err := b.Load(); err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
}

func TestBaselineSaveBadPath(t *testing.T) {
	b := NewBaseline("/nonexistent/dir/baseline.json")
	b.Add(80)
	if err := b.Save(); err == nil {
		t.Fatal("expected error saving to bad path")
	}
}

func TestBaselineFirstSeenSet(t *testing.T) {
	b := NewBaseline("")
	b.Add(9090)
	e, ok := b.entries[9090]
	if !ok {
		t.Fatal("entry not found")
	}
	if e.FirstSeen.IsZero() {
		t.Fatal("FirstSeen should not be zero")
	}
	_ = os.DevNull
}
