package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/yasinbozat/ecommerce-platform/services/search-service/internal/cache"
	"github.com/yasinbozat/ecommerce-platform/services/search-service/internal/config"
)

type SearchResult struct {
	Total int64       `json:"total"`
	Hits  []SearchHit `json:"hits"`
}

type SearchHit struct {
	ID     string                 `json:"id"`
	Index  string                 `json:"index"`
	Source map[string]interface{} `json:"source"`
}

type esSearchClient interface {
	Search(...func(*esapi.SearchRequest)) (*esapi.Response, error)
}

type elasticsearchClient struct {
	client *elasticsearch.Client
}

func NewElasticsearchClient(client *elasticsearch.Client) esSearchClient {
	return elasticsearchClient{client: client}
}

func (c elasticsearchClient) Search(opts ...func(*esapi.SearchRequest)) (*esapi.Response, error) {
	return c.client.Search(opts...)
}

type SearchService interface {
	Search(ctx context.Context, query string) (*SearchResult, error)
}

type searchService struct {
	es    esSearchClient
	cache cache.SearchCache
	cfg   *config.Config
}

func NewSearchService(es esSearchClient, cacheStore cache.SearchCache, cfg *config.Config) SearchService {
	if cacheStore == nil {
		cacheStore = cache.NewNoopSearchCache()
	}
	return &searchService{es: es, cache: cacheStore, cfg: cfg}
}

type esSearchResponse struct {
	Hits struct {
		Total struct {
			Value int64 `json:"value"`
		} `json:"total"`
		Hits []struct {
			ID     string          `json:"_id"`
			Index  string          `json:"_index"`
			Source json.RawMessage `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

func (s *searchService) Search(ctx context.Context, query string) (*SearchResult, error) {
	if query == "" {
		return nil, fmt.Errorf("query is required")
	}

	key := cache.CacheKey(s.cfg.Redis.KeyPrefix, query)
	if cached, err := s.cache.Get(ctx, key); err == nil && len(cached) != 0 {
		var result SearchResult
		if err := json.Unmarshal(cached, &result); err == nil {
			return &result, nil
		}
	}

	payload := map[string]interface{}{
		"query": map[string]interface{}{
			"query_string": map[string]interface{}{
				"query":  query,
				"fields": []string{"*"},
			},
		},
	}

	body := new(bytes.Buffer)
	if err := json.NewEncoder(body).Encode(payload); err != nil {
		return nil, err
	}

	res, err := s.es.Search(func(r *esapi.SearchRequest) {
		r.Body = body
		r.Index = []string{s.cfg.Elastic.Index}
		if s.cfg.Elastic.Timeout > 0 {
			r.Timeout = s.cfg.Elastic.Timeout
		}
	})
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("elasticsearch search failed: %s", res.String())
	}

	var raw esSearchResponse
	if err := json.NewDecoder(res.Body).Decode(&raw); err != nil {
		return nil, err
	}

	result := &SearchResult{Total: raw.Hits.Total.Value}
	for _, hit := range raw.Hits.Hits {
		source := map[string]interface{}{}
		if len(hit.Source) > 0 {
			_ = json.Unmarshal(hit.Source, &source)
		}
		result.Hits = append(result.Hits, SearchHit{ID: hit.ID, Index: hit.Index, Source: source})
	}

	if ttl := s.cfg.Redis.CacheTTL; ttl > 0 {
		if payload, err := json.Marshal(result); err == nil {
			_ = s.cache.Set(ctx, key, payload, ttl)
		}
	}

	return result, nil
}
