package cache

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func TestRedisSearchCacheRoundTrip(t *testing.T) {
	server, err := miniredis.Run()
	if err != nil {
		t.Fatalf("miniredis.Run returned error: %v", err)
	}
	defer server.Close()

	client := redis.NewClient(&redis.Options{Addr: server.Addr()})
	defer client.Close()

	cacheStore := NewRedisSearchCache(client, "search-service").(*RedisSearchCache)
	payload := []byte("payload")
	key := "test:key"

	if err := cacheStore.Set(context.Background(), key, payload, time.Minute); err != nil {
		t.Fatalf("Set returned error: %v", err)
	}

	val, err := cacheStore.Get(context.Background(), key)
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if string(val) != "payload" {
		t.Fatalf("unexpected value %s", val)
	}
}

func TestCacheKeyHashesQuery(t *testing.T) {
	key := CacheKey("search-service", "my query")
	if !strings.HasPrefix(key, "search-service:search:") {
		t.Fatalf("unexpected key %s", key)
	}
	if strings.Contains(key, "my query") {
		t.Fatalf("key leaked raw query %s", key)
	}
}
