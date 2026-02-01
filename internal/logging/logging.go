package logging

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/jerryluo/cmd/internal/terminal"
)

// FinalStatus represents the outcome of a session
type FinalStatus string

const (
	StatusAccepted FinalStatus = "accepted"
	StatusRejected FinalStatus = "rejected"
	StatusQuit     FinalStatus = "quit"
)

// ContextSources holds the context data fed into the prompt
type ContextSources struct {
	ClaudeMdContent string `json:"claude_md_content"`
	TerminalContext string `json:"terminal_context"`
}

// ModelInput holds the prompts sent to Claude
type ModelInput struct {
	SystemPrompt string `json:"system_prompt"`
	UserPrompt   string `json:"user_prompt"`
}

// ModelOutput holds Claude's response
type ModelOutput struct {
	RawResponse string `json:"raw_response"`
	Command     string `json:"command"`
	Explanation string `json:"explanation"`
}

// Iteration represents a single generate-feedback cycle
type Iteration struct {
	Feedback    string      `json:"feedback"`
	ModelInput  ModelInput  `json:"model_input"`
	ModelOutput ModelOutput `json:"model_output"`
	Timestamp   time.Time   `json:"timestamp"`
}

// Metadata holds session metadata
type Metadata struct {
	Timestamp      time.Time           `json:"timestamp"`
	Model          string              `json:"model"`
	FinalStatus    FinalStatus         `json:"final_status"`
	FinalFeedback  string              `json:"final_feedback,omitempty"`
	IterationCount int                 `json:"iteration_count"`
	TmuxInfo       terminal.TmuxInfo   `json:"tmux_info"`
}

// SessionLog is the complete log for one CLI invocation
type SessionLog struct {
	UserQuery      string         `json:"user_query"`
	ContextSources ContextSources `json:"context_sources"`
	Iterations     []Iteration    `json:"iterations"`
	Metadata       Metadata       `json:"metadata"`
}

// Logger manages the session log file
type Logger struct {
	log      *SessionLog
	filePath string
	mu       sync.Mutex
}

// GetLogDir returns the log directory path: ~/.local/share/cmd/logs/
func GetLogDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".local", "share", "cmd", "logs"), nil
}

// ensureLogDir creates the log directory if it doesn't exist
func ensureLogDir() error {
	logDir, err := GetLogDir()
	if err != nil {
		return err
	}
	return os.MkdirAll(logDir, 0755)
}

// generateFilename creates a filename from the given UTC time.
// Format: 2024-01-15T14-30-45Z.json
func generateFilename(t time.Time) string {
	return t.UTC().Format("2006-01-02T15-04-05Z") + ".json"
}

// NewLogger creates a new logger and initializes the session log.
// Returns nil with logged warning if logging setup fails.
func NewLogger(
	userQuery string,
	claudeMdContent string,
	terminalContext string,
	model string,
	tmuxInfo terminal.TmuxInfo,
) *Logger {
	if err := ensureLogDir(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Could not create log directory: %v\n", err)
		return nil
	}

	logDir, err := GetLogDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Could not get log directory: %v\n", err)
		return nil
	}

	now := time.Now()
	filePath := filepath.Join(logDir, generateFilename(now))

	logger := &Logger{
		log: &SessionLog{
			UserQuery: userQuery,
			ContextSources: ContextSources{
				ClaudeMdContent: claudeMdContent,
				TerminalContext: terminalContext,
			},
			Iterations: []Iteration{},
			Metadata: Metadata{
				Timestamp:      now.UTC(),
				Model:          model,
				FinalStatus:    StatusQuit, // Default, will be updated on finalize
				IterationCount: 0,
				TmuxInfo:       tmuxInfo,
			},
		},
		filePath: filePath,
	}

	// Write initial log file
	logger.save()

	return logger
}

// AddIteration records a new generation attempt.
func (l *Logger) AddIteration(
	feedback string,
	systemPrompt string,
	userPrompt string,
	rawResponse string,
	command string,
	explanation string,
) {
	if l == nil {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	iteration := Iteration{
		Feedback: feedback,
		ModelInput: ModelInput{
			SystemPrompt: systemPrompt,
			UserPrompt:   userPrompt,
		},
		ModelOutput: ModelOutput{
			RawResponse: rawResponse,
			Command:     command,
			Explanation: explanation,
		},
		Timestamp: time.Now().UTC(),
	}

	l.log.Iterations = append(l.log.Iterations, iteration)
	l.log.Metadata.IterationCount = len(l.log.Iterations)

	l.save()
}

// Finalize records the final status and writes the complete log.
func (l *Logger) Finalize(status FinalStatus, finalFeedback string) {
	if l == nil {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	l.log.Metadata.FinalStatus = status
	if finalFeedback != "" {
		l.log.Metadata.FinalFeedback = finalFeedback
	}

	l.save()
}

// save writes the current log state to the JSON file using atomic write.
func (l *Logger) save() error {
	data, err := json.MarshalIndent(l.log, "", "  ")
	if err != nil {
		return err
	}

	// Atomic write: write to temp file then rename
	tmpPath := l.filePath + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return err
	}

	return os.Rename(tmpPath, l.filePath)
}
