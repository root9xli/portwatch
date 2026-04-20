package monitor

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

// PortAge tracks how long each port has been observed as open.
type PortAge struct {
	mu      sync.Mutex
	firstSeen map[int]time.Time
	clock     func() time.Time
}

// NewPortAge creates a new PortAge tracker.
func NewPortAge() *PortAge {
	return &PortAge{
		firstSeen: make(map[int]time.Time),
		clock:     time.Now,
	}
}

// Observe records the first time a port is seen if not already tracked.
func (p *PortAge) Observe(port int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, ok := p.firstSeen[port]; !ok {
		p.firstSeen[port] = p.clock()
	}
}

// Age returns how long the port has been observed, and whether it is known.
func (p *PortAge) Age(port int) (time.Duration, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	t, ok := p.firstSeen[port]
	if !ok {
		return 0, false
	}
	return p.clock().Sub(t), true
}

// Forget removes a port from tracking (e.g. when it closes).
func (p *PortAge) Forget(port int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.firstSeen, port)
}

// Summary returns a human-readable list of ports and their ages, sorted by port.
func (p *PortAge) Summary() string {
	p.mu.Lock()
	defer p.mu.Unlock()
	if len(p.firstSeen) == 0 {
		return "no ports tracked"
	}
	ports := make([]int, 0, len(p.firstSeen))
	for port := range p.firstSeen {
		ports = append(ports, port)
	}
	sort.Ints(ports)
	out := ""
	for _, port := range ports {
		age := p.clock().Sub(p.firstSeen[port]).Truncate(time.Second)
		out += fmt.Sprintf("port %d: open for %s\n", port, age)
	}
	return out
}

// Len returns the number of currently tracked ports.
func (p *PortAge) Len() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.firstSeen)
}
