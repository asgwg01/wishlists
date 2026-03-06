package redis

import (
	"authService/internal/cache"
	"authService/internal/config"
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	log         *slog.Logger
	redisClient *redis.Client
}

func NewCache(log *slog.Logger, cfg *config.RedisConfig) (*Cache, error) {
	l := log.With(
		slog.String("cashe", "Redis"),
	)

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       0,
	})

	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	// if err := client.Ping(ctx).Err(); err != nil {
	//     log.Fatalf("Failed to connect to Redis: %v", err)
	// }

	//return client

	return &Cache{log: l, redisClient: client}, nil
}

func (c *Cache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	const logPrefix = "cache.redis.Set"
	log := c.log.With(
		slog.String("where", logPrefix),
	)

	log.Debug("Set key", slog.String("key", key))

	if err := c.redisClient.Set(ctx, key, value, expiration).Err(); err != nil {
		return fmt.Errorf("failed to set key %s: %w", key, err)
	}
	return nil
}
func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	const logPrefix = "cache.redis.Get"
	log := c.log.With(
		slog.String("where", logPrefix),
	)

	log.Debug("Get key", slog.String("key", key))

	val, err := c.redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", cache.ErrorTokenNotSet
	} else if err != nil {
		return "", fmt.Errorf("failed to get key %s: %w", key, err)
	}
	return val, nil
}

func (c *Cache) Delete(ctx context.Context, key string) error {
	const logPrefix = "cache.redis.Delete"
	log := c.log.With(
		slog.String("where", logPrefix),
	)

	log.Debug("Delete key", slog.String("key", key))
	if err := c.redisClient.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete key %s: %w", key, err)
	}
	return nil
}
func (c *Cache) Exists(ctx context.Context, key string) (bool, error) {
	const logPrefix = "cache.redis.Exists"
	log := c.log.With(
		slog.String("where", logPrefix),
	)

	log.Debug("Check exists key")

	n, err := c.redisClient.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check key %s: %w", key, err)
	}
	return n > 0, nil
}
func (c *Cache) Ping(ctx context.Context) error {
	const logPrefix = "cache.redis.Ping"
	log := c.log.With(
		slog.String("where", logPrefix),
	)

	log.Debug("Ping cache")

	if err := c.redisClient.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis ping failed: %w", err)
	}
	return nil
}
func (c *Cache) Close() error {
	const logPrefix = "cache.redis.Close"
	log := c.log.With(
		slog.String("where", logPrefix),
	)

	log.Debug("Close cache")

	if err := c.redisClient.Close(); err != nil {
		return fmt.Errorf("failed to close redis connection: %w", err)
	}
	return nil
}
