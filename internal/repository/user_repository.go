package repository

import (
	"fitness-club-1/internal/models"

	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	GetAll() ([]models.User, error)
	GetByID(id int) (*models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(id int) error
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetAll() ([]models.User, error) {
	var users []models.User
	query := `SELECT id, name, email, phone, created_at FROM users ORDER BY id`
	err := r.db.Select(&users, query)
	return users, err
}

func (r *userRepository) GetByID(id int) (*models.User, error) {
	var user models.User
	query := `SELECT id, name, email, phone, created_at FROM users WHERE id = $1`
	err := r.db.Get(&user, query, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(user *models.User) error {
	query := `INSERT INTO users (name, email, phone) VALUES ($1, $2, $3) RETURNING id`
	return r.db.QueryRow(query, user.Name, user.Email, user.Phone).Scan(&user.ID)
}

func (r *userRepository) Update(user *models.User) error {
	query := `UPDATE users SET name = $1, email = $2, phone = $3 WHERE id = $4`
	_, err := r.db.Exec(query, user.Name, user.Email, user.Phone, user.ID)
	return err
}

func (r *userRepository) Delete(id int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
