package config

func NewKeyBindingsDefault() map[string]string {
    return map[string]string{
        "<Up>":              "scrollup",
        "k":                 "scrollup",
        "<Down>":            "scrolldown",
        "j":                 "scrolldown",
        "qq":                "quit",
        "qa":                "quitall",
        "pathline:<Esc>":    "normalmod",
        "pager:<Esc>":       "normalmod",
        "pager:q":           "normalmod",
        "commandline:<Esc>": "normalmod",
        "commandline:<Enter>": "cmdrun",
        "filelist:<Esc>":    "unselectall",
        "<Enter>":           "opencurrent",
        "<Backspace2>":      "openparent",
        "h":                 "openprevious",
        "H":                 "opennext",
        "l":                 "opencurrent",
        ":":                 "opencommandline",
    }
}

func NewMouseBindingsDefault() map[string]string {
	return map[string]string{
		"filelist:<MouseScrollUp>":        "scrollup",
		"filelist:<MouseScrollDown>":      "scrolldown",
		"filelist:<MouseLeftClick>":       "setcurrent",
		"pathline:<MouseLeftClick>":       "editpath",
		"pathline:<MouseMiddleClick>":     "openpath",
		"filelist:<MouseLeftDoubleClick>": "setcurrent opencurrent",
	}
}
