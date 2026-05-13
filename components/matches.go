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

	fileLineStyle = style.Default.
			Underline(true).
			UnderlineSpaces(true).
			UnderlineColor(style.Gray0).
			UnderlineStyle(lipgloss.UnderlineDashed)

	fileDirStyle  = fileLineStyle.Foreground(style.Gray0)
	filenameStyle = fileLineStyle.Bold(true)
	lineNumStyle  = style.Default.Foreground(style.Gray0)
	selectedBg    = style.BgDimRed
)

func Matches(
	files []string,
	byFile map[string][]search.ContentMatch,
	selected *MatchEntry,
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
			icon := matches[0].Icon
			dir, file := filepath.Split(path)
			iconStyle := fileLineStyle.Foreground(lipgloss.Color(icon.Color))
			myFileLineStyle := fileLineStyle
			myFileDirStyle := fileDirStyle
			myFilenameStyle := filenameStyle
			if selected != nil && selected.Match == nil && selected.Path == path {
				iconStyle = iconStyle.Background(selectedBg).UnderlineColor(style.Accent)
				myFileLineStyle = myFileLineStyle.Background(selectedBg).UnderlineColor(style.Accent)
				myFileDirStyle = myFileDirStyle.Background(selectedBg).UnderlineColor(style.Accent)
				myFilenameStyle = myFilenameStyle.Background(selectedBg).UnderlineColor(style.Accent)
			}

			writeLine(
				&b, width,
				myFileLineStyle,
				[]styledString{
					{Style: iconStyle, Text: icon.Icon + " "},
					{Style: myFileDirStyle, Text: dir},
					{Style: myFilenameStyle, Text: file},
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
				padStyle := style.Default
				myLineNumStyle := lineNumStyle
				myBg := style.Bg0

				isSelectedMatch := selected != nil && selected.Match != nil && selected.Path == path && selected.Match.LineNum == m.LineNum
				if isSelectedMatch {
					padStyle = padStyle.Background(selectedBg)
					myLineNumStyle = myLineNumStyle.Background(selectedBg).Foreground(style.Accent).Bold(true)
					myBg = selectedBg
				}

				writeLine(
					&b, width,
					padStyle,
					[]styledString{
						{Style: myLineNumStyle, Text: fmt.Sprintf("%5d ", m.LineNum)},
						{StyledText: cachedHighlight(m.Line, path, m.LineNum, myBg)},
					},
				)
			}
			currentLine++
		}
	}

	return b.String()
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
