package monitor

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

func newTestAlerter(buf *bytes.Buffer) *Alerter {
	return &Alerter{
		logger: log.New(buf, "", 0),
	}
}

func TestAlertSendContainsPort(t *testing.T) {
	var buf bytes.Buffer
	a := newTestAlerter(&buf)
	a.Send(8080, "nginx")
	output := buf.String()
	if !strings.Contains(output, "8080") {
		t.Errorf("expected port 8080 in alert output, got: %s", output)
	}
	if !strings.Contains(output, "nginx") {
		t.Errorf("expected process name in alert output, got: %s", output)
	}
}

func TestAlertSendLevel(t *testing.T) {
	var buf bytes.Buffer
	a := newTestAlerter(&buf)
	a.Send(9090, "unknown")
	if !strings.Contains(buf.String(), string(AlertWarning)) {
		t.Errorf("expected WARNING level in output")
	}
}

func TestAlertSendCriticalLevel(t *testing.T) {
	var buf bytes.Buffer
	a := newTestAlerter(&buf)
	a.SendCritical(22, "sshd")
	output := buf.String()
	if !strings.Contains(output, string(AlertCritical)) {
		t.Errorf("expected CRITICAL level in output")
	}
	if !strings.Contains(output, "22") {
		t.Errorf("expected port 22 in critical alert output")
	}
}

func TestNewAlerter(t *testing.T) {
	a := NewAlerter()
	if a == nil {
		t.Fatal("expected non-nil Alerter")
	}
	if a.logger == nil {
		t.Fatal("expected non-nil logger")
	}
}
