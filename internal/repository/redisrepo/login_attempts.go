package redisrepo

import (
	"context"
	"errors"
	"fmt"
	"soccer_manager_service/internal/config"
	"strconv"

	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type LoginAttempt struct {
	client *redis.Client
	config *config.Config
	logger *zap.Logger
}

type LoginAttemptParams struct {
	fx.In

	Redis  *redis.Client
	Config *config.Config
	Logger *zap.Logger
}

func NewLoginAttempt(params LoginAttemptParams) *LoginAttempt {
	return &LoginAttempt{
		client: params.Redis,
		config: params.Config,
		logger: params.Logger.With(zap.String("repository", "LoginAttempt")),
	}
}

func (r *LoginAttempt) Increment(ctx context.Context, email string) (count int, err error) {
	if email == "" {
		return 0, errors.New("empty email")
	}

	key := createLoginAttemptsKey(email)

	pipe := r.client.Pipeline()

	incrCmd := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, r.config.Login.LoginAttemptTTL)

	if _, err := pipe.Exec(ctx); err != nil {
		r.logger.Error("failed to increment login attempts", zap.Error(err), zap.String("email", email))

		return 0, fmt.Errorf("increment login attempts: %w", err)
	}

	count64, err := incrCmd.Result()
	if err != nil {
		return 0, fmt.Errorf("get incr result: %w", err)
	}

	return int(count64), nil
}

func (r *LoginAttempt) Get(ctx context.Context, email string) (attempts int, err error) {
	if email == "" {
		return 0, errors.New("empty email")
	}

	key := createLoginAttemptsKey(email)

	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, nil
		}

		r.logger.Error("failed to get login attempts", zap.Error(err), zap.String("email", email))

		return 0, fmt.Errorf("get login attempts: %w", err)
	}

	attempts, err = strconv.Atoi(val)
	if err != nil {
		return 0, fmt.Errorf("parse attempts: %w", err)
	}

	return attempts, nil
}

func (r *LoginAttempt) Reset(ctx context.Context, email string) (err error) {
	if email == "" {
		return errors.New("empty email")
	}

	key := createLoginAttemptsKey(email)

	if err := r.client.Del(ctx, key).Err(); err != nil {
		r.logger.Error("failed to reset login attempts", zap.Error(err), zap.String("email", email))

		return fmt.Errorf("reset login attempts: %w", err)
	}

	return nil
}

func createLoginAttemptsKey(email string) string {
	return fmt.Sprintf("login_attempts:%s", email)
}
