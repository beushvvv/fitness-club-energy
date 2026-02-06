package service

import (
	"fitness-club-1/internal/models"
	"fitness-club-1/internal/repository"
)

type MembershipService struct {
	membershipRepo repository.MembershipRepository
}

func NewMembershipService(membershipRepo repository.MembershipRepository) *MembershipService {
	return &MembershipService{membershipRepo: membershipRepo}
}

func (s *MembershipService) GetAllMemberships() ([]models.Membership, error) {
	return s.membershipRepo.GetAll()
}

func (s *MembershipService) GetMembershipsByUser(userID int) ([]models.Membership, error) {
	return s.membershipRepo.GetByUserID(userID)
}

func (s *MembershipService) CreateMembership(membership *models.Membership) error {
	return s.membershipRepo.Create(membership)
}
