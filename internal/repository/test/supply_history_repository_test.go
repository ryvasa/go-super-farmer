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

type SupplyHistoryIDs struct {
	SupplyHistoryID uuid.UUID
	RegionID        uuid.UUID
	CommodityID     uuid.UUID
}

type SupplyHistoryMockRows struct {
	SupplyHistory   *sqlmock.Rows
	SupplyHistories *sqlmock.Rows
}

type SupplyHistoryMocDomain struct {
	SupplyHistory *domain.SupplyHistory
}

func SupplyHistoryRepositorySetup(t *testing.T) (*sql.DB, sqlmock.Sqlmock, repository_interface.SupplyHistoryRepository, SupplyHistoryIDs, SupplyHistoryMockRows, SupplyHistoryMocDomain) {

	sqlDB, db, mock := database.DbMock(t)

	repo := repository_implementation.NewSupplyHistoryRepository(db)

	supplyHistoryID := uuid.New()
	regionID := uuid.New()
	commodityID := uuid.New()

	ids := SupplyHistoryIDs{
		SupplyHistoryID: supplyHistoryID,
		RegionID:        regionID,
		CommodityID:     commodityID,
	}

	rows := SupplyHistoryMockRows{
		SupplyHistory: sqlmock.NewRows([]string{"id", "region_id", "commodity_id", "quantity", "created_at", "updated_at"}).
			AddRow(supplyHistoryID, regionID, commodityID, float64(10), time.Now(), time.Now()),

		SupplyHistories: sqlmock.NewRows([]string{"id", "region_id", "commodity_id", "quantity", "created_at", "updated_at"}).
			AddRow(supplyHistoryID, regionID, commodityID, float64(10), time.Now(), time.Now()).
			AddRow(uuid.New(), regionID, commodityID, float64(10), time.Now(), time.Now()),
	}

	domains := SupplyHistoryMocDomain{
		SupplyHistory: &domain.SupplyHistory{
			ID:          supplyHistoryID,
			RegionID:    regionID,
			CommodityID: commodityID,
			Quantity:    float64(10),
		},
	}

	return sqlDB, mock, repo, ids, rows, domains
}

func TestSupplyHistoryRepository_Create(t *testing.T) {
	db, mock, repo, ids, _, domains := SupplyHistoryRepositorySetup(t)

	defer db.Close()

	expectedSQL := `INSERT INTO "supply_histories" ("id","commodity_id","region_id","quantity","unit","created_at","updated_at","deleted_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`

	t.Run("should not return error when create successfully", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.SupplyHistoryID, ids.CommodityID, ids.RegionID, float64(10), "kg", sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Create(context.TODO(), domains.SupplyHistory)
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when create failed", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.SupplyHistoryID, ids.CommodityID, ids.RegionID, float64(10), "kg", sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
			WillReturnError(errors.New("database error"))
		mock.ExpectRollback()

		err := repo.Create(context.TODO(), domains.SupplyHistory)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestSupplyHistoryRepository_FindAll(t *testing.T) {
	db, mock, repo, ids, rows, _ := SupplyHistoryRepositorySetup(t)

	defer db.Close()

	expectedSQL := `SELECT * FROM "supply_histories"`

	t.Run("should return supply histories when find all successfully", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WillReturnRows(rows.SupplyHistories)

		result, err := repo.FindAll(context.TODO())
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Len(t, *result, 2)
		assert.Equal(t, ids.SupplyHistoryID, (*result)[0].ID)
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

func TestSupplyHistoryRepository_FindByID(t *testing.T) {
	db, mock, repo, ids, rows, _ := SupplyHistoryRepositorySetup(t)

	defer db.Close()

	expectedSQL := `SELECT * FROM "supply_histories" WHERE "supply_histories"."id" = $1 AND "supply_histories"."deleted_at" IS NULL ORDER BY "supply_histories"."id" LIMIT $2`

	t.Run("should return supply history when find by id successfully", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.SupplyHistoryID, 1).WillReturnRows(rows.SupplyHistory)

		result, err := repo.FindByID(context.TODO(), ids.SupplyHistoryID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.SupplyHistoryID, result.ID)
		assert.Equal(t, ids.RegionID, result.RegionID)
		assert.Equal(t, ids.CommodityID, result.CommodityID)
		assert.Equal(t, float64(10), result.Quantity)

		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when find by id failed", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.SupplyHistoryID, 1).WillReturnError(errors.New("database error"))

		result, err := repo.FindByID(context.TODO(), ids.SupplyHistoryID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")

		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestSupplyHistoryRepository_FindByCommodityID(t *testing.T) {
	db, mock, repo, ids, rows, _ := SupplyHistoryRepositorySetup(t)

	defer db.Close()

	expectedSQL := `SELECT * FROM "supply_histories" WHERE commodity_id = $1 AND "supply_histories"."deleted_at" IS NULL`

	t.Run("should return supply histories when find by commodity id successfully", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.CommodityID).WillReturnRows(rows.SupplyHistories)

		result, err := repo.FindByCommodityID(context.TODO(), ids.CommodityID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.SupplyHistoryID, (*result)[0].ID)
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

func TestSupplyHistoryRepository_FindByRegionID(t *testing.T) {
	db, mock, repo, ids, rows, _ := SupplyHistoryRepositorySetup(t)

	defer db.Close()

	expectedSQL := `SELECT * FROM "supply_histories" WHERE region_id = $1 AND "supply_histories"."deleted_at" IS NULL`

	t.Run("should return supply histories when find by region id successfully", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.RegionID).WillReturnRows(rows.SupplyHistories)

		result, err := repo.FindByRegionID(context.TODO(), ids.RegionID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.SupplyHistoryID, (*result)[0].ID)
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

func TestSupplyHistoryRepository_FindByCommodityIDAndRegionID(t *testing.T) {
	db, mock, repo, ids, rows, _ := SupplyHistoryRepositorySetup(t)

	defer db.Close()

	expectedSQL := `SELECT * FROM "supply_histories" WHERE (commodity_id = $1 AND region_id = $2) AND "supply_histories"."deleted_at" IS NULL`

	t.Run("should return supply history when find by commodity id and region id successfully", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.CommodityID, ids.RegionID).WillReturnRows(rows.SupplyHistory)

		result, err := repo.FindByCommodityIDAndRegionID(context.TODO(), ids.CommodityID, ids.RegionID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.SupplyHistoryID, (*result)[0].ID)
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
