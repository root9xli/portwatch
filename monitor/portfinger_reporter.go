package monitor

import (
	"fmt"
	"sort"
	"strings"
)

// PortFingerReporter produces human-readable summaries of fingerprint results.
type PortFingerReporter struct {
	finger *PortFinger
}

// NewPortFingerReporter wraps a PortFinger for reporting.
func NewPortFingerReporter(pf *PortFinger) *PortFingerReporter {
	return &PortFingerReporter{finger: pf}
}

// Summary returns a multi-line string describing all probed ports.
func (r *PortFingerReporter) Summary() string {
	r.finger.mu.Lock()
	results := make([]FingerResult, 0, len(r.finger.results))
	for _, fr := range r.finger.results {
		results = append(results, fr)
	}
	r.finger.mu.Unlock()

	if len(results) == 0 {
		return "portfinger: no probes recorded"
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Port < results[j].Port
	})

	var sb strings.Builder
	sb.WriteString("portfinger summary:\n")
	for _, fr := range results {
		status := "unreachable"
		if fr.Reached {
			status = "reached"
		}
		banner := fr.Banner
		if banner == "" {
			banner = "(no banner)"
		}
		sb.WriteString(fmt.Sprintf("  port %d: %s banner=%q probed=%s\n",
			fr.Port, status, banner, fr.Probed.Format("15:04:05")))
	}
	return strings.TrimRight(sb.String(), "\n")
}

// HasBanner returns true if the given port returned a non-empty banner.
func (r *PortFingerReporter) HasBanner(port int) bool {
	fr, ok := r.finger.Get(port)
	return ok && fr.Banner != ""
}

// BannerFor returns the banner string for a port, or empty string.
func (r *PortFingerReporter) BannerFor(port int) string {
	fr, ok := r.finger.Get(port)
	if !ok {
		return ""
	}
	return fr.Banner
}
