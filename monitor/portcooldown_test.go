package monitor

import (
	"testing"
	"time"
)

func newTestCooldown(d time.Duration) *PortCooldown {
	c := NewPortCooldown(d)
	return c
}

func TestPortCooldownNotInCooldownInitially(t *testing.T) {
	c := newTestCooldown(time.Second)
	if c.InCooldown(8080) {
		t.Fatal("expected port not in cooldown initially")
	}
}

func TestPortCooldownInCooldownAfterEnter(t *testing.T) {
	c := newTestCooldown(time.Minute)
	c.Enter(8080)
	if !c.InCooldown(8080) {
		t.Fatal("expected port to be in cooldown after Enter")
	}
}

func TestPortCooldownExpiresAfterDuration(t *testing.T) {
	base := time.Now()
	c := newTestCooldown(5 * time.Second)
	c.now = func() time.Time { return base }
	c.Enter(9000)

	c.now = func() time.Time { return base.Add(6 * time.Second) }
	if c.InCooldown(9000) {
		t.Fatal("expected cooldown to have expired")
	}
}

func TestPortCooldownClearRemovesPort(t *testing.T) {
	c := newTestCooldown(time.Minute)
	c.Enter(3000)
	c.Clear(3000)
	if c.InCooldown(3000) {
		t.Fatal("expected port to be cleared from cooldown")
	}
}

func TestPortCooldownIndependentPorts(t *testing.T) {
	c := newTestCooldown(time.Minute)
	c.Enter(1111)
	if c.InCooldown(2222) {
		t.Fatal("port 2222 should not be in cooldown")
	}
	if !c.InCooldown(1111) {
		t.Fatal("port 1111 should be in cooldown")
	}
}

func TestPortCooldownEvictRemovesExpired(t *testing.T) {
	base := time.Now()
	c := newTestCooldown(2 * time.Second)
	c.now = func() time.Time { return base }
	c.Enter(5000)
	c.Enter(6000)

	c.now = func() time.Time { return base.Add(3 * time.Second) }
	c.Evict()
	if c.Len() != 0 {
		t.Fatalf("expected 0 entries after evict, got %d", c.Len())
	}
}
