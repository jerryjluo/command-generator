package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

// Colors.
var (
	colorGreen  = lipgloss.Color("2")
	colorRed    = lipgloss.Color("1")
	colorGray   = lipgloss.Color("8")
	colorPurple = lipgloss.Color("5")
	colorBlue   = lipgloss.Color("4")
	colorCyan   = lipgloss.Color("6")
	colorSubtle = lipgloss.Color("241")
)

// Layout styles.
var (
	titleStyle    = lipgloss.NewStyle().Bold(true)
	subtitleStyle = lipgloss.NewStyle().Foreground(colorSubtle)

	selectedTabStyle = lipgloss.NewStyle().Bold(true).Underline(true).Padding(0, 1)
	tabStyle         = lipgloss.NewStyle().Foreground(colorSubtle).Padding(0, 1)
	tabGapStyle      = lipgloss.NewStyle().SetString("  ")

	helpStyle = lipgloss.NewStyle().Foreground(colorSubtle)

	codeBlockStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("235")).
			Padding(0, 1)
)

// StatusStyle returns styled status text for the given status string.
func StatusStyle(status string) string {
	s := strings.ToLower(status)
	switch s {
	case "accepted":
		return lipgloss.NewStyle().Foreground(colorGreen).Render("✓ accepted")
	case "rejected":
		return lipgloss.NewStyle().Foreground(colorRed).Render("✗ rejected")
	case "quit":
		return lipgloss.NewStyle().Foreground(colorGray).Render("- quit")
	default:
		return s
	}
}

// ModelStyle returns styled model text for the given model name.
func ModelStyle(model string) string {
	m := strings.ToLower(model)
	switch {
	case strings.Contains(m, "opus"):
		return lipgloss.NewStyle().Foreground(colorPurple).Render(model)
	case strings.Contains(m, "sonnet"):
		return lipgloss.NewStyle().Foreground(colorBlue).Render(model)
	case strings.Contains(m, "haiku"):
		return lipgloss.NewStyle().Foreground(colorCyan).Render(model)
	default:
		return model
	}
}

// TableStyles returns custom table.Styles for the list view.
func TableStyles() table.Styles {
	s := table.Styles{
		Header: lipgloss.NewStyle().
			Bold(true).
			Padding(0, 1).
			Foreground(colorSubtle).
			BorderBottom(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(colorSubtle),
		Cell: lipgloss.NewStyle().
			Padding(0, 1),
		Selected: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("57")).
			Padding(0, 1),
	}
	return s
}
