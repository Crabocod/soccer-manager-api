package bootstrap

import (
	"context"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"soccer_manager_service/internal/api/rest"
	"soccer_manager_service/internal/config"
)

func startHTTPServer(lc fx.Lifecycle, server *rest.Server, config *config.Config, logger *zap.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				addr := config.Server.Address()
				logger.Info("starting HTTP server", zap.String("address", addr))

				if err := server.GetRouter().Run(addr); err != nil {
					logger.Fatal("failed to start HTTP server", zap.Error(err))
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("stopping HTTP server")

			return nil
		},
	})
}
