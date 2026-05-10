package main

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

type keyMap struct {
	Help          key.Binding
	Quit          key.Binding
	Reset         key.Binding
	TogglePreview key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Reset, k.Help}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Help},
		{k.Quit},
	}
}

var keys = keyMap{
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("esc", "quit"),
	),
	Reset: key.NewBinding(
		key.WithKeys("ctrl+r"),
		key.WithHelp("ctrl+r", "reset"),
	),
	TogglePreview: key.NewBinding(
		key.WithKeys("ctrl+p"),
		key.WithHelp("ctrl+p", "toggle preview"),
	),
}

func handleKeyPressMsg(m *model, msg tea.KeyPressMsg) (model, tea.Cmd, bool) {
	var done bool

	switch {
	case key.Matches(msg, keys.Help):
		m.help.ShowAll = !m.help.ShowAll
		done = true
	case key.Matches(msg, keys.Quit):
		return *m, tea.Quit, true
	case key.Matches(msg, keys.Reset):
		m.input.SetValue("")
		m.generation++
		m.resetMatches()
		m.renderLayout()
		return *m, nil, true
	case key.Matches(msg, keys.TogglePreview):
		m.previewOpen = !m.previewOpen
		m.renderLayout()
		done = true
	}

	return *m, nil, done
}
