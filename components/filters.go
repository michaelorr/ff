package components

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/michaelorr/ff/search"
	"github.com/michaelorr/ff/style"
)

func RenderFilters(icons map[search.FileIcon]int) string {
	var b strings.Builder
	b.WriteByte('\n')
	types := slices.Collect(maps.Keys(icons))
	slices.SortFunc(types, func(a, b search.FileIcon) int {
		return strings.Compare(a.Icon, b.Icon)
	})

	defaultStyle := lipgloss.NewStyle().Foreground(style.Fg0).Background(style.Bg0)
	for _, icon := range types {
		iconStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(icon.Color)).Background(style.Bg0)

		fmt.Fprintf(
			&b, "%s%s%s",
			iconStyle.Render(icon.Icon),
			defaultStyle.Render(fmt.Sprintf(" %d", icons[icon])),
			"\n",
		)
	}
	return b.String()
}
