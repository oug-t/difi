package ui

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"

	"github.com/oug-t/difi/internal/config"
	"github.com/oug-t/difi/internal/git"
	"github.com/oug-t/difi/internal/tree"
)

type Focus int

const (
	FocusTree Focus = iota
	FocusDiff
)

type StatsMsg struct {
	Added   int
	Deleted int
}

type Model struct {
	fileList     list.Model
	treeState    *tree.FileTree
	treeDelegate TreeDelegate
	diffViewport viewport.Model

	selectedPath  string
	currentBranch string
	targetBranch  string
	repoName      string

	statsAdded   int
	statsDeleted int

	diffContent string
	diffLines   []string
	diffCursor  int

	inputBuffer string
	pendingZ    bool

	focus    Focus
	showHelp bool

	width, height int
}

func NewModel(cfg config.Config, targetBranch string) Model {
	InitStyles(cfg)

	files, _ := git.ListChangedFiles(targetBranch)

	t := tree.New(files)
	items := t.Items()

	delegate := TreeDelegate{Focused: true}
	l := list.New(items, delegate, 0, 0)

	l.SetShowTitle(false)
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowPagination(false)
	l.DisableQuitKeybindings()

	m := Model{
		fileList:      l,
		treeState:     t,
		treeDelegate:  delegate,
		diffViewport:  viewport.New(0, 0),
		focus:         FocusTree,
		currentBranch: git.GetCurrentBranch(),
		targetBranch:  targetBranch,
		repoName:      git.GetRepoName(),
		showHelp:      false,
		inputBuffer:   "",
		pendingZ:      false,
	}

	if len(items) > 0 {
		if first, ok := items[0].(tree.TreeItem); ok {
			if !first.IsDir {
				m.selectedPath = first.FullPath
			}
		}
	}
	return m
}

func (m Model) Init() tea.Cmd {
	var cmds []tea.Cmd

	if m.selectedPath != "" {
		cmds = append(cmds, git.DiffCmd(m.targetBranch, m.selectedPath))
	}

	cmds = append(cmds, fetchStatsCmd(m.targetBranch))

	return tea.Batch(cmds...)
}

