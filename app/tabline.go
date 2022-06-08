package app

import (
	"atfm/app/config"
	"fmt"
	"path"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Tabline struct {
	*tview.Box

	inputHandler     *InputHandler
	pane             UiPane
	GetInstanceIndex func(int) *Instance

	SelectedTab int
	Tabs        []int

	displayConfig config.DisplayConfig

	showCloseButton bool

	canCloseTab bool

	tabShowOffset int
	tablen        int
}

func NewTabline(pane UiPane, getInstanceIndex func(int) *Instance, inputHandler *InputHandler, displayConfig config.DisplayConfig) *Tabline {
	b := tview.NewBox().SetBackgroundColor(GetColorWeb(displayConfig.Theme.Background_default))
	tablen := 16
	tl := Tabline{
		Box:              b,
		inputHandler:     inputHandler,
		pane:             pane,
		GetInstanceIndex: getInstanceIndex,
		SelectedTab:      0,
		Tabs:             []int{},
		showCloseButton:  false,
		canCloseTab:      false,
		tabShowOffset:    0,
		tablen:           tablen,
		displayConfig:    displayConfig,
	}
	return &tl
}

func (m *Tabline) Draw(screen tcell.Screen) {
	x, y, width, _ := m.GetInnerRect()

	tablen := m.getTabSizeCompute()

	nbTabWidth := width / tablen
	if m.SelectedTab >= m.tabShowOffset+nbTabWidth {
		m.tabShowOffset = (m.SelectedTab + 1) - nbTabWidth
	}
	if m.SelectedTab < m.tabShowOffset {
		m.tabShowOffset = m.SelectedTab
	}
	if m.displayConfig.DynamicTabSize {
		m.resizeTabs(width)
	}
	if len(m.Tabs)-(m.tabShowOffset+1) < nbTabWidth {
		m.tabShowOffset = len(m.Tabs) - (nbTabWidth)
	}
	if m.tabShowOffset < 0 {
		m.tabShowOffset = 0
	}

	tabtext := ""

	for i := m.tabShowOffset; i <= m.tabShowOffset+nbTabWidth; i++ {
		if i >= len(m.Tabs) {
			break
		}
		tab := m.Tabs[i]
		ins := m.GetInstanceIndex(tab)

		t := ""

		if m.displayConfig.ShowTabNumber {
			t += fmt.Sprintf("%d ", i+1)
		}
		if m.displayConfig.ShowTabTitle {
			t += path.Base(ins.DirPath) + " "
		}
		tabtext += t
	}

	tview.Print(screen, tabtext, x, y, width, tview.AlignLeft, tcell.ColorDefault)
}

func (m *Tabline) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return m.WrapInputHandler(func(event *tcell.EventKey, _ func(p tview.Primitive)) {
		m.inputHandler.listenInputKey(event, "tabline", false)
	})
}

func (m *Tabline) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (bool, tview.Primitive) {
	return m.WrapMouseHandler(func(action tview.MouseAction, event *tcell.EventMouse, _ func(p tview.Primitive)) (bool, tview.Primitive) {
		if !m.InRect(event.Position()) {
			return false, nil
		}
		return m.inputHandler.listenInputMouse(event, action, "tabline"), m.Box
	})
}

func (m *Tabline) SelectTab(index int) {
	if m.SelectedTab != index && index >= 0 && index < len(m.Tabs) {
		m.SelectedTab = index
	}
}

func (m *Tabline) TabNext() {
	i := m.SelectedTab + 1
	if i > len(m.Tabs)-1 {
		i = len(m.Tabs) - 1
	}
	m.SelectTab(i)
}

func (m *Tabline) TabPrev() {
	i := m.SelectedTab - 1
	if i < 0 {
		i = 0
	}
	m.SelectTab(i)
}

func (m *Tabline) TabFirst() {
	m.SelectTab(0)
}

func (m *Tabline) TabLast() {
	m.SelectTab(len(m.Tabs) - 1)
}

func (m *Tabline) AddTab(instanceIndex int) int {
	m.Tabs = append(m.Tabs, instanceIndex)
	m.setCanClose()
	return len(m.Tabs) - 1
}

func (m *Tabline) AddTabs(instancesIndex []int) int {
	m.Tabs = append(m.Tabs, instancesIndex...)
	m.setCanClose()
	return len(m.Tabs) - 1
}

func (m *Tabline) CloseTab(index int) {
	if !m.canCloseTab {
		return
	}

	m.Tabs = append(m.Tabs[:index], m.Tabs[index+1:]...)
	m.setCanClose()
	if index == m.SelectedTab {
		if m.SelectedTab >= len(m.Tabs) {
			m.SelectTab(m.SelectedTab - 1)
		} else {
			m.SelectTab(m.SelectedTab)
		}
	} else {
		if index < m.SelectedTab {
			m.SelectedTab -= 1
		}
		m.SelectTab(m.SelectedTab)
	}
}

func (m *Tabline) GetSelectedTab() int {
	return m.Tabs[m.SelectedTab]
}

func (m *Tabline) GetUnderMouseTabIndex(mousePosX, tablinePosX, tablineWidth int) int {
	tablen := m.getTabSizeCompute()
	nbTabWidth := tablineWidth / tablen
	xx := mousePosX - tablinePosX
	if xx < 0 {
		return m.tabShowOffset
	}
	if xx > tablen*len(m.Tabs)-1 {
		return m.tabShowOffset + nbTabWidth
	}
	ind := ((m.tabShowOffset * tablen) + xx) / tablen
	if ind >= len(m.Tabs) {
		ind = len(m.Tabs) - 1
	}
	return ind
}

func (m *Tabline) GetUnderNoTab(mousePosX, tablinePosX, tablineWidth int) bool {
	tablen := m.getTabSizeCompute()
	xx := mousePosX - tablinePosX
	if xx < 0 {
		return true
	}
	if xx > tablineWidth {
		return true
	}
	if xx >= tablen*len(m.Tabs) {
		return true
	}
	return false
}

func (m *Tabline) setCanClose() {
	if len(m.Tabs) > 1 {
		m.canCloseTab = true
	} else {
		m.canCloseTab = false
	}
}

func (m *Tabline) clearTabs() {
	m.Tabs = []int{}
	m.SelectedTab = -1
}

func (m *Tabline) resizeTabs(tablineWidth int) {
	if !m.displayConfig.DynamicTabSize {
		return
	}

	rw := len(m.Tabs) * m.displayConfig.TabLen
	if rw > tablineWidth {
		m.tablen = tablineWidth / len(m.Tabs)
	} else {
		m.tablen = m.displayConfig.TabLen
	}
}

func (m *Tabline) getTabSizeCompute() int {
	CountDigits := func(i int) int {
		count := 0
		for i > 0 {
			i = i / 10
			count++
		}
		return count
	}
	tablen := m.tablen
	if !m.doShowTabName() {
		numLen := CountDigits(len(m.Tabs))
		tablen = 2 + numLen
		if m.canCloseTab {
			tablen += 1
		}
	}
	return tablen
}

func (m *Tabline) doShowTabName() bool {
	sLimit := 7
	if m.showCloseButton {
		sLimit++
	}
	return m.displayConfig.ShowTabTitle && m.tablen > sLimit
}
