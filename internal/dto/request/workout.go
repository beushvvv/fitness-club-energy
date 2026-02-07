package request

type CreateWorkoutRequest struct {
	UserID          int    `json:"user_id" binding:"required"`
	Type            string `json:"type" binding:"required"`
	DurationMinutes int    `json:"duration_minutes" binding:"required"`
	CaloriesBurned  int    `json:"calories_burned,omitempty"`
	Notes           string `json:"notes,omitempty"`
}

type FilterWorkoutsRequest struct {
	Type   string `form:"type"`
	Date   string `form:"date"`
	UserID int    `form:"user_id"`
}
