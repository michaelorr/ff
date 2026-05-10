package main

import (
	"time"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/viewport"

	"charm.land/bubbles/v2/textinput"
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

func newModel() model {
	m := model{
		input:            components.NewSearchInput(),
		scanner:          search.NewScanner(),
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

	return m
}

func (m model) Init() tea.Cmd {
	search.StartWalking(".")

	cmds := []tea.Cmd{m.scanner.NextCmd()}

	// If query is non-empty after loading state, immediately do a search
	query, gen := m.input.Value(), m.generation
	if query != "" {
		cmds = append(cmds, func() tea.Msg {
			m.scanner.Search(query, gen)
			return nil
		})
	}

	return tea.Batch(cmds...)
}
