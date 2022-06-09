package app

import (
	"atfm/app/config"
	"atfm/app/models"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type UiPane int

const (
	LEFT UiPane = iota
	RIGHT
)

type Tui struct {
	instances *InstancePool

	appConfig *config.Config

	inputHandler *InputHandler

	app    *tview.Application
	grid   *tview.Grid
	layers *tview.Pages
	screen *tcell.Screen

	filelists []*Filelist
	tablines  []*Tabline
	pathlines []*Pathline

	showDoublePane bool

	selectedPane UiPane
}

func NewTui(instances *InstancePool, appConfig *config.Config) *Tui {
	SetAppColors(appConfig.Display.Theme)
	s, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}
	if err := s.Init(); err != nil {
		panic(err)
	}
	app := tview.NewApplication().SetScreen(s)
	appGrid := tview.NewGrid()
	tui := Tui{
		instances:      instances,
		appConfig:      appConfig,
		inputHandler:   NewInputHandler(appConfig.KeyBindings, appConfig.MouseBindings),
		app:            app,
		grid:           appGrid,
		layers:         tview.NewPages(),
		screen:         &s,
		filelists:      []*Filelist{},
		tablines:       []*Tabline{},
		pathlines:      []*Pathline{},
		showDoublePane: false,
		selectedPane:   LEFT,
	}
	app.EnableMouse(appConfig.EnableMouse)

	fll := NewFileList(LEFT, tui.getInstancePane, tui.inputHandler, appConfig.Display)
	fll.SetFocusFunc(func() {
		tui.selectedPane = LEFT
	})
	flr := NewFileList(RIGHT, tui.getInstancePane, tui.inputHandler, appConfig.Display)
	flr.SetFocusFunc(func() {
		tui.selectedPane = RIGHT
	})
	tui.filelists = append(tui.filelists, fll, flr)

	tll := NewTabline(LEFT, tui.getInstanceIndex, tui.inputHandler, appConfig.Display)
	tll.SetFocusFunc(func() {
		tui.selectedPane = LEFT
	})
	tlr := NewTabline(RIGHT, tui.getInstanceIndex, tui.inputHandler, appConfig.Display)
	tlr.SetFocusFunc(func() {
		tui.selectedPane = RIGHT
	})
	tui.tablines = append(tui.tablines, tll, tlr)

	pll := NewPathline(LEFT, tui.getInstancePane, tui.inputHandler, appConfig.Display)
	//tll.SetFocusFunc(func() {
	//	tui.selectedPane = LEFT
	//})
	plr := NewPathline(RIGHT, tui.getInstancePane, tui.inputHandler, appConfig.Display)
	//tlr.SetFocusFunc(func() {
	//	tui.selectedPane = RIGHT
	//})
	tui.pathlines = append(tui.pathlines, pll, plr)

	tui.setAppGridSinglePane()
	tui.layers.AddPage("main", tui.grid, true, true)

	commands := tui.GetAppCommands()
	tui.inputHandler.RegisterKeyActions(tui.GetActionsKey(commands)...)
	tui.inputHandler.RegisterMouseActions(tui.GetActionsMouse(commands)...)

	return &tui
}

func (t *Tui) NewInstance(openPath, basePath string, mod models.FsMod, setCurrent bool) {
	_, insId, err := t.instances.AddInstance(openPath, basePath, mod)
	if err != nil {
		return
	}
	tabline := t.tablines[t.selectedPane]
	tabline.AddTab(insId)
	if setCurrent {
		tabline.SelectedTab = insId
	}
}

func (t *Tui) StartApp() {
	if err := t.app.SetRoot(t.layers, true).SetFocus(t.layers).Run(); err != nil {
		panic(err)
	}
}

func (t *Tui) getInstanceIndex(index int) *Instance {
	return t.instances.GetInstance(index)
}

func (t *Tui) getInstancePane(view UiPane) *Instance {
	tl := t.tablines[view]
	i := tl.Tabs[tl.SelectedTab]
	return t.instances.GetInstance(i)
}

