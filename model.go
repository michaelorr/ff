package main

import (
	"time"

	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"

	"github.com/michaelorr/ff/components"
	"github.com/michaelorr/ff/search"
	"github.com/michaelorr/ff/state"
)

const debounceDuration = 250 * time.Millisecond

type debounceMsg int

type model struct {
	input            textinput.Model
	scanner          *search.Scanner
	previewOpen      bool
	matchedFileNames []string
	matchedFileIcons map[search.FileIcon]int
	matchesByFile    map[string][]search.ContentMatch
	flatEntries      []components.MatchEntry
	selectedMatchIdx int
	filtersViewport  viewport.Model
	matchesViewport  viewport.Model
	previewViewport  viewport.Model
	width, height    int
	generation       int
}

func newModel() model {
	m := model{
		input:            components.NewSearchInput(),
		scanner:          search.NewScanner(),
		previewOpen:      true,
		matchedFileIcons: make(map[search.FileIcon]int),
		matchesByFile:    make(map[string][]search.ContentMatch),
		selectedMatchIdx: -1,
		filtersViewport:  viewport.New(),
		matchesViewport:  viewport.New(),
		previewViewport:  viewport.New(),
		generation:       0,
	}

	if s, ok := state.Load(); ok {
		m.input.SetValue(s.Query)
	}

	return m
}

func (m model) Init() tea.Cmd {
	search.StartWalking(".")

	cmds := []tea.Cmd{m.scanner.NextCmd()}

	if m.input.Value() != "" {
		cmds = append(cmds, m.searchCmd())
	}

	return tea.Batch(cmds...)
}

func (m model) State() state.AppState {
	return state.AppState{
		Query: m.input.Value(),
	}
}

func (m model) searchCmd() tea.Cmd {
	return func() tea.Msg {
		m.scanner.Search(m.input.Value(), m.generation)
		return nil
	}
}
