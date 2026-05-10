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
		m.width = msg.Width
		m.height = msg.Height

		m.renderLayout()

		return m, nil

	case tea.KeyPressMsg:
		if m, cmd, done := handleKeyPressMsg(&m, msg); done {
			_ = state.Save(state.AppState{Query: m.input.Value()})
			return m, cmd
		}

	case search.ContentBatchMsg:
		m.addToMatches(msg.Matches)
		m.renderLayout()
		return m, m.scanner.NextCmd()

	case debounceMsg:
		// The value of `debounceMsg` corresponds to the `generation` from when the message was created.
		// If this matches the current generation, then the input hasn't changed and we should do a search.
		if int(msg) == m.generation {
			m.resetMatches()
			m.renderLayout()
			query, gen := m.input.Value(), m.generation
			if query != "" {
				return m, func() tea.Msg {
					m.scanner.Search(query, gen)
					return nil
				}
			}
		}
		return m, nil
	}

	// None of the specific message types matched, send to the textinput component
	prev := m.input.Value()
	var inputCmd tea.Cmd
	m.input, inputCmd = m.input.Update(msg)
	if m.input.Value() != prev {
		m.generation++
		deferredDebounceCmd := tea.Tick(debounceDuration, func(_ time.Time) tea.Msg {
			return debounceMsg(m.generation)
		})

		return m, tea.Batch(inputCmd, deferredDebounceCmd)
	}
	_ = state.Save(state.AppState{Query: m.input.Value()})
	return m, inputCmd
}
