package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisClient клиент для работы с Redis
type RedisClient struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedisClient создаёт новый клиент Redis
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

// Ping проверяет подключение к Redis
func (r *RedisClient) Ping() error {
	return r.client.Ping(r.ctx).Err()
}

// Close закрывает соединение с Redis
func (r *RedisClient) Close() error {
	return r.client.Close()
}

// Set сохраняет данные в Redis с TTL
func (r *RedisClient) Set(key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(r.ctx, key, data, ttl).Err()
}

// Get получает данные из Redis
func (r *RedisClient) Get(key string, dest interface{}) error {
	data, err := r.client.Get(r.ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

// Delete удаляет ключ из Redis
func (r *RedisClient) Delete(key string) error {
	return r.client.Del(r.ctx, key).Err()
}

// Exists проверяет существование ключа
func (r *RedisClient) Exists(key string) (bool, error) {
	n, err := r.client.Exists(r.ctx, key).Result()
	return n > 0, err
}
