package models

import "time"

type Workout struct {
	ID              int       `json:"id" db:"id"`
	UserID          int       `json:"user_id" db:"user_id"`
	Type            string    `json:"type" db:"type"`
	DurationMinutes int       `json:"duration_minutes" db:"duration_minutes"`
	CaloriesBurned  int       `json:"calories_burned,omitempty" db:"calories_burned"`
	Notes           string    `json:"notes,omitempty" db:"notes"`
	WorkoutDate     time.Time `json:"workout_date" db:"workout_date"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}
