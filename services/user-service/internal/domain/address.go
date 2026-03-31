package domain

import (
	"time"

	"github.com/google/uuid"
)

type Address struct {
	Id        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	User      User      `gorm:"foreignKey:UserID"`
	Title     string    `gorm:"type:varchar(100);not null"`
	FullName  string    `gorm:"type:varchar(255);not null"`
	Phone     string    `gorm:"type:varchar(20);not null"`
	Street    string    `gorm:"type:text;not null"`
	District  string    `gorm:"type:varchar(100);not null"`
	City      string    `gorm:"type:varchar(100);not null"`
	ZipCode   string    `gorm:"type:varchar(10)"`
	IsDefault bool      `gorm:"default:false"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
