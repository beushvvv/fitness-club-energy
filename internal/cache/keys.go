package cache

import "fmt"

const (
	UsersAllPrefix       = "users:all"
	MembershipsAllPrefix = "memberships:all"
	WorkoutsAllPrefix    = "workouts:all"
)

// UserKey возвращает ключ для кеширования пользователя по ID
func UserKey(id int) string {
	return fmt.Sprintf("user:%d", id)
}

// MembershipKey возвращает ключ для кеширования абонемента по ID
func MembershipKey(id int) string {
	return fmt.Sprintf("membership:%d", id)
}

// WorkoutKey возвращает ключ для кеширования тренировки по ID
func WorkoutKey(id int) string {
	return fmt.Sprintf("workout:%d", id)
}

// KeysToInvalidate возвращает ключи для инвалидации при изменении сущности
func KeysToInvalidate(entity string) []string {
	switch entity {
	case "user":
		return []string{UsersAllPrefix}
	case "membership":
		return []string{MembershipsAllPrefix}
	case "workout":
		return []string{WorkoutsAllPrefix}
	default:
		return []string{}
	}
}
