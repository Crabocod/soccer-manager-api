package postgresrepo

import (
	"context"
	"errors"
	"soccer_manager_service/internal/entity"
	"soccer_manager_service/pkg/errors"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type User struct {
	logger  *zap.Logger
	builder *goqu.SelectDataset
	db      *pgxpool.Pool
}

type UserParams struct {
	Postgres *pgxpool.Pool
	Logger   *zap.Logger
}

func NewUserRepository(params UserParams) *User {
	return &User{
		builder: goqu.Dialect(postgresdb).From(usersTable),
		logger:  params.Logger.With(zap.String("layer", "UserRepository")),
		db:      params.Postgres,
	}
}

func (r *User) Create(ctx context.Context, email, passwordHash string) (*entity.User, error) {
	query := r.builder.
		Insert().
		Rows(goqu.Record{
			"email":         email,
			"password_hash": passwordHash,
		}).
		Returning(goqu.Star())

	sql, args, err := query.ToSQL()
	if err != nil {
		return nil, apperr.SQLError("Create", err)
	}

	var user entity.User

	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, apperr.ErrUserAlreadyExists
		}

		return nil, apperr.SQLQueryError("Create", err)
	}

	return &user, nil
}

func (r *User) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	query := r.builder.
		Select(goqu.Star()).
		Where(goqu.C("id").Eq(id))

	sql, args, err := query.ToSQL()
	if err != nil {
		return nil, apperr.SQLError("GetByID", err)
	}

	var user entity.User

	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperr.ErrUserNotFound
		}

		return nil, apperr.SQLQueryError("GetByID", err)
	}

	return &user, nil
}

func (r *User) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := r.builder.
		Select(goqu.Star()).
		Where(goqu.C("email").Eq(email))

	sql, args, err := query.ToSQL()
	if err != nil {
		return nil, apperr.SQLError("GetByEmail", err)
	}

	var user entity.User

	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperr.ErrUserNotFound
		}

		return nil, apperr.SQLQueryError("GetByEmail", err)
	}

	return &user, nil
}
