package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        uuid.UUID      `gorm:"type:uuid;primary_key; default:uuid_generate_v4()"`
	Name      string         `gorm:"size:100; not null; type:varchar(100)"`
	Email     string         `gorm:"unique; not null; type:varchar(255)"`
	Password  string         `gorm:"not null type:varchar(255)"`
	RoleID    int            `gorm:"not null;default:1"` // Foreign key
	Role      Role           `gorm:"foreignKey:RoleID"`  //untuk melakukan eager loading
	Phone     *string        `gorm:"type:varchar(20)"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
