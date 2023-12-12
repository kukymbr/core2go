package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"

	"github.com/kukymbr/core2go/di"
	"github.com/kukymbr/core2go/logtools"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

var (
	errTerminated = errors.New("terminated")
	errCanceled   = errors.New("canceled")
)

// New creates new Service instance with the specified DI container.
func New(ctn *di.Container, log *zap.Logger) *Service {
	if log == nil {
		log = zap.Must(zap.NewProduction())
	}

	return &Service{
		ctn:     ctn,
		log:     log.With(zap.String("who", "core2go.Service")),
		runners: make([]Runner, 0),
	}
}

// Service is an application based on core2go kit.
type Service struct {
	ctn     *di.Container
	log     *zap.Logger
	runners []Runner

	executed atomic.Bool
}

// RegisterRunner register the Runner instances in the Service.
func (s *Service) RegisterRunner(runners ...Runner) {
	if s.executed.Load() {
		s.log.Panic("service is already executed, cannot register the runner")
	}

	s.runners = append(s.runners, runners...)
}

// Run starts the Service. Returns the exist code.
//
//nolint:funlen
func (s *Service) Run(ctx context.Context) int {
	var cancel context.CancelFunc

	s.executed.Store(true)

	ctx, cancel = context.WithCancel(ctx)
	defer func() {
		cancel()
	}()

	defer s.close()

	if len(s.runners) == 0 {
		s.log.Error("no runners registered in the Service instance")

		return 1
	}

	exitCode := atomic.Int32{}
	runnersDone := atomic.Int32{}

	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM)

	eg := errgroup.Group{}

	eg.Go(func() error {
		return s.listenCancel(ctx, cancel, &exitCode, shutdownChan)
	})

	for _, runner := range s.runners {
		runner := runner

		eg.Go(func() error {
			if err := ctx.Err(); err != nil {
				exitCode.Store(2)

				return err
			}

			if err := s.executeRunner(ctx, runner); err != nil {
				exitCode.Store(3)
				cancel()

				return err
			}

			done := runnersDone.Add(1)
			if int(done) == len(s.runners) {
				cancel()
			}

			return nil
		})
	}

	err := eg.Wait()

	cancel()

	if err != nil && !errors.Is(err, errTerminated) && !errors.Is(err, errCanceled) {
		s.log.Error(err.Error())
	}

	s.log.Sugar().Debugf("Got exit code: %d", exitCode.Load())

	return int(exitCode.Load())
}

func (s *Service) listenCancel(
	ctx context.Context,
	cancel context.CancelFunc,
	exitCode *atomic.Int32,
	shutdownChan chan os.Signal,
) error {
	for {
		select {
		case sig := <-shutdownChan:
			s.log.Info("Got shutdown signal: " + sig.String())
			exitCode.Store(128)

			cancel()

			return errTerminated

		case <-ctx.Done():
			s.log.Debug("Finalizing the Service")

			return errCanceled
		}
	}
}

func (s *Service) executeRunner(ctx context.Context, runner Runner) error {
	var panicRecovered any

	err := func() error {
		defer logtools.CatchPanic(s.log, func(recovered any) {
			panicRecovered = recovered
		})

		if err := runner.Run(ctx, s.ctn); err != nil {
			return fmt.Errorf("failed to execute runner: %w", err)
		}

		return nil
	}()
	if err != nil {
		return err
	}

	if panicRecovered != nil {
		return fmt.Errorf("runner panicked: %v", panicRecovered)
	}

	return nil
}

// Close finalizes the Service.
func (s *Service) close() {
	if s.ctn != nil {
		if err := s.ctn.Close(); err != nil {
			s.log.Warn("close container: " + err.Error())
		}
	}

	_ = s.log.Sync()
}
