package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/oug-t/difi/internal/config"
	"github.com/oug-t/difi/internal/ui"
)

func main() {
	// Define flags
	help := flag.Bool("help", false, "Show help")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: difi [flags] [target-branch]\n")
		fmt.Fprintf(os.Stderr, "\nFlags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  difi             # Diff against main\n")
		fmt.Fprintf(os.Stderr, "  difi develop     # Diff against develop\n")
		fmt.Fprintf(os.Stderr, "  difi HEAD~1      # Diff against last commit\n")
	}
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	// Determine Target Branch
	// If user provides an argument (e.g., "difi develop"), use it.
	// Otherwise default to "main".
	target := "main"
	if flag.NArg() > 0 {
		target = flag.Arg(0)
	}

	// Load Config
	cfg := config.Load()

	// Pass target to the model
	p := tea.NewProgram(ui.NewModel(cfg, target), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
