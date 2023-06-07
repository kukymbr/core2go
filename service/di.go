package service

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kukymbr/core2go/di"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Default dependencies names
const (
	diKeyPrefix = "core2go_"

	// DIKeyBaseContext is a key for base service context
	DIKeyBaseContext = diKeyPrefix + "base_context"

	// DIKeyLogger is a key for logger instance
	DIKeyLogger = diKeyPrefix + "logger"

	// DIKeyConfigRAW is key for non-structured config in viper object
	DIKeyConfigRAW = diKeyPrefix + "config_raw"

	// DIKeyConfig is a key for initialized Config instance
	DIKeyConfig = diKeyPrefix + "config"

	// DIKeyRouter is a key for gin router (gin.Engine) instance
	DIKeyRouter = diKeyPrefix + "router"

	// DIKeyRunner is a key for service runner instance
	DIKeyRunner = diKeyPrefix + "runner"

	// DIKeyRedis is a key for default redis instance
	DIKeyRedis = diKeyPrefix + "redis"

	// DIKeyPostgresPgx is a key for default postgres pgx connection instance
	DIKeyPostgresPgx = diKeyPrefix + "postgres_pgx"
)

const (
	DefaultConfigFileName = "app.yml"
	DefaultConfigFileType = "yml"
)

// GetDefaultDIBuilder returns default DI builder
func GetDefaultDIBuilder() (*di.Builder, error) {
	builder := &di.Builder{}

	err := builder.Add(
		DIDefBaseContext(),
		DIDefConfigRAW(DefaultConfigFileName, DefaultConfigFileType),
		DIDefConfig(),
		DIDefLogger(),
	)
	if err != nil {
		return nil, fmt.Errorf("add di definitions: %w", err)
	}

	return builder, nil
}

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

// DIGetBaseContext returns base service context
func DIGetBaseContext(ctn *di.Container) ContextWithCancel {
	return ctn.Get(DIKeyBaseContext).(ContextWithCancel)
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

// DIGetRedis returns Redis client from the DI container
func DIGetRedis(ctn *di.Container) *redis.Client {
	return ctn.Get(DIKeyRedis).(*redis.Client)
}

// DIGetPostgresPgx returns postgres pgx connection from the DI container
func DIGetPostgresPgx(ctn *di.Container) *pgxpool.Pool {
	return ctn.Get(DIKeyPostgresPgx).(*pgxpool.Pool)
}

// endregion DEFAULT GETTERS
