package ui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/oug-t/difi/internal/git"
	"github.com/oug-t/difi/internal/tree"
)

const TargetBranch = "main"

type Model struct {
	fileTree      list.Model
	diffViewport  viewport.Model
	selectedPath  string
	width, height int
}

func NewModel() Model {
	files, _ := git.ListChangedFiles(TargetBranch)
	items := tree.Build(files)

	l := list.New(items, listDelegate{}, 0, 0)
	l.SetShowTitle(false)
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	m := Model{
		fileTree:     l,
		diffViewport: viewport.New(0, 0),
	}

	if len(items) > 0 {
		if first, ok := items[0].(tree.TreeItem); ok {
			m.selectedPath = first.FullPath
		}
	}
	return m
}

func (m Model) Init() tea.Cmd {
	if m.selectedPath != "" {
		return git.DiffCmd(TargetBranch, m.selectedPath)
	}
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		treeWidth := int(float64(m.width) * 0.25)
		m.fileTree.SetSize(treeWidth, m.height)
		m.diffViewport.Width = m.width - treeWidth - 2
		m.diffViewport.Height = m.height

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "k", "down", "j":
			m.fileTree, cmd = m.fileTree.Update(msg)
			cmds = append(cmds, cmd)
			if item, ok := m.fileTree.SelectedItem().(tree.TreeItem); ok && !item.IsDir {
				if item.FullPath != m.selectedPath {
					m.selectedPath = item.FullPath
					cmds = append(cmds, git.DiffCmd(TargetBranch, m.selectedPath))
				}
			}
		case "e":
			if m.selectedPath != "" {
				return m, git.OpenEditorCmd(m.selectedPath)
			}
		}

	case git.DiffMsg:
		m.diffViewport.SetContent(msg.Content)

	case git.EditorFinishedMsg:
		if msg.Err != nil {
		}
		return m, git.DiffCmd(TargetBranch, m.selectedPath)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	treeView := PaneStyle.Copy().
		Width(m.fileTree.Width()).
		Height(m.fileTree.Height()).
		Render(m.fileTree.View())

	diffView := DiffStyle.Copy().
		Width(m.diffViewport.Width).
		Height(m.diffViewport.Height).
		Render(m.diffViewport.View())

	return lipgloss.JoinHorizontal(lipgloss.Top, treeView, diffView)
}

type listDelegate struct{}

func (d listDelegate) Height() int                               { return 1 }
func (d listDelegate) Spacing() int                              { return 0 }
func (d listDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d listDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i, ok := item.(tree.TreeItem)
	if !ok {
		return
	}

	str := i.Title()
	if index == m.Index() {
		fmt.Fprint(w, SelectedItemStyle.Render("│ "+str))
	} else {
		fmt.Fprint(w, ItemStyle.Render(str))
	}
}
