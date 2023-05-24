package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/kukymbr/core2go/di"
)

// SetDIContainer is a middleware to set a DI container instance to the context
func SetDIContainer(ctn *di.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(ContextKeyDIContainer, ctn)
	}
}
