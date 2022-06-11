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
		Name:   "opencommandline",
		Source: "",
		Action: func() {
			t.app.SetFocus(t.commandLine)
		},
	}

	runcommand := KeyAction{
		Name:   "cmdrun",
		Source: "",
		Action: func() {
			command := t.commandLine.GetText()
			t.cmdManager.RunCommand(command)
			ct := t.cmdManager.CmdOut
			if ct.Len() > 0 {
				t.pager.SetText(ct.String())
				nblines := len(strings.Split(ct.String(), "\n"))
				x, y, width, height := t.layers.GetRect()
				t.pager.SetRect(x, y+(height-nblines), width, nblines+1)
				t.layers.ShowPage("pager")
			} else {
				t.app.SetFocus(t.filelists[t.selectedPane])
			}
		},
	}

	return append(acs, &normalmod, &opencommandline, &runcommand)
}
