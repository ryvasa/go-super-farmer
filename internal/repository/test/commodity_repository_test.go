package repository_test

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/pkg/database"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	repository_implementation "github.com/ryvasa/go-super-farmer/internal/repository/implementation"
	repository_interface "github.com/ryvasa/go-super-farmer/internal/repository/interface"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type CommodityIDs struct {
	CommodityID uuid.UUID
}

type CommodityMockRows struct {
	Commodity   *sqlmock.Rows
	Commodities *sqlmock.Rows
}

type CommodityMocDomain struct {
	Commodity *domain.Commodity
}

type CommodityDTO struct {
	Pagination *dto.PaginationDTO
	Filter     *dto.ParamFilterDTO
}

func CommodityRepositorySetup(t *testing.T) (*database.MockDB, repository_interface.CommodityRepository, CommodityIDs, CommodityMockRows, CommodityMocDomain, CommodityDTO) {

	mockDB := database.NewMockDB(t)

	repo := repository_implementation.NewCommodityRepository(mockDB.DB)

	commodityID := uuid.New()

	ids := CommodityIDs{
		CommodityID: commodityID,
	}

	rows := CommodityMockRows{
		Commodity: sqlmock.NewRows([]string{"id", "name", "description", "created_at", "updated_at"}).
			AddRow(commodityID, "commodity", "commodity description", time.Now(), time.Now()),

		Commodities: sqlmock.NewRows([]string{"id", "name", "description", "created_at", "updated_at"}).
			AddRow(commodityID, "commodity", "commodity description", time.Now(), time.Now()).
			AddRow(uuid.New(), "commodity", "commodity description", time.Now(), time.Now()),
	}

	domains := CommodityMocDomain{
		Commodity: &domain.Commodity{
			ID:          commodityID,
			Name:        "commodity",
			Description: "commodity description",
			Code:        "12345",
		},
	}

	dtos := CommodityDTO{
		Pagination: &dto.PaginationDTO{
			Page:  1,
			Limit: 10,
			Sort:  "created_at desc",
		},
		Filter: &dto.ParamFilterDTO{},
	}

	return mockDB, repo, ids, rows, domains, dtos
}

