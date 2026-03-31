package service

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/yasinbozat/ecommerce-platform/services/user-service/internal/cache"
	"github.com/yasinbozat/ecommerce-platform/services/user-service/internal/config"
	"github.com/yasinbozat/ecommerce-platform/services/user-service/internal/domain"
	"github.com/yasinbozat/ecommerce-platform/services/user-service/repository"
)

type IAuthService interface {
	ValidateToken(ctx context.Context, token string) (*domain.ValidateTokenResponse, error)
	SyncUser(ctx context.Context, keycloakID, fullName, email string) (*domain.User, error)
}

type authService struct {
	userRepo repository.IUserRepository
	cfg      *config.Config
	cache    cache.AuthCache
}

func NewAuthService(userRepo repository.IUserRepository, cfg *config.Config, authCache cache.AuthCache) IAuthService {
	if authCache == nil {
		authCache = cache.NewNoopAuthCache()
	}

	return &authService{
		userRepo: userRepo,
		cfg:      cfg,
		cache:    authCache,
	}
}

type keycloakIntrospectResponse struct {
	Active bool   `json:"active"`
	Sub    string `json:"sub"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	Exp    int64  `json:"exp"`
}

func (s *authService) ValidateToken(ctx context.Context, token string) (*domain.ValidateTokenResponse, error) {
	token = strings.TrimSpace(token)
	if token == "" {
		return nil, domain.ErrInvalidToken
	}

	if cachedResponse, err := s.cache.GetValidatedToken(ctx, token); err == nil && cachedResponse != nil {
		return cachedResponse, nil
	}

	keycloak := s.cfg.Keycloak
	introspectURL := keycloak.Url + "/realms/" + keycloak.Realm + "/protocol/openid-connect/token/introspect"

	data := url.Values{}
	data.Set("client_id", keycloak.ClientID)
	data.Set("client_secret", keycloak.ClientSecret)
	data.Set("token", token)

	req, err := http.NewRequestWithContext(ctx, "POST", introspectURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 5 * time.Second}
	response, err := client.Do(req)
	if err != nil {
		return nil, domain.ErrKeycloakUnreachable
	}
	defer response.Body.Close()

	var introspect keycloakIntrospectResponse
	if err := json.NewDecoder(response.Body).Decode(&introspect); err != nil {
		return nil, err
	}

	if !introspect.Active {
		return nil, domain.ErrInvalidToken
	}

	if introspect.Exp > 0 && time.Until(time.Unix(introspect.Exp, 0)) <= 0 {
		return nil, domain.ErrTokenExpired
	}

	user, err := s.SyncUser(ctx, introspect.Sub, introspect.Name, introspect.Email)
	if err != nil {
		return nil, err
	}

	validateTokenResponse := &domain.ValidateTokenResponse{
		UserId:     user.Id,
		Email:      user.Email,
		Role:       user.Role,
		KeycloakId: user.KeycloakId,
	}

	if ttl := s.cacheTTL(introspect.Exp); ttl > 0 {
		_ = s.cache.SetValidatedToken(ctx, token, validateTokenResponse, ttl)
	}

	return validateTokenResponse, nil
}

func (s *authService) SyncUser(ctx context.Context, keycloakID, fullName, email string) (*domain.User, error) {
	user, err := s.userRepo.FindByKeycloakID(ctx, keycloakID)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return user, nil
	}

	// kullanıcı yoksa oluştur
	newUser := &domain.User{
		KeycloakId: keycloakID,
		Email:      email,
		FullName:   fullName,
		Role:       domain.RoleCustomer,
		IsActive:   true,
	}

	if err := s.userRepo.Create(ctx, newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}

func (s *authService) cacheTTL(exp int64) time.Duration {
	maxTTL := s.cfg.Redis.AuthCacheTTL
	if exp <= 0 {
		return maxTTL
	}

	untilExpiry := time.Until(time.Unix(exp, 0))
	if untilExpiry <= 0 {
		return 0
	}

	if maxTTL > 0 && untilExpiry > maxTTL {
		return maxTTL
	}

	return untilExpiry
}
