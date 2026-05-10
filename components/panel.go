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
	titleStyle         = lipgloss.NewStyle().Foreground(style.Accent).Bold(true).Background(style.Bg0)
	blurredBorderColor = style.Gray0
	focusedBorderColor = style.Accent
)

// RenderPanel contains rendering helpers for panel components
//
//	╭─ title ───────────╮
//	│  body...          │
//	╰───────────────────╯
func RenderPanel(title string, width, height int, body Viewable, focused bool) string {
	borderColor := blurredBorderColor
	if focused {
		borderColor = focusedBorderColor
	}

	topBorderStyle := lipgloss.NewStyle().Foreground(borderColor).Background(style.Bg0)
	title = titleStyle.Render(" " + title + " ")
	border := lipgloss.RoundedBorder()
	innerW := max(width-2, 0)
	titleW := lipgloss.Width(title)
	remaining := max(innerW-titleW-1, 0)

	topRow := "" +
		topBorderStyle.Render(border.TopLeft+border.Top) +
		title +
		topBorderStyle.Render(strings.Repeat(border.Top, remaining)+border.TopRight)

	box := lipgloss.NewStyle().
		BorderForeground(borderColor).
		BorderBackground(style.Bg0).
		Background(style.Bg0).
		BorderStyle(border).
		BorderTop(false).
		BorderLeft(true).
		BorderRight(true).
		BorderBottom(true).
		Width(max(width, 2)).
		Height(height-1).
		Padding(0, 1).
		Render(body.View())

	return lipgloss.JoinVertical(lipgloss.Left, topRow, box)
}
