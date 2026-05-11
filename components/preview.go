package components

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/michaelorr/ff/style"
)

var (
	previewLineNumStyle     = style.Default.Foreground(style.Gray0)
	previewSelectedNumStyle = style.Default.Foreground(style.Accent).Bold(true)

	previewGutterStyle         = style.Default.Foreground(style.Gray0)
	previewSelectedGutterStyle = style.Default.Foreground(style.Accent)
)

func Preview(entry *MatchEntry, width, height int) string {
	if entry == nil || height <= 0 {
		return ""
	}

	data, err := os.ReadFile(entry.Path)
	if err != nil {
		return style.Default.Foreground(style.Red).Render("error: " + err.Error())
	}

	lines := strings.Split(strings.TrimRight(string(data), "\n"), "\n")
	numLines := len(lines)
	lineNumWidth := len(strconv.Itoa(numLines))

	var targetLine int
	if entry.Match != nil {
		targetLine = entry.Match.LineNum
	}

	start, end := visibleWindow(numLines, targetLine, height)

	var b strings.Builder
	for i := start; i < end; i++ {
		lineNum := i + 1
		isSelected := entry.Match != nil && lineNum == targetLine

		var partial string
		var lineLen int

		if isSelected {
			partial = previewSelectedGutterStyle.Render("▸ ") + previewSelectedNumStyle.Render(fmt.Sprintf("%*d ", lineNumWidth, lineNum))
			lineLen += lipgloss.Width(partial)
			b.WriteString(partial)
		} else {
			partial = previewGutterStyle.Render("  ") + previewLineNumStyle.Render(fmt.Sprintf("%*d ", lineNumWidth, lineNum))
			lineLen += lipgloss.Width(partial)
			b.WriteString(partial)
		}

		partial = cachedHighlight(lines[i], entry.Path, lineNum, style.Bg0)
		lineLen += lipgloss.Width(partial)
		b.WriteString(partial)
		b.WriteString(style.Default.Render(strings.Repeat(" ", max(0, width-lineLen))))

		b.WriteByte('\n')
	}

	return b.String()
}

// visibleWindow computes the [start, end) line indices (0-based) to display,
// centering targetLine (1-based) in the viewport. When targetLine is 0 (file
// selection), the window starts at the top of the file.
func visibleWindow(total, targetLine, height int) (start, end int) {
	if targetLine == 0 {
		return 0, min(total, height)
	}
	half := height / 2
	start = max(0, targetLine-1-half)
	end = min(total, start+height)
	if end-start < height {
		start = max(0, end-height)
	}
	return start, end
}
