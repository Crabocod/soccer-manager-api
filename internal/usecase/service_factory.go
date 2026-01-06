package usecase

import (
	"soccer_manager_service/internal/usecase/adapters"
)

type serviceFactory struct {
	params Params
}

func newServiceFactory(params Params) *serviceFactory {
	return &serviceFactory{params: params}
}

func (f *serviceFactory) CreateAuthService() adapters.AuthService {
	return NewAuthService(AuthServiceParams{
		UserRepository:         f.params.Repository.User,
		TeamRepository:         f.params.Repository.Team,
		PlayerRepository:       f.params.Repository.Player,
		LoginAttemptRepository: f.params.Repository.LoginAttempt,
		JWTManager:             f.params.JWTManager,
		Logger:                 f.params.Logger,
		Config:                 f.params.Config,
	})
}

func (f *serviceFactory) CreateTeamService() adapters.TeamService {
	return NewTeamService(TeamServiceParams{
		TeamRepository:      f.params.Repository.Team,
		PlayerRepository:    f.params.Repository.Player,
		TeamCacheRepository: f.params.Repository.TeamCache,
		Logger:              f.params.Logger,
	})
}

func (f *serviceFactory) CreatePlayerService() adapters.PlayerService {
	return NewPlayerService(PlayerServiceParams{
		PlayerRepository:    f.params.Repository.Player,
		TeamRepository:      f.params.Repository.Team,
		TeamCacheRepository: f.params.Repository.TeamCache,
		Logger:              f.params.Logger,
	})
}

func (f *serviceFactory) CreateTransferService() adapters.TransferService {
	return NewTransferService(TransferServiceParams{
		TransferRepository:  f.params.Repository.Transfer,
		PlayerRepository:    f.params.Repository.Player,
		TeamRepository:      f.params.Repository.Team,
		TeamCacheRepository: f.params.Repository.TeamCache,
		Logger:              f.params.Logger,
	})
}
