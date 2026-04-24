package monitor

import (
	"fmt"
	"strings"
	"sync"
)

// PortTagging assigns and retrieves user-defined tags for ports.
type PortTagging struct {
	mu   sync.RWMutex
	tags map[int][]string
}

// NewPortTagging creates a new PortTagging instance.
func NewPortTagging() *PortTagging {
	return &PortTagging{
		tags: make(map[int][]string),
	}
}

// Tag adds one or more tags to a port. Duplicate tags are ignored.
func (pt *PortTagging) Tag(port int, newTags ...string) {
	pt.mu.Lock()
	defer pt.mu.Unlock()

	existing := pt.tags[port]
	set := make(map[string]struct{}, len(existing))
	for _, t := range existing {
		set[t] = struct{}{}
	}
	for _, t := range newTags {
		t = strings.TrimSpace(t)
		if t == "" {
			continue
		}
		if _, ok := set[t]; !ok {
			existing = append(existing, t)
			set[t] = struct{}{}
		}
	}
	pt.tags[port] = existing
}

// Untag removes a specific tag from a port.
func (pt *PortTagging) Untag(port int, tag string) {
	pt.mu.Lock()
	defer pt.mu.Unlock()

	current := pt.tags[port]
	updated := current[:0]
	for _, t := range current {
		if t != tag {
			updated = append(updated, t)
		}
	}
	if len(updated) == 0 {
		delete(pt.tags, port)
	} else {
		pt.tags[port] = updated
	}
}

// Get returns all tags for a port.
func (pt *PortTagging) Get(port int) []string {
	pt.mu.RLock()
	defer pt.mu.RUnlock()

	copy := append([]string(nil), pt.tags[port]...)
	return copy
}

// HasTag reports whether a port has a specific tag.
func (pt *PortTagging) HasTag(port int, tag string) bool {
	pt.mu.RLock()
	defer pt.mu.RUnlock()

	for _, t := range pt.tags[port] {
		if t == tag {
			return true
		}
	}
	return false
}

// Summary returns a human-readable summary of all tagged ports.
func (pt *PortTagging) Summary() string {
	pt.mu.RLock()
	defer pt.mu.RUnlock()

	if len(pt.tags) == 0 {
		return "port tags: none"
	}
	var sb strings.Builder
	sb.WriteString("port tags:\n")
	for port, tags := range pt.tags {
		fmt.Fprintf(&sb, "  :%d -> [%s]\n", port, strings.Join(tags, ", "))
	}
	return strings.TrimRight(sb.String(), "\n")
}
