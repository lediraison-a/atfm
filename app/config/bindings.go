package config

func NewKeyBindingsDefault() map[string]string {
	return map[string]string{
		"<Up>":           "scrollup",
		"k":              "scrollup",
		"<Down>":         "scrolldown",
		"j":              "scrolldown",
		"qq":             "quit",
		"qa":             "quitall",
		"pathline:<Esc>": "normalmod",
		"<Enter>":        "opencurrent",
	}
}

func NewMouseBindingsDefault() map[string]string {
	return map[string]string{
		"filelist:<MouseScrollUp>":   "scrollup",
		"filelist:<MouseScrollDown>": "scrolldown",
		"pathline:<MouseLeftClick>":  "editpath",
	}
}
