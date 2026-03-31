package service

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/yasinbozat/ecommerce-platform/services/search-service/internal/config"
)

type stubEsClient struct {
	response []byte
	err      error
	calls    int
}

func (c *stubEsClient) Search(opts ...func(*esapi.SearchRequest)) (*esapi.Response, error) {
	c.calls++
	if c.err != nil {
		return nil, c.err
	}
	var req esapi.SearchRequest
	for _, opt := range opts {
		opt(&req)
	}
	return &esapi.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewReader(c.response)),
		Header:     http.Header{},
	}, nil
}

type stubSearchCache struct {
	getData []byte
	getErr  error
	setKey  string
	setTTL  time.Duration
	calls   int
}

func (c *stubSearchCache) Get(ctx context.Context, key string) ([]byte, error) {
	return c.getData, c.getErr
}

func (c *stubSearchCache) Set(ctx context.Context, key string, payload []byte, ttl time.Duration) error {
	c.calls++
	c.setKey = key
	c.setTTL = ttl
	return nil
}

func TestSearchServiceUsesCache(t *testing.T) {
	cached := &SearchResult{Total: 1, Hits: []SearchHit{{ID: "1", Index: "products"}}}
	data, _ := json.Marshal(cached)

	cacheStore := &stubSearchCache{getData: data}
	svc := NewSearchService(&stubEsClient{}, cacheStore, newTestConfig())

	result, err := svc.Search(context.Background(), "query")
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	if result.Total != cached.Total {
		t.Fatalf("unexpected total %d", result.Total)
	}
	if cacheStore.calls != 0 {
		t.Fatalf("expected cache Set not called, got %d", cacheStore.calls)
	}
}

func TestSearchServiceCachesElasticsearchResult(t *testing.T) {
	response := []byte(`{"hits":{"total":{"value":1},"hits":[{"_id":"1","_index":"products","_source":{"name":"test"}}]}}`)
	cacheStore := &stubSearchCache{}
	svc := NewSearchService(&stubEsClient{response: response}, cacheStore, newTestConfig())

	_, err := svc.Search(context.Background(), "query")
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	if cacheStore.calls != 1 {
		t.Fatalf("expected cache Set called once, got %d", cacheStore.calls)
	}
	if cacheStore.setTTL != newTestConfig().Redis.CacheTTL {
		t.Fatalf("unexpected ttl %v", cacheStore.setTTL)
	}
}

func newTestConfig() *config.Config {
	return &config.Config{
		Elastic: config.ElasticConfig{Index: "products"},
		Redis:   config.RedisConfig{KeyPrefix: "search-service", CacheTTL: 2 * time.Minute},
	}
}
