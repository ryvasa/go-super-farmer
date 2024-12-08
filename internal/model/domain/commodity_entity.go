package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Commodity struct {
	ID          uuid.UUID      `gorm:"primary_key;"`
	Name        string         `gorm:"size:100; not null; type:varchar(100);uniqueIndex" validate:"min=5"`
	Description string         `gorm:"size:255; not null; type:varchar(255);uniqueIndex" validate:"min=5"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}