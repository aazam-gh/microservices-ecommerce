package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/yasinbozat/ecommerce-platform/services/user-service/internal/domain"
)

type IAddressRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Address, error)
	FindAllByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Address, error)
	Create(ctx context.Context, address *domain.Address) error
	Update(ctx context.Context, address *domain.Address) error
	Delete(ctx context.Context, id uuid.UUID) error
	SetDefault(ctx context.Context, id, userID uuid.UUID) error
}
