package models

import (
	"time"
)

type Trainer struct {
	ID               uint              `gorm:"primaryKey" json:"id"`
	FirstName        string            `gorm:"not null" json:"first_name"`
	LastName         string            `gorm:"not null" json:"last_name"`
	Specialization   string            `json:"specialization"`
	Phone            string            `json:"phone"`
	Email            string            `gorm:"unique" json:"email"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
	TrainingSessions []TrainingSession `gorm:"foreignKey:TrainerID" json:"training_sessions,omitempty"`
}
