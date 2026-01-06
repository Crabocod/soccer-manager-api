package bootstrap

import (
	"context"
	"soccer_manager_service/internal/config"

	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func newRedis(lc fx.Lifecycle, config *config.Config, logger *zap.Logger) (*redis.Client, error) {
	logger.Info("connecting to Redis", zap.String("host", config.Redis.Host))

	client := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Address(),
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
	})

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err := client.Ping(ctx).Err(); err != nil {
				logger.Error("failed to ping Redis", zap.Error(err))

				return err
			}

			logger.Info("Redis connected successfully")

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("closing Redis connection")

			return client.Close()
		},
	})

	return client, nil
}
