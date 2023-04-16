package middleware_infrastructure

import (
	"time"

	"github.com/labstack/echo/v4/middleware"
)

type TimeoutOption func(*middleware.TimeoutConfig)

func WithTimeout(t int) func(*middleware.TimeoutConfig) {
	return func(c *middleware.TimeoutConfig) {
		c.Timeout = time.Duration(t) * time.Second
	}
}

func WithErrorMessage(msg string) func(*middleware.TimeoutConfig) {
	return func(c *middleware.TimeoutConfig) {
		c.ErrorMessage = msg
	}
}

func NewTimeout(opts ...TimeoutOption) middleware.TimeoutConfig {
	cfg := middleware.TimeoutConfig{
		Skipper: middleware.DefaultSkipper,
	}

	// Apply options
	for _, opt := range opts {
		opt(&cfg)
	}

	return cfg
}
