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

func TestLandCommodityRepository_Create(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewLandCommodityRepository(db)

	expectedSQL := `INSERT INTO "land_commodities" ("id","land_area","commodity_id","land_id","created_at","updated_at","deleted_at") VALUES ($1,$2,$3,$4,$5,$6,$7)`

	landCommodityID := uuid.New()
	landID := uuid.New()
	commodityID := uuid.New()

	t.Run("Test Create, successfully", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(landCommodityID, float64(100), commodityID, landID, sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		err := repoImpl.Create(context.TODO(), &domain.LandCommodity{
			ID:          landCommodityID,
			LandArea:    float64(100),
			LandID:      landID,
			CommodityID: commodityID,
		})

		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test Create, error database", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(landCommodityID, float64(100), commodityID, landID, sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
			WillReturnError(errors.New("database error"))

		mock.ExpectRollback()

		err := repoImpl.Create(context.TODO(), &domain.LandCommodity{
			ID:          landCommodityID,
			LandArea:    float64(100),
			LandID:      landID,
			CommodityID: commodityID,
		})

		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestLandCommodityRepository_FindByID(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewLandCommodityRepository(db)

	expectedSQL := `SELECT * FROM "land_commodities" WHERE "land_commodities"."id" = $1 AND "land_commodities"."deleted_at" IS NULL ORDER BY "land_commodities"."id" LIMIT $2`

	expectedSQL2 := `SELECT "commodities"."id","commodities"."name" FROM "commodities" WHERE "commodities"."id" = $1 AND "commodities"."deleted_at" IS NULL`

	expectedSQL3 := `SELECT "lands"."id","lands"."user_id","lands"."land_area","lands"."certificate" FROM "lands" WHERE "lands"."id" = $1 AND "lands"."deleted_at" IS NULL`

	landCommodityID := uuid.New()
	landID := uuid.New()
	commodityID := uuid.New()
	userID := uuid.New()

	t.Run("Test FindByID, successfully", func(t *testing.T) {
		mock.MatchExpectationsInOrder(true)
		landCommoditie := sqlmock.NewRows([]string{"id", "land_area", "commodity_id", "land_id", "created_at", "updated_at"}).
			AddRow(landCommodityID, float64(100), commodityID, landID, time.Now(), time.Now())

		commodityMock := sqlmock.NewRows([]string{"id", "name", "description", "created_at", "updated_at"}).AddRow(commodityID, "commodity", "commodity description", time.Now(), time.Now())

		landMock := sqlmock.NewRows([]string{"id", "user_id", "land_area", "certificate", "created_at", "updated_at"}).AddRow(landID, userID, float64(100), "certificate", time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
			WithArgs(landCommodityID, 1).
			WillReturnRows(landCommoditie)

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL2)).WithArgs(commodityID).WillReturnRows(commodityMock)

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL3)).WithArgs(landID).WillReturnRows(landMock)

		result, err := repoImpl.FindByID(context.TODO(), landCommodityID)

		// Assertions
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, landCommodityID, result.ID)
		assert.Equal(t, float64(100), result.LandArea)
		assert.Equal(t, commodityID, result.CommodityID)
		assert.Equal(t, landID, result.LandID)

		// Verify mock expectations
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test FindByID, error database", func(t *testing.T) {

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
			WithArgs(landCommodityID, 1).
			WillReturnError(errors.New("database error"))

		result, err := repoImpl.FindByID(context.TODO(), landCommodityID)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")

		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test FindByID, not found", func(t *testing.T) {
		commodity := sqlmock.NewRows([]string{"id", "land_area", "commodity_id", "land_id", "created_at", "updated_at", "deleted_at"})

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
			WithArgs(landCommodityID, 1).
			WillReturnRows(commodity)

		result, err := repoImpl.FindByID(context.TODO(), landCommodityID)

		assert.Nil(t, result)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))

		assert.Nil(t, mock.ExpectationsWereMet())
	})

}

