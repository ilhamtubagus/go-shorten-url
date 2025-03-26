package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ilhamtubagus/go-shorten-url/constants"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

type Cache[T any] interface {
	Get(ctx context.Context, key string) (T, error)
	Put(ctx context.Context, key string, val T, ttl uint64) error
	IsExist(ctx context.Context, key string) (bool, error)
	Delete(ctx context.Context, key string) error
	Flush(ctx context.Context) error
}

type RedisCache[T any] struct {
	client *redis.Client
}

func NewRedisCache[T any](rc *redis.Client) *RedisCache[T] {
	return &RedisCache[T]{
		client: rc,
	}
}

func (rc *RedisCache[T]) Get(ctx context.Context, key string) (T, error) {
	var result T
	val, err := rc.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return result, constants.ErrorCacheNotFound
	} else if err != nil {
		return result, err
	}

	// Unmarshal the JSON string into the struct
	err = json.Unmarshal([]byte(val), &result)
	if err != nil {
		return result, fmt.Errorf("failed to unmarshal data: %w", err)
	}

	log.Println("cache", result)

	return result, nil
}

func (rc *RedisCache[T]) Put(ctx context.Context, key string, val T, ttl uint64) error {
	// Marshal the struct to JSON
	jsonData, err := json.Marshal(val)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	expiration := time.Duration(ttl) * time.Second
	err = rc.client.Set(ctx, key, jsonData, expiration).Err()

	return err
}

func (rc *RedisCache[T]) IsExist(ctx context.Context, key string) (bool, error) {
	r, err := rc.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return r > 0, nil
}

func (rc *RedisCache[T]) Delete(ctx context.Context, key string) error {
	return rc.client.Del(ctx, key).Err()
}

func (rc *RedisCache[T]) Flush(ctx context.Context) error {
	return rc.client.FlushAll(ctx).Err()
}
