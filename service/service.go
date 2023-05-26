package service

import (
	"fmt"

	"github.com/kukymbr/core2go/di"
	"go.uber.org/zap"
)

// New creates new Service instance with the specified DI container.
func New(ctn di.Container) *Service {
	return &Service{
		ctn: &ctn,
	}
}

// NewWithDefaultContainer creates new Service instance using the default container.
func NewWithDefaultContainer() (*Service, error) {
	builder, err := GetDefaultDIBuilder()
	if err != nil {
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

	logger *zap.Logger
	config *Config
	runner Runner

	initialized bool
}

// Init initializes App without starting the runner
func (s *Service) Init() error {
	if s.initialized {
		return nil
	}

	s.initialized = true

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

	s.logger.Info("starting the runner")

	err = s.runner.Run()
	if err != nil {
		return fmt.Errorf("runner run: %w", err)
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
