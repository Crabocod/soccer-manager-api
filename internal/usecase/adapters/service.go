package adapters

import (
	"context"
	"soccer_manager_service/internal/entity"

	"github.com/google/uuid"
)

type AuthService interface {
	Register(ctx context.Context, req *entity.RegisterRequest) (accessToken, refreshToken string, err error)
	Login(ctx context.Context, req *entity.LoginRequest) (accessToken, refreshToken string, err error)
}

type TeamService interface {
	GetMyTeam(ctx context.Context, userID uuid.UUID) (*entity.TeamWithPlayers, error)
	UpdateTeam(ctx context.Context, userID uuid.UUID, req *entity.UpdateTeamRequest) (*entity.Team, error)
}

type PlayerService interface {
	UpdatePlayer(ctx context.Context, userID, playerID uuid.UUID, req *entity.UpdatePlayerRequest) (*entity.Player, error)
}

type TransferService interface {
	ListPlayer(ctx context.Context, userID, playerID uuid.UUID, req *entity.ListPlayerRequest) (*entity.Transfer, error)
	GetTransferList(ctx context.Context) ([]entity.TransferListItem, error)
	BuyPlayer(ctx context.Context, userID, transferID uuid.UUID) error
}
