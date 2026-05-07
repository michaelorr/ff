package main

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/michaelorr/ff/panel"
)

func (m model) View() tea.View {
	var v tea.View
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion

	v.SetContent(
		lipgloss.JoinVertical(
			lipgloss.Left,
			searchView(m),
			lipgloss.JoinHorizontal(
				lipgloss.Left,
				filterView(m),
				resultsView(m),
				previewView(m),
			),
		),
	)

	return v
}

type foo struct {
	s string
}

func (f foo) View() string {
	return f.s
}

var (
	searchHeight = 3
	filterWidth  = 20
)

func searchView(m model) string {
	return panel.Render("search", m.width, searchHeight, m.input, m.mode == InsertMode)
}

func filterView(m model) string {
	return panel.Render("filters", filterWidth, m.height-searchHeight, foo{"foo"}, false)
}

func resultsView(m model) string {
	return panel.Render("results", resultsWidth(m), m.height-searchHeight, foo{"bar"}, false)
}

func previewView(m model) string {
	return panel.Render("preview", previewWidth(m), m.height-searchHeight, foo{"baz"}, false)
}

func resultsWidth(m model) int {
	return ((m.width - filterWidth) / 2) + (m.width-filterWidth)%2
}

func previewWidth(m model) int {
	return (m.width - filterWidth) / 2
}
