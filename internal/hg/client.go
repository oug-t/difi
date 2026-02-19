package hg

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// hgRoot caches the repository root directory.
var hgRoot string

// ansiRe matches ANSI escape sequences for stripping from terminal output.
var ansiRe = regexp.MustCompile("[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))")

// hunkHeaderRe matches unified diff hunk headers: @@ -l,s +l,s @@
var hunkHeaderRe = regexp.MustCompile(`^.*?@@ \-\d+(?:,\d+)? \+(\d+)(?:,\d+)? @@`)

func getHgRoot() string {
	if hgRoot != "" {
		return hgRoot
	}
	cmd := exec.Command("hg", "root")
	cmd.Env = append(os.Environ(), "HGRCPATH="+os.DevNull)
	out, err := cmd.Output()
	if err == nil {
		hgRoot = strings.TrimSpace(string(out))
	}
	return hgRoot
}

// hgCmd creates an hg command with HGRCPATH set to the platform null device
// to ignore all user config (pagers, aliases, hooks, etc.) without suppressing
// --color=always. Commands run from the repo root so relative file paths
// resolve correctly.
func hgCmd(args ...string) *exec.Cmd {
	cmd := exec.Command("hg", args...)
	cmd.Env = append(os.Environ(), "HGRCPATH="+os.DevNull)
	if root := getHgRoot(); root != "" {
		cmd.Dir = root
	}
	return cmd
}

func GetCurrentBranch() string {
	out, err := hgCmd("branch").Output()
	if err != nil {
		return "default"
	}
	return strings.TrimSpace(string(out))
}

func GetRepoName() string {
	out, err := hgCmd("root").Output()
	if err != nil {
		return "Repo"
	}
	path := strings.TrimSpace(string(out))
	parts := strings.Split(path, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return "Repo"
}

func ListChangedFiles(targetBranch string) ([]string, error) {
	var cmd *exec.Cmd
	// For working directory changes, use status without --rev
	// For specific revisions, use status with --rev
	if targetBranch == "tip" || targetBranch == "." || targetBranch == "" {
		cmd = hgCmd("status", "--no-status")
	} else {
		cmd = hgCmd("status", "--rev", targetBranch, "--no-status")
	}

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	files := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(files) == 1 && files[0] == "" {
		return []string{}, nil
	}
	return files, nil
}

func DiffCmd(targetBranch, path string) tea.Cmd {
	return func() tea.Msg {
		var cmd *exec.Cmd
		// For working directory diffs, don't use --rev
		// For specific revisions, use --rev
		if targetBranch == "tip" || targetBranch == "." || targetBranch == "" {
			cmd = hgCmd("diff", "--color=always", path)
		} else {
			cmd = hgCmd("diff", "--color=always", "--rev", targetBranch, path)
		}

		out, err := cmd.Output()
		if err != nil {
			return DiffMsg{Content: "Error fetching diff: " + err.Error()}
		}
		return DiffMsg{Content: string(out)}
	}
}

func OpenEditorCmd(path string, lineNumber int, targetBranch string) tea.Cmd {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		if _, err := exec.LookPath("nvim"); err == nil {
			editor = "nvim"
		} else {
			editor = "vim"
		}
	}

	var args []string
	if lineNumber > 0 {
		args = append(args, fmt.Sprintf("+%d", lineNumber))
	}
	args = append(args, path)

	c := exec.Command(editor, args...)
	c.Stdin, c.Stdout, c.Stderr = os.Stdin, os.Stdout, os.Stderr
	if root := getHgRoot(); root != "" {
		c.Dir = root
	}

	c.Env = append(os.Environ(), fmt.Sprintf("DIFI_TARGET=%s", targetBranch))

	return tea.ExecProcess(c, func(err error) tea.Msg {
		return EditorFinishedMsg{Err: err}
	})
}

func DiffStats(targetBranch string) (added int, deleted int, err error) {
	var cmd *exec.Cmd
	// For working directory diffs, don't use --rev
	if targetBranch == "tip" || targetBranch == "." || targetBranch == "" {
		cmd = hgCmd("diff", "--stat")
	} else {
		cmd = hgCmd("diff", "--rev", targetBranch, "--stat")
	}

	out, err := cmd.Output()
	if err != nil {
		return 0, 0, fmt.Errorf("hg diff stats error: %w", err)
	}

	// Parse Mercurial stat output format: " N files changed, M insertions(+), K deletions(-)"
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, line := range lines {
		if strings.Contains(line, "changed") && (strings.Contains(line, "insertion") || strings.Contains(line, "deletion")) {
			// Extract insertions and deletions from summary line
			re := regexp.MustCompile(`(\d+) insertion[s]?\(\+\)`)
			if matches := re.FindStringSubmatch(line); len(matches) > 1 {
				if n, err := strconv.Atoi(matches[1]); err == nil {
					added = n
				}
			}

			re = regexp.MustCompile(`(\d+) deletion[s]?\(\-\)`)
			if matches := re.FindStringSubmatch(line); len(matches) > 1 {
				if n, err := strconv.Atoi(matches[1]); err == nil {
					deleted = n
				}
			}
			break
		}
	}
	return added, deleted, nil
}

