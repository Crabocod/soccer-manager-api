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

type Transfer struct {
	logger  *zap.Logger
	builder *goqu.SelectDataset
	db      *pgxpool.Pool
}

type TransferParams struct {
	Postgres *pgxpool.Pool
	Logger   *zap.Logger
}

func NewTransferRepository(params TransferParams) *Transfer {
	return &Transfer{
		builder: goqu.Dialect(postgresdb).From(transfersTable),
		logger:  params.Logger.With(zap.String("layer", "TransferRepository")),
		db:      params.Postgres,
	}
}

func (r *Transfer) Create(ctx context.Context, playerID, sellerID uuid.UUID, askingPrice int64) (*entity.Transfer, error) {
	query := r.builder.
		Insert().
		Rows(goqu.Record{
			"player_id":    playerID,
			"seller_id":    sellerID,
			"asking_price": askingPrice,
			"status":       entity.TransferStatusActive,
		}).
		Returning(goqu.Star())

	sql, args, err := query.ToSQL()
	if err != nil {
		return nil, apperr.SQLError("Create", err)
	}

	var transfer entity.Transfer

	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&transfer.ID,
		&transfer.PlayerID,
		&transfer.SellerID,
		&transfer.BuyerID,
		&transfer.AskingPrice,
		&transfer.Status,
		&transfer.CreatedAt,
		&transfer.CompletedAt,
	)
	if err != nil {
		return nil, apperr.SQLQueryError("Create", err)
	}

	return &transfer, nil
}

func (r *Transfer) GetByID(ctx context.Context, id uuid.UUID) (*entity.Transfer, error) {
	query := r.builder.
		Select(goqu.Star()).
		Where(goqu.C("id").Eq(id))

	sql, args, err := query.ToSQL()
	if err != nil {
		return nil, apperr.SQLError("GetByID", err)
	}

	var transfer entity.Transfer

	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&transfer.ID,
		&transfer.PlayerID,
		&transfer.SellerID,
		&transfer.BuyerID,
		&transfer.AskingPrice,
		&transfer.Status,
		&transfer.CreatedAt,
		&transfer.CompletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperr.ErrTransferNotFound
		}

		return nil, apperr.SQLQueryError("GetByID", err)
	}

	return &transfer, nil
}

func (r *Transfer) GetActiveTransfers(ctx context.Context) ([]entity.Transfer, error) {
	query := r.builder.
		Select(goqu.Star()).
		Where(goqu.C("status").Eq(entity.TransferStatusActive)).
		Order(goqu.C("created_at").Desc())

	sql, args, err := query.ToSQL()
	if err != nil {
		return nil, apperr.SQLError("GetActiveTransfers", err)
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, apperr.SQLQueryError("GetActiveTransfers", err)
	}
	defer rows.Close()

	var transfers []entity.Transfer

	for rows.Next() {
		var transfer entity.Transfer

		err := rows.Scan(
			&transfer.ID,
			&transfer.PlayerID,
			&transfer.SellerID,
			&transfer.BuyerID,
			&transfer.AskingPrice,
			&transfer.Status,
			&transfer.CreatedAt,
			&transfer.CompletedAt,
		)
		if err != nil {
			return nil, apperr.SQLQueryError("GetActiveTransfers", err)
		}

		transfers = append(transfers, transfer)
	}

	return transfers, nil
}

func (r *Transfer) Complete(ctx context.Context, id, buyerID uuid.UUID) error {
	now := time.Now()

	query := r.builder.
		Update().
		Set(goqu.Record{
			"buyer_id":     buyerID,
			"status":       entity.TransferStatusCompleted,
			"completed_at": now,
		}).
		Where(goqu.C("id").Eq(id))

	sql, args, err := query.ToSQL()
	if err != nil {
		return apperr.SQLError("Complete", err)
	}

	result, err := r.db.Exec(ctx, sql, args...)
	if err != nil {
		return apperr.SQLExecError("Complete", err)
	}

	if result.RowsAffected() == 0 {
		return apperr.ErrTransferNotFound
	}

	return nil
}

func (r *Transfer) Cancel(ctx context.Context, id uuid.UUID) error {
	query := r.builder.
		Update().
		Set(goqu.Record{
			"status": entity.TransferStatusCancelled,
		}).
		Where(goqu.C("id").Eq(id))

	sql, args, err := query.ToSQL()
	if err != nil {
		return apperr.SQLError("Cancel", err)
	}

	result, err := r.db.Exec(ctx, sql, args...)
	if err != nil {
		return apperr.SQLExecError("Cancel", err)
	}

	if result.RowsAffected() == 0 {
		return apperr.ErrTransferNotFound
	}

	return nil
}

func (r *Transfer) GetByPlayerID(ctx context.Context, playerID uuid.UUID) (*entity.Transfer, error) {
	query := r.builder.
		Select(goqu.Star()).
		Where(
			goqu.C("player_id").Eq(playerID),
			goqu.C("status").Eq(entity.TransferStatusActive),
		)

	sql, args, err := query.ToSQL()
	if err != nil {
		return nil, apperr.SQLError("GetByPlayerID", err)
	}

	var transfer entity.Transfer

	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&transfer.ID,
		&transfer.PlayerID,
		&transfer.SellerID,
		&transfer.BuyerID,
		&transfer.AskingPrice,
		&transfer.Status,
		&transfer.CreatedAt,
		&transfer.CompletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperr.ErrTransferNotFound
		}

		return nil, apperr.SQLQueryError("GetByPlayerID", err)
	}

	return &transfer, nil
}
