package adapters

import (
	"context"
	"soccer_manager_service/internal/dto"
	"soccer_manager_service/internal/entity"

	"github.com/google/uuid"
)

type AuthService interface {
	Register(ctx context.Context, req *dto.RegisterRequest) (accessToken, refreshToken string, err error)
	Login(ctx context.Context, req *dto.LoginRequest) (accessToken, refreshToken string, err error)
}

type TeamService interface {
	GetMyTeam(ctx context.Context, userID uuid.UUID) (*dto.TeamWithPlayersResponse, error)
	UpdateTeam(ctx context.Context, userID uuid.UUID, req *dto.UpdateTeamRequest) (*entity.Team, error)
}

type PlayerService interface {
	UpdatePlayer(ctx context.Context, userID, playerID uuid.UUID, req *dto.UpdatePlayerRequest) (*entity.Player, error)
}

type TransferService interface {
	ListPlayer(ctx context.Context, userID, playerID uuid.UUID, req *dto.ListPlayerRequest) (*entity.Transfer, error)
	GetTransferList(ctx context.Context) ([]dto.TransferListItemResponse, error)
	BuyPlayer(ctx context.Context, userID, transferID uuid.UUID) error
}
