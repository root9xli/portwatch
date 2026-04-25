package monitor

import (
	"strings"
	"testing"
)

func newTestClassifyReporter() (*PortClassifier, *PortClassifyReporter) {
	c := NewPortClassifier()
	return c, NewPortClassifyReporter(c)
}

func TestPortClassifyReporterEmpty(t *testing.T) {
	_, r := newTestClassifyReporter()
	summary := r.Report(nil)
	if !strings.Contains(summary, "none") {
		t.Errorf("expected 'none' in empty report, got: %s", summary)
	}
}

func TestPortClassifyReporterContainsClass(t *testing.T) {
	_, r := newTestClassifyReporter()
	summary := r.Report([]int{22})
	if !strings.Contains(summary, "system") {
		t.Errorf("expected 'system' in report, got: %s", summary)
	}
}

func TestPortClassifyReporterContainsPort(t *testing.T) {
	_, r := newTestClassifyReporter()
	summary := r.Report([]int{8080})
	if !strings.Contains(summary, "8080") {
		t.Errorf("expected port 8080 in report, got: %s", summary)
	}
}

func TestPortClassifyReporterMultipleClasses(t *testing.T) {
	_, r := newTestClassifyReporter()
	summary := r.Report([]int{80, 3000, 55000})
	for _, cls := range []string{"system", "registered", "ephemeral"} {
		if !strings.Contains(summary, cls) {
			t.Errorf("expected class %q in report, got: %s", cls, summary)
		}
	}
}

func TestPortClassifyReporterCountByClass(t *testing.T) {
	_, r := newTestClassifyReporter()
	counts := r.CountByClass([]int{22, 80, 8080, 9090, 55000})
	if counts[ClassSystem] != 2 {
		t.Errorf("expected 2 system ports, got %d", counts[ClassSystem])
	}
	if counts[ClassRegistered] != 2 {
		t.Errorf("expected 2 registered ports, got %d", counts[ClassRegistered])
	}
	if counts[ClassEphemeral] != 1 {
		t.Errorf("expected 1 ephemeral port, got %d", counts[ClassEphemeral])
	}
}

func TestPortClassifyReporterCustomOverrideInReport(t *testing.T) {
	c, r := newTestClassifyReporter()
	c.SetOverride(9090, ClassCustom)
	summary := r.Report([]int{9090})
	if !strings.Contains(summary, "custom") {
		t.Errorf("expected 'custom' in report, got: %s", summary)
	}
}
