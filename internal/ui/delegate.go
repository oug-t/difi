package ui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/oug-t/difi/internal/tree"
)

type TreeDelegate struct {
	Focused bool
}

func (d TreeDelegate) Height() int                               { return 1 }
func (d TreeDelegate) Spacing() int                              { return 0 }
func (d TreeDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

func (d TreeDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i, ok := item.(tree.TreeItem)
	if !ok {
		return
	}

	title := i.Title()

	// Truncate to list width to prevent line wrapping that inflates pane height.
	// Subtract 2 for safety margin: Nerd Font icons may render wider than
	// ansi.StringWidth reports (double-width glyphs).
	maxWidth := m.Width() - 2
	if maxWidth < 4 {
		maxWidth = 4
	}
	title = ansi.Truncate(title, maxWidth, "…")

	if index == m.Index() {
		// Width (not MaxWidth) ensures every row renders to a consistent width,
		// filling the selection highlight background fully and preventing
		// stale-character artifacts when content shrinks between frames.
		style := lipgloss.NewStyle().
			Background(lipgloss.Color("237")). // Dark gray background
			Foreground(lipgloss.Color("255")). // White text
			Bold(true).
			Width(maxWidth)

		if !d.Focused {
			style = style.Foreground(lipgloss.Color("245"))
		}

		fmt.Fprint(w, style.Render(title))
	} else {
		style := lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")).
			Width(maxWidth)
		fmt.Fprint(w, style.Render(title))
	}
}
