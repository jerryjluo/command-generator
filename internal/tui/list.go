package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/jerryluo/cmd/internal/logging"
)

// logsLoadedMsg is sent when logs have been loaded from disk.
type logsLoadedMsg struct {
	logs []logging.LogSummary
	err  error
}

// listModel is the model for the log list view.
type listModel struct {
	table         table.Model
	searchInput   textinput.Model
	searching     bool
	statusFilter  string
	allLogs       []logging.LogSummary
	filteredLogs  []logging.LogSummary
	width         int
	height        int
	err           error
	showHelp      bool
	statusMessage string
}

// newListModel initializes a new list model with default state.
func newListModel() listModel {
	ti := textinput.New()
	ti.Placeholder = "search logs..."
	ti.Prompt = "Search: "

	t := table.New(
		table.WithFocused(true),
		table.WithStyles(TableStyles()),
	)

	return listModel{
		table:       t,
		searchInput: ti,
	}
}

// loadLogs returns a tea.Cmd that loads logs from disk.
func (m listModel) loadLogs() tea.Cmd {
	return func() tea.Msg {
		logs, err := logging.ListLogs()
		return logsLoadedMsg{logs: logs, err: err}
	}
}

// SetSize updates the dimensions and recalculates table column widths.
func (m *listModel) SetSize(w, h int) {
	m.width = w
	m.height = h

	// Reserve space for title (1), search bar (1 if searching), help bar (1-2), padding
	tableHeight := h - 3
	if m.searching {
		tableHeight--
	}
	if m.showHelp {
		tableHeight--
	}

	m.table.SetWidth(w)
	m.table.SetHeight(tableHeight)
	m.table.SetColumns(m.columns())
	m.rebuildRows()
}

// columns computes column widths based on current terminal width.
func (m listModel) columns() []table.Column {
	w := m.width
	if w < 40 {
		w = 40
	}

	statusW := 10
	modelW := 8
	timeW := 10
	fixedW := statusW + modelW + timeW
	// Account for cell padding (1 on each side per column = 2 per column, 5 columns)
	remaining := w - fixedW - 10
	if remaining < 20 {
		remaining = 20
	}
	queryW := remaining * 40 / 100
	cmdW := remaining - queryW

	return []table.Column{
		{Title: "Query", Width: queryW},
		{Title: "Status", Width: statusW},
		{Title: "Model", Width: modelW},
		{Title: "Time", Width: timeW},
		{Title: "Command", Width: cmdW},
	}
}

// Update handles messages for the list view.
func (m listModel) Update(msg tea.Msg) (listModel, tea.Cmd) {
	switch msg := msg.(type) {
	case logsLoadedMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.allLogs = msg.logs
		m.applyFilters()
		m.table.SetCursor(0)
		return m, nil

	case tea.KeyMsg:
		if m.searching {
			return m.updateSearching(msg)
		}
		return m.updateNormal(msg)
	}

	return m, nil
}

// updateSearching handles key events when the search input is focused.
func (m listModel) updateSearching(msg tea.KeyMsg) (listModel, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.searching = false
		m.searchInput.SetValue("")
		m.searchInput.Blur()
		m.applyFilters()
		// Recalculate table height since search bar is hidden
		m.SetSize(m.width, m.height)
		return m, nil
	case "enter":
		m.searching = false
		m.searchInput.Blur()
		// Keep the filter text active
		return m, nil
	}

	var cmd tea.Cmd
	prevValue := m.searchInput.Value()
	m.searchInput, cmd = m.searchInput.Update(msg)
	if m.searchInput.Value() != prevValue {
		m.applyFilters()
	}
	return m, cmd
}

