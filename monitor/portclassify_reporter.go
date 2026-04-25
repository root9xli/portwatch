package monitor

import (
	"fmt"
	"sort"
	"strings"
)

// PortClassifyReporter summarises observed ports grouped by their class.
type PortClassifyReporter struct {
	classifier *PortClassifier
}

// NewPortClassifyReporter returns a reporter backed by the given classifier.
func NewPortClassifyReporter(c *PortClassifier) *PortClassifyReporter {
	return &PortClassifyReporter{classifier: c}
}

// Report returns a human-readable summary of ports grouped by class.
func (r *PortClassifyReporter) Report(ports []int) string {
	if len(ports) == 0 {
		return "port classes: none"
	}

	groups := make(map[PortClass][]int)
	for _, p := range ports {
		cls := r.classifier.Classify(p)
		groups[cls] = append(groups[cls], p)
	}

	// Sort classes for deterministic output.
	classes := make([]string, 0, len(groups))
	for cls := range groups {
		classes = append(classes, string(cls))
	}
	sort.Strings(classes)

	var sb strings.Builder
	sb.WriteString("port classes:\n")
	for _, cls := range classes {
		ps := groups[PortClass(cls)]
		sort.Ints(ps)
		strs := make([]string, len(ps))
		for i, p := range ps {
			strs[i] = fmt.Sprintf("%d", p)
		}
		fmt.Fprintf(&sb, "  %s: [%s]\n", cls, strings.Join(strs, ", "))
	}
	return strings.TrimRight(sb.String(), "\n")
}

// CountByClass returns a map of class → count for the given ports.
func (r *PortClassifyReporter) CountByClass(ports []int) map[PortClass]int {
	out := make(map[PortClass]int)
	for _, p := range ports {
		out[r.classifier.Classify(p)]++
	}
	return out
}
