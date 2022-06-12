package app

import (
	"atfm/app/config"
	"atfm/app/icons"
	"atfm/app/models"
	"atfm/app/style"
	"atfm/generics"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Filelist struct {
	*tview.Box

	inputHandler    *InputHandler
	pane            UiPane
	GetInstancePane func(UiPane) *Instance

	listOffset int

	displayConfig config.DisplayConfig

	isDragSelect   bool
	lastDragSelect int

	DragSelectionSupport bool
}

func NewFileList(pane UiPane, getInstancePane func(UiPane) *Instance, inputHandler *InputHandler, displayConfig config.DisplayConfig) *Filelist {
	b := tview.NewBox().SetBackgroundColor(style.GetColorWeb(displayConfig.Theme.Background_default))
	m := Filelist{
		Box:                  b,
		inputHandler:         inputHandler,
		pane:                 pane,
		GetInstancePane:      getInstancePane,
		listOffset:           0,
		displayConfig:        displayConfig,
		isDragSelect:         false,
		lastDragSelect:       0,
		DragSelectionSupport: true,
	}
	return &m
}

func (m *Filelist) Draw(screen tcell.Screen) {
	x, y, width, height := m.GetInnerRect()
	ins := m.GetInstancePane(m.pane)
	current := ins.CurrentItem

	if current < m.listOffset {
		m.listOffset = current
	} else if ins.CurrentItem >= m.listOffset+(height-1) {
		m.listOffset = (current - height) + 1
	}
	if m.listOffset < 0 {
		m.listOffset = 0
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

		iconStyle := style.NewStyle().
			Foreground(icol).
			Background(bg).
			Padding(1)
		lineMainTextStyle := style.NewStyle().
			Foreground(m.displayConfig.Theme.Text_default).
			Background(bg).
			PaddingRight(1)
		return iconStyle.Render(ic) + lineMainTextStyle.Render(t)
	}

	printInfoText := func(item models.FileInfo, _ int) string {
		infoTextStyle := style.NewStyle().
			Foreground(m.displayConfig.Theme.Text_light).
			Background(m.displayConfig.Theme.Background_default)
		return infoTextStyle.Render(RenderFileInfo(item, m.displayConfig))
	}

	if itemCount == 0 {
		infoTextStyle := style.NewStyle().
			Foreground(m.displayConfig.Theme.Text_light).
			Background(m.displayConfig.Theme.Background_default)
		t := ""
		if m.displayConfig.ShowIcons {
			t = icons.DirEmptyIcon + "  Empty"
		} else {
			t = " - Empty - "
		}
		tview.Print(screen, infoTextStyle.Render(t), x, y+(height/2), width, tview.AlignCenter, tcell.ColorDefault)
		return
	}

	for i := y; i <= height+1; i++ {
		li := i - y
		cIndex := m.listOffset + li
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
	return m.WrapMouseHandler(func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (bool, tview.Primitive) {
		x, y := event.Position()
		if !m.InRect(x, y) {
			m.isDragSelect = false
			return false, nil
		}
		omh := false
		if handler := m.Box.MouseHandler(); handler != nil {
			omh, _ = handler(action, event, setFocus)
		}
		redraw := (m.DragSelectionSupport && m.HandleDragSelection(action, x, y)) ||
			m.inputHandler.listenInputMouse(event, action, "filelist") ||
			omh
		return redraw, m.Box
	})
}

func (m *Filelist) ScrollUp() {
	m.MoveCurrentItem(true)
}

func (m *Filelist) ScrollDown() {
	m.MoveCurrentItem(false)
}

func (m *Filelist) ScrollFirst() {
	ins := m.GetInstancePane(m.pane)
	ins.CurrentItem = 0
}

func (m *Filelist) ScrollLast() {
	ins := m.GetInstancePane(m.pane)
	ins.CurrentItem = len(ins.ShownContent) - 1
}

func (m *Filelist) MoveCurrentItem(up bool) {
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

func (m *Filelist) GetUnderMouseIndex(mousePosY int) int {
	ins := m.GetInstancePane(m.pane)
	a := m.listOffset
	pos := 2
	ii := mousePosY - pos + a
	if ii >= ins.FileCount() {
		ii = ins.FileCount() - 1
	}
	return ii
}

func (m *Filelist) HandleDragSelection(action tview.MouseAction, posX, posY int) bool {
	ins := m.GetInstancePane(m.pane)
	underMouseIndex := m.GetUnderMouseIndex(posY)
	switch action {

	case tview.MouseLeftDown:
		ins.UnselectAll()
		m.isDragSelect = true
		m.lastDragSelect = underMouseIndex
		return false

	case tview.MouseLeftUp:
		m.isDragSelect = false
		return false

	case tview.MouseMove:
		if !m.isDragSelect {
			break
		}
		if underMouseIndex > m.lastDragSelect {
			for i := m.lastDragSelect; i < underMouseIndex; i++ {
				ins.SelectItem(i, true)
			}
		} else if underMouseIndex < m.lastDragSelect {
			for i := m.lastDragSelect; i > underMouseIndex; i-- {
				ins.SelectItem(i, true)
			}
		} else {
			m.isDragSelect = false
			return false
		}
		ins.CurrentItem = underMouseIndex
		m.isDragSelect = true
		m.lastDragSelect = underMouseIndex
		return true
	}
	return false
}
