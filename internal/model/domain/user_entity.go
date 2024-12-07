package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID      `gorm:"primaryKey;type:varchar(255)"`
	Name      string         `gorm:"size:100;not null;type:varchar(100)"`
	Email     string         `gorm:"unique;not null;type:varchar(255)"`
	Password  string         `gorm:"not null;type:varchar(255)"`
	RoleID    int            `gorm:"not null;default:1"`
	Role      Role           `gorm:"foreignKey:RoleID"`
	Phone     *string        `gorm:"type:varchar(20)"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
