package config

import "atfm/app/style"

type Config struct {
	Start         StartConfig
	Display       DisplayConfig
	Preview       PreviewConfig
	KeyBindings   map[string]string
	MouseBindings map[string]string
	LeaderKey     string
	EnableMouse   bool
	IncSearch     bool
	SearchIgnCase bool
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
	ShowIcons      bool
	ShowOpenParent bool
	TabLen         int
	DynamicTabSize bool
	ShowTabTitle   bool
	ShowTabNumber  bool
	Theme          ThemeConfig
	DateFormat     string
	FileInfoFormat []string
	FileInfoExtendedFormat []string
	InfoSeparator  string
    StatusLineElements []StatusLineElement
}

func NewConfigDefault() *Config {
    defaultTheme := NewThemeDefault()
	c := Config{
		Start: StartConfig{
			StartDir:      "/home/alban",
			StartBasepath: "/",
			Readonly:      false,
		},
		Display: DisplayConfig{
			ShowIcons:      true,
			ShowOpenParent: true,
			TabLen:         14,
			DynamicTabSize: true,
			ShowTabTitle:   true,
			ShowTabNumber:  true,
			Theme:          defaultTheme,
			DateFormat:     "Jan _2 15:04:05",
			FileInfoFormat: []string{"~> {symlink}", "{size}", "{date}"},
			FileInfoExtendedFormat: []string{"{name}", "~> {symlink}", "{mod}", "{size}", "{date}"},
			InfoSeparator:  " • ",
            StatusLineElements: []StatusLineElement{
                {
                	Style:     *style.NewStyle().
                        Background(defaultTheme.Background_primary).
                        Foreground(defaultTheme.Text_default).
                        Padding(1),
                	Name:      "INDEX",
                	Alignment: style.ALIGN_RIGHT,
                },
                {
                	Style:     *style.NewStyle().
                        Background(defaultTheme.Background_default).
                        Foreground(defaultTheme.Text_light).
                        Padding(1),
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
		KeyBindings:   NewKeyBindingsDefault(),
		MouseBindings: NewMouseBindingsDefault(),
		EnableMouse:   true,
		IncSearch:     true,
		SearchIgnCase: true,
	}
	return &c
}
