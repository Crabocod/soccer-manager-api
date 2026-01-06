package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"soccer_manager_service/internal/entity"
	apperr "soccer_manager_service/pkg/errors"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type MockTransferRepository struct {
	mock.Mock
}

func (m *MockTransferRepository) Create(ctx context.Context, playerID, sellerID uuid.UUID, askingPrice int64) (*entity.Transfer, error) {
	args := m.Called(ctx, playerID, sellerID, askingPrice)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*entity.Transfer), args.Error(1)
}

func (m *MockTransferRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Transfer, error) {
	args := m.Called(ctx, id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*entity.Transfer), args.Error(1)
}

func (m *MockTransferRepository) GetActiveTransfers(ctx context.Context) ([]entity.Transfer, error) {
	args := m.Called(ctx)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]entity.Transfer), args.Error(1)
}

func (m *MockTransferRepository) Complete(ctx context.Context, id, buyerID uuid.UUID) error {
	args := m.Called(ctx, id, buyerID)

	return args.Error(0)
}

func (m *MockTransferRepository) Cancel(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)

	return args.Error(0)
}

func (m *MockTransferRepository) GetByPlayerID(ctx context.Context, playerID uuid.UUID) (*entity.Transfer, error) {
	args := m.Called(ctx, playerID)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*entity.Transfer), args.Error(1)
}

type MockPlayerRepository struct {
	mock.Mock
}

func (m *MockPlayerRepository) Create(ctx context.Context, teamID uuid.UUID, firstName, lastName, country string, age int, position entity.PlayerPosition, marketValue int64) (*entity.Player, error) {
	args := m.Called(ctx, teamID, firstName, lastName, country, age, position, marketValue)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*entity.Player), args.Error(1)
}

func (m *MockPlayerRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Player, error) {
	args := m.Called(ctx, id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*entity.Player), args.Error(1)
}

func (m *MockPlayerRepository) GetByTeamID(ctx context.Context, teamID uuid.UUID) ([]entity.Player, error) {
	args := m.Called(ctx, teamID)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]entity.Player), args.Error(1)
}

func (m *MockPlayerRepository) Update(ctx context.Context, id uuid.UUID, firstName, lastName, country string) (*entity.Player, error) {
	args := m.Called(ctx, id, firstName, lastName, country)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*entity.Player), args.Error(1)
}

func (m *MockPlayerRepository) UpdateMarketValue(ctx context.Context, id uuid.UUID, marketValue int64) error {
	args := m.Called(ctx, id, marketValue)

	return args.Error(0)
}

func (m *MockPlayerRepository) TransferPlayer(ctx context.Context, playerID, newTeamID uuid.UUID) error {
	args := m.Called(ctx, playerID, newTeamID)

	return args.Error(0)
}

type MockTeamRepository struct {
	mock.Mock
}

func (m *MockTeamRepository) Create(ctx context.Context, userID uuid.UUID, name, country string, budget int64) (*entity.Team, error) {
	args := m.Called(ctx, userID, name, country, budget)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*entity.Team), args.Error(1)
}

func (m *MockTeamRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Team, error) {
	args := m.Called(ctx, id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*entity.Team), args.Error(1)
}

func (m *MockTeamRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.Team, error) {
	args := m.Called(ctx, userID)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*entity.Team), args.Error(1)
}

func (m *MockTeamRepository) Update(ctx context.Context, id uuid.UUID, name, country string) (*entity.Team, error) {
	args := m.Called(ctx, id, name, country)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*entity.Team), args.Error(1)
}

func (m *MockTeamRepository) UpdateBudget(ctx context.Context, id uuid.UUID, budget int64) error {
	args := m.Called(ctx, id, budget)

	return args.Error(0)
}

func (m *MockTeamRepository) UpdateTotalValue(ctx context.Context, id uuid.UUID, totalValue int64) error {
	args := m.Called(ctx, id, totalValue)

	return args.Error(0)
}

type MockTeamCacheRepository struct {
	mock.Mock
}

func (m *MockTeamCacheRepository) SetTeam(ctx context.Context, userID uuid.UUID, team *entity.TeamWithPlayers) error {
	args := m.Called(ctx, userID, team)

	return args.Error(0)
}

func (m *MockTeamCacheRepository) GetTeam(ctx context.Context, userID uuid.UUID) (*entity.TeamWithPlayers, error) {
	args := m.Called(ctx, userID)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*entity.TeamWithPlayers), args.Error(1)
}

func (m *MockTeamCacheRepository) InvalidateTeam(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)

	return args.Error(0)
}

