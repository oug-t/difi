package config

type Config struct {
	UI UIConfig
}

type UIConfig struct {
	LineNumbers string // "relative", "absolute", "hybrid", "hidden"
	Theme       string
}

func Load() Config {
	return Config{
		UI: UIConfig{
			LineNumbers: "hybrid",
			Theme:       "default",
		},
	}
}
