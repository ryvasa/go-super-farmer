package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	repository_implementation "github.com/ryvasa/go-super-farmer/internal/repository/implementation"
	repository_interface "github.com/ryvasa/go-super-farmer/internal/repository/interface"
	"github.com/ryvasa/go-super-farmer/pkg/database"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type CityIDs struct {
	CityID     int64
	ProvinceID int64
}

type CityMockRows struct {
	City     *sqlmock.Rows
	Notfound *sqlmock.Rows
}

type CityMockDomains struct {
	City *domain.City
}

func CityRepositorySetup(t *testing.T) (*sql.DB, sqlmock.Sqlmock, repository_interface.CityRepository, CityIDs, CityMockRows, CityMockDomains) {

	sqlDB, db, mock := database.DbMock(t)

	repo := repository_implementation.NewCityRepository(db)

	cityID := int64(1)
	provinceID := int64(2)

	ids := CityIDs{
		CityID:     cityID,
		ProvinceID: provinceID,
	}

	rows := CityMockRows{
		City: sqlmock.NewRows([]string{
			"id", "name",
		}).AddRow(ids.CityID, "city1"),
		Notfound: sqlmock.NewRows([]string{
			"id", "name",
		}),
	}

	domains := CityMockDomains{
		City: &domain.City{
			ID:         ids.CityID,
			Name:       "city1",
			ProvinceID: provinceID,
		},
	}

	return sqlDB, mock, repo, ids, rows, domains
}

func TestCityRepository_Create(t *testing.T) {
	db, mock, repo, ids, _, domains := CityRepositorySetup(t)

	defer db.Close()

	expectedSQL := `INSERT INTO "cities" ("name","province_id","id") VALUES ($1,$2,$3) RETURNING "id"`

	t.Run("should return no error when create successfully", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
			WithArgs(domains.City.Name, ids.ProvinceID, ids.CityID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))

		mock.ExpectCommit()

		err := repo.Create(context.TODO(), domains.City)
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when failed create city", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
			WithArgs(domains.City.Name, ids.ProvinceID, ids.CityID).
			WillReturnError(errors.New("database error"))

		mock.ExpectRollback()

		err := repo.Create(context.TODO(), domains.City)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestCityRepository_FindAll(t *testing.T) {
	db, mock, repo, ids, rows, domains := CityRepositorySetup(t)

	defer db.Close()

	expectedSQL := `SELECT * FROM "cities"`

	t.Run("should return cities commodities when find all successfully", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WillReturnRows(rows.City)

		result, err := repo.FindAll(context.TODO())
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 1, len(*result))
		assert.Equal(t, ids.CityID, (*result)[0].ID)
		assert.Equal(t, domains.City.Name, (*result)[0].Name)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
	t.Run("should return error when find all failed", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WillReturnError(errors.New("database error"))

		result, err := repo.FindAll(context.TODO())
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestCityRepository_FindByID(t *testing.T) {
	db, mock, repo, ids, rows, domains := CityRepositorySetup(t)

	defer db.Close()

	expectedSQL := `SELECT * FROM "cities" WHERE "cities"."id" = $1 ORDER BY "cities"."id" LIMIT $2`

	t.Run("should return city when find by id successfully", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.CityID, 1).
			WillReturnRows(rows.City)

		result, err := repo.FindByID(context.TODO(), ids.CityID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.CityID, result.ID)
		assert.Equal(t, domains.City.Name, result.Name)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when find by id failed", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.CityID, 1).WillReturnError(errors.New("database error"))

		result, err := repo.FindByID(context.TODO(), ids.CityID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when find by id not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.CityID, 1).WillReturnRows(rows.Notfound)

		result, err := repo.FindByID(context.TODO(), ids.CityID)
		assert.Nil(t, result)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))

		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestCityRepository_Update(t *testing.T) {
	db, mock, repo, ids, _, domains := CityRepositorySetup(t)

	defer db.Close()

	expectedSQL := `UPDATE "cities" SET "id"=$1,"name"=$2,"province_id"=$3 WHERE id = $4`

	t.Run("should return no error when update successfully", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.CityID, domains.City.Name,
				ids.ProvinceID, ids.CityID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		err := repo.Update(context.TODO(), ids.CityID, domains.City)
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when failed update province", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.CityID, domains.City.Name,
				ids.ProvinceID, ids.CityID).
			WillReturnError(errors.New("database error"))

		mock.ExpectRollback()

		err := repo.Update(context.TODO(), ids.CityID, domains.City)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when update not found", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.CityID, domains.City.Name,
				ids.ProvinceID, ids.CityID).
			WillReturnError(gorm.ErrRecordNotFound)

		mock.ExpectRollback()

		err := repo.Update(context.TODO(), ids.CityID, domains.City)
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestCityRepository_Delete(t *testing.T) {
	db, mock, repo, ids, _, _ := CityRepositorySetup(t)

	defer db.Close()

	expectedSQL := `DELETE FROM "cities" WHERE id = $1`

	t.Run("should return no error when delete successfully", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.CityID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		err := repo.Delete(context.TODO(), ids.CityID)
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when failed delete", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.CityID).
			WillReturnError(errors.New("database error"))

		mock.ExpectRollback()

		err := repo.Delete(context.TODO(), ids.CityID)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when delete not found", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.CityID).
			WillReturnError(gorm.ErrRecordNotFound)

		mock.ExpectRollback()

		err := repo.Delete(context.TODO(), ids.CityID)
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}