func TestTransferService_ListPlayer(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()
	
	userID := uuid.New()
	playerID := uuid.New()
	teamID := uuid.New()

	t.Run("success", func(t *testing.T) {
		mockTransferRepo := new(MockTransferRepository)
		mockPlayerRepo := new(MockPlayerRepository)
		mockTeamRepo := new(MockTeamRepository)
		mockCacheRepo := new(MockTeamCacheRepository)

		player := &entity.Player{
			ID:     playerID,
			TeamID: teamID,
		}

		team := &entity.Team{
			ID:     teamID,
			UserID: userID,
		}

		transfer := &entity.Transfer{
			ID:          uuid.New(),
			PlayerID:    playerID,
			SellerID:    teamID,
			AskingPrice: 1000000,
			Status:      entity.TransferStatusActive,
		}

		mockPlayerRepo.On("GetByID", ctx, playerID).Return(player, nil)
		mockTeamRepo.On("GetByUserID", ctx, userID).Return(team, nil)
		mockTransferRepo.On("GetByPlayerID", ctx, playerID).Return(nil, apperr.ErrTransferNotFound)
		mockTransferRepo.On("Create", ctx, playerID, teamID, int64(1000000)).Return(transfer, nil)

		service := NewTransferService(TransferServiceParams{
			TransferRepository:  mockTransferRepo,
			PlayerRepository:    mockPlayerRepo,
			TeamRepository:      mockTeamRepo,
			TeamCacheRepository: mockCacheRepo,
			Logger:              logger,
		})

		req := &entity.ListPlayerRequest{AskingPrice: 1000000}
		result, err := service.ListPlayer(ctx, userID, playerID, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, transfer.ID, result.ID)
		mockPlayerRepo.AssertExpectations(t)
		mockTeamRepo.AssertExpectations(t)
		mockTransferRepo.AssertExpectations(t)
	})

	t.Run("player not found", func(t *testing.T) {
		mockTransferRepo := new(MockTransferRepository)
		mockPlayerRepo := new(MockPlayerRepository)
		mockTeamRepo := new(MockTeamRepository)
		mockCacheRepo := new(MockTeamCacheRepository)

		mockPlayerRepo.On("GetByID", ctx, playerID).Return(nil, apperr.ErrPlayerNotFound)

		service := NewTransferService(TransferServiceParams{
			TransferRepository:  mockTransferRepo,
			PlayerRepository:    mockPlayerRepo,
			TeamRepository:      mockTeamRepo,
			TeamCacheRepository: mockCacheRepo,
			Logger:              logger,
		})

		req := &entity.ListPlayerRequest{AskingPrice: 1000000}
		result, err := service.ListPlayer(ctx, userID, playerID, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, apperr.ErrPlayerNotFound, err)
		mockPlayerRepo.AssertExpectations(t)
	})

	t.Run("player not owned by user", func(t *testing.T) {
		mockTransferRepo := new(MockTransferRepository)
		mockPlayerRepo := new(MockPlayerRepository)
		mockTeamRepo := new(MockTeamRepository)
		mockCacheRepo := new(MockTeamCacheRepository)

		differentTeamID := uuid.New()
		player := &entity.Player{
			ID:     playerID,
			TeamID: differentTeamID,
		}

		team := &entity.Team{
			ID:     teamID,
			UserID: userID,
		}

		mockPlayerRepo.On("GetByID", ctx, playerID).Return(player, nil)
		mockTeamRepo.On("GetByUserID", ctx, userID).Return(team, nil)

		service := NewTransferService(TransferServiceParams{
			TransferRepository:  mockTransferRepo,
			PlayerRepository:    mockPlayerRepo,
			TeamRepository:      mockTeamRepo,
			TeamCacheRepository: mockCacheRepo,
			Logger:              logger,
		})

		req := &entity.ListPlayerRequest{AskingPrice: 1000000}
		result, err := service.ListPlayer(ctx, userID, playerID, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, apperr.ErrForbidden, err)
		mockPlayerRepo.AssertExpectations(t)
		mockTeamRepo.AssertExpectations(t)
	})

	t.Run("player already listed", func(t *testing.T) {
		mockTransferRepo := new(MockTransferRepository)
		mockPlayerRepo := new(MockPlayerRepository)
		mockTeamRepo := new(MockTeamRepository)
		mockCacheRepo := new(MockTeamCacheRepository)

		player := &entity.Player{
			ID:     playerID,
			TeamID: teamID,
		}

		team := &entity.Team{
			ID:     teamID,
			UserID: userID,
		}

		existingTransfer := &entity.Transfer{
			ID:       uuid.New(),
			PlayerID: playerID,
			Status:   entity.TransferStatusActive,
		}

		mockPlayerRepo.On("GetByID", ctx, playerID).Return(player, nil)
		mockTeamRepo.On("GetByUserID", ctx, userID).Return(team, nil)
		mockTransferRepo.On("GetByPlayerID", ctx, playerID).Return(existingTransfer, nil)

		service := NewTransferService(TransferServiceParams{
			TransferRepository:  mockTransferRepo,
			PlayerRepository:    mockPlayerRepo,
			TeamRepository:      mockTeamRepo,
			TeamCacheRepository: mockCacheRepo,
			Logger:              logger,
		})

		req := &entity.ListPlayerRequest{AskingPrice: 1000000}
		result, err := service.ListPlayer(ctx, userID, playerID, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, apperr.ErrPlayerAlreadyListed, err)
		mockPlayerRepo.AssertExpectations(t)
		mockTeamRepo.AssertExpectations(t)
		mockTransferRepo.AssertExpectations(t)
	})
}

