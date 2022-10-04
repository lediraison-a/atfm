package config

type ThemeConfig struct {
	Background_default   string
	Background_primary   string
	Background_light     string
	Background_secondary string
	Text_default         string
	Text_light           string
	Text_primary         string
	Text_warning         string
	Text_error           string
}

func NewThemeDefault() ThemeConfig {
	return ThemeConfig{
		Background_default:   "#090c10",
		Background_light:     "#0d1117",
		Background_primary:   "#1a5dad",
		Background_secondary: "#3e5673",
		Text_default:         "#c9d1d9",
		Text_light:           "#8b949e",
		Text_primary:         "#73b7f2",
		Text_warning:         "#ffa348",
		Text_error:           "#ed333b",
	}
}
