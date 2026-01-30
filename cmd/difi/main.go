package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/oug-t/difi/internal/config"
	"github.com/oug-t/difi/internal/ui"
)

func main() {
	// Load config (defaults used if file missing)
	cfg := config.Load()

	// Pass config to model
	p := tea.NewProgram(ui.NewModel(cfg), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running difi: %v", err)
		os.Exit(1)
	}
}
