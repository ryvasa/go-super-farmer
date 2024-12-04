package repository_test

import (
	"testing"

	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/repository"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestUserRepository_Create(t *testing.T) {
	dsn := "host=localhost user=postgres password=123 dbname=go_super_farmer_test port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	repo := repository.NewUserRepository(db)

	t.Run("Test Create user successfully", func(t *testing.T) {
		// Create user
		user := &domain.User{
			ID:       int64(1),
			Name:     "Test User",
			Email:    "test@example1.com",
			Password: "securepassword",
		}

		err := repo.Create(user)
		assert.NoError(t, err)

		// Cleanup: delete the user after the test is complete
		t.Cleanup(func() {
			db.Delete(&domain.User{}, user.ID)
		})
	})

	t.Run("Test Create user with invalid user data", func(t *testing.T) {
		// Create user with invalid data
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
		// Simulate a database error (e.g., by closing the database connection)
		// db.Close() // Uncomment to simulate the error

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
