package repository_test

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	repository_implementation "github.com/ryvasa/go-super-farmer/internal/repository/implementation"
	repository_interface "github.com/ryvasa/go-super-farmer/internal/repository/interface"
	"github.com/ryvasa/go-super-farmer/pkg/database"
	"github.com/stretchr/testify/assert"
)

type SupplyHistoryIDs struct {
	SupplyHistoryID uuid.UUID
	CityID          int64
	CommodityID     uuid.UUID
}

type SupplyHistoryMockRows struct {
	SupplyHistory   *sqlmock.Rows
	SupplyHistories *sqlmock.Rows
}

type SupplyHistoryMocDomain struct {
	SupplyHistory *domain.SupplyHistory
}

func SupplyHistoryRepositorySetup(t *testing.T) (*database.MockDB, repository_interface.SupplyHistoryRepository, SupplyHistoryIDs, SupplyHistoryMockRows, SupplyHistoryMocDomain) {

	mockDB := database.NewMockDB(t)

	repo := repository_implementation.NewSupplyHistoryRepository(mockDB.BaseRepo)

	supplyHistoryID := uuid.New()
	cityID := int64(1)
	commodityID := uuid.New()

	ids := SupplyHistoryIDs{
		SupplyHistoryID: supplyHistoryID,
		CityID:          cityID,
		CommodityID:     commodityID,
	}

	rows := SupplyHistoryMockRows{
		SupplyHistory: sqlmock.NewRows([]string{"id", "city_id", "commodity_id", "quantity", "created_at", "updated_at"}).
			AddRow(supplyHistoryID, cityID, commodityID, float64(10), time.Now(), time.Now()),

		SupplyHistories: sqlmock.NewRows([]string{"id", "city_id", "commodity_id", "quantity", "created_at", "updated_at"}).
			AddRow(supplyHistoryID, cityID, commodityID, float64(10), time.Now(), time.Now()).
			AddRow(uuid.New(), cityID, commodityID, float64(10), time.Now(), time.Now()),
	}

	domains := SupplyHistoryMocDomain{
		SupplyHistory: &domain.SupplyHistory{
			ID:          supplyHistoryID,
			CityID:      cityID,
			CommodityID: commodityID,
			Quantity:    float64(10),
		},
	}

	return mockDB, repo, ids, rows, domains
}

func TestSupplyHistoryRepository_Create(t *testing.T) {
	mockDB, repo, ids, _, domains := SupplyHistoryRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL := `INSERT INTO "supply_histories" ("id","commodity_id","city_id","quantity","unit","created_at","updated_at","deleted_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`

	t.Run("should not return error when create successfully", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.SupplyHistoryID, ids.CommodityID, ids.CityID, float64(10), "kg", sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mockDB.Mock.ExpectCommit()

		err := repo.Create(context.TODO(), domains.SupplyHistory)
		assert.Nil(t, err)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when create failed", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.SupplyHistoryID, ids.CommodityID, ids.CityID, float64(10), "kg", sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
			WillReturnError(errors.New("database error"))
		mockDB.Mock.ExpectRollback()

		err := repo.Create(context.TODO(), domains.SupplyHistory)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}

func TestSupplyHistoryRepository_FindAll(t *testing.T) {
	mockDB, repo, ids, rows, _ := SupplyHistoryRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL := `SELECT * FROM "supply_histories"`

	t.Run("should return supply histories when find all successfully", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WillReturnRows(rows.SupplyHistories)

		result, err := repo.FindAll(context.TODO())
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 2)
		assert.Equal(t, ids.SupplyHistoryID, (result)[0].ID)
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

func TestSupplyHistoryRepository_FindByID(t *testing.T) {
	mockDB, repo, ids, rows, _ := SupplyHistoryRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL := `SELECT * FROM "supply_histories" WHERE "supply_histories"."id" = $1 AND "supply_histories"."deleted_at" IS NULL ORDER BY "supply_histories"."id" LIMIT $2`

	t.Run("should return supply history when find by id successfully", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.SupplyHistoryID, 1).WillReturnRows(rows.SupplyHistory)

		result, err := repo.FindByID(context.TODO(), ids.SupplyHistoryID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.SupplyHistoryID, result.ID)
		assert.Equal(t, ids.CityID, result.CityID)
		assert.Equal(t, ids.CommodityID, result.CommodityID)
		assert.Equal(t, float64(10), result.Quantity)

		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when find by id failed", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.SupplyHistoryID, 1).WillReturnError(errors.New("database error"))

		result, err := repo.FindByID(context.TODO(), ids.SupplyHistoryID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")

		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}

func TestSupplyHistoryRepository_FindByCommodityID(t *testing.T) {
	mockDB, repo, ids, rows, _ := SupplyHistoryRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL := `SELECT * FROM "supply_histories" WHERE commodity_id = $1 AND "supply_histories"."deleted_at" IS NULL`

	t.Run("should return supply histories when find by commodity id successfully", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.CommodityID).WillReturnRows(rows.SupplyHistories)

		result, err := repo.FindByCommodityID(context.TODO(), ids.CommodityID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.SupplyHistoryID, (result)[0].ID)
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

func TestSupplyHistoryRepository_FindByCityID(t *testing.T) {
	mockDB, repo, ids, rows, _ := SupplyHistoryRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL := `SELECT * FROM "supply_histories" WHERE city_id = $1 AND "supply_histories"."deleted_at" IS NULL`

	t.Run("should return supply histories when find by city id successfully", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.CityID).WillReturnRows(rows.SupplyHistories)

		result, err := repo.FindByCityID(context.TODO(), ids.CityID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.SupplyHistoryID, (result)[0].ID)
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

func TestSupplyHistoryRepository_FindByCommodityIDAndCityID(t *testing.T) {
	mockDB, repo, ids, rows, _ := SupplyHistoryRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL := `SELECT * FROM "supply_histories" WHERE (commodity_id = $1 AND city_id = $2) AND "supply_histories"."deleted_at" IS NULL`

	t.Run("should return supply history when find by commodity id and city id successfully", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.CommodityID, ids.CityID).WillReturnRows(rows.SupplyHistory)

		result, err := repo.FindByCommodityIDAndCityID(context.TODO(), ids.CommodityID, ids.CityID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.SupplyHistoryID, (result)[0].ID)
		assert.Equal(t, ids.CityID, (result)[0].CityID)
		assert.Equal(t, ids.CommodityID, (result)[0].CommodityID)
		assert.Equal(t, float64(10), (result)[0].Quantity)

		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when find by commodity id and city id failed", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.CommodityID, ids.CityID).WillReturnError(errors.New("database error"))

		result, err := repo.FindByCommodityIDAndCityID(context.TODO(), ids.CommodityID, ids.CityID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")

		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}
