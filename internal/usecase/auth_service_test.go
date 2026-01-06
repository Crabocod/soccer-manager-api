package usecase

import (
	"context"
	"errors"
	"testing"

	"soccer_manager_service/internal/config"
	"soccer_manager_service/internal/entity"
	apperr "soccer_manager_service/pkg/errors"
	"soccer_manager_service/pkg/jwt"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash)
}

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, email, passwordHash string) (*entity.User, error) {
	args := m.Called(ctx, email, passwordHash)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	args := m.Called(ctx, id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*entity.User), args.Error(1)
}

type MockLoginAttemptRepository struct {
	mock.Mock
}

func (m *MockLoginAttemptRepository) Increment(ctx context.Context, email string) (int, error) {
	args := m.Called(ctx, email)

	return args.Int(0), args.Error(1)
}

func (m *MockLoginAttemptRepository) Get(ctx context.Context, email string) (int, error) {
	args := m.Called(ctx, email)

	return args.Int(0), args.Error(1)
}

func (m *MockLoginAttemptRepository) Reset(ctx context.Context, email string) error {
	args := m.Called(ctx, email)

	return args.Error(0)
}

func TestAuthService_Register(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()
	jwtManager := jwt.NewManager("test-secret", 0, 0)

	cfg := &config.Config{
		Login: config.LoginConfig{
			MaxLoginAttempts: 5,
		},
	}

	t.Run("success", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockTeamRepo := new(MockTeamRepository)
		mockPlayerRepo := new(MockPlayerRepository)
		mockLoginAttemptRepo := new(MockLoginAttemptRepository)

		userID := uuid.New()
		teamID := uuid.New()

		user := &entity.User{
			ID:    userID,
			Email: "test@example.com",
		}

		team := &entity.Team{
			ID:     teamID,
			UserID: userID,
		}

		mockUserRepo.On("GetByEmail", ctx, "test@example.com").Return(nil, apperr.ErrUserNotFound)
		mockUserRepo.On("Create", ctx, "test@example.com", mock.AnythingOfType("string")).Return(user, nil)
		mockTeamRepo.On("Create", ctx, userID, "Test Team", "England", int64(5000000)).Return(team, nil)
		mockPlayerRepo.On("Create", ctx, teamID, mock.AnythingOfType("string"), mock.AnythingOfType("string"), 
			mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("entity.PlayerPosition"), 
			int64(1000000)).Return(&entity.Player{}, nil).Times(20)
		mockTeamRepo.On("UpdateTotalValue", ctx, teamID, int64(20000000)).Return(nil)

		service := NewAuthService(AuthServiceParams{
			UserRepository:         mockUserRepo,
			TeamRepository:         mockTeamRepo,
			PlayerRepository:       mockPlayerRepo,
			LoginAttemptRepository: mockLoginAttemptRepo,
			JWTManager:             jwtManager,
			Logger:                 logger,
			Config:                 cfg,
		})

		req := &entity.RegisterRequest{
			Email:    "test@example.com",
			Password: "password123",
			TeamName: "Test Team",
			Country:  "England",
		}

		accessToken, refreshToken, err := service.Register(ctx, req)

		assert.NoError(t, err)
		assert.NotEmpty(t, accessToken)
		assert.NotEmpty(t, refreshToken)
		mockUserRepo.AssertExpectations(t)
		mockTeamRepo.AssertExpectations(t)
		mockPlayerRepo.AssertExpectations(t)
	})

	t.Run("user already exists", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockTeamRepo := new(MockTeamRepository)
		mockPlayerRepo := new(MockPlayerRepository)
		mockLoginAttemptRepo := new(MockLoginAttemptRepository)

		existingUser := &entity.User{
			ID:    uuid.New(),
			Email: "test@example.com",
		}

		mockUserRepo.On("GetByEmail", ctx, "test@example.com").Return(existingUser, nil)

		service := NewAuthService(AuthServiceParams{
			UserRepository:         mockUserRepo,
			TeamRepository:         mockTeamRepo,
			PlayerRepository:       mockPlayerRepo,
			LoginAttemptRepository: mockLoginAttemptRepo,
			JWTManager:             jwtManager,
			Logger:                 logger,
			Config:                 cfg,
		})

		req := &entity.RegisterRequest{
			Email:    "test@example.com",
			Password: "password123",
			TeamName: "Test Team",
			Country:  "England",
		}

		accessToken, refreshToken, err := service.Register(ctx, req)

		assert.Error(t, err)
		assert.Empty(t, accessToken)
		assert.Empty(t, refreshToken)
		assert.Equal(t, apperr.ErrUserAlreadyExists, err)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("failed to create user", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockTeamRepo := new(MockTeamRepository)
		mockPlayerRepo := new(MockPlayerRepository)
		mockLoginAttemptRepo := new(MockLoginAttemptRepository)

		mockUserRepo.On("GetByEmail", ctx, "test@example.com").Return(nil, apperr.ErrUserNotFound)
		mockUserRepo.On("Create", ctx, "test@example.com", mock.AnythingOfType("string")).Return(nil, errors.New("database error"))

		service := NewAuthService(AuthServiceParams{
			UserRepository:         mockUserRepo,
			TeamRepository:         mockTeamRepo,
			PlayerRepository:       mockPlayerRepo,
			LoginAttemptRepository: mockLoginAttemptRepo,
			JWTManager:             jwtManager,
			Logger:                 logger,
			Config:                 cfg,
		})

		req := &entity.RegisterRequest{
			Email:    "test@example.com",
			Password: "password123",
			TeamName: "Test Team",
			Country:  "England",
		}

		accessToken, refreshToken, err := service.Register(ctx, req)

		assert.Error(t, err)
		assert.Empty(t, accessToken)
		assert.Empty(t, refreshToken)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("failed to create team", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockTeamRepo := new(MockTeamRepository)
		mockPlayerRepo := new(MockPlayerRepository)
		mockLoginAttemptRepo := new(MockLoginAttemptRepository)

		userID := uuid.New()
		user := &entity.User{
			ID:    userID,
			Email: "test@example.com",
		}

		mockUserRepo.On("GetByEmail", ctx, "test@example.com").Return(nil, apperr.ErrUserNotFound)
		mockUserRepo.On("Create", ctx, "test@example.com", mock.AnythingOfType("string")).Return(user, nil)
		mockTeamRepo.On("Create", ctx, userID, "Test Team", "England", int64(5000000)).Return(nil, errors.New("database error"))

		service := NewAuthService(AuthServiceParams{
			UserRepository:         mockUserRepo,
			TeamRepository:         mockTeamRepo,
			PlayerRepository:       mockPlayerRepo,
			LoginAttemptRepository: mockLoginAttemptRepo,
			JWTManager:             jwtManager,
			Logger:                 logger,
			Config:                 cfg,
		})

		req := &entity.RegisterRequest{
			Email:    "test@example.com",
			Password: "password123",
			TeamName: "Test Team",
			Country:  "England",
		}

		accessToken, refreshToken, err := service.Register(ctx, req)

		assert.Error(t, err)
		assert.Empty(t, accessToken)
		assert.Empty(t, refreshToken)
		mockUserRepo.AssertExpectations(t)
		mockTeamRepo.AssertExpectations(t)
	})
}

