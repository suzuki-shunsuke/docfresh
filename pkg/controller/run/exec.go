package run

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)

type CommandResult struct {
	Stdout         string
	Stderr         string
	CombinedOutput string
	ExitCode       int
}

const waitDelay = 1000 * time.Hour

func setCancel(cmd *exec.Cmd) {
	cmd.Cancel = func() error {
		return cmd.Process.Signal(os.Interrupt)
	}
	cmd.WaitDelay = waitDelay
}

func (c *Controller) exec(ctx context.Context, block *Block) (*CommandResult, error) {
	shell := block.Input.Command.Shell
	if shell == nil {
		shell = []string{"bash", "-c"}
	}
	cmd := exec.CommandContext(ctx, shell[0], append(shell[1:], block.Input.Command.Command)...) //nolint:gosec
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	combinedOutput := &bytes.Buffer{}
	cmd.Stdout = io.MultiWriter(os.Stdout, stdout, combinedOutput)
	cmd.Stderr = io.MultiWriter(os.Stderr, stderr, combinedOutput)
	setCancel(cmd)
	fmt.Fprintln(os.Stderr, "+", block.Input.Command.Command)
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("execute a command: %w", err)
	}
	return &CommandResult{
		Stdout:         stdout.String(),
		Stderr:         stderr.String(),
		CombinedOutput: combinedOutput.String(),
		ExitCode:       cmd.ProcessState.ExitCode(),
	}, nil
}
