// Package style
package style

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

var (
	Bg0 = lipgloss.Color("#1d2021")
	Bg1 = lipgloss.Color("#282828")
	Bg3 = lipgloss.Color("#3c3836")
	Bg5 = lipgloss.Color("#504945")

	BgDim      = lipgloss.Color("#141617")
	BgDimRed   = lipgloss.Color("#3c1f1e")
	BgDimGreen = lipgloss.Color("#32361a")
	BgDimBlue  = lipgloss.Color("#0d3138")

	Fg0 = lipgloss.Color("#ebdbb2")

	Red    = lipgloss.Color("#fb4934")
	Green  = lipgloss.Color("#b8bb26")
	Blue   = lipgloss.Color("#83a598")
	Yellow = lipgloss.Color("#fabd2f")
	Purple = lipgloss.Color("#d3869b")
	Orange = lipgloss.Color("#fe8019")
	Aqua   = lipgloss.Color("#8ec07c")

	BgRed    = lipgloss.Color("#cc241d")
	BgGreen  = lipgloss.Color("#98971a")
	BgBlue   = lipgloss.Color("#458588")
	BgYellow = lipgloss.Color("#d79921")
	BgPurple = lipgloss.Color("#b16286")
	BgOrange = lipgloss.Color("#fe8019")
	BgAqua   = lipgloss.Color("#689d6a")

	Gray0 = lipgloss.Color("#7c6f64")
	Gray1 = lipgloss.Color("#928374")
	Gray2 = lipgloss.Color("#a89984")

	Accent = BgYellow

	ChicletColors = []color.Color{Red, Green, Blue, Yellow, Purple, Orange, Aqua}
)
