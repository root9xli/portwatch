package monitor

import (
	"fmt"
	"sort"
	"strings"
)

// PortGroup represents a named collection of ports that can be referenced
// together in configuration (e.g. "web" -> [80, 443, 8080]).
type PortGroup struct {
	groups map[string][]int
}

// NewPortGroup creates an empty PortGroup registry.
func NewPortGroup() *PortGroup {
	return &PortGroup{
		groups: make(map[string][]int),
	}
}

// Add registers a named group with the given ports. Duplicate ports within a
// group are silently deduplicated.
func (pg *PortGroup) Add(name string, ports []int) {
	seen := make(map[int]struct{}, len(ports))
	uniq := make([]int, 0, len(ports))
	for _, p := range ports {
		if _, ok := seen[p]; !ok {
			seen[p] = struct{}{}
			uniq = append(uniq, p)
		}
	}
	sort.Ints(uniq)
	pg.groups[name] = uniq
}

// Get returns the ports associated with the named group and whether it exists.
func (pg *PortGroup) Get(name string) ([]int, bool) {
	ports, ok := pg.groups[name]
	if !ok {
		return nil, false
	}
	cp := make([]int, len(ports))
	copy(cp, ports)
	return cp, true
}

// Contains reports whether the given port belongs to the named group.
func (pg *PortGroup) Contains(name string, port int) bool {
	ports, ok := pg.groups[name]
	if !ok {
		return false
	}
	for _, p := range ports {
		if p == port {
			return true
		}
	}
	return false
}

// Names returns all registered group names in sorted order.
func (pg *PortGroup) Names() []string {
	names := make([]string, 0, len(pg.groups))
	for n := range pg.groups {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

// Summary returns a human-readable description of all groups.
func (pg *PortGroup) Summary() string {
	if len(pg.groups) == 0 {
		return "no port groups defined"
	}
	var sb strings.Builder
	for _, name := range pg.Names() {
		ports := pg.groups[name]
		strs := make([]string, len(ports))
		for i, p := range ports {
			strs[i] = fmt.Sprintf("%d", p)
		}
		fmt.Fprintf(&sb, "%s=[%s]\n", name, strings.Join(strs, ","))
	}
	return strings.TrimRight(sb.String(), "\n")
}
