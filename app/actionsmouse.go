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
	return acs
}
