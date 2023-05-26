package service

import "context"

// ContextWithCancel is a context with cancel function
type ContextWithCancel interface {
	GetContext() context.Context
	GetCancelFn() context.CancelFunc
}

// NewBaseContext returns new default base service context
func NewBaseContext(parent context.Context) ContextWithCancel {
	ctx, cancel := context.WithCancel(parent)

	return &BaseContext{
		ctx:    ctx,
		cancel: cancel,
	}
}

// BaseContext is a default service context
type BaseContext struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func (c *BaseContext) GetContext() context.Context {
	return c.ctx
}

func (c *BaseContext) GetCancelFn() context.CancelFunc {
	return c.cancel
}
