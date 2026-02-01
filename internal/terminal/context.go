package terminal

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const (
	ScrollbackLines = 500
)

// InTmux returns true if running inside a tmux session
func InTmux() bool {
	return os.Getenv("TMUX") != ""
}

// CaptureContext captures the terminal scrollback via tmux
// Returns the context string and a warning message if not in tmux
func CaptureContext() (context string, warning string, err error) {
	if !InTmux() {
		return "", "Warning: Not running in tmux. Terminal context capture is unavailable.", nil
	}

	// Capture scrollback from tmux
	// -p: print to stdout
	// -S -N: start from N lines back (negative = scrollback)
	cmd := exec.Command("tmux", "capture-pane", "-p", "-S", fmt.Sprintf("-%d", ScrollbackLines))
	output, err := cmd.Output()
	if err != nil {
		return "", "", fmt.Errorf("failed to capture tmux pane: %w", err)
	}

	context = strings.TrimSpace(string(output))
	return context, "", nil
}
