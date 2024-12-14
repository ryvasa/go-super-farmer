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

type LandCommodityIDs struct {
	LandCommodityID uuid.UUID
	LandID          uuid.UUID
	CommodityID     uuid.UUID
}

type LandCommodityMockRows struct {
	LandCommodity   *sqlmock.Rows
	LandCommodities *sqlmock.Rows
	Notfound        *sqlmock.Rows
	Commodity       *sqlmock.Rows
	Land            *sqlmock.Rows
}

type LandCommodityMocDomain struct {
	LandCommodity *domain.LandCommodity
}

func LandCommodityRepositorySetup(t *testing.T) (*sql.DB, sqlmock.Sqlmock, repository.LandCommodityRepository, LandCommodityIDs, LandCommodityMockRows, LandCommodityMocDomain) {

	sqlDB, db, mock := database.DbMock(t)

	repo := repository.NewLandCommodityRepository(db)

	landCommodityID := uuid.New()
	landID := uuid.New()
	commodityID := uuid.New()
	userID := uuid.New()

	ids := LandCommodityIDs{
		LandCommodityID: landCommodityID,
		LandID:          landID,
		CommodityID:     commodityID,
	}

	rows := LandCommodityMockRows{
		LandCommodity: sqlmock.NewRows([]string{"id", "land_area", "commodity_id", "land_id", "created_at", "updated_at"}).
			AddRow(landCommodityID, float64(100), commodityID, landID, time.Now(), time.Now()),

		LandCommodities: sqlmock.NewRows([]string{"id", "land_area", "commodity_id", "land_id", "created_at", "updated_at"}).
			AddRow(landCommodityID, float64(100), commodityID, landID, time.Now(), time.Now()).
			AddRow(uuid.New(), float64(100), commodityID, landID, time.Now(), time.Now()),

		Notfound: sqlmock.NewRows([]string{"id", "land_area", "commodity_id", "land_id", "created_at", "updated_at"}),

		Commodity: sqlmock.NewRows([]string{"id", "name"}).
			AddRow(commodityID, "commodity name"),

		Land: sqlmock.NewRows([]string{"id", "user_id", "land_area", "certificate"}).
			AddRow(landID, userID, float64(100), "certificate"),
	}

	domains := LandCommodityMocDomain{
		LandCommodity: &domain.LandCommodity{
			ID:          landCommodityID,
			LandArea:    float64(100),
			LandID:      landID,
			CommodityID: commodityID,
		},
	}

	return sqlDB, mock, repo, ids, rows, domains
}

