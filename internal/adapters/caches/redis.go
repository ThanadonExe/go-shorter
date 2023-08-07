package caches

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

var ctx = context.Background()

func NewRedisCache(host string, port string, password string, db int) (*RedisCache, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       db,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &RedisCache{client: rdb}, nil
}

func (r *RedisCache) Get(key string, value interface{}) error {
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(data), value); err != nil {
		return err
	}

	return nil
}

func (r *RedisCache) Set(key string, value interface{}, duration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = r.client.Set(ctx, key, data, duration).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *RedisCache) Delete(key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *RedisCache) Close() error {
	return r.client.Close()
}
