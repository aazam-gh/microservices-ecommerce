package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/yasinbozat/ecommerce-platform/services/user-service/internal/domain"
	"github.com/yasinbozat/ecommerce-platform/services/user-service/repository"
	"gorm.io/gorm"
)

type addressPostgres struct {
	db *gorm.DB
}

func NewAddressRepository(db *gorm.DB) repository.IAddressRepository {
	return &addressPostgres{db: db}
}

func (r *addressPostgres) FindByID(ctx context.Context, id uuid.UUID) (*domain.Address, error) {
	var address domain.Address
	result := r.db.WithContext(ctx).First(address, "id = ?", id)
	if result.Error != gorm.ErrRecordNotFound {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &address, nil
}
func (r *addressPostgres) FindAllByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Address, error) {
	var address []*domain.Address
	result := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(address)
	if result.Error != nil {
		return nil, result.Error
	}
	return address, nil
}
func (r *addressPostgres) Create(ctx context.Context, address *domain.Address) error {
	return r.db.WithContext(ctx).Create(address).Error
}
func (r *addressPostgres) Update(ctx context.Context, address *domain.Address) error {
	return r.db.WithContext(ctx).Save(address).Error
}
func (r *addressPostgres) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.Address{}).Error
}
func (r *addressPostgres) SetDefault(ctx context.Context, id, userID uuid.UUID) error {
	// first set all to false
	err := r.db.WithContext(ctx).
		Model(&domain.Address{}).
		Where("user_id = ?", userID).
		Update("is_default", false).Error
	if err != nil {
		return err
	}

	// then set the selected one to true
	return r.db.WithContext(ctx).
		Model(&domain.Address{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("is_default", true).Error
}
