package ui

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

	// 1. Setup Indentation
	indentSize := i.Depth * 2
	indent := strings.Repeat(" ", indentSize)

	// 2. Get Icon and Raw Name
	iconStr, iconStyle := getIconInfo(i.Path, i.IsDir)

	// 3. Truncation (Safety)
	availableWidth := m.Width() - indentSize - 4
	displayName := i.Path
	if availableWidth > 0 && len(displayName) > availableWidth {
		displayName = displayName[:max(0, availableWidth-1)] + "…"
	}

	// 4. Render Logic ("Oil" Block Cursor)
	var row string
	isSelected := index == m.Index()

	if isSelected && d.Focused {
		// -- SELECTED STATE (Oil Style) --
		// We do NOT use iconStyle here. We want the icon to inherit the
		// selection text color so the background block is unbroken.
		// Content: Icon + Space + Name
		content := fmt.Sprintf("%s %s", iconStr, displayName)

		// Apply the solid block style to the whole content
		renderedContent := SelectedBlockStyle.Render(content)

		// Combine: Indent (unhighlighted) + Block (highlighted)
		row = fmt.Sprintf("%s%s", indent, renderedContent)

	} else {
		// -- NORMAL / INACTIVE STATE --
		// Render icon with its specific color
		renderedIcon := iconStyle.Render(iconStr)

		// Combine
		row = fmt.Sprintf("%s%s %s", indent, renderedIcon, displayName)

		// Apply generic padding/style
		row = ItemStyle.Render(row)
	}

	fmt.Fprint(w, row)
}

// Helper: Returns raw icon string and its preferred style
func getIconInfo(name string, isDir bool) (string, lipgloss.Style) {
	if isDir {
		return "", FolderIconStyle
	}

	ext := filepath.Ext(name)
	icon := ""

	switch strings.ToLower(ext) {
	case ".go":
		icon = ""
	case ".js", ".ts", ".tsx", ".jsx":
		icon = ""
	case ".md":
		icon = ""
	case ".json", ".yml", ".yaml", ".toml":
		icon = ""
	case ".css", ".scss":
		icon = ""
	case ".html":
		icon = ""
	case ".git":
		icon = ""
	case ".dockerfile":
		icon = ""
	case ".svelte":
		icon = ""
	}

	return icon, FileIconStyle
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
