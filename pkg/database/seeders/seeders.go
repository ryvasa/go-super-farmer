package seeders

import (
	"fmt"
	"log"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/pkg/logrus"
	"gorm.io/gorm"
)

// SeedRoles populates the roles table with predefined roles.
func SeedRoles(db *gorm.DB) []*domain.Role {
	predefinedRoles := []*domain.Role{
		{ID: 1, Name: "Admin"},
		{ID: 2, Name: "Farmer"},
	}

	if isDataSeeded(db, &domain.Role{}) {
		log.Println("Roles already seeded")
		return predefinedRoles
	}

	if err := db.Create(&predefinedRoles).Error; err != nil {
		log.Fatalf("Error seeding roles: %v", err)
	}
	return predefinedRoles
}

// SeedUsers populates the users table with fake data.
func SeedUsers(db *gorm.DB, roles []*domain.Role) []*domain.User {

	if err := db.Create(&users).Error; err != nil {
		log.Fatalf("Error seeding users: %v", err)
	}
	return users
}

// SeedCommodities populates the commodities table with fake data.
func SeedCommodities(db *gorm.DB) []*domain.Commodity {

	if err := db.Create(&commodities).Error; err != nil {
		log.Fatalf("Error seeding commodities: %v", err)
	}
	return commodities
}

func SeedLands(db *gorm.DB, users []*domain.User, cities []*domain.City) []*domain.Land {
	gofakeit.Seed(0)

	if err := db.Create(&lands).Error; err != nil {
		log.Fatalf("Error seeding lands: %v", err)
	}
	return lands
}

func SeedLandCommodities(db *gorm.DB, lands []*domain.Land, commodities []*domain.Commodity) []*domain.LandCommodity {
	var landCommodities []*domain.LandCommodity

	for _, land := range lands {
		var totalLandAreaUsed float64

		// Seed LandCommodity for harvested = false
		for _, commodity := range commodities {
			landArea := gofakeit.Float64Range(1.0, land.LandArea)

			harvested := gofakeit.Bool()
			if harvested == false {
				if totalLandAreaUsed+landArea > land.LandArea {
					landArea = land.LandArea - totalLandAreaUsed
				}
				totalLandAreaUsed += landArea
			}
			if landArea != 0 {
				createdAt := time.Time{}
				if harvested == true {
					createdAt = gofakeit.DateRange(time.Now().AddDate(0, 0, -120), time.Now().AddDate(0, 0, -90))
				} else {
					createdAt = gofakeit.DateRange(time.Now().AddDate(0, 0, -20), time.Now())
				}

				logrus.Log.Info(landArea, totalLandAreaUsed, land.LandArea)

				landCommodity := &domain.LandCommodity{
					ID:          uuid.New(),
					LandArea:    landArea,
					Unit:        "ha",
					CommodityID: commodity.ID,
					Commodity:   commodity,
					LandID:      land.ID,
					Land:        land,
					Harvested:   harvested,
					CreatedAt:   createdAt,
				}

				landCommodities = append(landCommodities, landCommodity)
			}
		}
	}

	if err := db.Create(&landCommodities).Error; err != nil {
		log.Fatalf("Error seeding land commodities: %v", err)
	}

	return landCommodities
}

func SeedProvinces(db *gorm.DB) []*domain.Province {
	gofakeit.Seed(0)

	if err := db.Create(&provinces).Error; err != nil {
		log.Fatalf("Error seeding provinces: %v", err)
	}
	return provinces
}

func SeedCities(db *gorm.DB) []*domain.City {
	gofakeit.Seed(0)

	if err := db.Create(&cities).Error; err != nil {
		log.Fatalf("Error seeding cities: %v", err)
	}
	return cities
}

func SeedHarvest(db *gorm.DB, landCommodity []*domain.LandCommodity) []*domain.Harvest {
	gofakeit.Seed(0)
	var harvestsList []*domain.Harvest
	for i := 0; i < len(landCommodity); i++ {
		if landCommodity[i].Harvested == true {

			var quantity float64
			switch landCommodity[i].Commodity.Code {
			case "RICE":
				quantity = landCommodity[i].LandArea * 5790
			case "CORN":
				quantity = landCommodity[i].LandArea * 7352
			case "SOYBEAN":
				quantity = landCommodity[i].LandArea * 2000
			case "WHEAT":
				quantity = landCommodity[i].LandArea * 1800
			}

			harvestsList = append(harvestsList, &domain.Harvest{
				ID:              uuid.New(),
				LandCommodityID: landCommodity[i].ID,
				LandCommodity:   landCommodity[i],
				Quantity:        quantity,
				Unit:            "kg",
				HarvestDate:     gofakeit.DateRange(time.Now().AddDate(0, 0, -30), time.Now()),
			})
		}
	}
	if err := db.Create(&harvestsList).Error; err != nil {
		log.Fatalf("Error seeding harvest: %v", err)
	}
	return harvestsList
}

