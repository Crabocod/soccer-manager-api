package usecase

import (
	"soccer_manager_service/internal/config"
	"soccer_manager_service/internal/repository"
	"soccer_manager_service/internal/usecase/adapters"
	"soccer_manager_service/pkg/jwt"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Service struct {
	Auth     adapters.AuthService
	Team     adapters.TeamService
	Player   adapters.PlayerService
	Transfer adapters.TransferService
}

type Params struct {
	fx.In

	Logger     *zap.Logger
	Config     *config.Config
	Repository *repository.Repository
	JWTManager *jwt.Manager
}

func NewUsecase(params Params) *Service {
	factory := newServiceFactory(params)

	return &Service{
		Auth:     factory.CreateAuthService(),
		Team:     factory.CreateTeamService(),
		Player:   factory.CreatePlayerService(),
		Transfer: factory.CreateTransferService(),
	}
}
