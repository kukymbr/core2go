package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/kukymbr/core2go/di"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

var (
	errTerminated = errors.New("terminated")
	errCanceled   = errors.New("canceled")
)

// New creates new Service instance with the specified DI container.
func New(ctn di.Container) *Service {
	return &Service{
		ctn: &ctn,
	}
}

// NewWithDefaultContainer creates a new Service instance using the default container.
func NewWithDefaultContainer() (*Service, error) {
	builder, err := GetDefaultDIBuilder()
	if err != nil {
		return nil, err
	}

	if err := builder.Add(DIDefRouter(), DIDefRunnerGin()); err != nil {
		return nil, err
	}

	ctn, err := builder.Build()
	if err != nil {
		return nil, fmt.Errorf("build default di container: %w", err)
	}

	return New(*ctn), nil
}

// Service is an application based on core2go kit
type Service struct {
	ctn *di.Container

	baseContext ContextWithCancel
	logger      *zap.Logger
	config      *Config
	runner      Runner

	initialized bool
}

// Init initializes App without starting the runner
func (s *Service) Init() error {
	if s.initialized {
		return nil
	}

	s.initialized = true

	s.baseContext = DIGetBaseContext(s.ctn)
	s.logger = DIGetLogger(s.ctn)
	s.config = DIGetConfig(s.ctn)
	s.runner = DIGetRunner(s.ctn)

	return nil
}

// Run starts the App's server
func (s *Service) Run() error {
	err := s.Init()
	if err != nil {
		return err
	}

	defer func() {
		err := s.Close()
		if err != nil {
			s.logger.Error(err.Error())
		}
	}()

	return s.run()
}

func (s *Service) run() error {
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM)

	errGroup, ctx := errgroup.WithContext(s.baseContext.GetContext())

	errGroup.Go(func() error {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("context error: %w", err)
		}

		defer s.baseContext.GetCancelFn()()

		s.logger.Info("starting the runner")

		err := s.runner.Run()
		if err != nil {
			return fmt.Errorf("runner run: %w", err)
		}

		return nil
	})

	errGroup.Go(func() error {
		for {
			select {
			case sig := <-shutdownChan:
				s.logger.Info("got shutdown signal: " + sig.String())
				s.baseContext.GetCancelFn()()

				return errTerminated

			case <-s.baseContext.GetContext().Done():
				err := s.baseContext.GetContext().Err()
				if errors.Is(err, context.Canceled) {
					err = errCanceled
				}

				return err
			}
		}
	})

	err := errGroup.Wait()

	if err != nil && !errors.Is(err, errTerminated) && !errors.Is(err, errCanceled) {
		s.logger.Error(err.Error())

		return fmt.Errorf("error group: %w", err)
	}

	return nil
}

// Close finalizes the App
func (s *Service) Close() error {
	err := s.ctn.Close()
	if err != nil {
		return fmt.Errorf("close container: %w", err)
	}

	return nil
}
