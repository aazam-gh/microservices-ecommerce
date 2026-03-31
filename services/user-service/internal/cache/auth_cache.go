package cache

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yasinbozat/ecommerce-platform/services/user-service/internal/domain"
)

type AuthCache interface {
	GetValidatedToken(ctx context.Context, token string) (*domain.ValidateTokenResponse, error)
	SetValidatedToken(ctx context.Context, token string, response *domain.ValidateTokenResponse, ttl time.Duration) error
}

type NoopAuthCache struct{}

func NewNoopAuthCache() AuthCache {
	return NoopAuthCache{}
}

func (NoopAuthCache) GetValidatedToken(context.Context, string) (*domain.ValidateTokenResponse, error) {
	return nil, nil
}

func (NoopAuthCache) SetValidatedToken(context.Context, string, *domain.ValidateTokenResponse, time.Duration) error {
	return nil
}

type RedisAuthCache struct {
	client    *redis.Client
	keyPrefix string
}

func NewRedisAuthCache(client *redis.Client, keyPrefix string) AuthCache {
	if client == nil {
		return NewNoopAuthCache()
	}

	return &RedisAuthCache{
		client:    client,
		keyPrefix: keyPrefix,
	}
}

func (c *RedisAuthCache) GetValidatedToken(ctx context.Context, token string) (*domain.ValidateTokenResponse, error) {
	value, err := c.client.Get(ctx, c.keyForToken(token)).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var response domain.ValidateTokenResponse
	if err := json.Unmarshal([]byte(value), &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *RedisAuthCache) SetValidatedToken(ctx context.Context, token string, response *domain.ValidateTokenResponse, ttl time.Duration) error {
	payload, err := json.Marshal(response)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, c.keyForToken(token), payload, ttl).Err()
}

func (c *RedisAuthCache) keyForToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return c.keyPrefix + ":auth:token:" + hex.EncodeToString(sum[:])
}
