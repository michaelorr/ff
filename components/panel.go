package components

import (
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/michaelorr/ff/style"
)

type Viewable interface {
	View() string
}

var (
	titleStyle         = style.Default.Foreground(style.Accent).Bold(true)
	blurredBorderColor = style.Gray0
	focusedBorderColor = style.Accent

	border = lipgloss.RoundedBorder()
)

type Size struct {
	Width, Height int
}

// Panel contains rendering helpers for panel components
//
//	╭─ title ───────────╮
//	│  body...          │
//	╰───────────────────╯
func Panel(title string, s Size, body Viewable, focused bool) string {
	borderColor := blurredBorderColor
	topBorderStyle := style.Default.Foreground(blurredBorderColor)
	if focused {
		borderColor = focusedBorderColor
		topBorderStyle = topBorderStyle.Foreground(focusedBorderColor)
	}

	title = titleStyle.Render(" " + title + " ")
	innerW := max(s.Width-2, 0)
	titleW := lipgloss.Width(title)
	remaining := max(innerW-titleW-1, 0)

	topRow := "" +
		topBorderStyle.Render(border.TopLeft+border.Top) +
		title +
		topBorderStyle.Render(strings.Repeat(border.Top, remaining)+border.TopRight)

	box := style.Default.
		BorderForeground(borderColor).
		BorderBackground(style.Bg0).
		BorderStyle(border).
		BorderTop(false).
		BorderLeft(true).
		BorderRight(true).
		BorderBottom(true).
		Width(max(s.Width, 2)).
		Height(s.Height-1).
		Padding(0, 1).
		Render(body.View())

	return lipgloss.JoinVertical(lipgloss.Left, topRow, box)
}
