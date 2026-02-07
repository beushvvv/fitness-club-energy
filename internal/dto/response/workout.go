package response

import "time"

type WorkoutResponse struct {
	ID              int       `json:"id"`
	UserID          int       `json:"user_id"`
	Type            string    `json:"type"`
	DurationMinutes int       `json:"duration_minutes"`
	CaloriesBurned  int       `json:"calories_burned,omitempty"`
	Notes           string    `json:"notes,omitempty"`
	WorkoutDate     time.Time `json:"workout_date"`
	CreatedAt       time.Time `json:"created_at"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Message string `json:"message"`
	ID      int    `json:"id,omitempty"`
}