func fetchStatsCmd(target string) tea.Cmd {
	return func() tea.Msg {
		added, deleted, err := git.DiffStats(target)
		if err != nil {
			return nil
		}
		return StatsMsg{Added: added, Deleted: deleted}
	}
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

	case StatsMsg:
		m.statsAdded = msg.Added
		m.statsDeleted = msg.Deleted

	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		if len(m.fileList.Items()) == 0 {
			return m, nil
		}

		if m.pendingZ {
			m.pendingZ = false
			if m.focus == FocusDiff {
				switch msg.String() {
				case "z", ".":
					m.centerDiffCursor()
				case "t":
					m.diffViewport.SetYOffset(m.diffCursor)
				case "b":
					offset := m.diffCursor - m.diffViewport.Height + 1
					if offset < 0 {
						offset = 0
					}
					m.diffViewport.SetYOffset(offset)
				}
			}
			return m, nil
		}

		if len(msg.String()) == 1 && strings.ContainsAny(msg.String(), "0123456789") {
			m.inputBuffer += msg.String()
			return m, nil
		}

		if msg.String() == "?" {
			m.showHelp = !m.showHelp
			m.updateSizes()
			return m, nil
		}

		switch msg.String() {
		case "tab":
			if m.focus == FocusTree {
				if item, ok := m.fileList.SelectedItem().(tree.TreeItem); ok && item.IsDir {
					return m, nil
				}
				m.focus = FocusDiff
			} else {
				m.focus = FocusTree
			}
			m.updateTreeFocus()
			m.inputBuffer = ""

		case "l", "]", "ctrl+l", "right":
			if m.focus == FocusTree {
				if item, ok := m.fileList.SelectedItem().(tree.TreeItem); ok && item.IsDir {
					return m, nil
				}
			}
			m.focus = FocusDiff
			m.updateTreeFocus()
			m.inputBuffer = ""

		case "h", "[", "ctrl+h", "left":
			m.focus = FocusTree
			m.updateTreeFocus()
			m.inputBuffer = ""

		case "enter":
			if m.focus == FocusTree {
				if i, ok := m.fileList.SelectedItem().(tree.TreeItem); ok && i.IsDir {
					m.treeState.ToggleExpand(i.FullPath)
					m.fileList.SetItems(m.treeState.Items())
					return m, nil
				}
			}
			if m.selectedPath != "" {
				if i, ok := m.fileList.SelectedItem().(tree.TreeItem); ok && !i.IsDir {
					// proceed
				} else {
					return m, nil
				}
			}
			fallthrough

		case "e":
			if m.selectedPath != "" {
				if i, ok := m.fileList.SelectedItem().(tree.TreeItem); ok && i.IsDir {
					return m, nil
				}

				line := 0
				if m.focus == FocusDiff {
					line = git.CalculateFileLine(m.diffContent, m.diffCursor)
				} else {
					line = git.CalculateFileLine(m.diffContent, 0)
				}
				m.inputBuffer = ""
				return m, git.OpenEditorCmd(m.selectedPath, line, m.targetBranch)
			}

		case "z":
			if m.focus == FocusDiff {
				m.pendingZ = true
				return m, nil
			}

		case "H":
			if m.focus == FocusDiff {
				m.diffCursor = m.diffViewport.YOffset
				if m.diffCursor >= len(m.diffLines) {
					m.diffCursor = len(m.diffLines) - 1
				}
			}

		case "M":
			if m.focus == FocusDiff {
				half := m.diffViewport.Height / 2
				m.diffCursor = m.diffViewport.YOffset + half
				if m.diffCursor >= len(m.diffLines) {
					m.diffCursor = len(m.diffLines) - 1
				}
			}

		case "L":
			if m.focus == FocusDiff {
				m.diffCursor = m.diffViewport.YOffset + m.diffViewport.Height - 1
				if m.diffCursor >= len(m.diffLines) {
					m.diffCursor = len(m.diffLines) - 1
				}
			}

		case "ctrl+d":
			if m.focus == FocusDiff {
				halfScreen := m.diffViewport.Height / 2
				m.diffCursor += halfScreen
				if m.diffCursor >= len(m.diffLines) {
					m.diffCursor = len(m.diffLines) - 1
				}
				m.centerDiffCursor()
			}
			m.inputBuffer = ""

		case "ctrl+u":
			if m.focus == FocusDiff {
				halfScreen := m.diffViewport.Height / 2
				m.diffCursor -= halfScreen
				if m.diffCursor < 0 {
					m.diffCursor = 0
				}
				m.centerDiffCursor()
			}
			m.inputBuffer = ""

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
					m.fileList.CursorDown()
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
					m.fileList.CursorUp()
				}
			}
			m.inputBuffer = ""

		default:
			m.inputBuffer = ""
		}
	}

	if len(m.fileList.Items()) > 0 && m.focus == FocusTree {
		if !keyHandled {
			m.fileList, cmd = m.fileList.Update(msg)
			cmds = append(cmds, cmd)
		}

		if item, ok := m.fileList.SelectedItem().(tree.TreeItem); ok {
			if !item.IsDir && item.FullPath != m.selectedPath {
				m.selectedPath = item.FullPath
				m.diffCursor = 0
				m.diffViewport.GotoTop()
				cmds = append(cmds, git.DiffCmd(m.targetBranch, m.selectedPath))
			}
		}
	}

	switch msg := msg.(type) {
	case git.DiffMsg:
		m.diffContent = msg.Content
		m.diffLines = strings.Split(msg.Content, "\n")
		m.diffViewport.SetContent(msg.Content)

	case git.EditorFinishedMsg:
		return m, git.DiffCmd(m.targetBranch, m.selectedPath)
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) centerDiffCursor() {
	halfScreen := m.diffViewport.Height / 2
	targetOffset := m.diffCursor - halfScreen
	if targetOffset < 0 {
		targetOffset = 0
	}
	m.diffViewport.SetYOffset(targetOffset)
}

func (m *Model) updateSizes() {
	// 1 line Top Bar + 1 line Bottom Bar = 2 reserved
	reservedHeight := 2
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

	// Subtract border height (2) from contentHeight
	listHeight := contentHeight - 2
	if listHeight < 1 {
		listHeight = 1
	}
	m.fileList.SetSize(treeWidth, listHeight)

	m.diffViewport.Width = m.width - treeWidth - 4 // border (2) + padding (2) from tree pane
	m.diffViewport.Height = listHeight
}

func (m *Model) updateTreeFocus() {
	m.treeDelegate.Focused = (m.focus == FocusTree)
	m.fileList.SetDelegate(m.treeDelegate)
}