func TestLandCommodityRepository_FindByLandID(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewLandCommodityRepository(db)

	expectedSQL := `SELECT * FROM "land_commodities" WHERE land_id = $1 AND "land_commodities"."deleted_at" IS NULL`

	landID := uuid.New()
	commodityID := uuid.New()

	t.Run("Test FindByLandID, error database", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(landID).WillReturnError(errors.New("database error"))

		result, err := repoImpl.FindByLandID(context.TODO(), landID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")

		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test FindByLandID, successfully", func(t *testing.T) {
		commodity := sqlmock.NewRows([]string{"id", "land_area", "commodity_id", "land_id", "created_at", "updated_at"}).
			AddRow(uuid.New(), float64(100), commodityID, landID, time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(landID).WillReturnRows(commodity)

		result, err := repoImpl.FindByLandID(context.TODO(), landID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, landID, (*result)[0].LandID)
		assert.Equal(t, float64(100), (*result)[0].LandArea)
		assert.Equal(t, commodityID, (*result)[0].CommodityID)

		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestLandCommodityRepository_FindAll(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewLandCommodityRepository(db)

	expectedSQL := `SELECT * FROM "land_commodities"`

	commodityID := uuid.New()
	landID := uuid.New()

	t.Run("Test FindAll, successfully", func(t *testing.T) {
		commodities := sqlmock.NewRows([]string{"id", "land_area", "commodity_id", "land_id", "created_at", "updated_at"}).AddRow(commodityID, float64(100), commodityID, landID, time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
			WillReturnRows(commodities)

		result, err := repoImpl.FindAll(context.TODO())

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 1, len(*result))
		assert.Equal(t, commodityID, (*result)[0].ID)
		assert.Equal(t, float64(100), (*result)[0].LandArea)
		assert.Equal(t, commodityID, (*result)[0].CommodityID)
		assert.Equal(t, landID, (*result)[0].LandID)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test FindAll, error database", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
			WillReturnError(errors.New("database error"))

		result, err := repoImpl.FindAll(context.TODO())

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")

		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestLandCommodityRepository_FindByCommodityID(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewLandCommodityRepository(db)

	expectedSQL := `SELECT * FROM "land_commodities" WHERE commodity_id = $1 AND "land_commodities"."deleted_at" IS NULL`

	commodityID := uuid.New()
	landID := uuid.New()

	t.Run("Test FindByCommodityID, error database", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(commodityID).WillReturnError(errors.New("database error"))

		result, err := repoImpl.FindByCommodityID(context.TODO(), commodityID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")

		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test FindByCommodityID, successfully", func(t *testing.T) {
		commodity := sqlmock.NewRows([]string{"id", "land_area", "commodity_id", "land_id", "created_at", "updated_at"}).
			AddRow(uuid.New(), float64(100), commodityID, landID, time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(commodityID).WillReturnRows(commodity)

		result, err := repoImpl.FindByCommodityID(context.TODO(), commodityID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, commodityID, (*result)[0].CommodityID)
		assert.Equal(t, float64(100), (*result)[0].LandArea)
		assert.Equal(t, commodityID, (*result)[0].CommodityID)
		assert.Equal(t, landID, (*result)[0].LandID)

		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestLandCommodityRepository_Update(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewLandCommodityRepository(db)

	expectedSQL := `UPDATE "land_commodities" SET "land_area"=$1,"commodity_id"=$2,"land_id"=$3,"updated_at"=$4 WHERE id = $5 AND "land_commodities"."deleted_at" IS NULL`

	landCommodityID := uuid.New()
	commodityID := uuid.New()
	landID := uuid.New()

	t.Run("Test Update, successfully", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(float64(100), commodityID, landID, sqlmock.AnyArg(), landCommodityID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repoImpl.Update(context.TODO(), landCommodityID, &domain.LandCommodity{LandArea: float64(100), CommodityID: commodityID, LandID: landID})
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test Update, error database", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(float64(100), commodityID, landID, sqlmock.AnyArg(), landCommodityID).
			WillReturnError(errors.New("database error"))
		mock.ExpectRollback()

		err := repoImpl.Update(context.TODO(), landCommodityID, &domain.LandCommodity{LandArea: float64(100), CommodityID: commodityID, LandID: landID})
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test Update, not found", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(float64(100), commodityID, landID, sqlmock.AnyArg(), landCommodityID).
			WillReturnError(gorm.ErrRecordNotFound)
		mock.ExpectRollback()

		err := repoImpl.Update(context.TODO(), landCommodityID, &domain.LandCommodity{LandArea: float64(100), CommodityID: commodityID, LandID: landID})
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestLandCommodityRepository_Delete(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewLandCommodityRepository(db)

	expectedSQL := `UPDATE "land_commodities" SET "deleted_at"=$1 WHERE id = $2 AND "land_commodities"."deleted_at" IS NULL`

	landCommodityID := uuid.New()

	t.Run("Test Delete, successfully", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), landCommodityID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repoImpl.Delete(context.TODO(), landCommodityID)
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test Delete, error database", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), landCommodityID).
			WillReturnError(errors.New("database error"))
		mock.ExpectRollback()

		err := repoImpl.Delete(context.TODO(), landCommodityID)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test Delete, not found", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), landCommodityID).
			WillReturnError(gorm.ErrRecordNotFound)
		mock.ExpectRollback()

		err := repoImpl.Delete(context.TODO(), landCommodityID)
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestLandCommodityRepository_Restore(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewLandCommodityRepository(db)

	expectedSQL := `UPDATE "land_commodities" SET "deleted_at"=$1,"updated_at"=$2 WHERE id = $3`

	landCommodityID := uuid.New()

	t.Run("Test Restore, successfully", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), landCommodityID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repoImpl.Restore(context.TODO(), landCommodityID)
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test Restore, error database", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), landCommodityID).
			WillReturnError(errors.New("database error"))
		mock.ExpectRollback()

		err := repoImpl.Restore(context.TODO(), landCommodityID)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test Restore, not found", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), landCommodityID).
			WillReturnError(gorm.ErrRecordNotFound)
		mock.ExpectRollback()

		err := repoImpl.Restore(context.TODO(), landCommodityID)
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestLandCommodityRepository_FindDeletedByID(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewLandCommodityRepository(db)

	expectedSQL := `SELECT * FROM "land_commodities" WHERE id = $1 AND deleted_at IS NOT NULL ORDER BY "land_commodities"."id" LIMIT $2`

	landCommodityID := uuid.New()
	commodityID := uuid.New()
	landID := uuid.New()

	t.Run("Test FindDeletedByID, error notfound", func(t *testing.T) {
		commodity := sqlmock.NewRows([]string{"id", "land_area", "commodity_id", "land_id", "created_at", "updated_at", "deleted_at"})

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(landCommodityID, 1).WillReturnRows(commodity)

		result, err := repoImpl.FindDeletedByID(context.TODO(), landCommodityID)
		assert.Nil(t, result)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))

		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test FindDeletedByID, error database", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(landCommodityID, 1).WillReturnError(errors.New("database error"))

		result, err := repoImpl.FindDeletedByID(context.TODO(), landCommodityID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")

		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test FindDeletedByID, successfully", func(t *testing.T) {
		commodity := sqlmock.NewRows([]string{"id", "land_area", "commodity_id", "land_id", "created_at", "updated_at"}).
			AddRow(landCommodityID, float64(100), commodityID, landID, time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(landCommodityID, 1).WillReturnRows(commodity)

		result, err := repoImpl.FindDeletedByID(context.TODO(), landCommodityID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, landCommodityID, result.ID)
		assert.Equal(t, float64(100), result.LandArea)
		assert.Equal(t, commodityID, result.CommodityID)
		assert.Equal(t, landID, result.LandID)

		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestLandCommodityRepository_SumLandAreaByLandID(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewLandCommodityRepository(db)

	landID := uuid.New()

	expectedSQL := `SELECT COALESCE(SUM(land_area), 0) FROM "land_commodities" WHERE land_id = $1 AND "land_commodities"."deleted_at" IS NULL`
	t.Run("Test SumLandAreaByLandID, error database", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(landID).WillReturnError(errors.New("database error"))

		result, err := repoImpl.SumLandAreaByLandID(context.TODO(), landID)
		assert.Equal(t, float64(0), result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")

		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test SumLandAreaByLandID, successfully", func(t *testing.T) {
		commodity := sqlmock.NewRows([]string{"land_area"}).
			AddRow(float64(100))

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(landID).WillReturnRows(commodity)

		result, err := repoImpl.SumLandAreaByLandID(context.TODO(), landID)
		assert.Equal(t, float64(100), result)
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestLandCommodityRepository_SumLandAreaByCommodityID(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewLandCommodityRepository(db)

	commodityID := uuid.New()

	expectedSQL := `SELECT SUM(land_area) FROM "land_commodities" WHERE commodity_id = $1`
	t.Run("Test SumLandAreaByCommodityID, error database", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(commodityID).WillReturnError(errors.New("database error"))

		result, err := repoImpl.SumLandAreaByCommodityID(context.TODO(), commodityID)
		assert.Equal(t, float64(0), result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")

		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test SumLandAreaByCommodityID, successfully", func(t *testing.T) {
		commodity := sqlmock.NewRows([]string{"land_area"}).
			AddRow(float64(100))

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(commodityID).WillReturnRows(commodity)

		result, err := repoImpl.SumLandAreaByCommodityID(context.TODO(), commodityID)
		assert.Equal(t, float64(100), result)
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}
