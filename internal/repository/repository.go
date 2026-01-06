package repository

import (
	"soccer_manager_service/internal/config"
	"soccer_manager_service/internal/ports"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Params struct {
	fx.In

	Logger   *zap.Logger
	Postgres *pgxpool.Pool
	Redis    *redis.Client
	Config   *config.Config
}

type Repository struct {
	User         ports.UserRepository
	Team         ports.TeamRepository
	Player       ports.PlayerRepository
	Transfer     ports.TransferRepository
	LoginAttempt ports.LoginAttemptRepository
	TeamCache    ports.TeamCacheRepository
}

func NewRepository(deps Params) *Repository {
	f := newRepositoryFactory(deps)

	return &Repository{
		User:         f.CreateUserRepository(),
		Team:         f.CreateTeamRepository(),
		Player:       f.CreatePlayerRepository(),
		Transfer:     f.CreateTransferRepository(),
		LoginAttempt: f.CreateLoginAttemptRepository(),
		TeamCache:    f.CreateTeamCacheRepository(),
	}
}
