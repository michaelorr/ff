// Package components provides helpers for rendering UI components
package components

import (
	"fmt"
	"path/filepath"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/michaelorr/ff/colors"
	"github.com/michaelorr/ff/search"
)

var (
	defaultStyle = lipgloss.NewStyle().
			Foreground(colors.Fg0).
			Background(colors.Bg0)

	fileLineStyle = defaultStyle.
			Background(colors.Bg1)

	fileDirStyle = fileLineStyle.
			Foreground(colors.Gray0)

	filenameStyle = fileLineStyle.
			Bold(true).
			Underline(true)

	lineNumStyle = defaultStyle.
			Foreground(colors.Gray0)
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
				SyntaxHighlight(m.Line, path)
			filler := defaultStyle.Render(strings.Repeat(" ", max(0, width-lipgloss.Width(line))))
			fmt.Fprint(&b, line, filler, "\n")
		}
	}

	return b.String()
}
