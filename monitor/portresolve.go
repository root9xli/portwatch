package monitor

import (
	"fmt"
	"net"
	"sync"
	"time"
)

// PortResolve performs reverse DNS lookups for observed ports and caches results.
type PortResolve struct {
	mu      sync.RWMutex
	cache   map[int]resolveEntry
	ttl     time.Duration
	nowFunc func() time.Time
}

type resolveEntry struct {
	hostnames []string
	resolvedAt time.Time
}

// NewPortResolve creates a PortResolve with the given cache TTL.
func NewPortResolve(ttl time.Duration) *PortResolve {
	return &PortResolve{
		cache:   make(map[int]resolveEntry),
		ttl:     ttl,
		nowFunc: time.Now,
	}
}

// Resolve performs a reverse DNS lookup for the local address on the given port.
// Results are cached for the configured TTL.
func (r *PortResolve) Resolve(port int) ([]string, error) {
	r.mu.RLock()
	entry, ok := r.cache[port]
	r.mu.RUnlock()

	if ok && r.nowFunc().Before(entry.resolvedAt.Add(r.ttl)) {
		return entry.hostnames, nil
	}

	addr := fmt.Sprintf("127.0.0.1:%d", port)
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}

	hostnames, err := net.LookupAddr(host)
	if err != nil {
		hostnames = []string{}
	}

	r.mu.Lock()
	r.cache[port] = resolveEntry{
		hostnames:  hostnames,
		resolvedAt: r.nowFunc(),
	}
	r.mu.Unlock()

	return hostnames, nil
}

// Get returns cached hostnames for a port without triggering a new lookup.
// Returns nil if the entry is absent or expired.
func (r *PortResolve) Get(port int) []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	entry, ok := r.cache[port]
	if !ok || r.nowFunc().After(entry.resolvedAt.Add(r.ttl)) {
		return nil
	}
	return entry.hostnames
}

// Evict removes expired cache entries.
func (r *PortResolve) Evict() {
	now := r.nowFunc()
	r.mu.Lock()
	defer r.mu.Unlock()
	for port, entry := range r.cache {
		if now.After(entry.resolvedAt.Add(r.ttl)) {
			delete(r.cache, port)
		}
	}
}

// Len returns the number of cached entries.
func (r *PortResolve) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.cache)
}
