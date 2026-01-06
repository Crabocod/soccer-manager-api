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

type Team struct {
	ID          uuid.UUID `db:"id" json:"id" goqu:"omitempty"`
	UserID      uuid.UUID `db:"user_id" json:"user_id" goqu:"omitempty"`
	Name        string    `db:"name" json:"name" goqu:"omitempty"`
	Country     string    `db:"country" json:"country" goqu:"omitempty"`
	Budget      int64     `db:"budget" json:"budget" goqu:"omitempty"`
	TotalValue  int64     `db:"total_value" json:"total_value" goqu:"omitempty"`
	CreatedAt   time.Time `db:"created_at" json:"created_at" goqu:"omitempty"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at" goqu:"omitempty"`
}

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

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	TeamName string `json:"team_name" binding:"required,min=3,max=50"`
	Country  string `json:"country" binding:"required,min=2,max=50"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UpdateTeamRequest struct {
	Name    string `json:"name" binding:"omitempty,min=3,max=50"`
	Country string `json:"country" binding:"omitempty,min=2,max=50"`
}

type UpdatePlayerRequest struct {
	FirstName string `json:"first_name" binding:"omitempty,min=2,max=50"`
	LastName  string `json:"last_name" binding:"omitempty,min=2,max=50"`
	Country   string `json:"country" binding:"omitempty,min=2,max=50"`
}

type ListPlayerRequest struct {
	AskingPrice int64 `json:"asking_price" binding:"required,min=1"`
}

type TeamWithPlayers struct {
	Team    Team     `json:"team"`
	Players []Player `json:"players"`
}

type TransferListItem struct {
	Transfer Transfer `json:"transfer"`
	Player   Player   `json:"player"`
	Team     Team     `json:"seller_team"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type TransfersResponse struct {
	Transfers []TransferListItem `json:"transfers"`
}
