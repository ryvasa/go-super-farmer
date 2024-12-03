package domain

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	ID    int    `gorm:"primary_key; auto_increment"`
	Name  string `gorm:"size:100; not null; type:varchar(100);uniqueIndex"`
	Users []User `gorm:"foreignKey:RoleID"`
}
