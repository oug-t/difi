package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/oug-t/difi/internal/config"
)

// Global UI colors and styles.
// Initialized once at startup via InitStyles.
var (
	ColorBorder   lipgloss.AdaptiveColor
	ColorFocus    lipgloss.AdaptiveColor
	ColorText     = lipgloss.AdaptiveColor{Light: "#1F1F1F", Dark: "#F8F8F2"}
	ColorSubtle   = lipgloss.AdaptiveColor{Light: "#A8A8A8", Dark: "#626262"}
	ColorCursorBg = lipgloss.AdaptiveColor{Light: "#E5E5E5", Dark: "#3E3E3E"}

	ColorBarBg = lipgloss.AdaptiveColor{Light: "#F2F2F2", Dark: "#1F1F1F"}
	ColorBarFg = lipgloss.AdaptiveColor{Light: "#6E6E6E", Dark: "#9E9E9E"}

	PaneStyle          lipgloss.Style
	FocusedPaneStyle   lipgloss.Style
	DiffStyle          lipgloss.Style
	ItemStyle          lipgloss.Style
	SelectedItemStyle  lipgloss.Style
	LineNumberStyle    lipgloss.Style
	StatusBarStyle     lipgloss.Style
	StatusKeyStyle     lipgloss.Style
	StatusDividerStyle lipgloss.Style
	HelpTextStyle      lipgloss.Style
	HelpDrawerStyle    lipgloss.Style

	CurrentConfig config.Config
)

// InitStyles initializes global styles based on the provided config.
// This should be called once during application startup.
func InitStyles(cfg config.Config) {
	CurrentConfig = cfg

	// Colors derived from user config
	ColorBorder = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: cfg.Colors.Border}
	ColorFocus = lipgloss.AdaptiveColor{Light: "#000000", Dark: cfg.Colors.Focus}

	// Pane styles
	PaneStyle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, cfg.UI.ShowGuide, false, false).
		BorderForeground(ColorBorder)

	FocusedPaneStyle = PaneStyle.Copy().
		BorderForeground(ColorFocus)

	// Diff and list item styles
	DiffStyle = lipgloss.NewStyle().Padding(0, 0)
	ItemStyle = lipgloss.NewStyle().PaddingLeft(2)

	SelectedItemStyle = lipgloss.NewStyle().
		PaddingLeft(1).
		Background(ColorCursorBg).
		Foreground(ColorText).
		Bold(true).
		Width(1000)

	// Line numbers
	LineNumberStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(cfg.Colors.LineNumber)).
		PaddingRight(1).
		Width(4)

	// Status bar
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
		Background(ColorBarBg)

	// Help drawer
	HelpTextStyle = lipgloss.NewStyle().
		Foreground(ColorSubtle).
		Padding(0, 1)

	HelpDrawerStyle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), true, false, false, false).
		BorderForeground(ColorBorder).
		Padding(1, 2)
}
