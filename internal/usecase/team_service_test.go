package usecase

import (
	"context"
	"errors"
	"github.com/stretchr/testify/mock"
	"testing"

	"soccer_manager_service/internal/entity"
	apperr "soccer_manager_service/pkg/errors"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestTeamService_GetMyTeam(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	userID := uuid.New()
	teamID := uuid.New()

	t.Run("success from cache", func(t *testing.T) {
		mockTeamRepo := new(MockTeamRepository)
		mockPlayerRepo := new(MockPlayerRepository)
		mockCacheRepo := new(MockTeamCacheRepository)

		cachedTeam := &entity.TeamWithPlayers{
			Team: entity.Team{
				ID:     teamID,
				UserID: userID,
				Name:   "Cached Team",
			},
			Players: []entity.Player{
				{ID: uuid.New(), TeamID: teamID},
			},
		}

		mockCacheRepo.On("GetTeam", ctx, userID).Return(cachedTeam, nil)

		service := NewTeamService(TeamServiceParams{
			TeamRepository:      mockTeamRepo,
			PlayerRepository:    mockPlayerRepo,
			TeamCacheRepository: mockCacheRepo,
			Logger:              logger,
		})

		result, err := service.GetMyTeam(ctx, userID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Cached Team", result.Team.Name)
		assert.Len(t, result.Players, 1)
		mockCacheRepo.AssertExpectations(t)
	})

	t.Run("success from database", func(t *testing.T) {
		mockTeamRepo := new(MockTeamRepository)
		mockPlayerRepo := new(MockPlayerRepository)
		mockCacheRepo := new(MockTeamCacheRepository)

		team := &entity.Team{
			ID:         teamID,
			UserID:     userID,
			Name:       "Test Team",
			TotalValue: 0,
		}

		players := []entity.Player{
			{ID: uuid.New(), TeamID: teamID, MarketValue: 1000000},
			{ID: uuid.New(), TeamID: teamID, MarketValue: 2000000},
		}

		mockCacheRepo.On("GetTeam", ctx, userID).Return(nil, errors.New("cache miss"))
		mockTeamRepo.On("GetByUserID", ctx, userID).Return(team, nil)
		mockPlayerRepo.On("GetByTeamID", ctx, teamID).Return(players, nil)
		mockTeamRepo.On("UpdateTotalValue", ctx, teamID, int64(3000000)).Return(nil)
		mockCacheRepo.On("SetTeam", ctx, userID, mock.MatchedBy(func(t *entity.TeamWithPlayers) bool {
			return t.Team.TotalValue == 3000000
		})).Return(nil)

		service := NewTeamService(TeamServiceParams{
			TeamRepository:      mockTeamRepo,
			PlayerRepository:    mockPlayerRepo,
			TeamCacheRepository: mockCacheRepo,
			Logger:              logger,
		})

		result, err := service.GetMyTeam(ctx, userID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Test Team", result.Team.Name)
		assert.Len(t, result.Players, 2)
		assert.Equal(t, int64(3000000), result.Team.TotalValue)
		mockTeamRepo.AssertExpectations(t)
		mockPlayerRepo.AssertExpectations(t)
		mockCacheRepo.AssertExpectations(t)
	})

	t.Run("team not found", func(t *testing.T) {
		mockTeamRepo := new(MockTeamRepository)
		mockPlayerRepo := new(MockPlayerRepository)
		mockCacheRepo := new(MockTeamCacheRepository)

		mockCacheRepo.On("GetTeam", ctx, userID).Return(nil, errors.New("cache miss"))
		mockTeamRepo.On("GetByUserID", ctx, userID).Return(nil, apperr.ErrTeamNotFound)

		service := NewTeamService(TeamServiceParams{
			TeamRepository:      mockTeamRepo,
			PlayerRepository:    mockPlayerRepo,
			TeamCacheRepository: mockCacheRepo,
			Logger:              logger,
		})

		result, err := service.GetMyTeam(ctx, userID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, apperr.ErrTeamNotFound, err)
		mockTeamRepo.AssertExpectations(t)
		mockCacheRepo.AssertExpectations(t)
	})

	t.Run("failed to get players", func(t *testing.T) {
		mockTeamRepo := new(MockTeamRepository)
		mockPlayerRepo := new(MockPlayerRepository)
		mockCacheRepo := new(MockTeamCacheRepository)

		team := &entity.Team{
			ID:     teamID,
			UserID: userID,
			Name:   "Test Team",
		}

		mockCacheRepo.On("GetTeam", ctx, userID).Return(nil, errors.New("cache miss"))
		mockTeamRepo.On("GetByUserID", ctx, userID).Return(team, nil)
		mockPlayerRepo.On("GetByTeamID", ctx, teamID).Return(nil, errors.New("database error"))

		service := NewTeamService(TeamServiceParams{
			TeamRepository:      mockTeamRepo,
			PlayerRepository:    mockPlayerRepo,
			TeamCacheRepository: mockCacheRepo,
			Logger:              logger,
		})

		result, err := service.GetMyTeam(ctx, userID)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockTeamRepo.AssertExpectations(t)
		mockPlayerRepo.AssertExpectations(t)
		mockCacheRepo.AssertExpectations(t)
	})

	t.Run("success with correct total value", func(t *testing.T) {
		mockTeamRepo := new(MockTeamRepository)
		mockPlayerRepo := new(MockPlayerRepository)
		mockCacheRepo := new(MockTeamCacheRepository)

		team := &entity.Team{
			ID:         teamID,
			UserID:     userID,
			Name:       "Test Team",
			TotalValue: 5000000,
		}

		players := []entity.Player{
			{ID: uuid.New(), TeamID: teamID, MarketValue: 2500000},
			{ID: uuid.New(), TeamID: teamID, MarketValue: 2500000},
		}

		mockCacheRepo.On("GetTeam", ctx, userID).Return(nil, errors.New("cache miss"))
		mockTeamRepo.On("GetByUserID", ctx, userID).Return(team, nil)
		mockPlayerRepo.On("GetByTeamID", ctx, teamID).Return(players, nil)
		mockCacheRepo.On("SetTeam", ctx, userID, mock.MatchedBy(func(t *entity.TeamWithPlayers) bool {
			return t.Team.TotalValue == 5000000
		})).Return(nil)

		service := NewTeamService(TeamServiceParams{
			TeamRepository:      mockTeamRepo,
			PlayerRepository:    mockPlayerRepo,
			TeamCacheRepository: mockCacheRepo,
			Logger:              logger,
		})

		result, err := service.GetMyTeam(ctx, userID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(5000000), result.Team.TotalValue)
		mockTeamRepo.AssertExpectations(t)
		mockPlayerRepo.AssertExpectations(t)
		mockCacheRepo.AssertExpectations(t)
	})
}

func TestTeamService_UpdateTeam(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	userID := uuid.New()
	teamID := uuid.New()

	t.Run("success", func(t *testing.T) {
		mockTeamRepo := new(MockTeamRepository)
		mockPlayerRepo := new(MockPlayerRepository)
		mockCacheRepo := new(MockTeamCacheRepository)

		existingTeam := &entity.Team{
			ID:     teamID,
			UserID: userID,
			Name:   "Old Team",
		}

		updatedTeam := &entity.Team{
			ID:     teamID,
			UserID: userID,
			Name:   "New Team",
		}

		mockTeamRepo.On("GetByUserID", ctx, userID).Return(existingTeam, nil)
		mockTeamRepo.On("Update", ctx, teamID, "New Team", "Spain").Return(updatedTeam, nil)
		mockCacheRepo.On("InvalidateTeam", ctx, userID).Return(nil)

		service := NewTeamService(TeamServiceParams{
			TeamRepository:      mockTeamRepo,
			PlayerRepository:    mockPlayerRepo,
			TeamCacheRepository: mockCacheRepo,
			Logger:              logger,
		})

		req := &entity.UpdateTeamRequest{
			Name:    "New Team",
			Country: "Spain",
		}

		result, err := service.UpdateTeam(ctx, userID, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "New Team", result.Name)
		mockTeamRepo.AssertExpectations(t)
		mockCacheRepo.AssertExpectations(t)
	})

	t.Run("team not found", func(t *testing.T) {
		mockTeamRepo := new(MockTeamRepository)
		mockPlayerRepo := new(MockPlayerRepository)
		mockCacheRepo := new(MockTeamCacheRepository)

		mockTeamRepo.On("GetByUserID", ctx, userID).Return(nil, apperr.ErrTeamNotFound)

		service := NewTeamService(TeamServiceParams{
			TeamRepository:      mockTeamRepo,
			PlayerRepository:    mockPlayerRepo,
			TeamCacheRepository: mockCacheRepo,
			Logger:              logger,
		})

		req := &entity.UpdateTeamRequest{
			Name:    "New Team",
			Country: "Spain",
		}

		result, err := service.UpdateTeam(ctx, userID, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, apperr.ErrTeamNotFound, err)
		mockTeamRepo.AssertExpectations(t)
	})

	t.Run("update fails", func(t *testing.T) {
		mockTeamRepo := new(MockTeamRepository)
		mockPlayerRepo := new(MockPlayerRepository)
		mockCacheRepo := new(MockTeamCacheRepository)

		existingTeam := &entity.Team{
			ID:     teamID,
			UserID: userID,
			Name:   "Old Team",
		}

		mockTeamRepo.On("GetByUserID", ctx, userID).Return(existingTeam, nil)
		mockTeamRepo.On("Update", ctx, teamID, "New Team", "Spain").Return(nil, errors.New("database error"))

		service := NewTeamService(TeamServiceParams{
			TeamRepository:      mockTeamRepo,
			PlayerRepository:    mockPlayerRepo,
			TeamCacheRepository: mockCacheRepo,
			Logger:              logger,
		})

		req := &entity.UpdateTeamRequest{
			Name:    "New Team",
			Country: "Spain",
		}

		result, err := service.UpdateTeam(ctx, userID, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockTeamRepo.AssertExpectations(t)
	})

	t.Run("cache invalidation fails but update succeeds", func(t *testing.T) {
		mockTeamRepo := new(MockTeamRepository)
		mockPlayerRepo := new(MockPlayerRepository)
		mockCacheRepo := new(MockTeamCacheRepository)

		existingTeam := &entity.Team{
			ID:     teamID,
			UserID: userID,
			Name:   "Old Team",
		}

		updatedTeam := &entity.Team{
			ID:     teamID,
			UserID: userID,
			Name:   "New Team",
		}

		mockTeamRepo.On("GetByUserID", ctx, userID).Return(existingTeam, nil)
		mockTeamRepo.On("Update", ctx, teamID, "New Team", "Spain").Return(updatedTeam, nil)
		mockCacheRepo.On("InvalidateTeam", ctx, userID).Return(errors.New("cache error"))

		service := NewTeamService(TeamServiceParams{
			TeamRepository:      mockTeamRepo,
			PlayerRepository:    mockPlayerRepo,
			TeamCacheRepository: mockCacheRepo,
			Logger:              logger,
		})

		req := &entity.UpdateTeamRequest{
			Name:    "New Team",
			Country: "Spain",
		}

		result, err := service.UpdateTeam(ctx, userID, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "New Team", result.Name)
		mockTeamRepo.AssertExpectations(t)
		mockCacheRepo.AssertExpectations(t)
	})
}
