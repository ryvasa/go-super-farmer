package seeders

import (
	"fmt"
	"log"

	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"gorm.io/gorm"
)

func SeedProvinces(db *gorm.DB) {
	var count int64
	db.Model(&domain.Province{}).Count(&count)
	if count > 0 {
		fmt.Println("Province data already seeded")
		return
	}

	provinces := []domain.Province{
		{Name: "Nanggroe Aceh Darussalam"},
		{Name: "Sumatera Utara"},
		{Name: "Sumatera Selatan"},
		{Name: "Sumatera Barat"},
		{Name: "Bengkulu"},
		{Name: "Riau"},
		{Name: "Kepulauan Riau"},
		{Name: "Jambi"},
		{Name: "Lampung"},
		{Name: "Bangka Belitung"},
		{Name: "Kalimantan Barat"},
		{Name: "Kalimantan Timur"},
		{Name: "Kalimantan Selatan"},
		{Name: "Kalimantan Tengah"},
		{Name: "Kalimantan Utara"},
		{Name: "Banten"},
		{Name: "DKI Jakarta"},
		{Name: "Jawa Barat"},
		{Name: "Jawa Tengah"},
		{Name: "Daerah Istimewa Yogyakarta"},
		{Name: "Jawa Timur"},
		{Name: "Bali"},
		{Name: "Nusa Tenggara Timur"},
		{Name: "Nusa Tenggara Barat"},
		{Name: "Gorontalo"},
		{Name: "Sulawesi Barat"},
		{Name: "Sulawesi Tengah"},
		{Name: "Sulawesi Utara"},
		{Name: "Sulawesi Tenggara"},
		{Name: "Sulawesi Selatan"},
		{Name: "Maluku Utara"},
		{Name: "Maluku"},
		{Name: "Papua Barat"},
		{Name: "Papua"},
		{Name: "Papua Tengah"},
		{Name: "Papua Pegunungan"},
		{Name: "Papua Selatan"},
		{Name: "Papua Barat Daya"},
	}

	for _, province := range provinces {
		result := db.Create(&province)
		if result.Error != nil {
			log.Fatalf("error seeding data Province: %v", result.Error)
		}
	}

	fmt.Println("Province data seeded successfully")
}
