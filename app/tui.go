package app

import (
	"atfm/app/config"
	"atfm/app/models"
	"atfm/app/style"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
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

	cmdManager *CommandManager

	app    *tview.Application
	grid   *tview.Grid
	layers *tview.Pages
	screen *tcell.Screen

	filelists   []*Filelist
	tablines    []*Tabline
	pathlines   []*Pathline
	statuslines []*StatusLine

	pager *Pager

	inputLine *InputLine

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

	commandManager := NewCommandManager()
	inputHandler := NewInputHandler(appConfig.KeyBindings, appConfig.MouseBindings)

	tui := Tui{
		instances:      instances,
		appConfig:      appConfig,
		inputHandler:   inputHandler,
		cmdManager:     commandManager,
		app:            app,
		grid:           appGrid,
		layers:         tview.NewPages(),
		screen:         &s,
		filelists:      []*Filelist{},
		tablines:       []*Tabline{},
		pathlines:      []*Pathline{},
		statuslines:    []*StatusLine{},
		pager:          NewPager(inputHandler, appConfig.Display),
		inputLine:      &InputLine{},
		showDoublePane: false,
		selectedPane:   LEFT,
	}
	app.EnableMouse(appConfig.EnableMouse)

	tui.inputLine = NewInputLine(inputHandler, func() *Instance {
		return tui.getInstancePane(tui.selectedPane)
	}, *appConfig)

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
	tll.SetFocusFunc(func() {
		tui.selectedPane = LEFT
	})
	plr := NewPathline(RIGHT, tui.getInstancePane, tui.inputHandler, appConfig.Display)
	tlr.SetFocusFunc(func() {
		tui.selectedPane = RIGHT
	})
	tui.pathlines = append(tui.pathlines, pll, plr)

	sll := NewStatusline(LEFT, tui.getInstancePane, tui.getInstanceGlobal, tui.inputHandler, appConfig.Display)
	tll.SetFocusFunc(func() {
		tui.selectedPane = LEFT
	})
	slr := NewStatusline(RIGHT, tui.getInstancePane, tui.getInstanceGlobal, tui.inputHandler, appConfig.Display)
	tlr.SetFocusFunc(func() {
		tui.selectedPane = RIGHT
	})
	tui.statuslines = append(tui.statuslines, sll, slr)

	tui.setAppGridSinglePane()
	tui.layers.AddPage("main", tui.grid, true, true)

	tui.layers.AddPage("pager", tui.pager, false, false)

	commands := tui.GetAppCommands()
	commandManager.RootCmd.AddCommand(commands...)
	tui.inputHandler.RegisterKeyActions(tui.GetActionsKey(commands)...)
	tui.inputHandler.RegisterMouseActions(tui.GetActionsMouse(commands)...)

	return &tui
}

func (t *Tui) NewInstance(openPath, basePath string, mod models.FsMod, setCurrent bool) error {
	_, insId, err := t.instances.AddInstance(openPath, basePath, mod)
	if err != nil {
		return err
	}
	tabline := t.tablines[t.selectedPane]
	tid := tabline.AddTab(insId)
	if setCurrent {
		tabline.SelectedTab = tid
	}
	return nil
}

func (t *Tui) RefreshInstances(path string, content []models.FileInfo, selfDelete bool) {
	t.app.QueueUpdateDraw(func() {
		t.instances.RefreshInstances(path, content, selfDelete)
	})
}

func (t *Tui) StartApp() {
	if err := t.app.SetRoot(t.layers, true).SetFocus(t.layers).Run(); err != nil {
		panic(err)
	}
}

func (t *Tui) getInstanceIndex(index int) *Instance {
	return t.instances.GetInstance(index)
}

func (t *Tui) getInstanceGlobal() *Instance {
	return t.instances.instances[t.selectedPane]
}

func (t *Tui) getInstancePane(view UiPane) *Instance {
	tl := t.tablines[view]
	i := tl.Tabs[tl.SelectedTab]
	return t.instances.GetInstance(i)
}

