package seeders

import (
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/pkg/logrus"
	"gorm.io/gorm"
)

func SeedCommodities(db *gorm.DB) {
	var count int64
	db.Model(&domain.Commodity{}).Count(&count)
	if count > 0 {
		logrus.Log.Info("Commodity data already seeded")
		return
	}

	for _, commodity := range commodities {
		result := db.Create(&commodity)
		if result.Error != nil {
			logrus.Log.Errorf("error seeding data Commodity: %v", result.Error)
		}
	}

	logrus.Log.Info("Commodity data seeded successfully")
}
