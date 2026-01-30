package config

type Config struct {
	UI UIConfig
}

type UIConfig struct {
	LineNumbers string // "relative", "absolute", "hybrid", "hidden"
	Theme       string
}

func Load() Config {
	// Default configuration
	return Config{
		UI: UIConfig{
			LineNumbers: "relative", // Default to relative numbers (vim style)
			Theme:       "default",
		},
	}
}
