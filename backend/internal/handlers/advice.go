package handlers

import (
	"net/http"
	as "sentimenta/internal/adviceService"
	"sentimenta/internal/utils"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type AdviceHandler struct {
	service as.AdviceService
	logger  *zap.SugaredLogger
}

func (h *AdviceHandler) GetAdvice(c echo.Context) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		h.logger.Errorf("Ошибка. Требуется аутентификация: %v", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "требуется аутентификация"})
	}

	dateStr := c.QueryParam("date")
	if dateStr != "" {
		layout := "2006-01-02"
		date, err := time.Parse(layout, dateStr)
		if err != nil {
			h.logger.Errorf("Ошибка при попытке получить Advice по date: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "не удалось получить Advice"})
		}
		advice, err := h.service.GetAdvice(userID, date)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "не удалось получить Advice"})
		}
		return c.JSON(http.StatusOK, advice)
	}
	advices, err := h.service.GetAdvices(userID)
	if err != nil {
		h.logger.Errorf("Ошибка при попытке получить Advices: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "не удалось получить Advices"})
	}
	return c.JSON(http.StatusOK, advices)

}

func NewAdviceHandler(service as.AdviceService, logger *zap.SugaredLogger) *AdviceHandler {
	return &AdviceHandler{service: service, logger: logger}
}
