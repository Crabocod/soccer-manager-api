package usecase

import (
	"context"
	"soccer_manager_service/internal/entity"
	"soccer_manager_service/internal/ports"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type PlayerService struct {
	playerRepository    ports.PlayerRepository
	teamRepository      ports.TeamRepository
	teamCacheRepository ports.TeamCacheRepository
	logger              *zap.Logger
}

type PlayerServiceParams struct {
	PlayerRepository    ports.PlayerRepository
	TeamRepository      ports.TeamRepository
	TeamCacheRepository ports.TeamCacheRepository
	Logger              *zap.Logger
}

func NewPlayerService(params PlayerServiceParams) *PlayerService {
	return &PlayerService{
		playerRepository:    params.PlayerRepository,
		teamRepository:      params.TeamRepository,
		teamCacheRepository: params.TeamCacheRepository,
		logger:              params.Logger.With(zap.String("service", "PlayerService")),
	}
}

func (s *PlayerService) UpdatePlayer(ctx context.Context, userID, playerID uuid.UUID, req *entity.UpdatePlayerRequest) (*entity.Player, error) {
	s.logger.Info("updating player", zap.String("player_id", playerID.String()), zap.String("user_id", userID.String()))

	player, err := s.playerRepository.GetByID(ctx, playerID)
	if err != nil {
		s.logger.Error("failed to get player", zap.Error(err))

		return nil, err
	}

	updatedPlayer, err := s.playerRepository.Update(ctx, player.ID, req.FirstName, req.LastName, req.Country)
	if err != nil {
		s.logger.Error("failed to update player", zap.Error(err))

		return nil, err
	}

	if err := s.teamCacheRepository.InvalidateTeam(ctx, userID); err != nil {
		s.logger.Warn("failed to invalidate team cache", zap.Error(err))
	}

	return updatedPlayer, nil
}
