package usecase

import (
	"context"
	"math/rand"
	"soccer_manager_service/internal/dto"
	"soccer_manager_service/internal/entity"
	"soccer_manager_service/internal/ports"
	apperr "soccer_manager_service/pkg/errors"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type TransferService struct {
	transferRepository  ports.TransferRepository
	playerRepository    ports.PlayerRepository
	teamRepository      ports.TeamRepository
	teamCacheRepository ports.TeamCacheRepository
	logger              *zap.Logger
}

type TransferServiceParams struct {
	TransferRepository  ports.TransferRepository
	PlayerRepository    ports.PlayerRepository
	TeamRepository      ports.TeamRepository
	TeamCacheRepository ports.TeamCacheRepository
	Logger              *zap.Logger
}

func NewTransferService(params TransferServiceParams) *TransferService {
	return &TransferService{
		transferRepository:  params.TransferRepository,
		playerRepository:    params.PlayerRepository,
		teamRepository:      params.TeamRepository,
		teamCacheRepository: params.TeamCacheRepository,
		logger:              params.Logger.With(zap.String("service", "TransferService")),
	}
}

func (s *TransferService) ListPlayer(ctx context.Context, userID, playerID uuid.UUID, req *dto.ListPlayerRequest) (*entity.Transfer, error) {
	s.logger.Info("listing player for transfer",
		zap.String("user_id", userID.String()),
		zap.String("player_id", playerID.String()),
		zap.Int64("asking_price", req.AskingPrice))

	player, err := s.playerRepository.GetByID(ctx, playerID)
	if err != nil {
		s.logger.Error("failed to get player", zap.Error(err))

		return nil, err
	}

	team, err := s.teamRepository.GetByUserID(ctx, userID)
	if err != nil {
		s.logger.Error("failed to get team", zap.Error(err))

		return nil, err
	}

	if player.TeamID != team.ID {
		s.logger.Warn("player does not belong to user's team",
			zap.String("player_team_id", player.TeamID.String()),
			zap.String("user_team_id", team.ID.String()))

		return nil, apperr.ErrForbidden
	}

	existingTransfer, err := s.transferRepository.GetByPlayerID(ctx, playerID)
	if err == nil && existingTransfer != nil {
		s.logger.Warn("player already listed for transfer", zap.String("transfer_id", existingTransfer.ID.String()))

		return nil, apperr.ErrPlayerAlreadyListed
	}

	transfer, err := s.transferRepository.Create(ctx, playerID, team.ID, req.AskingPrice)
	if err != nil {
		s.logger.Error("failed to create transfer", zap.Error(err))

		return nil, err
	}

	s.logger.Info("player listed for transfer successfully", zap.String("transfer_id", transfer.ID.String()))

	return transfer, nil
}

func (s *TransferService) GetTransferList(ctx context.Context) ([]dto.TransferListItemResponse, error) {
	s.logger.Info("getting transfer list")

	transfers, err := s.transferRepository.GetActiveTransfers(ctx)
	if err != nil {
		s.logger.Error("failed to get active transfers", zap.Error(err))

		return nil, err
	}

	var items []dto.TransferListItemResponse

	for _, transfer := range transfers {
		player, err := s.playerRepository.GetByID(ctx, transfer.PlayerID)
		if err != nil {
			s.logger.Warn("failed to get player for transfer",
				zap.String("player_id", transfer.PlayerID.String()),
				zap.Error(err))

			continue
		}

		team, err := s.teamRepository.GetByID(ctx, transfer.SellerID)
		if err != nil {
			s.logger.Warn("failed to get team for transfer",
				zap.String("team_id", transfer.SellerID.String()),
				zap.Error(err))

			continue
		}

		items = append(items, dto.TransferListItemResponse{
			Transfer: transfer,
			Player:   *player,
			Team:     *team,
		})
	}

	return items, nil
}

func (s *TransferService) BuyPlayer(ctx context.Context, userID, transferID uuid.UUID) error {
	s.logger.Info("buying player",
		zap.String("user_id", userID.String()),
		zap.String("transfer_id", transferID.String()))

	transfer, err := s.transferRepository.GetByID(ctx, transferID)
	if err != nil {
		s.logger.Error("failed to get transfer", zap.Error(err))

		return err
	}

	if transfer.Status != entity.TransferStatusActive {
		s.logger.Warn("transfer is not active", zap.String("status", string(transfer.Status)))

		return apperr.ErrTransferNotActive
	}

	buyerTeam, err := s.teamRepository.GetByUserID(ctx, userID)
	if err != nil {
		s.logger.Error("failed to get buyer team", zap.Error(err))

		return err
	}

	if transfer.SellerID == buyerTeam.ID {
		s.logger.Warn("cannot buy own player")

		return apperr.ErrCannotBuyOwnPlayer
	}

	if buyerTeam.Budget < transfer.AskingPrice {
		s.logger.Warn("insufficient funds",
			zap.Int64("budget", buyerTeam.Budget),
			zap.Int64("asking_price", transfer.AskingPrice))

		return apperr.ErrInsufficientFunds
	}

	sellerTeam, err := s.teamRepository.GetByID(ctx, transfer.SellerID)
	if err != nil {
		s.logger.Error("failed to get seller team", zap.Error(err))

		return err
	}

	player, err := s.playerRepository.GetByID(ctx, transfer.PlayerID)
	if err != nil {
		s.logger.Error("failed to get player", zap.Error(err))

		return err
	}

	increasePercentage := 10 + rand.Intn(91)

	newMarketValue := player.MarketValue + (player.MarketValue * int64(increasePercentage) / 100)

	if err := s.playerRepository.TransferPlayer(ctx, player.ID, buyerTeam.ID); err != nil {
		s.logger.Error("failed to transfer player", zap.Error(err))

		return err
	}

	if err := s.playerRepository.UpdateMarketValue(ctx, player.ID, newMarketValue); err != nil {
		s.logger.Error("failed to update player market value", zap.Error(err))

		return err
	}

	newBuyerBudget := buyerTeam.Budget - transfer.AskingPrice

	if err := s.teamRepository.UpdateBudget(ctx, buyerTeam.ID, newBuyerBudget); err != nil {
		s.logger.Error("failed to update buyer budget", zap.Error(err))

		return err
	}

	newSellerBudget := sellerTeam.Budget + transfer.AskingPrice

	if err := s.teamRepository.UpdateBudget(ctx, sellerTeam.ID, newSellerBudget); err != nil {
		s.logger.Error("failed to update seller budget", zap.Error(err))

		return err
	}

	if err := s.transferRepository.Complete(ctx, transfer.ID, buyerTeam.ID); err != nil {
		s.logger.Error("failed to complete transfer", zap.Error(err))

		return err
	}

	if err := s.teamCacheRepository.InvalidateTeam(ctx, userID); err != nil {
		s.logger.Warn("failed to invalidate buyer team cache", zap.Error(err))
	}

	if err := s.teamCacheRepository.InvalidateTeam(ctx, sellerTeam.UserID); err != nil {
		s.logger.Warn("failed to invalidate seller team cache", zap.Error(err))
	}

	s.logger.Info("player purchased successfully",
		zap.String("player_id", player.ID.String()),
		zap.String("buyer_team", buyerTeam.Name),
		zap.String("seller_team", sellerTeam.Name),
		zap.Int64("price", transfer.AskingPrice),
		zap.Int64("new_market_value", newMarketValue))

	return nil
}
