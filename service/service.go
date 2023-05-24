package service

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/kukymbr/core2go/di"
	"github.com/kukymbr/core2go/ginrouter/middlewares"
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

// PrepareRouterFn is a function to prepare router before run
type PrepareRouterFn func(router *gin.Engine, ctn *di.Container) error

// Service is an application based on core2go kit
type Service struct {
	ctn *di.Container

	logger *zap.Logger
	router *gin.Engine
	config *Config

	prepareRouterFn PrepareRouterFn

	initialized bool
}

// Init initializes App without starting the server
func (s *Service) Init() error {
	if s.initialized {
		return nil
	}

	s.initialized = true

	s.logger = DIGetLogger(s.ctn)
	s.config = DIGetConfig(s.ctn)

	s.router = DIGetRouter(s.ctn)
	s.router.Use(middlewares.SetDIContainer(s.ctn))

	if s.prepareRouterFn != nil {
		err := s.prepareRouterFn(s.router, s.ctn)
		if err != nil {
			return fmt.Errorf("failed to prepare router: %w", err)
		}
	}

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

	s.logger.Info("starting gin server")

	err = s.router.Run(fmt.Sprintf(":%d", s.config.Service.Port))
	if err != nil {
		s.logger.Error(err.Error())
	}

	return fmt.Errorf("run gin server: %w", err)
}

// Close finalizes the App
func (s *Service) Close() error {
	err := s.ctn.Close()
	if err != nil {
		return fmt.Errorf("close container: %w", err)
	}

	return nil
}

// SetPrepareRouterFn sets init router hook
func (s *Service) SetPrepareRouterFn(fn PrepareRouterFn) {
	s.prepareRouterFn = fn
}
