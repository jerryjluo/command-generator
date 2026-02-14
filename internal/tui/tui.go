package tui

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/jerryluo/cmd/internal/clipboard"
	"github.com/jerryluo/cmd/internal/logging"
)

type viewState int

const (
	listView viewState = iota
	detailView
)

type model struct {
	state         viewState
	list          listModel
	detail        detailModel
	keys          keyMap
	help          help.Model
	width         int
	height        int
	ready         bool
	showHelp      bool
	statusMessage string
}

// clipboardCopyMsg is sent after a clipboard copy attempt.
type clipboardCopyMsg struct {
	err error
}

// clearStatusMsg clears the status message after a delay.
type clearStatusMsg struct{}

func newModel() model {
	return model{
		state:  listView,
		list:   newListModel(),
		detail: newDetailModel(),
		keys:   newKeyMap(),
		help:   help.New(),
	}
}

func (m model) Init() tea.Cmd {
	return m.list.loadLogs()
}

// copyToClipboard copies text to the system clipboard.
func copyToClipboard(text string) tea.Cmd {
	return func() tea.Msg {
		return clipboardCopyMsg{err: clipboard.Copy(text)}
	}
}

// copyCommandFromLog loads a log by ID and copies its command to clipboard.
func copyCommandFromLog(id string) tea.Cmd {
	return func() tea.Msg {
		log, err := logging.ReadLog(id)
		if err != nil {
			return clipboardCopyMsg{err: err}
		}
		if len(log.Iterations) == 0 {
			return clipboardCopyMsg{err: fmt.Errorf("no command available")}
		}
		cmd := log.Iterations[len(log.Iterations)-1].ModelOutput.Command
		return clipboardCopyMsg{err: clipboard.Copy(cmd)}
	}
}

// clearStatusAfter returns a command that clears the status message after a delay.
func clearStatusAfter(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(time.Time) tea.Msg {
		return clearStatusMsg{}
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		m.help.Width = msg.Width
		m.list.showHelp = m.showHelp
		m.list.SetSize(msg.Width, msg.Height)
		m.detail.showHelp = m.showHelp
		m.detail.SetSize(msg.Width, msg.Height)
		return m, nil

	case clipboardCopyMsg:
		if msg.err != nil {
			m.statusMessage = "Copy failed: " + msg.err.Error()
		} else {
			m.statusMessage = "Copied!"
		}
		return m, clearStatusAfter(2 * time.Second)

	case clearStatusMsg:
		m.statusMessage = ""
		return m, nil

	case tea.KeyMsg:
		// Global quit (but not when searching in list view)
		if key.Matches(msg, m.keys.Quit) && !(m.state == listView && m.list.searching) {
			return m, tea.Quit
		}

		// Global help toggle (but not when searching in list view)
		if key.Matches(msg, m.keys.Help) && !(m.state == listView && m.list.searching) {
			m.showHelp = !m.showHelp
			switch m.state {
			case listView:
				m.list.showHelp = m.showHelp
				m.list.SetSize(m.width, m.height)
			case detailView:
				m.detail.showHelp = m.showHelp
				m.detail.SetSize(m.width, m.height)
			}
			return m, nil
		}

		// View-specific navigation
		switch m.state {
		case listView:
			if key.Matches(msg, m.keys.Enter) && !m.list.searching {
				id := m.list.selectedLogID()
				if id != "" {
					m.state = detailView
					m.detail.loading = true
					m.detail.err = nil
					m.detail.log = nil
					return m, m.detail.loadLog(id)
				}
				return m, nil
			}
			if key.Matches(msg, m.keys.Copy) && !m.list.searching {
				id := m.list.selectedLogID()
				if id != "" {
					return m, copyCommandFromLog(id)
				}
				return m, nil
			}

		case detailView:
			if key.Matches(msg, m.keys.Back) {
				m.state = listView
				m.detail.log = nil
				m.detail.ready = false
				return m, nil
			}
			if key.Matches(msg, m.keys.Copy) {
				content := m.detail.copyableContent()
				if content != "" {
					return m, copyToClipboard(content)
				}
				return m, nil
			}
		}
	}

	// Route to active view
	var cmd tea.Cmd
	switch m.state {
	case listView:
		m.list, cmd = m.list.Update(msg)
	case detailView:
		m.detail, cmd = m.detail.Update(msg)
	}

	return m, cmd
}

func (m model) View() string {
	if !m.ready {
		return "Loading..."
	}

	switch m.state {
	case listView:
		m.list.showHelp = m.showHelp
		m.list.statusMessage = m.statusMessage
		return m.list.View()
	case detailView:
		m.detail.showHelp = m.showHelp
		m.detail.statusMessage = m.statusMessage
		return m.detail.View()
	}

	return ""
}

// Run starts the TUI log viewer.
func Run() {
	p := tea.NewProgram(newModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
		os.Exit(1)
	}
}
