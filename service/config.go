package service

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/go-version"
	"github.com/spf13/viper"
)

type Strings map[string]string

// Config is a Service's configuration struct.
type Config struct {
	Service  CommonServiceConfig `mapstructure:"service"`
	Redis    RedisConfig         `mapstructure:"redis"`
	Postgres PostgresConfig      `mapstructure:"postgres"`
}

// CommonServiceConfig is a common service options
type CommonServiceConfig struct {
	// Name is a service name.
	Name string `mapstructure:"name" example:"service-name"`

	// Title is a human-friendly service title.
	Title string `mapstructure:"title" example:"Service Name"`

	// Version is a service version in SemVer format.
	Version string `mapstructure:"version" example:"0.1.1-beta"`

	// IsDebug mode enables debug flag in gin and etc.
	IsDebug bool `mapstructure:"is_debug"`

	// Port is a port number to listen to.
	Port int `mapstructure:"port"`

	version *version.Version
}

// GetVersionSafe returns parsed service's version.
func (c *CommonServiceConfig) GetVersionSafe() (v *version.Version, err error) {
	if c.version != nil {
		return c.version, nil
	}

	c.version, err = version.NewVersion(c.Version)
	if err != nil {
		return nil, fmt.Errorf("parse version: %w", err)
	}

	return c.version, nil
}

// GetVersion returns parsed service's version and panics in case of error.
func (c *CommonServiceConfig) GetVersion() *version.Version {
	v, err := c.GetVersionSafe()
	if err != nil {
		panic(err.Error())
	}

	return v
}

// RedisConfig is a redis connection config
type RedisConfig struct {
	Host     string `mapstructure:"host" example:"127.0.0.1:6379"`
	Username string `mapstructure:"username" example:"redis_user"`
	Password string `mapstructure:"password" example:"redis_password"`
	DB       int    `mapstructure:"db" example:"0"`
}

// PostgresConfig is a Postgres connection config
type PostgresConfig struct {
	Host     string  `mapstructure:"host" example:"127.0.0.1:5432"`
	Username string  `mapstructure:"username" example:"postgres"`
	Password string  `mapstructure:"password" example:"postgres_password"`
	DB       string  `mapstructure:"db" example:"dbname"`
	Options  Strings `mapstructure:"options" example:"{sslmode:verify-ca, pool_max_conns:10}"`
}

// AsURL returns PostgresConfig as an url.URL object
func (c *PostgresConfig) AsURL() *url.URL {
	query := &url.Values{}

	for key, val := range c.Options {
		query.Add(key, val)
	}

	return &url.URL{
		Scheme:   "postgres",
		Host:     c.Host,
		Path:     "/" + c.DB,
		RawPath:  "/" + url.QueryEscape(c.DB),
		User:     url.UserPassword(c.Username, c.Password),
		RawQuery: query.Encode(),
	}
}

// ReadConfigFromFile reads configuration file values to the new Config instance
func ReadConfigFromFile(file string, fileType string) (vpr *viper.Viper, err error) {
	vpr = viper.New()

	vpr.SetConfigFile(file)
	vpr.SetConfigType(fileType)

	err = vpr.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s (type=%s): %w", file, fileType, err)
	}

	return vpr, nil
}

// UnmarshalConfig unmarshal viper to Config
func UnmarshalConfig(raw *viper.Viper) (conf *Config, err error) {
	err = raw.Unmarshal(&conf)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return conf, nil
}
