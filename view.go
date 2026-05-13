package main

import (
	"path/filepath"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/michaelorr/ff/components"
)

var (
	borders           = 2
	bordersAndPadding = borders + 2
)

func (m model) View() tea.View {
	v := tea.View{}
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion

	m.filtersViewport.SetContent(components.Filters(m.matchedFileIcons, m.filtersViewport.Width()))
	m.matchesViewport.SetContent(components.Matches(m.matchedFileNames, m.matchesByFile, m.selectedEntry(), m.matchesViewport.Width(), m.matchesYOffset, m.matchesViewport.Height()))
	if m.previewOpen {
		m.previewViewport.SetContent(components.Preview(m.selectedEntry(), m.previewViewport.Width(), m.previewViewport.Height()))
	}

	v.SetContent(
		lipgloss.JoinVertical(
			lipgloss.Left, m.searchPanel(), lipgloss.JoinHorizontal(
				lipgloss.Top, m.filterPanel(), m.matchesPanel(), m.previewPanel(),
			),
		),
	)

	return v
}

func (m *model) updateDimensions() {
	m.filtersViewport.SetHeight(m.componentSize("filters").Height - borders)
	m.filtersViewport.SetWidth(m.componentSize("filters").Width - bordersAndPadding)

	m.matchesViewport.SetHeight(m.componentSize("matches").Height - borders)
	m.matchesViewport.SetWidth(m.componentSize("matches").Width - bordersAndPadding)

	m.previewViewport.SetHeight(m.componentSize("preview").Height - borders)
	m.previewViewport.SetWidth(m.componentSize("preview").Width - bordersAndPadding)
}

func (m model) searchPanel() string {
	return components.Panel("search", "", m.componentSize("search"), m.input, true)
}

func (m model) filterPanel() string {
	return components.Panel("filters", "", m.componentSize("filters"), &m.filtersViewport, false)
}

func (m model) matchesPanel() string {
	return components.Panel("matches", "", m.componentSize("matches"), &m.matchesViewport, false)
}

func (m model) previewPanel() string {
	if !m.previewOpen {
		return ""
	}

	var subtitle string
	if e := m.selectedEntry(); e != nil {
		subtitle = filepath.Base(e.Path)
	}

	return components.Panel("preview", subtitle, m.componentSize("preview"), &m.previewViewport, false)
}

func (m model) componentSize(component string) components.Size {
	searchHeight := 3
	filterWidth := 20
	bodyComponentHeight := m.height - searchHeight

	previewWidth := func() int {
		if m.previewOpen {
			return (m.width - filterWidth) / 2
		}
		return 0
	}

	components := map[string]components.Size{
		"search": {
			Width:  m.width,
			Height: searchHeight,
		},
		"filters": {
			Width:  filterWidth,
			Height: bodyComponentHeight,
		},
		"matches": {
			Width:  m.width - filterWidth - previewWidth(),
			Height: bodyComponentHeight,
		},
		"preview": {
			Width:  previewWidth(),
			Height: bodyComponentHeight,
		},
	}

	return components[component]
}
