package service

import (
	"strconv"
	"time"

	"fitness-club-energy/internal/cache"
	"fitness-club-energy/internal/logger"
	"fitness-club-energy/internal/model"
	"fitness-club-energy/internal/repository"
)

type MembershipService struct {
	membershipRepo repository.MembershipRepositoryInterface
	cacheWrapper   cache.CacheWrapperInterface
}

func NewMembershipService(membershipRepo repository.MembershipRepositoryInterface, cacheWrapper cache.CacheWrapperInterface) *MembershipService {
	return &MembershipService{
		membershipRepo: membershipRepo,
		cacheWrapper:   cacheWrapper,
	}
}

// GetAllMemberships с кэшированием
func (s *MembershipService) GetAllMemberships() ([]model.Membership, error) {
	var memberships []model.Membership
	sugar := logger.Log.Sugar()

	err := s.cacheWrapper.Get("memberships:all", &memberships)
	if err == nil {
		sugar.Debugw("📦 MEMBERSHIPS FROM CACHE", "key", "memberships:all")
		return memberships, nil
	}

	sugar.Debug("💾 MEMBERSHIPS FROM DATABASE")
	memberships, err = s.membershipRepo.FindAll()
	if err != nil {
		sugar.Errorw("Failed to get memberships from DB", "error", err)
		return nil, err
	}

	if len(memberships) > 0 {
		s.cacheWrapper.Set("memberships:all", memberships, 5*time.Minute)
		sugar.Debugw("✅ Memberships cached", "count", len(memberships))
	}

	return memberships, nil
}

// GetMembershipByID с кэшированием
func (s *MembershipService) GetMembershipByID(id int) (*model.Membership, error) {
	var membership model.Membership
	sugar := logger.Log.Sugar()
	key := "membership:" + strconv.Itoa(id)

	err := s.cacheWrapper.Get(key, &membership)
	if err == nil {
		sugar.Debugw("📦 MEMBERSHIP FROM CACHE", "membership_id", id)
		return &membership, nil
	}

	sugar.Debugw("💾 MEMBERSHIP FROM DATABASE", "membership_id", id)
	membershipFromDB, err := s.membershipRepo.FindByID(id)
	if err != nil {
		sugar.Errorw("Failed to get membership by ID from DB", "membership_id", id, "error", err)
		return nil, err
	}

	s.cacheWrapper.Set(key, membershipFromDB, 10*time.Minute)
	sugar.Debugw("✅ Membership cached", "membership_id", id)

	return membershipFromDB, nil
}

// CreateMembership создаёт абонемент
func (s *MembershipService) CreateMembership(membership *model.Membership) error {
	sugar := logger.Log.Sugar()

	if err := s.membershipRepo.Create(membership); err != nil {
		sugar.Errorw("Failed to create membership",
			"error", err,
			"user_id", membership.UserID,
			"type", membership.Type)
		return err
	}

	s.cacheWrapper.Delete("memberships:all")
	sugar.Infow("Membership created and cache cleared",
		"membership_id", membership.ID,
		"user_id", membership.UserID,
		"type", membership.Type)

	return nil
}

// UpdateMembership обновляет абонемент
func (s *MembershipService) UpdateMembership(membership *model.Membership) error {
	sugar := logger.Log.Sugar()

	if err := s.membershipRepo.Update(membership); err != nil {
		sugar.Errorw("Failed to update membership",
			"membership_id", membership.ID,
			"error", err)
		return err
	}

	key := "membership:" + strconv.Itoa(membership.ID)
	s.cacheWrapper.Delete(key)
	s.cacheWrapper.Delete("memberships:all")
	sugar.Infow("Membership updated and cache cleared",
		"membership_id", membership.ID,
		"type", membership.Type)

	return nil
}

// DeleteMembership удаляет абонемент
func (s *MembershipService) DeleteMembership(id int) error {
	sugar := logger.Log.Sugar()

	if err := s.membershipRepo.Delete(id); err != nil {
		sugar.Errorw("Failed to delete membership",
			"membership_id", id,
			"error", err)
		return err
	}

	key := "membership:" + strconv.Itoa(id)
	s.cacheWrapper.Delete(key)
	s.cacheWrapper.Delete("memberships:all")
	sugar.Infow("Membership deleted and cache cleared", "membership_id", id)

	return nil
}
