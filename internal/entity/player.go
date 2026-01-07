package entity

import (
	"time"

	"github.com/google/uuid"
)

type PlayerPosition string

const (
	PositionGoalkeeper PlayerPosition = "goalkeeper"
	PositionDefender   PlayerPosition = "defender"
	PositionMidfielder PlayerPosition = "midfielder"
	PositionAttacker   PlayerPosition = "attacker"
)

type Player struct {
	ID          uuid.UUID      `db:"id" json:"id" goqu:"omitempty"`
	TeamID      uuid.UUID      `db:"team_id" json:"team_id" goqu:"omitempty"`
	FirstName   string         `db:"first_name" json:"first_name" goqu:"omitempty"`
	LastName    string         `db:"last_name" json:"last_name" goqu:"omitempty"`
	Country     string         `db:"country" json:"country" goqu:"omitempty"`
	Age         int            `db:"age" json:"age" goqu:"omitempty"`
	Position    PlayerPosition `db:"position" json:"position" goqu:"omitempty"`
	MarketValue int64          `db:"market_value" json:"market_value" goqu:"omitempty"`
	CreatedAt   time.Time      `db:"created_at" json:"created_at" goqu:"omitempty"`
	UpdatedAt   time.Time      `db:"updated_at" json:"updated_at" goqu:"omitempty"`
}
