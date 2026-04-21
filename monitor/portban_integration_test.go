package monitor

import (
	"testing"
	"time"
)

// TestPortBanSuppressesAlertForBannedPort verifies that a banned port does not
// produce an alert message when processed through a diff pipeline.
func TestPortBanSuppressesAlertForBannedPort(t *testing.T) {
	pb := NewPortBan(time.Minute)

	before := makeSnapshot([]int{80})
	after := makeSnapshot([]int{80, 9999})
	msgs := Diff(before, after)

	if len(msgs) != 1 {
		t.Fatalf("expected 1 diff message, got %d", len(msgs))
	}

	pb.Ban(9999, "known scanner")

	var allowed []DiffMessage
	for _, m := range msgs {
		if !pb.IsBanned(m.Port) {
			allowed = append(allowed, m)
		}
	}

	if len(allowed) != 0 {
		t.Fatalf("expected banned port to be filtered, got %d messages", len(allowed))
	}
}

// TestPortBanAllowsUnbannedPorts verifies that ports not in the ban list pass through.
func TestPortBanAllowsUnbannedPorts(t *testing.T) {
	pb := NewPortBan(time.Minute)

	before := makeSnapshot([]int{80})
	after := makeSnapshot([]int{80, 8080, 9090})
	msgs := Diff(before, after)

	pb.Ban(9090, "blocked")

	var allowed []DiffMessage
	for _, m := range msgs {
		if !pb.IsBanned(m.Port) {
			allowed = append(allowed, m)
		}
	}

	if len(allowed) != 1 {
		t.Fatalf("expected 1 allowed message (8080), got %d", len(allowed))
	}
	if allowed[0].Port != 8080 {
		t.Fatalf("expected port 8080, got %d", allowed[0].Port)
	}
}
