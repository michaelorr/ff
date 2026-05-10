package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	_ "github.com/ryboe/q"
)

func main() {
	p := tea.NewProgram(newModel())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Keep refactoring:
// - view.go
// - matches.go
// - style/
// - components/
