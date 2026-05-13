package components

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/michaelorr/ff/style"
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

		previewLineHighlighted := isHighlighted(lineNum, targetLine, entry)
		writeLine(
			&b, width, previewDefault(previewLineHighlighted),
			[]styledString{
				{Style: gutterStyles(previewLineHighlighted), Text: gutterText(previewLineHighlighted)},
				{Style: previewLineNumStyles(previewLineHighlighted), Text: fmt.Sprintf("%*d ", lineNumWidth, lineNum)},
				{StyledText: cachedHighlight(lines[i], entry.Path, lineNum, syntaxBgColors(previewLineHighlighted))},
			},
		)
	}

	return b.String()
}

func isHighlighted(lineNum, targetLine int, entry *MatchEntry) bool {
	return entry.Match != nil && lineNum == targetLine
}

func previewDefault(isSelected bool) lipgloss.Style {
	myStyle := style.Default.Foreground(style.Gray0)
	if isSelected {
		myStyle = myStyle.Background(style.BgDimRed).Foreground(style.Accent)
	}
	return myStyle
}

func previewLineNumStyles(isSelected bool) lipgloss.Style {
	return previewDefault(isSelected).Bold(isSelected)
}

func gutterStyles(isSelected bool) lipgloss.Style {
	return previewDefault(isSelected)
}

func gutterText(isSelected bool) string {
	if isSelected {
		return "▸ "
	}
	return "  "
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
