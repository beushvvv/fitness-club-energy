package cache

import "fmt"

// Ключи для кеширования
const (
	// Префиксы для разных типов данных
	UsersAllPrefix       = "users:all"
	UserByIDPrefix       = "user:%d"
	MembershipsAllPrefix = "memberships:all"
	MembershipByIDPrefix = "membership:%d"
	WorkoutsFilterPrefix = "workouts:filter:%s"
)

// UserKey возвращает ключ для кеша пользователя по ID
func UserKey(id int) string {
	return fmt.Sprintf(UserByIDPrefix, id)
}

// MembershipKey возвращает ключ для кеша абонемента по ID
func MembershipKey(id int) string {
	return fmt.Sprintf(MembershipByIDPrefix, id)
}

// WorkoutsFilterKey возвращает ключ для кеша тренировок с фильтром
func WorkoutsFilterKey(filterHash string) string {
	return fmt.Sprintf(WorkoutsFilterPrefix, filterHash)
}

// KeysToInvalidate возвращает ключи, которые нужно очистить при изменении данных
func KeysToInvalidate(entityType string) []string {
	switch entityType {
	case "user":
		return []string{UsersAllPrefix}
	case "membership":
		return []string{MembershipsAllPrefix}
	case "workout":
		return []string{} // Для тренировок нужна более сложная логика
	default:
		return []string{}
	}
}
