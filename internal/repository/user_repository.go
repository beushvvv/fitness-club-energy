package repository

import (
	"fitness-club-energy/internal/model"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetAll() ([]model.User, error) {
	var users []model.User
	query := `SELECT id, name, email, phone, created_at FROM users ORDER BY id`
	err := r.db.Select(&users, query)
	return users, err
}

func (r *UserRepository) GetByID(id int) (*model.User, error) {
	var user model.User
	query := `SELECT id, name, email, phone, created_at FROM users WHERE id = $1`
	err := r.db.Get(&user, query, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Create(user *model.User) error {
	query := `INSERT INTO users (name, email, phone) VALUES ($1, $2, $3) RETURNING id, created_at`
	return r.db.QueryRow(query, user.Name, user.Email, user.Phone).Scan(&user.ID, &user.CreatedAt)
}

func (r *UserRepository) Update(user *model.User) error {
	query := `UPDATE users SET name = $1, email = $2, phone = $3 WHERE id = $4`
	_, err := r.db.Exec(query, user.Name, user.Email, user.Phone, user.ID)
	return err
}

func (r *UserRepository) Delete(id int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
