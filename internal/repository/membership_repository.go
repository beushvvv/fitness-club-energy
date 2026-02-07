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

func (r *MembershipRepository) GetAll() ([]model.Membership, error) {
	var memberships []model.Membership
	query := `
        SELECT id, user_id, type, price, start_date, end_date, is_active, created_at 
        FROM memberships 
        ORDER BY created_at DESC
    `
	err := r.db.Select(&memberships, query)
	return memberships, err
}

func (r *MembershipRepository) GetByUserID(userID int) ([]model.Membership, error) {
	var memberships []model.Membership
	query := `
        SELECT id, user_id, type, price, start_date, end_date, is_active, created_at 
        FROM memberships 
        WHERE user_id = $1 AND is_active = true
        ORDER BY end_date DESC
    `
	err := r.db.Select(&memberships, query, userID)
	return memberships, err
}

func (r *MembershipRepository) Create(membership *model.Membership) error {
	query := `
        INSERT INTO memberships (user_id, type, price, start_date, end_date, is_active)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, created_at
    `
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

func (r *MembershipRepository) Update(membership *model.Membership) error {
	query := `
        UPDATE memberships 
        SET type = $1, price = $2, start_date = $3, end_date = $4, is_active = $5 
        WHERE id = $6
    `
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
