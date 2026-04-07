// Package theme defines the Ayu Dark color scheme and reusable styles.
package theme

import "github.com/charmbracelet/lipgloss"

// Ayu Dark color palette.
var (
	ColorBg      = lipgloss.Color("#0D1017")
	ColorFg      = lipgloss.Color("#BFBDB6")
	ColorAccent  = lipgloss.Color("#E6B450")
	ColorOrange  = lipgloss.Color("#FF8F40")
	ColorYellow  = lipgloss.Color("#FFB454")
	ColorGreen   = lipgloss.Color("#7FD962")
	ColorBlue    = lipgloss.Color("#59C2FF")
	ColorPurple  = lipgloss.Color("#D2A6FF")
	ColorRed     = lipgloss.Color("#D95757")
	ColorCyan    = lipgloss.Color("#95E6CB")
	ColorComment = lipgloss.Color("#636A72")
	ColorLine    = lipgloss.Color("#131721")
)

// Reusable styles for consistent output across the CLI.
var (
	Title     = lipgloss.NewStyle().Foreground(ColorAccent).Bold(true)
	Bold      = lipgloss.NewStyle().Foreground(ColorFg).Bold(true)
	Subtle    = lipgloss.NewStyle().Foreground(ColorComment)
	Success   = lipgloss.NewStyle().Foreground(ColorGreen)
	Error     = lipgloss.NewStyle().Foreground(ColorRed).Bold(true)
	Warning   = lipgloss.NewStyle().Foreground(ColorOrange)
	Info      = lipgloss.NewStyle().Foreground(ColorBlue)
	Highlight = lipgloss.NewStyle().Foreground(ColorPurple)
	Accent    = lipgloss.NewStyle().Foreground(ColorAccent)
	Muted     = lipgloss.NewStyle().Foreground(ColorComment)
)
