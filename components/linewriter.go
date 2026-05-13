package components

import (
	"strings"

	"charm.land/lipgloss/v2"
)

// Provide either a raw string and a style, or a pre-styled string.
type styledString struct {
	Style      lipgloss.Style
	Text       string
	StyledText string
}

func writeLine(b *strings.Builder, length int, fillStyle lipgloss.Style, inputs []styledString) {
	soFar := 0
	for _, input := range inputs {
		partial := input.StyledText
		if partial == "" {
			partial = input.Style.Render(input.Text)
		}
		b.WriteString(partial)
		soFar += lipgloss.Width(partial)
	}
	b.WriteString(fillStyle.Render(strings.Repeat(" ", max(0, length-soFar))))
	b.WriteByte('\n')
}
