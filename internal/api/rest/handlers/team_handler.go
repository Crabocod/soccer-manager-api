package handlers

import (
	"errors"
	"net/http"
	"soccer_manager_service/internal/api/rest/middleware"
	"soccer_manager_service/internal/dto"
	"soccer_manager_service/internal/usecase/adapters"
	apperr "soccer_manager_service/pkg/errors"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"go.uber.org/zap"
)

type TeamHandler struct {
	teamService adapters.TeamService
	logger      *zap.Logger
}

func NewTeamHandler(teamService adapters.TeamService, logger *zap.Logger) *TeamHandler {
	return &TeamHandler{
		teamService: teamService,
		logger:      logger.With(zap.String("handler", "TeamHandler")),
	}
}

// GetMyTeam
// @Summary Get my team
// @Description Get current user's team with players
// @ID get-my-team
// @Tags team
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dto.TeamWithPlayersResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/team [get]
func (h *TeamHandler) GetMyTeam(c *gin.Context) {
	localizer := c.MustGet(middleware.LocalizerKey).(*i18n.Localizer)

	userID, ok := middleware.GetUserID(c)
	if !ok {
		msg := apperr.LocalizeError(apperr.ErrUnauthorized, localizer)
		c.JSON(http.StatusUnauthorized, gin.H{"error": msg})

		return
	}

	teamWithPlayers, err := h.teamService.GetMyTeam(c.Request.Context(), userID)
	if err != nil {
		h.logger.Error("failed to get team", zap.Error(err))

		if errors.Is(err, apperr.ErrTeamNotFound) {
			msg := apperr.LocalizeError(err, localizer)
			c.JSON(http.StatusNotFound, gin.H{"error": msg})

			return
		}

		msg := apperr.LocalizeError(apperr.ErrInternal, localizer)
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})

		return
	}

	c.JSON(http.StatusOK, teamWithPlayers)
}

// UpdateTeam
// @Summary Update team
// @Description Update team name and country
// @ID update-team
// @Tags team
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.UpdateTeamRequest true "Team update data"
// @Success 200 {object} entity.Team
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/team [patch]
func (h *TeamHandler) UpdateTeam(c *gin.Context) {
	localizer := c.MustGet(middleware.LocalizerKey).(*i18n.Localizer)

	userID, ok := middleware.GetUserID(c)
	if !ok {
		msg := apperr.LocalizeError(apperr.ErrUnauthorized, localizer)
		c.JSON(http.StatusUnauthorized, gin.H{"error": msg})

		return
	}

	var req dto.UpdateTeamRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid update team request", zap.Error(err))

		msg, _ := localizer.Localize(&i18n.LocalizeConfig{MessageID: "errors.invalid_request"})
		c.JSON(http.StatusBadRequest, gin.H{"error": msg, "details": err.Error()})

		return
	}

	team, err := h.teamService.UpdateTeam(c.Request.Context(), userID, &req)
	if err != nil {
		h.logger.Error("failed to update team", zap.Error(err))

		if errors.Is(err, apperr.ErrTeamNotFound) {
			msg := apperr.LocalizeError(err, localizer)
			c.JSON(http.StatusNotFound, gin.H{"error": msg})

			return
		}

		msg := apperr.LocalizeError(apperr.ErrInternal, localizer)
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})

		return
	}

	c.JSON(http.StatusOK, team)
}
