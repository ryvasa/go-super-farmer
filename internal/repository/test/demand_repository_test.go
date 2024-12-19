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

type DemandIDs struct {
	DemandID    uuid.UUID
	RegionID    uuid.UUID
	CommodityID uuid.UUID
}

type DemandMockRows struct {
	Demand  *sqlmock.Rows
	Demands *sqlmock.Rows
}

type DemandMocDomain struct {
	Demand *domain.Demand
}

func DemandRepositorySetup(t *testing.T) (*database.MockDB, repository_interface.DemandRepository, DemandIDs, DemandMockRows, DemandMocDomain) {

	mockDB := database.NewMockDB(t)

	repo := repository_implementation.NewDemandRepository(mockDB.BaseRepo)

	demandID := uuid.New()
	regionID := uuid.New()
	commodityID := uuid.New()

	ids := DemandIDs{
		DemandID:    demandID,
		RegionID:    regionID,
		CommodityID: commodityID,
	}

	rows := DemandMockRows{
		Demand: sqlmock.NewRows([]string{"id", "region_id", "commodity_id", "quantity", "created_at", "updated_at"}).
			AddRow(demandID, regionID, commodityID, float64(10), time.Now(), time.Now()),

		Demands: sqlmock.NewRows([]string{"id", "region_id", "commodity_id", "quantity", "created_at", "updated_at"}).
			AddRow(demandID, regionID, commodityID, float64(10), time.Now(), time.Now()).
			AddRow(uuid.New(), regionID, commodityID, float64(10), time.Now(), time.Now()),
	}

	domains := DemandMocDomain{
		Demand: &domain.Demand{
			ID:          demandID,
			RegionID:    regionID,
			CommodityID: commodityID,
			Quantity:    float64(10),
		},
	}

	return mockDB, repo, ids, rows, domains
}

