package config

import (
	"github.com/caarlos0/env/v6"
)

type Config struct {
	Host string `env:"APP_HOST" envDefault:"localhost"`
	Port int    `env:"APP_PORT" envDefault:"8080"`

	LogLevel string `env:"LOG_LEVEL"`

	DBHost     string `env:"DATABASE_HOST"`
	DBPort     int    `env:"DATABASE_PORT" envDefault:"3306"`
	DBUsername string `env:"DATABASE_USERNAME" envDefault:"root"`
	DBPassword string `env:"DATABASE_PASSWORD"`
	DBSchema   string `env:"DATABASE_SCHEMA"`
	DBDebug    bool   `env:"DATABASE_DEBUG" envDefault:"false"`
	DBDialect  string `env:"DATABASE_DIALECT"`
}

func Get() (*Config, error) {
	cfg := Config{}
	err := env.Parse(&cfg)
	return &cfg, err
}
