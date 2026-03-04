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

// FindAll - получить всех пользователей
func (r *UserRepository) FindAll() ([]model.User, error) {
	var users []model.User
	err := r.db.Select(&users, "SELECT id, name, email, phone, created_at FROM users ORDER BY id")
	return users, err
}

// FindByID - найти пользователя по ID
func (r *UserRepository) FindByID(id int) (*model.User, error) {
	var user model.User
	err := r.db.Get(&user, "SELECT id, name, email, phone, created_at FROM users WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Create - создать пользователя
func (r *UserRepository) Create(user *model.User) error {
	query := `INSERT INTO users (name, email, phone, created_at) 
              VALUES ($1, $2, $3, NOW()) 
              RETURNING id, created_at`

	return r.db.QueryRow(query, user.Name, user.Email, user.Phone).Scan(&user.ID, &user.CreatedAt)
}

// Update - обновить пользователя
func (r *UserRepository) Update(user *model.User) error {
	query := `UPDATE users SET name = $1, email = $2, phone = $3 WHERE id = $4`
	_, err := r.db.Exec(query, user.Name, user.Email, user.Phone, user.ID)
	return err
}

// Delete - удалить пользователя
func (r *UserRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM users WHERE id = $1", id)
	return err
}
