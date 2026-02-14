package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jerryluo/cmd/internal/logging"
)

// Tab indices
const (
	tabResponse = iota
	tabSystemPrompt
	tabUserPrompt
	tabUserQuery
	tabTmuxContext
	tabDocContext
	tabBuildTools
	tabPreferences
)

// logLoadedMsg is sent when a log has been loaded from disk.
type logLoadedMsg struct {
	log *logging.SessionLogWithID
	err error
}

type detailModel struct {
	log           *logging.SessionLogWithID
	viewport      viewport.Model
	activeTab     int
	tabs          []string
	width         int
	height        int
	ready         bool
	loading       bool
	err           error
	showHelp      bool
	statusMessage string
}

func newDetailModel() detailModel {
	return detailModel{
		tabs: []string{"Response", "System", "User Prompt", "Query", "Tmux", "Docs", "Build Tools", "Preferences"},
	}
}

// loadLog returns a tea.Cmd that reads a log by ID from disk.
func (m detailModel) loadLog(id string) tea.Cmd {
	return func() tea.Msg {
		log, err := logging.ReadLogWithID(id)
		return logLoadedMsg{log: log, err: err}
	}
}

func (m detailModel) Update(msg tea.Msg) (detailModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case logLoadedMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.log = msg.log
		m.activeTab = tabResponse
		if !m.ready {
			m.viewport = viewport.New(m.width, m.height-7)
			m.ready = true
		}
		m.viewport.SetContent(m.renderTabContent())
		m.viewport.GotoTop()
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "l":
			m.activeTab = (m.activeTab + 1) % len(m.tabs)
			m.viewport.SetContent(m.renderTabContent())
			m.viewport.GotoTop()
			return m, nil
		case "shift+tab", "h":
			m.activeTab = (m.activeTab - 1 + len(m.tabs)) % len(m.tabs)
			m.viewport.SetContent(m.renderTabContent())
			m.viewport.GotoTop()
			return m, nil
		case "1", "2", "3", "4", "5", "6", "7", "8":
			idx := int(msg.String()[0] - '1')
			if idx >= 0 && idx < len(m.tabs) {
				m.activeTab = idx
				m.viewport.SetContent(m.renderTabContent())
				m.viewport.GotoTop()
			}
			return m, nil
		}
	}

	if m.ready {
		m.viewport, cmd = m.viewport.Update(msg)
	}
	return m, cmd
}

// SetSize updates the dimensions and recalculates viewport size.
func (m *detailModel) SetSize(w, h int) {
	m.width = w
	m.height = h
	headerHeight := 6 // title + metadata + tab bar + borders
	footerHeight := 1 // help line
	if m.showHelp {
		footerHeight = 2
	}
	m.viewport.Width = w
	m.viewport.Height = h - headerHeight - footerHeight
}

