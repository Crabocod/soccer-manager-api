package bootstrap

import (
	"go.uber.org/zap"
	i18nPkg "soccer_manager_service/pkg/i18n"
)

func newI18nManager(logger *zap.Logger) (*i18nPkg.Manager, error) {
	manager, err := i18nPkg.NewManager()
	if err != nil {
		logger.Error("failed to initialize i18n manager", zap.Error(err))

		return nil, err
	}

	logger.Info("i18n manager initialized")

	return manager, nil
}
