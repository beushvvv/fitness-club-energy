package service

import (
	"fitness-club-energy/internal/model"
	"fitness-club-energy/internal/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetAllUsers() ([]model.User, error) {
	return s.userRepo.GetAll()
}

func (s *UserService) GetUserByID(id int) (*model.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *UserService) CreateUser(user *model.User) error {
	return s.userRepo.Create(user)
}

func (s *UserService) UpdateUser(user *model.User) error {
	return s.userRepo.Update(user)
}

func (s *UserService) DeleteUser(id int) error {
	return s.userRepo.Delete(id)
}
