package model

import "time"

type Membership struct {
	ID        int       `db:"id"`
	UserID    int       `db:"user_id"`
	Type      string    `db:"type"`
	Price     float64   `db:"price"`
	StartDate time.Time `db:"start_date"`
	EndDate   time.Time `db:"end_date"`
	IsActive  bool      `db:"is_active"`
	CreatedAt time.Time `db:"created_at"`
}
