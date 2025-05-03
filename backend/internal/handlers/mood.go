package handlers

import (
	"net/http"
	c "sentimenta/internal/config"
	ms "sentimenta/internal/moodService"
	"sentimenta/internal/utils"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type MoodHandler struct {
	service ms.MoodService
	config  c.Config
	logger  *zap.SugaredLogger
}

func (h *MoodHandler) PostAddMood(c echo.Context) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		h.logger.Errorf("Ошибка. Требуется аутентификация: %v", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "требуется аутентификация"})
	}

	var reqMood ms.MoodAdd
	if err := c.Bind(&reqMood); err != nil {
		h.logger.Errorf("Ошибка при привязке данных MoodAdd: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "неверная форма данных"})
	}

	mood, err := h.service.CreateMood(userID, reqMood.Score, reqMood.Emotions, reqMood.Description, reqMood.Date)
	if err != nil {
		h.logger.Errorf("Ошибка при создании mood (%v) у пользователя с id: %v: %v", mood, userID, err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "не удалось создать mood"})
	}

	return c.JSON(http.StatusOK, mood)
}

func (h *MoodHandler) GetMoods(c echo.Context) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		h.logger.Errorf("Ошибка. Требуется аутентификация: %v", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "требуется аутентификация"})
	}

	moods, err := h.service.GetMoods(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "не удалось получить moods"})
	}

	return c.JSON(http.StatusOK, moods)

}

func (h *MoodHandler) PutUpdateMood(c echo.Context) error {
	var reqMood ms.MoodUpdate
	if err := c.Bind(&reqMood); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "неверная форма данных"})
	}

	mood := ms.Mood{
		Uid:         reqMood.Uid,
		Score:       *reqMood.Score,
		Emotions:    *reqMood.Emotions,
		Description: *reqMood.Description,
	}

	if err := h.service.UpdateMood(&mood); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "не удалось обновить mood"})
	}

	return c.JSON(http.StatusOK, mood)
}

func NewMoodHandler(s ms.MoodService, cfg c.Config, logger *zap.SugaredLogger) *MoodHandler {
	return &MoodHandler{service: s, config: cfg, logger: logger}
}
