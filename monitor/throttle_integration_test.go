package monitor

import (
	"testing"
	"time"
)

// TestThrottleWithDiffMessages verifies throttle integrates with diff output.
func TestThrottleWithDiffMessages(t *testing.T) {
	th := NewThrottle(2, time.Minute)

	snap1 := makeSnapshot(80, 443)
	snap2 := makeSnapshot(80, 443, 9999)
	result := Diff(snap1, snap2)

	allowed := 0
	for _, msg := range result.Added {
		if th.Allow(msg.Port) {
			allowed++
		}
	}
	if allowed != 1 {
		t.Fatalf("expected 1 allowed alert, got %d", allowed)
	}
}

// TestThrottleExhaustedDoesNotAlert confirms no alerts pass once throttled.
func TestThrottleExhaustedDoesNotAlert(t *testing.T) {
	th := NewThrottle(1, time.Minute)

	snap1 := makeSnapshot(80)
	snap2 := makeSnapshot(80, 8080)

	result := Diff(snap1, snap2)
	_ = result

	th.Allow(8080) // consume token

	if th.Allow(8080) {
		t.Fatal("expected throttle to block second alert for same port")
	}
}
