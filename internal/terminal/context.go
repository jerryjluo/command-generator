package terminal

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const (
	ScrollbackLines = 100
)

// TmuxInfo holds information about the current tmux session
type TmuxInfo struct {
	InTmux  bool   `json:"in_tmux"`
	Session string `json:"session,omitempty"`
	Window  string `json:"window,omitempty"`
	Pane    string `json:"pane,omitempty"`
}

// InTmux returns true if running inside a tmux session
func InTmux() bool {
	return os.Getenv("TMUX") != ""
}

// GetTmuxInfo returns information about the current tmux session, window, and pane
func GetTmuxInfo() TmuxInfo {
	if !InTmux() {
		return TmuxInfo{InTmux: false}
	}

	info := TmuxInfo{InTmux: true}

	// Get session name
	if out, err := exec.Command("tmux", "display-message", "-p", "#S").Output(); err == nil {
		info.Session = strings.TrimSpace(string(out))
	}

	// Get window name
	if out, err := exec.Command("tmux", "display-message", "-p", "#W").Output(); err == nil {
		info.Window = strings.TrimSpace(string(out))
	}

	// Get pane index
	if out, err := exec.Command("tmux", "display-message", "-p", "#P").Output(); err == nil {
		info.Pane = strings.TrimSpace(string(out))
	}

	return info
}

// CaptureContext captures the terminal scrollback via tmux
// Returns the context string and a warning message if not in tmux
func CaptureContext(lines int) (context string, warning string, err error) {
	if !InTmux() {
		return "", "Warning: Not running in tmux. Terminal context capture is unavailable.", nil
	}

	// Capture scrollback from tmux
	// -p: print to stdout
	// -S -N: start from N lines back (negative = scrollback)
	cmd := exec.Command("tmux", "capture-pane", "-p", "-S", fmt.Sprintf("-%d", lines))
	output, err := cmd.Output()
	if err != nil {
		return "", "", fmt.Errorf("failed to capture tmux pane: %w", err)
	}

	context = strings.TrimSpace(string(output))
	return context, "", nil
}