func (t *Tui) ToggleDoublePane() {
	// t.grid.Clear()
	if t.showDoublePane {
		t.selectedPane = LEFT
		t.tablines[LEFT].AddTabs(t.tablines[RIGHT].Tabs)
		t.tablines[RIGHT].clearTabs()
		t.setAppGridSinglePane()
	} else {
		t.selectedPane = RIGHT
		t.NewInstance(t.appConfig.StartDir, t.appConfig.StartBasepath, models.LOCALFM, true)
		t.setAppGridDoublePane()
	}
	t.showDoublePane = !t.showDoublePane
}

func (t *Tui) setAppGridDoublePane() {
	t.grid.SetRows(1, 1, 0, 1, 1)
	t.grid.SetColumns(0, 1, 0)
	t.grid.AddItem(t.filelists[LEFT], 2, 0, 1, 1, 1, 1, true)
	t.grid.AddItem(t.tablines[LEFT], 0, 0, 1, 1, 1, 1, false)
	t.grid.AddItem(t.pathlines[LEFT], 1, 0, 1, 1, 1, 1, false)
	t.grid.AddItem(t.statuslines[LEFT], 3, 0, 1, 1, 1, 1, false)

	t.grid.AddItem(t.filelists[RIGHT], 2, 2, 1, 1, 1, 1, true)
	t.grid.AddItem(t.tablines[RIGHT], 0, 2, 1, 1, 1, 1, false)
	t.grid.AddItem(t.pathlines[RIGHT], 1, 2, 1, 1, 1, 1, false)
	t.grid.AddItem(t.statuslines[RIGHT], 3, 2, 1, 1, 1, 1, false)

	t.grid.AddItem(t.inputLine, 4, 0, 1, 3, 1, 1, false)
	// This next line is added to fix a very wierd bug :
	// Last element added to the grid don't receive click event, so I added a useless item
	t.grid.AddItem(tview.NewBox(), 0, 0, 0, 0, 1, 1, false)
}

func (t *Tui) setAppGridSinglePane() {
	t.grid.SetRows(1, 1, 0, 1, 1)
	t.grid.SetColumns(0)
	t.grid.AddItem(t.filelists[LEFT], 2, 0, 1, 1, 1, 1, true)
	t.grid.AddItem(t.tablines[LEFT], 0, 0, 1, 1, 1, 1, false)
	t.grid.AddItem(t.pathlines[LEFT], 1, 0, 1, 1, 1, 1, false)
	t.grid.AddItem(t.statuslines[LEFT], 3, 0, 1, 1, 1, 1, false)

	t.grid.AddItem(t.inputLine, 4, 0, 1, 1, 1, 1, false)
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
		PrimitiveBackgroundColor: style.GetColorWeb(COLOR_BACKGROUND),
		// Main background color for primitives.
		ContrastBackgroundColor: style.GetColorWeb(COLOR_CONTRAST_BACKGROUND),
		// Background color for contrasting elements.
		MoreContrastBackgroundColor: style.GetColorWeb(COLOR_CONTRAST_BACKGROUND_PLUS),
		// Background color for even more contrasting elements.
		BorderColor: style.GetColorWeb(COLOR_BORDER),
		// Box borders.
		TitleColor: style.GetColorWeb(COLOR_TITLE),
		// Box titles.
		GraphicsColor: style.GetColorWeb(COLOR_GRAPHICS),
		// Graphics.
		PrimaryTextColor: style.GetColorWeb(COLOR_TEXT_PRIMARY),
		// Primary text.
		SecondaryTextColor: style.GetColorWeb(COLOR_TEXT_SECONDARY),
		// Secondary text (e.g. labels).
		TertiaryTextColor: style.GetColorWeb(COLOR_TEXT_TERTIARY),
		// Tertiary text (e.g. subtitles, notes).
		InverseTextColor: style.GetColorWeb(COLOR_TEXT_INVERSE),
		// Text on primary-colored backgrounds.
		ContrastSecondaryTextColor: style.GetColorWeb(COLOR_TEXT_SECONDARY_CONTRAST),
		// Secondary text on ContrastBackgroundColor-colored backgrounds.
	}
}

func (t *Tui) GetAppCommands() []*cobra.Command {
	return append(
		append(getCommandsMove(t),
			append(getCommandsTabs(t),
				getCommandsApp(t)...)...),
		getCommandsFile(t)...)
}
