package monitor

import "sync"

// PortScore holds a risk score and contributing factors for a port.
type PortScore struct {
	Port    int
	Score   int
	Reasons []string
}

// PortScorer computes a composite risk score for newly seen ports.
type PortScorer struct {
	mu       sync.Mutex
	labeler  *PortLabeler
	anomaly  *AnomalyDetector
	deadport *DeadPort
	baseline *Baseline
}

// NewPortScorer creates a PortScorer using existing subsystems.
func NewPortScorer(l *PortLabeler, a *AnomalyDetector, d *DeadPort, b *Baseline) *PortScorer {
	return &PortScorer{
		labeler:  l,
		anomaly:  a,
		deadport: d,
		baseline: b,
	}
}

// Score computes a risk score for a port and returns a PortScore.
func (ps *PortScorer) Score(port int) PortScore {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	result := PortScore{Port: port}

	// Unknown port label adds risk.
	if ps.labeler != nil {
		label := ps.labeler.Label(port)
		if label == "unknown" {
			result.Score += 30
			result.Reasons = append(result.Reasons, "unlabeled port")
		}
	}

	// Anomalous frequency adds risk.
	if ps.anomaly != nil && ps.anomaly.IsAnomalous(port) {
		result.Score += 40
		result.Reasons = append(result.Reasons, "anomalous frequency")
	}

	// Port recycled quickly after disappearing adds risk.
	if ps.deadport != nil && ps.deadport.IsRecycled(port) {
		result.Score += 20
		result.Reasons = append(result.Reasons, "rapid port recycle")
	}

	// Port not in baseline adds risk.
	if ps.baseline != nil && !ps.baseline.Contains(port) {
		result.Score += 10
		result.Reasons = append(result.Reasons, "not in baseline")
	}

	return result
}

// Level returns a human-readable risk level for a score.
func (ps *PortScorer) Level(score int) string {
	switch {
	case score >= 70:
		return "critical"
	case score >= 40:
		return "warning"
	default:
		return "info"
	}
}
