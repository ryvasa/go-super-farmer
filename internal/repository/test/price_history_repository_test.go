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

type PriceHistoryIDs struct {
	PriceHistoryID uuid.UUID
	CommodityID    uuid.UUID
	RegionID       uuid.UUID
}

type PriceHistoryMockRows struct {
	PriceHistory *sqlmock.Rows
	Notfound     *sqlmock.Rows
}

type PriceHistoryMocDomain struct {
	PriceHistory *domain.PriceHistory
}

func PriceHistoryRepositorySetup(t *testing.T) (*sql.DB, sqlmock.Sqlmock, repository_interface.PriceHistoryRepository, PriceHistoryIDs, PriceHistoryMockRows, PriceHistoryMocDomain) {

	sqlDB, db, mock := database.DbMock(t)

	repo := repository_implementation.NewPriceHistoryRepository(db)

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
	}

	domains := PriceHistoryMocDomain{
		PriceHistory: &domain.PriceHistory{
			ID:          priceHistoryID,
			CommodityID: commodityID,
			RegionID:    regionID,
			Price:       float64(100),
		},
	}

	return sqlDB, mock, repo, ids, rows, domains
}

func TestPriceHistoryRepository_Create(t *testing.T) {
	db, mock, repo, ids, _, domains := PriceHistoryRepositorySetup(t)

	defer db.Close()

	expectedSQL := `INSERT INTO "price_histories" ("id","commodity_id","region_id","price","unit","created_at","updated_at","deleted_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`

	t.Run("should not return error when create successfully", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.PriceHistoryID, ids.CommodityID, ids.RegionID, float64(100), "idr", sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Create(context.TODO(), domains.PriceHistory)
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when create failed", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.PriceHistoryID, ids.CommodityID, ids.RegionID, float64(100), "idr", sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
			WillReturnError(errors.New("database error"))
		mock.ExpectRollback()

		err := repo.Create(context.TODO(), domains.PriceHistory)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestPriceHistoryRepository_FindAll(t *testing.T) {
	db, mock, repo, ids, rows, _ := PriceHistoryRepositorySetup(t)

	defer db.Close()

	expectedSQL := `SELECT * FROM "price_histories"`

	t.Run("should return price histories when find all successfully", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WillReturnRows(rows.PriceHistory)

		result, err := repo.FindAll(context.TODO())
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Len(t, *result, 1)
		assert.Equal(t, ids.PriceHistoryID, (*result)[0].ID)
		assert.Equal(t, ids.CommodityID, (*result)[0].CommodityID)
		assert.Equal(t, ids.RegionID, (*result)[0].RegionID)
		assert.Equal(t, float64(100), (*result)[0].Price)
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

func TestPriceHistoryRepository_FindByID(t *testing.T) {
	db, mock, repo, ids, rows, _ := PriceHistoryRepositorySetup(t)

	defer db.Close()

	expectedSQL := `SELECT * FROM "price_histories" WHERE "price_histories"."id" = $1 AND "price_histories"."deleted_at" IS NULL ORDER BY "price_histories"."id" LIMIT $2`

	t.Run("should return price history when find by id successfully", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.PriceHistoryID, 1).WillReturnRows(rows.PriceHistory)

		result, err := repo.FindByID(context.TODO(), ids.PriceHistoryID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.PriceHistoryID, result.ID)
		assert.Equal(t, ids.CommodityID, result.CommodityID)
		assert.Equal(t, ids.RegionID, result.RegionID)
		assert.Equal(t, float64(100), result.Price)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
	t.Run("should return error when find by id failed", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.PriceHistoryID, 1).WillReturnError(errors.New("database error"))

		result, err := repo.FindByID(context.TODO(), ids.PriceHistoryID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})
	t.Run("should return error when find by id not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.PriceHistoryID, 1).WillReturnRows(rows.Notfound)

		result, err := repo.FindByID(context.TODO(), ids.PriceHistoryID)
		assert.Nil(t, result)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestPriceHistoryRepository_FindByCommodityIDAndRegionID(t *testing.T) {
	db, mock, repo, ids, rows, _ := PriceHistoryRepositorySetup(t)

	defer db.Close()

	expectedSQL := `SELECT * FROM "price_histories" WHERE (commodity_id = $1 AND region_id = $2) AND "price_histories"."deleted_at" IS NULL`

	t.Run("should return price history when find by commodity id and region id successfully", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.CommodityID, ids.RegionID).WillReturnRows(rows.PriceHistory)

		result, err := repo.FindByCommodityIDAndRegionID(context.TODO(), ids.CommodityID, ids.RegionID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Len(t, *result, 1)
		assert.Equal(t, ids.PriceHistoryID, (*result)[0].ID)
		assert.Equal(t, ids.CommodityID, (*result)[0].CommodityID)
		assert.Equal(t, ids.RegionID, (*result)[0].RegionID)
		assert.Equal(t, float64(100), (*result)[0].Price)
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
