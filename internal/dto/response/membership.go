package response

import (
	"fitness-club-energy/internal/model"
	"time"
)

// MembershipResponse представляет абонемент для ответа API
type MembershipResponse struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Type      string    `json:"type"`
	Price     float64   `json:"price"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

// FromMembership конвертирует модель Membership в MembershipResponse
func FromMembership(membership *model.Membership) MembershipResponse {
	return MembershipResponse{
		ID:        membership.ID,
		UserID:    membership.UserID,
		Type:      membership.Type,
		Price:     membership.Price,
		StartDate: membership.StartDate,
		EndDate:   membership.EndDate,
		IsActive:  membership.IsActive,
		CreatedAt: membership.CreatedAt,
	}
}

// FromMemberships конвертирует срез моделей Membership в срез MembershipResponse
func FromMemberships(memberships []model.Membership) []MembershipResponse {
	result := make([]MembershipResponse, len(memberships))
	for i, membership := range memberships {
		result[i] = FromMembership(&membership)
	}
	return result
}
