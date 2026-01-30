package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Colors struct {
		Border          string `yaml:"border"`
		Focus           string `yaml:"focus"`
		LineNumber      string `yaml:"line_number"`
		DiffSelectionBg string `yaml:"diff_selection_bg"` // New config
	} `yaml:"colors"`
	UI struct {
		LineNumbers string `yaml:"line_numbers"`
		ShowGuide   bool   `yaml:"show_guide"`
	} `yaml:"ui"`
}

func DefaultConfig() Config {
	var c Config
	c.Colors.Border = "#D9DCCF"
	c.Colors.Focus = "#6e7781"
	c.Colors.LineNumber = "#808080"

	// Default: "Neutral Light Transparent Blue"
	// Dark Mode: Deep subtle blue-grey | Light Mode: Very faint blue
	// We only set one default here, but AdaptiveColor handles the split in styles.go
	c.Colors.DiffSelectionBg = "" // Empty means use internal defaults

	c.UI.LineNumbers = "hybrid"
	c.UI.ShowGuide = true
	return c
}

func Load() Config {
	cfg := DefaultConfig()
	home, err := os.UserHomeDir()
	if err != nil {
		return cfg
	}

	configPath := filepath.Join(home, ".config", "difi", "config.yml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return cfg
	}

	_ = yaml.Unmarshal(data, &cfg)
	return cfg
}
