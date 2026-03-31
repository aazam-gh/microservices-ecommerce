package cache

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/redis/go-redis/v9"
)

type SearchCache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, payload []byte, ttl time.Duration) error
}

type NoopSearchCache struct{}

func NewNoopSearchCache() SearchCache {
	return NoopSearchCache{}
}

func (NoopSearchCache) Get(context.Context, string) ([]byte, error) {
	return nil, nil
}

func (NoopSearchCache) Set(context.Context, string, []byte, time.Duration) error {
	return nil
}

type RedisSearchCache struct {
	client    *redis.Client
	keyPrefix string
}

func NewRedisSearchCache(client *redis.Client, keyPrefix string) SearchCache {
	if client == nil {
		return NewNoopSearchCache()
	}
	return &RedisSearchCache{client: client, keyPrefix: keyPrefix}
}

func (c *RedisSearchCache) Get(ctx context.Context, key string) ([]byte, error) {
	val, err := c.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return val, nil
}

func (c *RedisSearchCache) Set(ctx context.Context, key string, payload []byte, ttl time.Duration) error {
	return c.client.Set(ctx, key, payload, ttl).Err()
}

func CacheKey(prefix, query string) string {
	sum := sha256.Sum256([]byte(query))
	return prefix + ":search:" + hex.EncodeToString(sum[:])
}
