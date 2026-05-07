package main

import (
	"slices"
	"time"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/viewport"
	"charm.land/lipgloss/v2"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"

	"github.com/michaelorr/ff/colors"
	"github.com/michaelorr/ff/internal/search"
	"github.com/michaelorr/ff/internal/state"
	"github.com/michaelorr/ff/matches"
)

const debounceDuration = 250 * time.Millisecond

type debounceMsg int

type model struct {
	input           textinput.Model
	scanner         *search.Scanner
	mode            string
	help            help.Model
	matchedFiles    []string
	matchesByFile   map[string][]search.ContentMatch
	matchesViewport viewport.Model
	width, height   int
	generation      int
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
		input:           newSearchInput(),
		scanner:         search.NewScanner(),
		mode:            InsertMode,
		help:            help.New(),
		matchesByFile:   make(map[string][]search.ContentMatch),
		matchesViewport: viewport.New(),
		generation:      0,
	}

	if s, ok := state.Load(); ok {
		m.input.SetValue(s.Query)
	}

	m.input.Focus()

	return m
}

func (m model) Init() tea.Cmd {
	search.StartWalking(".")
	return m.scanner.NextCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.matchesViewport.SetWidth(matchesWidth(m) - 4)
		m.matchesViewport.SetHeight(m.height - searchHeight - 2)
		return m, nil

	case tea.KeyPressMsg:
		if m, cmd, done := handleKeyPressMsg(&m, msg); done {
			return m, cmd
		}

	case search.ContentBatchMsg:
		for _, match := range msg.Matches {
			if match.Generation == m.generation {
				if _, seen := m.matchesByFile[match.Path]; !seen {
					m.matchedFiles = append(m.matchedFiles, match.Path)
					slices.Sort(m.matchedFiles)
					m.matchesByFile[match.Path] = nil
				}
				m.matchesByFile[match.Path] = append(m.matchesByFile[match.Path], match)
			}
		}
		m.matchesViewport.SetContent(matches.RenderMatches(m.matchedFiles, m.matchesByFile))
		return m, m.scanner.NextCmd()

	case debounceMsg:
		// The value of `debounceMsg` corresponds to the `generation` from when the message was created.
		// If this matches the current generation, then the input hasn't changed and we should do a search.
		if int(msg) == m.generation {
			m.matchedFiles = nil
			m.matchesByFile = make(map[string][]search.ContentMatch)
			query, gen := m.input.Value(), m.generation
			return m, func() tea.Msg {
				m.scanner.Search(query, gen)
				return nil
			}
		}
		return m, nil
	}

	// None of the other more specific message types matched, try sending to the textinput component.

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
	return m, inputCmd
}
