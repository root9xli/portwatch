package monitor

import (
	"encoding/json"
	"os"
	"time"
)

// Config holds the runtime configuration for portwatch.
type Config struct {
	// AllowedPorts lists ports that are expected and should not trigger alerts.
	AllowedPorts []uint16 `json:"allowed_ports"`

	// PollInterval is how often the monitor checks for changes.
	PollInterval time.Duration `json:"poll_interval"`

	// AlertWebhook is an optional HTTP endpoint to POST alerts to.
	AlertWebhook string `json:"alert_webhook"`

	// AlertLevel is the minimum severity level to send: "info", "warn", "critical".
	AlertLevel string `json:"alert_level"`

	// SuppressCooldown is the duration to suppress duplicate alerts per port.
	SuppressCooldown time.Duration `json:"suppress_cooldown"`
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() Config {
	return Config{
		PollInterval:     15 * time.Second,
		AlertLevel:       "warn",
		SuppressCooldown: 5 * time.Minute,
	}
}

// LoadConfig reads a JSON config file from path and returns a Config.
// Missing fields fall back to defaults.
func LoadConfig(path string) (Config, error) {
	cfg := DefaultConfig()
	f, err := os.Open(path)
	if err != nil {
		return cfg, err
	}
	defer f.Close()
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
