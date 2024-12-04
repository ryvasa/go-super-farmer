package repository_test

import (
	"testing"

	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	repository "github.com/ryvasa/go-super-farmer/internal/repository/user"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestUserRepository_Create(t *testing.T) {
	dsn := "host=localhost user=postgres password=123 dbname=go_super_farmer port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	repo := repository.NewUserRepository(db)

	t.Run("Test Create user successfully", func(t *testing.T) {
		user := &domain.User{
			ID:       int64(1),
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "securepassword",
		}

		err := repo.Create(user)
		assert.NoError(t, err)
	})

	t.Run("Test Create user with invalid user data", func(t *testing.T) {
		user := &domain.User{
			ID:       int64(2),
			Name:     "", // invalid name
			Email:    "test@example.com",
			Password: "securepassword",
		}

		err := repo.Create(user)
		assert.Error(t, err)
	})

	t.Run("Test Create user with database error", func(t *testing.T) {
		// db.Close() // close the database connection to simulate an error

		user := &domain.User{
			ID:       int64(3),
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "securepassword",
		}

		err := repo.Create(user)
		assert.Error(t, err)
	})
}
