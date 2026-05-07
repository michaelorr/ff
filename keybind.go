package main

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

type keyMap struct {
	Help        key.Binding
	Quit        key.Binding
	CommandMode key.Binding
	InsertMode  key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help}
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
	CommandMode: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "command mode"),
	),
	InsertMode: key.NewBinding(
		key.WithKeys("i"),
		key.WithHelp("i", "insert mode"),
	),
}

func handleKeyPressMsg(m *model, msg tea.KeyPressMsg) (model, tea.Cmd, bool) {
	var done bool

	switch {
	case key.Matches(msg, keys.Help):
		if m.mode == CommandMode {
			m.help.ShowAll = !m.help.ShowAll
			done = true
		}
	case key.Matches(msg, keys.Quit):
		return *m, tea.Quit, true
	case key.Matches(msg, keys.CommandMode):
		if m.mode == InsertMode {
			m.input.Blur()
			m.mode = CommandMode
			done = true
		}
	case key.Matches(msg, keys.InsertMode):
		if m.mode == CommandMode {
			m.input.Focus()
			m.mode = InsertMode
			done = true
		}
	}

	return *m, nil, done
}
