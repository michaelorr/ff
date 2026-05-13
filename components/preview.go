package components

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/michaelorr/ff/style"
)

var (
	previewLineNumStyle     = style.Default.Foreground(style.Gray0)
	previewSelectedNumStyle = style.Default.Foreground(style.Accent).Bold(true).Background(style.BgDimRed)

	previewGutterStyle         = style.Default.Foreground(style.Gray0)
	previewSelectedGutterStyle = style.Default.Foreground(style.Accent)

	previewSelectedLineStyle = style.Default.Background(style.BgDimRed)
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

		myLineStyle := style.Default
		myBg := style.Bg0
		myGutterStyle := previewGutterStyle
		gutterText := "  "
		myLineNumStyle := previewLineNumStyle
		if isSelected {
			myLineStyle = previewSelectedLineStyle
			myBg = style.BgDimRed
			myGutterStyle = previewSelectedGutterStyle
			gutterText = "▸ "
			myLineNumStyle = previewSelectedNumStyle
		}

		writeLine(
			&b, width, myLineStyle,
			[]styledString{
				{Style: myGutterStyle, Text: gutterText},
				{Style: myLineNumStyle, Text: fmt.Sprintf("%*d ", lineNumWidth, lineNum)},
				{StyledText: cachedHighlight(lines[i], entry.Path, lineNum, myBg)},
			},
		)
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
