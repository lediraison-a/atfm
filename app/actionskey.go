package app

import (
	"atfm/generics"
	"strings"

	"github.com/spf13/cobra"
)

func (t *Tui) GetActionsKey(commands []*cobra.Command) []*KeyAction {
	createCommandKeyFromCobra := func(cmd *cobra.Command) *KeyAction {
		cmdk := KeyAction{
			Name:   cmd.Use,
			Source: "",
			Action: func() {
				cmd.Run(cmd, []string{})
			},
		}
		return &cmdk
	}
	acs := generics.Map(commands, func(value *cobra.Command, _ int) *KeyAction {
		return createCommandKeyFromCobra(value)
	})

	normalmod := KeyAction{
		Name:   "normalmod",
		Source: "",
		Action: func() {
			t.layers.HidePage("pager")
			t.layers.ShowPage("main")

			t.app.SetFocus(t.filelists[t.selectedPane])
		},
	}

	opencommandline := KeyAction{
		Name:   "cmdinput",
		Source: "",
		Action: func() {
			t.app.SetFocus(t.inputLine)
			t.inputLine.OpenCommandLine()
		},
	}

	opensearchline := KeyAction{
		Name:   "searchinput",
		Source: "",
		Action: func() {
			t.app.SetFocus(t.inputLine)
			t.inputLine.OpenSearchLine()
		},
	}

	cmdprevious := KeyAction{
		Name:   "cmdprevious",
		Source: "",
		Action: func() {
			cmd := t.cmdManager.CmdPrevious()
			if t.inputLine.source == COMMAND_LINE {
				t.inputLine.SetText(cmd)
			}
		},
	}

	cmdnext := KeyAction{
		Name:   "cmdnext",
		Source: "",
		Action: func() {
			cmd := t.cmdManager.CmdNext()
			if t.inputLine.source == COMMAND_LINE {
				t.inputLine.SetText(cmd)
			}
		},
	}

	searchprevious := KeyAction{
		Name:   "searchprevious",
		Source: "",
		Action: func() {
			ins := t.getInstancePane(t.selectedPane)
			st := ins.QuickSearch.SearchPrevious()
			if t.inputLine.source == SEARCH_LINE {
				t.inputLine.SetText(st)
			}
		},
	}

	searchnext := KeyAction{
		Name:   "searchnext",
		Source: "",
		Action: func() {
			ins := t.getInstancePane(t.selectedPane)
			st := ins.QuickSearch.SearchNext()
			if t.inputLine.source == SEARCH_LINE {
				t.inputLine.SetText(st)
			}
		},
	}

	searchjumpforward := KeyAction{
		Name:   "searchjumpforward",
		Source: "",
		Action: func() {
			ins := t.getInstancePane(t.selectedPane)
			ins.QuickSearch.SearchJumpForward(ins)
		},
	}

	searchjumpbackward := KeyAction{
		Name:   "searchjumpbackward",
		Source: "",
		Action: func() {
			ins := t.getInstancePane(t.selectedPane)
			ins.QuickSearch.SearchJumpBackward(ins)
		},
	}

	runsearch := KeyAction{
		Name:   "searchrun",
		Source: "",
		Action: func() {
			searchText := t.inputLine.GetText()
			ins := t.getInstancePane(t.selectedPane)
			ins.QuickSearch.SearchContent(searchText, ins)
			t.app.SetFocus(t.filelists[t.selectedPane])
		},
	}

	// TODO implement this better
	runcommand := KeyAction{
		Name:   "cmdrun",
		Source: "",
		Action: func() {
			ins := t.getInstancePane(t.selectedPane)
			command := t.inputLine.GetText()
			if command == "" {
				return
			}
			if command[0] == '!' {
				t.cmdManager.RunCommandShell(command[1:], ins.DirPath, ins.Mod)
			} else {
				t.cmdManager.RunCommand(command)
			}
			ct := t.cmdManager.CmdOut
			if ct.Len() > 0 {
				t.pager.SetText(ct.String())
				nblines := len(strings.Split(ct.String(), "\n")) + 1
				x, y, width, height := t.layers.GetRect()
				t.pager.SetRect(x, y+(height-nblines), width, nblines+1)
				t.layers.ShowPage("pager")
			} else {
				t.app.SetFocus(t.filelists[t.selectedPane])
			}
		},
	}

	return append(acs,
		&normalmod,
		&opencommandline,
		&runcommand,
		&opensearchline,
		&runsearch,
		&searchjumpbackward,
		&searchjumpforward,
		&cmdnext,
		&cmdprevious,
		&searchnext,
		&searchprevious)
}
