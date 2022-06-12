package style

import (
	"testing"
)

func TestStyle(t *testing.T) {
	style := NewStyle()
	test, render := "", ""

	// test colors
	style = style.Background("red").Foreground("black").Italic(true).Bold(true).Underline(true)
	test = "[black:red:bui]test[-:-:-]"
	render = style.Render("test")
	if render != test {
		t.Fatalf("style does not render properly : expected %s, got %s", test, render)
	}

	// test paddings
	style = style.Padding(3)
	test = "[black:red:bui]   test   [-:-:-]"
	render = style.Render("test")
	if render != test {
		t.Fatalf("style does not render properly : expected %s, got %s", test, render)
	}
	style = style.PaddingLeft(2).PaddingRight(1)
	test = "[black:red:bui]  test [-:-:-]"
	render = style.Render("test")
	if render != test {
		t.Fatalf("style does not render properly : expected %s, got %s", test, render)
	}
	style = style.PaddingRune('*')
	test = "[black:red:bui]**test*[-:-:-]"
	render = style.Render("test")
	if render != test {
		t.Fatalf("style does not render properly : expected %s, got %s", test, render)
	}
}

func TestStyleSize(t *testing.T) {
	style := NewStyle()
	test, render := "", ""

	style = style.Width(6)
	test = "[-:-:-]test  [-:-:-]"
	render = style.Render("test")
	if render != test {
		t.Fatalf("style does not render properly : expected %s, got %s", test, render)
	}
	style = style.Width(2)
	test = "[-:-:-]te[-:-:-]"
	render = style.Render("test")
	if render != test {
		t.Fatalf("style does not render properly : expected %s, got %s", test, render)
	}
	style = style.Alignment(ALIGN_RIGHT)
	style = style.Width(6)
	test = "[-:-:-]  test[-:-:-]"
	render = style.Render("test")
	if render != test {
		t.Fatalf("style does not render properly : expected %s, got %s", test, render)
	}
	style = style.Width(2)
	test = "[-:-:-]st[-:-:-]"
	render = style.Render("test")
	if render != test {
		t.Fatalf("style does not render properly : expected %s, got %s", test, render)
	}
	style = style.Alignment(ALIGN_CENTER)
	style = style.Width(8)
	test = "[-:-:-]  test  [-:-:-]"
	render = style.Render("test")
	if render != test {
		t.Fatalf("style does not render properly : expected %s, got %s", test, render)
	}
	style = style.Width(2)
	test = "[-:-:-]te[-:-:-]"
	render = style.Render("test")
	if render != test {
		t.Fatalf("style does not render properly : expected %s, got %s", test, render)
	}
	style = style.Width(8)
	test = "[-:-:-] testo  [-:-:-]"
	render = style.Render("testo")
	if render != test {
		t.Fatalf("style does not render properly : expected %s, got %s", test, render)
	}
	style = style.Width(2)
	test = "[-:-:-]te[-:-:-]"
	render = style.Render("test")
	if render != test {
		t.Fatalf("style does not render properly : expected %s, got %s", test, render)
	}
}
