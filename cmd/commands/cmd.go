package commands

import (
	"fmt"
	"os/exec"
	"strings"
)

// Cmd represents an executable command
type Cmd struct {
	exec.Cmd
}

// NewCmd creates an executable command
func NewCmd(name string, arg ...string) *Cmd {
	return &Cmd{Cmd: *exec.Command(name, arg...)}
}

// ToString converts an executable command to a string
func (command *Cmd) ToString() string {
	var toString string
	if len(command.Args) > 0 {
		toString = command.Args[0]
		for _, arg := range command.Args[1:] {
			toString = strings.Join([]string{toString, arg}, " ")
		}
	} else {
		toString = command.Path
	}
	return toString
}

// Run executes an executable command
func (command *Cmd) Run() *string {
	var output string
	stdoutStderr, err := command.Cmd.CombinedOutput()
	if err != nil {
		output = fmt.Sprintf("%s", stdoutStderr)
	} else {
		output = fmt.Sprintf("%s", stdoutStderr)
	}

	return &output
}
