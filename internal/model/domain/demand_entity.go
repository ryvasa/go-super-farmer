package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DemandHistory struct {
	ID          uuid.UUID      `gorm:"primaryKey;type:varchar(36)"`
	CommodityID uuid.UUID      `gorm:"not null;type:varchar(36)"`
	Commodity   *Commodity     `gorm:"foreignKey:CommodityID" json:"commodity,omitempty"`
	RegionID    uuid.UUID      `gorm:"not null;type:varchar(36)"`
	Region      *Region        `gorm:"foreignKey:RegionID" json:"region,omitempty"`
	Quantity    float64        `gorm:"not null;type:int"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
