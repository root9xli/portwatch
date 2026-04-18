package monitor

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds the portwatch configuration.
type Config struct {
	AllowedPorts []int    `yaml:"allowed_ports"`
	AlertCommand string   `yaml:"alert_command"`
	LogFile      string   `yaml:"log_file"`
	IgnorePIDs   []int    `yaml:"ignore_pids"`
}

// LoadConfig reads and parses a YAML config file.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// IsAllowed returns true if the given port is in the allowed list.
func (c *Config) IsAllowed(port int) bool {
	for _, p := range c.AllowedPorts {
		if p == port {
			return true
		}
	}
	return false
}
