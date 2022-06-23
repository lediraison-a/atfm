package config

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
	InfoSeparator  string
}

func NewConfigDefault() *Config {
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
			Theme:          NewThemeDefault(),
			DateFormat:     "Jan _2 15:04:05",
			FileInfoFormat: []string{"~> {symlink}", "{size}", "{date}"},
			InfoSeparator:  " • ",
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
