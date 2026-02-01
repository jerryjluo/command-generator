package buildtools

import (
	"bufio"
	"bytes"
	"regexp"
	"strings"
)

// JustfileParser parses justfile recipes
type JustfileParser struct{}

// FileName returns the justfile filename
func (p *JustfileParser) FileName() string {
	return "justfile"
}

var (
	// Pattern to match recipe definitions
	justRecipeRegex = regexp.MustCompile(`^([a-zA-Z_][a-zA-Z0-9_-]*)\s*(?:[^:]*)?:`)
	// Pattern to match preceding comment for description
	justDescRegex = regexp.MustCompile(`^#\s*(.+)$`)
)

// Parse extracts just recipes from a justfile
func (p *JustfileParser) Parse(content []byte) (*Tool, error) {
	tool := &Tool{
		Name:     "just",
		File:     "justfile",
		Commands: []Command{},
	}

	scanner := bufio.NewScanner(bytes.NewReader(content))
	var lastComment string

	for scanner.Scan() {
		line := scanner.Text()

		// Skip empty lines but reset last comment
		if strings.TrimSpace(line) == "" {
			lastComment = ""
			continue
		}

		// Check for comment (potential description for next recipe)
		if matches := justDescRegex.FindStringSubmatch(line); len(matches) > 1 {
			lastComment = matches[1]
			continue
		}

		// Skip lines starting with whitespace (recipe body)
		if len(line) > 0 && (line[0] == ' ' || line[0] == '\t') {
			continue
		}

		// Skip variable assignments
		if strings.Contains(line, ":=") || strings.Contains(line, "=") && !strings.Contains(line, ":") {
			lastComment = ""
			continue
		}

		// Check for recipe
		if matches := justRecipeRegex.FindStringSubmatch(line); len(matches) > 1 {
			recipe := matches[1]

			// Skip private recipes (starting with _)
			if strings.HasPrefix(recipe, "_") {
				lastComment = ""
				continue
			}

			cmd := Command{Name: recipe}
			if lastComment != "" {
				cmd.Description = lastComment
			}
			tool.Commands = append(tool.Commands, cmd)
		}

		lastComment = ""
	}

	if len(tool.Commands) == 0 {
		return nil, nil
	}

	return tool, nil
}
