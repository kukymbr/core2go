package service

import (
	"context"

	"github.com/kukymbr/core2go/di"
)

// Runner is an interface for a service's runner,
// such as gin or something else.
type Runner interface {
	// Run runs the Runner.
	Run(ctx context.Context, ctn *di.Container) error
}
