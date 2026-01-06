package bootstrap

import (
	"soccer_manager_service/internal/api/rest"
	"soccer_manager_service/internal/config"
	"soccer_manager_service/internal/repository"
	"soccer_manager_service/internal/usecase"
	"time"

	"go.uber.org/fx"
)

func NewApp() *fx.App {
	return fx.New(
		fx.Provide(
			config.GetConfig,
			newLogger,
			initBreakers,
			newRedis,
			newPostgres,
			newJWTManager,
			newI18nManager,
			repository.NewRepository,
			usecase.NewUsecase,
			rest.NewServer,
		),

		fx.Invoke(
			runMigrations,
			startHTTPServer,
			errWrapInit,
		),

		fx.StartTimeout(30*time.Second),
		fx.StopTimeout(10*time.Second),
	)
}
