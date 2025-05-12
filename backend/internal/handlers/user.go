package handlers

import (
	"net/http"
	c "sentimenta/internal/config"
	errs "sentimenta/internal/errors"
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
	resp    *Responser
}

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

	userGet := models.UserGet{
		Uid:       user.Uid,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	return c.JSON(http.StatusOK, userGet)
}

func (h *UserHandler) PatchUpdateUser(c echo.Context) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		h.logger.Errorf("Ошибка. Требуется аутентификация: %v", err)
		return h.resp.newErrorResponse(c, http.StatusUnauthorized, err.Error())
	}

	var reqUser models.UserUpdateReq
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

func (h *UserHandler) PutUpdatePasswordUser(c echo.Context) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		h.logger.Errorf("Ошибка. Требуется аутентификация: %v", err)
		return h.resp.newErrorResponse(c, http.StatusUnauthorized, err.Error())
	}

	var reqUser models.UserChangePass
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
	return nil
}

func NewUserHandler(s service.UserService, config *c.Config, logger *zap.SugaredLogger, resp *Responser) *UserHandler {
	return &UserHandler{service: s, logger: logger, config: config, resp: resp}
}
