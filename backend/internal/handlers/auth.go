package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	c "sentimenta/internal/config"
	errs "sentimenta/internal/errors"
	JWT "sentimenta/internal/jwt"
	us "sentimenta/internal/userService"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

type AuthHandler struct {
	service  us.UserService
	config   c.Config
	logger   *zap.SugaredLogger
	oauthCfg *oauth2.Config
}

type OAuthCallbackRequest struct {
	Code         string `json:"code"`
	CodeVerifier string `json:"codeVerifier"`
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

func (h *AuthHandler) GoogleAuthCallback(c echo.Context) error {
	var req OAuthCallbackRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid payload")
	}

	ctx := context.Background()

	token, err := h.oauthCfg.Exchange(ctx, req.Code,
		oauth2.SetAuthURLParam("code_verifier", req.CodeVerifier),
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Token exchange failed: %v", err))
	}

	client := h.oauthCfg.Client(ctx, token)
	resp, err := client.Get("https://openidconnect.googleapis.com/v1/userinfo")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to get user info: %v", err))
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to decode user info")
	}

	// 💡 Здесь ты можешь:
	// - Создать JWT
	// - Найти/создать пользователя в БД
	// - Вернуть JWT/сессию/пользователя

	return c.JSON(http.StatusOK, userInfo)
}

func NewAuthHandler(s us.UserService, cfg c.Config, logger *zap.SugaredLogger, oauthConfig *oauth2.Config) *AuthHandler {
	return &AuthHandler{service: s, config: cfg, logger: logger, oauthCfg: oauthConfig}
}
