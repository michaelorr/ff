package main

import (
	"slices"

	"github.com/michaelorr/ff/components"
	"github.com/michaelorr/ff/search"
)

func (m model) selectedEntry() *components.MatchEntry {
	if m.selectedMatchIdx < 0 {
		return nil
	}
	e := m.flatEntries[m.selectedMatchIdx]
	return &e
}

func (m *model) addToMatches(matches []search.ContentMatch) {
	for _, match := range matches {
		if match.Generation == m.generation {
			if _, seenFile := m.matchesByFile[match.Path]; !seenFile {
				m.matchedFileNames = append(m.matchedFileNames, match.Path)
				m.matchesByFile[match.Path] = nil
				m.matchedFileIcons[match.Icon]++

				slices.Sort(m.matchedFileNames)
			}

			m.matchesByFile[match.Path] = append(m.matchesByFile[match.Path], match)
		}
	}
	m.rebuildFlatEntries()
}

func (m *model) rebuildFlatEntries() {
	m.flatEntries = m.flatEntries[:0]
	for _, path := range m.matchedFileNames {
		m.flatEntries = append(m.flatEntries, components.MatchEntry{Path: path})
		for _, match := range m.matchesByFile[path] {
			m.flatEntries = append(m.flatEntries, components.MatchEntry{Path: path, Match: &match})
		}
	}
}

func (m *model) resetMatches() {
	m.matchedFileNames = nil
	m.matchedFileIcons = make(map[search.FileIcon]int)
	m.matchesByFile = make(map[string][]search.ContentMatch)
	m.flatEntries = nil
	m.selectedMatchIdx = -1
}
