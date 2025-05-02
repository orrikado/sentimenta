package handlers

import (
	ms "sentimenta/internal/moodService"

	"github.com/labstack/echo/v4"
)

type MoodHandler struct {
	service ms.MoodService
}

func (h *MoodHandler) PostAddMood(c echo.Context) error {

}

func (h *MoodHandler) GetMoods(c echo.Context) error {

}

func (h *MoodHandler) PutUpdateMood(c echo.Context) error {

}

func NewMoodHandler(s ms.MoodService) *MoodHandler {
	return &MoodHandler{service: s}
}
