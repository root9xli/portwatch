package monitor

import (
	"testing"
)

func makeSnapshot(ports ...int) *Snapshot {
	var entries []PortEntry
	for _, p := range ports {
		entries = append(entries, PortEntry{Port: p, Protocol: "tcp"})
	}
	return &Snapshot{Ports: entries}
}

func TestDiffNoChange(t *testing.T) {
	prev := makeSnapshot(80, 443)
	curr := makeSnapshot(80, 443)
	result := Diff(prev, curr)
	if len(result.Added) != 0 || len(result.Removed) != 0 {
		t.Errorf("expected no diff, got added=%v removed=%v", result.Added, result.Removed)
	}
}

func TestDiffDetectsAdded(t *testing.T) {
	prev := makeSnapshot(80)
	curr := makeSnapshot(80, 9000)
	result := Diff(prev, curr)
	if len(result.Added) != 1 || result.Added[0].Port != 9000 {
		t.Errorf("expected port 9000 added, got %v", result.Added)
	}
	if len(result.Removed) != 0 {
		t.Errorf("expected no removed ports, got %v", result.Removed)
	}
}

func TestDiffDetectsRemoved(t *testing.T) {
	prev := makeSnapshot(80, 8080)
	curr := makeSnapshot(80)
	result := Diff(prev, curr)
	if len(result.Removed) != 1 || result.Removed[0].Port != 8080 {
		t.Errorf("expected port 8080 removed, got %v", result.Removed)
	}
	if len(result.Added) != 0 {
		t.Errorf("expected no added ports, got %v", result.Added)
	}
}

func TestDiffBothChanges(t *testing.T) {
	prev := makeSnapshot(80, 443)
	curr := makeSnapshot(443, 9000)
	result := Diff(prev, curr)
	if len(result.Added) != 1 || result.Added[0].Port != 9000 {
		t.Errorf("unexpected added: %v", result.Added)
	}
	if len(result.Removed) != 1 || result.Removed[0].Port != 80 {
		t.Errorf("unexpected removed: %v", result.Removed)
	}
}
