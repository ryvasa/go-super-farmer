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

func TestProvinceRepository_Create(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewProvinceRepository(db)

	expectedSQL := `INSERT INTO "provinces" ("name") VALUES ($1) RETURNING "id"`
	t.Run("Test Create, successfully", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
			WithArgs("Tasikmalaya").
			WillReturnRows(sqlmock.NewRows([]string{"id"}))
		mock.ExpectCommit()

		err := repoImpl.Create(context.TODO(), &domain.Province{Name: "Tasikmalaya"})

		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test Create, error database", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
			WithArgs("Tasikmalaya").
			WillReturnError(errors.New("database error"))
		mock.ExpectRollback()

		err := repoImpl.Create(context.TODO(), &domain.Province{Name: "Tasikmalaya"})

		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestProvinceRepository_FindAll(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewProvinceRepository(db)

	expectedSQL := `SELECT * FROM "provinces"`
	t.Run("Test FindAll, error database", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WillReturnError(errors.New("database error"))

		result, err := repoImpl.FindAll(context.TODO())
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")

		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test FindAll, successfully", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name"}).
			AddRow(1, "Tasikmalaya").
			AddRow(2, "Bangkok")

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WillReturnRows(rows)

		result, err := repoImpl.FindAll(context.TODO())
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Len(t, *result, 2)
		assert.Equal(t, "Tasikmalaya", (*result)[0].Name)
		assert.Equal(t, "Bangkok", (*result)[1].Name)

		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestProvinceRepository_FindByID(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewProvinceRepository(db)

	expectedSQL := `SELECT * FROM "provinces" WHERE "provinces"."id" = $1 ORDER BY "provinces"."id" LIMIT $2`
	t.Run("Test FindByID, error notfound", func(t *testing.T) {
		province := sqlmock.NewRows([]string{"id", "name"})

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(2, 1).WillReturnRows(province)

		result, err := repoImpl.FindByID(context.TODO(), 2)
		assert.Nil(t, result)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))

		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test FindByID, error database", func(t *testing.T) {

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(1, 1).WillReturnError(errors.New("database error"))

		result, err := repoImpl.FindByID(context.TODO(), 1)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")

		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test FindByID, successfully", func(t *testing.T) {
		province := sqlmock.NewRows([]string{"id", "name"}).
			AddRow(1, "Tasikmalaya")

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(1, 1).WillReturnRows(province)

		result, err := repoImpl.FindByID(context.TODO(), 1)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Tasikmalaya", result.Name)

		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestProvinceRepository_Delete(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewProvinceRepository(db)

	expectedSQL := `DELETE FROM "provinces" WHERE id = $1`

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

func TestProvinceRepository_Update(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewProvinceRepository(db)

	expectedSQL := `UPDATE "provinces" SET "name"=$1 WHERE id = $2`
	t.Run("Test Update, error notfound", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).WithArgs("Tasikmalaya", 1).WillReturnError(gorm.ErrRecordNotFound)
		mock.ExpectRollback()

		err := repoImpl.Update(context.TODO(), 1, &domain.Province{Name: "Tasikmalaya"})
		assert.NotNil(t, err)
		assert.EqualError(t, err, gorm.ErrRecordNotFound.Error())

		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test Update, error database", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).WithArgs("Tasikmalaya", 1).WillReturnError(errors.New("database error"))
		mock.ExpectRollback()

		err := repoImpl.Update(context.TODO(), 1, &domain.Province{Name: "Tasikmalaya"})
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")

		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test Update, successfully", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).WithArgs("Tasikmalaya", 1).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repoImpl.Update(context.TODO(), 1, &domain.Province{Name: "Tasikmalaya"})
		assert.Nil(t, err)

		assert.Nil(t, mock.ExpectationsWereMet())
	})
}
