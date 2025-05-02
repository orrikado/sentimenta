package handlers

import (
	"net/http"
	c "sentimenta/internal/config"
	"sentimenta/internal/jwt"
	us "sentimenta/internal/userService"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	service us.UserService
	config  c.Config
}

func (h *AuthHandler) Register(c echo.Context) error {
	var newUser us.UserRegister
	if err := c.Bind(&newUser); err != nil {
		return c.JSON(http.StatusBadRequest, "неверная форма данных")
	}
	result, err := h.service.CreateUser(newUser.Username, newUser.Email, newUser.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "не удалось создать пользователя")
	}

	jwtToken, err := jwt.GenerateJWT(string(result.Uid))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "не удалось сгенерировать токен")
	}

	jwt_cookie := http.Cookie{
		Name:     h.config.JWT_COOKIE_NAME,
		Value:    jwtToken,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	}

	c.SetCookie(&jwt_cookie)
	return c.JSON(http.StatusOK, result)
}

func (h *AuthHandler) Login(c echo.Context) error {
	var reqUser us.UserLogin
	if err := c.Bind(&reqUser); err != nil {
		return c.JSON(http.StatusBadRequest, "неверная форма данных")
	}

	user, err := h.service.Authenticate(reqUser.Email, reqUser.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "аутентификация не удалась")
	}

	jwtToken, err := jwt.GenerateJWT(string(user.Uid))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "не удалось сгенерировать токен")
	}

	jwt_cookie := http.Cookie{
		Name:     h.config.JWT_COOKIE_NAME,
		Value:    jwtToken,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	}

	c.SetCookie(&jwt_cookie)
	return c.JSON(http.StatusOK, user)

}

func NewAuthHandler(s us.UserService) *AuthHandler {
	return &AuthHandler{service: s}
}
