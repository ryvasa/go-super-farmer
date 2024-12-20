package seeders

import (
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/pkg/logrus"
	"gorm.io/gorm"
)

func SeedProvinces(db *gorm.DB) {
	var count int64
	db.Model(&domain.Province{}).Count(&count)
	if count > 0 {
		logrus.Log.Info("Province data already seeded")
		return
	}

	for _, province := range provinces {
		result := db.Create(&province)
		if result.Error != nil {
			logrus.Log.Errorf("error seeding data Province: %v", result.Error)
		}
	}

	logrus.Log.Info("Province data seeded successfully")
}
