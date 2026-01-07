package dto

type UpdatePlayerRequest struct {
	FirstName string `json:"first_name" binding:"omitempty,min=2,max=50"`
	LastName  string `json:"last_name" binding:"omitempty,min=2,max=50"`
	Country   string `json:"country" binding:"omitempty,min=2,max=50"`
}
