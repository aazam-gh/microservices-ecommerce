package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/yasinbozat/ecommerce-platform/services/user-service/internal/domain"
)

type IUserRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	FindByKeycloakID(ctx context.Context, keycloakID string) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) error
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}
