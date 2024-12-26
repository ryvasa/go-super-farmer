package helper

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/pkg/logrus"
	"github.com/ryvasa/go-super-farmer/service_api/model/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBHelper interface {
	TeardownTestDB()
	CreateUser(id uuid.UUID, name string, email string, password string) error
	FindUserByEmail(email string) (*domain.User, error)
	DeleteUser(id uuid.UUID) error
}
type DBHelperImpl struct {
}

func NewDBHelper() DBHelper {
	return &DBHelperImpl{}
}

func connectDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s", "localhost", "postgres", "123", "go_super_farmer", "5432", "Asia/Jakarta")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (e *DBHelperImpl) TeardownTestDB() {
	db, err := connectDB()
	if err != nil {
		logrus.Log.Fatal(err)
	}

	// Gunakan transaksi untuk memastikan atomicity
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err = tx.Exec("DELETE FROM harvests").Error
	if err != nil {
		tx.Rollback()
		panic(err) // Atau penanganan error yang lebih sesuai
	}

	err = tx.Exec("DELETE FROM land_commodities").Error
	if err != nil {
		tx.Rollback()
		panic(err) // Atau penanganan error yang lebih sesuai
	}

	err = tx.Exec("DELETE FROM lands").Error
	if err != nil {
		tx.Rollback()
		panic(err) // Atau penanganan error yang lebih sesuai
	}

	err = tx.Exec("DELETE FROM users").Error
	if err != nil {
		tx.Rollback()
		panic(err) // Atau penanganan error yang lebih sesuai
	}

	tx.Commit()
}

func (e *DBHelperImpl) CreateUser(id uuid.UUID, name string, email string, password string) error {
	db, err := connectDB()
	if err != nil {
		return err
	}
	user := domain.User{
		ID:       id,
		Name:     name,
		Email:    email,
		Password: password,
	}
	err = db.Create(&user).Error
	if err != nil {
		return err
	}
	return nil
}

func (e *DBHelperImpl) DeleteUser(id uuid.UUID) error {
	db, err := connectDB()
	if err != nil {
		return err
	}
	err = db.Where("id = ?", id).Delete(&domain.User{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (e *DBHelperImpl) FindUserByEmail(email string) (*domain.User, error) {
	db, err := connectDB()
	if err != nil {
		return nil, err
	}
	var user domain.User
	err = db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
