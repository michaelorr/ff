package components

import (
	"image/color"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"

	"github.com/michaelorr/ff/style"
)

var (
	gruvbox = styles.Get("gruvbox")
	base    = style.Default
)

// syntaxHighlightWithBg returns line with ANSI syntax-highlighting applied based on
// the language inferred from path. Returns line unchanged if no lexer matches.
func syntaxHighlightWithBg(line, path string, bg color.Color) string {
	b := base.Background(bg)
	lexer := lexers.Match(path)
	if lexer == nil {
		return b.Render(line)
	}

	iter, err := chroma.Coalesce(lexer).Tokenise(nil, line)
	if err != nil {
		return b.Render(line)
	}

	var sb strings.Builder
	for tok := iter(); tok != chroma.EOF; tok = iter() {
		value := strings.TrimRight(tok.Value, "\n")
		if value == "" {
			continue
		}
		entry := gruvbox.Get(tok.Type)
		s := b
		if entry.Colour.IsSet() {
			s = s.Foreground(lipgloss.Color(entry.Colour.String()))
		}
		if entry.Bold == chroma.Yes {
			s = s.Bold(true)
		}
		if entry.Italic == chroma.Yes {
			s = s.Italic(true)
		}
		sb.WriteString(s.Render(value))
	}

	return sb.String()
}