func SeedSales(db *gorm.DB, harvests []*domain.Harvest) []*domain.Sale {
	gofakeit.Seed(0)
	var sales []*domain.Sale

	for _, harvest := range harvests {
		if harvest.LandCommodity.Harvested == true {
			remainingQuantity := harvest.Quantity

			// Tentukan kemungkinan apakah seluruh hasil panen akan terjual
			// Jika gofakeit.Bool() menghasilkan true, maka seluruh hasil panen terjual
			// Jika false, maka kita akan sisa beberapa hasil panen
			fullySold := gofakeit.Bool()

			// Jika fullySold true, semua hasil panen akan terjual
			if fullySold {
				sales = append(sales, &domain.Sale{
					ID:          uuid.New(),
					CommodityID: harvest.LandCommodity.CommodityID,
					Quantity:    remainingQuantity,
					Unit:        "kg",
					SaleDate:    gofakeit.DateRange(harvest.HarvestDate.AddDate(0, 0, 1), harvest.HarvestDate.AddDate(0, 0, 30)),
					CityID:      harvest.LandCommodity.Land.CityID,
					Price:       float64(1), // Tentukan harga yang sesuai
				})
			} else {
				// Jika tidak fully sold, lakukan beberapa kali penjualan hingga ada sisa
				minimumRemainingQuantity := gofakeit.Float64Range(0.1, 0.2) * harvest.Quantity
				for remainingQuantity > minimumRemainingQuantity {
					// Generate sale quantity
					saleQuantity := gofakeit.Float64Range(1.0, remainingQuantity)
					if saleQuantity > remainingQuantity {
						saleQuantity = remainingQuantity
					}

					// Kurangi remainingQuantity
					remainingQuantity -= saleQuantity

					// Generate sale
					sale := &domain.Sale{
						ID:          uuid.New(),
						CommodityID: harvest.LandCommodity.CommodityID,
						Quantity:    saleQuantity,
						Unit:        "kg",
						SaleDate:    gofakeit.DateRange(harvest.HarvestDate.AddDate(0, 0, 1), harvest.HarvestDate.AddDate(0, 0, 30)),
						CityID:      harvest.LandCommodity.Land.CityID,
						Price:       float64(1), // Tentukan harga yang sesuai
					}

					sales = append(sales, sale)
				}
			}
		}
	}

	if err := db.Create(&sales).Error; err != nil {
		log.Fatalf("Error seeding sales: %v", err)
	}
	return sales
}

func SeedDemand(db *gorm.DB, sales []*domain.Sale) []*domain.Demand {
	var demands []*domain.Demand

	// Hitung demand berdasarkan penjualan satu bulan terakhir
	var demandSummaries []struct {
		CommodityID uuid.UUID
		CityID      int64
		TotalSold   float64
	}
	lastMonth := time.Now().AddDate(0, -1, 0)
	db.Table("sales").
		Select("commodity_id, city_id, SUM(quantity) as total_sold").
		Where("sale_date BETWEEN ? AND ?", lastMonth, time.Now()).
		Group("commodity_id, city_id").
		Scan(&demandSummaries)

	// Buat data demand berdasarkan hasil pengelompokan
	for _, summary := range demandSummaries {
		demands = append(demands, &domain.Demand{
			ID:          uuid.New(),
			CommodityID: summary.CommodityID,
			CityID:      summary.CityID,
			Quantity:    summary.TotalSold,
			Unit:        "kg",
		})
	}

	// Simpan data demand ke database
	if err := db.Create(&demands).Error; err != nil {
		log.Fatalf("Error seeding demand: %v", err)
	}
	return demands
}

func SeedDemandHistory(db *gorm.DB, demands []*domain.Demand) []*domain.DemandHistory {
	gofakeit.Seed(0)
	var demandHistories []*domain.DemandHistory

	for _, demand := range demands {
		numHistories := gofakeit.Number(5, 10) // Tentukan jumlah data sejarah per demand

		for i := 0; i < numHistories; i++ {
			quantity := gofakeit.Float64Range(demand.Quantity-3000, demand.Quantity+3000)
			createdAt := gofakeit.DateRange(demand.CreatedAt.AddDate(0, 0, -100), demand.CreatedAt.AddDate(0, 0, -10))

			demandHistories = append(demandHistories, &domain.DemandHistory{
				ID:          uuid.New(),
				Quantity:    quantity,
				Unit:        "kg",
				CommodityID: demand.CommodityID,
				CityID:      demand.CityID,
				CreatedAt:   createdAt,
			})
		}
	}

	if err := db.Create(&demandHistories).Error; err != nil {
		log.Fatalf("Error seeding demand history: %v", err)
	}
	return demandHistories
}

