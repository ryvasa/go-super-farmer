package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LandCommoditiesEntity struct {
	ID         uuid.UUID      `gorm:"primary_key;"`
	LandArea   float64        `gorm:"not null"`
	ComodityID uuid.UUID      `gorm:"not null;type:varchar(255)"`
	LandID     uuid.UUID      `gorm:"not null;type:varchar(255)"`
	Commodity  Commodity      `gorm:"foreignKey:ComodityID"`
	Land       Land           `gorm:"foreignKey:LandID"`
	CreatedAt  time.Time      `gorm:"autoCreateTime"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}
