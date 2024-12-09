package domain

type City struct {
	ID         int64    `gorm:"primary_key;"`
	ProvinceID int64    `gorm:"not null"`
	Province   Province `gorm:"foreignkey:ProvinceID"`
	Name       string   `gorm:"not null"`
}
