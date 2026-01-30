package ui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
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

	// If this item is selected
	if index == m.Index() {
		if d.Focused {
			// Render the whole line (including indent) with the selection background
			fmt.Fprint(w, SelectedBlockStyle.Render(title))
		} else {
			// Dimmed selection if focus is on the other panel
			fmt.Fprint(w, SelectedBlockStyle.Copy().Foreground(ColorSubtle).Render(title))
		}
	} else {
		// Normal Item (No icons added, just the text)
		fmt.Fprint(w, ItemStyle.Render(title))
	}
}
