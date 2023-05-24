package service

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kukymbr/core2go/di"
	ginsrv "github.com/kukymbr/core2go/ginrouter"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Default dependencies names
const (
	// DIKeyLogger is a key for logger instance
	DIKeyLogger = "core2go_logger"

	// DIKeyConfig is a key for initialized Config instance
	DIKeyConfig = "core2go_config"

	// DIKeyRouter is a key for gin router (gin.Engine) instance
	DIKeyRouter = "core2go_router"
)

const (
	defaultConfigFileName = ".env"
	defaultConfigFileType = "env"
)

// GetDefaultDIBuilder returns default DI builder
func GetDefaultDIBuilder() (*di.Builder, error) {
	builder := &di.Builder{}

	err := builder.Add(
		DIDefConfig(defaultConfigFileName, defaultConfigFileType),
		DIDefLogger(),
		DIDefRouter(),
	)
	if err != nil {
		return nil, fmt.Errorf("add di definitions: %w", err)
	}

	return builder, nil
}

// region DEFAULT DEFINITIONS

// DIDefConfig returns service Config definition
func DIDefConfig(file string, fileType string) di.Def {
	return di.Def{
		Name: DIKeyConfig,
		Build: func(ctn *di.Container) (interface{}, error) {
			return ReadConfigFromFile(file, fileType)
		},
	}
}

// DIDefLogger returns Logger definition
func DIDefLogger() di.Def {
	return di.Def{
		Name: DIKeyLogger,
		Build: func(ctn *di.Container) (interface{}, error) {
			var conf *Config

			minimalLevel := zapcore.DebugLevel

			v, err := ctn.SafeGet(DIKeyConfig)
			if err == nil {
				conf = v.(*Config)
				if !conf.Service.IsDebug {
					minimalLevel = zapcore.InfoLevel
				}
			}

			core := getZapCore(minimalLevel)

			logger := zap.New(core)

			if conf != nil {
				logger = logger.With(
					zap.Field{
						Key:    "service_name",
						String: conf.Service.Name,
					},
					zap.Field{
						Key:    "service_version",
						String: conf.Service.GetVersion().String(),
					},
				)
			}

			return logger, nil
		},
		Close: func(obj any) (err error) {
			logger := obj.(*zap.Logger)

			err = logger.Sync()
			if err != nil {
				return fmt.Errorf("sync logger: %w", err)
			}

			return nil
		},
	}
}

func getZapCore(minimalLevel zapcore.Level) zapcore.Core {
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})

	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel && lvl >= minimalLevel
	})

	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)

	jsonEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())

	return zapcore.NewTee(
		zapcore.NewCore(jsonEncoder, consoleErrors, highPriority),
		zapcore.NewCore(jsonEncoder, consoleDebugging, lowPriority),
	)
}

// DIDefRouter returns service router definition
func DIDefRouter() di.Def {
	return di.Def{
		Name: DIKeyRouter,
		Build: func(ctn *di.Container) (interface{}, error) {
			router := ginsrv.GetDefaultRouter(DIGetLogger(ctn))

			return router, nil
		},
	}
}

// endregion DEFAULT DEFINITIONS

// region DEFAULT GETTERS

// DIGetConfig returns Config from the DI container
func DIGetConfig(ctn *di.Container) *Config {
	return ctn.Get(DIKeyConfig).(*Config)
}

// DIGetLogger returns logger from the DI container
func DIGetLogger(ctn *di.Container) *zap.Logger {
	return ctn.Get(DIKeyLogger).(*zap.Logger)
}

// DIGetRouter returns gin Engine from the DI container
func DIGetRouter(ctn *di.Container) *gin.Engine {
	return ctn.Get(DIKeyRouter).(*gin.Engine)
}

// endregion DEFAULT GETTERS
