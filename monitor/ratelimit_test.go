package monitor

import (
	"testing"
	"time"
)

func TestRateLimiterAllowsUnderLimit(t *testing.T) {
	rl := NewRateLimiter(time.Minute, 3)
	for i := 0; i < 3; i++ {
		if !rl.Allow(8080) {
			t.Fatalf("expected Allow to return true on call %d", i+1)
		}
	}
}

func TestRateLimiterBlocksOverLimit(t *testing.T) {
	rl := NewRateLimiter(time.Minute, 2)
	rl.Allow(9000)
	rl.Allow(9000)
	if rl.Allow(9000) {
		t.Fatal("expected Allow to return false when limit exceeded")
	}
}

func TestRateLimiterIndependentPorts(t *testing.T) {
	rl := NewRateLimiter(time.Minute, 1)
	if !rl.Allow(1111) {
		t.Fatal("expected true for port 1111")
	}
	if !rl.Allow(2222) {
		t.Fatal("expected true for port 2222, independent of 1111")
	}
}

func TestRateLimiterResetsPort(t *testing.T) {
	rl := NewRateLimiter(time.Minute, 1)
	rl.Allow(3000)
	rl.Reset(3000)
	if !rl.Allow(3000) {
		t.Fatal("expected Allow to return true after Reset")
	}
}

func TestRateLimiterWindowExpiry(t *testing.T) {
	rl := NewRateLimiter(50*time.Millisecond, 1)
	rl.Allow(4000)
	time.Sleep(60 * time.Millisecond)
	if !rl.Allow(4000) {
		t.Fatal("expected Allow to return true after window expired")
	}
}
