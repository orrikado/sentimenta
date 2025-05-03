package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	c "sentimenta/internal/config"
	errs "sentimenta/internal/errors"
	JWT "sentimenta/internal/jwt"
	us "sentimenta/internal/userService"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AuthHandler struct {
	service us.UserService
	config  c.Config
	logger  *zap.SugaredLogger
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
	var reqBody OAuthCallbackRequest
	if err := c.Bind(&reqBody); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid request body",
		})
	}

	if reqBody.Code == "" || reqBody.CodeVerifier == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Missing code or code_verifier",
		})
	}

	// Настройка OAuth2 конфигурации
	conf := &oauth2.Config{
		ClientID:     h.config.GOOGLE_CLIENT_ID,
		ClientSecret: h.config.GOOGLE_CLIENT_SECRET,
		RedirectURL:  h.config.GOOGLE_CLIENT_CALLBACK,
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint:     google.Endpoint,
	}

	// Поддержка PKCE
	ctx := context.Background()
	token, err := conf.Exchange(ctx, reqBody.Code, oauth2.SetAuthURLParam("code_verifier", reqBody.CodeVerifier))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to exchange token",
		})
	}

	// Создание безопасной куки
	cookie := &http.Cookie{
		Name:     "session",
		Value:    token.AccessToken,
		HttpOnly: false, // !!!
		Secure:   false, // !!!
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int((time.Hour * 24 * 30).Seconds()),
		Path:     "/",
	}
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Login successful",
	})
}

func NewAuthHandler(s us.UserService, cfg c.Config, logger *zap.SugaredLogger) *AuthHandler {
	return &AuthHandler{service: s, config: cfg, logger: logger}
}