func TestLandCommodityRepository_Create(t *testing.T) {
	db, mock, repo, ids, _, domains := LandCommodityRepositorySetup(t)

	defer db.Close()

	expectedSQL := `INSERT INTO "land_commodities" ("id","land_area","unit","commodity_id","land_id","created_at","updated_at","deleted_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`

	t.Run("should not return error when create successfully", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.LandCommodityID, float64(100), "ha", ids.CommodityID, ids.LandID, sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Create(context.TODO(), domains.LandCommodity)
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when create failed", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.LandCommodityID, float64(100), "ha", ids.CommodityID, ids.LandID, sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
			WillReturnError(errors.New("database error"))
		mock.ExpectRollback()

		err := repo.Create(context.TODO(), domains.LandCommodity)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestLandCommodityRepository_FindByID(t *testing.T) {
	db, mock, repo, ids, rows, _ := LandCommodityRepositorySetup(t)

	defer db.Close()

	expectedSQL1 := `SELECT * FROM "land_commodities" WHERE "land_commodities"."id" = $1 AND "land_commodities"."deleted_at" IS NULL ORDER BY "land_commodities"."id" LIMIT $2`

	expectedSQL2 := `SELECT "commodities"."id","commodities"."name","commodities"."code","commodities"."duration" FROM "commodities" WHERE "commodities"."id" = $1 AND "commodities"."deleted_at" IS NULL`

	expectedSQL3 := `SELECT "lands"."id","lands"."user_id","lands"."land_area","lands"."unit","lands"."certificate" FROM "lands" WHERE "lands"."id" = $1 AND "lands"."deleted_at" IS NULL`

	t.Run("should return land commodity when find by id successfully", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL1)).WithArgs(ids.LandCommodityID, 1).WillReturnRows(rows.LandCommodity)

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL2)).WithArgs(ids.CommodityID).WillReturnRows(rows.Commodity)

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL3)).WithArgs(ids.LandID).WillReturnRows(rows.Land)

		result, err := repo.FindByID(context.TODO(), ids.LandCommodityID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.LandCommodityID, result.ID)
		assert.Equal(t, float64(100), result.LandArea)
		assert.Equal(t, ids.CommodityID, result.CommodityID)
		assert.Equal(t, ids.LandID, result.LandID)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
	t.Run("should return error when find by id failed", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL1)).WithArgs(ids.LandCommodityID, 1).WillReturnError(errors.New("database error"))

		result, err := repo.FindByID(context.TODO(), ids.LandCommodityID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})
	t.Run("should return error when find by id not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL1)).WithArgs(ids.LandCommodityID, 1).WillReturnRows(rows.Notfound)

		result, err := repo.FindByID(context.TODO(), ids.LandCommodityID)
		assert.Nil(t, result)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestLandCommodityRepository_FindAll(t *testing.T) {
	db, mock, repo, ids, rows, _ := LandCommodityRepositorySetup(t)

	defer db.Close()

	expectedSQL := `SELECT * FROM "land_commodities"`

	t.Run("should return land commodities when find all successfully", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WillReturnRows(rows.LandCommodities)

		result, err := repo.FindAll(context.TODO())
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 2, len(*result))
		assert.Equal(t, ids.LandCommodityID, (*result)[0].ID)
		assert.Equal(t, float64(100), (*result)[0].LandArea)
		assert.Equal(t, ids.CommodityID, (*result)[0].CommodityID)
		assert.Equal(t, ids.LandID, (*result)[0].LandID)
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

func TestLandCommodityRepository_FindByCommodityID(t *testing.T) {
	db, mock, repo, ids, rows, _ := LandCommodityRepositorySetup(t)

	defer db.Close()

	expectedSQL := `SELECT * FROM "land_commodities" WHERE commodity_id = $1 AND "land_commodities"."deleted_at" IS NULL`

	t.Run("should return land commodities when find by commodity id successfully", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.CommodityID).WillReturnRows(rows.LandCommodities)

		result, err := repo.FindByCommodityID(context.TODO(), ids.CommodityID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.LandCommodityID, (*result)[0].ID)
		assert.Equal(t, float64(100), (*result)[0].LandArea)
		assert.Equal(t, ids.CommodityID, (*result)[0].CommodityID)
		assert.Equal(t, ids.LandID, (*result)[0].LandID)
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

func TestLandCommodityRepository_FindByLandID(t *testing.T) {
	db, mock, repo, ids, rows, _ := LandCommodityRepositorySetup(t)

	defer db.Close()

	expectedSQL := `SELECT * FROM "land_commodities" WHERE land_id = $1 AND "land_commodities"."deleted_at" IS NULL`

	t.Run("should return land commodities when find by land id successfully", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.LandID).WillReturnRows(rows.LandCommodities)

		result, err := repo.FindByLandID(context.TODO(), ids.LandID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.LandCommodityID, (*result)[0].ID)
		assert.Equal(t, float64(100), (*result)[0].LandArea)
		assert.Equal(t, ids.CommodityID, (*result)[0].CommodityID)
		assert.Equal(t, ids.LandID, (*result)[0].LandID)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when find by land id failed", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.LandID).WillReturnError(errors.New("database error"))

		result, err := repo.FindByLandID(context.TODO(), ids.LandID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestLandCommodityRepository_Update(t *testing.T) {
	db, mock, repo, ids, _, domains := LandCommodityRepositorySetup(t)

	defer db.Close()

	expectedSQL := `UPDATE "land_commodities" SET "id"=$1,"land_area"=$2,"commodity_id"=$3,"land_id"=$4,"updated_at"=$5 WHERE id = $6 AND "land_commodities"."deleted_at" IS NULL`

	t.Run("should not return error when update successfully", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.LandCommodityID, float64(100), ids.CommodityID, ids.LandID, sqlmock.AnyArg(), ids.LandCommodityID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Update(context.TODO(), ids.LandCommodityID, domains.LandCommodity)
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when update failed", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.LandCommodityID, float64(100), ids.CommodityID, ids.LandID, sqlmock.AnyArg(), ids.LandCommodityID).
			WillReturnError(errors.New("database error"))
		mock.ExpectRollback()

		err := repo.Update(context.TODO(), ids.LandCommodityID, domains.LandCommodity)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when update not found", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.LandCommodityID, float64(100), ids.CommodityID, ids.LandID, sqlmock.AnyArg(), ids.LandCommodityID).
			WillReturnError(gorm.ErrRecordNotFound)
		mock.ExpectRollback()

		err := repo.Update(context.TODO(), ids.LandCommodityID, domains.LandCommodity)
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestLandCommodityRepository_Delete(t *testing.T) {
	db, mock, repo, ids, _, _ := LandCommodityRepositorySetup(t)

	defer db.Close()

	expectedSQL := `UPDATE "land_commodities" SET "deleted_at"=$1 WHERE id = $2 AND "land_commodities"."deleted_at" IS NULL`

	t.Run("should not return error when delete successfully", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), ids.LandCommodityID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Delete(context.TODO(), ids.LandCommodityID)
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when delete failed", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), ids.LandCommodityID).
			WillReturnError(errors.New("database error"))
		mock.ExpectRollback()

		err := repo.Delete(context.TODO(), ids.LandCommodityID)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when delete not found", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), ids.LandCommodityID).
			WillReturnError(gorm.ErrRecordNotFound)
		mock.ExpectRollback()

		err := repo.Delete(context.TODO(), ids.LandCommodityID)
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestLandCommodityRepository_Restore(t *testing.T) {
	db, mock, repo, ids, _, _ := LandCommodityRepositorySetup(t)

	defer db.Close()

	expectedSQL := `UPDATE "land_commodities" SET "deleted_at"=$1,"updated_at"=$2 WHERE id = $3`

	t.Run("should not return error when restore successfully", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), ids.LandCommodityID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Restore(context.TODO(), ids.LandCommodityID)
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when restore failed", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), ids.LandCommodityID).
			WillReturnError(errors.New("database error"))
		mock.ExpectRollback()

		err := repo.Restore(context.TODO(), ids.LandCommodityID)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when restore not found", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), ids.LandCommodityID).
			WillReturnError(gorm.ErrRecordNotFound)
		mock.ExpectRollback()

		err := repo.Restore(context.TODO(), ids.LandCommodityID)
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestLandCommodityRepository_FindDeletedByID(t *testing.T) {
	db, mock, repo, ids, rows, _ := LandCommodityRepositorySetup(t)

	defer db.Close()

	expectedSQL := `SELECT * FROM "land_commodities" WHERE id = $1 AND deleted_at IS NOT NULL ORDER BY "land_commodities"."id" LIMIT $2`

	t.Run("should return land commodity when find deleted by id successfully", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.LandCommodityID, 1).WillReturnRows(rows.LandCommodity)

		result, err := repo.FindDeletedByID(context.TODO(), ids.LandCommodityID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.LandCommodityID, result.ID)
		assert.Equal(t, float64(100), result.LandArea)
		assert.Equal(t, ids.CommodityID, result.CommodityID)
		assert.Equal(t, ids.LandID, result.LandID)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
	t.Run("should return error when find deleted by id failed", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.LandCommodityID, 1).WillReturnError(errors.New("database error"))

		result, err := repo.FindDeletedByID(context.TODO(), ids.LandCommodityID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})
	t.Run("should return error when find deleted by id not found", func(t *testing.T) {
		commodity := sqlmock.NewRows([]string{"id", "land_area", "commodity_id", "land_id", "created_at", "updated_at", "deleted_at"})
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.LandCommodityID, 1).WillReturnRows(commodity)

		result, err := repo.FindDeletedByID(context.TODO(), ids.LandCommodityID)
		assert.Nil(t, result)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}
