package monitor

import (
	"testing"
	"time"
)

// TestPortLockSuppressesMessageForLockedPort verifies that a locked port's
// diff messages are filtered out before alerting.
func TestPortLockSuppressesMessageForLockedPort(t *testing.T) {
	pl := NewPortLock()
	pl.Lock(8080, "integration-test", time.Minute)

	msgs := []Message{
		{Port: 8080, Action: "added", Level: "warn"},
		{Port: 443, Action: "added", Level: "info"},
	}

	var passed []Message
	for _, m := range msgs {
		if !pl.IsLocked(m.Port) {
			passed = append(passed, m)
		}
	}

	if len(passed) != 1 {
		t.Fatalf("expected 1 message to pass, got %d", len(passed))
	}
	if passed[0].Port != 443 {
		t.Fatalf("expected port 443 to pass, got %d", passed[0].Port)
	}
}

// TestPortLockAllowsAfterExpiry verifies that once a lock expires,
// messages for that port are no longer suppressed.
func TestPortLockAllowsAfterExpiry(t *testing.T) {
	now := time.Now()
	pl := NewPortLock()
	pl.clock = func() time.Time { return now }
	pl.Lock(8080, "temp", 3*time.Second)

	// Advance past expiry.
	now = now.Add(5 * time.Second)

	msgs := []Message{
		{Port: 8080, Action: "added", Level: "warn"},
	}

	var passed []Message
	for _, m := range msgs {
		if !pl.IsLocked(m.Port) {
			passed = append(passed, m)
		}
	}

	if len(passed) != 1 {
		t.Fatalf("expected message to pass after lock expiry, got %d passed", len(passed))
	}
}
