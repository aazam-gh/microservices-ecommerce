package domain

import "errors"

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrAddressNotFound     = errors.New("address not found")
	ErrAddressNotOwned     = errors.New("address not owned by user")
	ErrInvalidToken        = errors.New("invalid token")
	ErrTokenExpired        = errors.New("token expired")
	ErrKeycloakUnreachable = errors.New("keycloak unreachable")
	ErrInvalidInput        = errors.New("invalid input")
)
