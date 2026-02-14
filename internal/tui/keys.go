package tui

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Quit         key.Binding
	Help         key.Binding
	Enter        key.Binding
	Back         key.Binding
	Search       key.Binding
	StatusFilter key.Binding
	NextTab      key.Binding
	PrevTab      key.Binding
	Escape       key.Binding
	Copy         key.Binding
}

func newKeyMap() keyMap {
	return keyMap{
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "view details"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc", "backspace"),
			key.WithHelp("esc", "back to list"),
		),
		Search: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "search"),
		),
		StatusFilter: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "cycle status filter"),
		),
		NextTab: key.NewBinding(
			key.WithKeys("tab", "l"),
			key.WithHelp("tab/l", "next tab"),
		),
		PrevTab: key.NewBinding(
			key.WithKeys("shift+tab", "h"),
			key.WithHelp("shift+tab/h", "prev tab"),
		),
		Escape: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "clear search"),
		),
		Copy: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "copy"),
		),
	}
}

// ShortHelp returns keybindings to show in the mini help view.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Enter, k.Back, k.Search, k.StatusFilter},
		{k.NextTab, k.PrevTab, k.Copy},
		{k.Help, k.Quit},
	}
}
