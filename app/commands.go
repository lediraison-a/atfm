package app

import (
	"atfm/app/models"
	"bytes"

	"github.com/spf13/cobra"
)

type Commands struct {
	RootCmd *cobra.Command
	CmdOut  *bytes.Buffer
}

func NewCommands() *Commands {
	rcmd := cobra.Command{Use: ""}
	bout := &bytes.Buffer{}
	rcmd.SetErr(bout)
	rcmd.SetOut(bout)

	return &Commands{
		RootCmd: &rcmd,
		CmdOut:  bout,
	}
}

func (t *Tui) GetAppCommands() []*cobra.Command {
	toggleDoublePane := &cobra.Command{
		Use:  "toggledoublepane",
		Args: cobra.ExactArgs(0),
		Run: func(_ *cobra.Command, _ []string) {
			t.ToggleDoublePane()
		},
	}
	quitall := &cobra.Command{
		Use:     "quitall",
		Aliases: []string{"qa", "qall"},
		Args:    cobra.ExactArgs(0),
		Run: func(_ *cobra.Command, _ []string) {
			t.app.Stop()
		},
	}
	quit := &cobra.Command{
		Use:     "quit",
		Aliases: []string{"q"},
		Args:    cobra.ExactArgs(0),
		Run: func(_ *cobra.Command, _ []string) {
			thisTabline := t.tablines[t.selectedPane]
			if len(thisTabline.Tabs) == 1 {
				if t.showDoublePane {
					t.ToggleDoublePane()
					thisTabline.CloseTab(len(thisTabline.Tabs) - 1)
				} else {
					t.app.Stop()
				}
			} else {
				thisTabline.CloseTab(thisTabline.SelectedTab)
			}
		},
	}
	tabnew := &cobra.Command{
		Use:  "tabnew",
		Args: cobra.ExactArgs(0),
		Run: func(_ *cobra.Command, _ []string) {
			t.NewInstance(t.appConfig.Start.StartDir, t.appConfig.Start.StartBasepath, models.LOCALFM, true)
		},
	}
	tabclose := &cobra.Command{
		Use:  "tabclose",
		Args: cobra.ExactArgs(0),
		Run: func(_ *cobra.Command, _ []string) {
			tl := t.tablines[t.selectedPane]
			if tl.canCloseTab {
				tl.CloseTab(tl.SelectedTab)
			}
		},
	}
	tabnext := &cobra.Command{
		Use:  "tabnext",
		Args: cobra.ExactArgs(0),
		Run: func(_ *cobra.Command, _ []string) {
			t.tablines[t.selectedPane].TabNext()
		},
	}
	tabprevious := &cobra.Command{
		Use:  "tabprevious",
		Args: cobra.ExactArgs(0),
		Run: func(_ *cobra.Command, _ []string) {
			t.tablines[t.selectedPane].TabPrev()
		},
	}
	tabfirst := &cobra.Command{
		Use:  "tabfirst",
		Args: cobra.ExactArgs(0),
		Run: func(_ *cobra.Command, _ []string) {
			t.tablines[t.selectedPane].TabNext()
		},
	}
	tablast := &cobra.Command{
		Use:  "tablast",
		Args: cobra.ExactArgs(0),
		Run: func(_ *cobra.Command, _ []string) {
			t.tablines[t.selectedPane].TabPrev()
		},
	}
	scrollup := &cobra.Command{
		Use:  "scrollup",
		Args: cobra.ExactArgs(0),
		Run: func(_ *cobra.Command, _ []string) {
			t.filelists[t.selectedPane].ScrollUp()
		},
	}
	scrolldown := &cobra.Command{
		Use:  "scrolldown",
		Args: cobra.ExactArgs(0),
		Run: func(_ *cobra.Command, _ []string) {
			t.filelists[t.selectedPane].ScrollDown()
		},
	}
	scrollfirst := &cobra.Command{
		Use:  "scrollfirst",
		Args: cobra.ExactArgs(0),
		Run: func(_ *cobra.Command, _ []string) {
			t.filelists[t.selectedPane].ScrollFirst()
		},
	}
	scrolllast := &cobra.Command{
		Use:  "scrolllast",
		Args: cobra.ExactArgs(0),
		Run: func(_ *cobra.Command, _ []string) {
			t.filelists[t.selectedPane].ScrollLast()
		},
	}
	editpath := &cobra.Command{
		Use:  "editpath",
		Args: cobra.ExactArgs(0),
		Run: func(_ *cobra.Command, _ []string) {
			t.app.SetFocus(t.pathlines[t.selectedPane])
			t.pathlines[t.selectedPane].EditPath()
		},
	}
	opencurrent := &cobra.Command{
		Use:  "opencurrent",
		Args: cobra.ExactArgs(0),
		Run: func(_ *cobra.Command, _ []string) {
			ins := t.getInstancePane(t.selectedPane)
			if err := ins.OpenAtIndex(ins.CurrentItem); err != nil {

			}
		},
	}
	openparent := &cobra.Command{
		Use:  "openparent",
		Args: cobra.ExactArgs(0),
		Run: func(_ *cobra.Command, _ []string) {
			ins := t.getInstancePane(t.selectedPane)
			pd, pbp, pm := ins.GetParentInfo()
			if err := ins.OpenDirSaveHistory(pd, pbp, pm); err != nil {

			}
		},
	}
	openprevious := &cobra.Command{
		Use:  "openprevious",
		Args: cobra.ExactArgs(0),
		Run: func(_ *cobra.Command, _ []string) {
			ins := t.getInstancePane(t.selectedPane)
			ok, err := ins.OpenHistoryDir(models.HISTORY_BACK)
			if ok && err != nil {

			}
		},
	}
	opennext := &cobra.Command{
		Use:  "opennext",
		Args: cobra.ExactArgs(0),
		Run: func(_ *cobra.Command, _ []string) {
			ins := t.getInstancePane(t.selectedPane)
			ok, err := ins.OpenHistoryDir(models.HISTORY_FORWARD)
			if ok && err != nil {

			}
		},
	}
	unselectall := &cobra.Command{
		Use:  "unselectall",
		Args: cobra.ExactArgs(0),
		Run: func(_ *cobra.Command, _ []string) {
			ins := t.getInstancePane(t.selectedPane)
			ins.UnselectAll()
		},
	}

	return []*cobra.Command{
		quit,
		quitall,
		toggleDoublePane,
		tabnew,
		tabclose,
		tabfirst,
		tablast,
		tabprevious,
		tabnext,
		scrollfirst,
		scrolldown,
		scrollup,
		scrolllast,
		editpath,
		opencurrent,
		openparent,
		openprevious,
		opennext,
		unselectall,
	}
}
