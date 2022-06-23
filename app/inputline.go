package app

import (
	"atfm/app/config"
	"atfm/app/style"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type InputLineSource string

const (
	COMMAND_LINE InputLineSource = "commandline"
	SEARCH_LINE  InputLineSource = "searchline"
)

type InputLine struct {
	*tview.InputField

	appConfig config.Config

	inputHandler *InputHandler

	getInstance func() *Instance

	source InputLineSource
}

func NewInputLine(inputHandler *InputHandler, getInstance func() *Instance, appConfig config.Config) *InputLine {
	inputField := tview.NewInputField()
	ifst := tcell.StyleDefault
	ifst = ifst.Background(style.GetColorWeb(appConfig.Display.Theme.Background_default))
	ifst = ifst.Foreground(style.GetColorWeb(appConfig.Display.Theme.Text_default))
	ifst = ifst.Italic(true)
	inputField.SetFieldStyle(ifst)
	cl := InputLine{
		InputField:   inputField,
		appConfig:    appConfig,
		inputHandler: inputHandler,
		source:       COMMAND_LINE,
		getInstance:  getInstance,
	}
	cl.SetBlurFunc(func() {
		cl.SetText("")
		cl.SetLabel("")
	})
	return &cl
}

func (m *InputLine) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return m.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		if event.Key() == tcell.KeyRune ||
			!m.inputHandler.listenInputKey(event, string(m.source), true) {
			if handler := m.InputField.InputHandler(); handler != nil {
				handler(event, setFocus)
			}
		}
	})
}

func (m *InputLine) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (bool, tview.Primitive) {
	return m.WrapMouseHandler(func(action tview.MouseAction, event *tcell.EventMouse, _ func(p tview.Primitive)) (bool, tview.Primitive) {
		if !m.InRect(event.Position()) {
			return false, nil
		}
		return m.inputHandler.listenInputMouse(event, action, string(m.source)), m.InputField
	})
}

func (m *InputLine) OpenCommandLine() {
	m.source = COMMAND_LINE
	m.SetLabel(":")
}

func (m *InputLine) OpenSearchLine() {
	m.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if m.appConfig.IncSearch && event.Key() == tcell.KeyRune {
			ins := m.getInstance()
			st := m.GetText() + string(event.Rune())
			ins.QuickSearch.SearchContent(st, ins, false, m.appConfig.SearchIgnCase)
		}
		return event
	})
	m.source = SEARCH_LINE
	m.SetLabel("/")
}