func (m detailModel) View() string {
	if m.loading {
		return "\n  Loading..."
	}
	if m.err != nil {
		return fmt.Sprintf("\n  Error: %v", m.err)
	}
	if m.log == nil {
		return ""
	}

	var sections []string

	// Header: query title
	queryTitle := titleStyle.Width(m.width).Render(m.log.UserQuery)
	sections = append(sections, queryTitle)

	// Metadata line
	meta := m.renderMetadata()
	sections = append(sections, subtitleStyle.Render(meta))

	// Iteration history (only if >1 iteration)
	if len(m.log.Iterations) > 1 {
		sections = append(sections, m.renderIterationHistory())
	}

	// Tab bar
	sections = append(sections, m.renderTabBar())

	// Viewport
	if m.ready {
		sections = append(sections, m.viewport.View())
	}

	// Help bar
	var helpLine string
	if m.statusMessage != "" {
		status := lipgloss.NewStyle().Foreground(colorGreen).Render(m.statusMessage)
		helpLine = status + "  " + helpStyle.Render("esc back  tab/1-8 switch tab  c copy  ? help  q quit")
	} else if m.showHelp {
		helpLine = helpStyle.Render("esc back to list  tab/l next tab  shift+tab/h prev tab  1-8 jump to tab\nc copy content  ? toggle help  q quit  ↑/↓ scroll")
	} else {
		helpLine = helpStyle.Render("esc back  tab/1-8 switch tab  c copy  ? help  q quit")
	}
	sections = append(sections, helpLine)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (m detailModel) renderMetadata() string {
	ts := m.log.Metadata.Timestamp
	ago := timeAgo(ts)

	parts := []string{ago}

	if m.log.Metadata.IterationCount > 0 {
		iterLabel := "iteration"
		if m.log.Metadata.IterationCount != 1 {
			iterLabel = "iterations"
		}
		parts = append(parts, fmt.Sprintf("%d %s", m.log.Metadata.IterationCount, iterLabel))
	}

	parts = append(parts, StatusStyle(string(m.log.Metadata.FinalStatus)))
	parts = append(parts, ModelStyle(m.log.Metadata.Model))

	return strings.Join(parts, " · ")
}

func (m detailModel) renderIterationHistory() string {
	var lines []string
	for i, iter := range m.log.Iterations {
		cmd := iter.ModelOutput.Command
		if len(cmd) > 50 {
			cmd = cmd[:47] + "..."
		}
		if iter.Feedback != "" {
			lines = append(lines, fmt.Sprintf("  #%d: %q → %s", i+1, iter.Feedback, cmd))
		} else {
			lines = append(lines, fmt.Sprintf("  #%d: → %s", i+1, cmd))
		}
	}
	return strings.Join(lines, "\n")
}

func (m detailModel) renderTabBar() string {
	var tabs []string
	for i, t := range m.tabs {
		if i == m.activeTab {
			tabs = append(tabs, selectedTabStyle.Render(t))
		} else {
			tabs = append(tabs, tabStyle.Render(t))
		}
	}
	row := lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
	gap := tabGapStyle.Render(strings.Repeat(" ", max(0, m.width-lipgloss.Width(row))))
	return lipgloss.JoinHorizontal(lipgloss.Bottom, row, gap)
}

func (m detailModel) renderTabContent() string {
	if m.log == nil || len(m.log.Iterations) == 0 {
		return ""
	}

	lastIter := m.log.Iterations[len(m.log.Iterations)-1]

	switch m.activeTab {
	case tabResponse:
		return m.renderResponse(lastIter)
	case tabSystemPrompt:
		return renderTextBlock("System Prompt", lastIter.ModelInput.SystemPrompt)
	case tabUserPrompt:
		return renderTextBlock("User Prompt", lastIter.ModelInput.UserPrompt)
	case tabUserQuery:
		return renderTextBlock("User Query", m.log.UserQuery)
	case tabTmuxContext:
		return m.renderTmuxContext()
	case tabDocContext:
		return renderTextBlock("Documentation Context", m.log.ContextSources.DocumentationContext)
	case tabBuildTools:
		return m.renderBuildTools(lastIter)
	case tabPreferences:
		return renderTextBlock("Preferences (claude.md)", m.log.ContextSources.ClaudeMdContent)
	default:
		return ""
	}
}

func (m detailModel) renderResponse(iter logging.Iteration) string {
	var s strings.Builder

	s.WriteString("  Command:\n")
	cmd := codeBlockStyle.Width(m.width - 4).Render(iter.ModelOutput.Command)
	s.WriteString(cmd)
	s.WriteString("\n\n")

	if iter.ModelOutput.Explanation != "" {
		s.WriteString("  Explanation:\n")
		s.WriteString("  ")
		s.WriteString(iter.ModelOutput.Explanation)
		s.WriteString("\n")
	}

	return s.String()
}

func (m detailModel) renderTmuxContext() string {
	var s strings.Builder

	info := m.log.Metadata.TmuxInfo
	if info.InTmux {
		s.WriteString("  Tmux Session Info:\n")
		s.WriteString(fmt.Sprintf("    Session: %s\n", info.Session))
		s.WriteString(fmt.Sprintf("    Window:  %s\n", info.Window))
		s.WriteString(fmt.Sprintf("    Pane:    %s\n", info.Pane))
		s.WriteString("\n")
	} else {
		s.WriteString("  Not running in tmux\n\n")
	}

	ctx := m.log.ContextSources.TerminalContext
	if ctx != "" {
		s.WriteString("  Terminal Scrollback:\n")
		s.WriteString(ctx)
		s.WriteString("\n")
	} else {
		s.WriteString("  No terminal context captured\n")
	}

	return s.String()
}

func (m detailModel) renderBuildTools(iter logging.Iteration) string {
	prompt := iter.ModelInput.UserPrompt

	// Look for the build tools section in the user prompt
	markers := []string{"Available commands", "Build tools", "Available build"}
	var section string
	for _, marker := range markers {
		idx := strings.Index(strings.ToLower(prompt), strings.ToLower(marker))
		if idx >= 0 {
			section = prompt[idx:]
			// Find the end of the section (next major section or end)
			if end := strings.Index(section[1:], "\n\n\n"); end >= 0 {
				section = section[:end+1]
			}
			break
		}
	}

	if section != "" {
		return renderTextBlock("Build Tools", section)
	}
	return "  No build tools information available\n"
}

// copyableContent returns the raw text content for the active tab, suitable for clipboard copy.
func (m detailModel) copyableContent() string {
	if m.log == nil || len(m.log.Iterations) == 0 {
		return ""
	}
	lastIter := m.log.Iterations[len(m.log.Iterations)-1]
	switch m.activeTab {
	case tabResponse:
		return lastIter.ModelOutput.Command
	case tabSystemPrompt:
		return lastIter.ModelInput.SystemPrompt
	case tabUserPrompt:
		return lastIter.ModelInput.UserPrompt
	case tabUserQuery:
		return m.log.UserQuery
	case tabTmuxContext:
		return m.log.ContextSources.TerminalContext
	case tabDocContext:
		return m.log.ContextSources.DocumentationContext
	case tabBuildTools:
		return m.extractBuildToolsText(lastIter)
	case tabPreferences:
		return m.log.ContextSources.ClaudeMdContent
	}
	return ""
}

// extractBuildToolsText returns the raw build tools section from the user prompt.
func (m detailModel) extractBuildToolsText(iter logging.Iteration) string {
	prompt := iter.ModelInput.UserPrompt
	markers := []string{"Available commands", "Build tools", "Available build"}
	for _, marker := range markers {
		idx := strings.Index(strings.ToLower(prompt), strings.ToLower(marker))
		if idx >= 0 {
			section := prompt[idx:]
			if end := strings.Index(section[1:], "\n\n\n"); end >= 0 {
				section = section[:end+1]
			}
			return section
		}
	}
	return ""
}

func renderTextBlock(title, content string) string {
	if content == "" {
		return fmt.Sprintf("  No %s available\n", strings.ToLower(title))
	}
	var s strings.Builder
	s.WriteString(fmt.Sprintf("  %s:\n\n", title))
	s.WriteString(content)
	s.WriteString("\n")
	return s.String()
}

func timeAgo(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		m := int(d.Minutes())
		if m == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", m)
	case d < 24*time.Hour:
		h := int(d.Hours())
		if h == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", h)
	default:
		days := int(d.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	}
}
