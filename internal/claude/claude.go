package claude

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

const (
	jsonSchema = `{
		"type": "object",
		"properties": {
			"command": {
				"type": "string",
				"description": "The exact shell command to execute"
			},
			"explanation": {
				"type": "string",
				"description": "A breakdown explaining each tool, argument, and flag used"
			}
		},
		"required": ["command", "explanation"]
	}`

	systemPromptAddition = `You are a CLI command generator. Your task is to generate shell commands based on the user's natural language request.

IMPORTANT: You must respond with valid JSON matching the required schema. Do not include any text outside the JSON object.

When generating commands:
- Consider the terminal context provided to understand the user's current environment
- Generate a single, complete command that accomplishes the task
- In the explanation, break down each tool, argument, and flag used
- Format the explanation with bullet points for clarity`
)

// Response represents the JSON response from Claude
type Response struct {
	Command     string `json:"command"`
	Explanation string `json:"explanation"`
}

// ClaudeResponse represents the outer JSON response from claude CLI
type ClaudeResponse struct {
	Result           string    `json:"result"`
	StructuredOutput *Response `json:"structured_output"`
	Error            bool      `json:"is_error"`
}

// GenerateResult contains all data from a generation call for logging
type GenerateResult struct {
	Response     *Response
	SystemPrompt string
	UserPrompt   string
	RawOutput    string
}

// GenerateCommand calls the claude CLI to generate a command
func GenerateCommand(model, claudeMdContent, terminalContext, buildToolsContext, docsContext, userQuery string, feedback string) (*GenerateResult, error) {
	// Build the prompt
	prompt := buildPrompt(terminalContext, buildToolsContext, docsContext, userQuery, feedback)

	// Build the system prompt
	systemPrompt := systemPromptAddition
	if claudeMdContent != "" {
		systemPrompt = claudeMdContent + "\n\n" + systemPromptAddition
	}

	// Build claude command arguments
	args := []string{
		"-p",
		"--model", model,
		"--output-format", "json",
		"--append-system-prompt", systemPrompt,
		"--json-schema", jsonSchema,
		prompt,
	}

	cmd := exec.Command("claude", args...)
	output, err := cmd.CombinedOutput()
	rawOutput := string(output)

	if err != nil {
		// Try to parse output even on error - sometimes it contains useful info
		if len(output) > 0 {
			return nil, fmt.Errorf("claude CLI error: %s", rawOutput)
		}
		return nil, fmt.Errorf("failed to execute claude CLI: %w", err)
	}

	// Parse the outer JSON response
	var claudeResp ClaudeResponse
	if err := json.Unmarshal(output, &claudeResp); err != nil {
		return nil, fmt.Errorf("failed to parse claude response: %w", err)
	}

	if claudeResp.Error {
		return nil, fmt.Errorf("claude returned an error")
	}

	// Check for structured_output first (used when --json-schema is provided)
	if claudeResp.StructuredOutput != nil {
		return &GenerateResult{
			Response:     claudeResp.StructuredOutput,
			SystemPrompt: systemPrompt,
			UserPrompt:   prompt,
			RawOutput:    rawOutput,
		}, nil
	}

	// Fallback: parse the inner result (the actual command response)
	// Claude sometimes wraps JSON in markdown code blocks or adds extra text
	jsonStr := extractJSON(claudeResp.Result)
	if jsonStr == "" {
		return nil, fmt.Errorf("no valid JSON found in response: %s", claudeResp.Result)
	}

	var response Response
	if err := json.Unmarshal([]byte(jsonStr), &response); err != nil {
		return nil, fmt.Errorf("failed to parse command response: %w (json was: %s)", err, jsonStr)
	}

	return &GenerateResult{
		Response:     &response,
		SystemPrompt: systemPrompt,
		UserPrompt:   prompt,
		RawOutput:    rawOutput,
	}, nil
}

// extractJSON tries to extract a JSON object from text that may contain markdown or extra content
func extractJSON(text string) string {
	text = strings.TrimSpace(text)

	// First, try to parse as-is (in case it's already valid JSON)
	if strings.HasPrefix(text, "{") {
		// Find the matching closing brace
		if jsonStr := extractJSONObject(text); jsonStr != "" {
			return jsonStr
		}
	}

	// Try to extract from markdown code blocks
	// Match ```json ... ``` or ``` ... ```
	re := regexp.MustCompile("(?s)```(?:json)?\\s*\\n?(\\{.*?\\})\\s*\\n?```")
	matches := re.FindStringSubmatch(text)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}

	// Try to find a JSON object anywhere in the text
	startIdx := strings.Index(text, "{")
	if startIdx != -1 {
		if jsonStr := extractJSONObject(text[startIdx:]); jsonStr != "" {
			return jsonStr
		}
	}

	return ""
}

// extractJSONObject extracts a complete JSON object from a string starting with {
func extractJSONObject(text string) string {
	if !strings.HasPrefix(text, "{") {
		return ""
	}

	depth := 0
	inString := false
	escaped := false

	for i, ch := range text {
		if escaped {
			escaped = false
			continue
		}

		if ch == '\\' && inString {
			escaped = true
			continue
		}

		if ch == '"' {
			inString = !inString
			continue
		}

		if inString {
			continue
		}

		if ch == '{' {
			depth++
		} else if ch == '}' {
			depth--
			if depth == 0 {
				return text[:i+1]
			}
		}
	}

	return ""
}

// buildPrompt constructs the full prompt including context
func buildPrompt(terminalContext, buildToolsContext, docsContext, userQuery, feedback string) string {
	var sb strings.Builder

	if terminalContext != "" {
		sb.WriteString("Terminal context (recent scrollback):\n")
		sb.WriteString("---\n")
		sb.WriteString(terminalContext)
		sb.WriteString("\n---\n\n")
	}

	if buildToolsContext != "" {
		sb.WriteString("Available build tools and commands in current directory:\n")
		sb.WriteString("---\n")
		sb.WriteString(buildToolsContext)
		sb.WriteString("\n---\n\n")
	}

	if docsContext != "" {
		sb.WriteString("Project documentation (command-related sections):\n")
		sb.WriteString("---\n")
		sb.WriteString(docsContext)
		sb.WriteString("\n---\n\n")
	}

	sb.WriteString("User request: ")
	sb.WriteString(userQuery)

	if feedback != "" {
		sb.WriteString("\n\nUser feedback on previous command: ")
		sb.WriteString(feedback)
	}

	sb.WriteString("\n\nGenerate a single shell command that accomplishes this task.")

	return sb.String()
}

// CheckClaudeCLI verifies that the claude CLI is installed
func CheckClaudeCLI() error {
	_, err := exec.LookPath("claude")
	if err != nil {
		return fmt.Errorf("claude CLI not found. Please install Claude Code: https://claude.ai/code")
	}
	return nil
}
