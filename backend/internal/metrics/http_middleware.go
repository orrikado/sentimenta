package metrics

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func (r *Prometheus) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)
			duration := time.Since(start).Seconds()

			statusCode := c.Response().Status

			path := c.Path()
			method := c.Request().Method
			statusText := http.StatusText(statusCode)

			r.HttpRequestDuration.WithLabelValues(
				method,
				path,
				statusText,
			).Observe(duration)
			r.HttpRequestsTotal.WithLabelValues(method, path, statusText).Inc()

			return err
		}
	}
}
