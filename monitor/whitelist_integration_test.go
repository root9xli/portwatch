package monitor

import (
	"testing"
)

// TestWhitelistFiltersDiff verifies that whitelisted ports are removed from
// Diff results before alerting, simulating the monitor pipeline.
func TestWhitelistFiltersDiff(t *testing.T) {
	before := makeSnapshot([]int{80, 443})
	after := makeSnapshot([]int{80, 443, 8080, 9090})

	changes := Diff(before, after)
	if len(changes.Added) != 2 {
		t.Fatalf("expected 2 added ports, got %d", len(changes.Added))
	}

	w := NewWhitelist()
	w.Add(8080)

	var filtered []PortEvent
	for _, ev := range changes.Added {
		if !w.Contains(ev.Port) {
			filtered = append(filtered, ev)
		}
	}

	if len(filtered) != 1 {
		t.Fatalf("expected 1 event after whitelist filter, got %d", len(filtered))
	}
	if filtered[0].Port != 9090 {
		t.Errorf("expected port 9090, got %d", filtered[0].Port)
	}
}

// TestWhitelistAllowsRemoved verifies removed whitelisted ports are also filtered.
func TestWhitelistAllowsRemoved(t *testing.T) {
	before := makeSnapshot([]int{22, 80, 443})
	after := makeSnapshot([]int{80, 443})

	changes := Diff(before, after)

	w := NewWhitelist()
	w.Add(22)

	var filtered []PortEvent
	for _, ev := range changes.Removed {
		if !w.Contains(ev.Port) {
			filtered = append(filtered, ev)
		}
	}

	if len(filtered) != 0 {
		t.Errorf("expected 0 non-whitelisted removed ports, got %d", len(filtered))
	}
}
