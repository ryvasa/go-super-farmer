package database

import (
	"fmt"
	"os"

	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/pkg/database/seeders"
	"github.com/ryvasa/go-super-farmer/pkg/env"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ProvideDSN(cfg *env.Env) (string, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	// name := os.Getenv("DB_TEST")
	name := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	timezone := os.Getenv("DB_TIMEZONE")

	if host == "" || port == "" || name == "" || user == "" || password == "" {
		return "", fmt.Errorf("missing database environment variables")
	}

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s", host, user, password, name, port, timezone), nil
}

func ConnectDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(
		&domain.User{},
		&domain.Role{},
		&domain.Land{},
		&domain.Commodity{},
		&domain.LandCommodity{},
		&domain.Province{},
		&domain.City{},
		&domain.Region{},
		&domain.Price{},
		&domain.PriceHistory{},
		&domain.Supply{},
		&domain.SupplyHistory{},
		&domain.Demand{},
		&domain.DemandHistory{},
		&domain.Harvest{},
	)

	seeders.Seeders(db)

	return db, nil
}
