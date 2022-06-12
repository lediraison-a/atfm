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
	RootCmd            *cobra.Command
	CmdOut             *bytes.Buffer
	CmdHistory         []string
	selectedCmdHistory int
}

func NewCommandManager() *CommandManager {
	rcmd := cobra.Command{Use: ""}
	bout := &bytes.Buffer{}

	return &CommandManager{
		RootCmd:            &rcmd,
		CmdOut:             bout,
		CmdHistory:         []string{},
		selectedCmdHistory: 0,
	}
}

func (c *CommandManager) RunCommand(command string) error {
	c.CmdHistory = append(c.CmdHistory, command)
	c.selectedCmdHistory = len(c.CmdHistory)
	c.RootCmd.SetArgs(strings.Split(command, " "))
	c.RootCmd.SetErr(c.CmdOut)
	c.RootCmd.SetOut(c.CmdOut)
	c.CmdOut.Reset()
	return c.RootCmd.Execute()
}

func (c *CommandManager) RunCommandShell(command, workdir string, mod models.FsMod) error {
	if command == "" {
		return nil
	}
	c.CmdHistory = append(c.CmdHistory, "!"+command)
	c.selectedCmdHistory = len(c.CmdHistory)
	c.SetCurrentWorkingDir(workdir, mod)
	co := strings.Split(command, " ")
	cmd := exec.Command(co[0], co[1:]...)
	c.CmdOut.Reset()
	cmd.Stdout = c.CmdOut
	cmd.Stderr = c.CmdOut
	cmd.Run()
	return nil
}

func (c *CommandManager) SetCurrentWorkingDir(workdir string, mod models.FsMod) {
	wd, err := os.Getwd()
	if err != nil {
		return
	}
	if mod == models.LOCALFM && wd != workdir {
		err = os.Chdir(workdir)
		if err != nil {
			return
		}
	}
}

func (c *CommandManager) CmdNext() string {
	histLen := len(c.CmdHistory)
	if histLen == 0 {
		return ""
	}
	c.selectedCmdHistory++
	t := ""
	if c.selectedCmdHistory > histLen {
		c.selectedCmdHistory = histLen
	}
	if c.selectedCmdHistory != histLen {
		t = c.CmdHistory[c.selectedCmdHistory]
	}
	return t
}

func (c *CommandManager) CmdPrevious() string {
	if len(c.CmdHistory) == 0 {
		return ""
	}
	c.selectedCmdHistory--
	if c.selectedCmdHistory < 0 {
		c.selectedCmdHistory = 0
	}
	t := c.CmdHistory[c.selectedCmdHistory]
	return t
}
