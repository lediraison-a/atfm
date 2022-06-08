package app

import (
	"atfm/generics"

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
			t.app.SetFocus(t.filelists[t.selectedPane])
		},
	}

	return append(acs, &normalmod)
}
