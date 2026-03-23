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

func TestMembershipService_GetAllMemberships_CacheHit(t *testing.T) {
	mockRepo := new(MockMembershipRepository)
	mockCache := new(MockCacheWrapper)
	service := NewMembershipService(mockRepo, mockCache)

	expectedMemberships := []model.Membership{{ID: 1, Type: "standard", Price: 1000}}

	mockCache.On("Get", "memberships:all", mock.Anything).Run(func(args mock.Arguments) {
		dest := args.Get(1).(*[]model.Membership)
		*dest = expectedMemberships
	}).Return(nil)

	memberships, err := service.GetAllMemberships()

	assert.NoError(t, err)
	assert.Equal(t, expectedMemberships, memberships)
	mockCache.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "FindAll")
}

func TestMembershipService_GetAllMemberships_CacheMiss(t *testing.T) {
	mockRepo := new(MockMembershipRepository)
	mockCache := new(MockCacheWrapper)
	service := NewMembershipService(mockRepo, mockCache)

	expectedMemberships := []model.Membership{{ID: 1, Type: "standard", Price: 1000}}

	mockCache.On("Get", "memberships:all", mock.Anything).Return(errors.New("cache miss"))
	mockRepo.On("FindAll").Return(expectedMemberships, nil)
	mockCache.On("Set", "memberships:all", expectedMemberships, 5*time.Minute).Return(nil)

	memberships, err := service.GetAllMemberships()

	assert.NoError(t, err)
	assert.Equal(t, expectedMemberships, memberships)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestMembershipService_GetAllMemberships_DBError(t *testing.T) {
	mockRepo := new(MockMembershipRepository)
	mockCache := new(MockCacheWrapper)
	service := NewMembershipService(mockRepo, mockCache)

	dbError := errors.New("database error")
	mockCache.On("Get", "memberships:all", mock.Anything).Return(errors.New("cache miss"))
	mockRepo.On("FindAll").Return([]model.Membership{}, dbError)

	memberships, err := service.GetAllMemberships()

	assert.Error(t, err)
	assert.Equal(t, dbError, err)
	assert.Empty(t, memberships)
	mockCache.AssertNotCalled(t, "Set")
}

func TestMembershipService_GetMembershipByID_CacheHit(t *testing.T) {
	mockRepo := new(MockMembershipRepository)
	mockCache := new(MockCacheWrapper)
	service := NewMembershipService(mockRepo, mockCache)

	expectedMembership := &model.Membership{ID: 1, Type: "standard", Price: 1000}

	mockCache.On("Get", "membership:1", mock.Anything).Run(func(args mock.Arguments) {
		dest := args.Get(1).(*model.Membership)
		*dest = *expectedMembership
	}).Return(nil)

	membership, err := service.GetMembershipByID(1)

	assert.NoError(t, err)
	assert.Equal(t, expectedMembership, membership)
	mockCache.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "FindByID")
}

func TestMembershipService_GetMembershipByID_CacheMiss(t *testing.T) {
	mockRepo := new(MockMembershipRepository)
	mockCache := new(MockCacheWrapper)
	service := NewMembershipService(mockRepo, mockCache)

	expectedMembership := &model.Membership{ID: 1, Type: "standard", Price: 1000}

	mockCache.On("Get", "membership:1", mock.Anything).Return(errors.New("cache miss"))
	mockRepo.On("FindByID", 1).Return(expectedMembership, nil)
	mockCache.On("Set", "membership:1", expectedMembership, 10*time.Minute).Return(nil)

	membership, err := service.GetMembershipByID(1)

	assert.NoError(t, err)
	assert.Equal(t, expectedMembership, membership)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestMembershipService_GetMembershipByID_NotFound(t *testing.T) {
	mockRepo := new(MockMembershipRepository)
	mockCache := new(MockCacheWrapper)
	service := NewMembershipService(mockRepo, mockCache)

	notFoundErr := errors.New("sql: no rows in result set")
	mockCache.On("Get", "membership:999", mock.Anything).Return(errors.New("cache miss"))
	mockRepo.On("FindByID", 999).Return(nil, notFoundErr)

	membership, err := service.GetMembershipByID(999)

	assert.Error(t, err)
	assert.Nil(t, membership)
	assert.Equal(t, notFoundErr, err)
	mockCache.AssertNotCalled(t, "Set")
}

func TestMembershipService_CreateMembership_Success(t *testing.T) {
	mockRepo := new(MockMembershipRepository)
	mockCache := new(MockCacheWrapper)
	service := NewMembershipService(mockRepo, mockCache)

	membership := &model.Membership{UserID: 1, Type: "standard", Price: 1000}

	mockRepo.On("Create", membership).Run(func(args mock.Arguments) {
		m := args.Get(0).(*model.Membership)
		m.ID = 1
	}).Return(nil)
	mockCache.On("Delete", "memberships:all").Return(nil)

	err := service.CreateMembership(membership)

	assert.NoError(t, err)
	assert.Equal(t, 1, membership.ID)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestMembershipService_CreateMembership_DBError(t *testing.T) {
	mockRepo := new(MockMembershipRepository)
	mockCache := new(MockCacheWrapper)
	service := NewMembershipService(mockRepo, mockCache)

	membership := &model.Membership{UserID: 1, Type: "standard", Price: 1000}
	dbError := errors.New("duplicate key")

	mockRepo.On("Create", membership).Return(dbError)

	err := service.CreateMembership(membership)

	assert.Error(t, err)
	assert.Equal(t, dbError, err)
	assert.Equal(t, 0, membership.ID)
	mockCache.AssertNotCalled(t, "Delete")
}

func TestMembershipService_DeleteMembership_Success(t *testing.T) {
	mockRepo := new(MockMembershipRepository)
	mockCache := new(MockCacheWrapper)
	service := NewMembershipService(mockRepo, mockCache)

	mockRepo.On("Delete", 1).Return(nil)
	mockCache.On("Delete", "membership:1").Return(nil)
	mockCache.On("Delete", "memberships:all").Return(nil)

	err := service.DeleteMembership(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestMembershipService_DeleteMembership_DBError(t *testing.T) {
	mockRepo := new(MockMembershipRepository)
	mockCache := new(MockCacheWrapper)
	service := NewMembershipService(mockRepo, mockCache)

	dbError := errors.New("membership not found")
	mockRepo.On("Delete", 999).Return(dbError)

	err := service.DeleteMembership(999)

	assert.Error(t, err)
	assert.Equal(t, dbError, err)
	mockCache.AssertNotCalled(t, "Delete")
}