func SeedSupply(db *gorm.DB, harvests []*domain.Harvest) []*domain.Supply {
	gofakeit.Seed(0)
	var supplies []*domain.Supply

	// Hitung demand berdasarkan penjualan satu bulan terakhir
	var supplySummaries []struct {
		CommodityID uuid.UUID
		CityID      int64
		TotalSold   float64
	}

	lastMonth := time.Now().AddDate(0, -1, 0)
	db.Table("harvests").
		Select("land_commodities.commodity_id, lands.city_id, SUM(harvests.quantity) as total_sold").
		Joins("JOIN land_commodities ON land_commodities.id = harvests.land_commodity_id").
		Joins("JOIN lands ON lands.id = land_commodities.land_id").
		Where("harvest_date BETWEEN ? AND ?", lastMonth, time.Now()).
		Group("land_commodities.commodity_id, lands.city_id").
		Scan(&supplySummaries)

	// Peta untuk mencari data `commodity` dan `city` berdasarkan `CommodityID` dan `CityID`
	harvestMap := make(map[string]*domain.Harvest)
	for _, harvest := range harvests {
		key := fmt.Sprintf("%s-%d", harvest.LandCommodity.Commodity.ID, harvest.LandCommodity.Land.CityID)
		harvestMap[key] = harvest
	}

	// Buat data supply berdasarkan hasil pengelompokan dan informasi tambahan dari harvest
	for _, summary := range supplySummaries {
		key := fmt.Sprintf("%s-%d", summary.CommodityID, summary.CityID)
		harvest, exists := harvestMap[key]
		if !exists {
			log.Printf("Tidak ada data harvest untuk CommodityID %s dan CityID %d. Lewati.", summary.CommodityID, summary.CityID)
			continue
		}

		supplies = append(supplies, &domain.Supply{
			ID:          uuid.New(),
			CommodityID: summary.CommodityID,
			Commodity:   harvest.LandCommodity.Commodity,
			City:        harvest.LandCommodity.Land.City,
			CityID:      summary.CityID,
			Quantity:    summary.TotalSold,
			Unit:        "kg",
		})
	}

	// Simpan data supply ke database
	if err := db.Create(&supplies).Error; err != nil {
		log.Fatalf("Error seeding supply: %v", err)
	}
	return supplies
}

func SeedSupplyHistory(db *gorm.DB, supplies []*domain.Supply) []*domain.SupplyHistory {
	gofakeit.Seed(0)
	var supplyHistories []*domain.SupplyHistory

	for _, supply := range supplies {
		numHistories := gofakeit.Number(5, 10) // Tentukan jumlah data sejarah per supply

		for i := 0; i < numHistories; i++ {
			quantity := gofakeit.Float64Range(supply.Quantity-3000, supply.Quantity+3000)

			createdAt := gofakeit.DateRange(supply.CreatedAt.AddDate(0, 0, -100), supply.CreatedAt.AddDate(0, 0, -10))

			supplyHistories = append(supplyHistories, &domain.SupplyHistory{
				ID:          uuid.New(),
				Quantity:    quantity,
				Unit:        "kg",
				CommodityID: supply.CommodityID,
				CityID:      supply.CityID,
				CreatedAt:   createdAt,
			})
		}
	}

	if err := db.Create(&supplyHistories).Error; err != nil {
		log.Fatalf("Error seeding supply history: %v", err)
	}
	return supplyHistories
}
func SeedPrices(db *gorm.DB, supplies []*domain.Supply, demands []*domain.Demand) []*domain.Price {
	gofakeit.Seed(0)

	// Harga rata-rata untuk komoditas berdasarkan kode
	commodityPrices := map[string]float64{
		"RICE":    6431.11,
		"CORN":    5005.583,
		"SOYBEAN": 10020.68,
		"WHEAT":   23020.68,
		// Tambahkan komoditas lainnya sesuai kebutuhan
	}

	// Data untuk menyimpan harga yang akan di-seed
	var prices []*domain.Price

	// Gunakan kombinasi unik dari CommodityID dan CityID untuk memastikan hanya ada satu harga per komoditas per kota
	priceMap := make(map[string]bool)

	// Iterasi melalui supplies untuk menentukan harga
	for _, supply := range supplies {
		key := fmt.Sprintf("%s-%d", supply.CommodityID, supply.CityID) // Kombinasi unik
		if _, exists := priceMap[key]; exists {
			continue // Lewati jika harga sudah dibuat untuk kombinasi ini
		}

		priceMap[key] = true // Tandai kombinasi sebagai telah diproses

		// Ambil harga rata-rata berdasarkan kode komoditas
		averagePrice, ok := commodityPrices[supply.Commodity.Code]
		if !ok {
			log.Printf("Harga rata-rata untuk komoditas %s tidak ditemukan. Lewati.", supply.Commodity.Code)
			continue
		}

		prices = append(prices, &domain.Price{
			ID:          uuid.New(),
			CommodityID: supply.CommodityID,
			CityID:      supply.CityID,
			Price:       gofakeit.Float64Range(averagePrice-100, averagePrice+100), // Variasi kecil untuk realisme
			CreatedAt:   time.Now(),
		})
	}

	// Simpan data prices ke database
	if err := db.Create(&prices).Error; err != nil {
		log.Fatalf("Error seeding prices: %v", err)
	}
	return prices
}

