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
	m "sentimenta/internal/models"
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
	resp    *Responser
}

type OAuthCallbackRequest struct {
	Code         string `json:"code"`
	CodeVerifier string `json:"codeVerifier"`
	Timezone     string `json:"timezone"`
}

// @Summary		Register
// @Description	Create account
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			input	body		models.UserRegister	true	"credentials"
// @Success		200		{object}	models.User
// @Failure		400		{object}	errorResponse
// @Failure		404		{object}	errorResponse
// @Failure		500		{object}	errorResponse
// @Router			/api/auth/register [post]
func (h *AuthHandler) Register(c echo.Context) error {
	var newUser m.UserRegister
	if err := c.Bind(&newUser); err != nil {
		return h.resp.newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if len([]rune(newUser.Password)) < h.config.PASSWORD_LENGTH_MIN {
		return h.resp.newErrorResponse(c, http.StatusBadRequest, errs.ErrPasswordLength.Error())
	}

	result, err := h.service.CreateUser(newUser.Username, newUser.Email, &newUser.Password, newUser.Timezone)
	if err != nil {
		if errors.Is(err, errs.ErrUserAlreadyExists) {
			return h.resp.newErrorResponse(c, http.StatusConflict, err.Error())
		}
		return h.resp.newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	uidStr := fmt.Sprintf("%v", result.Uid)
	jwtToken, err := h.JWT.GenerateJWT(uidStr)
	if err != nil {
		return h.resp.newErrorResponse(c, http.StatusInternalServerError, err.Error())
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

// @Summary		Login
// @Description	SignIn with credentials
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			input	body		m.UserLogin	true	"credentials"
// @Success		200		{object}	m.User
// @Failure		401		{object}	errorResponse
// @Failure		400		{object}	errorResponse
// @Failure		500		{object}	errorResponse
// @Router			/api/auth/login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	var reqUser m.UserLogin
	if err := c.Bind(&reqUser); err != nil {
		return h.resp.newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	user, err := h.service.Authenticate(reqUser.Email, reqUser.Password)
	if err != nil {
		return h.resp.newErrorResponse(c, http.StatusUnauthorized, err.Error())
	}

	uidStr := fmt.Sprintf("%v", user.Uid)
	jwtToken, err := h.JWT.GenerateJWT(uidStr)
	if err != nil {
		return h.resp.newErrorResponse(c, http.StatusInternalServerError, err.Error())
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

// @Summary		Google
// @Description	SignIn with Google OAuth
// @Tags			OAuth
// @Accept			json
// @Produce		json
// @Param			input	body		OAuthCallbackRequest	true	"OAuth Codes & Timezone"
// @Success		200		{object}	m.User
// @Failure		400		{object}	errorResponse
// @Failure		500		{object}	errorResponse
// @Router			/api/auth/google/callback [post]
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
		return h.resp.newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Token exchange failed: %v", err))
	}

	client := h.oauth.GoogleConfig.Client(ctx, token)
	resp, err := client.Get("https://openidconnect.googleapis.com/v1/userinfo")
	if err != nil {
		return h.resp.newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Failed to get user info: %v", err))
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			h.logger.Errorf("Failed to close response body: %v", err)
		}
	}()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return h.resp.newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	email, ok := userInfo["email"].(string)
	if !ok {
		return h.resp.newErrorResponse(c, http.StatusInternalServerError, "Invalid email format")
	}

	name, ok := userInfo["name"].(string)
	if !ok {
		return h.resp.newErrorResponse(c, http.StatusInternalServerError, "Invalid name format")
	}

	user, err := h.service.CreateUser(name, email, nil, req.Timezone)

	if err != nil {
		if err == errs.ErrUserAlreadyExists {
			h.logger.Infof("Не удалось создать пользователя: %v", err)
			userInfo["just_registered"] = true
		} else {
			h.logger.Errorf("Не удалось создать пользователя: %v", err)
			return h.resp.newErrorResponse(c, http.StatusInternalServerError, "Failed to create user")
		}
	}

	uidStr := fmt.Sprintf("%v", user.Uid)
	jwtToken, err := h.JWT.GenerateJWT(uidStr)
	if err != nil {
		return h.resp.newErrorResponse(c, http.StatusInternalServerError, err.Error())
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

// @Summary		Github
// @Description	SignIn with Github OAuth
// @Tags			OAuth
// @Accept			json
// @Produce		json
// @Param			input	body		OAuthCallbackRequest	true	"OAuth Codes & Timezone"
// @Success		200		{object}	m.EmailList
// @Failure		400		{object}	errorResponse
// @Failure		500		{object}	errorResponse
// @Router			/api/auth/github/callback [post]
func (h *AuthHandler) GithubAuthCallback(c echo.Context) error {
	var req OAuthCallbackRequest
	if err := c.Bind(&req); err != nil {
		return h.resp.newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	ctx := context.Background()

	token, err := h.oauth.GithubConfig.Exchange(ctx, req.Code,
		oauth2.SetAuthURLParam("code_verifier", req.CodeVerifier),
	)
	if err != nil {
		return h.resp.newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Token exchange failed: %v", err))
	}

	client := h.oauth.GithubConfig.Client(ctx, token)
	emailResp, err := client.Get("https://api.github.com/user/emails")
	if err != nil {
		return h.resp.newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Failed to get user emails: %v", err))
	}
	defer func() {
		if err := emailResp.Body.Close(); err != nil {
			h.logger.Errorf("Failed to close response body: %v", err)
		}
	}()
	userResp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return h.resp.newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Failed to get user info: %v", err))
	}
	defer func() {
		if err := userResp.Body.Close(); err != nil {
			h.logger.Errorf("Failed to close response body: %v", err)
		}
	}()

	var emailInfo m.EmailList
	if err := json.NewDecoder(emailResp.Body).Decode(&emailInfo.Emails); err != nil {
		return h.resp.newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	var userInfo m.GithubUserInfo
	if err := json.NewDecoder(userResp.Body).Decode(&userInfo); err != nil {
		return h.resp.newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	fmt.Println(emailInfo)
	fmt.Println(userInfo)

	name := userInfo.Name
	var email string
	for _, e := range emailInfo.Emails {
		if e.Primary {
			email = e.Email
			break
		}
	}

	user, err := h.service.CreateUser(name, email, nil, req.Timezone)

	if err != nil {
		if err == errs.ErrUserAlreadyExists {
			h.logger.Infof("Не удалось создать пользователя: %v", err)
			truePtr := true
			emailInfo.JustRegistered = &truePtr
		} else {
			return h.resp.newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
	}

	uidStr := fmt.Sprintf("%v", user.Uid)
	jwtToken, err := h.JWT.GenerateJWT(uidStr)
	if err != nil {
		return h.resp.newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	jwt_cookie := http.Cookie{
		Name:     h.config.JWT_COOKIE_NAME,
		Value:    jwtToken,
		HttpOnly: false,
		Secure:   false,
		Path:     "/",
	}

	c.SetCookie(&jwt_cookie)
	return c.JSON(http.StatusOK, emailInfo)
}

func NewAuthHandler(s service.UserService, cfg *c.Config, logger *zap.SugaredLogger, oauthConfig *auth.OAuth, JWT *security.JWT, resp *Responser) *AuthHandler {
	return &AuthHandler{service: s, config: cfg, logger: logger, oauth: oauthConfig, JWT: JWT, resp: resp}
}
