package commands

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// CmdOutput represents the cached output of an executable command
type CmdOutput struct {
	Output  *string
	RunTime *time.Time
}

// Cmd represents an executable command
type Cmd struct {
	exec.Cmd
	CmdOutput
}

// NewCmd creates an executable command
func NewCmd(name string, arg ...string) *Cmd {
	return &Cmd{Cmd: *exec.Command(name, arg...)}
}

// ToString converts an executable command to a string
func (cmd *Cmd) ToString() string {
	var toString string
	if len(cmd.Args) > 0 {
		toString = cmd.Args[0]
		for _, arg := range cmd.Args[1:] {
			toString = strings.Join([]string{toString, arg}, " ")
		}
	} else {
		toString = cmd.Path
	}
	return toString
}

// Run executes an executable command
func (cmd *Cmd) Run(cacheFirst bool) *CmdOutput {
	if cacheFirst && cmd.CmdOutput.Output != nil {
		return &cmd.CmdOutput
	}

	var output string
	// We cannot call CombinedOutput twice on the same exec.Cmd
	cmd.Cmd = *exec.Command(cmd.Cmd.Args[0], cmd.Cmd.Args[1:]...)
	stdoutStderr, err := cmd.Cmd.CombinedOutput()
	if err != nil {
		output = fmt.Sprintf("%s", stdoutStderr)
	} else {
		output = fmt.Sprintf("%s", stdoutStderr)
	}

	cmd.CmdOutput.Output = &output
	now := time.Now()
	cmd.CmdOutput.RunTime = &now

	return &cmd.CmdOutput
}
