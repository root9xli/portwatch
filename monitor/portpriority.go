package monitor

import (
	"fmt"
	"sort"
	"sync"
)

// Priority levels for ports.
const (
	PriorityLow      = 1
	PriorityNormal   = 2
	PriorityHigh     = 3
	PriorityCritical = 4
)

// PortPriorityEntry holds priority metadata for a single port.
type PortPriorityEntry struct {
	Port     int
	Priority int
	Reason   string
}

// PortPriority assigns and retrieves priority levels for monitored ports.
type PortPriority struct {
	mu      sync.RWMutex
	entries map[int]PortPriorityEntry
}

// NewPortPriority creates a new PortPriority with an empty registry.
func NewPortPriority() *PortPriority {
	return &PortPriority{
		entries: make(map[int]PortPriorityEntry),
	}
}

// Set assigns a priority level and optional reason to a port.
func (p *PortPriority) Set(port, priority int, reason string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.entries[port] = PortPriorityEntry{Port: port, Priority: priority, Reason: reason}
}

// Get returns the priority entry for a port and whether it was found.
func (p *PortPriority) Get(port int) (PortPriorityEntry, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	e, ok := p.entries[port]
	return e, ok
}

// Remove deletes the priority assignment for a port.
func (p *PortPriority) Remove(port int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.entries, port)
}

// Len returns the number of ports with assigned priorities.
func (p *PortPriority) Len() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.entries)
}

// TopN returns up to n entries sorted by priority descending.
func (p *PortPriority) TopN(n int) []PortPriorityEntry {
	p.mu.RLock()
	defer p.mu.RUnlock()

	result := make([]PortPriorityEntry, 0, len(p.entries))
	for _, e := range p.entries {
		result = append(result, e)
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Priority != result[j].Priority {
			return result[i].Priority > result[j].Priority
		}
		return result[i].Port < result[j].Port
	})
	if n > 0 && n < len(result) {
		return result[:n]
	}
	return result
}

// Summary returns a human-readable summary of all priority assignments.
func (p *PortPriority) Summary() string {
	top := p.TopN(0)
	if len(top) == 0 {
		return "port priorities: none"
	}
	s := fmt.Sprintf("port priorities (%d):\n", len(top))
	for _, e := range top {
		s += fmt.Sprintf("  port=%d priority=%d reason=%q\n", e.Port, e.Priority, e.Reason)
	}
	return s
}
