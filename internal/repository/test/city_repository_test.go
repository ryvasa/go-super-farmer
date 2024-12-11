package repository_test

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/repository"
	"github.com/ryvasa/go-super-farmer/pkg/database"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCityRepository_Create(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewCityRepository(db)

	expectedSQL := `INSERT INTO "cities" ("name","province_id") VALUES ($1,$2) RETURNING "id"`
	t.Run("Test Create, successfully", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs("Tasikmalaya", 1).WillReturnRows(sqlmock.NewRows([]string{"id"}))
		mock.ExpectCommit()

		err := repoImpl.Create(context.TODO(), &domain.City{ProvinceID: 1, Name: "Tasikmalaya"})

		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test Create, error database", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs("Tasikmalaya", 1).WillReturnError(errors.New("database error"))
		mock.ExpectRollback()

		err := repoImpl.Create(context.TODO(), &domain.City{ProvinceID: 1, Name: "Tasikmalaya"})

		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestCityRepository_FindByID(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewCityRepository(db)

	expectedSQL := `SELECT * FROM "cities" WHERE "cities"."id" = $1 ORDER BY "cities"."id" LIMIT $2`
	t.Run("Test FindByID, error notfound", func(t *testing.T) {
		city := sqlmock.NewRows([]string{"id", "province_id", "name"})

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(2, 1).WillReturnRows(city)

		result, err := repoImpl.FindByID(context.TODO(), int64(2))
		assert.Nil(t, result)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))

		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test FindByID, error database", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(1, 1).WillReturnError(errors.New("database error"))

		result, err := repoImpl.FindByID(context.TODO(), int64(1))
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")

		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test FindByID, successfully", func(t *testing.T) {
		city := sqlmock.NewRows([]string{"id", "province_id", "name"}).
			AddRow(1, 1, "Tasikmalaya")

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(1, 1).WillReturnRows(city)

		result, err := repoImpl.FindByID(context.TODO(), int64(1))
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(1), result.ID)
		assert.Equal(t, int64(1), result.ProvinceID)
		assert.Equal(t, "Tasikmalaya", result.Name)

		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestCityRepository_FindAll(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewCityRepository(db)

	expectedSQL := `SELECT * FROM "cities"`
	t.Run("Test FindAll, error database", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WillReturnError(errors.New("database error"))

		result, err := repoImpl.FindAll(context.TODO())
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")

		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test FindAll, successfully", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "province_id", "name"}).
			AddRow(1, 1, "Tasikmalaya").
			AddRow(2, 1, "Bangkok")

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WillReturnRows(rows)

		result, err := repoImpl.FindAll(context.TODO())
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Len(t, *result, 2)
		assert.Equal(t, int64(1), (*result)[0].ID)
		assert.Equal(t, "Tasikmalaya", (*result)[0].Name)
		assert.Equal(t, int64(2), (*result)[1].ID)
		assert.Equal(t, "Bangkok", (*result)[1].Name)

		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestCityRepository_Update(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewCityRepository(db)

	expectedSQL := `UPDATE "cities" SET "name"=$1 WHERE id = $2`
	t.Run("Test Update, error notfound", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).WithArgs("Tasikmalaya", 1).WillReturnError(gorm.ErrRecordNotFound)
		mock.ExpectRollback()

		err := repoImpl.Update(context.TODO(), 1, &domain.City{Name: "Tasikmalaya"})
		assert.NotNil(t, err)
		assert.EqualError(t, err, gorm.ErrRecordNotFound.Error())

		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test Update, error database", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).WithArgs("Tasikmalaya", 1).WillReturnError(errors.New("database error"))
		mock.ExpectRollback()

		err := repoImpl.Update(context.TODO(), 1, &domain.City{Name: "Tasikmalaya"})
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")

		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test Update, successfully", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).WithArgs("Tasikmalaya", 1).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repoImpl.Update(context.TODO(), 1, &domain.City{Name: "Tasikmalaya"})
		assert.Nil(t, err)

		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestCityRepository_Delete(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewCityRepository(db)

	expectedSQL := `DELETE FROM "cities" WHERE id = $1`

	t.Run("Test Delete, error notfound", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).WithArgs(1).WillReturnError(gorm.ErrRecordNotFound)
		mock.ExpectRollback()

		err := repoImpl.Delete(context.TODO(), 1)
		assert.NotNil(t, err)
		assert.EqualError(t, err, gorm.ErrRecordNotFound.Error())

		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test Delete, error database", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).WithArgs(1).WillReturnError(errors.New("database error"))
		mock.ExpectRollback()

		err := repoImpl.Delete(context.TODO(), 1)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")

		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test Delete, successfully", func(t *testing.T) {

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repoImpl.Delete(context.TODO(), 1)
		assert.Nil(t, err)

		assert.Nil(t, mock.ExpectationsWereMet())
	})
}
