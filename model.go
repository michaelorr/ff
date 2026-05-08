package main

import (
	"time"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/viewport"
	"charm.land/lipgloss/v2"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"

	"github.com/michaelorr/ff/colors"
	"github.com/michaelorr/ff/search"
	"github.com/michaelorr/ff/state"
)

const debounceDuration = 250 * time.Millisecond

type debounceMsg int

type model struct {
	input            textinput.Model
	scanner          *search.Scanner
	mode             string
	previewOpen      bool
	help             help.Model
	matchedFileNames []string
	matchedFileIcons map[search.FileIcon]int
	matchesByFile    map[string][]search.ContentMatch
	filtersViewport  viewport.Model
	matchesViewport  viewport.Model
	previewViewport  viewport.Model
	width, height    int
	generation       int
}

func newSearchInput() textinput.Model {
	input := textinput.New()
	input.Prompt = "❯ "

	// Styles.[Focused|Blurred|Cursor].[Text|Placeholder|Suggestion|Prompt]
	s := textinput.DefaultDarkStyles()
	s.Focused.Prompt = lipgloss.NewStyle().Foreground(colors.BgBlue)
	s.Blurred.Prompt = s.Focused.Prompt
	s.Focused.Text = lipgloss.NewStyle().Foreground(colors.Fg0).Background(colors.Bg0)
	s.Blurred.Text = s.Focused.Text
	s.Cursor.Color = colors.Gray2
	input.SetStyles(s)
	return input
}

const (
	InsertMode  = "insert"
	CommandMode = "command"
)

func newModel() model {
	m := model{
		input:            newSearchInput(),
		scanner:          search.NewScanner(),
		mode:             InsertMode,
		previewOpen:      true,
		help:             help.New(),
		matchedFileIcons: make(map[search.FileIcon]int),
		matchesByFile:    make(map[string][]search.ContentMatch),
		filtersViewport:  viewport.New(),
		matchesViewport:  viewport.New(),
		previewViewport:  viewport.New(),
		generation:       0,
	}

	if s, ok := state.Load(); ok {
		m.input.SetValue(s.Query)
	}

	m.input.Focus()

	return m
}

func (m model) Init() tea.Cmd {
	search.StartWalking(".")
	cmds := []tea.Cmd{m.scanner.NextCmd()}

	query, gen := m.input.Value(), m.generation
	if query != "" {
		cmds = append(cmds, func() tea.Msg {
			m.scanner.Search(query, gen)
			return nil
		})
	}

	return tea.Batch(cmds...)
}

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
