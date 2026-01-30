package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Colors struct {
		Border     string `yaml:"border"`
		Focus      string `yaml:"focus"`
		LineNumber string `yaml:"line_number"`
	} `yaml:"colors"`
	UI struct {
		LineNumbers string `yaml:"line_numbers"` // "absolute", "relative", "hybrid", "hidden"
		ShowGuide   bool   `yaml:"show_guide"`   // The vertical separation line
	} `yaml:"ui"`
}

func DefaultConfig() Config {
	var c Config
	c.Colors.Border = "#D9DCCF"
	c.Colors.Focus = "#000000" // Default neutral focus
	c.Colors.LineNumber = "#808080"
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

	// Parse YAML
	_ = yaml.Unmarshal(data, &cfg)
	return cfg
}
