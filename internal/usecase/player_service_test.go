package usecase

import (
	"context"
	"errors"
	"testing"

	"soccer_manager_service/internal/entity"
	apperr "soccer_manager_service/pkg/errors"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestPlayerService_UpdatePlayer(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	userID := uuid.New()
	playerID := uuid.New()

	t.Run("success", func(t *testing.T) {
		mockPlayerRepo := new(MockPlayerRepository)
		mockTeamRepo := new(MockTeamRepository)
		mockCacheRepo := new(MockTeamCacheRepository)

		player := &entity.Player{
			ID:        playerID,
			FirstName: "John",
			LastName:  "Doe",
			Country:   "USA",
		}

		updatedPlayer := &entity.Player{
			ID:        playerID,
			FirstName: "Jane",
			LastName:  "Smith",
			Country:   "UK",
		}

		mockPlayerRepo.On("GetByID", ctx, playerID).Return(player, nil)
		mockPlayerRepo.On("Update", ctx, playerID, "Jane", "Smith", "UK").Return(updatedPlayer, nil)
		mockCacheRepo.On("InvalidateTeam", ctx, userID).Return(nil)

		service := NewPlayerService(PlayerServiceParams{
			PlayerRepository:    mockPlayerRepo,
			TeamRepository:      mockTeamRepo,
			TeamCacheRepository: mockCacheRepo,
			Logger:              logger,
		})

		req := &entity.UpdatePlayerRequest{
			FirstName: "Jane",
			LastName:  "Smith",
			Country:   "UK",
		}

		result, err := service.UpdatePlayer(ctx, userID, playerID, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Jane", result.FirstName)
		assert.Equal(t, "Smith", result.LastName)
		assert.Equal(t, "UK", result.Country)
		mockPlayerRepo.AssertExpectations(t)
		mockCacheRepo.AssertExpectations(t)
	})

	t.Run("player not found", func(t *testing.T) {
		mockPlayerRepo := new(MockPlayerRepository)
		mockTeamRepo := new(MockTeamRepository)
		mockCacheRepo := new(MockTeamCacheRepository)

		mockPlayerRepo.On("GetByID", ctx, playerID).Return(nil, apperr.ErrPlayerNotFound)

		service := NewPlayerService(PlayerServiceParams{
			PlayerRepository:    mockPlayerRepo,
			TeamRepository:      mockTeamRepo,
			TeamCacheRepository: mockCacheRepo,
			Logger:              logger,
		})

		req := &entity.UpdatePlayerRequest{
			FirstName: "Jane",
			LastName:  "Smith",
			Country:   "UK",
		}

		result, err := service.UpdatePlayer(ctx, userID, playerID, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, apperr.ErrPlayerNotFound, err)
		mockPlayerRepo.AssertExpectations(t)
	})

	t.Run("update fails", func(t *testing.T) {
		mockPlayerRepo := new(MockPlayerRepository)
		mockTeamRepo := new(MockTeamRepository)
		mockCacheRepo := new(MockTeamCacheRepository)

		player := &entity.Player{
			ID:        playerID,
			FirstName: "John",
			LastName:  "Doe",
			Country:   "USA",
		}

		mockPlayerRepo.On("GetByID", ctx, playerID).Return(player, nil)
		mockPlayerRepo.On("Update", ctx, playerID, "Jane", "Smith", "UK").Return(nil, errors.New("database error"))

		service := NewPlayerService(PlayerServiceParams{
			PlayerRepository:    mockPlayerRepo,
			TeamRepository:      mockTeamRepo,
			TeamCacheRepository: mockCacheRepo,
			Logger:              logger,
		})

		req := &entity.UpdatePlayerRequest{
			FirstName: "Jane",
			LastName:  "Smith",
			Country:   "UK",
		}

		result, err := service.UpdatePlayer(ctx, userID, playerID, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockPlayerRepo.AssertExpectations(t)
	})

	t.Run("cache invalidation fails but update succeeds", func(t *testing.T) {
		mockPlayerRepo := new(MockPlayerRepository)
		mockTeamRepo := new(MockTeamRepository)
		mockCacheRepo := new(MockTeamCacheRepository)

		player := &entity.Player{
			ID:        playerID,
			FirstName: "John",
			LastName:  "Doe",
			Country:   "USA",
		}

		updatedPlayer := &entity.Player{
			ID:        playerID,
			FirstName: "Jane",
			LastName:  "Smith",
			Country:   "UK",
		}

		mockPlayerRepo.On("GetByID", ctx, playerID).Return(player, nil)
		mockPlayerRepo.On("Update", ctx, playerID, "Jane", "Smith", "UK").Return(updatedPlayer, nil)
		mockCacheRepo.On("InvalidateTeam", ctx, userID).Return(errors.New("cache error"))

		service := NewPlayerService(PlayerServiceParams{
			PlayerRepository:    mockPlayerRepo,
			TeamRepository:      mockTeamRepo,
			TeamCacheRepository: mockCacheRepo,
			Logger:              logger,
		})

		req := &entity.UpdatePlayerRequest{
			FirstName: "Jane",
			LastName:  "Smith",
			Country:   "UK",
		}

		result, err := service.UpdatePlayer(ctx, userID, playerID, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Jane", result.FirstName)
		mockPlayerRepo.AssertExpectations(t)
		mockCacheRepo.AssertExpectations(t)
	})

	t.Run("partial update - only first name", func(t *testing.T) {
		mockPlayerRepo := new(MockPlayerRepository)
		mockTeamRepo := new(MockTeamRepository)
		mockCacheRepo := new(MockTeamCacheRepository)

		player := &entity.Player{
			ID:        playerID,
			FirstName: "John",
			LastName:  "Doe",
			Country:   "USA",
		}

		updatedPlayer := &entity.Player{
			ID:        playerID,
			FirstName: "Jane",
			LastName:  "Doe",
			Country:   "USA",
		}

		mockPlayerRepo.On("GetByID", ctx, playerID).Return(player, nil)
		mockPlayerRepo.On("Update", ctx, playerID, "Jane", "", "").Return(updatedPlayer, nil)
		mockCacheRepo.On("InvalidateTeam", ctx, userID).Return(nil)

		service := NewPlayerService(PlayerServiceParams{
			PlayerRepository:    mockPlayerRepo,
			TeamRepository:      mockTeamRepo,
			TeamCacheRepository: mockCacheRepo,
			Logger:              logger,
		})

		req := &entity.UpdatePlayerRequest{
			FirstName: "Jane",
		}

		result, err := service.UpdatePlayer(ctx, userID, playerID, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Jane", result.FirstName)
		mockPlayerRepo.AssertExpectations(t)
		mockCacheRepo.AssertExpectations(t)
	})
}
