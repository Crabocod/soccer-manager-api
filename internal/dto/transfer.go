package dto

import "soccer_manager_service/internal/entity"

type ListPlayerRequest struct {
	AskingPrice int64 `json:"asking_price" binding:"required,min=1"`
}

type TransferListItemResponse struct {
	Transfer entity.Transfer `json:"transfer"`
	Player   entity.Player   `json:"player"`
	Team     entity.Team     `json:"seller_team"`
}

type TransfersResponse struct {
	Transfers []TransferListItemResponse `json:"transfers"`
}
