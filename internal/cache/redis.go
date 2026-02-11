package cache

import (
	"context"

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

// Проверка подключения
func (r *RedisClient) Ping() error {
	return r.client.Ping(r.ctx).Err()
}

// Закрытие соединения
func (r *RedisClient) Close() error {
	return r.client.Close()
}
