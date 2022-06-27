package app

import (
	"atfm/app/config"
	"atfm/app/models"
	"atfm/app/style"
	"path"
	"path/filepath"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Pathline struct {
	*tview.InputField

	inputHandler *InputHandler

	displayConfig config.DisplayConfig

	pane UiPane

	getInstance func(UiPane) *Instance

	editing bool
}

func NewPathline(pane UiPane, getInstancePane func(UiPane) *Instance, inputHandler *InputHandler, displayConfig config.DisplayConfig) *Pathline {
	inputField := tview.NewInputField().SetText("test")
	inputField.SetBackgroundColor(style.GetColorWeb(displayConfig.Theme.Background_light))
	fstyle := tcell.StyleDefault.
		Foreground(style.GetColorWeb(displayConfig.Theme.Text_default)).
		Background(style.GetColorWeb(displayConfig.Theme.Background_light))
	inputField.SetFieldStyle(fstyle)
	pl := Pathline{
		InputField:    inputField,
		inputHandler:  inputHandler,
		displayConfig: displayConfig,
		pane:          pane,
		getInstance:   getInstancePane,
		editing:       false,
	}
	pl.SetDoneFunc(func(_ tcell.Key) {
		pl.OpenPath()
	})
	pl.SetAutocompleteFunc(pl.autocomplete)
	pl.SetBlurFunc(pl.blur)
	return &pl
}

func (m *Pathline) Draw(screen tcell.Screen) {
	x, y, width, _ := m.GetInnerRect()
	if m.editing {
		m.InputField.Draw(screen)
	} else {
		tbp, t := m.getPathlineTexts()
		td := path.Dir(t)
		if td != "/" {
			td += "/"
		}
		ta := path.Base(t)
		if td == ta && td == "/" {
			td = ""
		}
		bpStyle := style.NewStyle().
			Background(m.displayConfig.Theme.Background_light).
			Foreground(m.displayConfig.Theme.Text_light)
		dirStyle := style.NewStyle().
			Background(m.displayConfig.Theme.Background_light).
			Foreground(m.displayConfig.Theme.Text_default)
		baseStyle := style.NewStyle().
			Background(m.displayConfig.Theme.Background_light).
			Foreground(m.displayConfig.Theme.Text_primary).
			Bold(true)
		line := bpStyle.Render(tbp) + dirStyle.Render(td) + baseStyle.Render(ta)
		fstyle := style.NewStyle().
			Background(m.displayConfig.Theme.Background_light).
			Width(width - tview.TaggedStringWidth(line))

		line += fstyle.Render("")
		tview.Print(screen, line, x, y, width, tview.AlignLeft, tcell.ColorDefault)
	}
}

func (m *Pathline) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return m.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		if m.editing &&
			(event.Key() == tcell.KeyRune ||
				!m.inputHandler.listenInputKey(event, "pathline", true)) {
			if handler := m.InputField.InputHandler(); handler != nil {
				handler(event, setFocus)
			}
		}
	})
}

func (m *Pathline) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (bool, tview.Primitive) {
	return m.WrapMouseHandler(func(action tview.MouseAction, event *tcell.EventMouse, _ func(p tview.Primitive)) (bool, tview.Primitive) {
		if !m.InRect(event.Position()) {
			return false, nil
		}
		return m.inputHandler.listenInputMouse(event, action, "pathline"), m.InputField
	})
}

func (m *Pathline) OpenPath() {
	ins := m.getInstance(m.pane)
	p := path.Dir(m.GetText())
	if err := ins.OpenDirSaveHistory(p, ins.BasePath, ins.Mod); err != nil {

	}
}

func (m *Pathline) OpenPathPos(posX int) {
	x := posX
	px, _, _, _ := m.GetRect()
	x -= px - 1
	if x >= 0 {
		bst, tst := m.getPathlineTexts()
		x -= tview.TaggedStringWidth(bst)
		tst2 := tview.NewTextView().SetDynamicColors(true).SetText(tst).GetText(true)
		ind := len(tst2[:x]) + strings.Index(tst2[x:], "/")
		w := tst2[:ind]
		if x == 1 {
			w = "/"
		}
		ins := m.getInstance(m.pane)
		if err := ins.OpenDirSaveHistory(w, ins.BasePath, ins.Mod); err != nil {

		}
	}
}

func (m *Pathline) CompleteToParent() {
	t := m.GetText()
	m.SetText(filepath.Dir(t))
}

func (m *Pathline) EditPath() {
	tb, t := m.getPathlineTexts()
	m.SetLabel(tb)
	m.SetText(t)
	m.editing = true
	m.Autocomplete()
}

func (m *Pathline) blur() {
	m.editing = false
	m.Autocomplete()
}

func (m *Pathline) autocomplete(currentText string) []string {
	if !m.editing {
		return []string{}
	}
	ins := m.getInstance(m.pane)
	p := filepath.Clean(currentText)
	dc, err := ins.ReadDir(p, ins.BasePath, ins.Mod)
	if err != nil {
		return []string{}
	}
	var sc []string
	for _, fi := range dc {
		if fi.IsDir {
			sc = append(sc, path.Join(p, fi.Name))
		}
	}
	return sc
}

func (m *Pathline) getPathlineTexts() (string, string) {
	ins := m.getInstance(m.pane)
	tb := ""
	if ins.Mod == models.SYSTRASH || models.IsArchive(ins.Mod) {
		tb = filepath.Base(ins.BasePath)
	}
	t := ins.DirPath
	return tb, t
}
