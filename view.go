package main

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/michaelorr/ff/components"
)

func (m model) View() tea.View {
	v := tea.View{}
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion

	v.SetContent(
		lipgloss.JoinVertical(
			lipgloss.Left,
			searchPanel(m),
			lipgloss.JoinHorizontal(
				lipgloss.Left,
				filterPanel(m),
				matchesPanel(m),
				previewPanel(m),
			),
		),
	)

	return v
}

var (
	searchHeight = 3
	filterWidth  = 20
)

func (m *model) renderLayout() {
	m.filtersViewport.SetHeight(m.height - searchHeight - 2)
	m.filtersViewport.SetWidth(filterWidth - 2)
	m.matchesViewport.SetHeight(m.height - searchHeight - 2)
	m.matchesViewport.SetWidth(matchesWidth(*m) - 4)
	m.previewViewport.SetHeight(m.height - searchHeight - 2)
	m.previewViewport.SetWidth(previewWidth(*m) - 4)

	m.filtersViewport.SetContent(components.RenderFilters(m.matchedFileIcons))
	m.matchesViewport.SetContent(components.RenderMatches(m.matchedFileNames, m.matchesByFile, matchesWidth(*m)))
	m.previewViewport.SetContent("baz")
}

func searchPanel(m model) string {
	return components.RenderPanel("search", m.width, searchHeight, m.input, m.mode == InsertMode)
}

func filterPanel(m model) string {
	return components.RenderPanel("filters", filterWidth, m.height-searchHeight, &m.filtersViewport, false)
}

func matchesPanel(m model) string {
	return components.RenderPanel("matches", matchesWidth(m), m.height-searchHeight, &m.matchesViewport, false)
}

func previewPanel(m model) string {
	if !m.previewOpen {
		return ""
	}

	return components.RenderPanel("preview", previewWidth(m), m.height-searchHeight, &m.previewViewport, false)
}

func matchesWidth(m model) int {
	if !m.previewOpen {
		return m.width - filterWidth
	}
	return ((m.width - filterWidth) / 2) + (m.width-filterWidth)%2
}

func previewWidth(m model) int {
	if !m.previewOpen {
		return 0
	}

	return (m.width - filterWidth) / 2
}
