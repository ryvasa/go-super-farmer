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
)

type SupplyID struct {
	SupplyID    uuid.UUID
	RegionID    uuid.UUID
	CommodityID uuid.UUID
}

type SupplyMockRows struct {
	Supply   *sqlmock.Rows
	Supplies *sqlmock.Rows
}

type SupplyMocDomain struct {
	Supply *domain.Supply
}

func SupplyRepositorySetup(t *testing.T) (*sql.DB, sqlmock.Sqlmock, repository_interface.SupplyRepository, SupplyID, SupplyMockRows, SupplyMocDomain) {

	sqlDB, db, mock := database.DbMock(t)

	repo := repository_implementation.NewSupplyRepository(db)

	supplyID := uuid.New()
	regionID := uuid.New()
	commodityID := uuid.New()

	ids := SupplyID{
		SupplyID:    supplyID,
		RegionID:    regionID,
		CommodityID: commodityID,
	}

	rows := SupplyMockRows{
		Supply: sqlmock.NewRows([]string{"id", "region_id", "commodity_id", "quantity", "created_at", "updated_at"}).
			AddRow(supplyID, regionID, commodityID, float64(10), time.Now(), time.Now()),

		Supplies: sqlmock.NewRows([]string{"id", "region_id", "commodity_id", "quantity", "created_at", "updated_at"}).
			AddRow(supplyID, regionID, commodityID, float64(10), time.Now(), time.Now()).
			AddRow(uuid.New(), regionID, commodityID, float64(10), time.Now(), time.Now()),
	}

	domains := SupplyMocDomain{
		Supply: &domain.Supply{
			ID:          supplyID,
			RegionID:    regionID,
			CommodityID: commodityID,
			Quantity:    float64(10),
		},
	}

	return sqlDB, mock, repo, ids, rows, domains
}

