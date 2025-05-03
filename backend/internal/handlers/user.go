package handlers

import (
	"net/http"
	us "sentimenta/internal/userService"
	"sentimenta/internal/utils"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type UserHandler struct {
	service us.UserService
	logger  *zap.SugaredLogger
}

func (h *UserHandler) GetUser(c echo.Context) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		h.logger.Errorf("Ошибка. Требуется аутентификация: %v", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "требуется аутентификация"})
	}

	user, err := h.service.GetUser(userID)
	if err != nil {
		h.logger.Errorf("Ошибка при получении пользователя: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "не удалось получить данные пользователя"})
	}

	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) PatchUpdateUser(c echo.Context) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		h.logger.Errorf("Ошибка. Требуется аутентификация: %v", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "требуется аутентификация"})
	}

	var reqUser us.UserUpdate
	if err := c.Bind(&reqUser); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "неверная форма данных"})
	}

	user, err := h.service.UpdateUser(userID, reqUser)
	if err != nil {
		h.logger.Errorf("Ошибка при обновлении пользователя: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "не удалось обновить данные пользователя"})
	}

	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) PutUpdatePasswordUser(c echo.Context) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		h.logger.Errorf("Ошибка. Требуется аутентификация: %v", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "требуется аутентификация"})
	}

	var reqUser us.UserChangePass
	if err := c.Bind(&reqUser); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "неверная форма данных"})
	}

	if err := h.service.ChangePassword(userID, reqUser.Password, reqUser.NewPassword); err != nil {
		h.logger.Errorf("Ошибка при смене пароля: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Не удалось сменить пароль"})
	}
	return nil
}

func NewUserHandler(s us.UserService, logger *zap.SugaredLogger) *UserHandler {
	return &UserHandler{service: s, logger: logger}
}
