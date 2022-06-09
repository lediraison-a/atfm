package app

import (
	"atfm/generics"

	"github.com/spf13/cobra"
)

func (t *Tui) GetActionsMouse(commands []*cobra.Command) []*MouseAction {
	createCommandMouseFromCobra := func(cmd *cobra.Command) *MouseAction {
		cmdk := MouseAction{
			Name: cmd.Use,
			Action: func(_, _ int) {
				cmd.Run(cmd, []string{})
			},
		}
		return &cmdk
	}
	acs := generics.Map(commands, func(value *cobra.Command, _ int) *MouseAction {
		return createCommandMouseFromCobra(value)
	})

	setcurrent := MouseAction{
		Name: "setcurrent",
		Action: func(_, y int) {
			t.getInstancePane(t.selectedPane).
				CurrentItem = t.filelists[t.selectedPane].
				GetUnderMouseIndex(y)
		},
	}
	tabsetcurrent := MouseAction{
		Name: "tabsetcurrent",
		Action: func(x, _ int) {
			tl := t.tablines[t.selectedPane]
			ind := tl.GetUnderMouseTabIndex(x)
			if ind >= len(tl.Tabs) {
				ind = len(tl.Tabs) - 1
			}
			tl.SelectTab(ind)
		},
	}
	openpath := MouseAction{
		Name: "openpath",
		Action: func(x, _ int) {
			t.pathlines[t.selectedPane].OpenPathPos(x)
		},
	}
	return append(acs, &setcurrent, &tabsetcurrent, &openpath)
}
