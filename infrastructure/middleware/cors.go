package middleware_infrastructure

import (
	"net/http"

	"github.com/labstack/echo/v4/middleware"
)

type CORSOption func(*middleware.CORSConfig)

func WithAllowMethods(methods []string) func(*middleware.CORSConfig) {
	return func(c *middleware.CORSConfig) {
		c.AllowMethods = methods
	}
}

func WithAllowOrigins(origins []string) func(*middleware.CORSConfig) {
	return func(c *middleware.CORSConfig) {
		c.AllowOrigins = origins
	}
}

func NewCORS(opts ...CORSOption) middleware.CORSConfig {
	cors := middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}

	// Apply options
	for _, opt := range opts {
		opt(&cors)
	}

	return cors
}
