package domain

import (
	"time"

	"github.com/google/uuid"
)

type UpdateProfileRequest struct {
	FullName string `json:"full_name" validate:"required,max=255"`
	Phone    string `json:"phone" validate:"required,max=20"`
}

type UserResponse struct {
	Id       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	FullName string    `json:"full_name"`
	Phone    string    `json:"phone"`
	Role     Role      `json:"role"`
}

type CreateAddressRequest struct {
	Title     string `json:"title" validate:"required,max=100"`
	FullName  string `json:"full_name" validate:"required,max=255"`
	Phone     string `json:"phone" validate:"required,max=20"`
	Street    string `json:"street" validate:"required"`
	District  string `json:"district" validate:"required,max=100"`
	City      string `json:"city" validate:"required,max=100"`
	ZipCode   string `json:"zip_code" validate:"required,max=10"`
	IsDefault bool   `json:"is_default"`
}

type UpdateAddressRequest struct {
	Title     *string `json:"title" validate:"omitempty,max=100"`
	FullName  *string `json:"full_name" validate:"omitempty,max=255"`
	Phone     *string `json:"phone" validate:"omitempty,max=20"`
	Street    *string `json:"street" validate:"omitempty"`
	District  *string `json:"district" validate:"omitempty,max=100"`
	City      *string `json:"city" validate:"omitempty,max=100"`
	ZipCode   *string `json:"zip_code" validate:"omitempty,max=10"`
	IsDefault *bool   `json:"is_default"`
}

type AddressResponse struct {
	Id        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	FullName  string    `json:"full_name"`
	Phone     string    `json:"phone"`
	Street    string    `json:"street"`
	District  string    `json:"district"`
	City      string    `json:"city"`
	ZipCode   string    `json:"zip_code"`
	IsDefault bool      `json:"is_default"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ValidateTokenResponse struct {
	UserId     uuid.UUID `json:"user_id"`
	Email      string    `json:"email"`
	Role       Role      `json:"role"`
	KeycloakId string    `json:"keycloak_id"`
}
