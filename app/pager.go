package app

import (
	"atfm/app/config"
	"atfm/app/style"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Pager struct {
	*tview.TextView

	inputHandler *InputHandler

	displayConfig config.DisplayConfig
}

func NewPager(inputHandler *InputHandler, displayConfig config.DisplayConfig) *Pager {
	textView := tview.NewTextView()
	textView.SetBackgroundColor(style.GetColorWeb(displayConfig.Theme.Background_light))
	textView.SetTextColor(style.GetColorWeb(displayConfig.Theme.Text_default))
	textView.SetDynamicColors(true)
	p := Pager{
		TextView:      textView,
		inputHandler:  inputHandler,
		displayConfig: displayConfig,
	}
	p.SetBlurFunc(func() {
		p.Clear()
	})
	p.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		tview.Print(screen, "(Escape or q to exit)", x, y+height-2, width, tview.AlignLeft, tcell.ColorDefault)
		return x, y, width, height - 2
	})
	return &p
}

func (m *Pager) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return m.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		if !m.inputHandler.listenInputKey(event, "pager", true) {
			if handler := m.TextView.InputHandler(); handler != nil {
				handler(event, setFocus)
			}
		}
	})
}

func (m *Pager) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (bool, tview.Primitive) {
	return m.WrapMouseHandler(func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (bool, tview.Primitive) {
		if !m.InRect(event.Position()) {
			return false, nil
		}
		if !m.inputHandler.listenInputMouse(event, action, "pager") {
			if handler := m.TextView.MouseHandler(); handler != nil {
				return handler(action, event, setFocus)
			}
		}
		return m.inputHandler.listenInputMouse(event, action, "pager"), m.TextView
	})
}
