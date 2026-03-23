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

func TestUserService_GetAllUsers_CacheHit(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockCache := new(MockCacheWrapper)
	service := NewUserService(mockRepo, mockCache)

	expectedUsers := []model.User{{ID: 1, Name: "Test", Email: "test@test.com"}}

	mockCache.On("Get", "users:all", mock.Anything).Run(func(args mock.Arguments) {
		dest := args.Get(1).(*[]model.User)
		*dest = expectedUsers
	}).Return(nil)

	users, err := service.GetAllUsers()

	assert.NoError(t, err)
	assert.Equal(t, expectedUsers, users)
	mockCache.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "FindAll")
}

func TestUserService_GetAllUsers_CacheMiss(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockCache := new(MockCacheWrapper)
	service := NewUserService(mockRepo, mockCache)

	expectedUsers := []model.User{{ID: 1, Name: "Test", Email: "test@test.com"}}

	mockCache.On("Get", "users:all", mock.Anything).Return(errors.New("cache miss"))
	mockRepo.On("FindAll").Return(expectedUsers, nil)
	mockCache.On("Set", "users:all", expectedUsers, 5*time.Minute).Return(nil)

	users, err := service.GetAllUsers()

	assert.NoError(t, err)
	assert.Equal(t, expectedUsers, users)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestUserService_GetAllUsers_DBError(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockCache := new(MockCacheWrapper)
	service := NewUserService(mockRepo, mockCache)

	dbError := errors.New("database error")
	mockCache.On("Get", "users:all", mock.Anything).Return(errors.New("cache miss"))
	mockRepo.On("FindAll").Return([]model.User{}, dbError)

	users, err := service.GetAllUsers()

	assert.Error(t, err)
	assert.Equal(t, dbError, err)
	assert.Empty(t, users)
	mockCache.AssertNotCalled(t, "Set")
}

func TestUserService_GetUserByID_CacheHit(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockCache := new(MockCacheWrapper)
	service := NewUserService(mockRepo, mockCache)

	expectedUser := &model.User{ID: 1, Name: "Test", Email: "test@test.com"}

	mockCache.On("Get", "user:1", mock.Anything).Run(func(args mock.Arguments) {
		dest := args.Get(1).(*model.User)
		*dest = *expectedUser
	}).Return(nil)

	user, err := service.GetUserByID(1)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockCache.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "FindByID")
}

func TestUserService_GetUserByID_CacheMiss(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockCache := new(MockCacheWrapper)
	service := NewUserService(mockRepo, mockCache)

	expectedUser := &model.User{ID: 1, Name: "Test", Email: "test@test.com"}

	mockCache.On("Get", "user:1", mock.Anything).Return(errors.New("cache miss"))
	mockRepo.On("FindByID", 1).Return(expectedUser, nil)
	mockCache.On("Set", "user:1", expectedUser, 10*time.Minute).Return(nil)

	user, err := service.GetUserByID(1)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestUserService_GetUserByID_NotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockCache := new(MockCacheWrapper)
	service := NewUserService(mockRepo, mockCache)

	notFoundErr := errors.New("sql: no rows in result set")
	mockCache.On("Get", "user:999", mock.Anything).Return(errors.New("cache miss"))
	mockRepo.On("FindByID", 999).Return(nil, notFoundErr)

	user, err := service.GetUserByID(999)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, notFoundErr, err)
	mockCache.AssertNotCalled(t, "Set")
}

func TestUserService_CreateUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockCache := new(MockCacheWrapper)
	service := NewUserService(mockRepo, mockCache)

	user := &model.User{Name: "New User", Email: "new@test.com"}

	mockRepo.On("Create", user).Run(func(args mock.Arguments) {
		u := args.Get(0).(*model.User)
		u.ID = 1
	}).Return(nil)
	mockCache.On("Delete", "users:all").Return(nil)

	err := service.CreateUser(user)

	assert.NoError(t, err)
	assert.Equal(t, 1, user.ID)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestUserService_CreateUser_DBError(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockCache := new(MockCacheWrapper)
	service := NewUserService(mockRepo, mockCache)

	user := &model.User{Name: "New User", Email: "new@test.com"}
	dbError := errors.New("duplicate key value violates unique constraint")

	mockRepo.On("Create", user).Return(dbError)

	err := service.CreateUser(user)

	assert.Error(t, err)
	assert.Equal(t, dbError, err)
	assert.Equal(t, 0, user.ID)
	mockCache.AssertNotCalled(t, "Delete")
}

func TestUserService_DeleteUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockCache := new(MockCacheWrapper)
	service := NewUserService(mockRepo, mockCache)

	mockRepo.On("Delete", 1).Return(nil)
	mockCache.On("Delete", "user:1").Return(nil)
	mockCache.On("Delete", "users:all").Return(nil)

	err := service.DeleteUser(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestUserService_DeleteUser_DBError(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockCache := new(MockCacheWrapper)
	service := NewUserService(mockRepo, mockCache)

	dbError := errors.New("user not found")
	mockRepo.On("Delete", 999).Return(dbError)

	err := service.DeleteUser(999)

	assert.Error(t, err)
	assert.Equal(t, dbError, err)
	mockCache.AssertNotCalled(t, "Delete")
}
