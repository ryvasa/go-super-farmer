package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Land struct {
	ID          uuid.UUID      `gorm:"primaryKey;type:varchar(255)"`
	UserID      uuid.UUID      `gorm:"not null;type:varchar(255)"`
	User        User           `gorm:"foreignKey:UserID" json:"-"`
	LandArea    float64        `gorm:"not null;type:bigint"`
	Certificate string         `gorm:"not null;type:varchar(255)"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
