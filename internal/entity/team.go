package entity

import (
	"time"

	"github.com/google/uuid"
)

type Team struct {
	ID         uuid.UUID `db:"id" json:"id" goqu:"omitempty"`
	UserID     uuid.UUID `db:"user_id" json:"user_id" goqu:"omitempty"`
	Name       string    `db:"name" json:"name" goqu:"omitempty"`
	Country    string    `db:"country" json:"country" goqu:"omitempty"`
	Budget     int64     `db:"budget" json:"budget" goqu:"omitempty"`
	TotalValue int64     `db:"total_value" json:"total_value" goqu:"omitempty"`
	CreatedAt  time.Time `db:"created_at" json:"created_at" goqu:"omitempty"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at" goqu:"omitempty"`
}
