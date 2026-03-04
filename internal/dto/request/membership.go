package request

type CreateMembershipRequest struct {
	UserID    int     `json:"user_id" binding:"required"`
	Type      string  `json:"type" binding:"required"`
	Price     float64 `json:"price" binding:"required"`
	StartDate string  `json:"start_date" binding:"required"`
	EndDate   string  `json:"end_date" binding:"required"`
}
