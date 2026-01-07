package entity

import (
	"time"

	"github.com/google/uuid"
)

type TransferStatus string

const (
	TransferStatusActive    TransferStatus = "active"
	TransferStatusCompleted TransferStatus = "completed"
	TransferStatusCancelled TransferStatus = "cancelled"
)

type Transfer struct {
	ID          uuid.UUID      `db:"id" json:"id" goqu:"omitempty"`
	PlayerID    uuid.UUID      `db:"player_id" json:"player_id" goqu:"omitempty"`
	SellerID    uuid.UUID      `db:"seller_id" json:"seller_id" goqu:"omitempty"`
	BuyerID     *uuid.UUID     `db:"buyer_id" json:"buyer_id,omitempty" goqu:"omitempty"`
	AskingPrice int64          `db:"asking_price" json:"asking_price" goqu:"omitempty"`
	Status      TransferStatus `db:"status" json:"status" goqu:"omitempty"`
	CreatedAt   time.Time      `db:"created_at" json:"created_at" goqu:"omitempty"`
	CompletedAt *time.Time     `db:"completed_at" json:"completed_at,omitempty" goqu:"omitempty"`
}
