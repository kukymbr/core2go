package service

import (
	"context"

	"github.com/kukymbr/core2go/di"
)

// RunnerRunFn is a custom function to run in runner.
type RunnerRunFn func(ctx context.Context, ctn *di.Container) error

// NewCustomRunner returns a Runner executing the custom function.
func NewCustomRunner(fn RunnerRunFn) Runner {
	return &customRunner{fn: fn}
}

type customRunner struct {
	fn RunnerRunFn
}

func (r *customRunner) Run(ctx context.Context, ctn *di.Container) error {
	return r.fn(ctx, ctn)
}
