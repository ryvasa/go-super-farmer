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
	"github.com/ryvasa/go-super-farmer/internal/repository"
	"github.com/ryvasa/go-super-farmer/pkg/database"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestPriceHistoryRepository_Create(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewPriceHistoryRepository(db)

	expectedSQL := `INSERT INTO "price_histories" ("id","commodity_id","region_id","price","created_at","updated_at","deleted_at") VALUES ($1,$2,$3,$4,$5,$6,$7)`

	priceID := uuid.New()
	price := float64(100)
	commodityID := uuid.New()
	regionID := uuid.New()

	t.Run("Test Create, successfully", func(t *testing.T) {

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(priceID, commodityID, regionID, price, sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		req := domain.PriceHistory{
			ID:          priceID,
			CommodityID: commodityID,
			RegionID:    regionID,
			Price:       price,
		}
		err := repoImpl.Create(context.TODO(), &req)
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test Create, database error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(priceID, commodityID, regionID, price, sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
			WillReturnError(errors.New("database error"))
		mock.ExpectRollback()

		req := domain.PriceHistory{
			ID:          priceID,
			CommodityID: commodityID,
			RegionID:    regionID,
			Price:       price,
		}
		err := repoImpl.Create(context.TODO(), &req)

		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})

}

func TestPriceHistoryRepository_FindAll(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewPriceHistoryRepository(db)

	expectedSQL := `SELECT * FROM "price_histories" WHERE "price_histories"."deleted_at" IS NULL`

	priceID := uuid.New()
	price := float64(100)
	commodityID := uuid.New()
	regionID := uuid.New()

	t.Run("Test FindAll, error database", func(t *testing.T) {

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WillReturnError(errors.New("database error"))

		result, err := repoImpl.FindAll(context.TODO())
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")

		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test FindAll, successfully", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "commodity_id", "region_id", "price", "created_at", "updated_at"}).
			AddRow(priceID, commodityID, regionID, price, time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WillReturnRows(rows)

		result, err := repoImpl.FindAll(context.TODO())
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Len(t, *result, 1)
		assert.Equal(t, priceID, (*result)[0].ID)
		assert.Equal(t, commodityID, (*result)[0].CommodityID)
		assert.Equal(t, regionID, (*result)[0].RegionID)
		assert.Equal(t, price, (*result)[0].Price)

		assert.Nil(t, mock.ExpectationsWereMet())
	})

}

func TestPriceHistoryRepository_FindByID(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewPriceHistoryRepository(db)

	expectedSQL := `SELECT * FROM "price_histories" WHERE "price_histories"."id" = $1 AND "price_histories"."deleted_at" IS NULL ORDER BY "price_histories"."id" LIMIT $2`

	priceID := uuid.New()
	price := float64(100)
	commodityID := uuid.New()
	regionID := uuid.New()

	t.Run("Test FindByID, error notfound", func(t *testing.T) {
		priceHistory := sqlmock.NewRows([]string{"id", "commodity_id", "region_id", "price", "created_at", "updated_at"})

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(priceID, 1).WillReturnRows(priceHistory)

		result, err := repoImpl.FindByID(context.TODO(), priceID)
		assert.Nil(t, result)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))

		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test FindByID, error database", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(priceID, 1).WillReturnError(errors.New("database error"))

		result, err := repoImpl.FindByID(context.TODO(), priceID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")

		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test FindByID, successfully", func(t *testing.T) {
		priceHistory := sqlmock.NewRows([]string{"id", "commodity_id", "region_id", "price", "created_at", "updated_at"}).
			AddRow(priceID, commodityID, regionID, price, time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(priceID, 1).WillReturnRows(priceHistory)

		result, err := repoImpl.FindByID(context.TODO(), priceID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, priceID, result.ID)
		assert.Equal(t, commodityID, result.CommodityID)
		assert.Equal(t, regionID, result.RegionID)
		assert.Equal(t, price, result.Price)

		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestPriceHistoryRepository_FindByCommodityIDAndRegionID(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewPriceHistoryRepository(db)

	expectedSQL := `SELECT * FROM "price_histories" WHERE (commodity_id = $1 AND region_id = $2) AND "price_histories"."deleted_at" IS NULL`

	priceID := uuid.New()
	price := float64(100)
	commodityID := uuid.New()
	regionID := uuid.New()

	t.Run("Test FindByCommodityIDAndRegionID, error database", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(commodityID, regionID).WillReturnError(errors.New("database error"))

		result, err := repoImpl.FindByCommodityIDAndRegionID(context.TODO(), commodityID, regionID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")

		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test FindByCommodityIDAndRegionID, successfully", func(t *testing.T) {
		priceHistory := sqlmock.NewRows([]string{"id", "commodity_id", "region_id", "price", "created_at", "updated_at"}).
			AddRow(priceID, commodityID, regionID, price, time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(commodityID, regionID).WillReturnRows(priceHistory)

		result, err := repoImpl.FindByCommodityIDAndRegionID(context.TODO(), commodityID, regionID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Len(t, *result, 1)
		assert.Equal(t, priceID, (*result)[0].ID)
		assert.Equal(t, commodityID, (*result)[0].CommodityID)
		assert.Equal(t, regionID, (*result)[0].RegionID)
		assert.Equal(t, price, (*result)[0].Price)

		assert.Nil(t, mock.ExpectationsWereMet())
	})
}
