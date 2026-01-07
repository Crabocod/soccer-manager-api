package redisrepo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"soccer_manager_service/internal/config"
	"soccer_manager_service/internal/dto"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type TeamCache struct {
	client *redis.Client
	config *config.Config
	logger *zap.Logger
}

type TeamCacheParams struct {
	fx.In

	Redis  *redis.Client
	Config *config.Config
	Logger *zap.Logger
}

func NewTeamCache(params TeamCacheParams) *TeamCache {
	return &TeamCache{
		client: params.Redis,
		config: params.Config,
		logger: params.Logger.With(zap.String("repository", "TeamCache")),
	}
}

func (r *TeamCache) SetTeam(ctx context.Context, userID uuid.UUID, team *dto.TeamWithPlayersResponse) (err error) {
	if userID == uuid.Nil {
		return errors.New("empty user_id")
	}

	if team == nil {
		return errors.New("empty team")
	}

	key := createTeamCacheKey(userID)

	data, err := json.Marshal(team)
	if err != nil {
		r.logger.Error("failed to marshal team", zap.Error(err), zap.String("user_id", userID.String()))

		return fmt.Errorf("marshal team: %w", err)
	}

	if err := r.client.Set(ctx, key, data, r.config.Login.TeamCacheTTL).Err(); err != nil {
		r.logger.Error("failed to cache team", zap.Error(err), zap.String("user_id", userID.String()))

		return fmt.Errorf("cache team: %w", err)
	}

	r.logger.Debug("team cached successfully", zap.String("user_id", userID.String()))

	return nil
}

func (r *TeamCache) GetTeam(ctx context.Context, userID uuid.UUID) (team *dto.TeamWithPlayersResponse, err error) {
	if userID == uuid.Nil {
		return nil, errors.New("empty user_id")
	}

	key := createTeamCacheKey(userID)

	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			r.logger.Debug("team not found in cache", zap.String("user_id", userID.String()))

			return nil, nil
		}

		r.logger.Error("failed to get cached team", zap.Error(err), zap.String("user_id", userID.String()))

		return nil, fmt.Errorf("get cached team: %w", err)
	}

	var teamWithPlayers dto.TeamWithPlayersResponse

	if err := json.Unmarshal([]byte(data), &teamWithPlayers); err != nil {
		r.logger.Error("failed to unmarshal team", zap.Error(err), zap.String("user_id", userID.String()))

		return nil, fmt.Errorf("unmarshal team: %w", err)
	}

	r.logger.Debug("team retrieved from cache", zap.String("user_id", userID.String()))

	return &teamWithPlayers, nil
}

func (r *TeamCache) InvalidateTeam(ctx context.Context, userID uuid.UUID) (err error) {
	if userID == uuid.Nil {
		return errors.New("empty user_id")
	}

	key := createTeamCacheKey(userID)

	if err := r.client.Del(ctx, key).Err(); err != nil {
		r.logger.Error("failed to invalidate team cache", zap.Error(err), zap.String("user_id", userID.String()))

		return fmt.Errorf("invalidate team cache: %w", err)
	}

	r.logger.Debug("team cache invalidated", zap.String("user_id", userID.String()))

	return nil
}

func createTeamCacheKey(userID uuid.UUID) string {
	return fmt.Sprintf("team_cache:%s", userID.String())
}
