package repository

import (
	"fitness-club-energy/internal/model"

	"github.com/jmoiron/sqlx"
)

type WorkoutRepository struct {
	db *sqlx.DB
}

func NewWorkoutRepository(db *sqlx.DB) *WorkoutRepository {
	return &WorkoutRepository{db: db}
}

func (r *WorkoutRepository) FindAll() ([]model.Workout, error) {
	var workouts []model.Workout
	query := `SELECT id, user_id, type, duration_minutes, calories_burned, notes, workout_date, created_at FROM workouts ORDER BY workout_date DESC`
	err := r.db.Select(&workouts, query)
	return workouts, err
}

func (r *WorkoutRepository) FindByID(id int) (*model.Workout, error) {
	var workout model.Workout
	query := `SELECT id, user_id, type, duration_minutes, calories_burned, notes, workout_date, created_at FROM workouts WHERE id = $1`
	err := r.db.Get(&workout, query, id)
	if err != nil {
		return nil, err
	}
	return &workout, nil
}

func (r *WorkoutRepository) Create(workout *model.Workout) error {
	query := `INSERT INTO workouts (user_id, type, duration_minutes, calories_burned, notes, workout_date, created_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, NOW()) RETURNING id`
	err := r.db.QueryRow(query, workout.UserID, workout.Type, workout.DurationMinutes, workout.CaloriesBurned, workout.Notes, workout.WorkoutDate).Scan(&workout.ID)
	return err
}

func (r *WorkoutRepository) Update(workout *model.Workout) error {
	query := `UPDATE workouts SET user_id=$1, type=$2, duration_minutes=$3, calories_burned=$4, notes=$5, workout_date=$6 WHERE id=$7`
	_, err := r.db.Exec(query, workout.UserID, workout.Type, workout.DurationMinutes, workout.CaloriesBurned, workout.Notes, workout.WorkoutDate, workout.ID)
	return err
}

func (r *WorkoutRepository) Delete(id int) error {
	query := `DELETE FROM workouts WHERE id=$1`
	_, err := r.db.Exec(query, id)
	return err
}
