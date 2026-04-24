package monitor

import "sync"

// PortChainEntry records a sequence of ports observed together in the same
// diff cycle, capturing co-occurrence relationships between listeners.
type PortChainEntry struct {
	Ports []int
	Count int
}

// PortChain tracks which ports tend to appear together in the same diff
// snapshot, helping identify coordinated or related listener activity.
type PortChain struct {
	mu      sync.Mutex
	chains  map[string]*PortChainEntry
	maxKeys int
}

// NewPortChain returns a PortChain that retains at most maxKeys distinct
// co-occurrence groups before evicting the least-seen entry.
func NewPortChain(maxKeys int) *PortChain {
	if maxKeys <= 0 {
		maxKeys = 256
	}
	return &PortChain{
		chains:  make(map[string]*PortChainEntry),
		maxKeys: maxKeys,
	}
}

// Observe records a set of ports that appeared together in one diff cycle.
// Ports are sorted and joined into a canonical key before storage.
func (pc *PortChain) Observe(ports []int) {
	if len(ports) < 2 {
		return
	}
	sorted := sortedCopy(ports)
	key := intsKey(sorted)

	pc.mu.Lock()
	defer pc.mu.Unlock()

	if e, ok := pc.chains[key]; ok {
		e.Count++
		return
	}
	if len(pc.chains) >= pc.maxKeys {
		pc.evictLocked()
	}
	pc.chains[key] = &PortChainEntry{Ports: sorted, Count: 1}
}

// Get returns the entry for the given canonical port set, if present.
func (pc *PortChain) Get(ports []int) (*PortChainEntry, bool) {
	key := intsKey(sortedCopy(ports))
	pc.mu.Lock()
	defer pc.mu.Unlock()
	e, ok := pc.chains[key]
	return e, ok
}

// Len returns the number of distinct co-occurrence groups recorded.
func (pc *PortChain) Len() int {
	pc.mu.Lock()
	defer pc.mu.Unlock()
	return len(pc.chains)
}

// evictLocked removes the entry with the lowest count. Caller must hold mu.
func (pc *PortChain) evictLocked() {
	var minKey string
	minCount := -1
	for k, e := range pc.chains {
		if minCount < 0 || e.Count < minCount {
			minCount = e.Count
			minKey = k
		}
	}
	if minKey != "" {
		delete(pc.chains, minKey)
	}
}

// sortedCopy returns a sorted copy of the given int slice.
func sortedCopy(ports []int) []int {
	cp := make([]int, len(ports))
	copy(cp, ports)
	for i := 1; i < len(cp); i++ {
		for j := i; j > 0 && cp[j] < cp[j-1]; j-- {
			cp[j], cp[j-1] = cp[j-1], cp[j]
		}
	}
	return cp
}

// intsKey encodes a sorted int slice into a string key.
func intsKey(ports []int) string {
	b := make([]byte, 0, len(ports)*6)
	for i, p := range ports {
		if i > 0 {
			b = append(b, ',')
		}
		b = appendInt(b, p)
	}
	return string(b)
}

func appendInt(b []byte, n int) []byte {
	if n == 0 {
		return append(b, '0')
	}
	var tmp [10]byte
	pos := 10
	for n > 0 {
		pos--
		tmp[pos] = byte('0' + n%10)
		n /= 10
	}
	return append(b, tmp[pos:]...)
}
