package repository_test

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/service_api/model/domain"
	repository_implementation "github.com/ryvasa/go-super-farmer/service_api/repository/implementation"
	repository_interface "github.com/ryvasa/go-super-farmer/service_api/repository/interface"
	"github.com/ryvasa/go-super-farmer/pkg/database"
	"github.com/stretchr/testify/assert"
)

type SupplyID struct {
	SupplyID    uuid.UUID
	CityID      int64
	CommodityID uuid.UUID
}

type SupplyMockRows struct {
	Supply   *sqlmock.Rows
	Supplies *sqlmock.Rows
}

type SupplyMocDomain struct {
	Supply *domain.Supply
}

func SupplyRepositorySetup(t *testing.T) (*database.MockDB, repository_interface.SupplyRepository, SupplyID, SupplyMockRows, SupplyMocDomain) {

	mockDB := database.NewMockDB(t)

	repo := repository_implementation.NewSupplyRepository(mockDB.BaseRepo)

	supplyID := uuid.New()
	cityID := int64(1)
	commodityID := uuid.New()

	ids := SupplyID{
		SupplyID:    supplyID,
		CityID:      cityID,
		CommodityID: commodityID,
	}

	rows := SupplyMockRows{
		Supply: sqlmock.NewRows([]string{"id", "city_id", "commodity_id", "quantity", "created_at", "updated_at"}).
			AddRow(supplyID, cityID, commodityID, float64(10), time.Now(), time.Now()),

		Supplies: sqlmock.NewRows([]string{"id", "city_id", "commodity_id", "quantity", "created_at", "updated_at"}).
			AddRow(supplyID, cityID, commodityID, float64(10), time.Now(), time.Now()).
			AddRow(uuid.New(), cityID, commodityID, float64(10), time.Now(), time.Now()),
	}

	domains := SupplyMocDomain{
		Supply: &domain.Supply{
			ID:          supplyID,
			CityID:      cityID,
			CommodityID: commodityID,
			Quantity:    float64(10),
		},
	}

	return mockDB, repo, ids, rows, domains
}

