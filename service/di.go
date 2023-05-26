package service

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kukymbr/core2go/di"
	ginsrv "github.com/kukymbr/core2go/ginrouter"
	"github.com/kukymbr/core2go/ginrouter/middlewares"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Default dependencies names
const (
	diKeyPrefix = "core2go_"

	// DIKeyLogger is a key for logger instance
	DIKeyLogger = diKeyPrefix + "logger"

	// DIKeyConfigRAW is key for non-structured config in viper object
	DIKeyConfigRAW = diKeyPrefix + "config_raw"

	// DIKeyConfig is a key for initialized Config instance
	DIKeyConfig = diKeyPrefix + "config"

	// DIKeyRouter is a key for gin router (gin.Engine) instance
	DIKeyRouter = diKeyPrefix + "router"

	DIKeyRunner = diKeyPrefix + "runner"
)

const (
	DefaultConfigFileName = "app.yml"
	DefaultConfigFileType = "yml"
)

// GetDefaultDIBuilder returns default DI builder
func GetDefaultDIBuilder() (*di.Builder, error) {
	builder := &di.Builder{}

	err := builder.Add(
		DIDefConfigRAW(DefaultConfigFileName, DefaultConfigFileType),
		DIDefConfig(),
		DIDefLogger(),
		DIDefRouter(),
		DIDefRunnerGin(),
	)
	if err != nil {
		return nil, fmt.Errorf("add di definitions: %w", err)
	}

	return builder, nil
}

// region DEFAULT DEFINITIONS

// DIDefConfigRAW is a non-structured config in viper object
func DIDefConfigRAW(file string, fileType string) di.Def {
	return di.Def{
		Name: DIKeyConfigRAW,
		Build: func(ctn *di.Container) (interface{}, error) {
			return ReadConfigFromFile(file, fileType)
		},
	}
}

// DIDefConfig returns service Config definition
func DIDefConfig() di.Def {
	return di.Def{
		Name: DIKeyConfig,
		Build: func(ctn *di.Container) (interface{}, error) {
			raw := DIGetConfigRAW(ctn)

			return UnmarshalConfig(raw)
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
					zap.String("service_name", conf.Service.Name),
					zap.String("service_version", conf.Service.GetVersion().String()),
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
			mode := gin.ReleaseMode
			conf := diGetConfigNillable(ctn)

			if conf != nil && conf.Service.IsDebug {
				mode = gin.DebugMode
			}

			gin.SetMode(mode)

			router := ginsrv.GetDefaultRouter(DIGetLogger(ctn))

			router.Use(middlewares.SetDIContainer(ctn))

			return router, nil
		},
	}
}

// DIDefRunnerGin returns service gin router runner definition
func DIDefRunnerGin() di.Def {
	return di.Def{
		Name: DIKeyRunner,
		Build: func(ctn *di.Container) (interface{}, error) {
			router := DIGetRouter(ctn)
			conf := diGetConfigNillable(ctn)
			addr := ""

			if conf != nil {
				addr = fmt.Sprintf(":%d", conf.Service.Port)
			}

			return ginsrv.NewRunner(router, addr), nil
		},
	}
}

// endregion DEFAULT DEFINITIONS

// region DEFAULT GETTERS

func diGetConfigNillable(ctn *di.Container) *Config {
	v, err := ctn.SafeGet(DIKeyConfig)
	if err != nil {
		return nil
	}

	conf, ok := v.(*Config)
	if !ok {
		return nil
	}

	return conf
}

// DIGetConfigRAW returns non-structured viper config from the DI container
func DIGetConfigRAW(ctn *di.Container) *viper.Viper {
	return ctn.Get(DIKeyConfigRAW).(*viper.Viper)
}

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

// DIGetRunner returns service runner from the DI container
func DIGetRunner(ctn *di.Container) Runner {
	return ctn.Get(DIKeyRunner).(Runner)
}

// endregion DEFAULT GETTERS
