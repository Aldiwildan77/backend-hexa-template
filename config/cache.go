package config

import "fmt"

type Cache struct {
	Host     string `env:"CACHE_HOST" envDefault:"localhost"`
	Port     int    `env:"CACHE_PORT" envDefault:"6379"`
	Username string `env:"CACHE_USERNAME" envDefault:""`
	Password string `env:"CACHE_PASSWORD" envDefault:""`
	MaxIdle  int    `env:"CACHE_MAX_IDLE" envDefault:"10"`
	MaxRetry int    `env:"CACHE_MAX_RETRY" envDefault:"3"`
}

func (c Cache) GetAddress() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