func TestSupplyRepository_Create(t *testing.T) {
	mockDB, repo, ids, _, domains := SupplyRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL := `INSERT INTO "supplies" ("id","commodity_id","city_id","quantity","unit","created_at","updated_at","deleted_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`

	t.Run("should not return error when create successfully", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.SupplyID, ids.CommodityID, ids.CityID, float64(10), "kg", sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mockDB.Mock.ExpectCommit()

		err := repo.Create(context.TODO(), domains.Supply)
		assert.Nil(t, err)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when create failed", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.SupplyID, ids.CommodityID, ids.CityID, float64(10), "kg", sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
			WillReturnError(errors.New("database error"))
		mockDB.Mock.ExpectRollback()

		err := repo.Create(context.TODO(), domains.Supply)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}

func TestSupplyRepository_FindAll(t *testing.T) {
	mockDB, repo, ids, rows, _ := SupplyRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL := `SELECT * FROM "supplies"`

	t.Run("should return supplies when find all successfully", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WillReturnRows(rows.Supplies)

		result, err := repo.FindAll(context.TODO())
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 2)
		assert.Equal(t, ids.SupplyID, (result)[0].ID)
		assert.Equal(t, ids.CityID, (result)[0].CityID)
		assert.Equal(t, ids.CommodityID, (result)[0].CommodityID)
		assert.Equal(t, float64(10), (result)[0].Quantity)

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

func TestSupplyRepository_FindByID(t *testing.T) {
	mockDB, repo, ids, rows, _ := SupplyRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL := `SELECT * FROM "supplies" WHERE "supplies"."id" = $1 AND "supplies"."deleted_at" IS NULL ORDER BY "supplies"."id" LIMIT $2`

	t.Run("should return supply when find by id successfully", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.SupplyID, 1).WillReturnRows(rows.Supply)

		result, err := repo.FindByID(context.TODO(), ids.SupplyID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.SupplyID, result.ID)
		assert.Equal(t, ids.CityID, result.CityID)
		assert.Equal(t, ids.CommodityID, result.CommodityID)
		assert.Equal(t, float64(10), result.Quantity)

		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when find by id failed", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.SupplyID, 1).WillReturnError(errors.New("database error"))

		result, err := repo.FindByID(context.TODO(), ids.SupplyID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")

		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}

func TestSupplyRepository_FindByCommodityID(t *testing.T) {
	mockDB, repo, ids, rows, _ := SupplyRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL := `SELECT * FROM "supplies" WHERE commodity_id = $1 AND "supplies"."deleted_at" IS NULL`

	t.Run("should return supplies when find by commodity id successfully", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.CommodityID).WillReturnRows(rows.Supplies)

		result, err := repo.FindByCommodityID(context.TODO(), ids.CommodityID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.SupplyID, (result)[0].ID)
		assert.Equal(t, ids.CityID, (result)[0].CityID)
		assert.Equal(t, ids.CommodityID, (result)[0].CommodityID)
		assert.Equal(t, float64(10), (result)[0].Quantity)

		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when find by commodity id failed", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.CommodityID).WillReturnError(errors.New("database error"))

		result, err := repo.FindByCommodityID(context.TODO(), ids.CommodityID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")

		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}

func TestSupplyRepository_FindByCityID(t *testing.T) {
	mockDB, repo, ids, rows, _ := SupplyRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL := `SELECT * FROM "supplies" WHERE city_id = $1 AND "supplies"."deleted_at" IS NULL`

	t.Run("should return supplies when find by city id successfully", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.CityID).WillReturnRows(rows.Supplies)

		result, err := repo.FindByCityID(context.TODO(), ids.CityID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.SupplyID, (result)[0].ID)
		assert.Equal(t, ids.CityID, (result)[0].CityID)
		assert.Equal(t, ids.CommodityID, (result)[0].CommodityID)
		assert.Equal(t, float64(10), (result)[0].Quantity)

		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when find by city id failed", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.CityID).WillReturnError(errors.New("database error"))

		result, err := repo.FindByCityID(context.TODO(), ids.CityID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")

		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}

func TestSupplyRepository_Delete(t *testing.T) {
	mockDB, repo, ids, _, _ := SupplyRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL := `UPDATE "supplies" SET "deleted_at"=$1 WHERE id = $2 AND "supplies"."deleted_at" IS NULL`

	t.Run("should not return error when delete successfully", func(t *testing.T) {

		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).WithArgs(sqlmock.AnyArg(), ids.SupplyID).WillReturnResult(sqlmock.NewResult(1, 1))
		mockDB.Mock.ExpectCommit()

		err := repo.Delete(context.TODO(), ids.SupplyID)
		assert.Nil(t, err)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when delete failed", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).WithArgs(sqlmock.AnyArg(), ids.SupplyID).WillReturnError(errors.New("database error"))
		mockDB.Mock.ExpectRollback()

		err := repo.Delete(context.TODO(), ids.SupplyID)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}

func TestSupplyRepository_Update(t *testing.T) {
	mockDB, repo, ids, _, _ := SupplyRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL := `UPDATE "supplies" SET "quantity"=$1,"updated_at"=$2 WHERE id = $3 AND "supplies"."deleted_at" IS NULL`

	t.Run("should not return error when update successfully", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).WithArgs(float64(10), sqlmock.AnyArg(), ids.SupplyID).WillReturnResult(sqlmock.NewResult(1, 1))
		mockDB.Mock.ExpectCommit()

		err := repo.Update(context.TODO(), ids.SupplyID, &domain.Supply{Quantity: 10})
		assert.Nil(t, err)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when update failed", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).WithArgs(float64(10), sqlmock.AnyArg(), ids.SupplyID).WillReturnError(errors.New("database error"))
		mockDB.Mock.ExpectRollback()

		err := repo.Update(context.TODO(), ids.SupplyID, &domain.Supply{Quantity: 10})
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}

func TestSupplyRepository_FindByCommodityIDAndCityID(t *testing.T) {
	mockDB, repo, ids, rows, _ := SupplyRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL := `SELECT * FROM "supplies" WHERE (commodity_id = $1 AND city_id = $2) AND "supplies"."deleted_at" IS NULL ORDER BY "supplies"."id" LIMIT $3`

	t.Run("should return supply when find by commodity id and city id successfully", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.CommodityID, ids.CityID, 1).WillReturnRows(rows.Supply)

		result, err := repo.FindByCommodityIDAndCityID(context.TODO(), ids.CommodityID, ids.CityID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.SupplyID, result.ID)
		assert.Equal(t, ids.CityID, result.CityID)
		assert.Equal(t, ids.CommodityID, result.CommodityID)
		assert.Equal(t, float64(10), result.Quantity)

		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when find by commodity id and city id failed", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.CommodityID, ids.CityID, 1).WillReturnError(errors.New("database error"))

		result, err := repo.FindByCommodityIDAndCityID(context.TODO(), ids.CommodityID, ids.CityID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")

		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}
