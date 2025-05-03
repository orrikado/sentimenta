package handlers

import (
	"errors"
	"fmt"
	"net/http"
	c "sentimenta/internal/config"
	errs "sentimenta/internal/errors"
	JWT "sentimenta/internal/jwt"
	us "sentimenta/internal/userService"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type AuthHandler struct {
	service us.UserService
	config  c.Config
	logger  *zap.SugaredLogger
}

func (h *AuthHandler) Register(c echo.Context) error {
	var newUser us.UserRegister
	if err := c.Bind(&newUser); err != nil {
		h.logger.Errorf("Ошибка при Bind UserRegister: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "неверная форма данных"})
	}

	result, err := h.service.CreateUser(newUser.Username, newUser.Email, newUser.Password)
	if err != nil {
		if errors.Is(err, errs.ErrUserAlreadyExists) {
			h.logger.Infof("Регистрация отклонена: пользователь с почтой %s уже существует", newUser.Email)
			return c.JSON(http.StatusConflict, map[string]string{"error": "пользователь с такой почтой уже существует"})
		}
		h.logger.Errorf("Неизвестная ошибка при создании пользователя: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "не удалось создать пользователя"})
	}

	uidStr := fmt.Sprintf("%v", result.Uid)
	jwtToken, err := JWT.GenerateJWT(uidStr)
	if err != nil {
		h.logger.Errorf("Ошибка при генерации JWT-Токена: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "не удалось сгенерировать токен"})
	}

	jwt_cookie := http.Cookie{
		Name:     h.config.JWT_COOKIE_NAME,
		Value:    jwtToken,
		HttpOnly: false,
		Secure:   false,
		Path:     "/",
	}

	c.SetCookie(&jwt_cookie)
	return c.JSON(http.StatusCreated, result)
}

func (h *AuthHandler) Login(c echo.Context) error {
	var reqUser us.UserLogin
	if err := c.Bind(&reqUser); err != nil {
		h.logger.Errorf("Ошибка при Bind UserLogin: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "неверная форма данных"})
	}

	user, err := h.service.Authenticate(reqUser.Email, reqUser.Password)
	if err != nil {
		h.logger.Errorf("Ошибка аутентификации: %v", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "данные не верны"})
	}

	uidStr := fmt.Sprintf("%v", user.Uid)
	jwtToken, err := JWT.GenerateJWT(uidStr)
	if err != nil {
		h.logger.Errorf("Ошибка при генерации JWT-Токена: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "не удалось сгенерировать токен"})
	}

	jwt_cookie := http.Cookie{
		Name:     h.config.JWT_COOKIE_NAME,
		Value:    jwtToken,
		HttpOnly: false,
		Secure:   false,
		Path:     "/",
	}

	c.SetCookie(&jwt_cookie)
	return c.JSON(http.StatusOK, user)

}

func NewAuthHandler(s us.UserService, cfg c.Config, logger *zap.SugaredLogger) *AuthHandler {
	return &AuthHandler{service: s, config: cfg, logger: logger}
}
