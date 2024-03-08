package service

import (
	"context"
	"fmt"

	"github.com/kukymbr/core2go/di"
)

// Router interface.
// The graceful.Graceful gin's wrapper fits: https://github.com/gin-contrib/graceful.
type Router interface {
	RunWithContext(ctx context.Context) error
}

// NewRouterRunner creates new Runner executing the given http Router.
func NewRouterRunner(router Router) Runner {
	return NewCustomRunner(func(ctx context.Context, _ *di.Container) error {
		if err := router.RunWithContext(ctx); err != nil {
			return fmt.Errorf("run router: %w", err)
		}

		return nil
	})
}

// NopRouter is a Router doing nothing.
type NopRouter struct{}

func (r *NopRouter) RunWithContext(ctx context.Context) error {
	<-ctx.Done()

	return ctx.Err()
}
