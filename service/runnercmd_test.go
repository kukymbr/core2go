package service_test

import (
	"context"
	"testing"

	"github.com/kukymbr/core2go/service"
	"github.com/stretchr/testify/assert"
)

func TestCommandRunner_Run_WhenFailed_ExpectError(t *testing.T) {
	runner := service.NewCommandRunner(&service.NopCommand{})
	ctx, cancel := context.WithCancel(context.Background())

	cancel()

	err := runner.Run(ctx, nil)

	assert.Error(t, err)
}
