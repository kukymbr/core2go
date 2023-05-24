package ginrouter

import (
	"net/http"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const maxMultipartMemoryDft = 8 << 20

// GetDefaultRouter creates new gin router with default middlewares
func GetDefaultRouter(logger *zap.Logger) *gin.Engine {
	router := gin.New()

	router.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	router.Use(ginzap.RecoveryWithZap(logger, true))

	// 8 Mb limit for uploads
	router.MaxMultipartMemory = maxMultipartMemoryDft

	// Custom no-route error
	router.NoRoute(func(c *gin.Context) {
		ctx := NewContextHandler(c)
		ctx.ErrResponseS("", http.StatusNotFound)
	})

	return router
}
