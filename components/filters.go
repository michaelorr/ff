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
		var partial string
		var lineLen int
		iconStyle := style.Default.Foreground(lipgloss.Color(icon.Color))

		partial = iconStyle.Render(icon.Icon + " ")
		lineLen += lipgloss.Width(partial)
		b.WriteString(partial)

		partial = style.Default.Render(fmt.Sprintf("%d", icons[icon]))
		lineLen += lipgloss.Width(partial)
		b.WriteString(partial)
		b.WriteString(style.Default.Render(strings.Repeat(" ", max(0, width-lineLen))))
		b.WriteByte('\n')
	}
	return b.String()
}
