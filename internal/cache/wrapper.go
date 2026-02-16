package cache

import (
	"encoding/json"
	"time"
)

// CacheWrapper предоставляет удобные методы для работы с кешем
type CacheWrapper struct {
	client *RedisClient
}

func NewCacheWrapper(client *RedisClient) *CacheWrapper {
	return &CacheWrapper{client: client}
}

// GetOrSet пытается получить данные из кеша, если нет - вызывает функцию и сохраняет результат
func (cw *CacheWrapper) GetOrSet(key string, ttl time.Duration, dest interface{}, fn func() (interface{}, error)) error {
	// Пробуем получить из кеша
	err := cw.client.Get(key, dest)
	if err == nil {
		return nil // Данные найдены в кеше
	}

	// Вызываем функцию для получения данных
	data, err := fn()
	if err != nil {
		return err
	}

	// Сохраняем в кеш
	if err := cw.client.Set(key, data, ttl); err != nil {
		// Логируем ошибку, но не возвращаем - данные уже получены
		// log.Printf("Failed to set cache: %v", err)
	}

	// Копируем данные в dest
	// Здесь нужно использовать рефлексию или json.Marshal/Unmarshal
	// Для простоты предлагаю использовать json
	bytes, _ := json.Marshal(data)
	json.Unmarshal(bytes, dest)

	return nil
}
