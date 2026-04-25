package monitor

import (
	"strings"
	"testing"
)

// TestPortClassifyWithDiffMessages verifies that added ports from a Diff are
// correctly classified and appear in the reporter output.
func TestPortClassifyWithDiffMessages(t *testing.T) {
	before := makeSnapshot([]int{22, 80})
	after := makeSnapshot([]int{22, 80, 8080, 55000})
	diff := Diff(before, after)

	if len(diff.Added) != 2 {
		t.Fatalf("expected 2 added ports, got %d", len(diff.Added))
	}

	classifier := NewPortClassifier()
	reporter := NewPortClassifyReporter(classifier)

	var added []int
	for _, m := range diff.Added {
		added = append(added, m.Port)
	}

	summary := reporter.Report(added)
	if !strings.Contains(summary, "registered") {
		t.Errorf("expected 'registered' class in summary, got: %s", summary)
	}
	if !strings.Contains(summary, "ephemeral") {
		t.Errorf("expected 'ephemeral' class in summary, got: %s", summary)
	}
}

// TestPortClassifyAllClassesPresentAfterDiff checks CountByClass reflects
// the correct distribution after a diff that spans multiple port ranges.
func TestPortClassifyAllClassesPresentAfterDiff(t *testing.T) {
	before := makeSnapshot([]int{})
	after := makeSnapshot([]int{443, 5000, 60000})
	diff := Diff(before, after)

	classifier := NewPortClassifier()
	reporter := NewPortClassifyReporter(classifier)

	var ports []int
	for _, m := range diff.Added {
		ports = append(ports, m.Port)
	}

	counts := reporter.CountByClass(ports)
	if counts[ClassSystem] != 1 {
		t.Errorf("expected 1 system port, got %d", counts[ClassSystem])
	}
	if counts[ClassRegistered] != 1 {
		t.Errorf("expected 1 registered port, got %d", counts[ClassRegistered])
	}
	if counts[ClassEphemeral] != 1 {
		t.Errorf("expected 1 ephemeral port, got %d", counts[ClassEphemeral])
	}
}
