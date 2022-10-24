package config

import "atfm/app/style"

type Config struct {
	// Start         StartConfig
	Display                 DisplayConfig
	Preview                 PreviewConfig
	KeyBindings             map[string]string
	MouseBindings           map[string]string
	StartDir, StartBasepath string
	Readonly                bool
	LeaderKey               string
	EnableMouse             bool
	IncSearch               bool
	SearchIgnCase           bool
}

type PreviewConfig struct {
	FilePreviewer, DirPreviewer string
	FilePrevInternal,
	DirPrevInternal,
	ArchivePrevInternal bool
	PreviewFileMaxSize        int
	CacheSize, ProcessTimeout int32
}

type StartConfig struct {
	StartDir, StartBasepath string
	Readonly                bool
}

type DisplayConfig struct {
	ShowIcons              bool
	ShowOpenParent         bool
	TabLen                 int
	DynamicTabSize         bool
	ShowTabTitle           bool
	ShowTabNumber          bool
	Theme                  ThemeConfig
	DateFormat             string
	FileInfoFormat         []string
	FileInfoExtendedFormat []string
	InfoSeparator          string
	StatusLineElements     []StatusLineElement
}

func NewConfigDefault() *Config {
	defaultTheme := NewThemeDefault()
	c := Config{

		Display: DisplayConfig{
			ShowIcons:              true,
			ShowOpenParent:         true,
			TabLen:                 14,
			DynamicTabSize:         true,
			ShowTabTitle:           true,
			ShowTabNumber:          true,
			Theme:                  defaultTheme,
			DateFormat:             "Jan _2 15:04",
			FileInfoFormat:         []string{"~> {symlink}", "{size}", "{date}"},
			FileInfoExtendedFormat: []string{"{name}", "~> {symlink}", "{mod}", "{size}", "{date}"},
			InfoSeparator:          " • ",
			StatusLineElements: []StatusLineElement{
				{
					Highlight: *style.NewHighlight().
						Background(defaultTheme.Background_primary).
						Foreground(defaultTheme.Text_default),
					Name:      "INDEX",
					Alignment: style.ALIGN_RIGHT,
				},
				{
					Highlight: *style.NewHighlight().
						Background(defaultTheme.Background_default).
						Foreground(defaultTheme.Text_light),
					Name:      "FILEINFO",
					Alignment: style.ALIGN_LEFT,
				},
			},
		},
		Preview: PreviewConfig{
			FilePreviewer:       "pistol",
			DirPreviewer:        "pistol",
			FilePrevInternal:    false,
			ArchivePrevInternal: true,
			DirPrevInternal:     true,
		},
		StartDir:      "$HOME",
		StartBasepath: "/",
		Readonly:      false,
		KeyBindings:   NewKeyBindingsDefault(),
		MouseBindings: NewMouseBindingsDefault(),
		EnableMouse:   true,
		IncSearch:     true,
		SearchIgnCase: true,
		LeaderKey:     " ",
	}
	return &c
}
