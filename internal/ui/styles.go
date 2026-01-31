package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/oug-t/difi/internal/config"
)

var (
	// Config
	CurrentConfig config.Config

	// Theme colors
	ColorBorder   = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	ColorFocus    = lipgloss.AdaptiveColor{Light: "#000000", Dark: "#E5E5E5"}
	ColorText     = lipgloss.AdaptiveColor{Light: "#1F1F1F", Dark: "#F8F8F2"}
	ColorSubtle   = lipgloss.AdaptiveColor{Light: "#A8A8A8", Dark: "#626262"}
	ColorCursorBg = lipgloss.AdaptiveColor{Light: "#E5E5E5", Dark: "#3E3E3E"}
	ColorAccent   = lipgloss.AdaptiveColor{Light: "#00ADD8", Dark: "#00ADD8"} // Go blue

	// Pane styles
	PaneStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, true, false, false).
			BorderForeground(ColorBorder)

	FocusedPaneStyle = PaneStyle.Copy().
				BorderForeground(ColorFocus)

	DiffStyle = lipgloss.NewStyle().Padding(0, 0)
	ItemStyle = lipgloss.NewStyle().PaddingLeft(2)

	// List styles
	SelectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(1).
				Background(ColorCursorBg).
				Foreground(ColorText).
				Bold(true).
				Width(1000)

	SelectedBlockStyle = lipgloss.NewStyle().
				Background(ColorCursorBg).
				Foreground(ColorText).
				Bold(true).
				PaddingLeft(1)

	// Icon styles
	FolderIconStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#F7B96E", Dark: "#E5C07B"})
	FileIconStyle   = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#969696", Dark: "#ABB2BF"})

	// Diff view styles
	LineNumberStyle = lipgloss.NewStyle().
			Foreground(ColorSubtle).
			PaddingRight(1).
			Width(4)

	DiffSelectionStyle = lipgloss.NewStyle().
				Background(ColorCursorBg).
				Width(1000)

	// Status bar colors
	ColorBarBg = lipgloss.AdaptiveColor{Light: "#F2F2F2", Dark: "#1F1F1F"}
	ColorBarFg = lipgloss.AdaptiveColor{Light: "#6E6E6E", Dark: "#9E9E9E"}

	// Status bar styles
	StatusBarStyle = lipgloss.NewStyle().
			Foreground(ColorBarFg).
			Background(ColorBarBg).
			Padding(0, 1)

	StatusKeyStyle = lipgloss.NewStyle().
			Foreground(ColorText).
			Background(ColorBarBg).
			Bold(true).
			Padding(0, 1)

	StatusDividerStyle = lipgloss.NewStyle().
				Foreground(ColorSubtle).
				Background(ColorBarBg).
				Padding(0, 0)

	// Help styles
	HelpTextStyle = lipgloss.NewStyle().
			Foreground(ColorSubtle).
			Padding(0, 1)

	HelpDrawerStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), true, false, false, false).
			BorderForeground(ColorBorder).
			Padding(1, 2)

	// Empty/landing styles
	EmptyLogoStyle = lipgloss.NewStyle().
			Foreground(ColorAccent).
			Bold(true).
			PaddingBottom(1)

	EmptyDescStyle = lipgloss.NewStyle().
			Foreground(ColorSubtle).
			PaddingBottom(2)

	EmptyStatusStyle = lipgloss.NewStyle().
				Foreground(ColorText).
				Background(ColorCursorBg).
				Padding(0, 2).
				MarginBottom(2)

	EmptyCodeStyle = lipgloss.NewStyle().
			Foreground(ColorSubtle).
			MarginLeft(2)

	EmptyHeaderStyle = lipgloss.NewStyle().
				Foreground(ColorText).
				Bold(true).
				MarginBottom(1)
)

func InitStyles(cfg config.Config) {
	CurrentConfig = cfg
}
