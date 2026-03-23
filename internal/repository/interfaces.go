package repository

import "fitness-club-energy/internal/model"

type UserRepositoryInterface interface {
	FindAll() ([]model.User, error)
	FindByID(id int) (*model.User, error)
	Create(user *model.User) error
	Update(user *model.User) error
	Delete(id int) error
}

type MembershipRepositoryInterface interface {
	FindAll() ([]model.Membership, error)
	FindByID(id int) (*model.Membership, error)
	Create(membership *model.Membership) error
	Update(membership *model.Membership) error
	Delete(id int) error
}

type WorkoutRepositoryInterface interface {
	FindAll() ([]model.Workout, error)
	FindByID(id int) (*model.Workout, error)
	Create(workout *model.Workout) error
	Update(workout *model.Workout) error
	Delete(id int) error
}
