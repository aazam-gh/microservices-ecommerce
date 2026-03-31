package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/yasinbozat/ecommerce-platform/services/user-service/internal/cache"
	"github.com/yasinbozat/ecommerce-platform/services/user-service/internal/config"
	"github.com/yasinbozat/ecommerce-platform/services/user-service/internal/domain"
)

type stubUserRepo struct {
	userByKeycloakID map[string]*domain.User
	findErr          error
	createErr        error
	createCalls      int
}

func (r *stubUserRepo) FindByID(context.Context, uuid.UUID) (*domain.User, error) {
	return nil, nil
}

func (r *stubUserRepo) FindByKeycloakID(_ context.Context, keycloakID string) (*domain.User, error) {
	if r.findErr != nil {
		return nil, r.findErr
	}
	if r.userByKeycloakID == nil {
		return nil, nil
	}
	return r.userByKeycloakID[keycloakID], nil
}

func (r *stubUserRepo) FindByEmail(context.Context, string) (*domain.User, error) {
	return nil, nil
}

func (r *stubUserRepo) Create(_ context.Context, user *domain.User) error {
	if r.createErr != nil {
		return r.createErr
	}

	r.createCalls++
	if user.Id == uuid.Nil {
		user.Id = uuid.New()
	}

	if r.userByKeycloakID == nil {
		r.userByKeycloakID = make(map[string]*domain.User)
	}

	cloned := *user
	r.userByKeycloakID[user.KeycloakId] = &cloned
	return nil
}

func (r *stubUserRepo) Update(context.Context, *domain.User) error {
	return nil
}

func (r *stubUserRepo) Delete(context.Context, uuid.UUID) error {
	return nil
}

type stubAuthCache struct {
	getResponse *domain.ValidateTokenResponse
	getErr      error
	setErr      error
	setCalls    int
	setToken    string
	setTTL      time.Duration
	setResponse *domain.ValidateTokenResponse
}

func (c *stubAuthCache) GetValidatedToken(context.Context, string) (*domain.ValidateTokenResponse, error) {
	if c.getErr != nil {
		return nil, c.getErr
	}
	return c.getResponse, nil
}

func (c *stubAuthCache) SetValidatedToken(_ context.Context, token string, response *domain.ValidateTokenResponse, ttl time.Duration) error {
	c.setCalls++
	c.setToken = token
	c.setTTL = ttl
	if response != nil {
		cloned := *response
		c.setResponse = &cloned
	}
	return c.setErr
}

func TestValidateTokenReturnsCachedResponse(t *testing.T) {
	cachedResponse := &domain.ValidateTokenResponse{
		UserId:     uuid.New(),
		Email:      "cached@example.com",
		Role:       domain.RoleAdmin,
		KeycloakId: "kc-cached",
	}

	service := NewAuthService(&stubUserRepo{}, newTestConfig("http://example.com", 5*time.Minute), &stubAuthCache{
		getResponse: cachedResponse,
	})

	response, err := service.ValidateToken(context.Background(), "token-1")
	if err != nil {
		t.Fatalf("ValidateToken returned error: %v", err)
	}

	if *response != *cachedResponse {
		t.Fatalf("ValidateToken returned %+v, want %+v", response, cachedResponse)
	}
}

func TestValidateTokenCachesSuccessfulIntrospection(t *testing.T) {
	server := newIntrospectionServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method %s", r.Method)
		}
		if err := r.ParseForm(); err != nil {
			t.Fatalf("ParseForm returned error: %v", err)
		}
		if got := r.FormValue("token"); got != "token-1" {
			t.Fatalf("unexpected token %q", got)
		}

		fmt.Fprintf(w, `{"active":true,"sub":"kc-1","email":"user@example.com","name":"Jane Doe","exp":%d}`, time.Now().Add(10*time.Minute).Unix())
	})
	defer server.Close()

	cacheStore := &stubAuthCache{}
	repo := &stubUserRepo{}
	service := NewAuthService(repo, newTestConfig(server.URL, 5*time.Minute), cacheStore)

	response, err := service.ValidateToken(context.Background(), "token-1")
	if err != nil {
		t.Fatalf("ValidateToken returned error: %v", err)
	}

	if response.Email != "user@example.com" {
		t.Fatalf("unexpected email %q", response.Email)
	}
	if response.KeycloakId != "kc-1" {
		t.Fatalf("unexpected keycloak id %q", response.KeycloakId)
	}
	if response.Role != domain.RoleCustomer {
		t.Fatalf("unexpected role %q", response.Role)
	}
	if response.UserId == uuid.Nil {
		t.Fatal("expected synced user id to be set")
	}
	if repo.createCalls != 1 {
		t.Fatalf("Create was called %d times, want 1", repo.createCalls)
	}
	if cacheStore.setCalls != 1 {
		t.Fatalf("SetValidatedToken was called %d times, want 1", cacheStore.setCalls)
	}
	if cacheStore.setToken != "token-1" {
		t.Fatalf("unexpected cached token %q", cacheStore.setToken)
	}
	if cacheStore.setTTL != 5*time.Minute {
		t.Fatalf("unexpected cache ttl %s", cacheStore.setTTL)
	}
}