func TestCommodityRepository_Create(t *testing.T) {
	mockDB, repo, ids, _, domains, _ := CommodityRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL := `INSERT INTO "commodities" ("id","name","description","code","duration","created_at","updated_at","deleted_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`

	t.Run("should not return error when create successfully", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.CommodityID, "commodity", "commodity description", "12345", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mockDB.Mock.ExpectCommit()

		err := repo.Create(context.TODO(), domains.Commodity)
		assert.Nil(t, err)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when create failed", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.CommodityID, "commodity", "commodity description", "12345", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
			WillReturnError(errors.New("database error"))
		mockDB.Mock.ExpectRollback()

		err := repo.Create(context.TODO(), domains.Commodity)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}

func TestCommodityRepository_FindAll(t *testing.T) {
	mockDB, repo, ids, rows, _, dtos := CommodityRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL := `SELECT * FROM "commodities" WHERE "commodities"."deleted_at" IS NULL ORDER BY created_at desc LIMIT $1`

	t.Run("should return commodities when find all successfully", func(t *testing.T) {

		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WillReturnRows(rows.Commodities)

		result, err := repo.FindAll(context.TODO(), dtos.Pagination)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 2)
		assert.Equal(t, ids.CommodityID, (result)[0].ID)
		assert.Equal(t, "commodity", (result)[0].Name)
		assert.Equal(t, "commodity description", (result)[0].Description)

		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when find all failed", func(t *testing.T) {

		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
			WillReturnError(errors.New("database error"))

		result, err := repo.FindAll(context.TODO(), dtos.Pagination)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")

		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}

func TestCommodityRepository_FindByID(t *testing.T) {
	mockDB, repo, ids, rows, _, _ := CommodityRepositorySetup(t)

	defer mockDB.SqlDB.Close()
	expectedSQL := `SELECT * FROM "commodities" WHERE "commodities"."id" = $1 AND "commodities"."deleted_at" IS NULL ORDER BY "commodities"."id" LIMIT $2`
	t.Run("should return commodity when find by id successfully", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.CommodityID, 1).WillReturnRows(rows.Commodity)
		result, err := repo.FindByID(context.TODO(), ids.CommodityID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.CommodityID, result.ID)
		assert.Equal(t, "commodity", result.Name)
		assert.Equal(t, "commodity description", result.Description)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
	t.Run("should return error when find by id failed", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.CommodityID, 1).WillReturnError(errors.New("database error"))
		result, err := repo.FindByID(context.TODO(), ids.CommodityID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
	t.Run("should return error when find by id not found", func(t *testing.T) {
		commodity := sqlmock.NewRows([]string{"id", "name", "description", "created_at", "updated_at", "deleted_at"})
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.CommodityID, 1).WillReturnRows(commodity)
		result, err := repo.FindByID(context.TODO(), ids.CommodityID)
		assert.Nil(t, result)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}

func TestCommodityRepository_Update(t *testing.T) {
	mockDB, repo, ids, _, domains, _ := CommodityRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL := `UPDATE "commodities" SET "id"=$1,"name"=$2,"description"=$3,"code"=$4,"updated_at"=$5 WHERE id = $6 AND "commodities"."deleted_at" IS NULL`

	t.Run("should not return error when update successfully", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.CommodityID, "commodity", "commodity description", "12345", sqlmock.AnyArg(), ids.CommodityID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mockDB.Mock.ExpectCommit()
		err := repo.Update(context.TODO(), ids.CommodityID, domains.Commodity)
		assert.Nil(t, err)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
	t.Run("should return error when update failed", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.CommodityID, "commodity", "commodity description", "12345", sqlmock.AnyArg(), ids.CommodityID).
			WillReturnError(errors.New("database error"))
		mockDB.Mock.ExpectRollback()
		err := repo.Update(context.TODO(), ids.CommodityID, domains.Commodity)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
	t.Run("should return error when update not found", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.CommodityID, "commodity", "commodity description", "12345", sqlmock.AnyArg(), ids.CommodityID).
			WillReturnError(gorm.ErrRecordNotFound)
		mockDB.Mock.ExpectRollback()
		err := repo.Update(context.TODO(), ids.CommodityID, domains.Commodity)
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}

func TestCommodityRepository_Delete(t *testing.T) {
	mockDB, repo, ids, _, _, _ := CommodityRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL := `UPDATE "commodities" SET "deleted_at"=$1 WHERE id = $2 AND "commodities"."deleted_at" IS NULL`

	t.Run("should not return error when delete successfully", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), ids.CommodityID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mockDB.Mock.ExpectCommit()
		err := repo.Delete(context.TODO(), ids.CommodityID)
		assert.Nil(t, err)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
	t.Run("should return error when delete failed", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), ids.CommodityID).
			WillReturnError(errors.New("database error"))
		mockDB.Mock.ExpectRollback()
		err := repo.Delete(context.TODO(), ids.CommodityID)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
	t.Run("should return error when delete not found", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), ids.CommodityID).
			WillReturnError(gorm.ErrRecordNotFound)
		mockDB.Mock.ExpectRollback()
		err := repo.Delete(context.TODO(), ids.CommodityID)
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}

func TestCommodityRepository_Restore(t *testing.T) {
	mockDB, repo, ids, _, _, _ := CommodityRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL := `UPDATE "commodities" SET "deleted_at"=$1,"updated_at"=$2 WHERE id = $3`

	t.Run("should not return error when restore successfully", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), ids.CommodityID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mockDB.Mock.ExpectCommit()
		err := repo.Restore(context.TODO(), ids.CommodityID)
		assert.Nil(t, err)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
	t.Run("should return error when restore failed", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), ids.CommodityID).
			WillReturnError(errors.New("database error"))
		mockDB.Mock.ExpectRollback()
		err := repo.Restore(context.TODO(), ids.CommodityID)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
	t.Run("should return error when restore not found", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), ids.CommodityID).
			WillReturnError(gorm.ErrRecordNotFound)
		mockDB.Mock.ExpectRollback()
		err := repo.Restore(context.TODO(), ids.CommodityID)
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}

func TestCommodityRepository_FindDeletedByID(t *testing.T) {
	mockDB, repo, ids, rows, _, _ := CommodityRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL := `SELECT * FROM "commodities" WHERE id = $1 AND deleted_at IS NOT NULL ORDER BY "commodities"."id" LIMIT $2`

	t.Run("should return commodity when find deleted by id successfully", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.CommodityID, 1).WillReturnRows(rows.Commodity)
		result, err := repo.FindDeletedByID(context.TODO(), ids.CommodityID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.CommodityID, result.ID)
		assert.Equal(t, "commodity", result.Name)
		assert.Equal(t, "commodity description", result.Description)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
	t.Run("should return error when find deleted by id failed", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.CommodityID, 1).WillReturnError(errors.New("database error"))
		result, err := repo.FindDeletedByID(context.TODO(), ids.CommodityID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
	t.Run("should return error when find deleted by id not found", func(t *testing.T) {
		commodity := sqlmock.NewRows([]string{"id", "name", "description", "created_at", "updated_at", "deleted_at"})
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.CommodityID, 1).WillReturnRows(commodity)
		result, err := repo.FindDeletedByID(context.TODO(), ids.CommodityID)
		assert.Nil(t, result)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}

func TestCommodityRepository_Count(t *testing.T) {
	mockDB, repo, _, _, _, dtos := CommodityRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL := `SELECT count(*) FROM "commodities" WHERE "commodities"."deleted_at" IS NULL`

	t.Run("should return count when count successfully", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
		result, err := repo.Count(context.TODO(), dtos.Filter)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(1), result)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
	t.Run("should return error when count failed", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WillReturnError(errors.New("database error"))
		result, err := repo.Count(context.TODO(), dtos.Filter)
		assert.Equal(t, int64(0), result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}
