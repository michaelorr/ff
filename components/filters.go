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

func Filters(icons map[search.FileIcon]int, width int) string {
	types := slices.Collect(maps.Keys(icons))
	slices.SortFunc(types, func(a, b search.FileIcon) int {
		if icons[a] == icons[b] {
			return strings.Compare(a.Icon, b.Icon)
		}
		return icons[b] - icons[a]
	})

	var b strings.Builder
	b.WriteByte('\n')
	for _, icon := range types {
		iconStyle := style.Default.Foreground(lipgloss.Color(icon.Color))
		writeLine(
			&b, width,
			style.Default,
			[]styledString{
				{Style: iconStyle, Text: icon.Icon + " "},
				{Style: style.Default, Text: fmt.Sprintf("%d", icons[icon])},
			},
		)
	}
	return b.String()
}