func (t *Tui) ToggleDoublePane() {
	// t.grid.Clear()
	if t.showDoublePane {
		t.tablines[LEFT].AddTabs(t.tablines[RIGHT].Tabs)
		t.tablines[RIGHT].clearTabs()
		t.setAppGridSinglePane()
	} else {
		t.NewInstance(t.appConfig.Start.StartDir, t.appConfig.Start.StartBasepath, models.LOCALFM, true)
		t.setAppGridDoublePane()
	}
	t.showDoublePane = !t.showDoublePane
}

func (t *Tui) setAppGridDoublePane() {
	t.grid.SetRows(1, 1, 0)
	t.grid.SetColumns(0, 1, 0)

	t.grid.AddItem(t.pathlines[LEFT], 1, 0, 1, 1, 1, 1, false)
	t.grid.AddItem(t.tablines[LEFT], 0, 0, 1, 1, 1, 1, false)
	t.grid.AddItem(t.filelists[LEFT], 2, 0, 1, 1, 1, 1, true)

	t.grid.AddItem(t.pathlines[RIGHT], 1, 1, 1, 1, 1, 1, false)
	t.grid.AddItem(t.tablines[RIGHT], 0, 1, 1, 1, 1, 1, false)
	t.grid.AddItem(t.filelists[RIGHT], 2, 1, 1, 1, 1, 1, false)
}

func (t *Tui) setAppGridSinglePane() {
	t.grid.SetRows(1, 1, 0)
	t.grid.SetColumns(0)
	t.grid.AddItem(t.filelists[LEFT], 2, 0, 1, 1, 1, 1, true)
	t.grid.AddItem(t.tablines[LEFT], 0, 0, 1, 1, 1, 1, false)
	t.grid.AddItem(t.pathlines[LEFT], 1, 0, 1, 1, 1, 1, false)
	// This next line is added to fix a very wierd bug :
	// Last element added to the grid don't receive click event, so I added a useless item
	t.grid.AddItem(tview.NewBox(), 0, 0, 0, 0, 1, 1, false)
}

func SetAppColors(t config.ThemeConfig) {
	COLOR_BACKGROUND := t.Background_default
	COLOR_CONTRAST_BACKGROUND := t.Background_default
	COLOR_CONTRAST_BACKGROUND_PLUS := t.Background_default
	COLOR_BORDER := t.Text_light
	COLOR_TITLE := t.Text_default
	COLOR_GRAPHICS := t.Text_default
	COLOR_TEXT_PRIMARY := t.Background_primary
	COLOR_TEXT_SECONDARY := t.Text_light
	COLOR_TEXT_TERTIARY := t.Text_default
	COLOR_TEXT_INVERSE := t.Text_default
	COLOR_TEXT_SECONDARY_CONTRAST := t.Background_default

	tview.Styles = tview.Theme{
		PrimitiveBackgroundColor: GetColorWeb(COLOR_BACKGROUND),
		// Main background color for primitives.
		ContrastBackgroundColor: GetColorWeb(COLOR_CONTRAST_BACKGROUND),
		// Background color for contrasting elements.
		MoreContrastBackgroundColor: GetColorWeb(COLOR_CONTRAST_BACKGROUND_PLUS),
		// Background color for even more contrasting elements.
		BorderColor: GetColorWeb(COLOR_BORDER),
		// Box borders.
		TitleColor: GetColorWeb(COLOR_TITLE),
		// Box titles.
		GraphicsColor: GetColorWeb(COLOR_GRAPHICS),
		// Graphics.
		PrimaryTextColor: GetColorWeb(COLOR_TEXT_PRIMARY),
		// Primary text.
		SecondaryTextColor: GetColorWeb(COLOR_TEXT_SECONDARY),
		// Secondary text (e.g. labels).
		TertiaryTextColor: GetColorWeb(COLOR_TEXT_TERTIARY),
		// Tertiary text (e.g. subtitles, notes).
		InverseTextColor: GetColorWeb(COLOR_TEXT_INVERSE),
		// Text on primary-colored backgrounds.
		ContrastSecondaryTextColor: GetColorWeb(COLOR_TEXT_SECONDARY_CONTRAST),
		// Secondary text on ContrastBackgroundColor-colored backgrounds.
	}
}