func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	topBar := m.renderTopBar()

	var mainContent string
	contentHeight := m.height - 2
	if m.showHelp {
		contentHeight -= 6
	}
	if contentHeight < 0 {
		contentHeight = 0
	}

	if len(m.fileList.Items()) == 0 {
		mainContent = m.renderEmptyState(m.width, contentHeight, "No changes found against "+m.targetBranch)
	} else {
		treeStyle := PaneStyle
		if m.focus == FocusTree {
			treeStyle = FocusedPaneStyle
		} else {
			treeStyle = PaneStyle
		}

		treeView := treeStyle.Copy().
			Width(m.fileList.Width()).
			Height(m.fileList.Height()).
			Render(m.fileList.View())

		var rightPaneView string

		selectedItem, ok := m.fileList.SelectedItem().(tree.TreeItem)
		if ok && selectedItem.IsDir {
			rightPaneView = m.renderEmptyState(m.diffViewport.Width, m.diffViewport.Height, "Directory: "+selectedItem.Name)
		} else {
			var renderedDiff strings.Builder
			start := m.diffViewport.YOffset
			end := start + m.diffViewport.Height
			if end > len(m.diffLines) {
				end = len(m.diffLines)
			}

			// 5 for line number (Width 4 + MarginRight 1), 2 for indent
			maxLineWidth := m.diffViewport.Width - 7
			if maxLineWidth < 1 {
				maxLineWidth = 1
			}

			for i := start; i < end; i++ {
				line := ansi.Truncate(m.diffLines[i], maxLineWidth, "")

				var numStr string
				mode := "relative"

				if mode == "hidden" {
					numStr = ""
				} else {
					isCursor := (i == m.diffCursor)
					if isCursor && mode == "hybrid" {
						realLine := git.CalculateFileLine(m.diffContent, m.diffCursor)
						numStr = fmt.Sprintf("%d", realLine)
					} else if isCursor && mode == "relative" {
						numStr = "0"
					} else if mode == "absolute" {
						numStr = fmt.Sprintf("%d", i+1)
					} else {
						dist := int(math.Abs(float64(i - m.diffCursor)))
						numStr = fmt.Sprintf("%d", dist)
					}
				}

				lineNumRendered := ""
				if numStr != "" {
					lineNumRendered = LineNumberStyle.Render(numStr)
				}

				if m.focus == FocusDiff && i == m.diffCursor {
					cleanLine := stripAnsi(line)
					line = DiffSelectionStyle.Render("  " + cleanLine)
				} else {
					line = "  " + line
				}

				renderedDiff.WriteString(lineNumRendered + line + "\n")
			}

			diffContentStr := strings.TrimRight(renderedDiff.String(), "\n")

			rightPaneView = DiffStyle.Copy().
				Width(m.diffViewport.Width).
				Height(m.diffViewport.Height).
				Render(diffContentStr)
		}

		mainContent = lipgloss.JoinHorizontal(lipgloss.Top, treeView, rightPaneView)
	}

	var bottomBar string
	if m.showHelp {
		bottomBar = m.renderHelpDrawer()
	} else {
		bottomBar = m.viewStatusBar()
	}

	return lipgloss.JoinVertical(lipgloss.Top, topBar, mainContent, bottomBar)
}

func (m Model) renderTopBar() string {
	repo := fmt.Sprintf(" %s", m.repoName)
	branches := fmt.Sprintf(" %s ➜ %s", m.currentBranch, m.targetBranch)
	info := fmt.Sprintf("%s   %s", repo, branches)
	leftSide := TopInfoStyle.Render(info)

	middle := ""
	if m.selectedPath != "" {
		middle = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(m.selectedPath)
	}

	rightSide := ""
	if m.statsAdded > 0 || m.statsDeleted > 0 {
		added := TopStatsAddedStyle.Render(fmt.Sprintf("+%d", m.statsAdded))
		deleted := TopStatsDeletedStyle.Render(fmt.Sprintf("-%d", m.statsDeleted))
		rightSide = lipgloss.JoinHorizontal(lipgloss.Center, added, deleted)
	}

	availWidth := m.width - lipgloss.Width(leftSide) - lipgloss.Width(rightSide)
	if availWidth < 0 {
		availWidth = 0
	}

	midWidth := lipgloss.Width(middle)
	var centerBlock string
	if midWidth > availWidth {
		centerBlock = strings.Repeat(" ", availWidth)
	} else {
		padL := (availWidth - midWidth) / 2
		padR := availWidth - midWidth - padL
		centerBlock = strings.Repeat(" ", padL) + middle + strings.Repeat(" ", padR)
	}

	finalBar := lipgloss.JoinHorizontal(lipgloss.Top, leftSide, centerBlock, rightSide)
	return TopBarStyle.Width(m.width).Render(finalBar)
}

func (m Model) viewStatusBar() string {
	shortcuts := StatusKeyStyle.Render("? Help  q Quit  Tab Switch")
	return StatusBarStyle.Width(m.width).Render(shortcuts)
}

