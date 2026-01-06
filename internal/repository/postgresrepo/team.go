package postgresrepo

import (
	"context"
	"errors"
	"soccer_manager_service/internal/entity"
	"soccer_manager_service/pkg/errors"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Team struct {
	logger  *zap.Logger
	builder *goqu.SelectDataset
	db      *pgxpool.Pool
}

type TeamParams struct {
	Postgres *pgxpool.Pool
	Logger   *zap.Logger
}

func NewTeamRepository(params TeamParams) *Team {
	return &Team{
		builder: goqu.Dialect(postgresdb).From(teamsTable),
		logger:  params.Logger.With(zap.String("layer", "TeamRepository")),
		db:      params.Postgres,
	}
}

func (r *Team) Create(ctx context.Context, userID uuid.UUID, name, country string, budget int64) (*entity.Team, error) {
	query := r.builder.
		Insert().
		Rows(goqu.Record{
			"user_id": userID,
			"name":    name,
			"country": country,
			"budget":  budget,
		}).
		Returning(goqu.Star())

	sql, args, err := query.ToSQL()
	if err != nil {
		return nil, apperr.SQLError("Create", err)
	}

	var team entity.Team

	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&team.ID,
		&team.UserID,
		&team.Name,
		&team.Country,
		&team.Budget,
		&team.TotalValue,
		&team.CreatedAt,
		&team.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, apperr.ErrTeamAlreadyExists
		}

		return nil, apperr.SQLQueryError("Create", err)
	}

	return &team, nil
}

func (r *Team) GetByID(ctx context.Context, id uuid.UUID) (*entity.Team, error) {
	query := r.builder.
		Select(goqu.Star()).
		Where(goqu.C("id").Eq(id))

	sql, args, err := query.ToSQL()
	if err != nil {
		return nil, apperr.SQLError("GetByID", err)
	}

	var team entity.Team

	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&team.ID,
		&team.UserID,
		&team.Name,
		&team.Country,
		&team.Budget,
		&team.TotalValue,
		&team.CreatedAt,
		&team.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperr.ErrTeamNotFound
		}

		return nil, apperr.SQLQueryError("GetByID", err)
	}

	return &team, nil
}

func (r *Team) GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.Team, error) {
	query := r.builder.
		Select(goqu.Star()).
		Where(goqu.C("user_id").Eq(userID))

	sql, args, err := query.ToSQL()
	if err != nil {
		return nil, apperr.SQLError("GetByUserID", err)
	}

	var team entity.Team

	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&team.ID,
		&team.UserID,
		&team.Name,
		&team.Country,
		&team.Budget,
		&team.TotalValue,
		&team.CreatedAt,
		&team.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperr.ErrTeamNotFound
		}

		return nil, apperr.SQLQueryError("GetByUserID", err)
	}

	return &team, nil
}

func (r *Team) Update(ctx context.Context, id uuid.UUID, name, country string) (*entity.Team, error) {
	record := goqu.Record{"updated_at": time.Now()}

	if name != "" {
		record["name"] = name
	}

	if country != "" {
		record["country"] = country
	}

	query := r.builder.
		Update().
		Set(record).
		Where(goqu.C("id").Eq(id)).
		Returning(goqu.Star())

	sql, args, err := query.ToSQL()
	if err != nil {
		return nil, apperr.SQLError("Update", err)
	}

	var team entity.Team

	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&team.ID,
		&team.UserID,
		&team.Name,
		&team.Country,
		&team.Budget,
		&team.TotalValue,
		&team.CreatedAt,
		&team.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperr.ErrTeamNotFound
		}

		return nil, apperr.SQLQueryError("Update", err)
	}

	return &team, nil
}

func (r *Team) UpdateBudget(ctx context.Context, id uuid.UUID, budget int64) error {
	query := r.builder.
		Update().
		Set(goqu.Record{
			"budget":     budget,
			"updated_at": time.Now(),
		}).
		Where(goqu.C("id").Eq(id))

	sql, args, err := query.ToSQL()
	if err != nil {
		return apperr.SQLError("UpdateBudget", err)
	}

	result, err := r.db.Exec(ctx, sql, args...)
	if err != nil {
		return apperr.SQLExecError("UpdateBudget", err)
	}

	if result.RowsAffected() == 0 {
		return apperr.ErrTeamNotFound
	}

	return nil
}

func (r *Team) UpdateTotalValue(ctx context.Context, id uuid.UUID, totalValue int64) error {
	query := r.builder.
		Update().
		Set(goqu.Record{
			"total_value": totalValue,
			"updated_at":  time.Now(),
		}).
		Where(goqu.C("id").Eq(id))

	sql, args, err := query.ToSQL()
	if err != nil {
		return apperr.SQLError("UpdateTotalValue", err)
	}

	result, err := r.db.Exec(ctx, sql, args...)
	if err != nil {
		return apperr.SQLExecError("UpdateTotalValue", err)
	}

	if result.RowsAffected() == 0 {
		return apperr.ErrTeamNotFound
	}

	return nil
}
