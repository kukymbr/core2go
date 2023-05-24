package service

import (
	"fmt"

	"github.com/hashicorp/go-version"
	"github.com/spf13/viper"
)

// Config is a Service's configuration struct.
type Config struct {
	Service *CommonServiceConfig `mapstructure:"service"`
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

	return nil, fmt.Errorf("parse version: %w", err)
}

// GetVersion returns parsed service's version and panics in case of error.
func (c *CommonServiceConfig) GetVersion() *version.Version {
	v, err := c.GetVersionSafe()
	if err != nil {
		panic(err)
	}

	return v
}

// ReadConfigFromFile reads configuration file values to the new Config instance
func ReadConfigFromFile(file string, fileType string) (*Config, error) {
	vpr := viper.New()

	vpr.SetConfigFile(file)
	vpr.SetConfigType(fileType)

	err := vpr.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s (type=%s): %w", file, fileType, err)
	}

	conf := &Config{}

	err = vpr.Unmarshal(&conf)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file %s (type=%s): %w", file, fileType, err)
	}

	return conf, nil
}
