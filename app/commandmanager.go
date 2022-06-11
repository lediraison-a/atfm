package app

import (
	"atfm/app/models"
	"bytes"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

type CommandManager struct {
	RootCmd    *cobra.Command
	CmdOut     *bytes.Buffer
	CmdHistory []string
}

func NewCommandManager() *CommandManager {
	rcmd := cobra.Command{Use: ""}
	bout := &bytes.Buffer{}

	return &CommandManager{
		RootCmd: &rcmd,
		CmdOut:  bout,
	}
}

func (c *CommandManager) RunCommand(command string) error {
	if command == "" {
		return nil
	}
	if command[0] == '!' {
		return c.RunCommandShell(command[1:])
	}
	c.CmdHistory = append(c.CmdHistory, command)
	c.RootCmd.SetArgs(strings.Split(command, " "))
	c.RootCmd.SetErr(c.CmdOut)
	c.RootCmd.SetOut(c.CmdOut)
	c.CmdOut.Reset()
	return c.RootCmd.Execute()
}

func (c *CommandManager) RunCommandShell(command string) error {
	if command == "" {
		return nil
	}
	co := strings.Split(command, " ")
	cmd := exec.Command(co[0], co[1:]...)
	c.CmdOut.Reset()
	cmd.Stdout = c.CmdOut
	cmd.Stderr = c.CmdOut
	cmd.Run()
	return nil
}

func (c *CommandManager) SetCurrentWorkingDir(ins *Instance) {
	wd, err := os.Getwd()
	if err != nil {
		return
	}
	if ins.Mod == models.LOCALFM && wd != ins.DirPath {
		err = os.Chdir(ins.DirPath)
		if err != nil {
			return
		}
	}
}
