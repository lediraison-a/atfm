package config

func NewKeyBindingsDefault() map[string]string {
	return map[string]string{
		"<Up>":           "pageup",
		"k":              "scrollup",
		"<Down>":         "pagedown",
		"j":              "scrolldown",
		"qq":             "quit",
		"qa":             "quitall",
		"filelist:<Esc>": "unselectall",
		"<Enter>":        "opencurrent",
		"<Backspace2>":   "openparent",
		"h":              "openprevious",
		"H":              "opennext",
		"l":              "opencurrent",
		":":              "cmdinput",
		"/":              "searchinput",
		"n":              "searchjumpforward",
		"N":              "searchjumpbackward",

		"pager:<Esc>": "normalmod",
		"pager:q":     "normalmod",

		"commandline:<Esc>":   "normalmod",
		"commandline:<Enter>": "cmdrun",
		"commandline:<Up>":    "cmdprevious",
		"commandline:<Down>":  "cmdnext",

		"searchline:<Esc>":   "normalmod",
		"searchline:<Enter>": "searchrun",
		"searchline:<Up>":    "searchprevious",
		"searchline:<Down>":  "searchnext",
	}
}

func NewMouseBindingsDefault() map[string]string {
	return map[string]string{
		"filelist:<MouseScrollUp>":        "scrollup",
		"filelist:<MouseScrollDown>":      "scrolldown",
		"filelist:<MouseLeftClick>":       "setcurrent",
		"tabline:<MouseLeftClick>":        "tabsetcurrent",
		"pathline:<MouseLeftClick>":       "editpath",
		"pathline:<MouseMiddleClick>":     "openpath",
		"filelist:<MouseLeftDoubleClick>": "setcurrent opencurrent",
	}
}
