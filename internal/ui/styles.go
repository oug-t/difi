package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/oug-t/difi/internal/config"
)

var (
	// -- Colors --
	ColorText   = lipgloss.AdaptiveColor{Light: "#24292f", Dark: "#c9d1d9"}
	ColorSubtle = lipgloss.AdaptiveColor{Light: "#6e7781", Dark: "#8b949e"}

	// UNIFIED SELECTION COLOR (The "Neutral Light Transparent Blue")
	// This is used for BOTH the file tree and the diff panel background.
	// Dark: Deep subtle slate blue | Light: Pale selection blue
	ColorVisualBg = lipgloss.AdaptiveColor{Light: "#daeaff", Dark: "#3a4b5c"}

	// Tree Text Color (High Contrast for the block cursor)
	ColorVisualFg = lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"}

	ColorFolder = lipgloss.AdaptiveColor{Light: "#0969da", Dark: "#83a598"}
	ColorFile   = lipgloss.AdaptiveColor{Light: "#24292f", Dark: "#ebdbb2"}

	ColorBarBg = lipgloss.AdaptiveColor{Light: "#F2F2F2", Dark: "#1F1F1F"}
	ColorBarFg = lipgloss.AdaptiveColor{Light: "#6E6E6E", Dark: "#9E9E9E"}

	// -- Styles --
	PaneStyle        lipgloss.Style
	FocusedPaneStyle lipgloss.Style
	DiffStyle        lipgloss.Style

	ItemStyle          lipgloss.Style
	SelectedBlockStyle lipgloss.Style // Tree (Opaque)
	DiffSelectionStyle lipgloss.Style // Diff (Transparent/BG only)

	FolderIconStyle lipgloss.Style
	FileIconStyle   lipgloss.Style
	LineNumberStyle lipgloss.Style

	StatusBarStyle     lipgloss.Style
	StatusKeyStyle     lipgloss.Style
	StatusDividerStyle lipgloss.Style
	HelpTextStyle      lipgloss.Style
	HelpDrawerStyle    lipgloss.Style

	CurrentConfig config.Config
)

func InitStyles(cfg config.Config) {
	CurrentConfig = cfg

	ColorBorder := lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: cfg.Colors.Border}
	ColorFocus := lipgloss.AdaptiveColor{Light: "#6e7781", Dark: cfg.Colors.Focus}

	// Allow user override for the selection background
	var selectionBg lipgloss.TerminalColor
	if cfg.Colors.DiffSelectionBg != "" {
		selectionBg = lipgloss.Color(cfg.Colors.DiffSelectionBg)
	} else {
		selectionBg = ColorVisualBg
	}

	PaneStyle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, cfg.UI.ShowGuide, false, false).
		BorderForeground(ColorBorder)

	FocusedPaneStyle = PaneStyle.Copy().
		BorderForeground(ColorFocus)

	DiffStyle = lipgloss.NewStyle().Padding(0, 0)

	// Base Row
	ItemStyle = lipgloss.NewStyle().
		PaddingLeft(1).
		PaddingRight(1).
		Foreground(ColorText)

	// 1. LEFT PANE STYLE (Tree)
	// Uses the shared background + forces a foreground color for readability
	SelectedBlockStyle = lipgloss.NewStyle().
		Background(selectionBg).
		Foreground(ColorVisualFg).
		PaddingLeft(1).
		PaddingRight(1).
		Bold(true)

	// 2. RIGHT PANE STYLE (Diff)
	// Uses the SAME shared background, but NO foreground.
	// This makes it "transparent" so Green(+)/Red(-) text colors show through.
	DiffSelectionStyle = lipgloss.NewStyle().
		Background(selectionBg)

	FolderIconStyle = lipgloss.NewStyle().Foreground(ColorFolder)
	FileIconStyle = lipgloss.NewStyle().Foreground(ColorFile)

	LineNumberStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(cfg.Colors.LineNumber)).
		PaddingRight(1).
		Width(4)

	StatusBarStyle = lipgloss.NewStyle().Foreground(ColorBarFg).Background(ColorBarBg).Padding(0, 1)
	StatusKeyStyle = lipgloss.NewStyle().Foreground(ColorText).Background(ColorBarBg).Bold(true).Padding(0, 1)
	StatusDividerStyle = lipgloss.NewStyle().Foreground(ColorSubtle).Background(ColorBarBg).Padding(0, 0)

	HelpTextStyle = lipgloss.NewStyle().Foreground(ColorSubtle).Padding(0, 1)
	HelpDrawerStyle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), true, false, false, false).
		BorderForeground(ColorBorder).
		Padding(1, 2)
}
