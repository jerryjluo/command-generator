package buildtools

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Command represents a single available command/target from a build tool
type Command struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// Tool represents a detected build tool and its available commands
type Tool struct {
	Name     string    `json:"name"`
	File     string    `json:"file"`
	Commands []Command `json:"commands"`
}

// DetectionResult contains all detected build tools in a directory
type DetectionResult struct {
	Tools []Tool `json:"tools"`
}

// Parser interface for each build tool type
type Parser interface {
	// FileName returns the config file to look for
	FileName() string
	// Parse reads the file and extracts commands
	Parse(content []byte) (*Tool, error)
}

// Detect scans a directory for known build tool configuration files
// and returns a DetectionResult with all detected tools and their commands
func Detect(dir string) *DetectionResult {
	result := &DetectionResult{Tools: []Tool{}}

	parsers := []Parser{
		&MakefileParser{},
		&PackageJSONParser{},
		&MiseParser{},
		&JustfileParser{},
		&TaskfileParser{},
		&CargoParser{},
		&PyprojectParser{},
		&DockerComposeParser{},
	}

	for _, parser := range parsers {
		filePath := filepath.Join(dir, parser.FileName())
		content, err := os.ReadFile(filePath)
		if err != nil {
			continue // File doesn't exist or can't read, skip silently
		}

		tool, err := parser.Parse(content)
		if err != nil {
			continue // Parse error, skip silently
		}

		if tool != nil && len(tool.Commands) > 0 {
			result.Tools = append(result.Tools, *tool)
		}
	}

	return result
}

// FormatForPrompt returns a human-readable representation for the Claude prompt
func (r *DetectionResult) FormatForPrompt() string {
	if len(r.Tools) == 0 {
		return ""
	}

	var sb strings.Builder
	for _, tool := range r.Tools {
		sb.WriteString(fmt.Sprintf("%s (%s):\n", tool.Name, tool.File))
		for _, cmd := range tool.Commands {
			if cmd.Description != "" {
				sb.WriteString(fmt.Sprintf("  - %s: %s\n", cmd.Name, cmd.Description))
			} else {
				sb.WriteString(fmt.Sprintf("  - %s\n", cmd.Name))
			}
		}
		sb.WriteString("\n")
	}
	return strings.TrimSuffix(sb.String(), "\n")
}
