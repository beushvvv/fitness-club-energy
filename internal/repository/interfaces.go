package repository

import "fitness-club-energy/internal/model"

type UserRepositoryInterface interface {
	FindAll() ([]model.User, error)
	FindByID(id int) (*model.User, error)
	Create(user *model.User) error
	Update(user *model.User) error
	Delete(id int) error
}
