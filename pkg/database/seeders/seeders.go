package seeders

import "gorm.io/gorm"

func Seeders(db *gorm.DB) {
	SeedRoles(db)
	SeedProvinces(db)
	SeedCities(db)
	SeedUsers(db)
	SeedCommodities(db)
}
