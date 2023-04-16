package config

import (
	"github.com/caarlos0/env/v8"
)

type Config struct {
	Application Application
	Logger      Logger
	Database    Database
	Cache       Cache
	Middleware  Middleware
}

func New() (*Config, error) {
	cfg := Config{}
	err := env.Parse(&cfg)
	return &cfg, err
}
