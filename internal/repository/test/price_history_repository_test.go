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
	"gorm.io/gorm"
)

type PriceHistoryIDs struct {
	PriceHistoryID uuid.UUID
	CommodityID    uuid.UUID
	RegionID       uuid.UUID
}

type PriceHistoryMockRows struct {
	PriceHistory *sqlmock.Rows
	Notfound     *sqlmock.Rows
	Commodity    *sqlmock.Rows
	Region       *sqlmock.Rows
}

type PriceHistoryMocDomain struct {
	PriceHistory *domain.PriceHistory
}

func PriceHistoryRepositorySetup(t *testing.T) (*database.MockDB, repository_interface.PriceHistoryRepository, PriceHistoryIDs, PriceHistoryMockRows, PriceHistoryMocDomain) {

	mockDB := database.NewMockDB(t)

	repo := repository_implementation.NewPriceHistoryRepository(mockDB.BaseRepo)

	priceHistoryID := uuid.New()
	commodityID := uuid.New()
	regionID := uuid.New()

	ids := PriceHistoryIDs{
		PriceHistoryID: priceHistoryID,
		CommodityID:    commodityID,
		RegionID:       regionID,
	}

	rows := PriceHistoryMockRows{
		PriceHistory: sqlmock.NewRows([]string{"id", "commodity_id", "region_id", "price", "created_at", "updated_at"}).
			AddRow(priceHistoryID, commodityID, regionID, float64(100), time.Now(), time.Now()),

		Notfound: sqlmock.NewRows([]string{"id", "commodity_id", "region_id", "price", "created_at", "updated_at"}),

		Commodity: sqlmock.NewRows([]string{"id", "name", "code", "duration"}).
			AddRow(commodityID, "commodity name", "commodity code", time.Now()),
		Region: sqlmock.NewRows([]string{"id", "name"}).
			AddRow(regionID, "region name"),
	}

	domains := PriceHistoryMocDomain{
		PriceHistory: &domain.PriceHistory{
			ID:          priceHistoryID,
			CommodityID: commodityID,
			RegionID:    regionID,
			Price:       float64(100),
		},
	}

	return mockDB, repo, ids, rows, domains
}

func TestPriceHistoryRepository_Create(t *testing.T) {
	mockDB, repo, ids, _, domains := PriceHistoryRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL := `INSERT INTO "price_histories" ("id","commodity_id","region_id","price","unit","created_at","updated_at","deleted_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`

	t.Run("should not return error when create successfully", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.PriceHistoryID, ids.CommodityID, ids.RegionID, float64(100), "idr", sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mockDB.Mock.ExpectCommit()

		err := repo.Create(context.TODO(), domains.PriceHistory)
		assert.Nil(t, err)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when create failed", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.PriceHistoryID, ids.CommodityID, ids.RegionID, float64(100), "idr", sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
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
		assert.Equal(t, ids.RegionID, result.RegionID)
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

func TestPriceHistoryRepository_FindByCommodityIDAndRegionID(t *testing.T) {
	mockDB, repo, ids, rows, _ := PriceHistoryRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL1 := `SELECT * FROM "price_histories" WHERE (commodity_id = $1 AND region_id = $2) AND "price_histories"."deleted_at" IS NULL`

	expectedSQL2 := `SELECT "commodities"."id","commodities"."name","commodities"."code","commodities"."duration" FROM "commodities" WHERE "commodities"."id" = $1`

	expectedSQL3 := `SELECT "regions"."id","regions"."province_id","regions"."city_id" FROM "regions" WHERE "regions"."id" = $1`

	t.Run("should return price history when find by commodity id and region id successfully", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL1)).WithArgs(ids.CommodityID, ids.RegionID).WillReturnRows(rows.PriceHistory)

		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL2)).WithArgs(ids.CommodityID).WillReturnRows(rows.Commodity)

		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL3)).WithArgs(ids.RegionID).WillReturnRows(rows.Region)

		result, err := repo.FindByCommodityIDAndRegionID(context.TODO(), ids.CommodityID, ids.RegionID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 1)
		assert.Equal(t, ids.PriceHistoryID, (result)[0].ID)
		assert.Equal(t, ids.CommodityID, (result)[0].CommodityID)
		assert.Equal(t, ids.RegionID, (result)[0].RegionID)
		assert.Equal(t, float64(100), (result)[0].Price)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when find by commodity id and region id failed", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL1)).WithArgs(ids.CommodityID, ids.RegionID).WillReturnError(errors.New("database error"))

		result, err := repo.FindByCommodityIDAndRegionID(context.TODO(), ids.CommodityID, ids.RegionID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}
