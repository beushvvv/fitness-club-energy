package service

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fitness-club-energy/internal/cache"
	"fitness-club-energy/internal/model"
	"time"
)

type WorkoutService struct {
	cache     *cache.RedisClient
	cacheWrap *cache.CacheWrapper
}

func NewWorkoutService(redisClient *cache.RedisClient) *WorkoutService {
	return &WorkoutService{
		cache:     redisClient,
		cacheWrap: cache.NewCacheWrapper(redisClient),
	}
}

// GetWorkoutsWithFilters получает тренировки с фильтрацией и кешированием
func (s *WorkoutService) GetWorkoutsWithFilters(filter map[string]interface{}) ([]model.Workout, error) {
	// Создаём хеш из фильтров для уникального ключа
	filterBytes, _ := json.Marshal(filter)
	hash := md5.Sum(filterBytes)
	cacheKey := "workouts:filter:" + hex.EncodeToString(hash[:])

	var workouts []model.Workout

	err := s.cacheWrap.GetOrSet(
		cacheKey,
		2*time.Minute,
		&workouts,
		func() (interface{}, error) {
			return s.getWorkoutsFromDB(filter)
		},
	)

	return workouts, err
}

// getWorkoutsFromDB имитирует получение из БД с фильтрацией
func (s *WorkoutService) getWorkoutsFromDB(filter map[string]interface{}) ([]model.Workout, error) {
	// Здесь должна быть реальная логика запроса к БД с фильтрацией
	mockWorkouts := []model.Workout{
		{
			ID:              1,
			UserID:          1,
			Type:            "Йога",
			DurationMinutes: 60,
			CaloriesBurned:  300,
			Notes:           "Утренняя практика йоги",
			WorkoutDate:     time.Now(),
			CreatedAt:       time.Now(),
		},
		{
			ID:              2,
			UserID:          2,
			Type:            "Кардио",
			DurationMinutes: 45,
			CaloriesBurned:  500,
			Notes:           "Интервальная тренировка",
			WorkoutDate:     time.Now(),
			CreatedAt:       time.Now(),
		},
	}

	// Применяем фильтры
	var result []model.Workout
	for _, w := range mockWorkouts {
		match := true

		if workoutType, ok := filter["type"]; ok && w.Type != workoutType {
			match = false
		}
		if userID, ok := filter["user_id"]; ok && w.UserID != userID {
			match = false
		}

		if match {
			result = append(result, w)
		}
	}

	return result, nil
}

// CreateWorkout создаёт тренировку
func (s *WorkoutService) CreateWorkout(workout *model.Workout) error {
	// Здесь будет сохранение в БД
	// Пока просто возвращаем nil
	return nil
}
