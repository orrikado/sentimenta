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

type MoodHandler struct {
	service service.MoodService
	config  *c.Config
	logger  *zap.SugaredLogger
	resp    *Responser
}

// @Summary		Create
// @Description	Create new mood
// @Tags			Moods
// @Accept			json
// @Produce		json
//
// @Param			input	body		models.MoodAdd	true	"credentials"
//
// @Success		200		{object}	models.Mood
// @Failure		401		{object}	errorResponse
// @Failure		400		{object}	errorResponse
// @Failure		500		{object}	errorResponse
// @Router			/api/moods/add [post]
func (h *MoodHandler) PostAddMood(c echo.Context) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return h.resp.newErrorResponse(c, http.StatusUnauthorized, err.Error())
	}

	var reqMood models.MoodAdd
	if err := c.Bind(&reqMood); err != nil {
		return h.resp.newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if len([]rune(reqMood.Description)) > h.config.MOOD_DESC_LENGTH_MAX {
		return h.resp.newErrorResponse(c, http.StatusBadRequest, errs.ErrMoodDescLength.Error())
	}

	if len([]rune(reqMood.Emotions)) > h.config.MOOD_EMOTES_LENGTH_MAX {
		return h.resp.newErrorResponse(c, http.StatusBadRequest, errs.ErrMoodEmotesLength.Error())
	}

	mood, err := h.service.CreateMood(userID, reqMood.Score, reqMood.Emotions, reqMood.Description, reqMood.Date)
	if err != nil {
		return h.resp.newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, mood)
}

// @Summary		Get all
// @Description	Get all moods by user id in jwt-token
// @Tags			Moods
// @Accept			json
// @Produce			json
//
// @Success		200	{array}		models.Mood
// @Failure		401	{object}	errorResponse
// @Failure		400	{object}	errorResponse
// @Failure		500	{object}	errorResponse
// @Router			/api/moods/get [get]
func (h *MoodHandler) GetMoods(c echo.Context) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return h.resp.newErrorResponse(c, http.StatusUnauthorized, err.Error())
	}

	moods, err := h.service.GetMoods(userID)
	if err != nil {
		return h.resp.newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, moods)
}

// @Summary		Update mood
// @Description	Update something mood fields
// @Tags			Moods
// @Accept			json
// @Produce		json
//
// @Param			input	body		models.MoodUpdate	false	"credentials"
//
// @Success		200		{object}	models.Mood
// @Failure		401		{object}	errorResponse
// @Failure		400		{object}	errorResponse
// @Failure		500		{object}	errorResponse
// @Router			/api/moods/update [put]
func (h *MoodHandler) PutUpdateMood(c echo.Context) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return h.resp.newErrorResponse(c, http.StatusUnauthorized, err.Error())
	}

	var reqMood models.MoodUpdate
	if err := c.Bind(&reqMood); err != nil {
		return h.resp.newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	mood := models.Mood{
		Uid:         reqMood.Uid,
		Score:       *reqMood.Score,
		Emotions:    *reqMood.Emotions,
		Description: *reqMood.Description,
	}

	if err := h.service.UpdateMood(userID, &mood); err != nil {
		return h.resp.newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, mood)
}

func NewMoodHandler(s service.MoodService, cfg *c.Config, logger *zap.SugaredLogger, resp *Responser) *MoodHandler {
	return &MoodHandler{service: s, config: cfg, logger: logger, resp: resp}
}