func TestAuthService_Login(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()
	jwtManager := jwt.NewManager("test-secret", 0, 0)

	cfg := &config.Config{
		Login: config.LoginConfig{
			MaxLoginAttempts: 5,
		},
	}

	t.Run("success", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockTeamRepo := new(MockTeamRepository)
		mockPlayerRepo := new(MockPlayerRepository)
		mockLoginAttemptRepo := new(MockLoginAttemptRepository)

		userID := uuid.New()
		user := &entity.User{
			ID:           userID,
			Email:        "test@example.com",
			PasswordHash: hashPassword("password"),
		}

		mockLoginAttemptRepo.On("Get", ctx, "test@example.com").Return(0, nil)
		mockUserRepo.On("GetByEmail", ctx, "test@example.com").Return(user, nil)
		mockLoginAttemptRepo.On("Reset", ctx, "test@example.com").Return(nil)

		service := NewAuthService(AuthServiceParams{
			UserRepository:         mockUserRepo,
			TeamRepository:         mockTeamRepo,
			PlayerRepository:       mockPlayerRepo,
			LoginAttemptRepository: mockLoginAttemptRepo,
			JWTManager:             jwtManager,
			Logger:                 logger,
			Config:                 cfg,
		})

		req := &entity.LoginRequest{
			Email:    "test@example.com",
			Password: "password",
		}

		accessToken, refreshToken, err := service.Login(ctx, req)

		assert.NoError(t, err)
		assert.NotEmpty(t, accessToken)
		assert.NotEmpty(t, refreshToken)
		mockUserRepo.AssertExpectations(t)
		mockLoginAttemptRepo.AssertExpectations(t)
	})

	t.Run("too many attempts", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockTeamRepo := new(MockTeamRepository)
		mockPlayerRepo := new(MockPlayerRepository)
		mockLoginAttemptRepo := new(MockLoginAttemptRepository)

		mockLoginAttemptRepo.On("Get", ctx, "test@example.com").Return(5, nil)

		service := NewAuthService(AuthServiceParams{
			UserRepository:         mockUserRepo,
			TeamRepository:         mockTeamRepo,
			PlayerRepository:       mockPlayerRepo,
			LoginAttemptRepository: mockLoginAttemptRepo,
			JWTManager:             jwtManager,
			Logger:                 logger,
			Config:                 cfg,
		})

		req := &entity.LoginRequest{
			Email:    "test@example.com",
			Password: "password",
		}

		accessToken, refreshToken, err := service.Login(ctx, req)

		assert.Error(t, err)
		assert.Empty(t, accessToken)
		assert.Empty(t, refreshToken)
		assert.Equal(t, apperr.ErrTooManyAttempts, err)
		mockLoginAttemptRepo.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockTeamRepo := new(MockTeamRepository)
		mockPlayerRepo := new(MockPlayerRepository)
		mockLoginAttemptRepo := new(MockLoginAttemptRepository)

		mockLoginAttemptRepo.On("Get", ctx, "test@example.com").Return(0, nil)
		mockUserRepo.On("GetByEmail", ctx, "test@example.com").Return(nil, apperr.ErrUserNotFound)
		mockLoginAttemptRepo.On("Increment", ctx, "test@example.com").Return(1, nil)

		service := NewAuthService(AuthServiceParams{
			UserRepository:         mockUserRepo,
			TeamRepository:         mockTeamRepo,
			PlayerRepository:       mockPlayerRepo,
			LoginAttemptRepository: mockLoginAttemptRepo,
			JWTManager:             jwtManager,
			Logger:                 logger,
			Config:                 cfg,
		})

		req := &entity.LoginRequest{
			Email:    "test@example.com",
			Password: "password",
		}

		accessToken, refreshToken, err := service.Login(ctx, req)

		assert.Error(t, err)
		assert.Empty(t, accessToken)
		assert.Empty(t, refreshToken)
		assert.Equal(t, apperr.ErrInvalidCredentials, err)
		mockUserRepo.AssertExpectations(t)
		mockLoginAttemptRepo.AssertExpectations(t)
	})

	t.Run("invalid password", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockTeamRepo := new(MockTeamRepository)
		mockPlayerRepo := new(MockPlayerRepository)
		mockLoginAttemptRepo := new(MockLoginAttemptRepository)

		userID := uuid.New()
		user := &entity.User{
			ID:           userID,
			Email:        "test@example.com",
			PasswordHash: hashPassword("password"),
		}

		mockLoginAttemptRepo.On("Get", ctx, "test@example.com").Return(0, nil)
		mockUserRepo.On("GetByEmail", ctx, "test@example.com").Return(user, nil)
		mockLoginAttemptRepo.On("Increment", ctx, "test@example.com").Return(1, nil)

		service := NewAuthService(AuthServiceParams{
			UserRepository:         mockUserRepo,
			TeamRepository:         mockTeamRepo,
			PlayerRepository:       mockPlayerRepo,
			LoginAttemptRepository: mockLoginAttemptRepo,
			JWTManager:             jwtManager,
			Logger:                 logger,
			Config:                 cfg,
		})

		req := &entity.LoginRequest{
			Email:    "test@example.com",
			Password: "wrongpassword",
		}

		accessToken, refreshToken, err := service.Login(ctx, req)

		assert.Error(t, err)
		assert.Empty(t, accessToken)
		assert.Empty(t, refreshToken)
		assert.Equal(t, apperr.ErrInvalidCredentials, err)
		mockUserRepo.AssertExpectations(t)
		mockLoginAttemptRepo.AssertExpectations(t)
	})
}
