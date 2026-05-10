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
	var b strings.Builder
	b.WriteByte('\n')
	types := slices.Collect(maps.Keys(icons))
	slices.SortFunc(types, func(a, b search.FileIcon) int {
		return strings.Compare(a.Icon, b.Icon)
	})

	for _, icon := range types {
		iconStyle := style.DefaultStyle.Foreground(lipgloss.Color(icon.Color))

		line := iconStyle.Render(icon.Icon) + style.DefaultStyle.Render(fmt.Sprintf(" %d", icons[icon]))
		line = line + style.DefaultStyle.Render(strings.Repeat(" ", max(0, width-lipgloss.Width(line))), "\n")
		b.WriteString(line)
	}
	return b.String()
}
