package git

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func GetCurrentBranch() string {
	out, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	if err != nil {
		return "HEAD"
	}
	return strings.TrimSpace(string(out))
}

func GetRepoName() string {
	out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
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
	cmd := exec.Command("git", "diff", "--name-only", targetBranch)
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
		out, err := exec.Command("git", "diff", "--color=always", targetBranch, "--", path).Output()
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

	// Pass the diff target branch to the editor via environment variable
	// This enables plugins like difi.nvim to auto-configure the view
	c.Env = append(os.Environ(), fmt.Sprintf("DIFI_TARGET=%s", targetBranch))

	return tea.ExecProcess(c, func(err error) tea.Msg {
		return EditorFinishedMsg{Err: err}
	})
}

func CalculateFileLine(diffContent string, visualLineIndex int) int {
	lines := strings.Split(diffContent, "\n")
	if visualLineIndex >= len(lines) {
		return 0
	}

	re := regexp.MustCompile(`^.*?@@ \-\d+(?:,\d+)? \+(\d+)(?:,\d+)? @@`)

	currentLineNo := 0

	for i := 0; i <= visualLineIndex; i++ {
		line := lines[i]

		matches := re.FindStringSubmatch(line)
		if len(matches) > 1 {
			startLine, _ := strconv.Atoi(matches[1])
			currentLineNo = startLine
			continue
		}

		cleanLine := stripAnsi(line)
		if strings.HasPrefix(cleanLine, " ") || strings.HasPrefix(cleanLine, "+") {
			currentLineNo++
		}
	}

	if currentLineNo == 0 {
		return 1
	}
	return currentLineNo - 1
}

func stripAnsi(str string) string {
	const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"
	var re = regexp.MustCompile(ansi)
	return re.ReplaceAllString(str, "")
}

type DiffMsg struct{ Content string }
type EditorFinishedMsg struct{ Err error }