// updateNormal handles key events when the table is focused.
func (m listModel) updateNormal(msg tea.KeyMsg) (listModel, tea.Cmd) {
	switch msg.String() {
	case "/":
		m.searching = true
		m.searchInput.Focus()
		// Recalculate table height since search bar is shown
		m.SetSize(m.width, m.height)
		return m, nil
	case "s":
		m.cycleStatusFilter()
		m.applyFilters()
		return m, nil
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

// cycleStatusFilter cycles through "" -> "accepted" -> "rejected" -> "quit" -> "".
func (m *listModel) cycleStatusFilter() {
	switch m.statusFilter {
	case "":
		m.statusFilter = "accepted"
	case "accepted":
		m.statusFilter = "rejected"
	case "rejected":
		m.statusFilter = "quit"
	case "quit":
		m.statusFilter = ""
	}
}

// applyFilters filters allLogs based on search text and status filter,
// then rebuilds the table rows.
func (m *listModel) applyFilters() {
	search := strings.ToLower(m.searchInput.Value())

	m.filteredLogs = nil
	for _, log := range m.allLogs {
		if m.statusFilter != "" && string(log.FinalStatus) != m.statusFilter {
			continue
		}
		if search != "" {
			q := strings.ToLower(log.UserQuery)
			c := strings.ToLower(log.CommandPreview)
			if !strings.Contains(q, search) && !strings.Contains(c, search) {
				continue
			}
		}
		m.filteredLogs = append(m.filteredLogs, log)
	}

	m.rebuildRows()
}

// rebuildRows converts filteredLogs into table rows.
func (m *listModel) rebuildRows() {
	rows := make([]table.Row, len(m.filteredLogs))
	for i, log := range m.filteredLogs {
		rows[i] = table.Row{
			log.UserQuery,
			statusText(string(log.FinalStatus)),
			log.Model,
			shortTimeAgo(log.Timestamp),
			log.CommandPreview,
		}
	}
	m.table.SetRows(rows)
}

// View renders the list view.
func (m listModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error loading logs: %v", m.err)
	}

	var b strings.Builder

	// Title bar with filter status
	title := titleStyle.Render("cmd logs")
	filterInfo := m.filterStatusText()
	if filterInfo != "" {
		padding := m.width - lipgloss.Width(title) - lipgloss.Width(filterInfo) - 2
		if padding < 1 {
			padding = 1
		}
		b.WriteString(title + strings.Repeat(" ", padding) + filterInfo)
	} else {
		b.WriteString(title)
	}
	b.WriteString("\n")

	// Search bar (only when searching)
	if m.searching {
		b.WriteString(m.searchInput.View())
		b.WriteString("\n")
	}

	// Table
	b.WriteString(m.table.View())
	b.WriteString("\n")

	// Help bar
	var helpLine string
	if m.statusMessage != "" {
		status := lipgloss.NewStyle().Foreground(colorGreen).Render(m.statusMessage)
		helpLine = status + "  " + helpStyle.Render("/ search  s filter  c copy  enter view  ? help  q quit")
	} else if m.showHelp {
		helpLine = helpStyle.Render("/ search  s cycle status  c copy command  enter view details  ? toggle help  q quit\n↑/↓ navigate  page up/down scroll")
	} else {
		helpLine = helpStyle.Render("/ search  s filter  c copy  enter view  ? help  q quit")
	}
	b.WriteString(helpLine)

	return b.String()
}

// filterStatusText returns a string describing active filters.
func (m listModel) filterStatusText() string {
	parts := []string{}

	count := len(m.filteredLogs)
	total := len(m.allLogs)

	if m.statusFilter != "" {
		parts = append(parts, fmt.Sprintf("%d %s", count, m.statusFilter))
	} else if count != total {
		parts = append(parts, fmt.Sprintf("%d/%d", count, total))
	} else {
		parts = append(parts, fmt.Sprintf("%d logs", total))
	}

	search := m.searchInput.Value()
	if search != "" && !m.searching {
		parts = append(parts, fmt.Sprintf("search: %s", search))
	}

	if len(parts) > 0 {
		return "[" + strings.Join(parts, " | ") + "]"
	}
	return ""
}

// selectedLogID returns the ID of the currently selected log.
func (m listModel) selectedLogID() string {
	cursor := m.table.Cursor()
	if cursor < 0 || cursor >= len(m.filteredLogs) {
		return ""
	}
	return m.filteredLogs[cursor].ID
}

// statusText returns plain (unstyled) status text for table display.
func statusText(status string) string {
	switch strings.ToLower(status) {
	case "accepted":
		return "✓ accepted"
	case "rejected":
		return "✗ rejected"
	case "quit":
		return "- quit"
	default:
		return status
	}
}

// shortTimeAgo returns a compact relative time string for table display.
func shortTimeAgo(t time.Time) string {
	d := time.Since(t)

	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		m := int(d.Minutes())
		return fmt.Sprintf("%dm ago", m)
	case d < 24*time.Hour:
		h := int(d.Hours())
		return fmt.Sprintf("%dh ago", h)
	case d < 7*24*time.Hour:
		days := int(d.Hours() / 24)
		return fmt.Sprintf("%dd ago", days)
	default:
		weeks := int(d.Hours() / 24 / 7)
		return fmt.Sprintf("%dw ago", weeks)
	}
}
