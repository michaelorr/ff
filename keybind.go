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
	switch {
	case key.Matches(msg, keys.Quit):
		return *m, tea.Quit, true
	case key.Matches(msg, keys.Reset):
		return handleResetKey(m)
	case key.Matches(msg, keys.TogglePreview):
		return handleTogglePreviewKey(m)
	default:
		return *m, nil, false
	}
}

func handleResetKey(m *model) (model, tea.Cmd, bool) {
	m.generation++
	m.input.Reset()
	m.resetMatches()
	m.renderLayout()
	return *m, state.SaveCmd(m.State()), true
}

func handleTogglePreviewKey(m *model) (model, tea.Cmd, bool) {
	m.previewOpen = !m.previewOpen
	m.renderLayout()
	return *m, nil, true
}
