package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PriceHistory struct {
	ID          uuid.UUID      `gorm:"primary_key;"`
	CommodityID uuid.UUID      `gorm:"not null"`
	RegionID    uuid.UUID      `gorm:"not null"`
	Commodity   Commodity      `gorm:"foreignKey:CommodityID;references:ID"`
	Region      Region         `gorm:"foreignKey:RegionID;references:ID"`
	Price       float64        `gorm:"not null"`
	CreatedAt   string         `gorm:"autoCreateTime"`
	UpdatedAt   string         `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
