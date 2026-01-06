package repository

import (
	"soccer_manager_service/internal/ports"
	"soccer_manager_service/internal/repository/postgresrepo"
	"soccer_manager_service/internal/repository/redisrepo"
)

type repositoryFactory struct {
	deps Params
}

func newRepositoryFactory(deps Params) *repositoryFactory {
	return &repositoryFactory{deps: deps}
}

func (f *repositoryFactory) CreateUserRepository() ports.UserRepository {
	return postgresrepo.NewUserRepository(postgresrepo.UserParams{
		Postgres: f.deps.Postgres,
		Logger:   f.deps.Logger,
	})
}

func (f *repositoryFactory) CreateTeamRepository() ports.TeamRepository {
	return postgresrepo.NewTeamRepository(postgresrepo.TeamParams{
		Postgres: f.deps.Postgres,
		Logger:   f.deps.Logger,
	})
}

func (f *repositoryFactory) CreatePlayerRepository() ports.PlayerRepository {
	return postgresrepo.NewPlayerRepository(postgresrepo.PlayerParams{
		Postgres: f.deps.Postgres,
		Logger:   f.deps.Logger,
	})
}

func (f *repositoryFactory) CreateTransferRepository() ports.TransferRepository {
	return postgresrepo.NewTransferRepository(postgresrepo.TransferParams{
		Postgres: f.deps.Postgres,
		Logger:   f.deps.Logger,
	})
}

func (f *repositoryFactory) CreateLoginAttemptRepository() ports.LoginAttemptRepository {
	return redisrepo.NewLoginAttempt(redisrepo.LoginAttemptParams{
		Redis:  f.deps.Redis,
		Logger: f.deps.Logger,
		Config: f.deps.Config,
	})
}

func (f *repositoryFactory) CreateTeamCacheRepository() ports.TeamCacheRepository {
	return redisrepo.NewTeamCache(redisrepo.TeamCacheParams{
		Redis:  f.deps.Redis,
		Logger: f.deps.Logger,
		Config: f.deps.Config,
	})
}
