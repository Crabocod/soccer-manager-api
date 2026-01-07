package dto

import "soccer_manager_service/internal/entity"

type UpdateTeamRequest struct {
	Name    string `json:"name" binding:"omitempty,min=3,max=50"`
	Country string `json:"country" binding:"omitempty,min=2,max=50"`
}

type TeamWithPlayersResponse struct {
	Team    entity.Team     `json:"team"`
	Players []entity.Player `json:"players"`
}
