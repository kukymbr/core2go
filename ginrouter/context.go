package ginrouter

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kukymbr/core2go/di"
	"github.com/kukymbr/core2go/errs"
	"github.com/kukymbr/core2go/ginrouter/middlewares"
)

// NewContextHandler creates new ContextHandler instance
func NewContextHandler(ctx *gin.Context) *ContextHandler {
	return &ContextHandler{
		Context: ctx,
	}
}

// ContextHandler is a gin context wrapper
type ContextHandler struct {
	*gin.Context
}

// GetContainer returns di.Container instance
func (h *ContextHandler) GetContainer() *di.Container {
	ctn, ok := h.Get(middlewares.ContextKeyDIContainer)
	if !ok {
		panic("attempt to access non-initialized DI Container")
	}

	return ctn.(*di.Container)
}

// ErrResponse sends error as response
func (h *ContextHandler) ErrResponse(err error) {
	status := http.StatusInternalServerError

	var e *errs.ExposableError

	ok := errors.As(err, &e)
	if ok {
		if e.Code >= 400 && e.Code <= 599 {
			status = e.Code
		}
	}

	h.ErrResponseWithStatus(err, status)
}

// ErrResponseS sends error message as response
func (h *ContextHandler) ErrResponseS(msg string, status int) {
	h.ErrResponse(errs.NewExposable(status, msg))
}

// ErrResponseWithStatus sends error as response with custom status
func (h *ContextHandler) ErrResponseWithStatus(err error, status int) {
	var e *errs.ExposableError

	ok := errors.As(err, &e)
	if !ok {
		err = errs.NewExposable(status, err.Error())
		errors.As(err, &e)
	}

	h.JSON(
		status,
		e,
	)
}
