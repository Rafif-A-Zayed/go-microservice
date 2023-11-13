package config

import (
	"user-management/internal/config/db/postgress"
)

type Config struct {
	Environment string           `mapstructure:"APP_ENVIRONMENT"`
	Version     string           `mapstructure:"APP_VERSION"`
	LogLevel    string           `mapstructure:"LOG_LEVEL"`
	DB          postgress.Config `mapstructure:",squash"`
	HTTP        HTTP             `mapstructure:",squash"`
}

// HTTP contains HTTP server config.
type HTTP struct {
	Port int `mapstructure:"HTTP_PORT"`
}
