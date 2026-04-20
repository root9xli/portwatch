package monitor

import (
	"fmt"
	"sort"
	"strings"
)

// PortGroupReporter produces human-readable summaries of port group activity,
// showing which named groups have seen recent listener changes.
type PortGroupReporter struct {
	groups *PortGroup
	history *History
}

// NewPortGroupReporter creates a reporter that correlates history entries
// with port group membership to produce group-level summaries.
func NewPortGroupReporter(groups *PortGroup, history *History) *PortGroupReporter {
	return &PortGroupReporter{
		groups:  groups,
		history: history,
	}
}

// groupActivity holds aggregated counts for a single named port group.
type groupActivity struct {
	name    string
	added   int
	removed int
	ports   []int
}

// Summary returns a multi-line report of group-level port activity derived
// from the most recent n history entries. Groups with no activity are omitted.
func (r *PortGroupReporter) Summary(n int) string {
	entries := r.history.Last(n)
	if len(entries) == 0 {
		return "port-group report: no recent activity"
	}

	// Aggregate activity per group name.
	activity := map[string]*groupActivity{}

	for _, entry := range entries {
		port := entry.Port
		names := r.groups.GroupsForPort(port)
		if len(names) == 0 {
			names = []string{"ungrouped"}
		}
		for _, name := range names {
			ga, ok := activity[name]
			if !ok {
				ga = &groupActivity{name: name}
				activity[name] = ga
			}
			switch entry.Action {
			case "added":
				ga.added++
			case "removed":
				ga.removed++
			}
			// Track unique ports per group.
			found := false
			for _, p := range ga.ports {
				if p == port {
					found = true
					break
				}
			}
			if !found {
				ga.ports = append(ga.ports, port)
			}
		}
	}

	// Sort group names for deterministic output.
	names := make([]string, 0, len(activity))
	for name := range activity {
		names = append(names, name)
	}
	sort.Strings(names)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("port-group report (%d entries):\n", len(entries)))

	for _, name := range names {
		ga := activity[name]
		sort.Ints(ga.ports)
		portStrs := make([]string, len(ga.ports))
		for i, p := range ga.ports {
			portStrs[i] = fmt.Sprintf("%d", p)
		}
		sb.WriteString(fmt.Sprintf(
			"  %-20s added=%-4d removed=%-4d ports=[%s]\n",
			name, ga.added, ga.removed, strings.Join(portStrs, ","),
		))
	}

	return strings.TrimRight(sb.String(), "\n")
}

// ActiveGroups returns the names of groups that have at least one port
// appearing in the most recent n history entries.
func (r *PortGroupReporter) ActiveGroups(n int) []string {
	entries := r.history.Last(n)
	seen := map[string]struct{}{}
	for _, entry := range entries {
		names := r.groups.GroupsForPort(entry.Port)
		for _, name := range names {
			seen[name] = struct{}{}
		}
	}
	result := make([]string, 0, len(seen))
	for name := range seen {
		result = append(result, name)
	}
	sort.Strings(result)
	return result
}
