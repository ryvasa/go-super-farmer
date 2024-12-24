package seeders

import (
	"log"
	"strconv"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/pkg/logrus"
	"github.com/ryvasa/go-super-farmer/service_api/model/domain"
	"gorm.io/gorm"
)

// SeedRoles populates the roles table with predefined roles.
func SeedRoles(db *gorm.DB) []domain.Role {
	predefinedRoles := []domain.Role{
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
func SeedUsers(db *gorm.DB, roles []domain.Role, count int) []domain.User {
	gofakeit.Seed(0)
	var users []domain.User

	for i := 0; i < count; i++ {
		role := roles[gofakeit.Number(0, len(roles)-1)]
		users = append(users, domain.User{
			ID:       uuid.New(),
			RoleID:   role.ID,
			Name:     gofakeit.Name(),
			Email:    gofakeit.Email(),
			Password: gofakeit.Password(true, true, true, true, false, 12),
			Phone:    pointer(gofakeit.Phone()),
		})
	}

	if err := db.Create(&users).Error; err != nil {
		log.Fatalf("Error seeding users: %v", err)
	}
	return users
}

// SeedCommodities populates the commodities table with fake data.
func SeedCommodities(db *gorm.DB, count int) []domain.Commodity {
	gofakeit.Seed(0)
	var commodities []domain.Commodity
	usedNames := make(map[string]bool)

	for i := 0; i < count; i++ {
		uniqueName := generateUniqueName(usedNames)
		commodities = append(commodities, domain.Commodity{
			ID:          uuid.New(),
			Name:        uniqueName,
			Description: gofakeit.Sentence(10),
			Code:        gofakeit.LetterN(8),
			Duration:    strconv.Itoa(gofakeit.Number(1, 100)),
		})
	}

	if err := db.Create(&commodities).Error; err != nil {
		log.Fatalf("Error seeding commodities: %v", err)
	}
	return commodities
}

func SeedLands(db *gorm.DB, users []domain.User, count int) []domain.Land {
	gofakeit.Seed(0)
	var lands []domain.Land

	for i := 0; i < count; i++ {
		user := users[gofakeit.Number(0, len(users)-1)]
		lands = append(lands, domain.Land{
			ID:          uuid.New(),
			UserID:      user.ID,
			LandArea:    gofakeit.Float64Range(1.0, 100.0),
			Unit:        "ha",
			Certificate: gofakeit.Word() + "-" + gofakeit.LetterN(4),
		})
	}

	if err := db.Create(&lands).Error; err != nil {
		log.Fatalf("Error seeding lands: %v", err)
	}
	return lands
}

func SeedLandCommodities(db *gorm.DB, lands []domain.Land, commodities []domain.Commodity, count int) []domain.LandCommodity {
	gofakeit.Seed(0)
	var landCommodities []domain.LandCommodity
	for i := 0; i < count; i++ {
		land := lands[gofakeit.Number(0, len(lands)-1)]
		commodity := commodities[gofakeit.Number(0, len(commodities)-1)]
		landCommodity := domain.LandCommodity{
			ID:          uuid.New(),
			LandArea:    gofakeit.Float64Range(1.0, 10.0),
			Unit:        "ha",
			CommodityID: commodity.ID,
			LandID:      land.ID,
		}
		landCommodities = append(landCommodities, landCommodity)
		if err := db.Create(&landCommodity).Error; err != nil {
			log.Fatalf("Error seeding land commodities: %v", err)
		}
	}
	return landCommodities
}

func SeedProvinces(db *gorm.DB, count int) []domain.Province {
	gofakeit.Seed(0)
	var provinces []domain.Province
	for i := 0; i < count; i++ {
		provinces = append(provinces, domain.Province{
			Name: gofakeit.Word(),
		})
	}
	if err := db.Create(&provinces).Error; err != nil {
		log.Fatalf("Error seeding provinces: %v", err)
	}
	return provinces
}

func SeedCities(db *gorm.DB, provinces []domain.Province, count int) []domain.City {
	gofakeit.Seed(0)
	var cities []domain.City
	for i := 0; i < count; i++ {
		province := provinces[gofakeit.Number(0, len(provinces)-1)]
		cities = append(cities, domain.City{
			Name:       gofakeit.Word(),
			ProvinceID: province.ID,
		})
	}
	if err := db.Create(&cities).Error; err != nil {
		log.Fatalf("Error seeding cities: %v", err)
	}
	return cities
}

func SeedDemand(db *gorm.DB, commodities []domain.Commodity, citys []domain.City, count int) []domain.Demand {
	gofakeit.Seed(0)
	var demands []domain.Demand
	for i := 0; i < count; i++ {
		commodity := commodities[gofakeit.Number(0, len(commodities)-1)]
		city := citys[gofakeit.Number(0, len(citys)-1)]
		demands = append(demands, domain.Demand{
			ID:          uuid.New(),
			CommodityID: commodity.ID,
			CityID:      city.ID,
			Quantity:    gofakeit.Float64Range(1.0, 100000.0),
			Unit:        "kg",
		})
	}
	if err := db.Create(&demands).Error; err != nil {
		log.Fatalf("Error seeding demand: %v", err)
	}
	return demands
}

func SeedDemandHistory(db *gorm.DB, demands []domain.Demand, count int) []domain.DemandHistory {
	gofakeit.Seed(0)
	var demandHistories []domain.DemandHistory
	for i := 0; i < count; i++ {
		demand := demands[gofakeit.Number(0, len(demands)-1)]
		demandHistories = append(demandHistories, domain.DemandHistory{
			ID:          uuid.New(),
			Quantity:    gofakeit.Float64Range(1.0, 100000.0),
			Unit:        "kg",
			CommodityID: demand.CommodityID,
			CityID:      demand.CityID,
		})
	}
	if err := db.Create(&demandHistories).Error; err != nil {
		log.Fatalf("Error seeding demand history: %v", err)
	}
	return demandHistories
}

func SeedPrices(db *gorm.DB, commodities []domain.Commodity, cities []domain.City, count int) []domain.Price {
	gofakeit.Seed(0)
	var prices []domain.Price
	for i := 0; i < count; i++ {
		commodity := commodities[gofakeit.Number(0, len(commodities)-1)]
		city := cities[gofakeit.Number(0, len(cities)-1)]
		prices = append(prices, domain.Price{
			ID:          uuid.New(),
			CommodityID: commodity.ID,
			CityID:      city.ID,
			Price:       gofakeit.Float64Range(1.0, 100000.0),
			Unit:        "kg",
		})
	}
	if err := db.Create(&prices).Error; err != nil {
		log.Fatalf("Error seeding prices: %v", err)
	}
	return prices
}

func SeedPriceHistory(db *gorm.DB, prices []domain.Price, count int) []domain.PriceHistory {
	gofakeit.Seed(0)
	var priceHistories []domain.PriceHistory
	for i := 0; i < count; i++ {
		price := prices[gofakeit.Number(0, len(prices)-1)]
		priceHistories = append(priceHistories, domain.PriceHistory{
			ID:          uuid.New(),
			Price:       gofakeit.Float64Range(1.0, 100000.0),
			Unit:        "kg",
			CommodityID: price.CommodityID,
			CityID:      price.CityID,
		})
	}
	if err := db.Create(&priceHistories).Error; err != nil {
		log.Fatalf("Error seeding price history: %v", err)
	}
	return priceHistories
}

func SeedSupply(db *gorm.DB, commodities []domain.Commodity, cities []domain.City, count int) []domain.Supply {
	gofakeit.Seed(0)
	var supplies []domain.Supply
	for i := 0; i < count; i++ {
		commodity := commodities[gofakeit.Number(0, len(commodities)-1)]
		city := cities[gofakeit.Number(0, len(cities)-1)]
		supplies = append(supplies, domain.Supply{
			ID:          uuid.New(),
			CommodityID: commodity.ID,
			CityID:      city.ID,
			Quantity:    gofakeit.Float64Range(1.0, 100000.0),
			Unit:        "kg",
		})
	}
	if err := db.Create(&supplies).Error; err != nil {
		log.Fatalf("Error seeding supply: %v", err)
	}
	return supplies
}

func SeedSupplyHistory(db *gorm.DB, supplies []domain.Supply, count int) []domain.SupplyHistory {
	gofakeit.Seed(0)
	var supplyHistories []domain.SupplyHistory
	for i := 0; i < count; i++ {
		supply := supplies[gofakeit.Number(0, len(supplies)-1)]
		supplyHistories = append(supplyHistories, domain.SupplyHistory{
			ID:          uuid.New(),
			Quantity:    gofakeit.Float64Range(1.0, 100000.0),
			Unit:        "kg",
			CommodityID: supply.CommodityID,
			CityID:      supply.CityID,
		})
	}
	if err := db.Create(&supplyHistories).Error; err != nil {
		log.Fatalf("Error seeding supply history: %v", err)
	}
	return supplyHistories
}

func SeedHarvest(db *gorm.DB, harvests []domain.LandCommodity, cities []domain.City) []domain.Harvest {
	gofakeit.Seed(0)
	var harvestsList []domain.Harvest
	for i := 0; i < len(harvests); i++ {
		harvest := harvests[i]
		city := cities[gofakeit.Number(0, len(cities)-1)]
		harvestsList = append(harvestsList, domain.Harvest{
			ID:              uuid.New(),
			LandCommodityID: harvest.ID,
			Quantity:        gofakeit.Float64Range(1.0, 100000.0),
			Unit:            "kg",
			CityID:          city.ID,
		})
	}
	if err := db.Create(&harvestsList).Error; err != nil {
		log.Fatalf("Error seeding harvest: %v", err)
	}
	return harvestsList
}

// Common function to check if data is already seeded.
func isDataSeeded(db *gorm.DB, model interface{}) bool {
	var count int64
	db.Model(model).Count(&count)
	return count > 0
}

// Helper function to generate unique names.
func generateUniqueName(usedNames map[string]bool) string {
	for {
		name := gofakeit.Word() + "_" + gofakeit.LetterN(4)
		if !usedNames[name] {
			usedNames[name] = true
			return name
		}
	}
}

// Helper function to create a pointer for a value.
func pointer[T any](value T) *T {
	return &value
}

// Main Seeder function.
func Seeders(db *gorm.DB) {
	roles := SeedRoles(db)
	users := SeedUsers(db, roles, 5)
	commodities := SeedCommodities(db, 10)
	lands := SeedLands(db, users, 10)
	landCommodities := SeedLandCommodities(db, lands, commodities, 10)
	provinces := SeedProvinces(db, 10)
	cities := SeedCities(db, provinces, 5)
	demand := SeedDemand(db, commodities, cities, 10)
	SeedDemandHistory(db, demand, 10)
	prices := SeedPrices(db, commodities, cities, 10)
	SeedPriceHistory(db, prices, 10)
	supply := SeedSupply(db, commodities, cities, 10)
	SeedSupplyHistory(db, supply, 10)
	SeedHarvest(db, landCommodities, cities)
	logrus.Log.Info("Seeding commpleted successfylly!")
}
