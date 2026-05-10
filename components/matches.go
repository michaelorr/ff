// Package components provides helpers for rendering UI components
package components

import (
	"fmt"
	"path/filepath"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/michaelorr/ff/search"
	"github.com/michaelorr/ff/style"
)

var (
	highlightCache = map[string]map[int]string{}

	fileLineStyle = style.DefaultStyle.
			Underline(true).
			UnderlineSpaces(true).
			UnderlineColor(style.Gray0).
			UnderlineStyle(lipgloss.UnderlineDashed)

	fileDirStyle = fileLineStyle.
			Foreground(style.Gray0)

	filenameStyle = fileLineStyle.
			Bold(true)

	lineNumStyle = style.DefaultStyle.
			Foreground(style.Gray0)
)

func Matches(files []string, byFile map[string][]search.ContentMatch, width int) string {
	var b strings.Builder

	for _, path := range files {
		b.WriteByte('\n')
		icon := byFile[path][0].Icon
		iconStyle := fileLineStyle.Foreground(lipgloss.Color(icon.Color))

		dir, file := filepath.Split(path)
		line := iconStyle.Render(icon.Icon) + fileDirStyle.Render("", dir) + filenameStyle.Render(file)
		filler := fileLineStyle.Render(strings.Repeat(" ", max(0, width-lipgloss.Width(line))))
		fmt.Fprintf(&b, "%s%s\n", line, filler)

		for _, m := range byFile[path] {
			line := lineNumStyle.Render(fmt.Sprintf("%5d ", m.LineNum)) +
				cachedHighlight(m.Line, path, m.LineNum)
			filler := style.DefaultStyle.Render(strings.Repeat(" ", max(0, width-lipgloss.Width(line))))
			fmt.Fprint(&b, line, filler, "\n")
		}
	}

	return b.String()
}

func cachedHighlight(line, path string, lineNum int) string {
	if lines, ok := highlightCache[path]; ok {
		if cached, ok := lines[lineNum]; ok {
			return cached
		}
	} else {
		highlightCache[path] = map[int]string{}
	}
	result := syntaxHighlight(line, path)
	highlightCache[path][lineNum] = result
	return result
}
