package handlers

import (
	"net/http"
	"sentimenta/internal/service"
	"sentimenta/internal/utils"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type AdviceHandler struct {
	service service.AdviceService
	logger  *zap.SugaredLogger
	resp    *Responser
}

func (h *AdviceHandler) GetAdvice(c echo.Context) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		h.logger.Errorf("Ошибка. Требуется аутентификация: %v", err)
		return h.resp.newErrorResponse(c, http.StatusUnauthorized, err.Error())
	}

	dateStr := c.QueryParam("date")
	if dateStr != "" {
		layout := "2006-01-02"
		date, err := time.Parse(layout, dateStr)
		if err != nil {
			h.logger.Errorf("Ошибка при попытке получить Advice по date: %v", err)
			return h.resp.newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		advice, err := h.service.GetAdvice(userID, date)
		if err != nil {
			return h.resp.newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, advice)
	}
	advices, err := h.service.GetAdvices(userID)
	if err != nil {
		h.logger.Errorf("Ошибка при попытке получить Advices: %v", err)
		return h.resp.newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, advices)

}

func NewAdviceHandler(service service.AdviceService, logger *zap.SugaredLogger, resp *Responser) *AdviceHandler {
	return &AdviceHandler{service: service, logger: logger, resp: resp}
}
