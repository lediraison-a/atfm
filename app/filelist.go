package app

import (
	"atfm/app/config"
	"atfm/app/icons"
	"atfm/app/models"
	"atfm/generics"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Filelist struct {
	*tview.Box

	inputHandler    *InputHandler
	pane            UiPane
	GetInstancePane func(UiPane) *Instance

	offset int

	displayConfig config.DisplayConfig

	isDragSelect   bool
	lastDragSelect int
}

func NewFileList(pane UiPane, getInstancePane func(UiPane) *Instance, inputHandler *InputHandler, displayConfig config.DisplayConfig) *Filelist {
	b := tview.NewBox().SetBackgroundColor(GetColorWeb(displayConfig.Theme.Background_default))
	m := Filelist{
		Box:             b,
		inputHandler:    inputHandler,
		pane:            pane,
		GetInstancePane: getInstancePane,
		offset:          0,
		displayConfig:   displayConfig,
		isDragSelect:    false,
		lastDragSelect:  0,
	}
	return &m
}

func (m *Filelist) Draw(screen tcell.Screen) {
	x, y, width, height := m.GetInnerRect()
	ins := m.GetInstancePane(m.pane)
	current := ins.CurrentItem

	if current < m.offset {
		m.offset = current
	} else if ins.CurrentItem >= m.offset+(height-1) {
		m.offset = (current - height) + 1
	}
	if m.offset < 0 {
		m.offset = 0
	}
	itemCount := len(ins.ShownContent)

	printMainText := func(item models.FileInfo, index int) string {
		bg := m.displayConfig.Theme.Background_default
		ic, icol, t := "", "", ""
		if m.displayConfig.ShowIcons {
			if item.IsDir {
				ic = icons.GetDirIcon(item.Name)
				icol = m.displayConfig.Theme.Text_primary
			} else {
				ic, icol = icons.GetFileIcon(item.Name, m.displayConfig.Theme.Text_default)
			}
		}
		t += item.Name

		if generics.Contains(ins.SelectedIndexes, index) {
			bg = m.displayConfig.Theme.Background_secondary
		}
		if ins.CurrentItem == index {
			bg = m.displayConfig.Theme.Background_primary
		}

		iconStyle := NewStyle().
			Foreground(icol).
			Background(bg).
			Padding(1)
		lineMainTextStyle := NewStyle().
			Foreground(m.displayConfig.Theme.Text_default).
			Background(bg).
			PaddingRight(1)
		return iconStyle.Render(ic) + lineMainTextStyle.Render(t)
	}

	printInfoText := func(item models.FileInfo, _ int) string {
		infoTextStyle := NewStyle().
			Foreground(m.displayConfig.Theme.Text_light).
			Background(m.displayConfig.Theme.Background_default)
		return infoTextStyle.Render("")
	}

	for i := y; i <= height+1; i++ {
		li := i - y
		cIndex := m.offset + li
		if cIndex >= itemCount {
			tview.Print(screen, "", x, i, width, tview.AlignLeft, tcell.ColorDefault)
			continue
		}
		fi := ins.ShownContent[cIndex]
		tview.Print(screen, printInfoText(fi, cIndex), x, i, width, tview.AlignRight, tcell.ColorDefault)
		tview.Print(screen, printMainText(fi, cIndex), x, i, width, tview.AlignLeft, tcell.ColorDefault)

	}
}

func (m *Filelist) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return m.WrapInputHandler(func(event *tcell.EventKey, _ func(p tview.Primitive)) {
		m.inputHandler.listenInputKey(event, "filelist", false)
	})
}

func (m *Filelist) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (bool, tview.Primitive) {
	return m.WrapMouseHandler(func(action tview.MouseAction, event *tcell.EventMouse, _ func(p tview.Primitive)) (bool, tview.Primitive) {
		if !m.InRect(event.Position()) {
			return false, nil
		}
		return m.inputHandler.listenInputMouse(event, action, "filelist"), m.Box
	})
}

func (m *Filelist) ScrollUp() {
	m.moveCurrentItem(true)
}

func (m *Filelist) ScrollDown() {
	m.moveCurrentItem(false)
}

func (m *Filelist) ScrollFirst() {
	ins := m.GetInstancePane(m.pane)
	ins.CurrentItem = 0
}

func (m *Filelist) ScrollLast() {
	ins := m.GetInstancePane(m.pane)
	ins.CurrentItem = len(ins.ShownContent) - 1
}

func (m *Filelist) moveCurrentItem(up bool) {
	ins := m.GetInstancePane(m.pane)
	current := ins.CurrentItem
	if up {
		current--
		if current < 0 {
			current = ins.FileCount() - 1
		}
	} else {
		current++
		if current > ins.FileCount()-1 {
			current = 0
		}
	}
	ins.CurrentItem = current
}

func (m *Filelist) getUnderMouseIndex(mousePosY int) int {
	ins := m.GetInstancePane(m.pane)
	a := m.offset
	pos := 2
	ii := mousePosY - pos + a
	if ii >= ins.FileCount() {
		ii = ins.FileCount() - 1
	}
	return ii
}
