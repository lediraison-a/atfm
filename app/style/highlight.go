package style

type Highlight struct {
	background, foreground  string
	bold, underline, italic bool
}

func NewHighlight() *Highlight {
	return &Highlight{
		background: "-",
		foreground: "-",
		bold:       false,
		underline:  false,
		italic:     false,
	}
}

func (h *Highlight) Background(color string) *Highlight {
	h.background = color
	return h
}

func (h *Highlight) Foreground(color string) *Highlight {
	h.foreground = color
	return h
}

func (h *Highlight) Bold(is bool) *Highlight {
	h.bold = is
	return h
}

func (h *Highlight) Underline(is bool) *Highlight {
	h.underline = is
	return h
}

func (h *Highlight) Italic(is bool) *Highlight {
	h.italic = is
	return h
}
