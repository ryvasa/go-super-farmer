package domain

type Province struct {
	ID   int64  `gorm:"primary_key;"`
	Name string `gorm:"not null, unique"`
}
