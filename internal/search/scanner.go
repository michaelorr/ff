package search

import tea "charm.land/bubbletea/v2"

type Scanner struct{}

func NewScanner() *Scanner {
	return &Scanner{}
}

func (s Scanner) NextCmd() tea.Cmd {
	return nil
}