func TestDemandRepository_Create(t *testing.T) {
	mockDB, repo, ids, _, domains := DemandRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL := `INSERT INTO "demands" ("id","commodity_id","region_id","quantity","unit","created_at","updated_at","deleted_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`

	t.Run("should not return error when create successfully", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.DemandID, ids.CommodityID, ids.RegionID, float64(10), "kg", sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mockDB.Mock.ExpectCommit()

		err := repo.Create(context.TODO(), domains.Demand)
		assert.Nil(t, err)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when create failed", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.DemandID, ids.CommodityID, ids.RegionID, float64(10), "kg", sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
			WillReturnError(errors.New("database error"))
		mockDB.Mock.ExpectRollback()

		err := repo.Create(context.TODO(), domains.Demand)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}

func TestDemandRepository_FindAll(t *testing.T) {
	mockDB, repo, ids, rows, _ := DemandRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL := `SELECT * FROM "demands"`

	t.Run("should return demands when find all successfully", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WillReturnRows(rows.Demands)

		result, err := repo.FindAll(context.TODO())
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 2)
		assert.Equal(t, ids.DemandID, (result)[0].ID)
		assert.Equal(t, ids.RegionID, (result)[0].RegionID)
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

func TestDemandRepository_FindByID(t *testing.T) {
	mockDB, repo, ids, rows, _ := DemandRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL := `SELECT * FROM "demands" WHERE "demands"."id" = $1 AND "demands"."deleted_at" IS NULL ORDER BY "demands"."id" LIMIT $2`

	t.Run("should return demand when find by id successfully", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.DemandID, 1).WillReturnRows(rows.Demand)

		result, err := repo.FindByID(context.TODO(), ids.DemandID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.DemandID, result.ID)
		assert.Equal(t, ids.RegionID, result.RegionID)
		assert.Equal(t, ids.CommodityID, result.CommodityID)
		assert.Equal(t, float64(10), result.Quantity)

		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when find by id failed", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.DemandID, 1).WillReturnError(errors.New("database error"))

		result, err := repo.FindByID(context.TODO(), ids.DemandID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")

		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}

func TestDemandRepository_FindByCommodityID(t *testing.T) {
	mockDB, repo, ids, rows, _ := DemandRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL := `SELECT * FROM "demands" WHERE commodity_id = $1 AND "demands"."deleted_at" IS NULL`

	t.Run("should return demands when find by commodity id successfully", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.CommodityID).WillReturnRows(rows.Demands)

		result, err := repo.FindByCommodityID(context.TODO(), ids.CommodityID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.DemandID, (result)[0].ID)
		assert.Equal(t, ids.RegionID, (result)[0].RegionID)
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

func TestDemandRepository_FindByRegionID(t *testing.T) {
	mockDB, repo, ids, rows, _ := DemandRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL := `SELECT * FROM "demands" WHERE region_id = $1 AND "demands"."deleted_at" IS NULL`

	t.Run("should return demands when find by region id successfully", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.RegionID).WillReturnRows(rows.Demands)

		result, err := repo.FindByRegionID(context.TODO(), ids.RegionID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.DemandID, (result)[0].ID)
		assert.Equal(t, ids.RegionID, (result)[0].RegionID)
		assert.Equal(t, ids.CommodityID, (result)[0].CommodityID)
		assert.Equal(t, float64(10), (result)[0].Quantity)

		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when find by region id failed", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.RegionID).WillReturnError(errors.New("database error"))

		result, err := repo.FindByRegionID(context.TODO(), ids.RegionID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")

		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}

func TestDemandRepository_Delete(t *testing.T) {
	mockDB, repo, ids, _, _ := DemandRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL := `UPDATE "demands" SET "deleted_at"=$1 WHERE id = $2 AND "demands"."deleted_at" IS NULL`

	t.Run("should not return error when delete successfully", func(t *testing.T) {

		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).WithArgs(sqlmock.AnyArg(), ids.DemandID).WillReturnResult(sqlmock.NewResult(1, 1))
		mockDB.Mock.ExpectCommit()

		err := repo.Delete(context.TODO(), ids.DemandID)
		assert.Nil(t, err)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when delete failed", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).WithArgs(sqlmock.AnyArg(), ids.DemandID).WillReturnError(errors.New("database error"))
		mockDB.Mock.ExpectRollback()

		err := repo.Delete(context.TODO(), ids.DemandID)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}

func TestDemandRepository_Update(t *testing.T) {
	mockDB, repo, ids, _, _ := DemandRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL := `UPDATE "demands" SET "quantity"=$1,"updated_at"=$2 WHERE id = $3 AND "demands"."deleted_at" IS NULL`

	t.Run("should not return error when update successfully", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).WithArgs(float64(10), sqlmock.AnyArg(), ids.DemandID).WillReturnResult(sqlmock.NewResult(1, 1))
		mockDB.Mock.ExpectCommit()

		err := repo.Update(context.TODO(), ids.DemandID, &domain.Demand{Quantity: 10})
		assert.Nil(t, err)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when update failed", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).WithArgs(float64(10), sqlmock.AnyArg(), ids.DemandID).WillReturnError(errors.New("database error"))
		mockDB.Mock.ExpectRollback()

		err := repo.Update(context.TODO(), ids.DemandID, &domain.Demand{Quantity: 10})
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}

func TestDemandRepository_FindByCommodityIDAndRegionID(t *testing.T) {
	mockDB, repo, ids, rows, _ := DemandRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL := `SELECT * FROM "demands" WHERE (commodity_id = $1 AND region_id = $2) AND "demands"."deleted_at" IS NULL ORDER BY "demands"."id" LIMIT $3`

	t.Run("should return demand when find by commodity id and region id successfully", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.CommodityID, ids.RegionID, 1).WillReturnRows(rows.Demand)

		result, err := repo.FindByCommodityIDAndRegionID(context.TODO(), ids.CommodityID, ids.RegionID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.DemandID, result.ID)
		assert.Equal(t, ids.RegionID, result.RegionID)
		assert.Equal(t, ids.CommodityID, result.CommodityID)
		assert.Equal(t, float64(10), result.Quantity)

		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when find by commodity id and region id failed", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.CommodityID, ids.RegionID, 1).WillReturnError(errors.New("database error"))

		result, err := repo.FindByCommodityIDAndRegionID(context.TODO(), ids.CommodityID, ids.RegionID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")

		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}
