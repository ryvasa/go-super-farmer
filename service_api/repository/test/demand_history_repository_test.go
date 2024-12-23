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

type DemandHistoryIDs struct {
	DemandHistoryID uuid.UUID
	CityID          int64
	CommodityID     uuid.UUID
}

type DemandHistoryMockRows struct {
	DemandHistory   *sqlmock.Rows
	DemandHistories *sqlmock.Rows
}

type DemandHistoryMocDomain struct {
	DemandHistory *domain.DemandHistory
}

func DemandHistoryRepositorySetup(t *testing.T) (*database.MockDB, DemandHistoryIDs, DemandHistoryMockRows, DemandHistoryMocDomain, repository_interface.DemandHistoryRepository) {

	mockDB := database.NewMockDB(t)

	repo := repository_implementation.NewDemandHistoryRepository(mockDB.BaseRepo)

	cityID := int64(1)
	commodityID := uuid.New()
	demandHistoryID := uuid.New()

	ids := DemandHistoryIDs{
		CityID:          cityID,
		CommodityID:     commodityID,
		DemandHistoryID: demandHistoryID,
	}

	rows := DemandHistoryMockRows{
		DemandHistory: sqlmock.NewRows([]string{"id", "city_id", "commodity_id", "quantity", "created_at", "updated_at"}).
			AddRow(demandHistoryID, cityID, commodityID, float64(10), time.Now(), time.Now()),

		DemandHistories: sqlmock.NewRows([]string{"id", "city_id", "commodity_id", "quantity", "created_at", "updated_at"}).
			AddRow(demandHistoryID, cityID, commodityID, float64(10), time.Now(), time.Now()).
			AddRow(uuid.New(), cityID, commodityID, float64(10), time.Now(), time.Now()),
	}

	domains := DemandHistoryMocDomain{
		DemandHistory: &domain.DemandHistory{
			ID:          demandHistoryID,
			CityID:      cityID,
			CommodityID: commodityID,
			Quantity:    float64(10),
		},
	}

	return mockDB, ids, rows, domains, repo
}

func TestDemandHistoryRepository_Create(t *testing.T) {
	mockDB, ids, _, domains, repo := DemandHistoryRepositorySetup(t)
	defer mockDB.SqlDB.Close()

	expectedSQL := `INSERT INTO "demand_histories" ("id","commodity_id","city_id","quantity","unit","created_at","updated_at","deleted_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`

	t.Run("should not return error when create successfully", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.DemandHistoryID, ids.CommodityID, ids.CityID, float64(10), "kg", sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mockDB.Mock.ExpectCommit()
		err := mockDB.TxManager.WithTransaction(context.Background(), func(ctx context.Context) error {
			return repo.Create(ctx, domains.DemandHistory)
		})
		assert.Nil(t, err)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when create failed", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.DemandHistoryID, ids.CommodityID, ids.CityID, float64(10), "kg", sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
			WillReturnError(errors.New("database error"))
		mockDB.Mock.ExpectRollback()

		err := repo.Create(context.TODO(), domains.DemandHistory)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}

func TestDemandHistoryRepository_FindAll(t *testing.T) {
	mockDB, ids, rows, _, repo := DemandHistoryRepositorySetup(t)
	defer mockDB.SqlDB.Close()

	expectedSQL := `SELECT * FROM "demand_histories"`

	t.Run("should return demand histories when find all successfully", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WillReturnRows(rows.DemandHistories)

		result, err := repo.FindAll(context.TODO())
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 2)
		assert.Equal(t, ids.DemandHistoryID, (result)[0].ID)
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

func TestDemandHistoryRepository_FindByID(t *testing.T) {
	mockDB, ids, rows, _, repo := DemandHistoryRepositorySetup(t)
	defer mockDB.SqlDB.Close()

	expectedSQL := `SELECT * FROM "demand_histories" WHERE "demand_histories"."id" = $1 AND "demand_histories"."deleted_at" IS NULL ORDER BY "demand_histories"."id" LIMIT $2`

	t.Run("should return demand history when find by id successfully", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.DemandHistoryID, 1).WillReturnRows(rows.DemandHistory)

		result, err := repo.FindByID(context.TODO(), ids.DemandHistoryID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.DemandHistoryID, result.ID)
		assert.Equal(t, ids.CityID, result.CityID)
		assert.Equal(t, ids.CommodityID, result.CommodityID)
		assert.Equal(t, float64(10), result.Quantity)

		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when find by id failed", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.DemandHistoryID, 1).WillReturnError(errors.New("database error"))

		result, err := repo.FindByID(context.TODO(), ids.DemandHistoryID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")

		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}

func TestDemandHistoryRepository_FindByCommodityID(t *testing.T) {
	mockDB, ids, rows, _, repo := DemandHistoryRepositorySetup(t)
	defer mockDB.SqlDB.Close()

	expectedSQL := `SELECT * FROM "demand_histories" WHERE commodity_id = $1 AND "demand_histories"."deleted_at" IS NULL`

	t.Run("should return demand histories when find by commodity id successfully", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.CommodityID).WillReturnRows(rows.DemandHistories)

		result, err := repo.FindByCommodityID(context.TODO(), ids.CommodityID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.DemandHistoryID, (result)[0].ID)
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

func TestDemandHistoryRepository_FindByCityID(t *testing.T) {
	mockDB, ids, rows, _, repo := DemandHistoryRepositorySetup(t)
	defer mockDB.SqlDB.Close()

	expectedSQL := `SELECT * FROM "demand_histories" WHERE city_id = $1 AND "demand_histories"."deleted_at" IS NULL`

	t.Run("should return demand histories when find by city id successfully", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.CityID).WillReturnRows(rows.DemandHistories)

		result, err := repo.FindByCityID(context.TODO(), ids.CityID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.DemandHistoryID, (result)[0].ID)
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

func TestDemandHistoryRepository_FindByCommodityIDAndCityID(t *testing.T) {
	mockDB, ids, rows, _, repo := DemandHistoryRepositorySetup(t)
	defer mockDB.SqlDB.Close()

	expectedSQL := `SELECT * FROM "demand_histories" WHERE (commodity_id = $1 AND city_id = $2) AND "demand_histories"."deleted_at" IS NULL`

	t.Run("should return demands history when find by commodity id and city id successfully", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.CommodityID, ids.CityID).WillReturnRows(rows.DemandHistory)

		result, err := repo.FindByCommodityIDAndCityID(context.TODO(), ids.CommodityID, ids.CityID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.DemandHistoryID, (result)[0].ID)
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
