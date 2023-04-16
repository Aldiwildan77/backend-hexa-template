package middleware_infrastructure

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type RateLimiterOption func(*middleware.RateLimiterConfig)

func WithRateLimiterStore(store middleware.RateLimiterStore) func(*middleware.RateLimiterConfig) {
	return func(c *middleware.RateLimiterConfig) {
		c.Store = store
	}
}

func WithRateLimiterIdentifierExtractor(extractor func(ctx echo.Context) (string, error)) func(*middleware.RateLimiterConfig) {
	return func(c *middleware.RateLimiterConfig) {
		c.IdentifierExtractor = extractor
	}
}

func WithRateLimiterErrorHandler(handler func(ctx echo.Context, err error) error) func(*middleware.RateLimiterConfig) {
	return func(c *middleware.RateLimiterConfig) {
		c.ErrorHandler = handler
	}
}

func WithRateLimiterDenyHandler(handler func(ctx echo.Context, identifier string, err error) error) func(*middleware.RateLimiterConfig) {
	return func(c *middleware.RateLimiterConfig) {
		c.DenyHandler = handler
	}
}

func NewRateLimiter(opts ...RateLimiterOption) middleware.RateLimiterConfig {
	cfg := middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			id := ctx.RealIP()
			return id, nil
		},
		ErrorHandler: func(ctx echo.Context, err error) error {
			return ctx.JSON(http.StatusForbidden, nil)
		},
		DenyHandler: func(ctx echo.Context, identifier string, err error) error {
			return ctx.JSON(http.StatusTooManyRequests, nil)
		},
	}

	// Apply options
	for _, opt := range opts {
		opt(&cfg)
	}

	return cfg
}
