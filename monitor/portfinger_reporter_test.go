package monitor

import (
	"strings"
	"testing"
	"time"
)

func newTestFingerReporter(t *testing.T) (*PortFinger, *PortFingerReporter) {
	t.Helper()
	pf := NewPortFinger(time.Second)
	return pf, NewPortFingerReporter(pf)
}

func TestPortFingerReporterEmptySummary(t *testing.T) {
	_, rep := newTestFingerReporter(t)
	s := rep.Summary()
	if !strings.Contains(s, "no probes") {
		t.Fatalf("expected empty summary, got: %s", s)
	}
}

func TestPortFingerReporterSummaryContainsPort(t *testing.T) {
	port := startBannerListener(t, "SMTP")
	pf, rep := newTestFingerReporter(t)
	pf.Probe(port)
	s := rep.Summary()
	if !strings.Contains(s, "SMTP") {
		t.Fatalf("expected banner in summary, got: %s", s)
	}
}

func TestPortFingerReporterHasBanner(t *testing.T) {
	port := startBannerListener(t, "FTP-READY")
	pf, rep := newTestFingerReporter(t)
	pf.Probe(port)
	if !rep.HasBanner(port) {
		t.Fatal("expected HasBanner true")
	}
}

func TestPortFingerReporterNoBannerUnreachable(t *testing.T) {
	_, rep := newTestFingerReporter(t)
	if rep.HasBanner(9) {
		t.Skip("port 9 unexpectedly has banner")
	}
}

func TestPortFingerReporterBannerFor(t *testing.T) {
	port := startBannerListener(t, "HTTP/1.0")
	pf, rep := newTestFingerReporter(t)
	pf.Probe(port)
	b := rep.BannerFor(port)
	if b != "HTTP/1.0" {
		t.Fatalf("expected HTTP/1.0, got %q", b)
	}
}

func TestPortFingerReporterBannerForMissing(t *testing.T) {
	_, rep := newTestFingerReporter(t)
	if rep.BannerFor(65000) != "" {
		t.Fatal("expected empty banner for unprobed port")
	}
}
