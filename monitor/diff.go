package monitor

// DiffResult holds ports that appeared or disappeared between two snapshots.
type DiffResult struct {
	Added   []PortEntry
	Removed []PortEntry
}

// Diff compares two snapshots and returns new and removed port entries.
func Diff(prev, curr *Snapshot) DiffResult {
	prevMap := indexPorts(prev)
	currMap := indexPorts(curr)

	var added, removed []PortEntry

	for key, entry := range currMap {
		if _, exists := prevMap[key]; !exists {
			added = append(added, entry)
		}
	}

	for key, entry := range prevMap {
		if _, exists := currMap[key]; !exists {
			removed = append(removed, entry)
		}
	}

	return DiffResult{Added: added, Removed: removed}
}

// indexPorts builds a map keyed by "protocol:port" for quick lookup.
func indexPorts(s *Snapshot) map[string]PortEntry {
	m := make(map[string]PortEntry, len(s.Ports))
	for _, p := range s.Ports {
		key := p.Protocol + ":" + itoa(p.Port)
		m[key] = p
	}
	return m
}

// itoa is a minimal int-to-string helper to avoid importing strconv.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	buf := [20]byte{}
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte(n%10) + '0'
		n /= 10
	}
	return string(buf[pos:])
}
