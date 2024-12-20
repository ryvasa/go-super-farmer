package seeders

import (
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/pkg/logrus"
	"gorm.io/gorm"
)

func SeedUsers(db *gorm.DB) {
	var count int64
	db.Model(&domain.User{}).Count(&count)
	if count > 0 {
		logrus.Log.Info("Users already seeded")
		return
	}

	for _, user := range users {
		result := db.Create(&user)
		if result.Error != nil {
			logrus.Log.Errorf("error seeding data user: %v", result.Error)
		}
	}

	logrus.Log.Info("Users seeded successfully")
}
