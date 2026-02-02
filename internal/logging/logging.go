package logging

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
	ClaudeMdContent      string `json:"claude_md_content"`
	TerminalContext      string `json:"terminal_context"`
	DocumentationContext string `json:"documentation_context"`
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
	docsContext string,
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
				ClaudeMdContent:      claudeMdContent,
				TerminalContext:      terminalContext,
				DocumentationContext: docsContext,
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

// LogSummary is a lightweight representation of a log for list views
type LogSummary struct {
	ID             string      `json:"id"`
	UserQuery      string      `json:"user_query"`
	FinalStatus    FinalStatus `json:"final_status"`
	Model          string      `json:"model"`
	Timestamp      time.Time   `json:"timestamp"`
	IterationCount int         `json:"iteration_count"`
	CommandPreview string      `json:"command_preview"`
	TmuxSession    string      `json:"tmux_session,omitempty"`
}

// ListLogs returns summaries of all log files in the log directory.
// Results are sorted by timestamp descending (newest first).
func ListLogs() ([]LogSummary, error) {
	logDir, err := GetLogDir()
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(logDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []LogSummary{}, nil
		}
		return nil, err
	}

	var summaries []LogSummary
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		// Skip temp files
		if filepath.Ext(filepath.Base(entry.Name())[:len(entry.Name())-5]) == ".tmp" {
			continue
		}

		filePath := filepath.Join(logDir, entry.Name())
		log, err := readLogFile(filePath)
		if err != nil {
			continue // Skip files that can't be parsed
		}

		// Extract command preview from last iteration
		commandPreview := ""
		if len(log.Iterations) > 0 {
			cmd := log.Iterations[len(log.Iterations)-1].ModelOutput.Command
			if len(cmd) > 80 {
				commandPreview = cmd[:77] + "..."
			} else {
				commandPreview = cmd
			}
		}

		// ID is filename without .json extension
		id := entry.Name()[:len(entry.Name())-5]

		summaries = append(summaries, LogSummary{
			ID:             id,
			UserQuery:      log.UserQuery,
			FinalStatus:    log.Metadata.FinalStatus,
			Model:          log.Metadata.Model,
			Timestamp:      log.Metadata.Timestamp,
			IterationCount: log.Metadata.IterationCount,
			CommandPreview: commandPreview,
			TmuxSession:    log.Metadata.TmuxInfo.Session,
		})
	}

	// Sort by timestamp descending (newest first)
	for i := 0; i < len(summaries)-1; i++ {
		for j := i + 1; j < len(summaries); j++ {
			if summaries[j].Timestamp.After(summaries[i].Timestamp) {
				summaries[i], summaries[j] = summaries[j], summaries[i]
			}
		}
	}

	return summaries, nil
}

// ReadLog reads and parses a single log file by ID.
// ID is the filename without the .json extension.
func ReadLog(id string) (*SessionLog, error) {
	logDir, err := GetLogDir()
	if err != nil {
		return nil, err
	}

	filePath := filepath.Join(logDir, id+".json")
	return readLogFile(filePath)
}

// readLogFile reads and parses a log file from the given path.
func readLogFile(filePath string) (*SessionLog, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var log SessionLog
	if err := json.Unmarshal(data, &log); err != nil {
		return nil, err
	}

	return &log, nil
}

// SessionLogWithID wraps SessionLog with an ID field for API responses
type SessionLogWithID struct {
	ID string `json:"id"`
	SessionLog
}

// ReadLogWithID reads a log and returns it with its ID included
func ReadLogWithID(id string) (*SessionLogWithID, error) {
	log, err := ReadLog(id)
	if err != nil {
		return nil, err
	}

	return &SessionLogWithID{
		ID:         id,
		SessionLog: *log,
	}, nil
}

// SearchLogs searches through all logs for the given query string.
// It searches in user_query, command, and explanation fields.
func SearchLogs(query string) ([]LogSummary, error) {
	if query == "" {
		return ListLogs()
	}

	allLogs, err := ListLogs()
	if err != nil {
		return nil, err
	}

	query = strings.ToLower(query)
	var results []LogSummary

	for _, summary := range allLogs {
		// Check user query
		if strings.Contains(strings.ToLower(summary.UserQuery), query) {
			results = append(results, summary)
			continue
		}

		// Check command preview
		if strings.Contains(strings.ToLower(summary.CommandPreview), query) {
			results = append(results, summary)
			continue
		}

		// For deeper search, load the full log
		log, err := ReadLog(summary.ID)
		if err != nil {
			continue
		}

		found := false
		for _, iter := range log.Iterations {
			if strings.Contains(strings.ToLower(iter.ModelOutput.Command), query) ||
				strings.Contains(strings.ToLower(iter.ModelOutput.Explanation), query) {
				found = true
				break
			}
		}

		if found {
			results = append(results, summary)
		}
	}

	return results, nil
}