func DiffStatsByFile(targetBranch string) (map[string][2]int, error) {
	var cmd *exec.Cmd
	if targetBranch == "tip" || targetBranch == "." || targetBranch == "" {
		cmd = hgCmd("diff", "--stat")
	} else {
		cmd = hgCmd("diff", "--rev", targetBranch, "--stat")
	}

	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("hg diff stat error: %w", err)
	}

	result := make(map[string][2]int)
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, line := range lines {
		// Skip summary line like " 3 files changed, 10 insertions(+), 5 deletions(-)"
		if strings.Contains(line, "changed") && (strings.Contains(line, "insertion") || strings.Contains(line, "deletion")) {
			continue
		}
		// Per-file line: " path/to/file |  5 ++--"
		pipeIdx := strings.LastIndex(line, "|")
		if pipeIdx < 0 {
			continue
		}
		filePath := strings.TrimSpace(line[:pipeIdx])
		changesPart := strings.TrimSpace(line[pipeIdx+1:])
		// changesPart is like "5 ++-" or "3 +++"
		var a, d int
		for _, ch := range changesPart {
			if ch == '+' {
				a++
			} else if ch == '-' {
				d++
			}
		}
		if filePath != "" {
			result[filePath] = [2]int{a, d}
		}
	}
	return result, nil
}

func CalculateFileLine(diffContent string, visualLineIndex int) int {
	lines := strings.Split(diffContent, "\n")
	if visualLineIndex >= len(lines) {
		return 0
	}

	// Mercurial uses similar diff format to Git: @@ -l,s +l,s @@
	currentLineNo := 0
	lastWasHunk := false

	for i := 0; i <= visualLineIndex; i++ {
		line := lines[i]

		matches := hunkHeaderRe.FindStringSubmatch(line)
		if len(matches) > 1 {
			startLine, _ := strconv.Atoi(matches[1])
			currentLineNo = startLine
			lastWasHunk = true
			continue
		}

		lastWasHunk = false
		cleanLine := stripAnsi(line)
		if strings.HasPrefix(cleanLine, " ") || strings.HasPrefix(cleanLine, "+") {
			currentLineNo++
		}
	}

	if currentLineNo == 0 {
		return 1
	}
	// Only adjust by -1 when cursor is on a context/added line (which
	// incremented past the current position). On a hunk header the
	// line number is already exact.
	if lastWasHunk {
		return currentLineNo
	}
	return currentLineNo - 1
}

func stripAnsi(str string) string {
	return ansiRe.ReplaceAllString(str, "")
}

type DiffMsg struct{ Content string }
type EditorFinishedMsg struct{ Err error }

func ParseFilesFromDiff(diffText string) []string {
	var files []string
	seen := make(map[string]bool)
	lines := strings.Split(diffText, "\n")

	for _, line := range lines {
		if strings.HasPrefix(line, "diff -r ") {
			// Mercurial diff format: "diff -r <rev> <file>" (working dir)
			// or "diff -r <rev1> -r <rev2> <file>" (two revisions).
			// The file path is always the last whitespace-separated field.
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				file := parts[len(parts)-1]
				if !seen[file] {
					seen[file] = true
					files = append(files, file)
				}
			}
		}
	}
	return files
}

func ExtractFileDiff(diffText, targetPath string) string {
	lines := strings.Split(diffText, "\n")
	var out []string
	inTarget := false

	for _, line := range lines {
		if strings.HasPrefix(line, "diff -r ") {
			// Mercurial diff format: "diff -r <rev> <file>" (working dir)
			// or "diff -r <rev1> -r <rev2> <file>" (two revisions).
			// The file path is always the last whitespace-separated field.
			parts := strings.Fields(line)
			inTarget = len(parts) > 0 && parts[len(parts)-1] == targetPath
		}

		if inTarget {
			out = append(out, line)
		}
	}

	return strings.Join(out, "\n")
}
