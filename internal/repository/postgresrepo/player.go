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
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Player struct {
	logger  *zap.Logger
	builder *goqu.SelectDataset
	db      *pgxpool.Pool
}

type PlayerParams struct {
	Postgres *pgxpool.Pool
	Logger   *zap.Logger
}

func NewPlayerRepository(params PlayerParams) *Player {
	return &Player{
		builder: goqu.Dialect(postgresdb).From(playersTable),
		logger:  params.Logger.With(zap.String("layer", "PlayerRepository")),
		db:      params.Postgres,
	}
}

func (r *Player) Create(ctx context.Context, teamID uuid.UUID, firstName, lastName, country string, age int, position entity.PlayerPosition, marketValue int64) (*entity.Player, error) {
	query := r.builder.
		Insert().
		Rows(goqu.Record{
			"team_id":      teamID,
			"first_name":   firstName,
			"last_name":    lastName,
			"country":      country,
			"age":          age,
			"position":     position,
			"market_value": marketValue,
		}).
		Returning(goqu.Star())

	sql, args, err := query.ToSQL()
	if err != nil {
		return nil, apperr.SQLError("Create", err)
	}

	var player entity.Player

	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&player.ID,
		&player.TeamID,
		&player.FirstName,
		&player.LastName,
		&player.Country,
		&player.Age,
		&player.Position,
		&player.MarketValue,
		&player.CreatedAt,
		&player.UpdatedAt,
	)
	if err != nil {
		return nil, apperr.SQLQueryError("Create", err)
	}

	return &player, nil
}

func (r *Player) GetByID(ctx context.Context, id uuid.UUID) (*entity.Player, error) {
	query := r.builder.
		Select(goqu.Star()).
		Where(goqu.C("id").Eq(id))

	sql, args, err := query.ToSQL()
	if err != nil {
		return nil, apperr.SQLError("GetByID", err)
	}

	var player entity.Player

	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&player.ID,
		&player.TeamID,
		&player.FirstName,
		&player.LastName,
		&player.Country,
		&player.Age,
		&player.Position,
		&player.MarketValue,
		&player.CreatedAt,
		&player.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperr.ErrPlayerNotFound
		}

		return nil, apperr.SQLQueryError("GetByID", err)
	}

	return &player, nil
}

func (r *Player) GetByTeamID(ctx context.Context, teamID uuid.UUID) ([]entity.Player, error) {
	query := r.builder.
		Select(goqu.Star()).
		Where(goqu.C("team_id").Eq(teamID)).
		Order(goqu.C("position").Asc(), goqu.C("last_name").Asc())

	sql, args, err := query.ToSQL()
	if err != nil {
		return nil, apperr.SQLError("GetByTeamID", err)
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, apperr.SQLQueryError("GetByTeamID", err)
	}
	defer rows.Close()

	var players []entity.Player

	for rows.Next() {
		var player entity.Player

		err := rows.Scan(
			&player.ID,
			&player.TeamID,
			&player.FirstName,
			&player.LastName,
			&player.Country,
			&player.Age,
			&player.Position,
			&player.MarketValue,
			&player.CreatedAt,
			&player.UpdatedAt,
		)
		if err != nil {
			return nil, apperr.SQLQueryError("GetByTeamID", err)
		}

		players = append(players, player)
	}

	return players, nil
}

func (r *Player) Update(ctx context.Context, id uuid.UUID, firstName, lastName, country string) (*entity.Player, error) {
	record := goqu.Record{"updated_at": time.Now()}

	if firstName != "" {
		record["first_name"] = firstName
	}

	if lastName != "" {
		record["last_name"] = lastName
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

	var player entity.Player

	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&player.ID,
		&player.TeamID,
		&player.FirstName,
		&player.LastName,
		&player.Country,
		&player.Age,
		&player.Position,
		&player.MarketValue,
		&player.CreatedAt,
		&player.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperr.ErrPlayerNotFound
		}

		return nil, apperr.SQLQueryError("Update", err)
	}

	return &player, nil
}

func (r *Player) UpdateMarketValue(ctx context.Context, id uuid.UUID, marketValue int64) error {
	query := r.builder.
		Update().
		Set(goqu.Record{
			"market_value": marketValue,
			"updated_at":   time.Now(),
		}).
		Where(goqu.C("id").Eq(id))

	sql, args, err := query.ToSQL()
	if err != nil {
		return apperr.SQLError("UpdateMarketValue", err)
	}

	result, err := r.db.Exec(ctx, sql, args...)
	if err != nil {
		return apperr.SQLExecError("UpdateMarketValue", err)
	}

	if result.RowsAffected() == 0 {
		return apperr.ErrPlayerNotFound
	}

	return nil
}

func (r *Player) TransferPlayer(ctx context.Context, playerID, newTeamID uuid.UUID) error {
	query := r.builder.
		Update().
		Set(goqu.Record{
			"team_id":    newTeamID,
			"updated_at": time.Now(),
		}).
		Where(goqu.C("id").Eq(playerID))

	sql, args, err := query.ToSQL()
	if err != nil {
		return apperr.SQLError("TransferPlayer", err)
	}

	result, err := r.db.Exec(ctx, sql, args...)
	if err != nil {
		return apperr.SQLExecError("TransferPlayer", err)
	}

	if result.RowsAffected() == 0 {
		return apperr.ErrPlayerNotFound
	}

	return nil
}
