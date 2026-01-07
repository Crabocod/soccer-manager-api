package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `db:"id" json:"id" goqu:"omitempty"`
	Email        string    `db:"email" json:"email" goqu:"omitempty"`
	PasswordHash string    `db:"password_hash" json:"-" goqu:"omitempty"`
	CreatedAt    time.Time `db:"created_at" json:"created_at" goqu:"omitempty"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at" goqu:"omitempty"`
}
