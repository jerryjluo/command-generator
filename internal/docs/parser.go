package docs

import (
	"regexp"
	"strings"
)

// RelevantHeadingPatterns contains patterns to match in headings (case-insensitive)
var RelevantHeadingPatterns = []string{
	"build",
	"development",
	"dev",
	"installation",
	"install",
	"usage",
	"commands",
	"cli",
	"getting started",
	"quick start",
	"quickstart",
	"running",
	"run",
	"setup",
	"prerequisites",
	"requirements",
}

// ShellFences contains code fence languages indicating shell commands
var ShellFences = []string{"bash", "shell", "sh", "zsh", "console", "terminal", ""}

// headingRegex matches markdown headings (# to ######)
var headingRegex = regexp.MustCompile(`^(#{1,6})\s+(.+)$`)

// codeFenceRegex matches opening code fence with optional language
var codeFenceRegex = regexp.MustCompile("^```(\\w*)\\s*$")

// parseMarkdown extracts command-related sections from markdown content
func parseMarkdown(content []byte, filename string) []Section {
	lines := strings.Split(string(content), "\n")
	var sections []Section

	var currentSection *Section
	currentLevel := 0
	inCodeBlock := false
	var codeBlockBuilder strings.Builder

	for i := 0; i < len(lines); i++ {
		line := lines[i]

		// Track code blocks to avoid parsing headings inside them
		if strings.HasPrefix(line, "```") {
			inCodeBlock = !inCodeBlock
		}

		if inCodeBlock {
			if currentSection != nil {
				currentSection.Content += line + "\n"
			}
			continue
		}

		// Check for heading
		if match := headingRegex.FindStringSubmatch(line); match != nil {
			level := len(match[1])
			headingText := match[2]

			// Close previous section if we hit same or higher level heading
			if currentSection != nil && level <= currentLevel {
				currentSection.Content = strings.TrimRight(currentSection.Content, "\n")
				sections = append(sections, *currentSection)
				currentSection = nil
			}

			// Check if this heading is relevant
			if isRelevantHeading(headingText) {
				currentSection = &Section{
					Heading: line,
					Content: line + "\n",
					Source:  filename,
				}
				currentLevel = level
			}
			continue
		}

		// If we're in a relevant section, accumulate content
		if currentSection != nil {
			currentSection.Content += line + "\n"
			continue
		}

		// Check for standalone shell code blocks outside relevant sections
		if fenceMatch := codeFenceRegex.FindStringSubmatch(line); fenceMatch != nil {
			lang := fenceMatch[1]
			if isShellFence(lang) {
				// Extract the code block with context
				section := extractCodeBlockWithContext(lines, i, filename)
				if section != nil {
					sections = append(sections, *section)
				}
				// Skip past the code block
				for i++; i < len(lines) && !strings.HasPrefix(lines[i], "```"); i++ {
				}
			}
		}
	}

	// Close final section if still open
	if currentSection != nil {
		currentSection.Content = strings.TrimRight(currentSection.Content, "\n")
		sections = append(sections, *currentSection)
	}

	// Reset code block builder
	codeBlockBuilder.Reset()

	return sections
}

// isRelevantHeading checks if a heading contains relevant keywords
func isRelevantHeading(text string) bool {
	lowered := strings.ToLower(text)
	for _, pattern := range RelevantHeadingPatterns {
		if strings.Contains(lowered, pattern) {
			return true
		}
	}
	return false
}

// isShellFence checks if a code fence language indicates shell commands
func isShellFence(lang string) bool {
	lang = strings.ToLower(lang)
	for _, shellLang := range ShellFences {
		if lang == shellLang {
			return true
		}
	}
	return false
}

// extractCodeBlockWithContext extracts a code block with preceding context
func extractCodeBlockWithContext(lines []string, blockIdx int, filename string) *Section {
	// Find end of code block
	endIdx := blockIdx + 1
	for endIdx < len(lines) && !strings.HasPrefix(lines[endIdx], "```") {
		endIdx++
	}

	// Include closing fence
	if endIdx < len(lines) {
		endIdx++
	}

	// Look back for preceding context (non-empty lines before the code block)
	contextStart := blockIdx
	for contextStart > 0 {
		prevLine := strings.TrimSpace(lines[contextStart-1])
		// Stop at empty lines, headings, or other code fences
		if prevLine == "" || strings.HasPrefix(prevLine, "#") || strings.HasPrefix(prevLine, "```") {
			break
		}
		contextStart--
	}

	// Build the section content
	var content strings.Builder
	for i := contextStart; i < endIdx; i++ {
		content.WriteString(lines[i])
		content.WriteString("\n")
	}

	return &Section{
		Heading: "",
		Content: strings.TrimRight(content.String(), "\n"),
		Source:  filename,
	}
}
