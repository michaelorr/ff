// Package panel contains rendering helpers for panel components
//
//	╭─ title ───────────╮
//	│  body...          │
//	╰───────────────────╯
package panel

import (
	"strings"

	"charm.land/bubbles/v2/textinput"
	"charm.land/lipgloss/v2"
	"github.com/michaelorr/ff/colors"
)

type Viewable interface {
	View() string
}

var (
	titleStyle         = lipgloss.NewStyle().Foreground(colors.Accent).Bold(true).Background(colors.Bg0)
	blurredBorderColor = colors.Gray0
	focusedBorderColor = colors.Accent
)

func Render(title string, width, height int, body Viewable, focused bool) string {
	borderColor := blurredBorderColor
	if focused {
		borderColor = focusedBorderColor
	}

	topBorderStyle := lipgloss.NewStyle().Foreground(borderColor).Background(colors.Bg0)
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
		BorderBackground(colors.Bg0).
		Background(colors.Bg0).
		BorderStyle(border).
		BorderTop(false).
		BorderLeft(true).
		BorderRight(true).
		BorderBottom(true).
		Width(innerW+2).
		Height(height-1).
		Padding(0, 1).
		Render(body.View())

	return lipgloss.JoinVertical(lipgloss.Left, topRow, box)
}

func NewSearchPanel(textinput textinput.Model) string {
	return "foo"
}
