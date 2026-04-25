package monitor

import (
	"sync"
	"time"
)

// PortCooldown tracks ports that are in a cooldown period after being removed,
// preventing re-alert noise when a port flaps back quickly.
type PortCooldown struct {
	mu       sync.Mutex
	entries  map[int]time.Time
	cooldown time.Duration
	now      func() time.Time
}

// NewPortCooldown creates a PortCooldown with the given cooldown duration.
func NewPortCooldown(cooldown time.Duration) *PortCooldown {
	return &PortCooldown{
		entries:  make(map[int]time.Time),
		cooldown: cooldown,
		now:      time.Now,
	}
}

// Enter marks a port as entering cooldown starting now.
func (c *PortCooldown) Enter(port int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[port] = c.now().Add(c.cooldown)
}

// InCooldown reports whether the port is currently in cooldown.
func (c *PortCooldown) InCooldown(port int) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	expiry, ok := c.entries[port]
	if !ok {
		return false
	}
	if c.now().After(expiry) {
		delete(c.entries, port)
		return false
	}
	return true
}

// Clear removes a port from cooldown immediately.
func (c *PortCooldown) Clear(port int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, port)
}

// Evict removes all expired cooldown entries.
func (c *PortCooldown) Evict() {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := c.now()
	for port, expiry := range c.entries {
		if now.After(expiry) {
			delete(c.entries, port)
		}
	}
}

// Len returns the number of ports currently tracked (including expired).
func (c *PortCooldown) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.entries)
}
