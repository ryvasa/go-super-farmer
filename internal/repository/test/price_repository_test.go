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

func TestPriceRepository_Create(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewPriceRepository(db)

	priceID := uuid.New()
	price := float64(100)
	commodityID := uuid.New()
	regionID := uuid.New()

	expectedSQL := `INSERT INTO "prices" ("id","commodity_id","region_id","price","created_at","updated_at","deleted_at") VALUES ($1,$2,$3,$4,$5,$6,$7)`

	t.Run("Test Create, successfully", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(priceID, commodityID, regionID, price, sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		req := domain.Price{
			ID:          priceID,
			CommodityID: commodityID,
			RegionID:    regionID,
			Price:       price,
		}

		err := repoImpl.Create(context.TODO(), &req)
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test Create, error database", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(priceID, commodityID, regionID, price, sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
			WillReturnError(errors.New("database error"))
		mock.ExpectRollback()

		req := domain.Price{
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

func TestPriceRepository_FindAll(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewPriceRepository(db)

	priceID := uuid.New()
	price := float64(100)

	expectedSQL := `SELECT * FROM "prices" WHERE "prices"."deleted_at" IS NULL`

	t.Run("Test FindAll, error database", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WillReturnError(errors.New("database error"))

		result, err := repoImpl.FindAll(context.TODO())
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")

		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test FindAll, successfully", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "price", "created_at", "updated_at", "city_name", "province_name"}).
			AddRow(priceID, price, time.Now(), time.Now(), "city", "province")

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WillReturnRows(rows)

		result, err := repoImpl.FindAll(context.TODO())

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Len(t, *result, 1)
		assert.Equal(t, priceID, (*result)[0].ID)
		assert.Equal(t, price, (*result)[0].Price)

		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestPriceRepository_FindByID(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()
	repoImpl := repository.NewPriceRepository(db)
	priceID := uuid.New()
	expectedSQL := `SELECT * FROM "prices" WHERE "prices"."id" = $1 AND "prices"."deleted_at" IS NULL ORDER BY "prices"."id" LIMIT $2`
	t.Run("Test FindByID, successfully", func(t *testing.T) {
		price := sqlmock.NewRows([]string{"id", "price", "created_at", "updated_at"}).
			AddRow(priceID, float64(100), time.Now(), time.Now())
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(priceID, 1).WillReturnRows(price)
		result, err := repoImpl.FindByID(context.TODO(), priceID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, priceID, result.ID)
		assert.Equal(t, float64(100), result.Price)
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
	t.Run("Test FindByID, not found", func(t *testing.T) {
		price := sqlmock.NewRows([]string{"id", "price", "created_at", "updated_at"})
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(priceID, 1).WillReturnRows(price)
		result, err := repoImpl.FindByID(context.TODO(), priceID)
		assert.Nil(t, result)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestPriceRepository_FindByCommodityID(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()
	repoImpl := repository.NewPriceRepository(db)

	commodityID := uuid.New()

	expectedSQL := `SELECT * FROM "prices" WHERE prices.commodity_id = $1 AND "prices"."deleted_at" IS NULL`

	priceID := uuid.New()

	t.Run("Test FindByCommodityID, successfully", func(t *testing.T) {
		price := sqlmock.NewRows([]string{"id", "price", "created_at", "updated_at"}).
			AddRow(priceID, float64(100), time.Now(), time.Now())
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(commodityID).WillReturnRows(price)
		result, err := repoImpl.FindByCommodityID(context.TODO(), commodityID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, priceID, (*result)[0].ID)
		assert.Equal(t, float64(100), (*result)[0].Price)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
	t.Run("Test FindByCommodityID, error database", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(commodityID).WillReturnError(errors.New("database error"))
		result, err := repoImpl.FindByCommodityID(context.TODO(), commodityID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestPriceRepository_FindByRegionID(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()
	repoImpl := repository.NewPriceRepository(db)
	regionID := uuid.New()

	expectedSQL := `SELECT * FROM "prices" WHERE prices.region_id = $1 AND "prices"."deleted_at" IS NULL`

	priceID := uuid.New()
	t.Run("Test FindByRegionID, successfully", func(t *testing.T) {
		price := sqlmock.NewRows([]string{"id", "price", "created_at", "updated_at"}).
			AddRow(priceID, float64(100), time.Now(), time.Now())
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(regionID).WillReturnRows(price)
		result, err := repoImpl.FindByRegionID(context.TODO(), regionID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, priceID, (*result)[0].ID)
		assert.Equal(t, float64(100), (*result)[0].Price)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
	t.Run("Test FindByRegionID, error database", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(regionID).WillReturnError(errors.New("database error"))
		result, err := repoImpl.FindByRegionID(context.TODO(), regionID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestPriceRepository_Update(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()
	repoImpl := repository.NewPriceRepository(db)
	priceID := uuid.New()
	price := float64(100)
	commodityID := uuid.New()
	regionID := uuid.New()
	expectedSQL := `UPDATE "prices" SET "id"=$1,"commodity_id"=$2,"region_id"=$3,"price"=$4,"updated_at"=$5 WHERE id = $6 AND "prices"."deleted_at" IS NULL`
	t.Run("Test Update, successfully", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(priceID, commodityID, regionID, price, sqlmock.AnyArg(), priceID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
		req := domain.Price{
			ID:          priceID,
			CommodityID: commodityID,
			RegionID:    regionID,
			Price:       price,
		}
		err := repoImpl.Update(context.TODO(), priceID, &req)
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
	t.Run("Test Update, error database", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(priceID, commodityID, regionID, price, sqlmock.AnyArg(), priceID).
			WillReturnError(errors.New("database error"))
		mock.ExpectRollback()
		req := domain.Price{
			ID:          priceID,
			CommodityID: commodityID,
			RegionID:    regionID,
			Price:       price,
		}
		err := repoImpl.Update(context.TODO(), priceID, &req)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})
	t.Run("Test Update, not found", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(priceID, commodityID, regionID, price, sqlmock.AnyArg(), priceID).
			WillReturnError(gorm.ErrRecordNotFound)
		mock.ExpectRollback()
		req := domain.Price{
			ID:          priceID,
			CommodityID: commodityID,
			RegionID:    regionID,
			Price:       price,
		}
		err := repoImpl.Update(context.TODO(), priceID, &req)
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestPriceRepository_Delete(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()
	repoImpl := repository.NewPriceRepository(db)
	priceID := uuid.New()
	expectedSQL := `UPDATE "prices" SET "deleted_at"=$1 WHERE id = $2 AND "prices"."deleted_at" IS NULL`
	t.Run("Test Delete, successfully", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), priceID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
		err := repoImpl.Delete(context.TODO(), priceID)
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
	t.Run("Test Delete, error database", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), priceID).
			WillReturnError(errors.New("database error"))
		mock.ExpectRollback()
		err := repoImpl.Delete(context.TODO(), priceID)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})
	t.Run("Test Delete, not found", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), priceID).
			WillReturnError(gorm.ErrRecordNotFound)
		mock.ExpectRollback()
		err := repoImpl.Delete(context.TODO(), priceID)
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestPriceRepository_Restore(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()
	repoImpl := repository.NewPriceRepository(db)
	priceID := uuid.New()
	expectedSQL := `UPDATE "prices" SET "deleted_at"=$1,"updated_at"=$2 WHERE id = $3`
	t.Run("Test Restore, successfully", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), priceID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
		err := repoImpl.Restore(context.TODO(), priceID)
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
	t.Run("Test Restore, error database", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), priceID).
			WillReturnError(errors.New("database error"))
		mock.ExpectRollback()
		err := repoImpl.Restore(context.TODO(), priceID)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})
	t.Run("Test Restore, not found", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), priceID).
			WillReturnError(gorm.ErrRecordNotFound)
		mock.ExpectRollback()
		err := repoImpl.Restore(context.TODO(), priceID)
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestPriceRepository_FindDeletedByID(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()
	repoImpl := repository.NewPriceRepository(db)
	priceID := uuid.New()
	expectedSQL := `SELECT * FROM "prices" WHERE prices.id = $1 AND prices.deleted_at IS NOT NULL ORDER BY "prices"."id" LIMIT $2`

	t.Run("Test FindDeletedByID, successfully", func(t *testing.T) {
		price := sqlmock.NewRows([]string{"id", "price", "created_at", "updated_at"}).
			AddRow(priceID, float64(100), time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(priceID, 1).WillReturnRows(price)
		result, err := repoImpl.FindDeletedByID(context.TODO(), priceID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, priceID, result.ID)
		assert.Equal(t, float64(100), result.Price)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
	t.Run("Test FindDeletedByID, error database", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(priceID, 1).WillReturnError(errors.New("database error"))
		result, err := repoImpl.FindDeletedByID(context.TODO(), priceID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})
	t.Run("Test FindDeletedByID, not found", func(t *testing.T) {
		price := sqlmock.NewRows([]string{"id", "price", "created_at", "updated_at", "city_name", "province_name"})
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(priceID, 1).WillReturnRows(price)
		result, err := repoImpl.FindDeletedByID(context.TODO(), priceID)
		assert.Nil(t, result)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestPriceRepositoryImpl_FindByCommodityIDAndRegionID(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()
	repoImpl := repository.NewPriceRepository(db)
	priceID := uuid.New()
	regionID := uuid.New()
	commodityID := uuid.New()

	expectedSQL := `SELECT * FROM "prices" WHERE (prices.commodity_id = $1 AND prices.region_id = $2) AND "prices"."deleted_at" IS NULL ORDER BY "prices"."id" LIMIT $3`
	t.Run("Test FindByCommodityIDAndRegionID, successfully", func(t *testing.T) {
		price := sqlmock.NewRows([]string{"id", "price",
			"region_id", "commodity_id", "created_at", "updated_at"}).
			AddRow(priceID, float64(100), regionID, commodityID, time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(commodityID, regionID, 1).WillReturnRows(price)
		result, err := repoImpl.FindByCommodityIDAndRegionID(context.TODO(), commodityID, regionID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, priceID, result.ID)
		assert.Equal(t, commodityID, result.CommodityID)
		assert.Equal(t, regionID, result.RegionID)
		assert.Equal(t, float64(100), result.Price)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test FindByCommodityIDAndRegionID, data base error", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(commodityID, regionID, 1).WillReturnError(errors.New("error"))
		result, err := repoImpl.FindByCommodityIDAndRegionID(context.TODO(), commodityID, regionID)
		assert.NotNil(t, err)
		assert.Nil(t, result)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test FindByCommodityIDAndRegionID, not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(commodityID, regionID, 1).WillReturnRows(sqlmock.NewRows([]string{"id", "price",
			"region_id", "commodity_id", "created_at", "updated_at"}))
		result, err := repoImpl.FindByCommodityIDAndRegionID(context.TODO(), commodityID, regionID)
		assert.NotNil(t, err)
		assert.Nil(t, result)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}
