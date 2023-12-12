package service

import (
	"context"
	"fmt"

	"github.com/kukymbr/core2go/di"
)

// Command is a CLI command interface.
// The cobra.Command fits, see: https://github.com/spf13/cobra
type Command interface {
	ExecuteContext(ctx context.Context) error
}

// NewCommandRunner creates new CLI command Service Runner.
func NewCommandRunner(command Command) *CommandRunner {
	return &CommandRunner{command: command}
}

// CommandRunner is a Service Runner, executing the CLI command.
type CommandRunner struct {
	command Command
}

func (r *CommandRunner) Run(ctx context.Context, _ *di.Container) error {
	if err := r.command.ExecuteContext(ctx); err != nil {
		return fmt.Errorf("command execute: %w", err)
	}

	return nil
}

// NopCommand is a Command doing nothing.
type NopCommand struct {
	Executed bool
}

func (c *NopCommand) ExecuteContext(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	c.Executed = true

	return nil
}
