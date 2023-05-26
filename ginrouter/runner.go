package ginrouter

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// NewRunner returns new gin Runner instance
func NewRunner(router *gin.Engine, addr string) *Runner {
	return &Runner{
		router: router,
		addr:   addr,
	}
}

// Runner is a service gin runner
type Runner struct {
	router *gin.Engine
	addr   string
}

func (r *Runner) Run() error {
	addrArgs := make([]string, 0)

	if r.addr != "" {
		addrArgs = append(addrArgs, r.addr)
	}

	if err := r.router.Run(addrArgs...); err != nil {
		return fmt.Errorf("run gin router: %w", err)
	}

	return nil
}
