package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/yasinbozat/ecommerce-platform/services/user-service/internal/domain"
	"github.com/yasinbozat/ecommerce-platform/services/user-service/repository"
)

type IUserService interface {
	GetProfile(ctx context.Context, userID uuid.UUID) (*domain.UserResponse, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, req *domain.UpdateProfileRequest) (*domain.UserResponse, error)
	GetAddresses(ctx context.Context, userID uuid.UUID) ([]*domain.AddressResponse, error)
	CreateAddress(ctx context.Context, userID uuid.UUID, req *domain.CreateAddressRequest) (*domain.AddressResponse, error)
	UpdateAddress(ctx context.Context, userID, addressID uuid.UUID, req *domain.UpdateAddressRequest) (*domain.AddressResponse, error)
	DeleteAddress(ctx context.Context, userID, addressID uuid.UUID) error
	SetDefaultAddress(ctx context.Context, userID, addressID uuid.UUID) error
}

type userService struct {
	userRepo    repository.IUserRepository
	addressRepo repository.IAddressRepository
}

func NewUserService(userRepo repository.IUserRepository, addressRepo repository.IAddressRepository) IUserService {
	return &userService{
		userRepo:    userRepo,
		addressRepo: addressRepo,
	}
}

func (s *userService) GetProfile(ctx context.Context, userID uuid.UUID) (*domain.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, domain.ErrUserNotFound
	}
	return &domain.UserResponse{
		Id:       user.Id,
		Email:    user.Email,
		FullName: user.FullName,
		Phone:    user.Phone,
		Role:     user.Role,
	}, nil
}

func (s *userService) UpdateProfile(ctx context.Context, userID uuid.UUID, req *domain.UpdateProfileRequest) (*domain.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, domain.ErrUserNotFound
	}

	user.FullName = req.FullName
	user.Phone = req.Phone

	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return nil, err
	}
	return &domain.UserResponse{
		Id:       user.Id,
		Email:    user.Email,
		FullName: user.FullName,
		Phone:    user.Phone,
		Role:     user.Role,
	}, nil
}

func (s *userService) GetAddresses(ctx context.Context, userID uuid.UUID) ([]*domain.AddressResponse, error) {
	addresses, err := s.addressRepo.FindAllByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var response []*domain.AddressResponse
	for _, addr := range addresses {
		response = append(response, &domain.AddressResponse{
			Id:        addr.Id,
			Title:     addr.Title,
			FullName:  addr.FullName,
			Phone:     addr.Phone,
			Street:    addr.Street,
			District:  addr.District,
			City:      addr.City,
			ZipCode:   addr.ZipCode,
			IsDefault: addr.IsDefault,
			CreatedAt: addr.CreatedAt,
			UpdatedAt: addr.UpdatedAt,
		})
	}

	return response, nil
}

func (s *userService) CreateAddress(ctx context.Context, userID uuid.UUID, req *domain.CreateAddressRequest) (*domain.AddressResponse, error) {
	address := &domain.Address{
		UserID:    userID,
		Title:     req.Title,
		FullName:  req.FullName,
		Phone:     req.Phone,
		Street:    req.Street,
		District:  req.District,
		City:      req.City,
		ZipCode:   req.ZipCode,
		IsDefault: req.IsDefault,
	}

	err := s.addressRepo.Create(ctx, address)

	if err != nil {
		return nil, err
	}

	return &domain.AddressResponse{
		Id:        address.Id,
		Title:     address.Title,
		FullName:  address.FullName,
		Phone:     address.Phone,
		Street:    address.Street,
		District:  address.District,
		City:      address.City,
		ZipCode:   address.ZipCode,
		IsDefault: address.IsDefault,
	}, nil
}

func (s *userService) UpdateAddress(ctx context.Context, userID, addressID uuid.UUID, req *domain.UpdateAddressRequest) (*domain.AddressResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, domain.ErrUserNotFound
	}

	address, err := s.addressRepo.FindByID(ctx, addressID)
	if err != nil {
		return nil, err
	}

	if address == nil {
		return nil, domain.ErrAddressNotFound
	}

	if address.UserID != userID {
		return nil, domain.ErrAddressNotOwned
	}

	if req.Title != nil {
		address.Title = *req.Title
	}

	if req.FullName != nil {
		address.FullName = *req.FullName
	}
	if req.Phone != nil {
		address.Phone = *req.Phone
	}
	if req.Street != nil {
		address.Street = *req.Street
	}
	if req.District != nil {
		address.District = *req.District
	}
	if req.City != nil {
		address.City = *req.City
	}
	if req.ZipCode != nil {
		address.ZipCode = *req.ZipCode
	}
	if req.IsDefault != nil {
		address.IsDefault = *req.IsDefault
	}

	err = s.addressRepo.Update(ctx, address)
	if err != nil {
		return nil, err
	}
	return &domain.AddressResponse{
		Id:        address.Id,
		Title:     address.Title,
		FullName:  address.FullName,
		Phone:     address.Phone,
		Street:    address.Street,
		District:  address.District,
		City:      address.City,
		ZipCode:   address.ZipCode,
		IsDefault: address.IsDefault,
	}, nil
}

func (s *userService) DeleteAddress(ctx context.Context, userID, addressID uuid.UUID) error {
	if _, err := s.findAndValidateAddress(ctx, userID, addressID); err != nil {
		return err
	}
	return s.addressRepo.Delete(ctx, addressID)
}

func (s *userService) SetDefaultAddress(ctx context.Context, userID, addressID uuid.UUID) error {
	if _, err := s.findAndValidateAddress(ctx, userID, addressID); err != nil {
		return err
	}
	return s.addressRepo.SetDefault(ctx, userID, addressID)
}

func (s *userService) findAndValidateAddress(ctx context.Context, userID, addressID uuid.UUID) (*domain.Address, error) {
	address, err := s.addressRepo.FindByID(ctx, addressID)
	if err != nil {
		return nil, err
	}
	if address == nil {
		return nil, domain.ErrAddressNotFound
	}
	if address.UserID != userID {
		return nil, domain.ErrAddressNotOwned
	}
	return address, nil
}
