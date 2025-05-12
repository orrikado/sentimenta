package handlers

import (
	"net/http"
	c "sentimenta/internal/config"
	"sentimenta/internal/models"
	"sentimenta/internal/service"
	"sentimenta/internal/utils"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type MoodHandler struct {
	service service.MoodService
	config  *c.Config
	logger  *zap.SugaredLogger
	resp    *Responser
}

func (h *MoodHandler) PostAddMood(c echo.Context) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return h.resp.newErrorResponse(c, http.StatusUnauthorized, err.Error())
	}

	var reqMood models.MoodAdd
	if err := c.Bind(&reqMood); err != nil {
		return h.resp.newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	mood, err := h.service.CreateMood(userID, reqMood.Score, reqMood.Emotions, reqMood.Description, reqMood.Date)
	if err != nil {
		return h.resp.newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, mood)
}

func (h *MoodHandler) GetMoods(c echo.Context) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return h.resp.newErrorResponse(c, http.StatusUnauthorized, err.Error())
	}

	moods, err := h.service.GetMoods(userID)
	if err != nil {
		return h.resp.newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, moods)
}

func (h *MoodHandler) PutUpdateMood(c echo.Context) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return h.resp.newErrorResponse(c, http.StatusUnauthorized, err.Error())
	}

	var reqMood models.MoodUpdate
	if err := c.Bind(&reqMood); err != nil {
		return h.resp.newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	mood := models.Mood{
		Uid:         reqMood.Uid,
		Score:       *reqMood.Score,
		Emotions:    *reqMood.Emotions,
		Description: *reqMood.Description,
	}

	if err := h.service.UpdateMood(userID, &mood); err != nil {
		return h.resp.newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, mood)
}

func NewMoodHandler(s service.MoodService, cfg *c.Config, logger *zap.SugaredLogger, resp *Responser) *MoodHandler {
	return &MoodHandler{service: s, config: cfg, logger: logger, resp: resp}
}
