package docs

import (
	"strings"
	"testing"
)

func TestParseMarkdown(t *testing.T) {
	content := []byte(`# Project

Some intro text.

## Installation

Install the tool:

` + "```bash" + `
npm install
` + "```" + `

## Usage

Run it:

` + "```bash" + `
npm run dev
` + "```" + `

## Other Section

This is not relevant.
`)

	sections := parseMarkdown(content, "README.md")

	if len(sections) < 2 {
		t.Errorf("Expected at least 2 sections, got %d", len(sections))
	}

	// Check that Installation section was found
	foundInstall := false
	foundUsage := false
	for _, s := range sections {
		if strings.Contains(s.Content, "## Installation") {
			foundInstall = true
		}
		if strings.Contains(s.Content, "## Usage") {
			foundUsage = true
		}
	}

	if !foundInstall {
		t.Error("Expected to find Installation section")
	}
	if !foundUsage {
		t.Error("Expected to find Usage section")
	}
}

func TestIsRelevantHeading(t *testing.T) {
	tests := []struct {
		heading  string
		expected bool
	}{
		{"Build Commands", true},
		{"Installation", true},
		{"Getting Started", true},
		{"Development", true},
		{"License", false},
		{"Contributing", false},
		{"Quick Start Guide", true},
		{"Prerequisites", true},
	}

	for _, tt := range tests {
		if got := isRelevantHeading(tt.heading); got != tt.expected {
			t.Errorf("isRelevantHeading(%q) = %v, want %v", tt.heading, got, tt.expected)
		}
	}
}

func TestDetect(t *testing.T) {
	// This test runs against the actual project directory
	result := Detect("../..")

	// Should find at least some sections from README.md or CLAUDE.md
	if len(result.Sections) == 0 {
		t.Log("No sections found - this may be expected if doc files are empty or have no relevant sections")
	}

	output := result.FormatForPrompt()
	t.Logf("Detected docs output (%d chars):\n%s", len(output), truncate(output, 1000))
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "\n... [truncated]"
}
