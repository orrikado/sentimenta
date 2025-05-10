package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sentimenta/internal/auth"
	c "sentimenta/internal/config"
	errs "sentimenta/internal/errors"
	"sentimenta/internal/models"
	"sentimenta/internal/security"
	"sentimenta/internal/service"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

type AuthHandler struct {
	service service.UserService
	config  *c.Config
	logger  *zap.SugaredLogger
	oauth   *auth.OAuth
	JWT     *security.JWT
}

type OAuthCallbackRequest struct {
	Code         string `json:"code"`
	CodeVerifier string `json:"codeVerifier"`
	Timezone     string `json:"timezone"`
}

func (h *AuthHandler) Register(c echo.Context) error {
	var newUser models.UserRegister
	if err := c.Bind(&newUser); err != nil {
		h.logger.Errorf("Ошибка при Bind UserRegister: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "неверная форма данных"})
	}

	if len([]rune(newUser.Password)) < h.config.PASSWORD_LENGTH_MIN {
		h.logger.Infof("Регистрация отклонена: длина пароля меньше нужного")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "длина пароля меньше нужного"})
	}

	result, err := h.service.CreateUser(newUser.Username, newUser.Email, &newUser.Password, newUser.Timezone)
	if err != nil {
		if errors.Is(err, errs.ErrUserAlreadyExists) {
			h.logger.Infof("Регистрация отклонена: пользователь с почтой %s уже существует", newUser.Email)
			return c.JSON(http.StatusConflict, map[string]string{"error": "пользователь с такой почтой уже существует"})
		}
		h.logger.Errorf("Неизвестная ошибка при создании пользователя: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "не удалось создать пользователя"})
	}

	uidStr := fmt.Sprintf("%v", result.Uid)
	jwtToken, err := h.JWT.GenerateJWT(uidStr)
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
	var reqUser models.UserLogin
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
	jwtToken, err := h.JWT.GenerateJWT(uidStr)
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

	token, err := h.oauth.GoogleConfig.Exchange(ctx, req.Code,
		oauth2.SetAuthURLParam("code_verifier", req.CodeVerifier),
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Token exchange failed: %v", err))
	}

	client := h.oauth.GoogleConfig.Client(ctx, token)
	resp, err := client.Get("https://openidconnect.googleapis.com/v1/userinfo")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to get user info: %v", err))
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			h.logger.Errorf("Failed to close response body: %v", err)
		}
	}()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to decode user info")
	}

	email, ok := userInfo["email"].(string)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "Invalid email format")
	}

	name, ok := userInfo["name"].(string)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "Invalid name format")
	}

	user, err := h.service.CreateUser(name, email, nil, req.Timezone)

	if err != nil {
		if err == errs.ErrUserAlreadyExists {
			h.logger.Infof("Не удалось создать пользователя: %v", err)
		} else {
			h.logger.Errorf("Не удалось создать пользователя: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create user")
		}
	}

	uidStr := fmt.Sprintf("%v", user.Uid)
	jwtToken, err := h.JWT.GenerateJWT(uidStr)
	if err != nil {
		h.logger.Errorf("Ошибка при генерации JWT-Токена: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "не удалось сгенерировать токен")
	}

	jwt_cookie := http.Cookie{
		Name:     h.config.JWT_COOKIE_NAME,
		Value:    jwtToken,
		HttpOnly: false,
		Secure:   false,
		Path:     "/",
	}

	c.SetCookie(&jwt_cookie)
	return c.JSON(http.StatusOK, userInfo)
}

func (h *AuthHandler) GithubAuthCallback(c echo.Context) error {
	var req OAuthCallbackRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid payload")
	}

	ctx := context.Background()

	token, err := h.oauth.GithubConfig.Exchange(ctx, req.Code,
		oauth2.SetAuthURLParam("code_verifier", req.CodeVerifier),
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Token exchange failed: %v", err))
	}

	client := h.oauth.GithubConfig.Client(ctx, token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to get user info: %v", err))
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			h.logger.Errorf("Failed to close response body: %v", err)
		}
	}()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to decode user info")
	}

	email, ok := userInfo["email"].(string)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "Invalid email format")
	}

	name, ok := userInfo["name"].(string)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "Invalid name format")
	}

	user, err := h.service.CreateUser(name, email, nil, req.Timezone)

	if err != nil {
		if err == errs.ErrUserAlreadyExists {
			h.logger.Infof("Не удалось создать пользователя: %v", err)
		} else {
			h.logger.Errorf("Не удалось создать пользователя: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create user")
		}
	}

	uidStr := fmt.Sprintf("%v", user.Uid)
	jwtToken, err := h.JWT.GenerateJWT(uidStr)
	if err != nil {
		h.logger.Errorf("Ошибка при генерации JWT-Токена: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "не удалось сгенерировать токен")
	}

	jwt_cookie := http.Cookie{
		Name:     h.config.JWT_COOKIE_NAME,
		Value:    jwtToken,
		HttpOnly: false,
		Secure:   false,
		Path:     "/",
	}

	c.SetCookie(&jwt_cookie)
	return c.JSON(http.StatusOK, userInfo)
}

func NewAuthHandler(s service.UserService, cfg *c.Config, logger *zap.SugaredLogger, oauthConfig *auth.OAuth, JWT *security.JWT) *AuthHandler {
	return &AuthHandler{service: s, config: cfg, logger: logger, oauth: oauthConfig, JWT: JWT}
}
