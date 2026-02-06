package models

// Запросы
type CreateUserRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Phone string `json:"phone,omitempty"`
}

type CreateMembershipRequest struct {
	UserID    int     `json:"user_id" binding:"required"`
	Type      string  `json:"type" binding:"required"`
	Price     float64 `json:"price" binding:"required"`
	StartDate string  `json:"start_date" binding:"required"`
	EndDate   string  `json:"end_date" binding:"required"`
}

type WorkoutRequest struct {
	UserID          int    `json:"user_id" binding:"required"`
	Type            string `json:"type" binding:"required"`
	DurationMinutes int    `json:"duration_minutes" binding:"required"`
	CaloriesBurned  int    `json:"calories_burned,omitempty"`
	Notes           string `json:"notes,omitempty"`
}
