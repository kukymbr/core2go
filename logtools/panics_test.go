package logtools_test

import (
	"testing"

	"github.com/kukymbr/core2go/logtools"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestCatchPanic(t *testing.T) {
	log := zap.Must(zap.NewProduction())

	assert.NotPanics(t, func() {
		defer logtools.CatchPanic(log)

		panic("test panic")
	})
}

func TestCatchPanic_WhenWithCallback_ExpectExecuted(t *testing.T) {
	log := zap.Must(zap.NewProduction())
	executed := false

	assert.NotPanics(t, func() {
		defer logtools.CatchPanic(log, func(recovered any) {
			executed = true

			assert.Equal(t, "test panic", recovered)
		})

		panic("test panic")
	})

	assert.True(t, executed)
}

func TestCatchPanic_WhenNoPanic_ExpectNothing(t *testing.T) {
	log := zap.Must(zap.NewProduction())
	executed := false

	assert.NotPanics(t, func() {
		defer logtools.CatchPanic(log, func(_ any) {
			executed = true
		})
	})

	assert.False(t, executed)
}
