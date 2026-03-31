package cache

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/yasinbozat/ecommerce-platform/services/user-service/internal/domain"
)

func TestRedisAuthCacheRoundTrip(t *testing.T) {
	server, err := miniredis.Run()
	if err != nil {
		t.Fatalf("miniredis.Run returned error: %v", err)
	}
	defer server.Close()

	client := redis.NewClient(&redis.Options{Addr: server.Addr()})
	defer client.Close()

	cache := NewRedisAuthCache(client, "user-service").(*RedisAuthCache)
	response := &domain.ValidateTokenResponse{
		UserId:     uuid.New(),
		Email:      "user@example.com",
		Role:       domain.RoleAdmin,
		KeycloakId: "kc-1",
	}

	if err := cache.SetValidatedToken(context.Background(), "secret-token", response, time.Minute); err != nil {
		t.Fatalf("SetValidatedToken returned error: %v", err)
	}

	if !server.Exists(cache.keyForToken("secret-token")) {
		t.Fatalf("expected key %q to exist in redis", cache.keyForToken("secret-token"))
	}

	cachedResponse, err := cache.GetValidatedToken(context.Background(), "secret-token")
	if err != nil {
		t.Fatalf("GetValidatedToken returned error: %v", err)
	}

	if *cachedResponse != *response {
		t.Fatalf("GetValidatedToken returned %+v, want %+v", cachedResponse, response)
	}
}

func TestRedisAuthCacheKeyUsesTokenHash(t *testing.T) {
	cache := &RedisAuthCache{keyPrefix: "user-service"}

	key := cache.keyForToken("plain-text-token")

	if !strings.HasPrefix(key, "user-service:auth:token:") {
		t.Fatalf("unexpected key prefix %q", key)
	}
	if strings.Contains(key, "plain-text-token") {
		t.Fatalf("expected key %q to avoid storing the raw token", key)
	}
}
