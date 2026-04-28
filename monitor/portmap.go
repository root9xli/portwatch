package monitor

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

// PortMapEntry holds metadata about a mapped port relationship.
type PortMapEntry struct {
	Port    int
	Alias   string
	Comment string
}

// PortMap tracks named aliases and comments for ports, useful for
// correlating discovered ports with known services or deployments.
type PortMap struct {
	mu      sync.RWMutex
	entries map[int]PortMapEntry
}

// NewPortMap returns an empty PortMap.
func NewPortMap() *PortMap {
	return &PortMap{
		entries: make(map[int]PortMapEntry),
	}
}

// Set registers or updates an alias and comment for a port.
func (pm *PortMap) Set(port int, alias, comment string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.entries[port] = PortMapEntry{Port: port, Alias: alias, Comment: comment}
}

// Get returns the entry for a port and whether it was found.
func (pm *PortMap) Get(port int) (PortMapEntry, bool) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	e, ok := pm.entries[port]
	return e, ok
}

// Remove deletes the mapping for a port.
func (pm *PortMap) Remove(port int) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	delete(pm.entries, port)
}

// Len returns the number of mapped ports.
func (pm *PortMap) Len() int {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return len(pm.entries)
}

// All returns a sorted slice of all PortMapEntry values.
func (pm *PortMap) All() []PortMapEntry {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	ports := make([]int, 0, len(pm.entries))
	for p := range pm.entries {
		ports = append(ports, p)
	}
	sort.Ints(ports)
	result := make([]PortMapEntry, 0, len(ports))
	for _, p := range ports {
		result = append(result, pm.entries[p])
	}
	return result
}

// Summary returns a human-readable listing of all mapped ports.
func (pm *PortMap) Summary() string {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	if len(pm.entries) == 0 {
		return "portmap: no entries"
	}
	ports := make([]int, 0, len(pm.entries))
	for p := range pm.entries {
		ports = append(ports, p)
	}
	sort.Ints(ports)
	var sb strings.Builder
	for _, p := range ports {
		e := pm.entries[p]
		sb.WriteString(fmt.Sprintf("  %-6d %-20s %s\n", e.Port, e.Alias, e.Comment))
	}
	return sb.String()
}
