package service

import (
	"time"

	"fitness-club-energy/internal/model"

	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindAll() ([]model.User, error) {
	args := m.Called()
	return args.Get(0).([]model.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(id int) (*model.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) Create(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Update(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

type MockCacheWrapper struct {
	mock.Mock
}

func (m *MockCacheWrapper) Get(key string, dest interface{}) error {
	args := m.Called(key, dest)
	return args.Error(0)
}

func (m *MockCacheWrapper) Set(key string, value interface{}, ttl time.Duration) error {
	args := m.Called(key, value, ttl)
	return args.Error(0)
}

func (m *MockCacheWrapper) Delete(key string) error {
	args := m.Called(key)
	return args.Error(0)
}

// MockMembershipRepository
type MockMembershipRepository struct {
	mock.Mock
}

func (m *MockMembershipRepository) FindAll() ([]model.Membership, error) {
	args := m.Called()
	return args.Get(0).([]model.Membership), args.Error(1)
}

func (m *MockMembershipRepository) FindByID(id int) (*model.Membership, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Membership), args.Error(1)
}

func (m *MockMembershipRepository) Create(membership *model.Membership) error {
	args := m.Called(membership)
	return args.Error(0)
}

func (m *MockMembershipRepository) Update(membership *model.Membership) error {
	args := m.Called(membership)
	return args.Error(0)
}

func (m *MockMembershipRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

// MockWorkoutRepository
type MockWorkoutRepository struct {
	mock.Mock
}

func (m *MockWorkoutRepository) FindAll() ([]model.Workout, error) {
	args := m.Called()
	return args.Get(0).([]model.Workout), args.Error(1)
}

func (m *MockWorkoutRepository) FindByID(id int) (*model.Workout, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Workout), args.Error(1)
}

func (m *MockWorkoutRepository) Create(workout *model.Workout) error {
	args := m.Called(workout)
	return args.Error(0)
}

func (m *MockWorkoutRepository) Update(workout *model.Workout) error {
	args := m.Called(workout)
	return args.Error(0)
}

func (m *MockWorkoutRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}
