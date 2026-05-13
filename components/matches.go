// Package components provides helpers for rendering UI components
package components

import (
	"fmt"
	"image/color"
	"path/filepath"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/michaelorr/ff/search"
	"github.com/michaelorr/ff/style"
)

// MatchEntry is a navigable item in the matches pane — either a file or a match line.
type MatchEntry struct {
	Path  string
	Match *search.ContentMatch
}

var (
	highlightCache = map[color.Color]map[string]map[int]string{}

	lineNumStyle = style.Default.Foreground(style.Gray0)
	selectedBg   = style.BgDimRed
)

func Matches(
	files []string,
	byFile map[string][]search.ContentMatch,
	selectedEntry *MatchEntry,
	width, startLine, numLines int,
) string {
	endLine := startLine + numLines
	currentLine := 0
	var b strings.Builder

	for _, path := range files {
		if currentLine >= endLine {
			break
		}

		matches := byFile[path]

		// skip to the first file that will actually draw something
		fileEnd := currentLine + 2 + len(matches)
		if fileEnd <= startLine {
			currentLine = fileEnd
			continue
		}

		// blank separator line
		if currentLine >= startLine {
			b.WriteByte('\n')
		}
		currentLine++

		// file header line
		if currentLine >= startLine && currentLine < endLine {
			currentFileSelected := isSelectedFile(selectedEntry, path)

			icon := matches[0].Icon
			dir, file := filepath.Split(path)
			writeLine(
				&b, width,
				fileLineStyles(currentFileSelected),
				[]styledString{
					{Style: iconStyles(icon.Color, currentFileSelected), Text: icon.Icon + " "},
					{Style: fileDirStyles(currentFileSelected), Text: dir},
					{Style: fileNameStyles(currentFileSelected), Text: file},
				},
			)
		}
		currentLine++

		// match lines
		for _, m := range matches {
			if currentLine >= endLine {
				break
			}
			if currentLine >= startLine {
				currentMatchSelected := isSelectedMatch(selectedEntry, path, m.LineNum)
				writeLine(
					&b, width,
					whitespaceStyle(currentMatchSelected),
					[]styledString{
						{Style: lineNumStyles(currentMatchSelected), Text: fmt.Sprintf("%5d ", m.LineNum)},
						{StyledText: cachedHighlight(m.Line, path, m.LineNum, syntaxBgColors(currentMatchSelected))},
					},
				)
			}
			currentLine++
		}
	}

	return b.String()
}

func lineNumStyles(selected bool) lipgloss.Style {
	if selected {
		return lineNumStyle.Background(selectedBg).Foreground(style.Accent).Bold(true)
	}
	return lineNumStyle
}

func syntaxBgColors(selected bool) color.Color {
	if selected {
		return selectedBg
	}
	return style.Bg0
}

func iconStyles(color string, selected bool) lipgloss.Style {
	return fileLineStyles(selected).Foreground(lipgloss.Color(color))
}

func isSelectedFile(selectedEntry *MatchEntry, path string) bool {
	return selectedEntry != nil && selectedEntry.Match == nil && selectedEntry.Path == path
}

func isSelectedMatch(selectedEntry *MatchEntry, path string, lineNum int) bool {
	return selectedEntry != nil && selectedEntry.Match != nil && selectedEntry.Path == path && selectedEntry.Match.LineNum == lineNum
}

func fileDirStyles(selected bool) lipgloss.Style {
	fg := style.Gray0
	if selected {
		fg = style.Fg0
	}
	return fileLineStyles(selected).Foreground(fg)
}

func fileNameStyles(selected bool) lipgloss.Style {
	return fileLineStyles(selected).Bold(true).Foreground(style.Accent)
}

func fileLineStyles(selected bool) lipgloss.Style {
	fileLineStyle := style.Default.
		Underline(true).
		UnderlineSpaces(true).
		UnderlineColor(style.Gray0).
		UnderlineStyle(lipgloss.UnderlineDashed)

	if selected {
		return fileLineStyle.Background(selectedBg).UnderlineColor(style.Accent)
	}

	return fileLineStyle
}

func whitespaceStyle(selected bool) lipgloss.Style {
	if selected {
		return style.Default.Background(selectedBg)
	}
	return style.Default
}

func cachedHighlight(line, path string, lineNum int, bg color.Color) string {
	if highlightCache[bg] == nil {
		highlightCache[bg] = map[string]map[int]string{}
	}
	if highlightCache[bg][path] == nil {
		highlightCache[bg][path] = map[int]string{}
	}
	if cached, ok := highlightCache[bg][path][lineNum]; ok {
		return cached
	}
	result := syntaxHighlightWithBg(line, path, bg)
	highlightCache[bg][path][lineNum] = result
	return result
}
