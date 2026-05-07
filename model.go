package main

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/viewport"
	"charm.land/lipgloss/v2"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"

	"github.com/michaelorr/ff/colors"
	"github.com/michaelorr/ff/internal/search"
	"github.com/michaelorr/ff/internal/state"
)

type model struct {
	input           textinput.Model
	scanner         *search.Scanner
	mode            string
	help            help.Model
	results         []string
	resultsViewport viewport.Model
	width, height   int
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
		scanner:         *search.NewScanner(),
		mode:            InsertMode,
		help:            help.New(),
		results:         []string{},
		resultsViewport: viewport.New(),
	}

	if s, ok := state.Load(); ok {
		m.input.SetValue(s.Query)
		m.results = s.Results
	}

	m.input.Focus()

	return m
}

func (m model) Init() tea.Cmd {
	search.StartWalking()
	return m.scanner.NextCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyPressMsg:
		if m, cmd, done := handleKeyPressMsg(&m, msg); done {
			return m, cmd
		}
	}

	var inputCmd tea.Cmd
	m.input, inputCmd = m.input.Update(msg)
	return m, inputCmd
}
