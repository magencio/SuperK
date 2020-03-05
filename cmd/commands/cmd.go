package commands

import (
	"bytes"
	"io"
	"os/exec"
	"strings"
	"sync"
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

// The following code will wait until the command finishes to get the output, so launching a
// command with e.g. -w will hang the app
//
// Run executes an executable command
// func (cmd *Cmd) Run(cacheFirst bool) *CmdOutput {
// 	if cacheFirst && cmd.CmdOutput.Output != nil {
// 		return &cmd.CmdOutput
// 	}

// 	var output string
// 	// We cannot call CombinedOutput twice on the same exec.Cmd
// 	cmd.Cmd = *exec.Command(cmd.Cmd.Args[0], cmd.Cmd.Args[1:]...)
// 	// TODO: This will hang with -w commands
// 	stdoutStderr, err := cmd.Cmd.CombinedOutput()
// 	if err != nil {
// 		output = fmt.Sprintf("%s", stdoutStderr)
// 	} else {
// 		output = fmt.Sprintf("%s", stdoutStderr)
// 	}

// 	cmd.CmdOutput.Output = &output
// 	now := time.Now()
// 	cmd.CmdOutput.RunTime = &now

// 	return &cmd.CmdOutput
// }

// Run executes an executable command
// This code is based on: https://blog.kowalczyk.info/article/wOYk/advanced-command-execution-in-go-with-osexec.html
// The idea is that with a minor modification I could show progress somewhere else while getting the output
func (cmd *Cmd) Run(cacheFirst bool) *CmdOutput {
	if cacheFirst && cmd.CmdOutput.Output != nil {
		return &cmd.CmdOutput
	}

	var output string
	// We cannot run a command twice on the same exec.Cmd
	cmd.Cmd = *exec.Command(cmd.Cmd.Args[0], cmd.Cmd.Args[1:]...)

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return nil
	}

	var stdoutBuffer, stderrBuffer bytes.Buffer
	err = cmd.Start()
	if err != nil {
		return nil
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		_, err = io.Copy(&stdoutBuffer, stdoutPipe)
		if err != nil {
			return
		}
		wg.Done()
	}()

	_, err = io.Copy(&stderrBuffer, stderrPipe)
	if err != nil {
		return nil
	}
	wg.Wait()

	err = cmd.Wait()
	if err != nil || err == nil {
		builder := strings.Builder{}
		builder.WriteString(stderrBuffer.String())
		builder.WriteString(stdoutBuffer.String())
		output = builder.String()
	}

	cmd.CmdOutput.Output = &output
	now := time.Now()
	cmd.CmdOutput.RunTime = &now

	return &cmd.CmdOutput
}
