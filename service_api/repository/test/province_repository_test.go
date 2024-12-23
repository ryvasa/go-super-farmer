package repository_test

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ryvasa/go-super-farmer/service_api/model/domain"
	repository_implementation "github.com/ryvasa/go-super-farmer/service_api/repository/implementation"
	repository_interface "github.com/ryvasa/go-super-farmer/service_api/repository/interface"
	"github.com/ryvasa/go-super-farmer/pkg/database"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type ProvinceIDs struct {
	ProvinceID int64
}

type ProvinceMockRows struct {
	Province *sqlmock.Rows
	Notfound *sqlmock.Rows
}

type ProvinceMockDomains struct {
	Province *domain.Province
}

func ProvinceRepositorySetup(t *testing.T) (*database.MockDB, repository_interface.ProvinceRepository, ProvinceIDs, ProvinceMockRows, ProvinceMockDomains) {

	mockDB := database.NewMockDB(t)

	repo := repository_implementation.NewProvinceRepository(mockDB.DB)

	provinceID := int64(1)

	ids := ProvinceIDs{ProvinceID: provinceID}

	rows := ProvinceMockRows{
		Province: sqlmock.NewRows([]string{
			"id", "name",
		}).AddRow(ids.ProvinceID, "province1"),
		Notfound: sqlmock.NewRows([]string{
			"id", "name",
		}),
	}

	domains := ProvinceMockDomains{
		Province: &domain.Province{
			ID:   ids.ProvinceID,
			Name: "province1",
		},
	}

	return mockDB, repo, ids, rows, domains

}

func TestProvinceRepository_Create(t *testing.T) {
	mockDB, repo, ids, _, domains := ProvinceRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL := `INSERT INTO "provinces" ("name","id") VALUES ($1,$2) RETURNING "id"`

	t.Run("should return no error when create successfully", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
			WithArgs(domains.Province.Name, ids.ProvinceID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))

		mockDB.Mock.ExpectCommit()

		err := repo.Create(context.TODO(), domains.Province)
		assert.Nil(t, err)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when failed create province", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
			WithArgs(domains.Province.Name, ids.ProvinceID).
			WillReturnError(errors.New("database error"))

		mockDB.Mock.ExpectRollback()

		err := repo.Create(context.TODO(), domains.Province)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}

func TestProvinceRepository_FindAll(t *testing.T) {
	mockDB, repo, ids, rows, domains := ProvinceRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL := `SELECT * FROM "provinces"`

	t.Run("should return provinces commodities when find all successfully", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WillReturnRows(rows.Province)

		result, err := repo.FindAll(context.TODO())
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 1, len(result))
		assert.Equal(t, ids.ProvinceID, (result)[0].ID)
		assert.Equal(t, domains.Province.Name, (result)[0].Name)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
	t.Run("should return error when find all failed", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WillReturnError(errors.New("database error"))

		result, err := repo.FindAll(context.TODO())
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}

func TestProvinceRepository_FindByID(t *testing.T) {
	mockDB, repo, ids, rows, domains := ProvinceRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL := `SELECT * FROM "provinces" WHERE "provinces"."id" = $1 ORDER BY "provinces"."id" LIMIT $2`

	t.Run("should return province when find by id successfully", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.ProvinceID, 1).WillReturnRows(rows.Province)

		result, err := repo.FindByID(context.TODO(), ids.ProvinceID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.ProvinceID, result.ID)
		assert.Equal(t, domains.Province.Name, result.Name)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when find by id failed", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.ProvinceID, 1).WillReturnError(errors.New("database error"))

		result, err := repo.FindByID(context.TODO(), ids.ProvinceID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when find by id not found", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.ProvinceID, 1).WillReturnRows(rows.Notfound)

		result, err := repo.FindByID(context.TODO(), ids.ProvinceID)
		assert.Nil(t, result)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))

		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}

func TestProvinceRepository_Update(t *testing.T) {
	mockDB, repo, ids, _, domains := ProvinceRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL := `UPDATE "provinces" SET "id"=$1,"name"=$2 WHERE id = $3`

	t.Run("should return no error when update successfully", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.ProvinceID, domains.Province.Name, ids.ProvinceID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mockDB.Mock.ExpectCommit()

		err := repo.Update(context.TODO(), ids.ProvinceID, domains.Province)
		assert.Nil(t, err)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when failed update province", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.ProvinceID, domains.Province.Name, ids.ProvinceID).
			WillReturnError(errors.New("database error"))

		mockDB.Mock.ExpectRollback()

		err := repo.Update(context.TODO(), ids.ProvinceID, domains.Province)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when update not found", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.ProvinceID, domains.Province.Name, ids.ProvinceID).
			WillReturnError(gorm.ErrRecordNotFound)

		mockDB.Mock.ExpectRollback()

		err := repo.Update(context.TODO(), ids.ProvinceID, domains.Province)
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}

func TestProvinceRepository_Delete(t *testing.T) {
	mockDB, repo, ids, _, _ := ProvinceRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL := `DELETE FROM "provinces" WHERE id = $1`

	t.Run("should return no error when delete successfully", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.ProvinceID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mockDB.Mock.ExpectCommit()

		err := repo.Delete(context.TODO(), ids.ProvinceID)
		assert.Nil(t, err)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when delete failed", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.ProvinceID).
			WillReturnError(errors.New("database error"))

		mockDB.Mock.ExpectRollback()

		err := repo.Delete(context.TODO(), ids.ProvinceID)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when delete not found", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.ProvinceID).
			WillReturnError(gorm.ErrRecordNotFound)

		mockDB.Mock.ExpectRollback()

		err := repo.Delete(context.TODO(), ids.ProvinceID)
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}
