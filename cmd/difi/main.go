package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/oug-t/difi/internal/config"
	"github.com/oug-t/difi/internal/ui"
)

var version = "dev"

func main() {
	showVersion := flag.Bool("version", false, "Show version")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: difi [flags] [target-branch]\n")
		fmt.Fprintf(os.Stderr, "\nFlags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  difi              # Diff against main\n")
		fmt.Fprintf(os.Stderr, "  difi develop      # Diff against develop\n")
		fmt.Fprintf(os.Stderr, "  difi HEAD~1       # Diff against last commit\n")
	}

	flag.Parse()

	if *showVersion {
		fmt.Printf("difi version %s\n", version)
		os.Exit(0)
	}

	target := "main"
	if flag.NArg() > 0 {
		target = flag.Arg(0)
	}

	cfg := config.Load()

	p := tea.NewProgram(ui.NewModel(cfg, target), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
