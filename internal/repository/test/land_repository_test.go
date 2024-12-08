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

func TestLandRepository_Create(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewLandRepository(db)

	expectedSQL := `INSERT INTO "lands" ("id","user_id","land_area","certificate","created_at","updated_at","deleted_at") VALUES ($1,$2,$3,$4,$5,$6,$7)`

	landID := uuid.New()
	userID := uuid.New()

	t.Run("Test Create, successfully", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(landID, userID, 100, "certificate", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		err := repoImpl.Create(context.TODO(), &domain.Land{
			ID:          landID,
			UserID:      userID,
			LandArea:    100,
			Certificate: "certificate",
		})

		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test Create, error database", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(landID, userID, 100, "certificate", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnError(errors.New("database error"))

		mock.ExpectRollback()

		err := repoImpl.Create(context.TODO(), &domain.Land{
			ID:          landID,
			UserID:      userID,
			LandArea:    100,
			Certificate: "certificate",
		})

		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestLandRepository_FindByID(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewLandRepository(db)

	expectedSQL := `SELECT * FROM "lands" WHERE "lands"."id" = $1 AND "lands"."deleted_at" IS NULL ORDER BY "lands"."id" LIMIT $2`

	landID := uuid.New()
	userID := uuid.New()

	t.Run("Test FindByID, successfully", func(t *testing.T) {
		land := sqlmock.NewRows([]string{"id", "user_id", "land_area", "certificate", "created_at", "updated_at"}).AddRow(landID, userID, int64(100), "certificate", time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
			WithArgs(landID, 1).
			WillReturnRows(land)

		result, err := repoImpl.FindByID(context.TODO(), landID)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, landID, result.ID)
		assert.Equal(t, userID, result.UserID)
		assert.Equal(t, int64(100), result.LandArea)
		assert.Equal(t, "certificate", result.Certificate)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test FindByID, error database", func(t *testing.T) {

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
			WithArgs(landID, 1).
			WillReturnError(errors.New("database error"))

		result, err := repoImpl.FindByID(context.TODO(), landID)

		assert.NotNil(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test FindByID, not found", func(t *testing.T) {
		land := sqlmock.NewRows([]string{"id", "user_id", "land_area", "certificate", "created_at", "updated_at", "deleted_at"})

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
			WithArgs(landID, 1).
			WillReturnRows(land)

		result, err := repoImpl.FindByID(context.TODO(), landID)

		assert.Nil(t, result)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestLandRepository_FindByUserID(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewLandRepository(db)

	expectedSQL := `SELECT * FROM "lands" WHERE user_id = $1 AND "lands"."deleted_at" IS NULL`

	userID := uuid.New()

	t.Run("Test FindByUserID, error database", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(userID).WillReturnError(errors.New("database error"))

		result, err := repoImpl.FindByUserID(context.TODO(), userID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")

		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test FindByUserID, successfully", func(t *testing.T) {
		land := sqlmock.NewRows([]string{"id", "user_id", "land_area", "certificate", "created_at", "updated_at"}).
			AddRow(uuid.New(), userID, int64(100), "certificate", time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(userID).WillReturnRows(land)

		result, err := repoImpl.FindByUserID(context.TODO(), userID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, userID, (*result)[0].UserID)
		assert.Equal(t, int64(100), (*result)[0].LandArea)
		assert.Equal(t, "certificate", (*result)[0].Certificate)

		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestLandRepository_FindAll(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewLandRepository(db)

	expectedSQL := `SELECT * FROM "lands"`

	landID := uuid.New()
	userID := uuid.New()

	t.Run("Test FindAll, successfully", func(t *testing.T) {
		lands := sqlmock.NewRows([]string{"id", "user_id", "land_area", "certificate", "created_at", "updated_at"}).AddRow(landID, userID, int64(100), "certificate", time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
			WillReturnRows(lands)

		result, err := repoImpl.FindAll(context.TODO())

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 1, len(*result))
		assert.Equal(t, landID, (*result)[0].ID)
		assert.Equal(t, userID, (*result)[0].UserID)
		assert.Equal(t, int64(100), (*result)[0].LandArea)
		assert.Equal(t, "certificate", (*result)[0].Certificate)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestLandRepository_Update(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewLandRepository(db)

	expectedSQL := `UPDATE "lands" SET "land_area"=$1,"certificate"=$2,"updated_at"=$3 WHERE id = $4 AND "lands"."deleted_at" IS NULL`

	landID := uuid.New()

	t.Run("Test Update, successfully", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(int64(100), "certificate", sqlmock.AnyArg(), landID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repoImpl.Update(context.TODO(), landID, &domain.Land{LandArea: 100, Certificate: "certificate"})
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test Update, error database", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(int64(100), "certificate", sqlmock.AnyArg(), landID).
			WillReturnError(errors.New("database error"))
		mock.ExpectRollback()

		err := repoImpl.Update(context.TODO(), landID, &domain.Land{LandArea: 100, Certificate: "certificate"})
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test Update, not found", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(int64(100), "certificate", sqlmock.AnyArg(), landID).
			WillReturnError(gorm.ErrRecordNotFound)
		mock.ExpectRollback()

		err := repoImpl.Update(context.TODO(), landID, &domain.Land{LandArea: 100, Certificate: "certificate"})
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestLandRepository_Delete(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewLandRepository(db)

	expectedSQL := `UPDATE "lands" SET "deleted_at"=$1 WHERE id = $2 AND "lands"."deleted_at" IS NULL`

	landID := uuid.New()

	t.Run("Test Delete, successfully", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), landID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repoImpl.Delete(context.TODO(), landID)
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test Delete, error database", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), landID).
			WillReturnError(errors.New("database error"))
		mock.ExpectRollback()

		err := repoImpl.Delete(context.TODO(), landID)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test Delete, not found", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), landID).
			WillReturnError(gorm.ErrRecordNotFound)
		mock.ExpectRollback()

		err := repoImpl.Delete(context.TODO(), landID)
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestLandRepository_Restore(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewLandRepository(db)

	expectedSQL := `UPDATE "lands" SET "deleted_at"=$1,"updated_at"=$2 WHERE id = $3`

	landID := uuid.New()

	t.Run("Test Restore, successfully", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), landID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repoImpl.Restore(context.TODO(), landID)
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test Restore, error database", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), landID).
			WillReturnError(errors.New("database error"))
		mock.ExpectRollback()

		err := repoImpl.Restore(context.TODO(), landID)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test Restore, not found", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), landID).
			WillReturnError(gorm.ErrRecordNotFound)
		mock.ExpectRollback()

		err := repoImpl.Restore(context.TODO(), landID)
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestLandRepository_FindDeletedByID(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewLandRepository(db)

	expectedSQL := `SELECT * FROM "lands" WHERE id = $1 AND deleted_at IS NOT NULL ORDER BY "lands"."id" LIMIT $2`

	landID := uuid.New()
	userID := uuid.New()

	t.Run("Test FindDeletedByID, error notfound", func(t *testing.T) {
		land := sqlmock.NewRows([]string{"id", "user_id", "land_area", "certificate", "created_at", "updated_at", "deleted_at"})

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(landID, 1).WillReturnRows(land)

		result, err := repoImpl.FindDeletedByID(context.TODO(), landID)
		assert.Nil(t, result)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))

		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test FindDeletedByID, error database", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(landID, 1).WillReturnError(errors.New("database error"))

		result, err := repoImpl.FindDeletedByID(context.TODO(), landID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")

		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test FindDeletedByID, successfully", func(t *testing.T) {
		land := sqlmock.NewRows([]string{"id", "user_id", "land_area", "certificate", "created_at", "updated_at"}).
			AddRow(landID, userID, int64(100), "certificate", time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(landID, 1).WillReturnRows(land)

		result, err := repoImpl.FindDeletedByID(context.TODO(), landID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, landID, result.ID)
		assert.Equal(t, userID, result.UserID)
		assert.Equal(t, int64(100), result.LandArea)
		assert.Equal(t, "certificate", result.Certificate)

		assert.Nil(t, mock.ExpectationsWereMet())
	})
}
