package components

import (
	"charm.land/bubbles/v2/textinput"
	"charm.land/lipgloss/v2"
	"github.com/michaelorr/ff/style"
)

func NewSearchInput() textinput.Model {
	input := textinput.New()
	input.Prompt = "❯ "

	// Styles.[Focused|Blurred|Cursor].[Text|Placeholder|Suggestion|Prompt]
	s := textinput.DefaultDarkStyles()
	s.Focused.Prompt = lipgloss.NewStyle().Foreground(style.BgBlue)
	s.Blurred.Prompt = s.Focused.Prompt
	s.Focused.Text = lipgloss.NewStyle().Foreground(style.Fg0).Background(style.Bg0)
	s.Blurred.Text = s.Focused.Text
	s.Cursor.Color = style.Gray2
	input.SetStyles(s)
	input.Focus()

	return input
}
