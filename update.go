package main

import (
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/michaelorr/ff/search"
	"github.com/michaelorr/ff/state"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		return handleWindowSizeMsg(m, msg)
	case tea.KeyPressMsg:
		m, cmd, done := handleKeyPressMsg(&m, msg)
		if done {
			return m, cmd
		}
	case debounceMsg:
		return handleDebounceMsg(m, msg)
	case search.ContentBatchMsg:
		return handleContentBatchMsg(m, msg)
	}

	return sendMsgToInput(m, msg)
}

func handleWindowSizeMsg(m model, msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	m.width = msg.Width
	m.height = msg.Height
	m.renderLayout()

	return m, nil
}

func handleDebounceMsg(m model, msg debounceMsg) (tea.Model, tea.Cmd) {
	// The value of `debounceMsg` corresponds to the `generation` from when the message was created.
	// If this matches the current generation, then the input hasn't changed and we should do a search.
	// Otherwise, drop the message because it's outdated.
	if int(msg) == m.generation {
		m.resetMatches()
		m.renderLayout()
		if m.input.Value() != "" {
			return m, m.searchCmd()
		}
	}
	return m, nil
}

func handleContentBatchMsg(m model, msg search.ContentBatchMsg) (tea.Model, tea.Cmd) {
	m.addToMatches(msg.Matches)
	m.renderLayout()
	return m, m.scanner.NextCmd()
}

func sendMsgToInput(m model, msg tea.Msg) (model, tea.Cmd) {
	// Compare the value before and after the update.
	// If it changed then save the state and debounce a search.
	prev := m.input.Value()
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	cmds := []tea.Cmd{cmd}
	if prev != m.input.Value() {
		state.Save(m.State())
		m.generation++
		debounceCmd := tea.Tick(debounceDuration, func(_ time.Time) tea.Msg {
			return debounceMsg(m.generation)
		})
		cmds = append(cmds, debounceCmd)
	}
	return m, tea.Batch(cmds...)
}
