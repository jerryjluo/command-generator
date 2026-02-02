package docs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// DocFiles defines the documentation files to read in priority order
var DocFiles = []string{"README.md", "CLAUDE.md", "AGENTS.md"}

// Section represents an extracted command-related section from documentation
type Section struct {
	Heading string // The heading text (e.g., "## Build Commands")
	Content string // The full section content including code blocks
	Source  string // Source file (README.md, CLAUDE.md, etc.)
}

// Result contains all extracted documentation sections
type Result struct {
	Sections []Section
}

// Detect reads documentation files and extracts command-related sections
func Detect(dir string) *Result {
	result := &Result{Sections: []Section{}}

	for _, filename := range DocFiles {
		filePath := filepath.Join(dir, filename)
		content, err := os.ReadFile(filePath)
		if err != nil {
			continue // File doesn't exist or can't read, skip silently
		}

		sections := parseMarkdown(content, filename)
		result.Sections = append(result.Sections, sections...)
	}

	return result
}

// FormatForPrompt returns a human-readable representation for the Claude prompt
func (r *Result) FormatForPrompt() string {
	if len(r.Sections) == 0 {
		return ""
	}

	var sb strings.Builder
	currentSource := ""

	for _, section := range r.Sections {
		if section.Source != currentSource {
			if currentSource != "" {
				sb.WriteString("\n")
			}
			sb.WriteString(fmt.Sprintf("From %s:\n", section.Source))
			currentSource = section.Source
		}
		sb.WriteString(section.Content)
		sb.WriteString("\n")
	}

	return strings.TrimSpace(sb.String())
}
