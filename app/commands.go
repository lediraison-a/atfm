package app

import (
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
