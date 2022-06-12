package app

import (
	"atfm/app/config"
	"atfm/app/style"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type StatusLine struct {
    *tview.Box

    displayConfig config.DisplayConfig

    pane UiPane

    getInstance func(UiPane) *Instance

    inputHandler *InputHandler
}

func NewStatusline(pane UiPane, getInstance func(UiPane) *Instance, inputHander *InputHandler, displayConfig config.DisplayConfig) *StatusLine {
	b := tview.NewBox().SetBackgroundColor(style.GetColorWeb(displayConfig.Theme.Background_default))
    return &StatusLine{
    	Box:           b,
    	displayConfig: displayConfig,
    	pane:          pane,
    	getInstance:   getInstance,
    	inputHandler:  inputHander,
    }
}

func (m *StatusLine) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return m.WrapInputHandler(func(event *tcell.EventKey, _ func(p tview.Primitive)) {
		m.inputHandler.listenInputKey(event, "statusline", false)
	})
}

func (m *StatusLine) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (bool, tview.Primitive) {
	return m.WrapMouseHandler(func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (bool, tview.Primitive) {
		if !m.InRect(event.Position()) {
			return false, nil
		}
		return m.inputHandler.listenInputMouse(event, action, "statusline"), m.Box
	})
}
