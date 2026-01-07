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

type AuthHandler struct {
	authService adapters.AuthService
	logger      *zap.Logger
}

func NewAuthHandler(authService adapters.AuthService, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger.With(zap.String("handler", "AuthHandler")),
	}
}

// Register
// @Summary Register new user
// @Description Register new user and create team
// @ID register
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Registration data"
// @Success 201 {object} dto.TokenResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	localizer := c.MustGet(middleware.LocalizerKey).(*i18n.Localizer)

	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid register request", zap.Error(err))

		msg, _ := localizer.Localize(&i18n.LocalizeConfig{MessageID: "errors.invalid_request"})
		c.JSON(http.StatusBadRequest, gin.H{"error": msg, "details": err.Error()})

		return
	}

	accessToken, refreshToken, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("registration failed", zap.Error(err))

		if errors.Is(err, apperr.ErrUserAlreadyExists) {
			msg := apperr.LocalizeError(err, localizer)
			c.JSON(http.StatusConflict, gin.H{"error": msg})

			return
		}

		msg := apperr.LocalizeError(apperr.ErrInternal, localizer)
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// Login
// @Summary Login user
// @Description Login user with credentials
// @ID login
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login credentials"
// @Success 200 {object} dto.TokenResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 429 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	localizer := c.MustGet(middleware.LocalizerKey).(*i18n.Localizer)

	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid login request", zap.Error(err))

		msg, _ := localizer.Localize(&i18n.LocalizeConfig{MessageID: "errors.invalid_request"})
		c.JSON(http.StatusBadRequest, gin.H{"error": msg, "details": err.Error()})

		return
	}

	accessToken, refreshToken, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("login failed", zap.Error(err))

		if errors.Is(err, apperr.ErrInvalidCredentials) {
			msg := apperr.LocalizeError(err, localizer)
			c.JSON(http.StatusUnauthorized, gin.H{"error": msg})

			return
		}

		if errors.Is(err, apperr.ErrTooManyAttempts) {
			msg := apperr.LocalizeError(err, localizer)
			c.JSON(http.StatusTooManyRequests, gin.H{"error": msg})

			return
		}

		msg := apperr.LocalizeError(apperr.ErrInternal, localizer)
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
