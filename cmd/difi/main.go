package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/oug-t/difi/internal/config"
	"github.com/oug-t/difi/internal/ui"
)

var version = "dev"

func main() {
	showVersion := flag.Bool("version", false, "Show version")
	plain := flag.Bool("plain", false, "Print a plain summary")
	flag.Parse()

	if *showVersion {
		fmt.Printf("difi version %s\n", version)
		os.Exit(0)
	}

	var pipedDiff string
	if stat, _ := os.Stdin.Stat(); (stat.Mode() & os.ModeCharDevice) == 0 {
		b, _ := io.ReadAll(os.Stdin)
		pipedDiff = string(b)
	}

	target := "HEAD"
	if flag.NArg() > 0 {
		target = flag.Arg(0)
	}

	if *plain && pipedDiff == "" {
		cmd := exec.Command("git", "diff", "--name-status", fmt.Sprintf("%s...HEAD", target))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			os.Exit(1)
		}
		os.Exit(0)
	}

	cfg := config.Load()

	opts := []tea.ProgramOption{tea.WithAltScreen()}
	if pipedDiff != "" {
		if tty, err := os.Open("/dev/tty"); err == nil {
			opts = append(opts, tea.WithInput(tty))
		}
	}

	p := tea.NewProgram(ui.NewModel(cfg, target, pipedDiff), opts...)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
