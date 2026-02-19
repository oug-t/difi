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

	maxWidth := m.Width() - 2 // pane has Padding(0,1) = 2 chars horizontal
	if maxWidth < 1 {
		maxWidth = 1
	}
	title = ansi.Truncate(title, maxWidth, "…")

	if index == m.Index() {
		style := lipgloss.NewStyle().
			Width(maxWidth).
			Background(lipgloss.Color("237")). // Dark gray background
			Foreground(lipgloss.Color("255")). // White text
			Bold(true)

		if !d.Focused {
			style = style.Foreground(lipgloss.Color("245"))
		}

		fmt.Fprint(w, style.Render(title))
	} else {
		style := lipgloss.NewStyle().Width(maxWidth).Foreground(lipgloss.Color("252"))
		fmt.Fprint(w, style.Render(title))
	}
}
