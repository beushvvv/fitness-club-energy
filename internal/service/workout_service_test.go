package service

import (
	"testing"
	"time"

	"fitness-club-energy/internal/cache"
	"fitness-club-energy/internal/logger"
	"fitness-club-energy/internal/model"
	"fitness-club-energy/internal/repository"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	logger.Init("debug")
}

func setupTestDBWorkout(t *testing.T) (*sqlx.DB, func()) {
	db, err := sqlx.Connect("postgres", "host=localhost port=5433 user=postgres password=postgres123 dbname=fitness_club sslmode=disable")
	require.NoError(t, err)

	_, err = db.Exec("DELETE FROM workouts")
	require.NoError(t, err)

	cleanup := func() {
		db.Exec("DELETE FROM workouts")
		db.Close()
	}
	return db, cleanup
}

func setupTestCacheWorkout(t *testing.T) (*cache.RedisClient, func()) {
	client := cache.NewRedisClient("localhost:6379", "", 0)
	err := client.Ping()
	if err != nil {
		t.Skip("Redis not available, skipping test")
	}
	client.Delete("workouts:all")
	client.Delete("workout:1")

	cleanup := func() {
		client.Delete("workouts:all")
		client.Delete("workout:1")
		client.Close()
	}
	return client, cleanup
}

func TestWorkoutService_GetAllWorkouts_Integration(t *testing.T) {
	db, dbCleanup := setupTestDBWorkout(t)
	defer dbCleanup()

	redisClient, redisCleanup := setupTestCacheWorkout(t)
	defer redisCleanup()

	cacheWrapper := cache.NewCacheWrapper(redisClient)
	workoutRepo := repository.NewWorkoutRepository(db)
	service := NewWorkoutService(workoutRepo, cacheWrapper)

	// Создаём пользователя
	var userID int
	err := db.QueryRow("INSERT INTO users (name, email, phone, created_at) VALUES ($1, $2, $3, NOW()) RETURNING id", "Test", "test@test.com", "123").Scan(&userID)
	require.NoError(t, err)

	// Создаём тренировку
	testWorkout := &model.Workout{
		UserID:          userID,
		Type:            "cardio",
		DurationMinutes: 30,
		CaloriesBurned:  300,
		Notes:           "test",
		WorkoutDate:     time.Now(),
	}
	err = workoutRepo.Create(testWorkout)
	require.NoError(t, err)

	// Первый запрос — из БД
	workouts1, err := service.GetAllWorkouts()
	assert.NoError(t, err)
	assert.Len(t, workouts1, 1)

	// Второй запрос — из кэша
	workouts2, err := service.GetAllWorkouts()
	assert.NoError(t, err)
	assert.Len(t, workouts2, 1)

	assert.Equal(t, workouts1[0].Type, workouts2[0].Type)
}

func TestWorkoutService_CreateWorkout_Integration(t *testing.T) {
	db, dbCleanup := setupTestDBWorkout(t)
	defer dbCleanup()

	redisClient, redisCleanup := setupTestCacheWorkout(t)
	defer redisCleanup()

	cacheWrapper := cache.NewCacheWrapper(redisClient)
	workoutRepo := repository.NewWorkoutRepository(db)
	service := NewWorkoutService(workoutRepo, cacheWrapper)

	// Создаём пользователя
	var userID int
	err := db.QueryRow("INSERT INTO users (name, email, phone, created_at) VALUES ($1, $2, $3, NOW()) RETURNING id", "Test2", "test2@test.com", "456").Scan(&userID)
	require.NoError(t, err)

	workout := &model.Workout{
		UserID:          userID,
		Type:            "strength",
		DurationMinutes: 45,
		CaloriesBurned:  400,
		Notes:           "test workout",
		WorkoutDate:     time.Now(),
	}

	err = service.CreateWorkout(workout)
	assert.NoError(t, err)
	assert.NotZero(t, workout.ID)

	// Проверяем что создалось
	created, err := workoutRepo.FindByID(workout.ID)
	assert.NoError(t, err)
	assert.Equal(t, workout.Type, created.Type)
}
