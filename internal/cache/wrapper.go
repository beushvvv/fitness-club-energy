package cache

import (
	"encoding/json"
	"time"
)

// CacheWrapper обёртка над RedisClient для удобного использования в сервисах
type CacheWrapper struct {
	client *RedisClient
}

// GetOrSet получает данные из кэша или устанавливает их, если нет
func (cw *CacheWrapper) GetOrSet(key string, ttl time.Duration, dest interface{}, fallback func() (interface{}, error)) error {
	// Пробуем получить из кэша
	err := cw.Get(key, dest)
	if err == nil {
		return nil // Данные из кэша
	}

	// Нет в кэше - вызываем fallback функцию
	data, err := fallback()
	if err != nil {
		return err
	}

	// Сохраняем в кэш
	if err := cw.Set(key, data, ttl); err != nil {
		return err
	}

	// Заполняем dest полученными данными
	jsonData, _ := json.Marshal(data)
	return json.Unmarshal(jsonData, dest)
}

// NewCacheWrapper создаёт новый экземпляр CacheWrapper
func NewCacheWrapper(client *RedisClient) *CacheWrapper {
	return &CacheWrapper{client: client}
}

// Get получает данные из кэша по ключу
func (cw *CacheWrapper) Get(key string, dest interface{}) error {
	return cw.client.Get(key, dest)
}

// Set сохраняет данные в кэш с TTL
func (cw *CacheWrapper) Set(key string, value interface{}, ttl time.Duration) error {
	return cw.client.Set(key, value, ttl)
}

// Delete удаляет ключ из кэша
func (cw *CacheWrapper) Delete(key string) error {
	return cw.client.Delete(key)
}

// ClearUsersCache очищает кэш пользователей
func (cw *CacheWrapper) ClearUsersCache() error {
	return cw.Delete("users:all")
}

// ClearMembershipsCache очищает кэш абонементов
func (cw *CacheWrapper) ClearMembershipsCache() error {
	return cw.Delete("memberships:all")
}

// ClearWorkoutsCache очищает кэш тренировок
func (cw *CacheWrapper) ClearWorkoutsCache() error {
	return cw.Delete("workouts:all")
}

// ClearAllCaches очищает все кэши
func (cw *CacheWrapper) ClearAllCaches() error {
	if err := cw.ClearUsersCache(); err != nil {
		return err
	}
	if err := cw.ClearMembershipsCache(); err != nil {
		return err
	}
	return cw.ClearWorkoutsCache()
}
