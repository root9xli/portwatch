package monitor

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// Notifier sends alerts via a configured external command.
type Notifier struct {
	cmd    string
	args   []string
	logger *log.Logger
}

// NewNotifier creates a Notifier that invokes cmd with args for each alert.
// The message is appended as the final argument.
func NewNotifier(cmd string, args []string, logger *log.Logger) *Notifier {
	if logger == nil {
		logger = log.Default()
	}
	return &Notifier{cmd: cmd, args: args, logger: logger}
}

// Notify executes the configured command with the provided message.
func (n *Notifier) Notify(message string) error {
	if n.cmd == "" {
		return nil
	}
	fullArgs := append(append([]string{}, n.args...), message)
	cmd := exec.Command(n.cmd, fullArgs...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("notifier: command failed: %w (output: %s)", err, strings.TrimSpace(string(out)))
	}
	n.logger.Printf("notifier: sent alert via %q", n.cmd)
	return nil
}

// NotifyMultiple sends a notification for each message in msgs.
func (n *Notifier) NotifyMultiple(msgs []string) []error {
	var errs []error
	for _, m := range msgs {
		if err := n.Notify(m); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}
