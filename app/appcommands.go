package app

import (
	"atfm/app/models"
	"path"

	"github.com/spf13/cobra"
)

func getCommandsMove(t *Tui) []*cobra.Command {
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
	pagedown := &cobra.Command{
		Use:  "pagedown",
		Args: cobra.ExactArgs(0),
		Run: func(_ *cobra.Command, _ []string) {
			ins := t.getInstancePane(t.selectedPane)
			fl := t.filelists[t.selectedPane]
			_, _, _, h := fl.GetRect()
			pos := fl.listOffset + (h - 1)
			if ins.CurrentItem == pos {
				ins.CurrentItem += h - 1
			} else {
				ins.CurrentItem = pos
			}

			if ins.CurrentItem >= len(ins.ShownContent) {
				ins.CurrentItem = len(ins.ShownContent) - 1
			}
		},
	}
	pageup := &cobra.Command{
		Use:  "pageup",
		Args: cobra.ExactArgs(0),
		Run: func(_ *cobra.Command, _ []string) {
			ins := t.getInstancePane(t.selectedPane)
			fl := t.filelists[t.selectedPane]
			_, _, _, h := fl.GetRect()
			pos := fl.listOffset
			if ins.CurrentItem == pos {
				ins.CurrentItem -= h - 1
			} else {
				ins.CurrentItem = pos
			}

			if ins.CurrentItem < 0 {
				ins.CurrentItem = 0
			}
		},
	}
	return []*cobra.Command{
		scrollup,
		scrolldown,
		scrolllast,
		scrollfirst,
		pagedown,
		pageup,
	}
}

func getCommandsTabs(t *Tui) []*cobra.Command {
	tabnew := &cobra.Command{
		Use:  "tabnew",
		Args: cobra.ExactArgs(0),
		Run: func(_ *cobra.Command, _ []string) {
			err := t.NewInstance(t.appConfig.Start.StartDir, t.appConfig.Start.StartBasepath, models.LOCALFM, true)
			if err != nil {
				t.inputLine.LogError(err.Error())
			}
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
	return []*cobra.Command{
		tabnew,
		tabclose,
		tabfirst,
		tablast,
		tabnext,
		tabprevious,
		tablast,
	}
}

func getCommandsFile(t *Tui) []*cobra.Command {
	rename := &cobra.Command{
		Use:  "rename",
		Args: cobra.ExactArgs(0),
		Run: func(_ *cobra.Command, _ []string) {
			ins := t.getInstanceGlobal()
			dirPath := ins.DirPath
			filename := ins.ShownContent[ins.CurrentItem].Name
			filepath := path.Join(dirPath, filename)
			onRename := func(newName string) (string, error) {
				newNamePath := path.Join(dirPath, newName)
				if len(newName) > 0 && newName[0] == '/' {
					newNamePath = newName
				}
				logInfo := filepath + " renamed to " + newNamePath
				return logInfo, ins.RenameFile(filepath, newNamePath)
			}
			label := "rename " + filename + " > "
			t.app.SetFocus(t.inputLine)
			t.inputLine.OpenInput(label, filename, onRename)
		},
	}

	return []*cobra.Command{
		rename,
	}
}

func getCommandsApp(t *Tui) []*cobra.Command {
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
				t.inputLine.LogError(err.Error())
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
				t.inputLine.LogError(err.Error())
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
				t.inputLine.LogError(err.Error())
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
				t.inputLine.LogError(err.Error())
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
	togglehiddenfiles := &cobra.Command{
		Use:  "togglehiddenfiles",
		Args: cobra.ExactArgs(0),
		Run: func(_ *cobra.Command, _ []string) {
			ins := t.getInstancePane(t.selectedPane)
			ins.ShowHidden = !ins.ShowHidden
			ins.ShownContent = ins.GetShownContent(ins.Content)
		},
	}
	return []*cobra.Command{
		quit,
		quitall,
		toggleDoublePane,
		editpath,
		opencurrent,
		openparent,
		openprevious,
		opennext,
		unselectall,
		togglehiddenfiles,
	}
}