func TestValidateTokenFallsBackWhenCacheReadFails(t *testing.T) {
	server := newIntrospectionServer(t, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"active":true,"sub":"kc-2","email":"user2@example.com","name":"Jane Doe","exp":%d}`, time.Now().Add(5*time.Minute).Unix())
	})
	defer server.Close()

	repo := &stubUserRepo{}
	cacheStore := &stubAuthCache{getErr: errors.New("redis unavailable")}
	service := NewAuthService(repo, newTestConfig(server.URL, 5*time.Minute), cacheStore)

	response, err := service.ValidateToken(context.Background(), "token-2")
	if err != nil {
		t.Fatalf("ValidateToken returned error: %v", err)
	}

	if response.KeycloakId != "kc-2" {
		t.Fatalf("unexpected keycloak id %q", response.KeycloakId)
	}
	if cacheStore.setCalls != 1 {
		t.Fatalf("SetValidatedToken was called %d times, want 1", cacheStore.setCalls)
	}
}

func TestValidateTokenDoesNotCacheInactiveTokens(t *testing.T) {
	server := newIntrospectionServer(t, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"active":false}`)
	})
	defer server.Close()

	repo := &stubUserRepo{}
	cacheStore := &stubAuthCache{}
	service := NewAuthService(repo, newTestConfig(server.URL, 5*time.Minute), cacheStore)

	_, err := service.ValidateToken(context.Background(), "token-3")
	if !errors.Is(err, domain.ErrInvalidToken) {
		t.Fatalf("ValidateToken returned error %v, want %v", err, domain.ErrInvalidToken)
	}
	if repo.createCalls != 0 {
		t.Fatalf("Create was called %d times, want 0", repo.createCalls)
	}
	if cacheStore.setCalls != 0 {
		t.Fatalf("SetValidatedToken was called %d times, want 0", cacheStore.setCalls)
	}
}

func TestValidateTokenUsesTokenExpiryForCacheTTL(t *testing.T) {
	server := newIntrospectionServer(t, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"active":true,"sub":"kc-4","email":"user4@example.com","name":"Jane Doe","exp":%d}`, time.Now().Add(30*time.Second).Unix())
	})
	defer server.Close()

	cacheStore := &stubAuthCache{}
	service := NewAuthService(&stubUserRepo{}, newTestConfig(server.URL, 5*time.Minute), cacheStore)

	_, err := service.ValidateToken(context.Background(), "token-4")
	if err != nil {
		t.Fatalf("ValidateToken returned error: %v", err)
	}

	if cacheStore.setCalls != 1 {
		t.Fatalf("SetValidatedToken was called %d times, want 1", cacheStore.setCalls)
	}
	if cacheStore.setTTL <= 0 {
		t.Fatalf("expected positive ttl, got %s", cacheStore.setTTL)
	}
	if cacheStore.setTTL > 30*time.Second {
		t.Fatalf("expected ttl to be capped by token expiry, got %s", cacheStore.setTTL)
	}
}

func TestNewAuthServiceUsesNoopCacheByDefault(t *testing.T) {
	service := NewAuthService(&stubUserRepo{}, newTestConfig("http://example.com", 5*time.Minute), nil)

	if _, ok := service.(*authService); !ok {
		t.Fatal("expected concrete authService")
	}

	if _, ok := service.(*authService).cache.(cache.NoopAuthCache); !ok {
		t.Fatal("expected default cache to be NoopAuthCache")
	}
}

func newTestConfig(keycloakURL string, ttl time.Duration) *config.Config {
	return &config.Config{
		Keycloak: config.KeycloakConfig{
			Url:          keycloakURL,
			Realm:        "ecommerce",
			ClientID:     "ecommerce-service",
			ClientSecret: "secret",
		},
		Redis: config.RedisConfig{
			AuthCacheTTL: ttl,
		},
	}
}

func newIntrospectionServer(t *testing.T, handler func(http.ResponseWriter, *http.Request)) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/realms/ecommerce/protocol/openid-connect/token/introspect" {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
		handler(w, r)
	}))
}
