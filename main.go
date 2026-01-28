package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const targetBranch = "main"

// --- STYLES ---

var (
	// Modern, Clean Theme
	colorBorder   = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	colorSelected = lipgloss.AdaptiveColor{Light: "#1F1F1F", Dark: "#F8F8F2"}
	colorInactive = lipgloss.AdaptiveColor{Light: "#A8A8A8", Dark: "#626262"}
	colorAccent   = lipgloss.AdaptiveColor{Light: "#00BCF0", Dark: "#00BCF0"} // Dash-like Cyan

	// Panes
	treeStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, true, false, false). // Right border only
			BorderForeground(colorBorder).
			MarginRight(1)

	diffStyle = lipgloss.NewStyle().
			Padding(0, 1)

	// Tree Items
	itemStyle         = lipgloss.NewStyle().PaddingLeft(1)
	selectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(0).
				Foreground(colorSelected).
				Bold(true)
)

// --- DATA STRUCTURES ---

type node struct {
	name     string
	fullPath string
	children map[string]*node
	isDir    bool
}

type treeItem struct {
	path     string
	fullPath string
	isDir    bool
	depth    int
}

func (i treeItem) FilterValue() string { return i.fullPath }
func (i treeItem) Description() string { return "" }

// Title renders the Nerd Font icon + filename
func (i treeItem) Title() string {
	indent := strings.Repeat("  ", i.depth)
	icon := getIcon(i.path, i.isDir)
	return fmt.Sprintf("%s%s %s", indent, icon, i.path)
}

// --- NERD FONT ICON MAPPING ---

func getIcon(name string, isDir bool) string {
	if isDir {
		return "" // Folder icon
	}

	ext := filepath.Ext(name)
	switch strings.ToLower(ext) {
	case ".go":
		return "" // Go Gopher
	case ".js", ".ts":
		return "" // JS/Node
	case ".md":
		return "" // Markdown
	case ".json":
		return "" // JSON
	case ".yml", ".yaml":
		return "" // Settings/Config
	case ".html":
		return ""
	case ".css":
		return ""
	case ".git":
		return "" // Git
	case ".dockerfile", "dockerfile":
		return "" // Docker
	default:
		return "" // Default File
	}
}

// --- DELEGATE (CUSTOM RENDERER) ---

type delegate struct{}

func (d delegate) Height() int                               { return 1 }
func (d delegate) Spacing() int                              { return 0 }
func (d delegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d delegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i, ok := item.(treeItem)
	if !ok {
		return
	}

	str := i.Title()

	// If selected, add a border or indicator
	if index == m.Index() {
		fmt.Fprint(w, selectedItemStyle.Render("│ "+str))
	} else {
		// Render unselected items with subtle gray
		fmt.Fprint(w, itemStyle.Render(str))
	}
}

// --- MODEL ---

type model struct {
	fileTree      list.Model
	diffViewport  viewport.Model
	selectedPath  string
	width, height int
}

func initialModel() model {
	// 1. Fetch changed files
	// Compares HEAD against targetBranch
	cmd := exec.Command("git", "diff", "--name-only", targetBranch)
	out, _ := cmd.Output()
	filePaths := strings.Split(strings.TrimSpace(string(out)), "\n")

	// Handle empty diff case
	if len(filePaths) == 1 && filePaths[0] == "" {
		filePaths = []string{}
	}

	// 2. Build Tree
	items := buildTree(filePaths)

	// 3. Configure List
	l := list.New(items, delegate{}, 0, 0)
	l.SetShowTitle(false)
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	vp := viewport.New(0, 0)

	m := model{
		fileTree:     l,
		diffViewport: vp,
	}

	// Select first file if available
	if len(items) > 0 {
		if first, ok := items[0].(treeItem); ok {
			m.selectedPath = first.fullPath
		}
	}

	return m
}

func (m model) Init() tea.Cmd {
	if m.selectedPath != "" {
		return fetchDiff(m.selectedPath)
	}
	return nil
}

// --- UPDATE ---

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Layout: 25% Tree (Fixed Left), Rest Diff
		treeWidth := int(float64(m.width) * 0.25)
		diffWidth := m.width - treeWidth - 2

		m.fileTree.SetSize(treeWidth, m.height)
		m.diffViewport.Width = diffWidth
		m.diffViewport.Height = m.height

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "up", "k", "down", "j":
			m.fileTree, cmd = m.fileTree.Update(msg)
			cmds = append(cmds, cmd)

			// Fetch diff only if actual file selected
			if item, ok := m.fileTree.SelectedItem().(treeItem); ok && !item.isDir {
				if item.fullPath != m.selectedPath {
					m.selectedPath = item.fullPath
					cmds = append(cmds, fetchDiff(m.selectedPath))
				}
			}

		case "e":
			// Edit in NVIM
			if m.selectedPath != "" {
				return m, openEditor(m.selectedPath)
			}
		}

	case diffMsg:
		m.diffViewport.SetContent(string(msg))

	case editorFinishedMsg:
		return m, fetchDiff(m.selectedPath)
	}

	return m, tea.Batch(cmds...)
}

// --- VIEW ---

func (m model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	treeView := treeStyle.Copy().
		Width(m.fileTree.Width()).
		Height(m.fileTree.Height()).
		Render(m.fileTree.View())

	diffView := diffStyle.Copy().
		Width(m.diffViewport.Width).
		Height(m.diffViewport.Height).
		Render(m.diffViewport.View())

	return lipgloss.JoinHorizontal(lipgloss.Top, treeView, diffView)
}

// --- COMMANDS ---

type diffMsg string
type editorFinishedMsg struct{ err error }

func fetchDiff(path string) tea.Cmd {
	return func() tea.Msg {
		// Use --color=always to let git handle syntax highlighting
		out, _ := exec.Command("git", "diff", "--color=always", targetBranch, "--", path).Output()
		return diffMsg(out)
	}
}

func openEditor(path string) tea.Cmd {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}
	c := exec.Command(editor, path)
	c.Stdin, c.Stdout, c.Stderr = os.Stdin, os.Stdout, os.Stderr
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return editorFinishedMsg{err}
	})
}

// --- TREE ALGORITHMS ---

func buildTree(paths []string) []list.Item {
	root := &node{children: make(map[string]*node)}

	for _, path := range paths {
		parts := strings.Split(path, "/")
		current := root
		for i, part := range parts {
			if _, exists := current.children[part]; !exists {
				isDir := i < len(parts)-1
				fullPath := strings.Join(parts[:i+1], "/")
				current.children[part] = &node{
					name:     part,
					fullPath: fullPath,
					children: make(map[string]*node),
					isDir:    isDir,
				}
			}
			current = current.children[part]
		}
	}

	var items []list.Item
	flatten(root, 0, &items)
	return items
}

func flatten(n *node, depth int, items *[]list.Item) {
	keys := make([]string, 0, len(n.children))
	for k := range n.children {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		a, b := n.children[keys[i]], n.children[keys[j]]
		if a.isDir && !b.isDir {
			return true
		}
		if !a.isDir && b.isDir {
			return false
		}
		return a.name < b.name
	})

	for _, k := range keys {
		child := n.children[k]
		*items = append(*items, treeItem{
			path:     child.name,
			fullPath: child.fullPath,
			isDir:    child.isDir,
			depth:    depth,
		})
		if child.isDir {
			flatten(child, depth+1, items)
		}
	}
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
