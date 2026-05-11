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
	NextMatch     key.Binding
	NextFile      key.Binding
	PrevMatch     key.Binding
	PrevFile      key.Binding
}

var keys = keyMap{
	Quit:          key.NewBinding(key.WithKeys("ctrl+c")),
	Reset:         key.NewBinding(key.WithKeys("ctrl+r")),
	TogglePreview: key.NewBinding(key.WithKeys("ctrl+p")),
	NextMatch:     key.NewBinding(key.WithKeys("ctrl+j")),
	NextFile:      key.NewBinding(key.WithKeys("alt+j")),
	PrevMatch:     key.NewBinding(key.WithKeys("ctrl+k")),
	PrevFile:      key.NewBinding(key.WithKeys("alt+k")),
}

func handleKeyPressMsg(m *model, msg tea.KeyPressMsg) (model, tea.Cmd, bool) {
	switch {
	case key.Matches(msg, keys.Quit):
		return *m, tea.Quit, true
	case key.Matches(msg, keys.Reset):
		return handleResetKey(m)
	case key.Matches(msg, keys.TogglePreview):
		return handleTogglePreviewKey(m)
	case key.Matches(msg, keys.NextMatch):
		return handleNextMatchKey(m)
	case key.Matches(msg, keys.NextFile):
		return handleNextFileKey(m)
	case key.Matches(msg, keys.PrevMatch):
		return handlePrevMatchKey(m)
	case key.Matches(msg, keys.PrevFile):
		return handlePrevFileKey(m)
	default:
		return *m, nil, false
	}
}

func handleResetKey(m *model) (model, tea.Cmd, bool) {
	m.generation++
	m.input.Reset()
	m.resetMatches()
	m.updateDimensions()
	return *m, state.SaveCmd(m.State()), true
}

func handleTogglePreviewKey(m *model) (model, tea.Cmd, bool) {
	m.previewOpen = !m.previewOpen
	m.updateDimensions()
	return *m, nil, true
}

func handleNextFileKey(m *model) (model, tea.Cmd, bool) {
	for i := m.selectedMatchIdx + 1; i < len(m.flatEntries); i++ {
		if m.flatEntries[i].Match == nil {
			m.selectedMatchIdx = i
			break
		}
	}
	return *m, nil, true
}

func handlePrevFileKey(m *model) (model, tea.Cmd, bool) {
	if m.selectedMatchIdx == 0 {
		return *m, nil, true
	}
	for i := m.selectedMatchIdx - 1; i >= 0; i-- {
		if m.flatEntries[i].Match == nil {
			m.selectedMatchIdx = i
			break
		}
	}
	return *m, nil, true
}

func handleNextMatchKey(m *model) (model, tea.Cmd, bool) {
	if len(m.flatEntries) != 0 && m.selectedMatchIdx < len(m.flatEntries)-1 {
		m.selectedMatchIdx++
	}
	return *m, nil, true
}

func handlePrevMatchKey(m *model) (model, tea.Cmd, bool) {
	if m.selectedMatchIdx > 0 {
		m.selectedMatchIdx--
	}
	return *m, nil, true
}
