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

type DemandHistoryIDs struct {
	DemandHistoryID uuid.UUID
	RegionID        uuid.UUID
	CommodityID     uuid.UUID
}

type DemandHistoryMockRows struct {
	DemandHistory   *sqlmock.Rows
	DemandHistories *sqlmock.Rows
}

type DemandHistoryMocDomain struct {
	DemandHistory *domain.DemandHistory
}

func DemandHistoryRepositorySetup(t *testing.T) (*sql.DB, sqlmock.Sqlmock, repository_interface.DemandHistoryRepository, DemandHistoryIDs, DemandHistoryMockRows, DemandHistoryMocDomain) {

	sqlDB, db, mock := database.DbMock(t)

	repo := repository_implementation.NewDemandHistoryRepository(db)

	demandHistoryID := uuid.New()
	regionID := uuid.New()
	commodityID := uuid.New()

	ids := DemandHistoryIDs{
		DemandHistoryID: demandHistoryID,
		RegionID:        regionID,
		CommodityID:     commodityID,
	}

	rows := DemandHistoryMockRows{
		DemandHistory: sqlmock.NewRows([]string{"id", "region_id", "commodity_id", "quantity", "created_at", "updated_at"}).
			AddRow(demandHistoryID, regionID, commodityID, float64(10), time.Now(), time.Now()),

		DemandHistories: sqlmock.NewRows([]string{"id", "region_id", "commodity_id", "quantity", "created_at", "updated_at"}).
			AddRow(demandHistoryID, regionID, commodityID, float64(10), time.Now(), time.Now()).
			AddRow(uuid.New(), regionID, commodityID, float64(10), time.Now(), time.Now()),
	}

	domains := DemandHistoryMocDomain{
		DemandHistory: &domain.DemandHistory{
			ID:          demandHistoryID,
			RegionID:    regionID,
			CommodityID: commodityID,
			Quantity:    float64(10),
		},
	}

	return sqlDB, mock, repo, ids, rows, domains
}

func TestDemandHistoryRepository_Create(t *testing.T) {
	db, mock, repo, ids, _, domains := DemandHistoryRepositorySetup(t)

	defer db.Close()

	expectedSQL := `INSERT INTO "demand_histories" ("id","commodity_id","region_id","quantity","unit","created_at","updated_at","deleted_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`

	t.Run("should not return error when create successfully", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.DemandHistoryID, ids.CommodityID, ids.RegionID, float64(10), "kg", sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Create(context.TODO(), domains.DemandHistory)
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when create failed", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.DemandHistoryID, ids.CommodityID, ids.RegionID, float64(10), "kg", sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
			WillReturnError(errors.New("database error"))
		mock.ExpectRollback()

		err := repo.Create(context.TODO(), domains.DemandHistory)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestDemandHistoryRepository_FindAll(t *testing.T) {
	db, mock, repo, ids, rows, _ := DemandHistoryRepositorySetup(t)

	defer db.Close()

	expectedSQL := `SELECT * FROM "demand_histories"`

	t.Run("should return demand histories when find all successfully", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WillReturnRows(rows.DemandHistories)

		result, err := repo.FindAll(context.TODO())
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Len(t, *result, 2)
		assert.Equal(t, ids.DemandHistoryID, (*result)[0].ID)
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

func TestDemandHistoryRepository_FindByID(t *testing.T) {
	db, mock, repo, ids, rows, _ := DemandHistoryRepositorySetup(t)

	defer db.Close()

	expectedSQL := `SELECT * FROM "demand_histories" WHERE "demand_histories"."id" = $1 AND "demand_histories"."deleted_at" IS NULL ORDER BY "demand_histories"."id" LIMIT $2`

	t.Run("should return demand history when find by id successfully", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.DemandHistoryID, 1).WillReturnRows(rows.DemandHistory)

		result, err := repo.FindByID(context.TODO(), ids.DemandHistoryID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.DemandHistoryID, result.ID)
		assert.Equal(t, ids.RegionID, result.RegionID)
		assert.Equal(t, ids.CommodityID, result.CommodityID)
		assert.Equal(t, float64(10), result.Quantity)

		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when find by id failed", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.DemandHistoryID, 1).WillReturnError(errors.New("database error"))

		result, err := repo.FindByID(context.TODO(), ids.DemandHistoryID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")

		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestDemandHistoryRepository_FindByCommodityID(t *testing.T) {
	db, mock, repo, ids, rows, _ := DemandHistoryRepositorySetup(t)

	defer db.Close()

	expectedSQL := `SELECT * FROM "demand_histories" WHERE commodity_id = $1 AND "demand_histories"."deleted_at" IS NULL`

	t.Run("should return demand histories when find by commodity id successfully", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.CommodityID).WillReturnRows(rows.DemandHistories)

		result, err := repo.FindByCommodityID(context.TODO(), ids.CommodityID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.DemandHistoryID, (*result)[0].ID)
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

func TestDemandHistoryRepository_FindByRegionID(t *testing.T) {
	db, mock, repo, ids, rows, _ := DemandHistoryRepositorySetup(t)

	defer db.Close()

	expectedSQL := `SELECT * FROM "demand_histories" WHERE region_id = $1 AND "demand_histories"."deleted_at" IS NULL`

	t.Run("should return demand histories when find by region id successfully", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.RegionID).WillReturnRows(rows.DemandHistories)

		result, err := repo.FindByRegionID(context.TODO(), ids.RegionID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.DemandHistoryID, (*result)[0].ID)
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

func TestDemandHistoryRepository_FindByCommodityIDAndRegionID(t *testing.T) {
	db, mock, repo, ids, rows, _ := DemandHistoryRepositorySetup(t)

	defer db.Close()

	expectedSQL := `SELECT * FROM "demand_histories" WHERE (commodity_id = $1 AND region_id = $2) AND "demand_histories"."deleted_at" IS NULL`

	t.Run("should return demands history when find by commodity id and region id successfully", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.CommodityID, ids.RegionID).WillReturnRows(rows.DemandHistory)

		result, err := repo.FindByCommodityIDAndRegionID(context.TODO(), ids.CommodityID, ids.RegionID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.DemandHistoryID, (*result)[0].ID)
		assert.Equal(t, ids.RegionID, (*result)[0].RegionID)
		assert.Equal(t, ids.CommodityID, (*result)[0].CommodityID)
		assert.Equal(t, float64(10), (*result)[0].Quantity)

		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when find by commodity id and region id failed", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.CommodityID, ids.RegionID).WillReturnError(errors.New("database error"))

		result, err := repo.FindByCommodityIDAndRegionID(context.TODO(), ids.CommodityID, ids.RegionID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")

		assert.Nil(t, mock.ExpectationsWereMet())
	})
}
