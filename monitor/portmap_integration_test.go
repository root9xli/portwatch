package monitor

import (
	"strings"
	"testing"
)

// TestPortMapEnrichesMessage verifies that PortMap aliases can be used
// to annotate diff messages before they are dispatched.
func TestPortMapEnrichesMessage(t *testing.T) {
	pm := NewPortMap()
	pm.Set(9200, "elasticsearch", "search cluster")

	msgs := []Message{
		{Port: 9200, Action: "added", Level: "warn"},
		{Port: 8080, Action: "added", Level: "info"},
	}

	for i, m := range msgs {
		if e, ok := pm.Get(m.Port); ok {
			msgs[i].Text = e.Alias + ": " + m.Action
		} else {
			msgs[i].Text = m.Action
		}
	}

	if !strings.Contains(msgs[0].Text, "elasticsearch") {
		t.Errorf("expected elasticsearch alias in message, got: %s", msgs[0].Text)
	}
	if strings.Contains(msgs[1].Text, ":") {
		// port 8080 has no alias so text should just be the action
		t.Errorf("unexpected alias for unmapped port: %s", msgs[1].Text)
	}
}

// TestPortMapSummaryAfterMultipleSets ensures summary lists all ports
// sorted and contains expected content after several Set calls.
func TestPortMapSummaryAfterMultipleSets(t *testing.T) {
	pm := NewPortMap()
	pm.Set(443, "https", "TLS traffic")
	pm.Set(80, "http", "plain traffic")
	pm.Set(5432, "postgres", "database")

	s := pm.Summary()

	ports := []string{"80", "443", "5432"}
	for _, p := range ports {
		if !strings.Contains(s, p) {
			t.Errorf("expected port %s in summary", p)
		}
	}

	// Verify sorted order: 80 should appear before 443
	idx80 := strings.Index(s, "80")
	idx443 := strings.Index(s, "443")
	if idx80 > idx443 {
		t.Error("expected port 80 to appear before 443 in summary")
	}
}

// TestPortMapOverwrite verifies that calling Set on an already-mapped port
// replaces the existing entry rather than leaving stale data.
func TestPortMapOverwrite(t *testing.T) {
	pm := NewPortMap()
	pm.Set(6379, "redis-old", "old description")
	pm.Set(6379, "redis", "cache layer")

	e, ok := pm.Get(6379)
	if !ok {
		t.Fatal("expected port 6379 to be present after overwrite")
	}
	if e.Alias != "redis" {
		t.Errorf("expected alias 'redis' after overwrite, got: %s", e.Alias)
	}
	if e.Description != "cache layer" {
		t.Errorf("expected description 'cache layer' after overwrite, got: %s", e.Description)
	}
}
