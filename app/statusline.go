package app

import (
	"atfm/app/config"
	"atfm/app/style"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type StatusLine struct {
	*tview.Box

	displayConfig config.DisplayConfig

	pane UiPane

	getInstancePane   func(UiPane) *Instance
	getInstanceGlobal func() *Instance

	inputHandler *InputHandler

	isGlobal bool
}


func NewStatusline(pane UiPane, getInstancePane func(UiPane) *Instance, getInstanceGlobal func() *Instance, inputHander *InputHandler, displayConfig config.DisplayConfig) *StatusLine {
	b := tview.NewBox().SetBackgroundColor(style.GetColorWeb(displayConfig.Theme.Background_default))
	return &StatusLine{
		Box:               b,
		displayConfig:     displayConfig,
		pane:              pane,
		getInstancePane:   getInstancePane,
		getInstanceGlobal: getInstanceGlobal,
		inputHandler:      inputHander,
		isGlobal:          false,
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

func (m *StatusLine) Draw(screen tcell.Screen) {
    for _, element := range m.displayConfig.StatusLineElements {
        align := tview.AlignLeft
        t := ""

        switch element.Alignment {
        case style.ALIGN_CENTER:
        align = tview.AlignCenter
        case style.ALIGN_RIGHT:
        align = tview.AlignRight
        case style.ALIGN_LEFT:
        align = tview.AlignLeft
        }

        switch strings.ToUpper(element.Name) {
        case "INDEX":
            t = m.displayIndex(element.Style)
        case "FILENAME":
            t = m.displayCurrentFileName(element.Style)
        case "FILEINFO":
            t = m.displayCurrentFileInfo(element.Style)
        default:
            t = element.Name
        }

        x, y, w, _ := m.GetRect()
        tview.Print(screen, t, x, y, w, align, tcell.ColorDefault)
    }
}

func (m *StatusLine) getInstance() *Instance {
	if m.isGlobal {
		return m.getInstanceGlobal()
	}
	return m.getInstancePane(m.pane)
}

func (m *StatusLine) displayIndex(style style.Style) string {
    t :=  ""
    ins := m.getInstance()
    if ins.IsEmpty() {
        t = "-/-"
        return style.Render(t)
    }
    t = strconv.FormatInt(int64(ins.CurrentItem) + 1, 10) + "/" + strconv.FormatInt(int64(len(ins.ShownContent)), 10)
    return style.Render(t)
}

func (m *StatusLine) displayCurrentFileName(style style.Style) string {
    t :=  ""
    ins := m.getInstance()
    if ins.IsEmpty() {
        t = "-"
        return style.Render(t)
    }
    t = ins.ShownContent[ins.CurrentItem].Name
    return style.Render(t)
}

func (m *StatusLine) displayCurrentFileInfo(style style.Style) string {
    t :=  ""
    ins := m.getInstance()
    if ins.IsEmpty() {
        return style.Render(t)
    }
    t = RenderFileInfo(ins.ShownContent[ins.CurrentItem],
        config.NewConfigDefault().Display.FileInfoExtendedFormat,
        m.displayConfig)
    return style.Render(t)
}
