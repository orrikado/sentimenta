package handlers

import (
	us "sentimenta/internal/userservice"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	service us.UserService
}

func (h *UserHandler) GetUser(c echo.Context) error {
	// TODO implement
	// jwt, err := c.Cookie(cfg.Settings.JWT_COOKIE_NAME)
	// if err != nil {
	// 	return c.JSON(http.StatusUnauthorized, "требуется авторизация")
	// }
	return nil
}

func (h *UserHandler) PatchUpdateUser(c echo.Context) error {
	// TODO implement
	return nil
}

func (h *UserHandler) PutUpdatePasswordUser(c echo.Context) error {
	// TODO implement
	return nil
}

func NewUserHandler(s us.UserService) *UserHandler {
	return &UserHandler{service: s}
}
