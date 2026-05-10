package main

import (
	"slices"

	"github.com/michaelorr/ff/search"
)

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
}

func (m *model) resetMatches() {
	m.matchedFileNames = nil
	m.matchedFileIcons = make(map[search.FileIcon]int)
	m.matchesByFile = make(map[string][]search.ContentMatch)
}
