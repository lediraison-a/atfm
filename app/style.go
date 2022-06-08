package app

import (
	"fmt"
	"strings"
)

type Style struct {
	background, foreground    string
	bold, underline, italic   bool
	paddingLeft, paddingRight int
}

func NewStyle() *Style {
	return &Style{
		background:   "-",
		foreground:   "-",
		bold:         false,
		underline:    false,
		italic:       false,
		paddingLeft:  0,
		paddingRight: 0,
	}
}

func (s *Style) Render(text string) string {
	decoration := ":"
	if s.bold {
		decoration += "b"
	}
	if s.underline {
		decoration += "u"
	}
	if s.italic {
		decoration += "i"
	}
	ptext := strings.Repeat(" ", s.paddingLeft) + text + strings.Repeat(" ", s.paddingRight)
	style := fmt.Sprintf("[%s:%s%s]", s.foreground, s.background, decoration)
	return style + ptext + "[-:-:-]"
}

func (s *Style) Background(color string) *Style {
	s.background = color
	return s
}

func (s *Style) Foreground(color string) *Style {
	s.foreground = color
	return s
}

func (s *Style) Bold(is bool) *Style {
	s.bold = is
	return s
}

func (s *Style) Underline(is bool) *Style {
	s.underline = is
	return s
}

func (s *Style) Italic(is bool) *Style {
	s.italic = is
	return s
}

func (s *Style) Padding(value int) *Style {
	if value < 0 {
		return s
	}
	s.paddingLeft = value
	s.paddingRight = value
	return s
}

func (s *Style) PaddingLeft(value int) *Style {
	if value < 0 {
		return s
	}
	s.paddingLeft = value
	return s
}

func (s *Style) PaddingRight(value int) *Style {
	if value < 0 {
		return s
	}
	s.paddingRight = value
	return s
}