func TestTransferService_GetTransferList(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockTransferRepo := new(MockTransferRepository)
		mockPlayerRepo := new(MockPlayerRepository)
		mockTeamRepo := new(MockTeamRepository)
		mockCacheRepo := new(MockTeamCacheRepository)

		playerID := uuid.New()
		teamID := uuid.New()

		transfers := []entity.Transfer{
			{
				ID:          uuid.New(),
				PlayerID:    playerID,
				SellerID:    teamID,
				AskingPrice: 1000000,
				Status:      entity.TransferStatusActive,
			},
		}

		player := &entity.Player{
			ID:        playerID,
			FirstName: "John",
			LastName:  "Doe",
		}

		team := &entity.Team{
			ID:   teamID,
			Name: "Test Team",
		}

		mockTransferRepo.On("GetActiveTransfers", ctx).Return(transfers, nil)
		mockPlayerRepo.On("GetByID", ctx, playerID).Return(player, nil)
		mockTeamRepo.On("GetByID", ctx, teamID).Return(team, nil)

		service := NewTransferService(TransferServiceParams{
			TransferRepository:  mockTransferRepo,
			PlayerRepository:    mockPlayerRepo,
			TeamRepository:      mockTeamRepo,
			TeamCacheRepository: mockCacheRepo,
			Logger:              logger,
		})

		result, err := service.GetTransferList(ctx)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, playerID, result[0].Player.ID)
		assert.Equal(t, teamID, result[0].Team.ID)
		mockTransferRepo.AssertExpectations(t)
		mockPlayerRepo.AssertExpectations(t)
		mockTeamRepo.AssertExpectations(t)
	})

	t.Run("empty list", func(t *testing.T) {
		mockTransferRepo := new(MockTransferRepository)
		mockPlayerRepo := new(MockPlayerRepository)
		mockTeamRepo := new(MockTeamRepository)
		mockCacheRepo := new(MockTeamCacheRepository)

		mockTransferRepo.On("GetActiveTransfers", ctx).Return([]entity.Transfer{}, nil)

		service := NewTransferService(TransferServiceParams{
			TransferRepository:  mockTransferRepo,
			PlayerRepository:    mockPlayerRepo,
			TeamRepository:      mockTeamRepo,
			TeamCacheRepository: mockCacheRepo,
			Logger:              logger,
		})

		result, err := service.GetTransferList(ctx)

		assert.NoError(t, err)
		assert.Empty(t, result)
		mockTransferRepo.AssertExpectations(t)
	})
}

