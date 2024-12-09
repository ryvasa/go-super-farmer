package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Price struct {
	ID          uuid.UUID      `gorm:"primary_key;"`
	CommodityID uuid.UUID      `gorm:"not null,uniqueIndex"`
	Commodity   Commodity      `gorm:"foreignKey:CommodityID;references:ID"`
	RegionID    uuid.UUID      `gorm:"not null,uniqueIndex"`
	Region      Region         `gorm:"foreignKey:RegionID;references:ID"`
	Price       float64        `gorm:"not null"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
