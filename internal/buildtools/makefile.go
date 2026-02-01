package buildtools

import (
	"bufio"
	"bytes"
	"regexp"
	"strings"
)

// MakefileParser parses Makefile targets
type MakefileParser struct{}

// FileName returns the Makefile filename
func (p *MakefileParser) FileName() string {
	return "Makefile"
}

var (
	// Pattern to match Makefile targets (excludes variable assignments with =)
	makeTargetRegex = regexp.MustCompile(`^([a-zA-Z_][a-zA-Z0-9_-]*)\s*:(?:[^=]|$)`)
	// Pattern to match preceding comment for description
	makeDescRegex = regexp.MustCompile(`^#\s*(.+)$`)
)

// Parse extracts make targets from a Makefile
func (p *MakefileParser) Parse(content []byte) (*Tool, error) {
	tool := &Tool{
		Name:     "make",
		File:     "Makefile",
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

		// Check for comment (potential description for next target)
		if matches := makeDescRegex.FindStringSubmatch(line); len(matches) > 1 {
			lastComment = matches[1]
			continue
		}

		// Skip lines starting with tab (recipe lines)
		if strings.HasPrefix(line, "\t") {
			continue
		}

		// Skip special targets starting with .
		if strings.HasPrefix(strings.TrimSpace(line), ".") {
			lastComment = ""
			continue
		}

		// Check for target
		if matches := makeTargetRegex.FindStringSubmatch(line); len(matches) > 1 {
			target := matches[1]

			// Skip pattern rules (containing %)
			if strings.Contains(target, "%") {
				lastComment = ""
				continue
			}

			cmd := Command{Name: target}
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
