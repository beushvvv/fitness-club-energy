package response

import (
	"fitness-club-energy/internal/model"
	"time"
)

// WorkoutResponse представляет тренировку для ответа API
type WorkoutResponse struct {
	ID              int       `json:"id"`
	UserID          int       `json:"user_id"`
	Type            string    `json:"type"`
	DurationMinutes int       `json:"duration_minutes"`
	CaloriesBurned  int       `json:"calories_burned"`
	Notes           string    `json:"notes,omitempty"`
	WorkoutDate     time.Time `json:"workout_date"`
	CreatedAt       time.Time `json:"created_at"`
}

// FromWorkout конвертирует модель Workout в WorkoutResponse
func FromWorkout(workout *model.Workout) WorkoutResponse {
	return WorkoutResponse{
		ID:              workout.ID,
		UserID:          workout.UserID,
		Type:            workout.Type,
		DurationMinutes: workout.DurationMinutes,
		CaloriesBurned:  workout.CaloriesBurned,
		Notes:           workout.Notes,
		WorkoutDate:     workout.WorkoutDate,
		CreatedAt:       workout.CreatedAt,
	}
}

// FromWorkouts конвертирует срез моделей Workout в срез WorkoutResponse
func FromWorkouts(workouts []model.Workout) []WorkoutResponse {
	result := make([]WorkoutResponse, len(workouts))
	for i, workout := range workouts {
		result[i] = FromWorkout(&workout)
	}
	return result
}
