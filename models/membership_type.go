package models

import (
	"time"
)

type MembershipType struct {
	ID          uint         `gorm:"primaryKey" json:"id"`
	Name        string       `gorm:"not null" json:"name"`
	Description string       `json:"description"`
	Price       float64      `gorm:"not null" json:"price"`
	Duration    int          `gorm:"not null" json:"duration"` // в днях
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	Memberships []Membership `gorm:"foreignKey:MembershipTypeID" json:"memberships,omitempty"`
}
