package usecase

import (
	"context"
	"soccer_manager_service/internal/dto"
	"soccer_manager_service/internal/entity"
	"soccer_manager_service/internal/ports"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type TeamService struct {
	teamRepository      ports.TeamRepository
	playerRepository    ports.PlayerRepository
	teamCacheRepository ports.TeamCacheRepository
	logger              *zap.Logger
}

type TeamServiceParams struct {
	TeamRepository      ports.TeamRepository
	PlayerRepository    ports.PlayerRepository
	TeamCacheRepository ports.TeamCacheRepository
	Logger              *zap.Logger
}

func NewTeamService(params TeamServiceParams) *TeamService {
	return &TeamService{
		teamRepository:      params.TeamRepository,
		playerRepository:    params.PlayerRepository,
		teamCacheRepository: params.TeamCacheRepository,
		logger:              params.Logger.With(zap.String("service", "TeamService")),
	}
}

func (s *TeamService) GetMyTeam(ctx context.Context, userID uuid.UUID) (*dto.TeamWithPlayersResponse, error) {
	s.logger.Info("getting team", zap.String("user_id", userID.String()))

	cachedTeam, err := s.teamCacheRepository.GetTeam(ctx, userID)
	if err != nil {
		s.logger.Warn("failed to get cached team", zap.Error(err))
	}

	if cachedTeam != nil {
		s.logger.Debug("team retrieved from cache", zap.String("user_id", userID.String()))

		return cachedTeam, nil
	}

	team, err := s.teamRepository.GetByUserID(ctx, userID)
	if err != nil {
		s.logger.Error("failed to get team", zap.Error(err))

		return nil, err
	}

	players, err := s.playerRepository.GetByTeamID(ctx, team.ID)
	if err != nil {
		s.logger.Error("failed to get players", zap.Error(err))

		return nil, err
	}

	var totalValue int64

	for _, p := range players {
		totalValue += p.MarketValue
	}

	if team.TotalValue != totalValue {
		if err := s.teamRepository.UpdateTotalValue(ctx, team.ID, totalValue); err != nil {
			s.logger.Warn("failed to update team total value", zap.Error(err))
		} else {
			team.TotalValue = totalValue
		}
	}

	teamWithPlayers := &dto.TeamWithPlayersResponse{
		Team:    *team,
		Players: players,
	}

	if err := s.teamCacheRepository.SetTeam(ctx, userID, teamWithPlayers); err != nil {
		s.logger.Warn("failed to cache team", zap.Error(err))
	}

	return teamWithPlayers, nil
}

func (s *TeamService) UpdateTeam(ctx context.Context, userID uuid.UUID, req *dto.UpdateTeamRequest) (*entity.Team, error) {
	s.logger.Info("updating team", zap.String("user_id", userID.String()))

	existingTeam, err := s.teamRepository.GetByUserID(ctx, userID)
	if err != nil {
		s.logger.Error("failed to get team", zap.Error(err))

		return nil, err
	}

	team, err := s.teamRepository.Update(ctx, existingTeam.ID, req.Name, req.Country)
	if err != nil {
		s.logger.Error("failed to update team", zap.Error(err))

		return nil, err
	}

	if err := s.teamCacheRepository.InvalidateTeam(ctx, userID); err != nil {
		s.logger.Warn("failed to invalidate team cache", zap.Error(err))
	}

	return team, nil
}
