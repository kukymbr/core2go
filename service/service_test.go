package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/kukymbr/core2go/di"
	"github.com/kukymbr/core2go/service"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type panickingCmd struct{}

func (c *panickingCmd) ExecuteContext(_ context.Context) error {
	panic("test panic")
}

func getSrv() *service.Service {
	ctn := &di.Container{}
	log := zap.Must(zap.NewDevelopment())

	return service.New(ctn, log)
}

func TestService_Run(t *testing.T) {
	srv := getSrv()
	cmd1 := &service.NopCommand{}
	cmd2 := &service.NopCommand{}

	srv.RegisterRunner(service.NewCommandRunner(cmd1), service.NewCommandRunner(cmd2))
	code := srv.Run(context.Background())

	assert.Equal(t, 0, code)
	assert.True(t, cmd1.Executed)
	assert.True(t, cmd2.Executed)
}

func TestService_Run_WhenTimeout_ExpectErrorCode(t *testing.T) {
	srv := getSrv()
	cmd := &service.NopCommand{}
	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	srv.RegisterRunner(service.NewCommandRunner(cmd), service.NewRouterRunner(&service.NopRouter{}))
	code := srv.Run(ctx)

	assert.NotEqual(t, 0, code)
	assert.True(t, cmd.Executed)
	assert.GreaterOrEqual(t, time.Since(start), 200*time.Millisecond)
}

func TestService_Run_WhenNoRunners_ExpectErrorCode(t *testing.T) {
	srv := service.New(&di.Container{}, nil)
	code := srv.Run(context.Background())

	assert.NotEqual(t, 0, code)
}

func TestService_Run_WhenCanceled_ExpectErrorCode(t *testing.T) {
	srv := getSrv()
	cmd := &service.NopCommand{}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	srv.RegisterRunner(service.NewCommandRunner(cmd), service.NewRouterRunner(&service.NopRouter{}))
	code := srv.Run(ctx)

	assert.NotEqual(t, 0, code)
	assert.False(t, cmd.Executed)
}

func TestService_Run_WhenRegisteringRunnerAfterRun_ExpectPanic(t *testing.T) {
	srv := getSrv()

	srv.RegisterRunner(service.NewCommandRunner(&service.NopCommand{}))
	srv.Run(context.Background())

	assert.Panics(t, func() {
		srv.RegisterRunner(service.NewCommandRunner(&service.NopCommand{}))
	})
}

func TestService_Run_WhenRunnerPanicked_ExpectNoPanic(t *testing.T) {
	srv := getSrv()

	srv.RegisterRunner(service.NewCommandRunner(&panickingCmd{}))

	assert.NotPanics(t, func() {
		code := srv.Run(context.Background())

		assert.NotEqual(t, 0, code)
	})

}
