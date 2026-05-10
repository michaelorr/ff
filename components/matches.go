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
	defaultStyle = lipgloss.NewStyle().
			Foreground(style.Fg0).
			Background(style.Bg0)

	fileLineStyle = defaultStyle.
			Background(style.Bg1)

	fileDirStyle = fileLineStyle.
			Foreground(style.Gray0)

	filenameStyle = fileLineStyle.
			Bold(true).
			Underline(true)

	lineNumStyle = defaultStyle.
			Foreground(style.Gray0)
)

func RenderMatches(files []string, byFile map[string][]search.ContentMatch, width int) string {
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
				// SyntaxHighlight(m.Line, path)
				m.Line
			filler := defaultStyle.Render(strings.Repeat(" ", max(0, width-lipgloss.Width(line))))
			fmt.Fprint(&b, line, filler, "\n")
		}
	}

	return b.String()
}
