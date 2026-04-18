package monitor

import (
	"fmt"
	"strings"
	"time"
)

// AlertLevel represents the severity of an alert.
type AlertLevel string

const (
	LevelInfo     AlertLevel = "INFO"
	LevelWarning  AlertLevel = "WARNING"
	LevelCritical AlertLevel = "CRITICAL"
)

// AlertMessage holds structured data for a formatted alert.
type AlertMessage struct {
	Level     AlertLevel
	Port      int
	Action    string // "added" or "removed"
	Timestamp time.Time
	Hostname  string
}

// Format returns a human-readable alert string.
func (a AlertMessage) Format() string {
	ts := a.Timestamp.UTC().Format(time.RFC3339)
	host := a.Hostname
	if host == "" {
		host = "unknown"
	}
	return fmt.Sprintf(
		"[%s] %s | host=%s port=%d action=%s",
		ts, a.Level, host, a.Port, a.Action,
	)
}

// FormatMultiple returns a combined string for multiple alerts.
func FormatMultiple(msgs []AlertMessage) string {
	lines := make([]string, len(msgs))
	for i, m := range msgs {
		lines[i] = m.Format()
	}
	return strings.Join(lines, "\n")
}
