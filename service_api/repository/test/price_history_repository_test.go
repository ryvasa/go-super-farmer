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
	"gorm.io/gorm"
)

type PriceHistoryIDs struct {
	PriceHistoryID uuid.UUID
	CommodityID    uuid.UUID
	CityID         int64
}

type PriceHistoryMockRows struct {
	PriceHistory *sqlmock.Rows
	Notfound     *sqlmock.Rows
	Commodity    *sqlmock.Rows
	City         *sqlmock.Rows
}

type PriceHistoryMocDomain struct {
	PriceHistory *domain.PriceHistory
}

func PriceHistoryRepositorySetup(t *testing.T) (*database.MockDB, repository_interface.PriceHistoryRepository, PriceHistoryIDs, PriceHistoryMockRows, PriceHistoryMocDomain) {

	mockDB := database.NewMockDB(t)

	repo := repository_implementation.NewPriceHistoryRepository(mockDB.BaseRepo)

	priceHistoryID := uuid.New()
	commodityID := uuid.New()
	cityID := int64(1)

	ids := PriceHistoryIDs{
		PriceHistoryID: priceHistoryID,
		CommodityID:    commodityID,
		CityID:         cityID,
	}

	rows := PriceHistoryMockRows{
		PriceHistory: sqlmock.NewRows([]string{"id", "commodity_id", "city_id", "price", "created_at", "updated_at"}).
			AddRow(priceHistoryID, commodityID, cityID, float64(100), time.Now(), time.Now()),

		Notfound: sqlmock.NewRows([]string{"id", "commodity_id", "city_id", "price", "created_at", "updated_at"}),

		Commodity: sqlmock.NewRows([]string{"id", "name", "code", "duration"}).
			AddRow(commodityID, "commodity name", "commodity code", 3000),
		City: sqlmock.NewRows([]string{"id", "name"}).
			AddRow(cityID, "city name"),
	}

	domains := PriceHistoryMocDomain{
		PriceHistory: &domain.PriceHistory{
			ID:          priceHistoryID,
			CommodityID: commodityID,
			CityID:      cityID,
			Price:       float64(100),
		},
	}

	return mockDB, repo, ids, rows, domains
}

func TestPriceHistoryRepository_Create(t *testing.T) {
	mockDB, repo, ids, _, domains := PriceHistoryRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL := `INSERT INTO "price_histories" ("id","commodity_id","city_id","price","unit","created_at","updated_at","deleted_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`

	t.Run("should not return error when create successfully", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.PriceHistoryID, ids.CommodityID, ids.CityID, float64(100), "idr", sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mockDB.Mock.ExpectCommit()

		err := repo.Create(context.TODO(), domains.PriceHistory)
		assert.Nil(t, err)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when create failed", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.PriceHistoryID, ids.CommodityID, ids.CityID, float64(100), "idr", sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
			WillReturnError(errors.New("database error"))
		mockDB.Mock.ExpectRollback()

		err := repo.Create(context.TODO(), domains.PriceHistory)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}

func TestPriceHistoryRepository_FindByID(t *testing.T) {
	mockDB, repo, ids, rows, _ := PriceHistoryRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL := `SELECT * FROM "price_histories" WHERE "price_histories"."id" = $1 AND "price_histories"."deleted_at" IS NULL ORDER BY "price_histories"."id" LIMIT $2`

	t.Run("should return price history when find by id successfully", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.PriceHistoryID, 1).WillReturnRows(rows.PriceHistory)

		result, err := repo.FindByID(context.TODO(), ids.PriceHistoryID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.PriceHistoryID, result.ID)
		assert.Equal(t, ids.CommodityID, result.CommodityID)
		assert.Equal(t, ids.CityID, result.CityID)
		assert.Equal(t, float64(100), result.Price)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
	t.Run("should return error when find by id failed", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.PriceHistoryID, 1).WillReturnError(errors.New("database error"))

		result, err := repo.FindByID(context.TODO(), ids.PriceHistoryID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
	t.Run("should return error when find by id not found", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.PriceHistoryID, 1).WillReturnRows(rows.Notfound)

		result, err := repo.FindByID(context.TODO(), ids.PriceHistoryID)
		assert.Nil(t, result)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}

func TestPriceHistoryRepository_FindByCommodityIDAndCityID(t *testing.T) {
	mockDB, repo, ids, rows, _ := PriceHistoryRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL1 := `SELECT * FROM "price_histories" WHERE (commodity_id = $1 AND city_id = $2) AND "price_histories"."deleted_at" IS NULL`

	expectedSQL2 := `SELECT "commodities"."id","commodities"."name","commodities"."code","commodities"."duration" FROM "commodities" WHERE "commodities"."id" = $1`

	expectedSQL3 := `SELECT * FROM "cities" WHERE "cities"."id" = $1`

	t.Run("should return price history when find by commodity id and city id successfully", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL1)).WithArgs(ids.CommodityID, ids.CityID).WillReturnRows(rows.PriceHistory)

		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL3)).WithArgs(ids.CityID).WillReturnRows(rows.City)
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL2)).WithArgs(ids.CommodityID).WillReturnRows(rows.Commodity)

		result, err := repo.FindByCommodityIDAndCityID(context.TODO(), ids.CommodityID, ids.CityID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 1)
		assert.Equal(t, ids.PriceHistoryID, (result)[0].ID)
		assert.Equal(t, ids.CommodityID, (result)[0].CommodityID)
		assert.Equal(t, ids.CityID, (result)[0].CityID)
		assert.Equal(t, float64(100), (result)[0].Price)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when find by commodity id and city id failed", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL1)).WithArgs(ids.CommodityID, ids.CityID).WillReturnError(errors.New("database error"))

		result, err := repo.FindByCommodityIDAndCityID(context.TODO(), ids.CommodityID, ids.CityID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}
