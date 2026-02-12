package cache

import (
	"context"

	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisClient(addr, password string, db int) *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &RedisClient{
		client: client,
		ctx:    context.Background(),
	}
}

// Set сохраняет данные в Redis с TTL
func (r *RedisClient) Set(key string, value interface{}, ttl time.Duration) error {
	// Сериализуем в JSON
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	// Сохраняем в Redis
	return r.client.Set(r.ctx, key, data, ttl).Err()
}

// Get получает данные из Redis
func (r *RedisClient) Get(key string, dest interface{}) error {
	// Получаем данные
	data, err := r.client.Get(r.ctx, key).Bytes()
	if err != nil {
		return err
	}

	// Десериализуем
	return json.Unmarshal(data, dest)
}

// Delete удаляет ключ
func (r *RedisClient) Delete(key string) error {
	return r.client.Del(r.ctx, key).Err()
}

// Exists проверяет существование ключа
func (r *RedisClient) Exists(key string) (bool, error) {
	n, err := r.client.Exists(r.ctx, key).Result()
	return n > 0, err
}

// Проверка подключения
func (r *RedisClient) Ping() error {
	return r.client.Ping(r.ctx).Err()
}

// Закрытие соединения
func (r *RedisClient) Close() error {
	return r.client.Close()
}
