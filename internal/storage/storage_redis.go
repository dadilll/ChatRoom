package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type UserStorageRedis interface {
	SetVerificationCode(ctx context.Context, email, code string, expiration time.Duration) error
	GetVerificationCode(ctx context.Context, email string) (string, error)
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, key string) error
}

type RedisStorage struct {
	client *redis.Client
}

func NewRedisStorage(client *redis.Client) *RedisStorage {
	return &RedisStorage{
		client: client,
	}
}

func (r *RedisStorage) SetVerificationCode(ctx context.Context, email, code string, expiration time.Duration) error {
	key := fmt.Sprintf("verify_code:%s", email)
	return r.client.Set(ctx, key, code, expiration).Err()
}

func (r *RedisStorage) GetVerificationCode(ctx context.Context, email string) (string, error) {
	key := fmt.Sprintf("verify_code:%s", email)
	return r.client.Get(ctx, key).Result()
}

func (r *RedisStorage) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *RedisStorage) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisStorage) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
