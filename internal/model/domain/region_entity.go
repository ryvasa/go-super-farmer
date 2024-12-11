package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Region struct {
	ID         uuid.UUID      `gorm:"primary_key;"`
	ProvinceID int64          `gorm:"not null;type:bigint"`
	Province   *Province      `gorm:"foreignkey:ProvinceID" json:"province,omitempty"`
	CityID     int64          `gorm:"not null;type:bigint"`
	City       *City          `gorm:"foreignkey:CityID" json:"city,omitempty"`
	CreatedAt  time.Time      `gorm:"autoCreateTime"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}
