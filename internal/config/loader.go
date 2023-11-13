package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

const defaultShutdownTimeout = time.Second * 20

// InitConfig initializes configuration from .env file.
func InitConfig(configFile string) (Config, error) {
	vpr := viper.New()

	vpr.SetDefault("SHUTDOWN_TIMEOUT", defaultShutdownTimeout)

	vpr.SetConfigFile(configFile)
	vpr.SetConfigType("env")
	vpr.AutomaticEnv()

	var cfg Config
	if err := vpr.ReadInConfig(); err != nil {
		return cfg, fmt.Errorf("cannot read local config file: %v", err)
	}

	if err := vpr.Unmarshal(&cfg); err != nil {
		return cfg, fmt.Errorf("cannot unmarshal local config file: %v", err)
	}

	return cfg, nil
}
