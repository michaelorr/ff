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
	m.matchesYOffset = 0
}

// lineForMatchIdx returns the 0-based line number in the matches viewport content
// that corresponds to flatEntries[idx]. Each file entry occupies 2 lines (a blank
// separator followed by the header), and each match entry occupies 1 line.
func (m *model) lineForMatchIdx(idx int) int {
	line := 0
	for i, entry := range m.flatEntries {
		if i == idx {
			if entry.Match == nil {
				// entry is a file, add 1 for the blank line between files
				return line + 1
			}
			return line
		}
		if entry.Match == nil {
			line += 2
		} else {
			line++
		}
	}
	return line
}

const scrolloff = 3

func (m *model) updateMatchesViewport() {
	target := m.lineForMatchIdx(m.selectedMatchIdx)
	height := m.matchesViewport.Height()
	offset := m.matchesYOffset

	if target < offset+scrolloff {
		m.matchesYOffset = max(0, target-scrolloff)
	} else if target >= offset+height-scrolloff {
		m.matchesYOffset = target - height + scrolloff + 1
	}
}
