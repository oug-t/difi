package ui

import "github.com/charmbracelet/lipgloss"

var (
	ColorBorder   = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	ColorSelected = lipgloss.AdaptiveColor{Light: "#1F1F1F", Dark: "#F8F8F2"}
	ColorInactive = lipgloss.AdaptiveColor{Light: "#A8A8A8", Dark: "#626262"}
	ColorAccent   = lipgloss.AdaptiveColor{Light: "#00BCF0", Dark: "#00BCF0"} // Cyan

	PaneStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, true, false, false).
			BorderForeground(ColorBorder).
			MarginRight(1)

	DiffStyle = lipgloss.NewStyle().
			Padding(0, 1)

	ItemStyle = lipgloss.NewStyle().PaddingLeft(1)

	SelectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(0).
				Foreground(ColorSelected).
				Bold(true)
)
