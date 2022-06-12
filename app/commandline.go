package app

import (
	"atfm/app/config"
	"atfm/app/style"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type CommandLine struct {
	*tview.InputField

	displayConfig config.DisplayConfig

	inputHandler *InputHandler

	Commands *CommandManager

	page func(string)
}

func NewCommandLine(commands *CommandManager, inputHandler *InputHandler, displayConfig config.DisplayConfig) *CommandLine {
	inputField := tview.NewInputField().SetLabel(":")
	ifst := tcell.StyleDefault
	ifst = ifst.Background(style.GetColorWeb(displayConfig.Theme.Background_default))
	ifst = ifst.Foreground(style.GetColorWeb(displayConfig.Theme.Text_default))
	ifst = ifst.Italic(true)
	inputField.SetFieldStyle(ifst)
	cl := CommandLine{
		InputField:    inputField,
		displayConfig: displayConfig,
		inputHandler:  inputHandler,
		Commands:      commands,
	}
	cl.SetBlurFunc(func() {
		cl.SetText("")
	})
	return &cl
}

func (m *CommandLine) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return m.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		if event.Key() == tcell.KeyRune ||
			!m.inputHandler.listenInputKey(event, "commandline", true) {
			if handler := m.InputField.InputHandler(); handler != nil {
				handler(event, setFocus)
			}
		}
	})
}

func (m *CommandLine) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (bool, tview.Primitive) {
	return m.WrapMouseHandler(func(action tview.MouseAction, event *tcell.EventMouse, _ func(p tview.Primitive)) (bool, tview.Primitive) {
		if !m.InRect(event.Position()) {
			return false, nil
		}
		return m.inputHandler.listenInputMouse(event, action, "commandline"), m.InputField
	})
}