func (m Model) renderHelpDrawer() string {
	col1 := lipgloss.JoinVertical(lipgloss.Left,
		HelpTextStyle.Render("↑/k   Move Up"),
		HelpTextStyle.Render("↓/j   Move Down"),
	)
	col2 := lipgloss.JoinVertical(lipgloss.Left,
		HelpTextStyle.Render("←/h   Left Panel"),
		HelpTextStyle.Render("→/l   Right Panel"),
	)
	col3 := lipgloss.JoinVertical(lipgloss.Left,
		HelpTextStyle.Render("C-d/u Page Dn/Up"),
		HelpTextStyle.Render("zz/zt Scroll View"),
	)
	col4 := lipgloss.JoinVertical(lipgloss.Left,
		HelpTextStyle.Render("H/M/L Move Cursor"),
		HelpTextStyle.Render("e     Edit File"),
	)

	return HelpDrawerStyle.Copy().
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
}

func (m Model) renderEmptyState(w, h int, statusMsg string) string {
	logo := EmptyLogoStyle.Render("difi")
	desc := EmptyDescStyle.Render("A calm, focused way to review Git diffs.")
	status := EmptyStatusStyle.Render(statusMsg)

	usageHeader := EmptyHeaderStyle.Render("Usage Patterns")
	cmd1 := lipgloss.NewStyle().Foreground(ColorText).Render("difi")
	desc1 := EmptyCodeStyle.Render("Diff against main")
	cmd2 := lipgloss.NewStyle().Foreground(ColorText).Render("difi dev")
	desc2 := EmptyCodeStyle.Render("Diff against branch")

	usageBlock := lipgloss.JoinVertical(lipgloss.Left,
		usageHeader,
		lipgloss.JoinHorizontal(lipgloss.Left, cmd1, "    ", desc1),
		lipgloss.JoinHorizontal(lipgloss.Left, cmd2, "    ", desc2),
	)

	navHeader := EmptyHeaderStyle.Render("Navigation")
	key1 := lipgloss.NewStyle().Foreground(ColorText).Render("Tab")
	key2 := lipgloss.NewStyle().Foreground(ColorText).Render("j/k")
	keyDesc1 := EmptyCodeStyle.Render("Switch panels")
	keyDesc2 := EmptyCodeStyle.Render("Move cursor")

	navBlock := lipgloss.JoinVertical(lipgloss.Left,
		navHeader,
		lipgloss.JoinHorizontal(lipgloss.Left, key1, "    ", keyDesc1),
		lipgloss.JoinHorizontal(lipgloss.Left, key2, "    ", keyDesc2),
	)

	nvimHeader := EmptyHeaderStyle.Render("Neovim Integration")
	nvim1 := lipgloss.NewStyle().Foreground(ColorText).Render("oug-t/difi.nvim")
	nvimDesc1 := EmptyCodeStyle.Render("Install plugin")
	nvim2 := lipgloss.NewStyle().Foreground(ColorText).Render("Press 'e'")
	nvimDesc2 := EmptyCodeStyle.Render("Edit with context")

	nvimBlock := lipgloss.JoinVertical(lipgloss.Left,
		nvimHeader,
		lipgloss.JoinHorizontal(lipgloss.Left, nvim1, "  ", nvimDesc1),
		lipgloss.JoinHorizontal(lipgloss.Left, nvim2, "          ", nvimDesc2),
	)

	var guides string
	if w > 80 {
		guides = lipgloss.JoinHorizontal(lipgloss.Top,
			usageBlock,
			lipgloss.NewStyle().Width(6).Render(""),
			navBlock,
			lipgloss.NewStyle().Width(6).Render(""),
			nvimBlock,
		)
	} else {
		topRow := lipgloss.JoinHorizontal(lipgloss.Top, usageBlock, lipgloss.NewStyle().Width(4).Render(""), navBlock)
		guides = lipgloss.JoinVertical(lipgloss.Left,
			topRow,
			lipgloss.NewStyle().Height(1).Render(""),
			nvimBlock,
		)
	}

	content := lipgloss.JoinVertical(lipgloss.Center,
		logo,
		desc,
		status,
		lipgloss.NewStyle().Height(1).Render(""),
		guides,
	)

	// Use lipgloss.Place to center the content in the available space
	// This automatically handles vertical and horizontal centering.
	return lipgloss.Place(w, h, lipgloss.Center, lipgloss.Center, content)
}

func stripAnsi(str string) string {
	re := regexp.MustCompile("[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))")
	return re.ReplaceAllString(str, "")
}
