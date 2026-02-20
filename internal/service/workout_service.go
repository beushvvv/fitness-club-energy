package service

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fitness-club-energy/internal/cache"
	"fitness-club-energy/internal/model"
	"fitness-club-energy/internal/repository"
	"log"
	"time"
)

type WorkoutService struct {
	workoutRepo  *repository.WorkoutRepository
	cacheWrapper *cache.CacheWrapper
}

// ИСПРАВЛЕНО: принимает WorkoutRepository и CacheWrapper
func NewWorkoutService(workoutRepo *repository.WorkoutRepository, cacheWrapper *cache.CacheWrapper) *WorkoutService {
	return &WorkoutService{
		workoutRepo:  workoutRepo,
		cacheWrapper: cacheWrapper,
	}
}

// GetAllWorkouts с ЛОГАМИ
func (s *WorkoutService) GetAllWorkouts() ([]model.Workout, error) {
	var workouts []model.Workout

	err := s.cacheWrapper.Get("workouts:all", &workouts)
	if err == nil {
		log.Println("📦 WORKOUTS FROM CACHE")
		return workouts, nil
	}

	log.Println("💾 WORKOUTS FROM DATABASE")
	workouts, err = s.workoutRepo.FindAll()
	if err != nil {
		return nil, err
	}

	if len(workouts) > 0 {
		s.cacheWrapper.Set("workouts:all", workouts, 5*time.Minute)
	}

	return workouts, nil
}

// GetWorkoutByID получает тренировку по ID с кешированием
func (s *WorkoutService) GetWorkoutByID(id int) (*model.Workout, error) {
	var workout model.Workout
	key := "workout:" + string(rune(id))

	err := s.cacheWrapper.GetOrSet(
		key,
		10*time.Minute,
		&workout,
		func() (interface{}, error) {
			return s.workoutRepo.FindByID(id)
		},
	)

	if err != nil {
		return nil, err
	}
	return &workout, nil
}

// GetWorkoutsWithFilters получает тренировки с фильтрацией и кешированием
func (s *WorkoutService) GetWorkoutsWithFilters(filter map[string]interface{}) ([]model.Workout, error) {
	// Создаём хеш из фильтров для уникального ключа
	filterBytes, _ := json.Marshal(filter)
	hash := md5.Sum(filterBytes)
	cacheKey := "workouts:filter:" + hex.EncodeToString(hash[:])

	var workouts []model.Workout

	err := s.cacheWrapper.GetOrSet(
		cacheKey,
		2*time.Minute,
		&workouts,
		func() (interface{}, error) {
			return s.getWorkoutsFromDB(filter)
		},
	)

	return workouts, err
}

// CreateWorkout создаёт тренировку
func (s *WorkoutService) CreateWorkout(workout *model.Workout) error {
	if err := s.workoutRepo.Create(workout); err != nil {
		return err
	}

	// Очищаем кеш списка
	s.cacheWrapper.Delete("workouts:all")
	return nil
}

// UpdateWorkout обновляет тренировку
func (s *WorkoutService) UpdateWorkout(workout *model.Workout) error {
	if err := s.workoutRepo.Update(workout); err != nil {
		return err
	}

	// Очищаем кеш
	s.cacheWrapper.Delete("workout:" + string(rune(workout.ID)))
	s.cacheWrapper.Delete("workouts:all")
	return nil
}

// DeleteWorkout удаляет тренировку
func (s *WorkoutService) DeleteWorkout(id int) error {
	if err := s.workoutRepo.Delete(id); err != nil {
		return err
	}

	// Очищаем кеш
	s.cacheWrapper.Delete("workout:" + string(rune(id)))
	s.cacheWrapper.Delete("workouts:all")
	return nil
}

// getWorkoutsFromDB получает тренировки из БД с фильтрацией
func (s *WorkoutService) getWorkoutsFromDB(filter map[string]interface{}) ([]model.Workout, error) {
	// Здесь должна быть реальная логика запроса к БД с фильтрацией
	// Пока возвращаем все
	return s.workoutRepo.FindAll()
}
