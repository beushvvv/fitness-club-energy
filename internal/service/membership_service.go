package service

import (
	"fitness-club-energy/internal/cache"
	"fitness-club-energy/internal/model"
	"fitness-club-energy/internal/repository"
	"time"
)

type MembershipService struct {
	membershipRepo *repository.MembershipRepository
	cache          *cache.RedisClient
	cacheWrap      *cache.CacheWrapper
}

func NewMembershipService(membershipRepo *repository.MembershipRepository, redisClient *cache.RedisClient) *MembershipService {
	return &MembershipService{
		membershipRepo: membershipRepo,
		cache:          redisClient,
		cacheWrap:      cache.NewCacheWrapper(redisClient),
	}
}

// GetAllMemberships с кешированием
func (s *MembershipService) GetAllMemberships() ([]model.Membership, error) {
	var memberships []model.Membership

	err := s.cacheWrap.GetOrSet(
		cache.MembershipsAllPrefix,
		5*time.Minute,
		&memberships,
		func() (interface{}, error) {
			return s.membershipRepo.FindAll()
		},
	)

	return memberships, err
}

// GetMembershipByID с кешированием
func (s *MembershipService) GetMembershipByID(id int) (*model.Membership, error) {
	var membership model.Membership

	err := s.cacheWrap.GetOrSet(
		cache.MembershipKey(id),
		10*time.Minute,
		&membership,
		func() (interface{}, error) {
			return s.membershipRepo.FindByID(id)
		},
	)

	if err != nil {
		return nil, err
	}
	return &membership, nil
}

// CreateMembership создаёт абонемент и очищает кеш
func (s *MembershipService) CreateMembership(membership *model.Membership) error {
	if err := s.membershipRepo.Create(membership); err != nil {
		return err
	}

	// Очищаем кеш списка абонементов
	for _, key := range cache.KeysToInvalidate("membership") {
		_ = s.cache.Delete(key)
	}

	return nil
}

// UpdateMembership обновляет абонемент и очищает кеш
func (s *MembershipService) UpdateMembership(membership *model.Membership) error {
	if err := s.membershipRepo.Update(membership); err != nil {
		return err
	}

	// Очищаем кеш
	_ = s.cache.Delete(cache.MembershipKey(membership.ID))
	for _, key := range cache.KeysToInvalidate("membership") {
		_ = s.cache.Delete(key)
	}

	return nil
}
