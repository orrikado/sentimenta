package handlers

import (
	"net/http"
	"sentimenta/internal/metrics"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type errorResponse struct {
	ErrorMessage string `json:"error"`
}

type okResponse struct {
	Message string `json:"message"`
}

type Responser struct {
	prometheus *metrics.Prometheus
	logger     *zap.SugaredLogger
}

func (r *Responser) newErrorResponse(c echo.Context, statusCode int, message string) error {
	r.prometheus.HttpErrorsTotal.WithLabelValues(c.Request().Method, c.Path(), http.StatusText(statusCode)).Inc()
	return c.JSON(statusCode, errorResponse{ErrorMessage: message})
}

func NewResponser(prometheus *metrics.Prometheus, logger *zap.SugaredLogger) *Responser {
	return &Responser{prometheus: prometheus, logger: logger}
}
