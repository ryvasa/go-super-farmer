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

func TestCommodityRepository_Create(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewCommodityRepository(db)

	commodityID := uuid.New()

	expectedSQL := `INSERT INTO "commodities" ("id","name","description","created_at","updated_at","deleted_at") VALUES ($1,$2,$3,$4,$5,$6)`

	t.Run("Test Create, successfully", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(commodityID, "commodity", "commodity description", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repoImpl.Create(context.TODO(), &domain.Commodity{ID: commodityID, Name: "commodity", Description: "commodity description"})
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test Create, error database", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(commodityID, "commodity", "commodity description", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnError(errors.New("database error"))
		mock.ExpectRollback()

		err := repoImpl.Create(context.TODO(), &domain.Commodity{ID: commodityID, Name: "commodity", Description: "commodity description"})
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestCommodityRepository_FindAll(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewCommodityRepository(db)

	expectedSQL := `SELECT * FROM "commodities"`
	t.Run("Test FindAll, successfully", func(t *testing.T) {
		commodities := sqlmock.NewRows([]string{"id", "name", "description", "created_at", "updated_at"}).AddRow(uuid.New(), "commodity", "commodity description", time.Now(), time.Now()).AddRow(uuid.New(), "commodity", "commodity description", time.Now(), time.Now())
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WillReturnRows(commodities)
		result, err := repoImpl.FindAll(context.TODO())
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Len(t, *result, 2)
		assert.Equal(t, "commodity", (*result)[0].Name)
		assert.Equal(t, "commodity description", (*result)[0].Description)
		assert.Equal(t, "commodity", (*result)[1].Name)
		assert.Equal(t, "commodity description", (*result)[1].Description)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test FindAll, error database", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WillReturnError(errors.New("database error"))
		result, err := repoImpl.FindAll(context.TODO())
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestCommodityRepository_FindByID(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewCommodityRepository(db)

	commodityID := uuid.New()

	expectedSQL := `SELECT * FROM "commodities" WHERE "commodities"."id" = $1 AND "commodities"."deleted_at" IS NULL ORDER BY "commodities"."id" LIMIT $2`
	t.Run("Test FindByID, successfully", func(t *testing.T) {
		commodity := sqlmock.NewRows([]string{"id", "name", "description", "created_at", "updated_at"}).AddRow(commodityID, "commodity", "commodity description", time.Now(), time.Now())
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
			WithArgs(commodityID, 1).
			WillReturnRows(commodity)
		result, err := repoImpl.FindByID(context.TODO(), commodityID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, commodityID, result.ID)
		assert.Equal(t, "commodity", result.Name)
		assert.Equal(t, "commodity description", result.Description)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test FindByID, error database", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
			WithArgs(commodityID, 1).
			WillReturnError(errors.New("database error"))
		result, err := repoImpl.FindByID(context.TODO(), commodityID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test FindByID, not found", func(t *testing.T) {
		commodity := sqlmock.NewRows([]string{"id", "name", "description", "created_at", "updated_at", "deleted_at"})
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
			WithArgs(commodityID, 1).
			WillReturnRows(commodity)
		result, err := repoImpl.FindByID(context.TODO(), commodityID)
		assert.Nil(t, result)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestCommodityRepository_Update(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewCommodityRepository(db)

	commodityID := uuid.New()
	expectedSQL := `UPDATE "commodities" SET "name"=$1,"description"=$2,"updated_at"=$3 WHERE id = $4 AND "commodities"."deleted_at" IS NULL`

	t.Run("Test Update, successfully", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs("commodity", "commodity description", sqlmock.AnyArg(), commodityID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repoImpl.Update(context.TODO(), commodityID, &domain.Commodity{Name: "commodity", Description: "commodity description"})
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test Update, error database", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs("commodity", "commodity description", sqlmock.AnyArg(), commodityID).
			WillReturnError(errors.New("database error"))
		mock.ExpectRollback()

		err := repoImpl.Update(context.TODO(), commodityID, &domain.Commodity{Name: "commodity", Description: "commodity description"})
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test Update, not found", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs("commodity", "commodity description", sqlmock.AnyArg(), commodityID).
			WillReturnError(gorm.ErrRecordNotFound)
		mock.ExpectRollback()

		err := repoImpl.Update(context.TODO(), commodityID, &domain.Commodity{Name: "commodity", Description: "commodity description"})
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestCommodityRepository_Delete(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewCommodityRepository(db)

	commodityID := uuid.New()

	expectedSQL := `UPDATE "commodities" SET "deleted_at"=$1 WHERE id = $2 AND "commodities"."deleted_at" IS NULL`

	t.Run("Test Delete, successfully", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), commodityID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repoImpl.Delete(context.TODO(), commodityID)
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test Delete, error database", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), commodityID).
			WillReturnError(errors.New("database error"))
		mock.ExpectRollback()

		err := repoImpl.Delete(context.TODO(), commodityID)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test Delete, not found", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), commodityID).
			WillReturnError(gorm.ErrRecordNotFound)
		mock.ExpectRollback()

		err := repoImpl.Delete(context.TODO(), commodityID)
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestCommodityRepository_Restore(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewCommodityRepository(db)

	commodityID := uuid.New()
	expectedSQL := `UPDATE "commodities" SET "deleted_at"=$1,"updated_at"=$2 WHERE id = $3`

	t.Run("Test Restore, successfully", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), commodityID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repoImpl.Restore(context.TODO(), commodityID)
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test Restore, error database", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), commodityID).
			WillReturnError(errors.New("database error"))
		mock.ExpectRollback()

		err := repoImpl.Restore(context.TODO(), commodityID)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test Restore, not found", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), commodityID).
			WillReturnError(gorm.ErrRecordNotFound)
		mock.ExpectRollback()

		err := repoImpl.Restore(context.TODO(), commodityID)
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestCommodityRepository_FindDeletedByID(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewCommodityRepository(db)

	commodityID := uuid.New()
	expectedSQL := `SELECT * FROM "commodities" WHERE id = $1 AND deleted_at IS NOT NULL ORDER BY "commodities"."id" LIMIT $2`

	t.Run("Test FindDeletedByID, successfully", func(t *testing.T) {
		commodity := sqlmock.NewRows([]string{"id", "name", "description", "created_at", "updated_at"}).AddRow(commodityID, "commodity", "commodity description", time.Now(), time.Now())
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
			WithArgs(commodityID, 1).
			WillReturnRows(commodity)
		result, err := repoImpl.FindDeletedByID(context.TODO(), commodityID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, commodityID, result.ID)
		assert.Equal(t, "commodity", result.Name)
		assert.Equal(t, "commodity description", result.Description)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test FindDeletedByID, error database", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
			WithArgs(commodityID, 1).
			WillReturnError(errors.New("database error"))
		result, err := repoImpl.FindDeletedByID(context.TODO(), commodityID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test FindDeletedByID, not found", func(t *testing.T) {
		commodity := sqlmock.NewRows([]string{"id", "name", "description", "created_at", "updated_at", "deleted_at"})
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
			WithArgs(commodityID, 1).
			WillReturnRows(commodity)
		result, err := repoImpl.FindDeletedByID(context.TODO(), commodityID)
		assert.Nil(t, result)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}