func TestTransferService_BuyPlayer(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	userID := uuid.New()
	sellerUserID := uuid.New()
	transferID := uuid.New()
	playerID := uuid.New()
	buyerTeamID := uuid.New()
	sellerTeamID := uuid.New()

	t.Run("success", func(t *testing.T) {
		mockTransferRepo := new(MockTransferRepository)
		mockPlayerRepo := new(MockPlayerRepository)
		mockTeamRepo := new(MockTeamRepository)
		mockCacheRepo := new(MockTeamCacheRepository)

		transfer := &entity.Transfer{
			ID:          transferID,
			PlayerID:    playerID,
			SellerID:    sellerTeamID,
			AskingPrice: 1000000,
			Status:      entity.TransferStatusActive,
		}

		buyerTeam := &entity.Team{
			ID:     buyerTeamID,
			UserID: userID,
			Budget: 5000000,
		}

		sellerTeam := &entity.Team{
			ID:     sellerTeamID,
			UserID: sellerUserID,
			Budget: 3000000,
		}

		player := &entity.Player{
			ID:          playerID,
			TeamID:      sellerTeamID,
			MarketValue: 1000000,
		}

		mockTransferRepo.On("GetByID", ctx, transferID).Return(transfer, nil)
		mockTeamRepo.On("GetByUserID", ctx, userID).Return(buyerTeam, nil)
		mockTeamRepo.On("GetByID", ctx, sellerTeamID).Return(sellerTeam, nil)
		mockPlayerRepo.On("GetByID", ctx, playerID).Return(player, nil)
		mockPlayerRepo.On("TransferPlayer", ctx, playerID, buyerTeamID).Return(nil)
		mockPlayerRepo.On("UpdateMarketValue", ctx, playerID, mock.AnythingOfType("int64")).Return(nil)
		mockTeamRepo.On("UpdateBudget", ctx, buyerTeamID, int64(4000000)).Return(nil)
		mockTeamRepo.On("UpdateBudget", ctx, sellerTeamID, int64(4000000)).Return(nil)
		mockTransferRepo.On("Complete", ctx, transferID, buyerTeamID).Return(nil)
		mockCacheRepo.On("InvalidateTeam", ctx, userID).Return(nil)
		mockCacheRepo.On("InvalidateTeam", ctx, sellerUserID).Return(nil)

		service := NewTransferService(TransferServiceParams{
			TransferRepository:  mockTransferRepo,
			PlayerRepository:    mockPlayerRepo,
			TeamRepository:      mockTeamRepo,
			TeamCacheRepository: mockCacheRepo,
			Logger:              logger,
		})

		err := service.BuyPlayer(ctx, userID, transferID)

		assert.NoError(t, err)
		mockTransferRepo.AssertExpectations(t)
		mockPlayerRepo.AssertExpectations(t)
		mockTeamRepo.AssertExpectations(t)
		mockCacheRepo.AssertExpectations(t)
	})

	t.Run("transfer not found", func(t *testing.T) {
		mockTransferRepo := new(MockTransferRepository)
		mockPlayerRepo := new(MockPlayerRepository)
		mockTeamRepo := new(MockTeamRepository)
		mockCacheRepo := new(MockTeamCacheRepository)

		mockTransferRepo.On("GetByID", ctx, transferID).Return(nil, apperr.ErrTransferNotFound)

		service := NewTransferService(TransferServiceParams{
			TransferRepository:  mockTransferRepo,
			PlayerRepository:    mockPlayerRepo,
			TeamRepository:      mockTeamRepo,
			TeamCacheRepository: mockCacheRepo,
			Logger:              logger,
		})

		err := service.BuyPlayer(ctx, userID, transferID)

		assert.Error(t, err)
		assert.Equal(t, apperr.ErrTransferNotFound, err)
		mockTransferRepo.AssertExpectations(t)
	})

	t.Run("transfer not active", func(t *testing.T) {
		mockTransferRepo := new(MockTransferRepository)
		mockPlayerRepo := new(MockPlayerRepository)
		mockTeamRepo := new(MockTeamRepository)
		mockCacheRepo := new(MockTeamCacheRepository)

		completedTime := time.Now()
		transfer := &entity.Transfer{
			ID:          transferID,
			PlayerID:    playerID,
			SellerID:    sellerTeamID,
			AskingPrice: 1000000,
			Status:      entity.TransferStatusCompleted,
			CompletedAt: &completedTime,
		}

		mockTransferRepo.On("GetByID", ctx, transferID).Return(transfer, nil)

		service := NewTransferService(TransferServiceParams{
			TransferRepository:  mockTransferRepo,
			PlayerRepository:    mockPlayerRepo,
			TeamRepository:      mockTeamRepo,
			TeamCacheRepository: mockCacheRepo,
			Logger:              logger,
		})

		err := service.BuyPlayer(ctx, userID, transferID)

		assert.Error(t, err)
		assert.Equal(t, apperr.ErrTransferNotActive, err)
		mockTransferRepo.AssertExpectations(t)
	})

	t.Run("cannot buy own player", func(t *testing.T) {
		mockTransferRepo := new(MockTransferRepository)
		mockPlayerRepo := new(MockPlayerRepository)
		mockTeamRepo := new(MockTeamRepository)
		mockCacheRepo := new(MockTeamCacheRepository)

		transfer := &entity.Transfer{
			ID:          transferID,
			PlayerID:    playerID,
			SellerID:    buyerTeamID,
			AskingPrice: 1000000,
			Status:      entity.TransferStatusActive,
		}

		buyerTeam := &entity.Team{
			ID:     buyerTeamID,
			UserID: userID,
			Budget: 5000000,
		}

		mockTransferRepo.On("GetByID", ctx, transferID).Return(transfer, nil)
		mockTeamRepo.On("GetByUserID", ctx, userID).Return(buyerTeam, nil)

		service := NewTransferService(TransferServiceParams{
			TransferRepository:  mockTransferRepo,
			PlayerRepository:    mockPlayerRepo,
			TeamRepository:      mockTeamRepo,
			TeamCacheRepository: mockCacheRepo,
			Logger:              logger,
		})

		err := service.BuyPlayer(ctx, userID, transferID)

		assert.Error(t, err)
		assert.Equal(t, apperr.ErrCannotBuyOwnPlayer, err)
		mockTransferRepo.AssertExpectations(t)
		mockTeamRepo.AssertExpectations(t)
	})

	t.Run("insufficient funds", func(t *testing.T) {
		mockTransferRepo := new(MockTransferRepository)
		mockPlayerRepo := new(MockPlayerRepository)
		mockTeamRepo := new(MockTeamRepository)
		mockCacheRepo := new(MockTeamCacheRepository)

		transfer := &entity.Transfer{
			ID:          transferID,
			PlayerID:    playerID,
			SellerID:    sellerTeamID,
			AskingPrice: 10000000,
			Status:      entity.TransferStatusActive,
		}

		buyerTeam := &entity.Team{
			ID:     buyerTeamID,
			UserID: userID,
			Budget: 1000000,
		}

		mockTransferRepo.On("GetByID", ctx, transferID).Return(transfer, nil)
		mockTeamRepo.On("GetByUserID", ctx, userID).Return(buyerTeam, nil)

		service := NewTransferService(TransferServiceParams{
			TransferRepository:  mockTransferRepo,
			PlayerRepository:    mockPlayerRepo,
			TeamRepository:      mockTeamRepo,
			TeamCacheRepository: mockCacheRepo,
			Logger:              logger,
		})

		err := service.BuyPlayer(ctx, userID, transferID)

		assert.Error(t, err)
		assert.Equal(t, apperr.ErrInsufficientFunds, err)
		mockTransferRepo.AssertExpectations(t)
		mockTeamRepo.AssertExpectations(t)
	})

	t.Run("player transfer fails", func(t *testing.T) {
		mockTransferRepo := new(MockTransferRepository)
		mockPlayerRepo := new(MockPlayerRepository)
		mockTeamRepo := new(MockTeamRepository)
		mockCacheRepo := new(MockTeamCacheRepository)

		transfer := &entity.Transfer{
			ID:          transferID,
			PlayerID:    playerID,
			SellerID:    sellerTeamID,
			AskingPrice: 1000000,
			Status:      entity.TransferStatusActive,
		}

		buyerTeam := &entity.Team{
			ID:     buyerTeamID,
			UserID: userID,
			Budget: 5000000,
		}

		sellerTeam := &entity.Team{
			ID:     sellerTeamID,
			UserID: sellerUserID,
			Budget: 3000000,
		}

		player := &entity.Player{
			ID:          playerID,
			TeamID:      sellerTeamID,
			MarketValue: 1000000,
		}

		mockTransferRepo.On("GetByID", ctx, transferID).Return(transfer, nil)
		mockTeamRepo.On("GetByUserID", ctx, userID).Return(buyerTeam, nil)
		mockTeamRepo.On("GetByID", ctx, sellerTeamID).Return(sellerTeam, nil)
		mockPlayerRepo.On("GetByID", ctx, playerID).Return(player, nil)
		mockPlayerRepo.On("TransferPlayer", ctx, playerID, buyerTeamID).Return(errors.New("transfer failed"))

		service := NewTransferService(TransferServiceParams{
			TransferRepository:  mockTransferRepo,
			PlayerRepository:    mockPlayerRepo,
			TeamRepository:      mockTeamRepo,
			TeamCacheRepository: mockCacheRepo,
			Logger:              logger,
		})

		err := service.BuyPlayer(ctx, userID, transferID)

		assert.Error(t, err)
		mockTransferRepo.AssertExpectations(t)
		mockPlayerRepo.AssertExpectations(t)
		mockTeamRepo.AssertExpectations(t)
	})
}
