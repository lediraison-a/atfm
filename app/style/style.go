package style

import (
	"fmt"
	"strings"
)

type Alignment string

const (
	ALIGN_LEFT   Alignment = "LEFT"
	ALIGN_RIGHT  Alignment = "RIGHT"
	ALIGN_CENTER Alignment = "CENTER"
)

type Style struct {
	background, foreground    string
	bold, underline, italic   bool
	paddingLeft, paddingRight int
	paddingRune               rune
	minWidth, maxWidth        int
	alignment                 Alignment
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
		paddingRune:  ' ',
		minWidth:     -1,
		maxWidth:     -1,
		alignment:    ALIGN_LEFT,
	}
}

func (s *Style) Render(text string) string {
	resizeText := func(text string, size int) string {
		ttext := text
		switch s.alignment {
		case ALIGN_LEFT:
			return fmt.Sprintf("%*s", -size, text)
		case ALIGN_RIGHT:
			return fmt.Sprintf("%*s", size, text)
		case ALIGN_CENTER:
			mq := size - len(text)
			ttext = fmt.Sprintf("%*s", len(text)+(mq/2), ttext)
			ttext = fmt.Sprintf("%*s", -size, ttext)
		}
		return ttext
	}

	shrinkText := func(text string, size int) string {
		ttext := text
		switch s.alignment {
		case ALIGN_LEFT:
			return ttext[:size]
		case ALIGN_RIGHT:
			return ttext[len(text)-size:]
		case ALIGN_CENTER:
			return ttext[:size]
			// return ttext[(len(text)/2)-(size/2) : (len(text)/2)+(size/2)]
		}
		return ttext
	}

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
	if !s.bold && !s.italic && !s.underline {
		decoration += "-"
	}

	if s.minWidth != -1 && len(text) < (s.minWidth-(s.paddingRight+s.paddingLeft)) {
		text = resizeText(text, s.minWidth)
	} else if s.maxWidth != -1 && len(text) > (s.maxWidth-(s.paddingRight+s.paddingLeft)) {
		text = shrinkText(text, s.maxWidth)
	}

	pc := string(s.paddingRune)
	ptext := strings.Repeat(pc, s.paddingLeft) + StripAllTag(text) + strings.Repeat(pc, s.paddingRight)

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

func (s *Style) PaddingRune(value rune) *Style {
	s.paddingRune = value
	return s
}

func (s *Style) Alignment(value Alignment) *Style {
	s.alignment = value
	return s
}

func (s *Style) Width(value int) *Style {
	if value < 0 {
		value = -1
	}
	s.minWidth = value
	s.maxWidth = value
	return s
}

func (s *Style) MinWidth(value int) *Style {
	if value < 0 {
		value = -1
	}
	s.minWidth = value
	return s
}

func (s *Style) MaxWidth(value int) *Style {
	if value < 0 {
		value = -1
	}
	s.maxWidth = value
	return s
}
