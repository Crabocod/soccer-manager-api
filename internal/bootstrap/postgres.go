package bootstrap

import (
	"context"
	"soccer_manager_service/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func newPostgres(lc fx.Lifecycle, config *config.Config, logger *zap.Logger) (*pgxpool.Pool, error) {
	logger.Info("connecting to PostgreSQL", zap.String("host", config.Database.Host))

	poolConfig, err := pgxpool.ParseConfig(config.Database.DSN())
	if err != nil {
		logger.Error("failed to parse PostgreSQL config", zap.Error(err))

		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		logger.Error("failed to create PostgreSQL pool", zap.Error(err))

		return nil, err
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err := pool.Ping(ctx); err != nil {
				logger.Error("failed to ping PostgreSQL", zap.Error(err))

				return err
			}

			logger.Info("PostgreSQL connected successfully")

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("closing PostgreSQL connection")
			pool.Close()

			return nil
		},
	})

	return pool, nil
}
