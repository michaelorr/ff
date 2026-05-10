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

	fileLineStyle = style.Default.
			Underline(true).
			UnderlineSpaces(true).
			UnderlineColor(style.Gray0).
			UnderlineStyle(lipgloss.UnderlineDashed)

	fileDirStyle = fileLineStyle.
			Foreground(style.Gray0)

	filenameStyle = fileLineStyle.
			Bold(true)

	lineNumStyle = style.Default.
			Foreground(style.Gray0)
)

func Matches(files []string, byFile map[string][]search.ContentMatch, width int) string {
	var b strings.Builder
	for _, path := range files {
		var partial string
		var lineLen int
		icon := byFile[path][0].Icon
		iconStyle := fileLineStyle.Foreground(lipgloss.Color(icon.Color))
		dir, file := filepath.Split(path)

		b.WriteByte('\n')
		partial = iconStyle.Render(icon.Icon + " ")
		lineLen += lipgloss.Width(partial)
		b.WriteString(partial)

		partial = fileDirStyle.Render(dir)
		lineLen += lipgloss.Width(partial)
		b.WriteString(partial)

		partial = filenameStyle.Render(file)
		lineLen += lipgloss.Width(partial)
		b.WriteString(partial)

		b.WriteString(fileLineStyle.Render(strings.Repeat(" ", max(0, width-lineLen))))
		b.WriteByte('\n')

		lineLen = 0
		for _, m := range byFile[path] {
			partial = lineNumStyle.Render(fmt.Sprintf("%5d ", m.LineNum))
			lineLen += lipgloss.Width(partial)
			b.WriteString(partial)

			partial = cachedHighlight(m.Line, path, m.LineNum)
			lineLen += lipgloss.Width(partial)
			b.WriteString(partial)

			b.WriteString(style.Default.Render(strings.Repeat(" ", max(0, width-lineLen))))
			b.WriteByte('\n')

			lineLen = 0
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
