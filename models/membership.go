package models

import (
	"time"
)

type Membership struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	MemberID         uint      `gorm:"not null" json:"member_id"`
	MembershipTypeID uint      `gorm:"not null" json:"membership_type_id"`
	StartDate        time.Time `gorm:"not null" json:"start_date"`
	EndDate          time.Time `gorm:"not null" json:"end_date"`
	Status           string    `gorm:"not null;default:'active'" json:"status"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`

	Member         Member         `gorm:"foreignKey:MemberID" json:"member,omitempty"`
	MembershipType MembershipType `gorm:"foreignKey:MembershipTypeID" json:"membership_type,omitempty"`
}
