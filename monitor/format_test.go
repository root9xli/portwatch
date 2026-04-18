package monitor

import (
	"strings"
	"testing"
	"time"
)

func baseMsg() AlertMessage {
	return AlertMessage{
		Level:     LevelWarning,
		Port:      8080,
		Action:    "added",
		Timestamp: time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
		Hostname:  "testhost",
	}
}

func TestFormatContainsPort(t *testing.T) {
	msg := baseMsg()
	out := msg.Format()
	if !strings.Contains(out, "8080") {
		t.Errorf("expected port in output, got: %s", out)
	}
}

func TestFormatContainsLevel(t *testing.T) {
	msg := baseMsg()
	out := msg.Format()
	if !strings.Contains(out, "WARNING") {
		t.Errorf("expected level in output, got: %s", out)
	}
}

func TestFormatContainsAction(t *testing.T) {
	msg := baseMsg()
	out := msg.Format()
	if !strings.Contains(out, "added") {
		t.Errorf("expected action in output, got: %s", out)
	}
}

func TestFormatFallbackHostname(t *testing.T) {
	msg := baseMsg()
	msg.Hostname = ""
	out := msg.Format()
	if !strings.Contains(out, "unknown") {
		t.Errorf("expected fallback hostname, got: %s", out)
	}
}

func TestFormatMultipleJoinsLines(t *testing.T) {
	msgs := []AlertMessage{baseMsg(), baseMsg()}
	msgs[1].Port = 9090
	out := FormatMultiple(msgs)
	lines := strings.Split(out, "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 lines, got %d", len(lines))
	}
	if !strings.Contains(out, "9090") {
		t.Errorf("expected second port in output")
	}
}
