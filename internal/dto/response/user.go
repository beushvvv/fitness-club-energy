package response

import (
	"fitness-club-energy/internal/model"
	"time"
)

// UserResponse представляет пользователя для ответа API
type UserResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
}

// FromUser конвертирует модель User в UserResponse
func FromUser(user *model.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt,
	}
}

// FromUsers конвертирует срез моделей User в срез UserResponse
func FromUsers(users []model.User) []UserResponse {
	result := make([]UserResponse, len(users))
	for i, user := range users {
		result[i] = FromUser(&user)
	}
	return result
}
