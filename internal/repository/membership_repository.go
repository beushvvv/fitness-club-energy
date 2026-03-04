package repository

import (
	"fitness-club-energy/internal/model"

	"github.com/jmoiron/sqlx"
)

type MembershipRepository struct {
	db *sqlx.DB
}

func NewMembershipRepository(db *sqlx.DB) *MembershipRepository {
	return &MembershipRepository{db: db}
}

// FindAll - получить все абонементы
func (r *MembershipRepository) FindAll() ([]model.Membership, error) {
	var memberships []model.Membership
	query := `SELECT id, user_id, type, price, start_date, end_date, is_active, created_at 
              FROM memberships ORDER BY id`
	err := r.db.Select(&memberships, query)
	return memberships, err
}

// FindByID - найти абонемент по ID
func (r *MembershipRepository) FindByID(id int) (*model.Membership, error) {
	var membership model.Membership
	query := `SELECT id, user_id, type, price, start_date, end_date, is_active, created_at 
              FROM memberships WHERE id = $1`
	err := r.db.Get(&membership, query, id)
	if err != nil {
		return nil, err
	}
	return &membership, nil
}

// Create - создать абонемент
func (r *MembershipRepository) Create(membership *model.Membership) error {
	query := `INSERT INTO memberships (user_id, type, price, start_date, end_date, is_active, created_at) 
              VALUES ($1, $2, $3, $4, $5, $6, NOW()) 
              RETURNING id, created_at`

	return r.db.QueryRow(
		query,
		membership.UserID,
		membership.Type,
		membership.Price,
		membership.StartDate,
		membership.EndDate,
		membership.IsActive,
	).Scan(&membership.ID, &membership.CreatedAt)
}

// Update - обновить абонемент
func (r *MembershipRepository) Update(membership *model.Membership) error {
	query := `UPDATE memberships 
              SET type = $1, price = $2, start_date = $3, end_date = $4, is_active = $5 
              WHERE id = $6`
	_, err := r.db.Exec(
		query,
		membership.Type,
		membership.Price,
		membership.StartDate,
		membership.EndDate,
		membership.IsActive,
		membership.ID,
	)
	return err
}

// Delete удаляет абонемент по ID
func (r *MembershipRepository) Delete(id int) error {
	query := `DELETE FROM memberships WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
