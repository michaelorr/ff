package main

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"github.com/michaelorr/ff/state"
)

type keyMap struct {
	Quit          key.Binding
	Reset         key.Binding
	TogglePreview key.Binding
}

var keys = keyMap{
	Quit:          key.NewBinding(key.WithKeys("ctrl+c")),
	Reset:         key.NewBinding(key.WithKeys("ctrl+r")),
	TogglePreview: key.NewBinding(key.WithKeys("ctrl+p")),
}

func handleKeyPressMsg(m *model, msg tea.KeyPressMsg) (model, tea.Cmd, bool) {
	var done bool

	switch {
	case key.Matches(msg, keys.Quit):
		return *m, tea.Quit, true
	case key.Matches(msg, keys.Reset):
		m.generation++
		m.input.Reset()
		m.resetMatches()
		m.renderLayout()
		return *m, state.SaveCmd(m.State()), true
	case key.Matches(msg, keys.TogglePreview):
		m.previewOpen = !m.previewOpen
		m.renderLayout()
		done = true
	}

	return *m, nil, done
}
