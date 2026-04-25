package monitor

import (
	"sync"
	"time"
)

// PortExpiry tracks ports that are expected to disappear by a deadline
// and flags those that outlive their expected lifetime.
type PortExpiry struct {
	mu      sync.Mutex
	entries map[int]time.Time
	now     func() time.Time
}

// NewPortExpiry creates a new PortExpiry tracker.
func NewPortExpiry() *PortExpiry {
	return &PortExpiry{
		entries: make(map[int]time.Time),
		now:     time.Now,
	}
}

// Register marks a port as expected to be gone by deadline.
func (pe *PortExpiry) Register(port int, deadline time.Time) {
	pe.mu.Lock()
	defer pe.mu.Unlock()
	pe.entries[port] = deadline
}

// Forget removes a port from expiry tracking.
func (pe *PortExpiry) Forget(port int) {
	pe.mu.Lock()
	defer pe.mu.Unlock()
	delete(pe.entries, port)
}

// IsOverdue returns true if the port is registered and its deadline has passed.
func (pe *PortExpiry) IsOverdue(port int) bool {
	pe.mu.Lock()
	defer pe.mu.Unlock()
	deadline, ok := pe.entries[port]
	if !ok {
		return false
	}
	return pe.now().After(deadline)
}

// Overdue returns all ports whose deadlines have passed.
func (pe *PortExpiry) Overdue() []int {
	pe.mu.Lock()
	defer pe.mu.Unlock()
	now := pe.now()
	var result []int
	for port, deadline := range pe.entries {
		if now.After(deadline) {
			result = append(result, port)
		}
	}
	return result
}

// Evict removes all ports whose deadlines have passed and returns them.
func (pe *PortExpiry) Evict() []int {
	pe.mu.Lock()
	defer pe.mu.Unlock()
	now := pe.now()
	var evicted []int
	for port, deadline := range pe.entries {
		if now.After(deadline) {
			evicted = append(evicted, port)
			delete(pe.entries, port)
		}
	}
	return evicted
}

// Len returns the number of tracked ports.
func (pe *PortExpiry) Len() int {
	pe.mu.Lock()
	defer pe.mu.Unlock()
	return len(pe.entries)
}
