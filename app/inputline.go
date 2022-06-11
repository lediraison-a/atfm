package app

import "github.com/rivo/tview"

type InputLine struct {
	*tview.Flex

	prompt tview.Box

	inputHandler *InputHandler

	getInstance func() *Instance
}

func (m *InputLine) AddInput() {

}
