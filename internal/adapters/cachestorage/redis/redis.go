package redis

import (
	"context"
	"frog-go/internal/core/ports/outbound/cachestorage"
	"frog-go/internal/utils/logger"
	"fmt"
	"net/url"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
}

func NewRedis(password, host, port string) (cachestorage.CacheStorage, error) {
	lg := logger.NewLogger("Redis")

	uri := fmt.Sprintf(
		"redis://:%s@%s:%s/0",
		url.QueryEscape(password),
		host,
		port,
	)

	opts, err := redis.ParseURL(uri)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opts)

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	lg.Start(
		"Host: %s:%s | DB: %d",
		host,
		port,
		opts.DB,
	)

	return &Redis{client: client}, nil
}

func (r *Redis) Close() error {
	return r.client.Close()
}
func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key '%s' not found", key)
	}
	return val, err
}

func (r *Redis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) (string, error) {
	return r.client.Set(ctx, key, value, expiration).Result()
}

func (r *Redis) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	return r.client.SetNX(ctx, key, value, expiration).Result()
}

func (r *Redis) Incr(ctx context.Context, key string) (int64, error) {
	return r.client.Incr(ctx, key).Result()
}

func (r *Redis) Expire(ctx context.Context, key string, ttl time.Duration) error {
	return r.client.Expire(ctx, key, ttl).Err()
}

func (r *Redis) WaitForCacheValue(
	ctx context.Context,
	key string,
	interval time.Duration,
	timeout time.Duration,
	predicate func(string) (bool, error),
) (string, error) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	timeoutCh := time.After(timeout)

	for {
		select {
		case <-timeoutCh:
			return "", context.DeadlineExceeded
		case <-ticker.C:
			val, err := r.client.Get(ctx, key).Result()
			if err != nil {
				if err == redis.Nil {
					continue // Valor ainda não está no cache
				}
				return "", err // Outro erro inesperado
			}

			ok, err := predicate(val)
			if err != nil {
				return "", err
			}
			if ok {
				return val, nil
			}
		}
	}
}
