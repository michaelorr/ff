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

func Matches(files []string, byFile map[string][]search.ContentMatch, selected *MatchEntry, width int) string {
	var b strings.Builder
	for _, path := range files {
		var partial string
		var lineLen int
		icon := byFile[path][0].Icon
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

		b.WriteByte('\n')
		partial = iconStyle.Render(icon.Icon + " ")
		lineLen += lipgloss.Width(partial)
		b.WriteString(partial)

		partial = myFileDirStyle.Render(dir)
		lineLen += lipgloss.Width(partial)
		b.WriteString(partial)

		partial = myFilenameStyle.Render(file)
		lineLen += lipgloss.Width(partial)
		b.WriteString(partial)

		b.WriteString(myFileLineStyle.Render(strings.Repeat(" ", max(0, width-lineLen))))
		b.WriteByte('\n')

		lineLen = 0
		for _, m := range byFile[path] {
			padStyle := style.Default
			myLineNumStyle := lineNumStyle

			isSelectedMatch := selected != nil && selected.Match != nil && selected.Path == path && selected.Match.LineNum == m.LineNum
			if isSelectedMatch {
				padStyle = padStyle.Background(selectedBg)
				myLineNumStyle = myLineNumStyle.Background(selectedBg).Foreground(style.Accent).Bold(true)
			}

			partial = myLineNumStyle.Render(fmt.Sprintf("%5d ", m.LineNum))
			lineLen += lipgloss.Width(partial)
			b.WriteString(partial)

			if isSelectedMatch {
				partial = cachedHighlight(m.Line, path, m.LineNum, selectedBg)
			} else {
				partial = cachedHighlight(m.Line, path, m.LineNum, style.Bg0)
			}
			lineLen += lipgloss.Width(partial)
			b.WriteString(partial)

			b.WriteString(padStyle.Render(strings.Repeat(" ", max(0, width-lineLen))))
			b.WriteByte('\n')

			lineLen = 0
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
