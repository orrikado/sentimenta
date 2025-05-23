package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type StatusHandler struct {
}

func (h *StatusHandler) GetStatus(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func NewStatusHandler() *StatusHandler {
	return &StatusHandler{}
}
