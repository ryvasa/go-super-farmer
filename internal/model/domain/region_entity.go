package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Region struct {
	ID         uuid.UUID      `gorm:"primary_key;"`
	ProvinceID int64          `gorm:"not null"`
	Province   Province       `gorm:"foreignkey:ProvinceID"`
	CityID     int64          `gorm:"not null"`
	City       City           `gorm:"foreignkey:CityID"`
	CreatedAt  time.Time      `gorm:"autoCreateTime"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}
