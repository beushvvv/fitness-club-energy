package service

import (
	"fitness-club-energy/internal/model"
	"fitness-club-energy/internal/repository"
)

type MembershipService struct {
	membershipRepo *repository.MembershipRepository
}

func NewMembershipService(membershipRepo *repository.MembershipRepository) *MembershipService {
	return &MembershipService{membershipRepo: membershipRepo}
}

func (s *MembershipService) GetAllMemberships() ([]model.Membership, error) {
	return s.membershipRepo.GetAll()
}

func (s *MembershipService) GetMembershipsByUser(userID int) ([]model.Membership, error) {
	return s.membershipRepo.GetByUserID(userID)
}

func (s *MembershipService) CreateMembership(membership *model.Membership) error {
	return s.membershipRepo.Create(membership)
}
