package monitor

import (
	"testing"
	"time"
)

func TestThrottleAllowsUpToMax(t *testing.T) {
	th := NewThrottle(3, time.Minute)
	for i := 0; i < 3; i++ {
		if !th.Allow(8080) {
			t.Fatalf("expected allow on call %d", i+1)
		}
	}
}

func TestThrottleBlocksOverMax(t *testing.T) {
	th := NewThrottle(2, time.Minute)
	th.Allow(9000)
	th.Allow(9000)
	if th.Allow(9000) {
		t.Fatal("expected block after max tokens consumed")
	}
}

func TestThrottleIndependentPorts(t *testing.T) {
	th := NewThrottle(1, time.Minute)
	th.Allow(80)
	if !th.Allow(443) {
		t.Fatal("expected allow for different port")
	}
}

func TestThrottleRefillAfterWindow(t *testing.T) {
	th := NewThrottle(1, 10*time.Millisecond)
	th.Allow(8080)
	if th.Allow(8080) {
		t.Fatal("expected block before refill")
	}
	time.Sleep(15 * time.Millisecond)
	if !th.Allow(8080) {
		t.Fatal("expected allow after refill window")
	}
}

func TestThrottleResetRestoresTokens(t *testing.T) {
	th := NewThrottle(1, time.Minute)
	th.Allow(3000)
	th.Reset(3000)
	if !th.Allow(3000) {
		t.Fatal("expected allow after reset")
	}
}

func TestThrottlePurgeRemovesStale(t *testing.T) {
	th := NewThrottle(2, time.Minute)
	th.Allow(7070)
	th.Allow(7070) // exhaust
	cutoff := time.Now().Add(time.Second)
	th.Purge(cutoff)
	// after purge, tokens should be refilled on next Allow
	if !th.Allow(7070) {
		t.Fatal("expected allow after purge cleared state")
	}
}
