package components

import (
	"charm.land/bubbles/v2/textinput"
	"github.com/michaelorr/ff/style"
)

func NewSearchInput() textinput.Model {
	input := textinput.New()
	input.Prompt = "❯ "

	// Styles.[Focused|Blurred|Cursor].[Text|Placeholder|Suggestion|Prompt]
	s := textinput.DefaultDarkStyles()
	s.Focused.Prompt = style.Default.Foreground(style.BgBlue)
	s.Blurred.Prompt = s.Focused.Prompt
	s.Focused.Text = style.Default
	s.Blurred.Text = style.Default
	s.Cursor.Color = style.Gray2
	input.SetStyles(s)
	input.Focus()

	return input
}
