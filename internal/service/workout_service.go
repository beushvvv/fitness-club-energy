package service

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"strconv"
	"time"

	"fitness-club-energy/internal/cache"
	"fitness-club-energy/internal/logger"
	"fitness-club-energy/internal/model"
	"fitness-club-energy/internal/repository"
)

type WorkoutService struct {
	workoutRepo  repository.WorkoutRepositoryInterface
	cacheWrapper cache.CacheWrapperInterface
}

func NewWorkoutService(workoutRepo repository.WorkoutRepositoryInterface, cacheWrapper cache.CacheWrapperInterface) *WorkoutService {
	return &WorkoutService{
		workoutRepo:  workoutRepo,
		cacheWrapper: cacheWrapper,
	}
}

// GetAllWorkouts получает все тренировки с кэшированием
func (s *WorkoutService) GetAllWorkouts() ([]model.Workout, error) {
	var workouts []model.Workout
	sugar := logger.Log.Sugar()

	err := s.cacheWrapper.Get("workouts:all", &workouts)
	if err == nil {
		sugar.Debugw("📦 WORKOUTS FROM CACHE", "key", "workouts:all")
		return workouts, nil
	}

	sugar.Debug("💾 WORKOUTS FROM DATABASE")
	workouts, err = s.workoutRepo.FindAll()
	if err != nil {
		sugar.Errorw("Failed to get workouts from DB", "error", err)
		return nil, err
	}

	if len(workouts) > 0 {
		s.cacheWrapper.Set("workouts:all", workouts, 5*time.Minute)
		sugar.Debugw("✅ Workouts cached", "count", len(workouts))
	}

	return workouts, nil
}

// GetWorkoutByID получает тренировку по ID с кэшированием
func (s *WorkoutService) GetWorkoutByID(id int) (*model.Workout, error) {
	var workout model.Workout
	sugar := logger.Log.Sugar()
	key := "workout:" + strconv.Itoa(id)

	err := s.cacheWrapper.Get(key, &workout)
	if err == nil {
		sugar.Debugw("📦 WORKOUT FROM CACHE", "workout_id", id)
		return &workout, nil
	}

	sugar.Debugw("💾 WORKOUT FROM DATABASE", "workout_id", id)
	workoutFromDB, err := s.workoutRepo.FindByID(id)
	if err != nil {
		sugar.Errorw("Failed to get workout by ID from DB",
			"workout_id", id,
			"error", err)
		return nil, err
	}

	s.cacheWrapper.Set(key, workoutFromDB, 10*time.Minute)
	sugar.Debugw("✅ Workout cached", "workout_id", id)

	return workoutFromDB, nil
}

// GetWorkoutsWithFilters получает тренировки с фильтрацией и кэшированием
func (s *WorkoutService) GetWorkoutsWithFilters(filter map[string]interface{}) ([]model.Workout, error) {
	sugar := logger.Log.Sugar()

	// Создаём хеш из фильтров для уникального ключа кэша
	filterBytes, _ := json.Marshal(filter)
	hash := md5.Sum(filterBytes)
	cacheKey := "workouts:filter:" + hex.EncodeToString(hash[:])

	var workouts []model.Workout

	err := s.cacheWrapper.Get(cacheKey, &workouts)
	if err == nil {
		sugar.Debugw("📦 FILTERED WORKOUTS FROM CACHE",
			"filter", filter,
			"cache_key", cacheKey)
		return workouts, nil
	}

	sugar.Debugw("💾 FILTERED WORKOUTS FROM DATABASE", "filter", filter)
	// В реальном приложении здесь должен быть запрос к БД с фильтрацией
	workouts, err = s.workoutRepo.FindAll()
	if err != nil {
		sugar.Errorw("Failed to get filtered workouts from DB",
			"filter", filter,
			"error", err)
		return nil, err
	}

	// Применяем фильтрацию в памяти (для демонстрации)
	var filtered []model.Workout
	for _, w := range workouts {
		match := true
		if workoutType, ok := filter["type"]; ok && w.Type != workoutType {
			match = false
		}
		if userID, ok := filter["user_id"]; ok && w.UserID != userID {
			match = false
		}
		if match {
			filtered = append(filtered, w)
		}
	}

	if len(filtered) > 0 {
		s.cacheWrapper.Set(cacheKey, filtered, 2*time.Minute)
		sugar.Debugw("✅ Filtered workouts cached",
			"count", len(filtered),
			"cache_key", cacheKey)
	}

	return filtered, nil
}

// CreateWorkout создаёт новую тренировку
func (s *WorkoutService) CreateWorkout(workout *model.Workout) error {
	sugar := logger.Log.Sugar()

	if err := s.workoutRepo.Create(workout); err != nil {
		sugar.Errorw("Failed to create workout",
			"error", err,
			"user_id", workout.UserID,
			"type", workout.Type)
		return err
	}

	// Очищаем кэш списка тренировок
	s.cacheWrapper.Delete("workouts:all")
	sugar.Infow("Workout created and cache cleared",
		"workout_id", workout.ID,
		"user_id", workout.UserID,
		"type", workout.Type)

	return nil
}

// UpdateWorkout обновляет тренировку
func (s *WorkoutService) UpdateWorkout(workout *model.Workout) error {
	sugar := logger.Log.Sugar()

	if err := s.workoutRepo.Update(workout); err != nil {
		sugar.Errorw("Failed to update workout",
			"workout_id", workout.ID,
			"error", err)
		return err
	}

	// Очищаем кэш
	key := "workout:" + strconv.Itoa(workout.ID)
	s.cacheWrapper.Delete(key)
	s.cacheWrapper.Delete("workouts:all")
	sugar.Infow("Workout updated and cache cleared",
		"workout_id", workout.ID,
		"type", workout.Type)

	return nil
}

// DeleteWorkout удаляет тренировку
func (s *WorkoutService) DeleteWorkout(id int) error {
	sugar := logger.Log.Sugar()

	if err := s.workoutRepo.Delete(id); err != nil {
		sugar.Errorw("Failed to delete workout",
			"workout_id", id,
			"error", err)
		return err
	}

	// Очищаем кэш
	key := "workout:" + strconv.Itoa(id)
	s.cacheWrapper.Delete(key)
	s.cacheWrapper.Delete("workouts:all")
	sugar.Infow("Workout deleted and cache cleared", "workout_id", id)

	return nil
}
