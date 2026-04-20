package monitor_test

import (
	"testing"
)

// TestPortGroupWithDiffMessages verifies that port groups correctly categorize
// ports discovered during a diff cycle.
func TestPortGroupWithDiffMessages(t *testing.T) {
	group := NewPortGroup()
	group.Add("web", 80, 443, 8080, 8443)
	group.Add("database", 5432, 3306, 27017)

	snap1 := makeSnapshot(80, 443)
	snap2 := makeSnapshot(80, 443, 5432)

	msgs := Diff(snap1, snap2)
	if len(msgs) != 1 {
		t.Fatalf("expected 1 diff message, got %d", len(msgs))
	}

	port := msgs[0].Port
	groupName, ok := group.GroupOf(port)
	if !ok {
		t.Fatalf("expected port %d to belong to a group", port)
	}
	if groupName != "database" {
		t.Errorf("expected group 'database', got %q", groupName)
	}
}

// TestPortGroupUnknownPortInDiff verifies that ports not in any group are
// correctly reported as unknown during a diff cycle.
func TestPortGroupUnknownPortInDiff(t *testing.T) {
	group := NewPortGroup()
	group.Add("web", 80, 443)

	snap1 := makeSnapshot(80)
	snap2 := makeSnapshot(80, 9999)

	msgs := Diff(snap1, snap2)
	if len(msgs) != 1 {
		t.Fatalf("expected 1 diff message, got %d", len(msgs))
	}

	port := msgs[0].Port
	_, ok := group.GroupOf(port)
	if ok {
		t.Errorf("expected port %d to be unknown, but it matched a group", port)
	}
}

// TestPortGroupMultipleAddedPorts verifies group membership across multiple
// new ports appearing in a single diff cycle.
func TestPortGroupMultipleAddedPorts(t *testing.T) {
	group := NewPortGroup()
	group.Add("web", 80, 443, 8080)
	group.Add("cache", 6379, 11211)

	snap1 := makeSnapshot(80)
	snap2 := makeSnapshot(80, 8080, 6379)

	msgs := Diff(snap1, snap2)
	if len(msgs) != 2 {
		t.Fatalf("expected 2 diff messages, got %d", len(msgs))
	}

	groupCounts := map[string]int{}
	for _, msg := range msgs {
		if name, ok := group.GroupOf(msg.Port); ok {
			groupCounts[name]++
		}
	}

	if groupCounts["web"] != 1 {
		t.Errorf("expected 1 web port, got %d", groupCounts["web"])
	}
	if groupCounts["cache"] != 1 {
		t.Errorf("expected 1 cache port, got %d", groupCounts["cache"])
	}
}
