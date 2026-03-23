package service

import (
	"errors"
	"testing"
	"time"

	"fitness-club-energy/internal/logger"
	"fitness-club-energy/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func init() {
	logger.Init("debug")
}

func TestWorkoutService_GetAllWorkouts_CacheHit(t *testing.T) {
	mockRepo := new(MockWorkoutRepository)
	mockCache := new(MockCacheWrapper)
	service := NewWorkoutService(mockRepo, mockCache)

	expectedWorkouts := []model.Workout{{ID: 1, Type: "cardio", DurationMinutes: 30}}

	mockCache.On("Get", "workouts:all", mock.Anything).Run(func(args mock.Arguments) {
		dest := args.Get(1).(*[]model.Workout)
		*dest = expectedWorkouts
	}).Return(nil)

	workouts, err := service.GetAllWorkouts()

	assert.NoError(t, err)
	assert.Equal(t, expectedWorkouts, workouts)
	mockCache.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "FindAll")
}

func TestWorkoutService_GetAllWorkouts_CacheMiss(t *testing.T) {
	mockRepo := new(MockWorkoutRepository)
	mockCache := new(MockCacheWrapper)
	service := NewWorkoutService(mockRepo, mockCache)

	expectedWorkouts := []model.Workout{{ID: 1, Type: "cardio", DurationMinutes: 30}}

	mockCache.On("Get", "workouts:all", mock.Anything).Return(errors.New("cache miss"))
	mockRepo.On("FindAll").Return(expectedWorkouts, nil)
	mockCache.On("Set", "workouts:all", expectedWorkouts, 5*time.Minute).Return(nil)

	workouts, err := service.GetAllWorkouts()

	assert.NoError(t, err)
	assert.Equal(t, expectedWorkouts, workouts)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestWorkoutService_GetAllWorkouts_DBError(t *testing.T) {
	mockRepo := new(MockWorkoutRepository)
	mockCache := new(MockCacheWrapper)
	service := NewWorkoutService(mockRepo, mockCache)

	dbError := errors.New("database error")
	mockCache.On("Get", "workouts:all", mock.Anything).Return(errors.New("cache miss"))
	mockRepo.On("FindAll").Return([]model.Workout{}, dbError)

	workouts, err := service.GetAllWorkouts()

	assert.Error(t, err)
	assert.Equal(t, dbError, err)
	assert.Empty(t, workouts)
	mockCache.AssertNotCalled(t, "Set")
}

func TestWorkoutService_GetWorkoutByID_CacheHit(t *testing.T) {
	mockRepo := new(MockWorkoutRepository)
	mockCache := new(MockCacheWrapper)
	service := NewWorkoutService(mockRepo, mockCache)

	expectedWorkout := &model.Workout{ID: 1, Type: "cardio", DurationMinutes: 30}

	mockCache.On("Get", "workout:1", mock.Anything).Run(func(args mock.Arguments) {
		dest := args.Get(1).(*model.Workout)
		*dest = *expectedWorkout
	}).Return(nil)

	workout, err := service.GetWorkoutByID(1)

	assert.NoError(t, err)
	assert.Equal(t, expectedWorkout, workout)
	mockCache.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "FindByID")
}

func TestWorkoutService_GetWorkoutByID_CacheMiss(t *testing.T) {
	mockRepo := new(MockWorkoutRepository)
	mockCache := new(MockCacheWrapper)
	service := NewWorkoutService(mockRepo, mockCache)

	expectedWorkout := &model.Workout{ID: 1, Type: "cardio", DurationMinutes: 30}

	mockCache.On("Get", "workout:1", mock.Anything).Return(errors.New("cache miss"))
	mockRepo.On("FindByID", 1).Return(expectedWorkout, nil)
	mockCache.On("Set", "workout:1", expectedWorkout, 10*time.Minute).Return(nil)

	workout, err := service.GetWorkoutByID(1)

	assert.NoError(t, err)
	assert.Equal(t, expectedWorkout, workout)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestWorkoutService_GetWorkoutByID_NotFound(t *testing.T) {
	mockRepo := new(MockWorkoutRepository)
	mockCache := new(MockCacheWrapper)
	service := NewWorkoutService(mockRepo, mockCache)

	notFoundErr := errors.New("sql: no rows in result set")
	mockCache.On("Get", "workout:999", mock.Anything).Return(errors.New("cache miss"))
	mockRepo.On("FindByID", 999).Return(nil, notFoundErr)

	workout, err := service.GetWorkoutByID(999)

	assert.Error(t, err)
	assert.Nil(t, workout)
	assert.Equal(t, notFoundErr, err)
	mockCache.AssertNotCalled(t, "Set")
}

func TestWorkoutService_CreateWorkout_Success(t *testing.T) {
	mockRepo := new(MockWorkoutRepository)
	mockCache := new(MockCacheWrapper)
	service := NewWorkoutService(mockRepo, mockCache)

	workout := &model.Workout{UserID: 1, Type: "cardio", DurationMinutes: 30}

	mockRepo.On("Create", workout).Run(func(args mock.Arguments) {
		w := args.Get(0).(*model.Workout)
		w.ID = 1
	}).Return(nil)
	mockCache.On("Delete", "workouts:all").Return(nil)

	err := service.CreateWorkout(workout)

	assert.NoError(t, err)
	assert.Equal(t, 1, workout.ID)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestWorkoutService_CreateWorkout_DBError(t *testing.T) {
	mockRepo := new(MockWorkoutRepository)
	mockCache := new(MockCacheWrapper)
	service := NewWorkoutService(mockRepo, mockCache)

	workout := &model.Workout{UserID: 1, Type: "cardio", DurationMinutes: 30}
	dbError := errors.New("duplicate key")

	mockRepo.On("Create", workout).Return(dbError)

	err := service.CreateWorkout(workout)

	assert.Error(t, err)
	assert.Equal(t, dbError, err)
	assert.Equal(t, 0, workout.ID)
	mockCache.AssertNotCalled(t, "Delete")
}

func TestWorkoutService_DeleteWorkout_Success(t *testing.T) {
	mockRepo := new(MockWorkoutRepository)
	mockCache := new(MockCacheWrapper)
	service := NewWorkoutService(mockRepo, mockCache)

	mockRepo.On("Delete", 1).Return(nil)
	mockCache.On("Delete", "workout:1").Return(nil)
	mockCache.On("Delete", "workouts:all").Return(nil)

	err := service.DeleteWorkout(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestWorkoutService_DeleteWorkout_DBError(t *testing.T) {
	mockRepo := new(MockWorkoutRepository)
	mockCache := new(MockCacheWrapper)
	service := NewWorkoutService(mockRepo, mockCache)

	dbError := errors.New("workout not found")
	mockRepo.On("Delete", 999).Return(dbError)

	err := service.DeleteWorkout(999)

	assert.Error(t, err)
	assert.Equal(t, dbError, err)
	mockCache.AssertNotCalled(t, "Delete")
}
