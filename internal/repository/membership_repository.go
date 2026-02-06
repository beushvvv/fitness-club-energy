package repository

import (
	"fitness-club-1/internal/models"

	"github.com/jmoiron/sqlx"
)

type MembershipRepository interface {
	GetAll() ([]models.Membership, error)
	GetByUserID(userID int) ([]models.Membership, error)
	Create(membership *models.Membership) error
}

type membershipRepository struct {
	db *sqlx.DB
}

func NewMembershipRepository(db *sqlx.DB) MembershipRepository {
	return &membershipRepository{db: db}
}

func (r *membershipRepository) GetAll() ([]models.Membership, error) {
	var memberships []models.Membership
	query := `
        SELECT m.id, m.type, m.price, m.start_date, m.end_date, m.is_active, 
               m.user_id, m.created_at 
        FROM memberships m
        ORDER BY m.created_at DESC
    `
	err := r.db.Select(&memberships, query)
	return memberships, err
}

func (r *membershipRepository) GetByUserID(userID int) ([]models.Membership, error) {
	var memberships []models.Membership
	query := `
        SELECT m.id, m.type, m.price, m.start_date, m.end_date, m.is_active, 
               m.user_id, m.created_at 
        FROM memberships m
        WHERE m.user_id = $1 AND m.is_active = true
        ORDER BY m.end_date DESC
    `
	err := r.db.Select(&memberships, query, userID)
	return memberships, err
}

func (r *membershipRepository) Create(membership *models.Membership) error {
	query := `
        INSERT INTO memberships (user_id, type, price, start_date, end_date, is_active)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id
    `
	return r.db.QueryRow(
		query,
		membership.UserID,
		membership.Type,
		membership.Price,
		membership.StartDate,
		membership.EndDate,
		membership.IsActive,
	).Scan(&membership.ID)
}
