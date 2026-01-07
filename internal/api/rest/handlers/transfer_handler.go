package handlers

import (
	"errors"
	"net/http"
	"soccer_manager_service/internal/api/rest/middleware"
	"soccer_manager_service/internal/dto"
	"soccer_manager_service/internal/usecase/adapters"
	apperr "soccer_manager_service/pkg/errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"go.uber.org/zap"
)

type TransferHandler struct {
	transferService adapters.TransferService
	logger          *zap.Logger
}

func NewTransferHandler(transferService adapters.TransferService, logger *zap.Logger) *TransferHandler {
	return &TransferHandler{
		transferService: transferService,
		logger:          logger.With(zap.String("handler", "TransferHandler")),
	}
}

// ListPlayer
// @Summary List player for transfer
// @Description Put player on transfer market
// @ID list-player-transfer
// @Tags transfers
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Player ID"
// @Param request body dto.ListPlayerRequest true "Transfer data"
// @Success 201 {object} entity.Transfer
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/players/{id}/transfer [post]
func (h *TransferHandler) ListPlayer(c *gin.Context) {
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

	var req dto.ListPlayerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid list player request", zap.Error(err))

		msg, _ := localizer.Localize(&i18n.LocalizeConfig{MessageID: "errors.invalid_request"})
		c.JSON(http.StatusBadRequest, gin.H{"error": msg, "details": err.Error()})

		return
	}

	transfer, err := h.transferService.ListPlayer(c.Request.Context(), userID, playerID, &req)
	if err != nil {
		h.logger.Error("failed to list player", zap.Error(err))

		if errors.Is(err, apperr.ErrPlayerNotFound) {
			msg := apperr.LocalizeError(err, localizer)
			c.JSON(http.StatusNotFound, gin.H{"error": msg})

			return
		}

		if errors.Is(err, apperr.ErrPlayerAlreadyListed) {
			msg := apperr.LocalizeError(err, localizer)
			c.JSON(http.StatusConflict, gin.H{"error": msg})

			return
		}

		if errors.Is(err, apperr.ErrForbidden) {
			msg := apperr.LocalizeError(err, localizer)
			c.JSON(http.StatusForbidden, gin.H{"error": msg})

			return
		}

		msg := apperr.LocalizeError(apperr.ErrInternal, localizer)
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})

		return
	}

	c.JSON(http.StatusCreated, transfer)
}

// GetTransferList
// @Summary Get transfer list
// @Description Get all available transfers
// @ID get-transfer-list
// @Tags transfers
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dto.TransfersResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/transfers [get]
func (h *TransferHandler) GetTransferList(c *gin.Context) {
	localizer := c.MustGet(middleware.LocalizerKey).(*i18n.Localizer)

	transfers, err := h.transferService.GetTransferList(c.Request.Context())
	if err != nil {
		h.logger.Error("failed to get transfer list", zap.Error(err))

		msg := apperr.LocalizeError(apperr.ErrInternal, localizer)
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})

		return
	}

	c.JSON(http.StatusOK, gin.H{"transfers": transfers})
}

// BuyPlayer
// @Summary Buy player
// @Description Purchase player from transfer market
// @ID buy-player
// @Tags transfers
// @Security BearerAuth
// @Produce json
// @Param id path string true "Transfer ID"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/transfers/{id}/buy [post]
func (h *TransferHandler) BuyPlayer(c *gin.Context) {
	localizer := c.MustGet(middleware.LocalizerKey).(*i18n.Localizer)

	userID, ok := middleware.GetUserID(c)
	if !ok {
		msg := apperr.LocalizeError(apperr.ErrUnauthorized, localizer)
		c.JSON(http.StatusUnauthorized, gin.H{"error": msg})

		return
	}

	transferIDStr := c.Param("id")

	transferID, err := uuid.Parse(transferIDStr)
	if err != nil {
		msg, _ := localizer.Localize(&i18n.LocalizeConfig{MessageID: "errors.invalid_transfer_id"})
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})

		return
	}

	if err := h.transferService.BuyPlayer(c.Request.Context(), userID, transferID); err != nil {
		h.logger.Error("failed to buy player", zap.Error(err))

		if errors.Is(err, apperr.ErrTransferNotFound) {
			msg := apperr.LocalizeError(err, localizer)
			c.JSON(http.StatusNotFound, gin.H{"error": msg})

			return
		}

		if errors.Is(err, apperr.ErrTransferNotActive) {
			msg := apperr.LocalizeError(err, localizer)
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})

			return
		}

		if errors.Is(err, apperr.ErrCannotBuyOwnPlayer) {
			msg := apperr.LocalizeError(err, localizer)
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})

			return
		}

		if errors.Is(err, apperr.ErrInsufficientFunds) {
			msg := apperr.LocalizeError(err, localizer)
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})

			return
		}

		msg := apperr.LocalizeError(apperr.ErrInternal, localizer)
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})

		return
	}

	successMsg, _ := localizer.Localize(&i18n.LocalizeConfig{MessageID: "success.player_purchased"})
	c.JSON(http.StatusOK, gin.H{"message": successMsg})
}
