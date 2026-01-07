package usecase

import (
	"context"
	"fmt"
	"math/rand"
	"soccer_manager_service/internal/config"
	"soccer_manager_service/internal/dto"
	"soccer_manager_service/internal/entity"
	"soccer_manager_service/internal/ports"
	apperr "soccer_manager_service/pkg/errors"
	"soccer_manager_service/pkg/jwt"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var (
	firstNames = []string{"Oliver", "Jack", "Harry", "George", "Noah", "Charlie", "Leo", "Oscar", "Jacob", "Liam"}
	lastNames  = []string{"Smith", "Johnson", "Williams", "Brown", "Jones", "Garcia", "Miller", "Davis", "Martinez", "Hernandez"}
	countries  = []string{"England", "Spain", "Germany", "France", "Italy", "Brazil", "Argentina", "Portugal", "Netherlands", "Belgium"}
)

type AuthService struct {
	userRepository         ports.UserRepository
	teamRepository         ports.TeamRepository
	playerRepository       ports.PlayerRepository
	loginAttemptRepository ports.LoginAttemptRepository
	jwtManager             *jwt.Manager
	logger                 *zap.Logger
	config                 *config.Config
}

type AuthServiceParams struct {
	UserRepository         ports.UserRepository
	TeamRepository         ports.TeamRepository
	PlayerRepository       ports.PlayerRepository
	LoginAttemptRepository ports.LoginAttemptRepository
	JWTManager             *jwt.Manager
	Logger                 *zap.Logger
	Config                 *config.Config
}

func NewAuthService(params AuthServiceParams) *AuthService {
	return &AuthService{
		userRepository:         params.UserRepository,
		teamRepository:         params.TeamRepository,
		playerRepository:       params.PlayerRepository,
		loginAttemptRepository: params.LoginAttemptRepository,
		jwtManager:             params.JWTManager,
		logger:                 params.Logger.With(zap.String("service", "AuthService")),
		config:                 params.Config,
	}
}

func (s *AuthService) incrementLoginAttemptsOnError(ctx context.Context, email string) {
	if _, err := s.loginAttemptRepository.Increment(ctx, email); err != nil {
		s.logger.Error("failed to increment login attempts", zap.Error(err), zap.String("email", email))
	}
}

func (s *AuthService) Register(ctx context.Context, req *dto.RegisterRequest) (accessToken, refreshToken string, err error) {
	s.logger.Info("registering new user", zap.String("email", req.Email))

	existing, err := s.userRepository.GetByEmail(ctx, req.Email)
	if err == nil && existing != nil {
		return "", "", apperr.ErrUserAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("failed to hash password", zap.Error(err))

		return "", "", fmt.Errorf("hash password: %w", err)
	}

	user, err := s.userRepository.Create(ctx, req.Email, string(hashedPassword))
	if err != nil {
		s.logger.Error("failed to create user", zap.Error(err))

		return "", "", err
	}

	team, err := s.teamRepository.Create(ctx, user.ID, req.TeamName, req.Country, 5000000)
	if err != nil {
		s.logger.Error("failed to create team", zap.Error(err))

		return "", "", err
	}

	if err := s.createInitialPlayers(ctx, team.ID); err != nil {
		s.logger.Error("failed to create initial players", zap.Error(err))

		return "", "", err
	}

	totalValue := int64(20 * 1000000)

	if err := s.teamRepository.UpdateTotalValue(ctx, team.ID, totalValue); err != nil {
		s.logger.Error("failed to update team total value", zap.Error(err))

		return "", "", err
	}

	accessToken, err = s.jwtManager.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		s.logger.Error("failed to generate access token", zap.Error(err))

		return "", "", fmt.Errorf("generate access token: %w", err)
	}

	refreshToken, err = s.jwtManager.GenerateRefreshToken(user.ID, user.Email)
	if err != nil {
		s.logger.Error("failed to generate refresh token", zap.Error(err))

		return "", "", fmt.Errorf("generate refresh token: %w", err)
	}

	s.logger.Info("user registered successfully", zap.String("user_id", user.ID.String()))

	return accessToken, refreshToken, nil
}

func (s *AuthService) Login(ctx context.Context, req *dto.LoginRequest) (accessToken, refreshToken string, err error) {
	s.logger.Info("user login attempt", zap.String("email", req.Email))

	attempts, err := s.loginAttemptRepository.Get(ctx, req.Email)
	if err != nil {
		s.logger.Error("failed to get login attempts", zap.Error(err))

		return "", "", fmt.Errorf("get login attempts: %w", err)
	}

	if attempts >= s.config.Login.MaxLoginAttempts {
		s.logger.Warn("too many login attempts", zap.String("email", req.Email), zap.Int("attempts", attempts))

		return "", "", apperr.ErrTooManyAttempts
	}

	user, err := s.userRepository.GetByEmail(ctx, req.Email)
	if err != nil {
		s.logger.Warn("user not found", zap.String("email", req.Email))
		s.incrementLoginAttemptsOnError(ctx, req.Email)

		return "", "", apperr.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		s.logger.Warn("invalid password", zap.String("email", req.Email))
		s.incrementLoginAttemptsOnError(ctx, req.Email)

		return "", "", apperr.ErrInvalidCredentials
	}

	if err := s.loginAttemptRepository.Reset(ctx, req.Email); err != nil {
		s.logger.Error("failed to reset login attempts", zap.Error(err))
	}

	accessToken, err = s.jwtManager.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		s.logger.Error("failed to generate access token", zap.Error(err))

		return "", "", fmt.Errorf("generate access token: %w", err)
	}

	refreshToken, err = s.jwtManager.GenerateRefreshToken(user.ID, user.Email)
	if err != nil {
		s.logger.Error("failed to generate refresh token", zap.Error(err))

		return "", "", fmt.Errorf("generate refresh token: %w", err)
	}

	s.logger.Info("user logged in successfully", zap.String("user_id", user.ID.String()))

	return accessToken, refreshToken, nil
}

func (s *AuthService) createInitialPlayers(ctx context.Context, teamID uuid.UUID) error {
	positions := []struct {
		position entity.PlayerPosition
		count    int
	}{
		{entity.PositionGoalkeeper, 3},
		{entity.PositionDefender, 6},
		{entity.PositionMidfielder, 6},
		{entity.PositionAttacker, 5},
	}

	for _, p := range positions {
		for i := 0; i < p.count; i++ {
			firstName := firstNames[rand.Intn(len(firstNames))]
			lastName := lastNames[rand.Intn(len(lastNames))]
			country := countries[rand.Intn(len(countries))]
			age := 18 + rand.Intn(23)

			_, err := s.playerRepository.Create(ctx, teamID, firstName, lastName, country, age, p.position, 1000000)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
