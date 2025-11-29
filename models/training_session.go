package models

import (
	"time"
)

type TrainingSession struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	MemberID    uint      `gorm:"not null" json:"member_id"`
	TrainerID   uint      `gorm:"not null" json:"trainer_id"`
	SessionDate time.Time `gorm:"not null" json:"session_date"`
	Duration    int       `gorm:"not null" json:"duration"` // в минутах
	SessionType string    `gorm:"not null" json:"session_type"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Member  Member  `gorm:"foreignKey:MemberID" json:"member,omitempty"`
	Trainer Trainer `gorm:"foreignKey:TrainerID" json:"trainer,omitempty"`
}
