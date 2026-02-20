package service

import (
	"fitness-club-energy/internal/cache"
	"fitness-club-energy/internal/model"
	"fitness-club-energy/internal/repository"
	"log"
	"time"
)

type MembershipService struct {
	membershipRepo *repository.MembershipRepository
	cacheWrapper   *cache.CacheWrapper
}

// ИСПРАВЛЕНО: принимает CacheWrapper вместо RedisClient
func NewMembershipService(membershipRepo *repository.MembershipRepository, cacheWrapper *cache.CacheWrapper) *MembershipService {
	return &MembershipService{
		membershipRepo: membershipRepo,
		cacheWrapper:   cacheWrapper,
	}
}

// GetAllMemberships с ЛОГАМИ
func (s *MembershipService) GetAllMemberships() ([]model.Membership, error) {
	var memberships []model.Membership

	err := s.cacheWrapper.Get("memberships:all", &memberships)
	if err == nil {
		log.Println("📦 MEMBERSHIPS FROM CACHE")
		return memberships, nil
	}

	log.Println("💾 MEMBERSHIPS FROM DATABASE")
	memberships, err = s.membershipRepo.FindAll()
	if err != nil {
		return nil, err
	}

	if len(memberships) > 0 {
		s.cacheWrapper.Set("memberships:all", memberships, 5*time.Minute)
	}

	return memberships, nil
}

// GetMembershipByID с кешированием
func (s *MembershipService) GetMembershipByID(id int) (*model.Membership, error) {
	var membership model.Membership
	key := "membership:" + string(rune(id))

	err := s.cacheWrapper.GetOrSet(
		key,
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

	// Очищаем кеш списка
	s.cacheWrapper.Delete("memberships:all")
	return nil
}

// UpdateMembership обновляет абонемент и очищает кеш
func (s *MembershipService) UpdateMembership(membership *model.Membership) error {
	if err := s.membershipRepo.Update(membership); err != nil {
		return err
	}

	// Очищаем кеш
	s.cacheWrapper.Delete("membership:" + string(rune(membership.ID)))
	s.cacheWrapper.Delete("memberships:all")
	return nil
}
