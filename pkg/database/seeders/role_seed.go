package seeders

import (
	"fmt"
	"log"

	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"gorm.io/gorm"
)

func SeedRoles(db *gorm.DB) {
	var count int64
	db.Model(&domain.Role{}).Count(&count)
	if count > 0 {
		fmt.Println("Role data already seeded")
		return
	}

	roles := []domain.Role{
		{ID: 1, Name: "Admin"},
		{ID: 2, Name: "Farmer"},
	}

	for _, role := range roles {
		result := db.Create(&role)
		if result.Error != nil {
			log.Fatalf("error seeding data Role: %v", result.Error)
		}
	}

	fmt.Println("Role data seeded successfully")
}
