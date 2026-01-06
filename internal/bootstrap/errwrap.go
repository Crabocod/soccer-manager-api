package bootstrap

import (
	"soccer_manager_service/internal/config"

	"go.uber.org/zap"
)

func errWrapInit(cfg *config.Config, logger *zap.Logger) {
	logger.Info("error wrapper initialized", zap.String("app", cfg.App.Name))
}
