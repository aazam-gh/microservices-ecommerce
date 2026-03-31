package domain

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleCustomer Role = "customer"
	RoleAdmin    Role = "admin"
)

type User struct {
	Id         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Email      string    `gorm:"type:varchar(255);unique;not null"`
	FullName   string    `gorm:"type:varchar(255);not null"`
	Phone      string    `gorm:"type:varchar(20)"`
	Role       Role      `gorm:"type:varchar(20);default:'customer'"`
	KeycloakId string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	IsActive   bool      `gorm:"not null;default:true"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
