package ui

import (
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/oug-t/difi/internal/config"
	"github.com/oug-t/difi/internal/git"
	"github.com/oug-t/difi/internal/tree"
)

const TargetBranch = "main"

type Focus int

const (
	FocusTree Focus = iota
	FocusDiff
)

type Model struct {
	fileTree     list.Model
	diffViewport viewport.Model

	selectedPath  string
	currentBranch string
	repoName      string

	diffContent string
	diffLines   []string
	diffCursor  int

	inputBuffer string

	focus    Focus
	showHelp bool

	width, height int
}

func NewModel(cfg config.Config) Model {
	// Initialize styles with the loaded config
	InitStyles(cfg)

	files, _ := git.ListChangedFiles(TargetBranch)
	items := tree.Build(files)

	l := list.New(items, listDelegate{}, 0, 0)
	l.SetShowTitle(false)
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.DisableQuitKeybindings()

	m := Model{
		fileTree:      l,
		diffViewport:  viewport.New(0, 0),
		focus:         FocusTree,
		currentBranch: git.GetCurrentBranch(),
		repoName:      git.GetRepoName(),
		showHelp:      false,
		inputBuffer:   "",
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

func (m *Model) getRepeatCount() int {
	if m.inputBuffer == "" {
		return 1
	}
	count, err := strconv.Atoi(m.inputBuffer)
	if err != nil {
		return 1
	}
	m.inputBuffer = ""
	return count
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	keyHandled := false

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateSizes()

	case tea.KeyMsg:
		if len(msg.String()) == 1 && strings.ContainsAny(msg.String(), "0123456789") {
			m.inputBuffer += msg.String()
			return m, nil
		}

		if msg.String() == "?" {
			m.showHelp = !m.showHelp
			m.updateSizes()
			return m, nil
		}

		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		switch msg.String() {
		case "tab":
			if m.focus == FocusTree {
				m.focus = FocusDiff
			} else {
				m.focus = FocusTree
			}
			m.inputBuffer = ""

		case "l", "]", "ctrl+l", "right":
			m.focus = FocusDiff
			m.inputBuffer = ""

		case "h", "[", "ctrl+h", "left":
			m.focus = FocusTree
			m.inputBuffer = ""

		case "e", "enter":
			if m.selectedPath != "" {
				line := 0
				if m.focus == FocusDiff {
					line = git.CalculateFileLine(m.diffContent, m.diffCursor)
				} else {
					line = git.CalculateFileLine(m.diffContent, 0)
				}
				m.inputBuffer = ""
				return m, git.OpenEditorCmd(m.selectedPath, line)
			}

		case "j", "down":
			keyHandled = true
			count := m.getRepeatCount()
			for i := 0; i < count; i++ {
				if m.focus == FocusDiff {
					if m.diffCursor < len(m.diffLines)-1 {
						m.diffCursor++
						if m.diffCursor >= m.diffViewport.YOffset+m.diffViewport.Height {
							m.diffViewport.LineDown(1)
						}
					}
				} else {
					m.fileTree.CursorDown()
				}
			}
			m.inputBuffer = ""

		case "k", "up":
			keyHandled = true
			count := m.getRepeatCount()
			for i := 0; i < count; i++ {
				if m.focus == FocusDiff {
					if m.diffCursor > 0 {
						m.diffCursor--
						if m.diffCursor < m.diffViewport.YOffset {
							m.diffViewport.LineUp(1)
						}
					}
				} else {
					m.fileTree.CursorUp()
				}
			}
			m.inputBuffer = ""

		default:
			m.inputBuffer = ""
		}
	}

	if m.focus == FocusTree {
		if !keyHandled {
			m.fileTree, cmd = m.fileTree.Update(msg)
			cmds = append(cmds, cmd)
		}

		if item, ok := m.fileTree.SelectedItem().(tree.TreeItem); ok && !item.IsDir {
			if item.FullPath != m.selectedPath {
				m.selectedPath = item.FullPath
				m.diffCursor = 0
				m.diffViewport.GotoTop()
				cmds = append(cmds, git.DiffCmd(TargetBranch, m.selectedPath))
			}
		}
	}

	switch msg := msg.(type) {
	case git.DiffMsg:
		m.diffContent = msg.Content
		m.diffLines = strings.Split(msg.Content, "\n")
		m.diffViewport.SetContent(msg.Content)

	case git.EditorFinishedMsg:
		return m, git.DiffCmd(TargetBranch, m.selectedPath)
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) updateSizes() {
	reservedHeight := 1
	if m.showHelp {
		reservedHeight += 6
	}

	contentHeight := m.height - reservedHeight
	if contentHeight < 1 {
		contentHeight = 1
	}

	treeWidth := int(float64(m.width) * 0.20)
	if treeWidth < 20 {
		treeWidth = 20
	}

	m.fileTree.SetSize(treeWidth, contentHeight)
	m.diffViewport.Width = m.width - treeWidth - 2
	m.diffViewport.Height = contentHeight
}

func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	// 1. PANES
	treeStyle := PaneStyle
	if m.focus == FocusTree {
		treeStyle = FocusedPaneStyle
	} else {
		treeStyle = PaneStyle
	}

	treeView := treeStyle.Copy().
		Width(m.fileTree.Width()).
		Height(m.fileTree.Height()).
		Render(m.fileTree.View())

	var renderedDiff strings.Builder
	start := m.diffViewport.YOffset
	end := start + m.diffViewport.Height
	if end > len(m.diffLines) {
		end = len(m.diffLines)
	}

	// RENDER LOOP
	for i := start; i < end; i++ {
		line := m.diffLines[i]

		// --- LINE NUMBER LOGIC ---
		var numStr string
		mode := CurrentConfig.UI.LineNumbers

		if mode == "hidden" {
			numStr = ""
		} else {
			// Is this the cursor line?
			isCursor := (i == m.diffCursor)

			if isCursor && mode == "hybrid" {
				// HYBRID: Show Real File Line Number
				realLine := git.CalculateFileLine(m.diffContent, m.diffCursor)
				numStr = fmt.Sprintf("%d", realLine)
			} else if isCursor && mode == "relative" {
				numStr = "0"
			} else if mode == "absolute" {
				// Note: Calculating absolute for every line is expensive,
				// usually absolute view shows Diff Line Index or File Line.
				// For simple 'absolute' view, we often show viewport index + 1
				numStr = fmt.Sprintf("%d", i+1)
			} else {
				// Default / Hybrid-non-cursor: Show Relative Distance
				dist := int(math.Abs(float64(i - m.diffCursor)))
				numStr = fmt.Sprintf("%d", dist)
			}
		}

		lineNumRendered := ""
		if numStr != "" {
			lineNumRendered = LineNumberStyle.Render(numStr)
		}
		// -------------------------

		if m.focus == FocusDiff && i == m.diffCursor {
			line = SelectedItemStyle.Render(line)
		} else {
			line = "  " + line
		}

		renderedDiff.WriteString(lineNumRendered + line + "\n")
	}

	diffView := DiffStyle.Copy().
		Width(m.diffViewport.Width).
		Height(m.diffViewport.Height).
		Render(renderedDiff.String())

	mainPanes := lipgloss.JoinHorizontal(lipgloss.Top, treeView, diffView)

	// 2. BOTTOM AREA
	repoSection := StatusKeyStyle.Render(" " + m.repoName)
	divider := StatusDividerStyle.Render("│")

	statusText := fmt.Sprintf(" %s ↔ %s", m.currentBranch, TargetBranch)
	if m.inputBuffer != "" {
		statusText += fmt.Sprintf(" [Cmd: %s]", m.inputBuffer)
	}
	branchSection := StatusBarStyle.Render(statusText)

	leftStatus := lipgloss.JoinHorizontal(lipgloss.Center, repoSection, divider, branchSection)
	rightStatus := StatusBarStyle.Render("? Help")

	statusBar := StatusBarStyle.Copy().
		Width(m.width).
		Render(lipgloss.JoinHorizontal(lipgloss.Top,
			leftStatus,
			lipgloss.PlaceHorizontal(m.width-lipgloss.Width(leftStatus)-lipgloss.Width(rightStatus), lipgloss.Right, rightStatus),
		))

	var finalView string
	if m.showHelp {
		col1 := lipgloss.JoinVertical(lipgloss.Left,
			HelpTextStyle.Render("↑/k   Move Up"),
			HelpTextStyle.Render("↓/j   Move Down"),
		)
		col2 := lipgloss.JoinVertical(lipgloss.Left,
			HelpTextStyle.Render("←/h   Left Panel"),
			HelpTextStyle.Render("→/l   Right Panel"),
		)
		col3 := lipgloss.JoinVertical(lipgloss.Left,
			HelpTextStyle.Render("Tab   Switch Panel"),
			HelpTextStyle.Render("Num   Motion Count"),
		)
		col4 := lipgloss.JoinVertical(lipgloss.Left,
			HelpTextStyle.Render("e     Edit File"),
			HelpTextStyle.Render("?     Close Help"),
		)

		helpDrawer := HelpDrawerStyle.Copy().
			Width(m.width).
			Render(lipgloss.JoinHorizontal(lipgloss.Top,
				col1,
				lipgloss.NewStyle().Width(4).Render(""),
				col2,
				lipgloss.NewStyle().Width(4).Render(""),
				col3,
				lipgloss.NewStyle().Width(4).Render(""),
				col4,
			))

		finalView = lipgloss.JoinVertical(lipgloss.Top, mainPanes, helpDrawer, statusBar)
	} else {
		finalView = lipgloss.JoinVertical(lipgloss.Top, mainPanes, statusBar)
	}

	return finalView
}

// -- Delegates (unchanged) --
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
		fmt.Fprint(w, SelectedItemStyle.Render(str))
	} else {
		fmt.Fprint(w, ItemStyle.Render(str))
	}
}