func SeedPriceHistory(db *gorm.DB, prices []*domain.Price, supplies []*domain.Supply, demands []*domain.Demand) []*domain.PriceHistory {
	gofakeit.Seed(0)
	var priceHistories []*domain.PriceHistory

	// Buat map supply dan demand berdasarkan CommodityID dan CityID
	supplyMap := make(map[string]float64)
	demandMap := make(map[string]float64)

	for _, supply := range supplies {
		key := fmt.Sprintf("%s-%d", supply.CommodityID, supply.CityID)
		supplyMap[key] = supply.Quantity
	}

	for _, demand := range demands {
		key := fmt.Sprintf("%s-%d", demand.CommodityID, demand.CityID)
		demandMap[key] = demand.Quantity
	}

	// Iterasi melalui daftar harga
	for _, price := range prices {
		key := fmt.Sprintf("%s-%d", price.CommodityID, price.CityID)
		baseSupply := supplyMap[key]
		baseDemand := demandMap[key]

		// Tentukan jumlah riwayat per harga
		numHistories := gofakeit.Number(5, 10)

		for i := 0; i < numHistories; i++ {
			// Variasi harga berdasarkan supply dan demand
			supplyEffect := gofakeit.Float64Range(-0.05, 0.05) * baseSupply
			demandEffect := gofakeit.Float64Range(-0.05, 0.05) * baseDemand
			priceChange := demandEffect - supplyEffect

			// Harga historis yang baru
			historyPrice := price.Price + priceChange

			// Pastikan harga tetap masuk akal (tidak negatif)
			if historyPrice < 0 {
				historyPrice = gofakeit.Float64Range(10, 50) // Harga minimum jika hasil negatif
			}

			// Tentukan tanggal riwayat dalam rentang waktu tertentu
			createdAt := gofakeit.DateRange(price.CreatedAt.AddDate(0, 0, -100), price.CreatedAt)

			// Tambahkan riwayat harga ke daftar
			priceHistories = append(priceHistories, &domain.PriceHistory{
				ID:          uuid.New(),
				Price:       historyPrice,
				CommodityID: price.CommodityID,
				CityID:      price.CityID,
				CreatedAt:   createdAt,
			})
		}
	}

	// Simpan riwayat harga ke database
	if err := db.Create(&priceHistories).Error; err != nil {
		log.Fatalf("Error seeding price history: %v", err)
	}
	return priceHistories
}

// Common function to check if data is already seeded.
func isDataSeeded(db *gorm.DB, model interface{}) bool {
	var count int64
	db.Model(model).Count(&count)
	return count > 0
}

func CleanAll(db *gorm.DB) {
	tables := []string{
		"users",
		"roles",
		"commodities",
		"provinces",
		"cities",
		"lands",
		"land_commodities",
		"demands",
		"demand_histories",
		"prices",
		"price_histories",
		"supplies",
		"supply_histories",
		"harvests",
	}

	for _, table := range tables {
		db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
	}
}

// Main Seeder function.
func Seeders(db *gorm.DB) {
	CleanAll(db)
	roles := SeedRoles(db)
	users := SeedUsers(db, roles)
	commodities := SeedCommodities(db)
	SeedProvinces(db)
	cities := SeedCities(db)
	lands := SeedLands(db, users, cities)
	landCommodities := SeedLandCommodities(db, lands, commodities)
	harvests := SeedHarvest(db, landCommodities)
	sales := SeedSales(db, harvests)
	demands := SeedDemand(db, sales)
	supply := SeedSupply(db, harvests)
	SeedSupplyHistory(db, supply)
	SeedDemandHistory(db, demands)
	prices := SeedPrices(db, supply, demands)
	SeedPriceHistory(db, prices, supply, demands)
	logrus.Log.Info("Seeding commpleted successfylly!")
}
