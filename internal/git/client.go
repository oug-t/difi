package git

import (
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

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

func OpenEditorCmd(path string) tea.Cmd {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}
	c := exec.Command(editor, path)
	c.Stdin, c.Stdout, c.Stderr = os.Stdin, os.Stdout, os.Stderr
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return EditorFinishedMsg{Err: err}
	})
}

type DiffMsg struct{ Content string }
type EditorFinishedMsg struct{ Err error }
