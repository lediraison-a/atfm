package config

type ThemeConfig struct {
	Background_default   string
	Background_primary   string
	Background_secondary string
	Text_default         string
	Text_light           string
	Text_primary         string
}

func NewThemeDefault() ThemeConfig {
	return ThemeConfig{
		Background_default:   "#090c10",
		Background_primary:   "#1a5dad",
		Background_secondary: "#3e5673",
		Text_default:         "#c9d1d9",
		Text_light:           "#8b949e",
		Text_primary:         "#73b7f2",
	}
}
