package postgress

import (
	"strings"
	"time"
)

type Config struct {
	User                  string        `mapstructure:"DB_USER"`
	Password              string        `mapstructure:"DB_PASS"`
	Master                string        `mapstructure:"DB_HOST"`
	Replica               string        `mapstructure:"DB_REPLICA_HOST"`
	Port                  string        `mapstructure:"DB_PORT"`
	Name                  string        `mapstructure:"DB_NAME"`
	SSLMode               string        `mapstructure:"DB_SSL_MODE"`
	LogLevel              int           `mapstructure:"DB_LOG_LEVEL"`
	MaxOpenConnections    int           `mapstructure:"DB_MAX_OPEN_CONNECTIONS"`
	MaxIdleConnections    int           `mapstructure:"DB_MAX_IDLE_CONNECTIONS"`
	MaxConnectionLifetime time.Duration `mapstructure:"DB_MAX_CONNECTIONS_LIFETIME"`
}

func (cfg *Config) MasterDSN() string {
	return cfg.buildDSN(cfg.Master)
}

func (cfg *Config) ReplicaDSN() string {
	return cfg.buildDSN(cfg.Replica)
}

func (cfg *Config) buildDSN(host string) string {
	if host == "" {
		return ""
	}

	dbOptions := []string{
		"user=" + cfg.User,
		"host=" + host,
		"port=" + cfg.Port,
		"dbname=" + cfg.Name,
	}

	if cfg.SSLMode != "" {
		dbOptions = append(dbOptions, "sslmode="+cfg.SSLMode)
	}

	if cfg.Password != "" {
		dbOptions = append(dbOptions, "password="+cfg.Password)
	}

	return strings.Join(dbOptions, " ")
}
