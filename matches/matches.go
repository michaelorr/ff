// Package matches provides functions for rendering search matches.
package matches

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/michaelorr/ff/internal/search"
)

func RenderMatches(files []string, byFile map[string][]search.ContentMatch) string {
	var b strings.Builder
	for _, path := range files {
		b.WriteByte('\n')
		ext := strings.TrimPrefix(filepath.Ext(path), ".")
		if ext == "" {
			ext = "  "
		}
		fmt.Fprintf(&b, "[%s] %s\n", ext, path)
		for _, m := range byFile[path] {
			fmt.Fprintf(&b, "%5d %s\n", m.LineNum, m.Line)
		}
	}
	return b.String()
}
