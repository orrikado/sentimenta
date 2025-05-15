package handlers

import (
	"net/http"
	c "sentimenta/internal/config"
	errs "sentimenta/internal/errors"
	m "sentimenta/internal/models"
	"sentimenta/internal/service"
	"sentimenta/internal/utils"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type UserHandler struct {
	service service.UserService
	logger  *zap.SugaredLogger
	config  *c.Config
	resp    *Responser
}

// @Summary		User profile
// @Description	Get user information
// @Tags			User
// @Accept			json
// @Produce		json
//
// @Success		200	{object}	models.User
// @Failure		401	{object}	errorResponse
// @Failure		500	{object}	errorResponse
// @Router			/api/user/get [get]
func (h *UserHandler) GetUser(c echo.Context) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		h.logger.Errorf("Ошибка. Требуется аутентификация: %v", err)
		return h.resp.newErrorResponse(c, http.StatusUnauthorized, err.Error())
	}

	user, err := h.service.GetUser(userID)
	if err != nil {
		h.logger.Errorf("Ошибка при получении пользователя: %v", err)
		return h.resp.newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	userGet := m.UserGet{
		Uid:       user.Uid,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	return c.JSON(http.StatusOK, userGet)
}

// @Summary		Update user
// @Description	Update user
// @Tags			User
// @Accept			json
// @Produce		json
//
// @Param			input	body		models.UserUpdateReq	false	"credentials"
//
// @Success		200		{object}	models.User
// @Failure		401		{object}	errorResponse
// @Failure		400		{object}	errorResponse
// @Failure		500		{object}	errorResponse
// @Router			/api/user/update [patch]
func (h *UserHandler) PatchUpdateUser(c echo.Context) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		h.logger.Errorf("Ошибка. Требуется аутентификация: %v", err)
		return h.resp.newErrorResponse(c, http.StatusUnauthorized, err.Error())
	}

	var reqUser m.UserUpdateReq
	if err := c.Bind(&reqUser); err != nil {
		return h.resp.newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	user, err := h.service.UpdateUser(userID, reqUser)
	if err != nil {
		h.logger.Errorf("Ошибка при обновлении пользователя: %v", err)
		return h.resp.newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, user)
}

// @Summary		Update user
// @Description	Update user
// @Tags			User
// @Accept			json
// @Produce		json
//
// @Param			input	body		models.UserChangePass	true	"credentials"
//
// @Success		200		{object}	models.User
// @Failure		401		{object}	errorResponse
// @Failure		400		{object}	errorResponse
// @Failure		500		{object}	errorResponse
// @Router			/api/user/update/password [patch]
func (h *UserHandler) PutUpdatePasswordUser(c echo.Context) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		h.logger.Errorf("Ошибка. Требуется аутентификация: %v", err)
		return h.resp.newErrorResponse(c, http.StatusUnauthorized, err.Error())
	}

	var reqUser m.UserChangePass
	if err := c.Bind(&reqUser); err != nil {
		return h.resp.newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if len([]rune(reqUser.Password)) < h.config.PASSWORD_LENGTH_MIN {
		h.logger.Infof("Регистрация отклонена: длина пароля меньше нужного")
		return h.resp.newErrorResponse(c, http.StatusBadRequest, errs.ErrPasswordLength.Error())
	}

	if err := h.service.ChangePassword(userID, reqUser.Password, reqUser.NewPassword); err != nil {
		h.logger.Errorf("Ошибка при смене пароля: %v", err)
		return h.resp.newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, okResponse{"password changed successfully"})
}

func NewUserHandler(s service.UserService, config *c.Config, logger *zap.SugaredLogger, resp *Responser) *UserHandler {
	return &UserHandler{service: s, logger: logger, config: config, resp: resp}
}
