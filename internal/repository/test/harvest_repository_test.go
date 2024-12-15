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
	repository_implementation "github.com/ryvasa/go-super-farmer/internal/repository/implementation"
	repository_interface "github.com/ryvasa/go-super-farmer/internal/repository/interface"
	"github.com/ryvasa/go-super-farmer/pkg/database"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type HarvestIDs struct {
	HarvestID       uuid.UUID
	LandCommodityID uuid.UUID
	RegionID        uuid.UUID
	LandID          uuid.UUID
	CommodityID     uuid.UUID
}

type HarvestMockRows struct {
	Harvest       *sqlmock.Rows
	Notfound      *sqlmock.Rows
	LandCommodity *sqlmock.Rows
	Region        *sqlmock.Rows
}

type HarvestMocDomain struct {
	Harvest *domain.Harvest
}

func HarvestRepositorySetup(t *testing.T) (*sql.DB, sqlmock.Sqlmock, repository_interface.HarvestRepository, HarvestIDs, HarvestMockRows, HarvestMocDomain) {

	sqlDB, db, mock := database.DbMock(t)

	repo := repository_implementation.NewHarvestRepository(db)

	harvestID := uuid.New()
	landCommodityID := uuid.New()
	regionID := uuid.New()
	landID := uuid.New()
	commodityID := uuid.New()

	date, _ := time.Parse("2006-01-02", "2022-01-01")

	ids := HarvestIDs{
		HarvestID:       harvestID,
		LandCommodityID: landCommodityID,
		RegionID:        regionID,
		LandID:          landID,
		CommodityID:     commodityID,
	}

	rows := HarvestMockRows{
		Harvest: sqlmock.NewRows([]string{"id", "land_commodity_id", "region_id", "quantity", "unit", "harvest_date", "created_at", "updated_at", "deleted_at"}).
			AddRow(harvestID, landCommodityID, regionID, float64(100), "kg", date, date, date, nil),

		Notfound: sqlmock.NewRows([]string{"id", "land_commodity_id", "region_id", "quantity", "unit", "harvest_date", "created_at", "updated_at", "deleted_at"}),

		LandCommodity: sqlmock.NewRows([]string{"id", "land_area", "commodity_id", "land_id", "created_at", "updated_at"}).
			AddRow(landCommodityID, float64(100), commodityID, landID, date, date),

		Region: sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
			AddRow(regionID, "region name", date, date),
	}

	domains := HarvestMocDomain{
		Harvest: &domain.Harvest{
			ID:              harvestID,
			LandCommodityID: landCommodityID,
			RegionID:        regionID,
			Quantity:        float64(100),
			Unit:            "kg",
			HarvestDate:     date,
		},
	}

	return sqlDB, mock, repo, ids, rows, domains
}

