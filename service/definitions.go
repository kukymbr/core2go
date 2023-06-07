package service

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kukymbr/core2go/di"
	ginsrv "github.com/kukymbr/core2go/ginrouter"
	"github.com/kukymbr/core2go/ginrouter/middlewares"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// DIDefBaseContext is a base service context
func DIDefBaseContext() di.Def {
	return di.Def{
		Name: DIKeyBaseContext,
		Build: func(ctn *di.Container) (any, error) {
			return NewBaseContext(context.Background()), nil
		},
	}
}

// DIDefConfigRAW is a non-structured config in viper object
func DIDefConfigRAW(file string, fileType string) di.Def {
	return di.Def{
		Name: DIKeyConfigRAW,
		Build: func(ctn *di.Container) (any, error) {
			return ReadConfigFromFile(file, fileType)
		},
	}
}

// DIDefConfig returns service Config definition
func DIDefConfig() di.Def {
	return di.Def{
		Name: DIKeyConfig,
		Build: func(ctn *di.Container) (any, error) {
			raw := DIGetConfigRAW(ctn)

			return UnmarshalConfig(raw)
		},
	}
}

// DIDefLogger returns Logger definition
func DIDefLogger() di.Def {
	return di.Def{
		Name: DIKeyLogger,
		Build: func(ctn *di.Container) (any, error) {
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
		Build: func(ctn *di.Container) (any, error) {
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
		Build: func(ctn *di.Container) (any, error) {
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

// DIDefRedis returns default redis.Client dependency definition.
func DIDefRedis() di.Def {
	return di.Def{
		Name: DIKeyRedis,
		Build: func(ctn *di.Container) (any, error) {
			baseCtx := DIGetBaseContext(ctn)
			conf := DIGetConfig(ctn)

			client := redis.NewClient(&redis.Options{
				Addr:     conf.Redis.Host,
				Username: conf.Redis.Username,
				Password: conf.Redis.Password,
				DB:       conf.Redis.DB,
			})

			//nolint:gomnd
			ctx, cancel := context.WithTimeout(baseCtx.GetContext(), 3*time.Second)
			defer cancel()

			if err := client.Ping(ctx).Err(); err != nil {
				return nil, fmt.Errorf("ping redis %s: %w", conf.Redis.Host, err)
			}

			return client, nil
		},
		Close: func(obj any) error {
			client, ok := obj.(*redis.Client)
			if ok && client != nil {
				err := client.Close()

				return fmt.Errorf("close redis: %w", err)
			}

			return nil
		},
	}
}

// DIDefPostgresPgx returns default pgx connection dependency definition.
func DIDefPostgresPgx() di.Def {
	return di.Def{
		Name: DIKeyPostgresPgx,
		Build: func(ctn *di.Container) (any, error) {
			conf := DIGetConfig(ctn)

			pgpool, err := pgxpool.New(context.Background(), conf.Postgres.AsURL().String())
			if err != nil {
				return nil, fmt.Errorf("%s dependency build: %w", DIKeyPostgresPgx, err)
			}

			return pgpool, nil
		},
		Close: func(obj any) error {
			pgpool, ok := obj.(*pgxpool.Pool)
			if ok && pgpool != nil {
				pgpool.Close()
			}

			return nil
		},
	}
}
