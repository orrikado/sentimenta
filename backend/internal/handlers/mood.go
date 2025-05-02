package handlers

import (
	"net/http"
	c "sentimenta/internal/config"
	ms "sentimenta/internal/moodService"
	"sentimenta/internal/utils"

	"github.com/labstack/echo/v4"
)

type MoodHandler struct {
	service ms.MoodService
	config  c.Config
}

func (h *MoodHandler) PostAddMood(c echo.Context) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "ошибка аутентификации")
	}

	var reqMood ms.MoodAdd
	if err := c.Bind(&reqMood).Error; err != nil {
		return c.JSON(http.StatusBadRequest, "неверная форма данных")
	}

	mood, err := h.service.CreateMood(userID, reqMood.Score, reqMood.Emotions, reqMood.Description)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "не удалось создать mood")
	}

	return c.JSON(http.StatusOK, mood)
}

func (h *MoodHandler) GetMoods(c echo.Context) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "ошибка аутентификации")
	}

	moods, err := h.service.GetMoods(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "не удалось получить moods")
	}

	return c.JSON(http.StatusOK, moods)

}

func (h *MoodHandler) PutUpdateMood(c echo.Context) error {
	var reqMood ms.MoodUpdate
	if err := c.Bind(&reqMood); err != nil {
		return c.JSON(http.StatusBadRequest, "неверная форма данных")
	}

	mood := ms.Mood{
		Uid:         reqMood.Uid,
		Score:       *reqMood.Score,
		Emotions:    *reqMood.Emotions,
		Description: reqMood.Description,
	}

	if err := h.service.UpdateMood(&mood); err != nil {
		return c.JSON(http.StatusInternalServerError, "не удалось обновить mood")
	}

	return c.JSON(http.StatusOK, mood)
}

func NewMoodHandler(s ms.MoodService, cfg c.Config) *MoodHandler {
	return &MoodHandler{service: s, config: cfg}
}