func TestHarvestRepository_Create(t *testing.T) {
	db, mock, repo, ids, _, domains := HarvestRepositorySetup(t)

	defer db.Close()

	expectedSQL := `INSERT INTO "harvests" ("id","land_commodity_id","region_id","quantity","unit","harvest_date","created_at","updated_at","deleted_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`

	t.Run("should not return error when create successfully", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.HarvestID, ids.LandCommodityID, ids.RegionID, float64(100), "kg", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Create(context.TODO(), domains.Harvest)
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when create failed", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.HarvestID, ids.LandCommodityID, ids.RegionID, float64(100), "kg", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
			WillReturnError(errors.New("database error"))
		mock.ExpectRollback()

		err := repo.Create(context.TODO(), domains.Harvest)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestHarvestRepository_FindByID(t *testing.T) {
	db, mock, repo, ids, rows, domains := HarvestRepositorySetup(t)

	defer db.Close()

	expectedSQL := `SELECT * FROM "harvests" WHERE "harvests"."id" = $1 AND "harvests"."deleted_at" IS NULL ORDER BY "harvests"."id" LIMIT $2`

	t.Run("should return harvest when find by id successfully", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.HarvestID, 1).WillReturnRows(rows.Harvest)

		result, err := repo.FindByID(context.TODO(), ids.HarvestID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.HarvestID, result.ID)
		assert.Equal(t, ids.LandCommodityID, result.LandCommodityID)
		assert.Equal(t, ids.RegionID, result.RegionID)
		assert.Equal(t, float64(100), result.Quantity)
		assert.Equal(t, "kg", result.Unit)
		assert.Equal(t, domains.Harvest.HarvestDate, result.HarvestDate)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
	t.Run("should return error when find by id failed", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.HarvestID, 1).WillReturnError(errors.New("database error"))

		result, err := repo.FindByID(context.TODO(), ids.HarvestID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})
	t.Run("should return error when find by id not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.HarvestID, 1).WillReturnRows(rows.Notfound)

		result, err := repo.FindByID(context.TODO(), ids.HarvestID)
		assert.Nil(t, result)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestHarvestRepository_FindAll(t *testing.T) {
	db, mock, repo, ids, rows, domains := HarvestRepositorySetup(t)

	defer db.Close()

	expectedSQL := `SELECT * FROM "harvests"`

	t.Run("should return harvests when find all successfully", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WillReturnRows(rows.Harvest)

		result, err := repo.FindAll(context.TODO())
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 1, len(*result))
		assert.Equal(t, ids.HarvestID, (*result)[0].ID)
		assert.Equal(t, ids.LandCommodityID, (*result)[0].LandCommodityID)
		assert.Equal(t, ids.RegionID, (*result)[0].RegionID)
		assert.Equal(t, float64(100), (*result)[0].Quantity)
		assert.Equal(t, "kg", (*result)[0].Unit)
		assert.Equal(t, domains.Harvest.HarvestDate, (*result)[0].HarvestDate)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
	t.Run("should return error when find all failed", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WillReturnError(errors.New("database error"))

		result, err := repo.FindAll(context.TODO())
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestHarvestRepository_FindByLandCommodityID(t *testing.T) {
	db, mock, repo, ids, rows, domains := HarvestRepositorySetup(t)

	defer db.Close()

	expectedSQL := `SELECT * FROM "harvests" WHERE land_commodity_id = $1 AND "harvests"."deleted_at" IS NULL`

	t.Run("should return harvests when find by land commodity id successfully", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.LandCommodityID).WillReturnRows(rows.Harvest)

		result, err := repo.FindByLandCommodityID(context.TODO(), ids.LandCommodityID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.HarvestID, (*result)[0].ID)
		assert.Equal(t, ids.LandCommodityID, (*result)[0].LandCommodityID)
		assert.Equal(t, ids.RegionID, (*result)[0].RegionID)
		assert.Equal(t, float64(100), (*result)[0].Quantity)
		assert.Equal(t, "kg", (*result)[0].Unit)
		assert.Equal(t, domains.Harvest.HarvestDate, (*result)[0].HarvestDate)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when find by land commodity id failed", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.LandCommodityID).WillReturnError(errors.New("database error"))

		result, err := repo.FindByLandCommodityID(context.TODO(), ids.LandCommodityID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestHarvestRepository_FindByRegionID(t *testing.T) {
	db, mock, repo, ids, rows, domains := HarvestRepositorySetup(t)

	defer db.Close()

	expectedSQL := `SELECT * FROM "harvests" WHERE region_id = $1 AND "harvests"."deleted_at" IS NULL`

	t.Run("should return harvests when find by region id successfully", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.RegionID).WillReturnRows(rows.Harvest)

		result, err := repo.FindByRegionID(context.TODO(), ids.RegionID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.HarvestID, (*result)[0].ID)
		assert.Equal(t, ids.LandCommodityID, (*result)[0].LandCommodityID)
		assert.Equal(t, ids.RegionID, (*result)[0].RegionID)
		assert.Equal(t, float64(100), (*result)[0].Quantity)
		assert.Equal(t, "kg", (*result)[0].Unit)
		assert.Equal(t, domains.Harvest.HarvestDate, (*result)[0].HarvestDate)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when find by region id failed", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.RegionID).WillReturnError(errors.New("database error"))

		result, err := repo.FindByRegionID(context.TODO(), ids.RegionID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestHarvestRepository_FindByLandID(t *testing.T) {
	db, mock, repo, ids, rows, domains := HarvestRepositorySetup(t)

	defer db.Close()

	expectedSQL1 := `SELECT "harvests"."id","harvests"."land_commodity_id","harvests"."region_id","harvests"."quantity","harvests"."unit","harvests"."harvest_date","harvests"."created_at","harvests"."updated_at","harvests"."deleted_at" FROM "harvests" JOIN land_commodities ON harvests.land_commodity_id = land_commodities.id WHERE land_commodities.land_id = $1 AND "harvests"."deleted_at" IS NULL`

	expectedSQL2 := `SELECT * FROM "land_commodities" WHERE "land_commodities"."id" = $1 AND "land_commodities"."deleted_at" IS NULL`

	t.Run("should return harvests when find by land id successfully", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL1)).WithArgs(ids.LandID).WillReturnRows(rows.Harvest)

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL2)).
			WithArgs(ids.LandCommodityID).
			WillReturnRows(rows.LandCommodity)

		result, err := repo.FindByLandID(context.TODO(), ids.LandID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.HarvestID, (*result)[0].ID)
		assert.Equal(t, ids.LandCommodityID, (*result)[0].LandCommodityID)
		assert.Equal(t, ids.LandID, (*result)[0].LandCommodity.LandID)
		assert.Equal(t, float64(100), (*result)[0].Quantity)
		assert.Equal(t, "kg", (*result)[0].Unit)
		assert.Equal(t, domains.Harvest.HarvestDate, (*result)[0].HarvestDate)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when find by land id failed", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL1)).WithArgs(ids.LandID).WillReturnError(errors.New("database error"))

		result, err := repo.FindByLandID(context.TODO(), ids.LandID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestHarvestRepository_FindByCommodityID(t *testing.T) {
	db, mock, repo, ids, rows, domains := HarvestRepositorySetup(t)

	defer db.Close()

	expectedSQL1 := `SELECT "harvests"."id","harvests"."land_commodity_id","harvests"."region_id","harvests"."quantity","harvests"."unit","harvests"."harvest_date","harvests"."created_at","harvests"."updated_at","harvests"."deleted_at" FROM "harvests" JOIN land_commodities ON harvests.land_commodity_id = land_commodities.id WHERE land_commodities.commodity_id = $1 AND "harvests"."deleted_at" IS NULL`

	expectedSQL2 := `SELECT * FROM "land_commodities" WHERE "land_commodities"."id" = $1 AND "land_commodities"."deleted_at" IS NULL`

	t.Run("should return harvests when find by commodity id successfully", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL1)).WithArgs(ids.CommodityID).WillReturnRows(rows.Harvest)

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL2)).
			WithArgs(ids.LandCommodityID).
			WillReturnRows(rows.LandCommodity)

		result, err := repo.FindByCommodityID(context.TODO(), ids.CommodityID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.HarvestID, (*result)[0].ID)
		assert.Equal(t, ids.LandCommodityID, (*result)[0].LandCommodityID)
		assert.Equal(t, ids.CommodityID, (*result)[0].LandCommodity.CommodityID)
		assert.Equal(t, float64(100), (*result)[0].Quantity)
		assert.Equal(t, "kg", (*result)[0].Unit)
		assert.Equal(t, domains.Harvest.HarvestDate, (*result)[0].HarvestDate)
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

func TestHarvestRepository_Update(t *testing.T) {
	db, mock, repo, ids, _, domains := HarvestRepositorySetup(t)

	defer db.Close()

	expectedSQL := `UPDATE "harvests" SET "id"=$1,"land_commodity_id"=$2,"region_id"=$3,"quantity"=$4,"unit"=$5,"harvest_date"=$6,"updated_at"=$7 WHERE id = $8 AND "harvests"."deleted_at" IS NULL`

	t.Run("should not return error when update successfully", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.HarvestID, ids.LandCommodityID, ids.RegionID, float64(100), "kg", sqlmock.AnyArg(), sqlmock.AnyArg(), ids.HarvestID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Update(context.TODO(), ids.HarvestID, domains.Harvest)
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when update failed", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.HarvestID, ids.LandCommodityID, ids.RegionID, float64(100), "kg", sqlmock.AnyArg(), sqlmock.AnyArg(), ids.HarvestID).
			WillReturnError(errors.New("database error"))
		mock.ExpectRollback()

		err := repo.Update(context.TODO(), ids.HarvestID, domains.Harvest)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when update not found", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.HarvestID, ids.LandCommodityID, ids.RegionID, float64(100), "kg", sqlmock.AnyArg(), sqlmock.AnyArg(), ids.HarvestID).
			WillReturnError(gorm.ErrRecordNotFound)
		mock.ExpectRollback()

		err := repo.Update(context.TODO(), ids.HarvestID, domains.Harvest)
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestHarvestRepository_Delete(t *testing.T) {
	db, mock, repo, ids, _, _ := HarvestRepositorySetup(t)

	defer db.Close()

	expectedSQL := `UPDATE "harvests" SET "deleted_at"=$1 WHERE id = $2 AND "harvests"."deleted_at" IS NULL`

	t.Run("should not return error when delete successfully", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), ids.HarvestID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Delete(context.TODO(), ids.HarvestID)
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when delete failed", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), ids.HarvestID).
			WillReturnError(errors.New("database error"))
		mock.ExpectRollback()

		err := repo.Delete(context.TODO(), ids.HarvestID)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when delete not found", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), ids.HarvestID).
			WillReturnError(gorm.ErrRecordNotFound)
		mock.ExpectRollback()

		err := repo.Delete(context.TODO(), ids.HarvestID)
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestHarvestRepository_Restore(t *testing.T) {
	db, mock, repo, ids, _, _ := HarvestRepositorySetup(t)

	defer db.Close()

	expectedSQL := `UPDATE "harvests" SET "deleted_at"=$1,"updated_at"=$2 WHERE id = $3`

	t.Run("should not return error when restore successfully", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), ids.HarvestID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Restore(context.TODO(), ids.HarvestID)
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when restore failed", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), ids.HarvestID).
			WillReturnError(errors.New("database error"))
		mock.ExpectRollback()

		err := repo.Restore(context.TODO(), ids.HarvestID)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when restore not found", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), ids.HarvestID).
			WillReturnError(gorm.ErrRecordNotFound)
		mock.ExpectRollback()

		err := repo.Restore(context.TODO(), ids.HarvestID)
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestHarvestRepository_FindAllDeleted(t *testing.T) {
	db, mock, repo, ids, rows, domains := HarvestRepositorySetup(t)

	defer db.Close()

	expectedSQL := `SELECT * FROM "harvests"`

	t.Run("should return harvests when find all deleted successfully", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WillReturnRows(rows.Harvest)

		result, err := repo.FindAllDeleted(context.TODO())
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 1, len(*result))
		assert.Equal(t, ids.HarvestID, (*result)[0].ID)
		assert.Equal(t, ids.LandCommodityID, (*result)[0].LandCommodityID)
		assert.Equal(t, ids.RegionID, (*result)[0].RegionID)
		assert.Equal(t, float64(100), (*result)[0].Quantity)
		assert.Equal(t, "kg", (*result)[0].Unit)
		assert.Equal(t, domains.Harvest.HarvestDate, (*result)[0].HarvestDate)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
	t.Run("should return error when find all deleted failed", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WillReturnError(errors.New("database error"))

		result, err := repo.FindAllDeleted(context.TODO())
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestHarvestRepository_FindDeletedByID(t *testing.T) {
	db, mock, repo, ids, rows, domains := HarvestRepositorySetup(t)

	defer db.Close()

	expectedSQL := `SELECT * FROM "harvests" WHERE id = $1 AND deleted_at IS NOT NULL ORDER BY "harvests"."id" LIMIT $2`

	t.Run("should return harvest when find deleted by id successfully", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.HarvestID, 1).WillReturnRows(rows.Harvest)

		result, err := repo.FindDeletedByID(context.TODO(), ids.HarvestID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.HarvestID, result.ID)
		assert.Equal(t, ids.LandCommodityID, result.LandCommodityID)
		assert.Equal(t, ids.RegionID, result.RegionID)
		assert.Equal(t, float64(100), result.Quantity)
		assert.Equal(t, "kg", result.Unit)
		assert.Equal(t, domains.Harvest.HarvestDate, result.HarvestDate)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
	t.Run("should return error when find deleted by id failed", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.HarvestID, 1).WillReturnError(errors.New("database error"))

		result, err := repo.FindDeletedByID(context.TODO(), ids.HarvestID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})
	t.Run("should return error when find deleted by id not found", func(t *testing.T) {
		commodity := sqlmock.NewRows([]string{"id", "land_area", "commodity_id", "land_id", "created_at", "updated_at", "deleted_at"})
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.HarvestID, 1).WillReturnRows(commodity)

		result, err := repo.FindDeletedByID(context.TODO(), ids.HarvestID)
		assert.Nil(t, result)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}
