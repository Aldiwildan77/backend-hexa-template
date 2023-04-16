package middleware_infrastructure

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

type Statistic struct {
	Uptime       time.Time      `json:"uptime"`
	RequestCount int            `json:"request_count"`
	Statuses     map[string]int `json:"statuses"`

	mutex sync.RWMutex
}

func NewStatistic() *Statistic {
	return &Statistic{
		Uptime:   time.Now(),
		Statuses: make(map[string]int),
	}
}

// Process is the middleware function.
func (s *Statistic) Process(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := next(c); err != nil {
			c.Error(err)
		}

		s.mutex.Lock()
		defer s.mutex.Unlock()

		s.RequestCount++
		status := strconv.Itoa(c.Response().Status)
		s.Statuses[status]++

		return nil
	}
}

// Handle is the endpoint to get stats.
func (s *Statistic) Handle(c echo.Context) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return c.JSON(http.StatusOK, s)
}

// ServerHeader middleware adds a `Server` header to the response.
func ServerHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderServer, "Echo/3.0 - Custom")
		return next(c)
	}
}
