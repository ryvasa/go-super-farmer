package repository_test

import (
	"context"
	"database/sql"
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

type PriceRepositoryIDs struct {
	PriceID     uuid.UUID
	CommodityID uuid.UUID
	RegionID    uuid.UUID
}

type PriceRepositoryMockRows struct {
	Price     *sqlmock.Rows
	Notfound  *sqlmock.Rows
	Commodity *sqlmock.Rows
	Region    *sqlmock.Rows
}

type PriceRepositoryMocDomain struct {
	Price *domain.Price
}

func PriceRepositorySetup(t *testing.T) (*sql.DB, sqlmock.Sqlmock, repository.PriceRepository, PriceRepositoryIDs, PriceRepositoryMockRows, PriceRepositoryMocDomain) {

	sqlDB, db, mock := database.DbMock(t)

	repo := repository.NewPriceRepository(db)

	priceID := uuid.New()
	commodityID := uuid.New()
	regionID := uuid.New()

	ids := PriceRepositoryIDs{
		PriceID:     priceID,
		CommodityID: commodityID,
		RegionID:    regionID,
	}

	rows := PriceRepositoryMockRows{
		Price: sqlmock.NewRows([]string{"id", "price", "commodity_id", "region_id", "created_at", "updated_at"}).
			AddRow(priceID, float64(100), commodityID, regionID, time.Now(), time.Now()),

		Notfound: sqlmock.NewRows([]string{"id", "price", "created_at", "updated_at"}),

		Commodity: sqlmock.NewRows([]string{"id", "name"}).
			AddRow(commodityID, "commodity name"),

		Region: sqlmock.NewRows([]string{"id", "name"}).
			AddRow(regionID, "region name"),
	}

	domains := PriceRepositoryMocDomain{
		Price: &domain.Price{
			ID:          priceID,
			CommodityID: commodityID,
			RegionID:    regionID,
			Price:       float64(100),
		},
	}

	return sqlDB, mock, repo, ids, rows, domains
}

func TestPriceRepository_Create(t *testing.T) {
	db, mock, repo, ids, _, domains := PriceRepositorySetup(t)

	defer db.Close()

	expectedSQL := `INSERT INTO "prices" ("id","commodity_id","region_id","price","unit","created_at","updated_at","deleted_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`

	t.Run("should not return error when create successfully", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.PriceID, ids.CommodityID, ids.RegionID, float64(100), "idr", sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Create(context.TODO(), domains.Price)
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when create failed", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.PriceID, ids.CommodityID, ids.RegionID, float64(100), "idr", sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
			WillReturnError(errors.New("database error"))
		mock.ExpectRollback()

		err := repo.Create(context.TODO(), domains.Price)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestPriceRepository_FindAll(t *testing.T) {
	db, mock, repo, ids, rows, _ := PriceRepositorySetup(t)

	defer db.Close()

	expectedSQL1 := `SELECT * FROM "prices"`

	expectedSQL2 := `SELECT * FROM "commodities" WHERE "commodities"."id" = $1 AND "commodities"."deleted_at" IS NULL`

	expectedSQL3 := `SELECT * FROM "regions" WHERE "regions"."id" = $1 AND "regions"."deleted_at" IS NULL`

	t.Run("should return price histories when find all successfully", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL1)).WillReturnRows(rows.Price)

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL2)).WithArgs(ids.CommodityID).WillReturnRows(rows.Commodity)

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL3)).WithArgs(ids.RegionID).WillReturnRows(rows.Region)

		result, err := repo.FindAll(context.TODO())
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Len(t, *result, 1)
		assert.Equal(t, ids.PriceID, (*result)[0].ID)
		assert.Equal(t, float64(100), (*result)[0].Price)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
	t.Run("should return error when find all failed", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL1)).WillReturnError(errors.New("database error"))

		result, err := repo.FindAll(context.TODO())
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestPriceRepository_FindByID(t *testing.T) {
	db, mock, repo, ids, rows, _ := PriceRepositorySetup(t)

	defer db.Close()

	expectedSQL1 := `SELECT * FROM "prices" WHERE "prices"."id" = $1 AND "prices"."deleted_at" IS NULL ORDER BY "prices"."id" LIMIT $2`

	expectedSQL2 := `SELECT "commodities"."id","commodities"."name","commodities"."code","commodities"."duration" FROM "commodities" WHERE "commodities"."id" = $1 AND "commodities"."deleted_at" IS NULL`

	expectedSQL3 := `SELECT "regions"."id","regions"."province_id","regions"."city_id" FROM "regions" WHERE "regions"."id" = $1`

	t.Run("should return price when find by id successfully", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL1)).WithArgs(ids.PriceID, 1).WillReturnRows(rows.Price)

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL2)).WithArgs(ids.CommodityID).WillReturnRows(rows.Commodity)

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL3)).WithArgs(ids.RegionID).WillReturnRows(rows.Region)

		result, err := repo.FindByID(context.TODO(), ids.PriceID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.PriceID, result.ID)
		assert.Equal(t, float64(100), result.Price)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
	t.Run("should return error when find by id failed", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL1)).WithArgs(ids.PriceID, 1).WillReturnError(errors.New("database error"))
		result, err := repo.FindByID(context.TODO(), ids.PriceID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})
	t.Run("should return error when find by id not found", func(t *testing.T) {
		price := sqlmock.NewRows([]string{"id", "price", "created_at", "updated_at"})
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL1)).WithArgs(ids.PriceID, 1).WillReturnRows(price)
		result, err := repo.FindByID(context.TODO(), ids.PriceID)
		assert.Nil(t, result)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestPriceRepository_FindByCommodityID(t *testing.T) {
	db, mock, repo, ids, rows, _ := PriceRepositorySetup(t)

	defer db.Close()

	expectedSQL1 := `SELECT * FROM "prices" WHERE prices.commodity_id = $1 AND "prices"."deleted_at" IS NULL`

	expectedSQL2 := `SELECT "commodities"."id","commodities"."name","commodities"."code","commodities"."duration" FROM "commodities" WHERE "commodities"."id" = $1 AND "commodities"."deleted_at" IS NULL`

	expectedSQL3 := `SELECT "regions"."id","regions"."province_id","regions"."city_id" FROM "regions" WHERE "regions"."id" = $1 AND "regions"."deleted_at" IS NULL`

	t.Run("should return prices when find by commodity id successfully", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL1)).WithArgs(ids.CommodityID).WillReturnRows(rows.Price)

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL2)).WithArgs(ids.CommodityID).WillReturnRows(rows.Commodity)

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL3)).WithArgs(ids.RegionID).WillReturnRows(rows.Region)

		result, err := repo.FindByCommodityID(context.TODO(), ids.CommodityID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.PriceID, (*result)[0].ID)
		assert.Equal(t, ids.CommodityID, (*result)[0].CommodityID)
		assert.Equal(t, ids.RegionID, (*result)[0].RegionID)
		assert.Equal(t, float64(100), (*result)[0].Price)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when find by commodity id failed", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL1)).WithArgs(ids.CommodityID).WillReturnError(errors.New("database error"))
		result, err := repo.FindByCommodityID(context.TODO(), ids.CommodityID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestPriceRepository_FindByRegionID(t *testing.T) {
	db, mock, repo, ids, rows, _ := PriceRepositorySetup(t)

	defer db.Close()

	expectedSQL1 := `SELECT * FROM "prices" WHERE prices.region_id = $1 AND "prices"."deleted_at" IS NULL`

	expectedSQL3 := `SELECT "regions"."id","regions"."province_id","regions"."city_id" FROM "regions" WHERE "regions"."id" = $1 AND "regions"."deleted_at" IS NULL`

	expectedSQL2 := `SELECT "commodities"."id","commodities"."name","commodities"."code","commodities"."duration" FROM "commodities" WHERE "commodities"."id" = $1 AND "commodities"."deleted_at" IS NULL`

	t.Run("should return prices when find by region id successfully", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL1)).WithArgs(ids.RegionID).WillReturnRows(rows.Price)

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL2)).WithArgs(ids.CommodityID).WillReturnRows(rows.Commodity)

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL3)).WithArgs(ids.RegionID).WillReturnRows(rows.Region)

		result, err := repo.FindByRegionID(context.TODO(), ids.RegionID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.PriceID, (*result)[0].ID)
		assert.Equal(t, ids.CommodityID, (*result)[0].CommodityID)
		assert.Equal(t, ids.RegionID, (*result)[0].RegionID)
		assert.Equal(t, float64(100), (*result)[0].Price)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when find by region id failed", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL1)).WithArgs(ids.RegionID).WillReturnError(errors.New("database error"))
		result, err := repo.FindByRegionID(context.TODO(), ids.RegionID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestPriceRepository_Update(t *testing.T) {
	db, mock, repo, ids, _, domains := PriceRepositorySetup(t)

	defer db.Close()

	expectedSQL := `UPDATE "prices" SET "id"=$1,"commodity_id"=$2,"region_id"=$3,"price"=$4,"updated_at"=$5 WHERE id = $6 AND "prices"."deleted_at" IS NULL`

	t.Run("should not return error when update successfully", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.PriceID, ids.CommodityID, ids.RegionID, float64(100), sqlmock.AnyArg(), ids.PriceID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Update(context.TODO(), ids.PriceID, domains.Price)
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when update failed", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.PriceID, ids.CommodityID, ids.RegionID, float64(100), sqlmock.AnyArg(), ids.PriceID).
			WillReturnError(errors.New("database error"))
		mock.ExpectRollback()

		err := repo.Update(context.TODO(), ids.PriceID, domains.Price)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when update not found", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.PriceID, ids.CommodityID, ids.RegionID, float64(100), sqlmock.AnyArg(), ids.PriceID).
			WillReturnError(gorm.ErrRecordNotFound)
		mock.ExpectRollback()

		err := repo.Update(context.TODO(), ids.PriceID, domains.Price)
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestPriceRepository_Delete(t *testing.T) {
	db, mock, repo, ids, _, _ := PriceRepositorySetup(t)

	defer db.Close()

	expectedSQL := `UPDATE "prices" SET "deleted_at"=$1 WHERE id = $2 AND "prices"."deleted_at" IS NULL`

	t.Run("should not return error when delete successfully", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), ids.PriceID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Delete(context.TODO(), ids.PriceID)
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when delete failed", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), ids.PriceID).
			WillReturnError(errors.New("database error"))
		mock.ExpectRollback()

		err := repo.Delete(context.TODO(), ids.PriceID)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when delete not found", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), ids.PriceID).
			WillReturnError(gorm.ErrRecordNotFound)
		mock.ExpectRollback()

		err := repo.Delete(context.TODO(), ids.PriceID)
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestPriceRepository_Restore(t *testing.T) {
	db, mock, repo, ids, _, _ := PriceRepositorySetup(t)

	defer db.Close()

	expectedSQL := `UPDATE "prices" SET "deleted_at"=$1,"updated_at"=$2 WHERE id = $3`

	t.Run("should not return error when restore successfully", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), ids.PriceID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Restore(context.TODO(), ids.PriceID)
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when restore failed", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), ids.PriceID).
			WillReturnError(errors.New("database error"))
		mock.ExpectRollback()

		err := repo.Restore(context.TODO(), ids.PriceID)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when restore not found", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), ids.PriceID).
			WillReturnError(gorm.ErrRecordNotFound)
		mock.ExpectRollback()

		err := repo.Restore(context.TODO(), ids.PriceID)
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestPriceRepository_FindDeletedByID(t *testing.T) {
	db, mock, repo, ids, rows, _ := PriceRepositorySetup(t)

	defer db.Close()

	expectedSQL1 := `SELECT * FROM "prices" WHERE prices.id = $1 AND prices.deleted_at IS NOT NULL ORDER BY "prices"."id" LIMIT $2`

	expectedSQL2 := `SELECT "commodities"."id","commodities"."name","commodities"."code","commodities"."duration" FROM "commodities" WHERE "commodities"."id" = $1`

	expectedSQL3 := `SELECT "regions"."id","regions"."province_id","regions"."city_id" FROM "regions" WHERE "regions"."id" = $1`

	t.Run("should return price when find deleted by id successfully", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL1)).WithArgs(ids.PriceID, 1).WillReturnRows(rows.Price)

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL2)).WithArgs(ids.CommodityID).WillReturnRows(rows.Commodity)

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL3)).WithArgs(ids.RegionID).WillReturnRows(rows.Region)

		result, err := repo.FindDeletedByID(context.TODO(), ids.PriceID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.PriceID, result.ID)
		assert.Equal(t, float64(100), result.Price)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
	t.Run("should return error when find deleted by id failed", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL1)).WithArgs(ids.PriceID, 1).WillReturnError(errors.New("database error"))
		result, err := repo.FindDeletedByID(context.TODO(), ids.PriceID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})
	t.Run("should return error when find deleted by id not found", func(t *testing.T) {
		price := sqlmock.NewRows([]string{"id", "price", "created_at", "updated_at"})
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL1)).WithArgs(ids.PriceID, 1).WillReturnRows(price)
		result, err := repo.FindDeletedByID(context.TODO(), ids.PriceID)
		assert.Nil(t, result)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestPriceRepository_FindByCommodityIDAndRegionID(t *testing.T) {
	db, mock, repo, ids, rows, _ := PriceRepositorySetup(t)

	defer db.Close()

	expectedSQL := `SELECT * FROM "prices" WHERE (prices.commodity_id = $1 AND prices.region_id = $2) AND "prices"."deleted_at" IS NULL ORDER BY "prices"."id" LIMIT $3`

	t.Run("should return price when find by commodity id and region id successfully", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.CommodityID, ids.RegionID, 1).WillReturnRows(rows.Price)
		result, err := repo.FindByCommodityIDAndRegionID(context.TODO(), ids.CommodityID, ids.RegionID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.PriceID, result.ID)
		assert.Equal(t, ids.CommodityID, result.CommodityID)
		assert.Equal(t, ids.RegionID, result.RegionID)
		assert.Equal(t, float64(100), result.Price)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when find by commodity id and region id failed", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.CommodityID, ids.RegionID, 1).WillReturnError(errors.New("database error"))
		result, err := repo.FindByCommodityIDAndRegionID(context.TODO(), ids.CommodityID, ids.RegionID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when find by commodity id and region id not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.CommodityID, ids.RegionID, 1).WillReturnRows(rows.Notfound)
		result, err := repo.FindByCommodityIDAndRegionID(context.TODO(), ids.CommodityID, ids.RegionID)
		assert.Nil(t, result)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}
