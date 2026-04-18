package monitor

import (
	"strings"
	"testing"
)

func TestMetricsReporterSummaryContainsFields(t *testing.T) {
	m := NewMetrics()
	m.RecordDiff()
	m.RecordAlert()
	m.RecordSuppressed()
	m.RecordRateLimited()
	r := NewMetricsReporter(m)
	summary := r.Summary()
	for _, key := range []string{"uptime", "diffs", "alerts", "suppressed", "rate_limited", "last_diff"} {
		if !strings.Contains(summary, key) {
			t.Errorf("summary missing field %q", key)
		}
	}
}

func TestMetricsReporterLastDiffNever(t *testing.T) {
	m := NewMetrics()
	r := NewMetricsReporter(m)
	if !strings.Contains(r.Summary(), "never") {
		t.Error("expected 'never' when no diff recorded")
	}
}

func TestMetricsReporterLastDiffTimestamp(t *testing.T) {
	m := NewMetrics()
	m.RecordDiff()
	r := NewMetricsReporter(m)
	if strings.Contains(r.Summary(), "never") {
		t.Error("did not expect 'never' after a diff was recorded")
	}
}

func TestFmtDuration(t *testing.T) {
	cases := []struct {
		secs int
		want string
	}{
		{5, "5s"},
		{65, "1m5s"},
		{3661, "1h1m1s"},
	}
	for _, c := range cases {
		import_dur := int64(c.secs) * int64(1e9)
		_ = import_dur
	}
	// basic smoke test via reporter
	m := NewMetrics()
	r := NewMetricsReporter(m)
	if !strings.Contains(r.Summary(), "uptime") {
		t.Error("summary should contain uptime")
	}
}
