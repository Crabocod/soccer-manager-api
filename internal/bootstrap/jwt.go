package bootstrap

import (
	"soccer_manager_service/internal/config"
	"soccer_manager_service/pkg/jwt"
)

func newJWTManager(config *config.Config) *jwt.Manager {
	return jwt.NewManager(
		config.JWT.Secret,
		config.JWT.AccessTokenTTL,
		config.JWT.RefreshTokenTTL,
	)
}
