package monitor

import (
	"fmt"
	"log"
	"os"
	"time"
)

// AlertLevel represents the severity of an alert.
type AlertLevel string

const (
	AlertWarning  AlertLevel = "WARNING"
	AlertCritical AlertLevel = "CRITICAL"
)

// Alert represents a port alert event.
type Alert struct {
	Timestamp time.Time
	Level     AlertLevel
	Port      int
	Process   string
	Message   string
}

// Alerter handles sending alerts to configured outputs.
type Alerter struct {
	logger *log.Logger
}

// NewAlerter creates a new Alerter writing to the given output.
func NewAlerter() *Alerter {
	return &Alerter{
		logger: log.New(os.Stdout, "", 0),
	}
}

// Send emits an alert for an unexpected port listener.
func (a *Alerter) Send(port int, process string) {
	alert := Alert{
		Timestamp: time.Now().UTC(),
		Level:     AlertWarning,
		Port:      port,
		Process:   process,
		Message:   fmt.Sprintf("unexpected listener on port %d (process: %s)", port, process),
	}
	a.logger.Printf("[%s] %s %s", alert.Level, alert.Timestamp.Format(time.RFC3339), alert.Message)
}

// SendCritical emits a critical alert.
func (a *Alerter) SendCritical(port int, process string) {
	alert := Alert{
		Timestamp: time.Now().UTC(),
		Level:     AlertCritical,
		Port:      port,
		Process:   process,
		Message:   fmt.Sprintf("critical: unexpected listener on port %d (process: %s)", port, process),
	}
	a.logger.Printf("[%s] %s %s", alert.Level, alert.Timestamp.Format(time.RFC3339), alert.Message)
}