func TestSupplyRepository_Create(t *testing.T) {
	db, mock, repo, ids, _, domains := SupplyRepositorySetup(t)

	defer db.Close()

	expectedSQL := `INSERT INTO "supplies" ("id","commodity_id","region_id","quantity","unit","created_at","updated_at","deleted_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`

	t.Run("should not return error when create successfully", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.SupplyID, ids.CommodityID, ids.RegionID, float64(10), "kg", sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Create(context.TODO(), domains.Supply)
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when create failed", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.SupplyID, ids.CommodityID, ids.RegionID, float64(10), "kg", sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
			WillReturnError(errors.New("database error"))
		mock.ExpectRollback()

		err := repo.Create(context.TODO(), domains.Supply)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestSupplyRepository_FindAll(t *testing.T) {
	db, mock, repo, ids, rows, _ := SupplyRepositorySetup(t)

	defer db.Close()

	expectedSQL := `SELECT * FROM "supplies"`

	t.Run("should return supplies when find all successfully", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WillReturnRows(rows.Supplies)

		result, err := repo.FindAll(context.TODO())
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Len(t, *result, 2)
		assert.Equal(t, ids.SupplyID, (*result)[0].ID)
		assert.Equal(t, ids.RegionID, (*result)[0].RegionID)
		assert.Equal(t, ids.CommodityID, (*result)[0].CommodityID)
		assert.Equal(t, float64(10), (*result)[0].Quantity)

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

func TestSupplyRepository_FindByID(t *testing.T) {
	db, mock, repo, ids, rows, _ := SupplyRepositorySetup(t)

	defer db.Close()

	expectedSQL := `SELECT * FROM "supplies" WHERE "supplies"."id" = $1 AND "supplies"."deleted_at" IS NULL ORDER BY "supplies"."id" LIMIT $2`

	t.Run("should return supply when find by id successfully", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.SupplyID, 1).WillReturnRows(rows.Supply)

		result, err := repo.FindByID(context.TODO(), ids.SupplyID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.SupplyID, result.ID)
		assert.Equal(t, ids.RegionID, result.RegionID)
		assert.Equal(t, ids.CommodityID, result.CommodityID)
		assert.Equal(t, float64(10), result.Quantity)

		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when find by id failed", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.SupplyID, 1).WillReturnError(errors.New("database error"))

		result, err := repo.FindByID(context.TODO(), ids.SupplyID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")

		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestSupplyRepository_FindByCommodityID(t *testing.T) {
	db, mock, repo, ids, rows, _ := SupplyRepositorySetup(t)

	defer db.Close()

	expectedSQL := `SELECT * FROM "supplies" WHERE commodity_id = $1 AND "supplies"."deleted_at" IS NULL`

	t.Run("should return supplies when find by commodity id successfully", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.CommodityID).WillReturnRows(rows.Supplies)

		result, err := repo.FindByCommodityID(context.TODO(), ids.CommodityID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.SupplyID, (*result)[0].ID)
		assert.Equal(t, ids.RegionID, (*result)[0].RegionID)
		assert.Equal(t, ids.CommodityID, (*result)[0].CommodityID)
		assert.Equal(t, float64(10), (*result)[0].Quantity)

		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when find by commodity id failed", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.CommodityID).WillReturnError(errors.New("database error"))

		result, err := repo.FindByCommodityID(context.TODO(), ids.CommodityID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")

		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestSupplyRepository_FindByRegionID(t *testing.T) {
	db, mock, repo, ids, rows, _ := SupplyRepositorySetup(t)

	defer db.Close()

	expectedSQL := `SELECT * FROM "supplies" WHERE region_id = $1 AND "supplies"."deleted_at" IS NULL`

	t.Run("should return supplies when find by region id successfully", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.RegionID).WillReturnRows(rows.Supplies)

		result, err := repo.FindByRegionID(context.TODO(), ids.RegionID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.SupplyID, (*result)[0].ID)
		assert.Equal(t, ids.RegionID, (*result)[0].RegionID)
		assert.Equal(t, ids.CommodityID, (*result)[0].CommodityID)
		assert.Equal(t, float64(10), (*result)[0].Quantity)

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

func TestSupplyRepository_Delete(t *testing.T) {
	db, mock, repo, ids, _, _ := SupplyRepositorySetup(t)

	defer db.Close()

	expectedSQL := `UPDATE "supplies" SET "deleted_at"=$1 WHERE id = $2 AND "supplies"."deleted_at" IS NULL`

	t.Run("should not return error when delete successfully", func(t *testing.T) {

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).WithArgs(sqlmock.AnyArg(), ids.SupplyID).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Delete(context.TODO(), ids.SupplyID)
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when delete failed", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).WithArgs(sqlmock.AnyArg(), ids.SupplyID).WillReturnError(errors.New("database error"))
		mock.ExpectRollback()

		err := repo.Delete(context.TODO(), ids.SupplyID)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestSupplyRepository_Update(t *testing.T) {
	db, mock, repo, ids, _, _ := SupplyRepositorySetup(t)

	defer db.Close()

	expectedSQL := `UPDATE "supplies" SET "quantity"=$1,"updated_at"=$2 WHERE id = $3 AND "supplies"."deleted_at" IS NULL`

	t.Run("should not return error when update successfully", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).WithArgs(float64(10), sqlmock.AnyArg(), ids.SupplyID).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Update(context.TODO(), ids.SupplyID, &domain.Supply{Quantity: 10})
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when update failed", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).WithArgs(float64(10), sqlmock.AnyArg(), ids.SupplyID).WillReturnError(errors.New("database error"))
		mock.ExpectRollback()

		err := repo.Update(context.TODO(), ids.SupplyID, &domain.Supply{Quantity: 10})
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestSupplyRepository_FindByCommodityIDAndRegionID(t *testing.T) {
	db, mock, repo, ids, rows, _ := SupplyRepositorySetup(t)

	defer db.Close()

	expectedSQL := `SELECT * FROM "supplies" WHERE (commodity_id = $1 AND region_id = $2) AND "supplies"."deleted_at" IS NULL ORDER BY "supplies"."id" LIMIT $3`

	t.Run("should return supply when find by commodity id and region id successfully", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.CommodityID, ids.RegionID, 1).WillReturnRows(rows.Supply)

		result, err := repo.FindByCommodityIDAndRegionID(context.TODO(), ids.CommodityID, ids.RegionID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.SupplyID, result.ID)
		assert.Equal(t, ids.RegionID, result.RegionID)
		assert.Equal(t, ids.CommodityID, result.CommodityID)
		assert.Equal(t, float64(10), result.Quantity)

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
}
