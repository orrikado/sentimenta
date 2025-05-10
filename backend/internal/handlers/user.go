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

type UserHandler struct {
	service service.UserService
	logger  *zap.SugaredLogger
	config  *c.Config
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

	var reqUser models.UserUpdateReq
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

	var reqUser models.UserChangePass
	if err := c.Bind(&reqUser); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "неверная форма данных"})
	}

	if len([]rune(reqUser.Password)) < h.config.PASSWORD_LENGTH_MIN {
		h.logger.Infof("Регистрация отклонена: длина пароля меньше нужного")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "длина пароля меньше нужного"})
	}

	if err := h.service.ChangePassword(userID, reqUser.Password, reqUser.NewPassword); err != nil {
		h.logger.Errorf("Ошибка при смене пароля: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Не удалось сменить пароль"})
	}
	return nil
}

func NewUserHandler(s service.UserService, config *c.Config, logger *zap.SugaredLogger) *UserHandler {
	return &UserHandler{service: s, logger: logger, config: config}
}
