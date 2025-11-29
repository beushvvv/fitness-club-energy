package models

import (
	"time"
)

type Member struct {
	ID               uint              `gorm:"primaryKey" json:"id"`
	FirstName        string            `gorm:"not null" json:"first_name"`
	LastName         string            `gorm:"not null" json:"last_name"`
	Email            string            `gorm:"unique;not null" json:"email"`
	Phone            string            `json:"phone"`
	DateOfBirth      time.Time         `json:"date_of_birth"`
	RegistrationDate time.Time         `gorm:"default:CURRENT_TIMESTAMP" json:"registration_date"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
	Memberships      []Membership      `gorm:"foreignKey:MemberID" json:"memberships,omitempty"`
	TrainingSessions []TrainingSession `gorm:"foreignKey:MemberID" json:"training_sessions,omitempty"`
}
