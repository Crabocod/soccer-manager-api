package ports

import (
	"context"
	"soccer_manager_service/internal/dto"
	"soccer_manager_service/internal/entity"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, email, passwordHash string) (*entity.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
}

type TeamRepository interface {
	Create(ctx context.Context, userID uuid.UUID, name, country string, budget int64) (*entity.Team, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Team, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.Team, error)
	Update(ctx context.Context, id uuid.UUID, name, country string) (*entity.Team, error)
	UpdateBudget(ctx context.Context, id uuid.UUID, budget int64) error
	UpdateTotalValue(ctx context.Context, id uuid.UUID, totalValue int64) error
}

type PlayerRepository interface {
	Create(ctx context.Context, teamID uuid.UUID, firstName, lastName, country string, age int, position entity.PlayerPosition, marketValue int64) (*entity.Player, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Player, error)
	GetByTeamID(ctx context.Context, teamID uuid.UUID) ([]entity.Player, error)
	Update(ctx context.Context, id uuid.UUID, firstName, lastName, country string) (*entity.Player, error)
	UpdateMarketValue(ctx context.Context, id uuid.UUID, marketValue int64) error
	TransferPlayer(ctx context.Context, playerID, newTeamID uuid.UUID) error
}

type TransferRepository interface {
	Create(ctx context.Context, playerID, sellerID uuid.UUID, askingPrice int64) (*entity.Transfer, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Transfer, error)
	GetActiveTransfers(ctx context.Context) ([]entity.Transfer, error)
	Complete(ctx context.Context, id, buyerID uuid.UUID) error
	Cancel(ctx context.Context, id uuid.UUID) error
	GetByPlayerID(ctx context.Context, playerID uuid.UUID) (*entity.Transfer, error)
}

type LoginAttemptRepository interface {
	Increment(ctx context.Context, email string) (count int, err error)
	Get(ctx context.Context, email string) (attempts int, err error)
	Reset(ctx context.Context, email string) (err error)
}

type TeamCacheRepository interface {
	SetTeam(ctx context.Context, userID uuid.UUID, team *dto.TeamWithPlayersResponse) (err error)
	GetTeam(ctx context.Context, userID uuid.UUID) (team *dto.TeamWithPlayersResponse, err error)
	InvalidateTeam(ctx context.Context, userID uuid.UUID) (err error)
}
