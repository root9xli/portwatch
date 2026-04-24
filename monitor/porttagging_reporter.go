package monitor

import (
	"fmt"
	"sort"
	"strings"
)

// PortTaggingReporter generates reports from a PortTagging instance.
type PortTaggingReporter struct {
	tagging *PortTagging
}

// NewPortTaggingReporter creates a reporter backed by the given PortTagging.
func NewPortTaggingReporter(pt *PortTagging) *PortTaggingReporter {
	return &PortTaggingReporter{tagging: pt}
}

// TaggedPorts returns a sorted list of ports that have at least one tag.
func (r *PortTaggingReporter) TaggedPorts() []int {
	r.tagging.mu.RLock()
	defer r.tagging.mu.RUnlock()

	ports := make([]int, 0, len(r.tagging.tags))
	for p := range r.tagging.tags {
		ports = append(ports, p)
	}
	sort.Ints(ports)
	return ports
}

// PortsWithTag returns all ports that carry the given tag, sorted.
func (r *PortTaggingReporter) PortsWithTag(tag string) []int {
	r.tagging.mu.RLock()
	defer r.tagging.mu.RUnlock()

	var result []int
	for port, tags := range r.tagging.tags {
		for _, t := range tags {
			if t == tag {
				result = append(result, port)
				break
			}
		}
	}
	sort.Ints(result)
	return result
}

// Summary returns a formatted string listing all ports and their tags.
func (r *PortTaggingReporter) Summary() string {
	ports := r.TaggedPorts()
	if len(ports) == 0 {
		return "port tagging report: no tagged ports"
	}
	var sb strings.Builder
	sb.WriteString("port tagging report:\n")
	for _, p := range ports {
		tags := r.tagging.Get(p)
		fmt.Fprintf(&sb, "  port %-6d tags: %s\n", p, strings.Join(tags, ", "))
	}
	return strings.TrimRight(sb.String(), "\n")
}
