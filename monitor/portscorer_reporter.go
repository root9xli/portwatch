package monitor

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

// ScoredPort pairs a port with its computed PortScore.
type ScoredPort struct {
	PortScore
	Level string
}

// PortScorerReporter collects scored ports and produces summaries.
type PortScorerReporter struct {
	mu     sync.Mutex
	scorer *PortScorer
	scores []ScoredPort
}

// NewPortScorerReporter creates a reporter backed by the given scorer.
func NewPortScorerReporter(scorer *PortScorer) *PortScorerReporter {
	return &PortScorerReporter{scorer: scorer}
}

// Record scores port and stores the result.
func (r *PortScorerReporter) Record(port int) ScoredPort {
	ps := r.scorer.Score(port)
	sp := ScoredPort{PortScore: ps, Level: r.scorer.Level(ps.Score)}
	r.mu.Lock()
	r.scores = append(r.scores, sp)
	r.mu.Unlock()
	return sp
}

// Top returns up to n highest-scored ports.
func (r *PortScorerReporter) Top(n int) []ScoredPort {
	r.mu.Lock()
	cp := make([]ScoredPort, len(r.scores))
	copy(cp, r.scores)
	r.mu.Unlock()

	sort.Slice(cp, func(i, j int) bool {
		return cp[i].Score > cp[j].Score
	})
	if n > len(cp) {
		n = len(cp)
	}
	return cp[:n]
}

// Summary returns a human-readable summary of the top scored ports.
func (r *PortScorerReporter) Summary() string {
	top := r.Top(5)
	if len(top) == 0 {
		return "port scorer: no ports recorded"
	}
	var sb strings.Builder
	sb.WriteString("port scorer top risks:\n")
	for _, sp := range top {
		sb.WriteString(fmt.Sprintf("  port=%d score=%d level=%s reasons=[%s]\n",
			sp.Port, sp.Score, sp.Level, strings.Join(sp.Reasons, ", ")))
	}
	return sb.String()
}
