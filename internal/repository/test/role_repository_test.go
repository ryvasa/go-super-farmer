package repository_test

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/repository"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestRoleRepository_Create(t *testing.T) {
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewRoleRepository(db)

	t.Run("Test Create, successfully", func(t *testing.T) {
		expectedSQL := `INSERT INTO "roles" ("name") VALUES ($1) RETURNING "id"`

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
			WithArgs("admin").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		err := repoImpl.Create(context.TODO(), &domain.Role{Name: "admin"})
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test Create, error database", func(t *testing.T) {
		expectedSQL := `INSERT INTO "roles" ("name") VALUES ($1) RETURNING "id"`

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
			WithArgs("admin").
			WillReturnError(errors.New("database error"))
		mock.ExpectRollback()

		err := repoImpl.Create(context.TODO(), &domain.Role{Name: "admin"})
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestRoleRepository_FindAll(t *testing.T) {
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewRoleRepository(db)
	t.Run("Test FindAll, error database", func(t *testing.T) {
		expectedSQL := `SELECT * FROM "roles"`

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WillReturnError(errors.New("database error"))

		result, err := repoImpl.FindAll(context.TODO())
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")

		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test FindAll, successfully", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name"}).
			AddRow(1, "admin").
			AddRow(2, "user")

		expectedSQL := `SELECT * FROM "roles"`
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WillReturnRows(rows)

		result, err := repoImpl.FindAll(context.TODO())
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Len(t, *result, 2)
		assert.Equal(t, "admin", (*result)[0].Name)
		assert.Equal(t, "user", (*result)[1].Name)

		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestRoleRepository_FindByID(t *testing.T) {
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewRoleRepository(db)

	t.Run("Test FindByID, error notfound", func(t *testing.T) {
		role := sqlmock.NewRows([]string{"id", "name"})

		expectedSQL := `SELECT * FROM "roles" WHERE "roles"."id" = $1 ORDER BY "roles"."id" LIMIT $2`

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(2, 1).WillReturnRows(role)

		result, err := repoImpl.FindByID(context.TODO(), 2)
		assert.Nil(t, result)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))

		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test FindByID, error database", func(t *testing.T) {
		expectedSQL := `SELECT * FROM "roles" WHERE "roles"."id" = $1 ORDER BY "roles"."id" LIMIT $2`

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(1, 1).WillReturnError(errors.New("database error"))

		result, err := repoImpl.FindByID(context.TODO(), 1)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")

		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test FindByID, successfully", func(t *testing.T) {
		role := sqlmock.NewRows([]string{"id", "name"}).
			AddRow(1, "admin")

		expectedSQL := `SELECT * FROM "roles" WHERE "roles"."id" = $1 ORDER BY "roles"."id" LIMIT $2`
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(1, 1).WillReturnRows(role)

		result, err := repoImpl.FindByID(context.TODO(), 1)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "admin", result.Name)

		assert.Nil(t, mock.ExpectationsWereMet())
	})
}
