package model

import "time"

type Workout struct {
	ID              int       `db:"id"`
	UserID          int       `db:"user_id"`
	Type            string    `db:"type"`
	DurationMinutes int       `db:"duration_minutes"`
	CaloriesBurned  int       `db:"calories_burned"`
	Notes           string    `db:"notes"`
	WorkoutDate     time.Time `db:"workout_date"`
	CreatedAt       time.Time `db:"created_at"`
}
