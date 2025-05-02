package handlers

import (
	"net/http"
	us "sentimenta/internal/userService"
	"sentimenta/internal/utils"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	service us.UserService
}

func (h *UserHandler) GetUser(c echo.Context) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "требуется аутентификация")
	}

	user, err := h.service.GetUser(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "не удалось получить данные пользователя")
	}

	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) PatchUpdateUser(c echo.Context) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "требуется аутентификация")
	}

	var reqUser us.UserUpdate
	if err := c.Bind(&reqUser); err != nil {
		return c.JSON(http.StatusBadRequest, "неверная форма данных")
	}

	user, err := h.service.UpdateUser(userID, reqUser)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) PutUpdatePasswordUser(c echo.Context) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "требуется аутентификация")
	}

	var reqUser us.UserChangePass
	if err := c.Bind(&reqUser); err != nil {
		return c.JSON(http.StatusBadRequest, "неверная форма данных")
	}

	if err := h.service.ChangePassword(userID, reqUser.Password, reqUser.NewPassword); err != nil {
		return c.JSON(http.StatusInternalServerError, "Не удалось сменить пароль")
	}
	return nil
}

func NewUserHandler(s us.UserService) *UserHandler {
	return &UserHandler{service: s}
}
