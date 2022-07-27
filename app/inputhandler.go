package app

import (
	"atfm/generics"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var mouseEventNames = []string{
	"MouseUnknown",
	"MouseLeft",
	"MouseRight",
	"MouseMiddle",
	"MouseRelease",
	"MouseWheelUp",
	"MouseWheelDown",
	"MouseMotion",
}

var MouseActionNames = map[tview.MouseAction]string{
	tview.MouseMove:              "MouseMove",
	tview.MouseLeftDown:          "MouseLeftDown",
	tview.MouseLeftUp:            "MouseLeftUp",
	tview.MouseLeftClick:         "MouseLeftClick",
	tview.MouseLeftDoubleClick:   "MouseLeftDoubleClick",
	tview.MouseMiddleDown:        "MouseMiddleDown",
	tview.MouseMiddleUp:          "MouseMiddleUp",
	tview.MouseMiddleClick:       "MouseMiddleClick",
	tview.MouseMiddleDoubleClick: "MouseMiddleDoubleClick",
	tview.MouseRightDown:         "MouseRightDown",
	tview.MouseRightUp:           "MouseRightUp",
	tview.MouseRightClick:        "MouseRightClick",
	tview.MouseRightDoubleClick:  "MouseRightDoubleClick",
	tview.MouseScrollUp:          "MouseScrollUp",
	tview.MouseScrollDown:        "MouseScrollDown",
	tview.MouseScrollLeft:        "MouseScrollLeft",
	tview.MouseScrollRight:       "MouseScrollRight",
}

var ModifiersNames = map[tcell.ModMask]string{
	tcell.ModCtrl: "C-",
	tcell.ModAlt:  "A-",
	tcell.ModMeta: "M-",
	tcell.ModNone: "",
}

type InputHandler struct {
	listening bool
	cancel    context.CancelFunc
	ctx       context.Context

	inputPatternKey string

	keyActions   []*KeyAction
	mouseActions []*MouseAction

	KeyBindings   map[string]string
	MouseBindings map[string]string

	InputTimeout int
}

func NewInputHandler(keyBindings, mouseBindings map[string]string) *InputHandler {
	il := InputHandler{
		listening:       false,
		inputPatternKey: "",
		keyActions:      []*KeyAction{},
		mouseActions:    []*MouseAction{},
		KeyBindings:     keyBindings,
		MouseBindings:   mouseBindings,
		InputTimeout:    500,
	}
	return &il
}

type KeyAction struct {
	Name, Source string
	Action       func()
}

type MouseAction struct {
	Name, Source string
	Action       func(int, int)
}

func (il *InputHandler) RegisterKeyActions(actions ...*KeyAction) {
	il.keyActions = append(il.keyActions, actions...)
}

func (il *InputHandler) RegisterMouseActions(actions ...*MouseAction) {
	il.mouseActions = append(il.mouseActions, actions...)
}

func (il *InputHandler) startListen() {
	il.listening = true
	dur := time.Millisecond * time.Duration(il.InputTimeout)
	ctx, cancel := context.WithTimeout(context.Background(), dur)
	il.cancel = cancel
	il.ctx = ctx
	go il.endContext()
}

func (il *InputHandler) endContext() {
	defer il.cancel()
	select {
	case <-il.ctx.Done():
		il.inputPatternKey = ""
		il.listening = false
		return
	}
}

func (il *InputHandler) listenInputKey(event *tcell.EventKey, source string, ignoreGlobals bool) bool {
	key := event.Key()
	rune := event.Rune()
	p := ""
	if rune == ' ' {
		p += fmt.Sprintf("<%s%s>", ModifiersNames[event.Modifiers()], "Space")
	} else if key == tcell.KeyRune {
		p += fmt.Sprintf("%s%s", ModifiersNames[event.Modifiers()], string(rune))
	} else {
		p += fmt.Sprintf("<%s%s>", ModifiersNames[event.Modifiers()], tcell.KeyNames[event.Key()])
	}

	if !il.listening {
		il.startListen()
	}
	if il.tryKeyInput(il.inputPatternKey+p, source, ignoreGlobals) {
		il.cancel()
		return true
	} else {
		il.inputPatternKey += p
		return false
	}
}

func (il *InputHandler) tryKeyInput(pattern string, source string, ignoreGlobals bool) bool {
	k := source + ":" + pattern
	ac, ok := il.KeyBindings[k]
	if !ok && !ignoreGlobals {
		ac, ok = il.KeyBindings[pattern]
	}

	if ok {
		cs := strings.Split(ac, " ")
		for _, v := range cs {
			css := generics.Filter(il.keyActions, func(value *KeyAction, index int) bool {
				return (value.Name == v && value.Source == "") ||
					(value.Name == v && value.Source == source)
			})
			if len(css) > 0 {
				css[0].Action()
			}
		}
		return true
	}
	return false
}

func (il *InputHandler) listenInputMouse(event *tcell.EventMouse, action tview.MouseAction, source string) bool {
	actions := []tview.MouseAction{
		tview.MouseLeftClick,
		tview.MouseRightClick,
		tview.MouseMiddleClick,
		tview.MouseScrollUp,
		tview.MouseScrollDown,
		tview.MouseLeftDoubleClick,
		tview.MouseRightDoubleClick,
		tview.MouseMiddleDoubleClick,
		tview.MouseScrollLeft,
		tview.MouseScrollRight,
	}
	if !generics.Contains(actions, action) {
		return false
	}
	posX, posY := event.Position()
	p := fmt.Sprintf("<%s%s>", ModifiersNames[event.Modifiers()], MouseActionNames[action])
	return il.tryMouseInput(p, source, posX, posY)
}

func (il *InputHandler) tryMouseInput(pattern string, source string, posX, posY int) bool {
	k := source + ":" + pattern
	ac, ok := il.MouseBindings[k]
	if ok {
		cs := strings.Split(ac, " ")
		for _, v := range cs {
			css := generics.Filter(il.mouseActions, func(value *MouseAction, index int) bool {
				return (value.Name == v && value.Source == "") ||
					(value.Name == v && value.Source == source)
			})
			if len(css) > 0 {
				css[0].Action(posX, posY)
			}
		}
		return true
	}
	return false
}
