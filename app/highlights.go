package app

type Highlight struct {
	Background, Foreground, Decoration string
}

type Highlights map[string]*Highlight

func NewHighlights() Highlights {
	return map[string]*Highlight{}
}

func (h *Highlights) AddHighlight(name, bg, fg, decor string) {
	(*h)[name] = &Highlight{
		Background: bg,
		Foreground: fg,
		Decoration: decor,
	}
}
