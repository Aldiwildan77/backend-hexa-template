package middleware_infrastructure

import (
	"fmt"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/VictoriaMetrics/metrics"
	"github.com/labstack/echo/v4"
)

var inFlightCount atomic.Int64

func init() {
	metrics.NewGauge("http_requests_in_flight", func() float64 {
		return float64(inFlightCount.Load())
	})
}

type Metrics struct{}

func NewMetrics() *Metrics {
	return &Metrics{}
}

// Process is the Echo middleware that records HTTP metrics.
func (m *Metrics) Process(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		inFlightCount.Add(1)
		defer inFlightCount.Add(-1)

		start := time.Now()

		err := next(c)
		if err != nil {
			c.Error(err)
		}

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Response().Status)
		method := c.Request().Method
		path := c.Path()

		counterName := fmt.Sprintf(`http_requests_total{method=%q, path=%q, status=%q}`, method, path, status)
		metrics.GetOrCreateCounter(counterName).Inc()

		summaryName := fmt.Sprintf(`http_request_duration_seconds{method=%q, path=%q, status=%q}`, method, path, status)
		metrics.GetOrCreateSummary(summaryName).Update(duration)

		return nil
	}
}

// Handle serves the /metrics endpoint in Prometheus exposition format.
func (m *Metrics) Handle(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, "text/plain; charset=utf-8")
	metrics.WritePrometheus(c.Response().Writer, true)
	return nil
}
