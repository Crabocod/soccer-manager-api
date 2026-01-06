package handlers

import (
	"errors"
	"net/http"
	"soccer_manager_service/internal/api/rest/middleware"
	"soccer_manager_service/internal/entity"
	"soccer_manager_service/internal/usecase/adapters"
	apperr "soccer_manager_service/pkg/errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"go.uber.org/zap"
)

type PlayerHandler struct {
	playerService adapters.PlayerService
	logger        *zap.Logger
}

func NewPlayerHandler(playerService adapters.PlayerService, logger *zap.Logger) *PlayerHandler {
	return &PlayerHandler{
		playerService: playerService,
		logger:        logger.With(zap.String("handler", "PlayerHandler")),
	}
}

// @Summary Update player
// @Description Update player information
// @ID update-player
// @Tags players
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Player ID"
// @Param request body entity.UpdatePlayerRequest true "Player update data"
// @Success 200 {object} entity.Player
// @Failure 400 {object} entity.ErrorResponse
// @Failure 401 {object} entity.ErrorResponse
// @Failure 404 {object} entity.ErrorResponse
// @Failure 500 {object} entity.ErrorResponse
// @Router /api/v1/players/{id} [patch]
func (h *PlayerHandler) UpdatePlayer(c *gin.Context) {
	localizer := c.MustGet(middleware.LocalizerKey).(*i18n.Localizer)

	userID, ok := middleware.GetUserID(c)
	if !ok {
		msg := apperr.LocalizeError(apperr.ErrUnauthorized, localizer)
		c.JSON(http.StatusUnauthorized, gin.H{"error": msg})

		return
	}

	playerIDStr := c.Param("id")

	playerID, err := uuid.Parse(playerIDStr)
	if err != nil {
		msg, _ := localizer.Localize(&i18n.LocalizeConfig{MessageID: "errors.invalid_player_id"})
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})

		return
	}

	var req entity.UpdatePlayerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid update player request", zap.Error(err))

		msg, _ := localizer.Localize(&i18n.LocalizeConfig{MessageID: "errors.invalid_request"})
		c.JSON(http.StatusBadRequest, gin.H{"error": msg, "details": err.Error()})

		return
	}

	player, err := h.playerService.UpdatePlayer(c.Request.Context(), userID, playerID, &req)
	if err != nil {
		h.logger.Error("failed to update player", zap.Error(err))

		if errors.Is(err, apperr.ErrPlayerNotFound) {
			msg := apperr.LocalizeError(err, localizer)
			c.JSON(http.StatusNotFound, gin.H{"error": msg})

			return
		}

		msg := apperr.LocalizeError(apperr.ErrInternal, localizer)
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})

		return
	}

	c.JSON(http.StatusOK, player)
}
