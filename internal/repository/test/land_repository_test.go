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
	"gorm.io/gorm"
)

type LandIDs struct {
	LandID uuid.UUID
	UserID uuid.UUID
}

type LandMockRows struct {
	Land  *sqlmock.Rows
	Lands *sqlmock.Rows
}

type LandMocDomain struct {
	Land *domain.Land
}

func LandRepositorySetup(t *testing.T) (*database.MockDB, repository_interface.LandRepository, LandIDs, LandMockRows, LandMocDomain) {

	mockDB := database.NewMockDB(t)

	repo := repository_implementation.NewLandRepository(mockDB.DB)

	landID := uuid.New()
	userID := uuid.New()

	ids := LandIDs{
		LandID: landID,
		UserID: userID,
	}

	rows := LandMockRows{
		Land: sqlmock.NewRows([]string{"id", "user_id", "land_area", "certificate", "created_at", "updated_at"}).
			AddRow(landID, userID, float64(100), "certificate", time.Now(), time.Now()),

		Lands: sqlmock.NewRows([]string{"id", "user_id", "land_area", "certificate", "created_at", "updated_at"}).
			AddRow(landID, userID, float64(100), "certificate", time.Now(), time.Now()).
			AddRow(uuid.New(), userID, float64(100), "certificate", time.Now(), time.Now()),
	}

	domains := LandMocDomain{
		Land: &domain.Land{
			ID:          landID,
			UserID:      userID,
			LandArea:    float64(100),
			Certificate: "certificate",
		},
	}

	return mockDB, repo, ids, rows, domains
}

func TestLandRepository_Create(t *testing.T) {
	mockDB, repo, ids, _, domains := LandRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL := `INSERT INTO "lands" ("id","user_id","land_area","unit","certificate","created_at","updated_at","deleted_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`

	t.Run("should not return error when create successfully", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.LandID, ids.UserID, float64(100), "ha", "certificate", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mockDB.Mock.ExpectCommit()

		err := repo.Create(context.TODO(), domains.Land)
		assert.Nil(t, err)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when create failed", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.LandID, ids.UserID, float64(100), "ha", "certificate", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnError(errors.New("database error"))
		mockDB.Mock.ExpectRollback()

		err := repo.Create(context.TODO(), domains.Land)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}

func TestLandRepository_FindByID(t *testing.T) {
	mockDB, repo, ids, rows, _ := LandRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL := `SELECT * FROM "lands" WHERE "lands"."id" = $1 AND "lands"."deleted_at" IS NULL ORDER BY "lands"."id" LIMIT $2`

	t.Run("should return land when find by id successfully", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.LandID, 1).WillReturnRows(rows.Land)

		result, err := repo.FindByID(context.TODO(), ids.LandID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.LandID, result.ID)
		assert.Equal(t, ids.UserID, result.UserID)
		assert.Equal(t, float64(100), result.LandArea)
		assert.Equal(t, "certificate", result.Certificate)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when find by id failed", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.LandID, 1).WillReturnError(errors.New("database error"))

		result, err := repo.FindByID(context.TODO(), ids.LandID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when find by id not found", func(t *testing.T) {
		land := sqlmock.NewRows([]string{"id", "user_id", "land_area", "certificate", "created_at", "updated_at", "deleted_at"})
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.LandID, 1).WillReturnRows(land)

		result, err := repo.FindByID(context.TODO(), ids.LandID)
		assert.Nil(t, result)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))

		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}

func TestLandRepository_FindByUserID(t *testing.T) {
	mockDB, repo, ids, rows, _ := LandRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL := `SELECT * FROM "lands" WHERE user_id = $1 AND "lands"."deleted_at" IS NULL`

	t.Run("should return lands when find by user id successfully", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.UserID).WillReturnRows(rows.Lands)

		result, err := repo.FindByUserID(context.TODO(), ids.UserID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.LandID, (result)[0].ID)
		assert.Equal(t, ids.UserID, (result)[0].UserID)
		assert.Equal(t, float64(100), (result)[0].LandArea)
		assert.Equal(t, "certificate", (result)[0].Certificate)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when find by user id failed", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.UserID).WillReturnError(errors.New("database error"))

		result, err := repo.FindByUserID(context.TODO(), ids.UserID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

}

func TestLandRepository_FindAll(t *testing.T) {
	mockDB, repo, ids, rows, _ := LandRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL := `SELECT * FROM "lands"`

	t.Run("should return lands when find all successfully", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WillReturnRows(rows.Lands)

		result, err := repo.FindAll(context.TODO())
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 2, len(result))
		assert.Equal(t, ids.LandID, (result)[0].ID)
		assert.Equal(t, ids.UserID, (result)[0].UserID)
		assert.Equal(t, float64(100), (result)[0].LandArea)
		assert.Equal(t, "certificate", (result)[0].Certificate)
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

func TestLandRepository_Update(t *testing.T) {
	mockDB, repo, ids, _, domains := LandRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL := `UPDATE "lands" SET "id"=$1,"user_id"=$2,"land_area"=$3,"certificate"=$4,"updated_at"=$5 WHERE id = $6 AND "lands"."deleted_at" IS NULL`

	t.Run("should not return error when update successfully", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.LandID, ids.UserID, float64(100), "certificate", sqlmock.AnyArg(), ids.LandID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mockDB.Mock.ExpectCommit()

		err := repo.Update(context.TODO(), ids.LandID, domains.Land)
		assert.Nil(t, err)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when update failed", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.LandID, ids.UserID, float64(100), "certificate", sqlmock.AnyArg(), ids.LandID).
			WillReturnError(errors.New("database error"))
		mockDB.Mock.ExpectRollback()

		err := repo.Update(context.TODO(), ids.LandID, domains.Land)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when update not found", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.LandID, ids.UserID, float64(100), "certificate", sqlmock.AnyArg(), ids.LandID).
			WillReturnError(gorm.ErrRecordNotFound)
		mockDB.Mock.ExpectRollback()

		err := repo.Update(context.TODO(), ids.LandID, domains.Land)
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}

func TestLandRepository_Delete(t *testing.T) {
	mockDB, repo, ids, _, _ := LandRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL := `UPDATE "lands" SET "deleted_at"=$1 WHERE id = $2 AND "lands"."deleted_at" IS NULL`

	t.Run("should not return error when delete successfully", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), ids.LandID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mockDB.Mock.ExpectCommit()

		err := repo.Delete(context.TODO(), ids.LandID)
		assert.Nil(t, err)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when delete failed", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), ids.LandID).
			WillReturnError(errors.New("database error"))
		mockDB.Mock.ExpectRollback()

		err := repo.Delete(context.TODO(), ids.LandID)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when delete not found", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), ids.LandID).
			WillReturnError(gorm.ErrRecordNotFound)
		mockDB.Mock.ExpectRollback()

		err := repo.Delete(context.TODO(), ids.LandID)
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}

func TestLandRepository_Restore(t *testing.T) {
	mockDB, repo, ids, _, _ := LandRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL := `UPDATE "lands" SET "deleted_at"=$1,"updated_at"=$2 WHERE id = $3`

	t.Run("should not return error when restore successfully", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), ids.LandID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mockDB.Mock.ExpectCommit()

		err := repo.Restore(context.TODO(), ids.LandID)
		assert.Nil(t, err)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when restore failed", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), ids.LandID).
			WillReturnError(errors.New("database error"))
		mockDB.Mock.ExpectRollback()

		err := repo.Restore(context.TODO(), ids.LandID)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when restore not found", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), ids.LandID).
			WillReturnError(gorm.ErrRecordNotFound)
		mockDB.Mock.ExpectRollback()

		err := repo.Restore(context.TODO(), ids.LandID)
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}

func TestLandRepository_FindDeletedByID(t *testing.T) {
	mockDB, repo, ids, rows, _ := LandRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL := `SELECT * FROM "lands" WHERE id = $1 AND deleted_at IS NOT NULL ORDER BY "lands"."id" LIMIT $2`

	t.Run("should return land when find deleted by id successfully", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.LandID, 1).WillReturnRows(rows.Land)

		result, err := repo.FindDeletedByID(context.TODO(), ids.LandID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.LandID, result.ID)
		assert.Equal(t, ids.UserID, result.UserID)
		assert.Equal(t, float64(100), result.LandArea)
		assert.Equal(t, "certificate", result.Certificate)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when find deleted by id failed", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.LandID, 1).WillReturnError(errors.New("database error"))

		result, err := repo.FindDeletedByID(context.TODO(), ids.LandID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
	t.Run("should return error when find deleted by id not found", func(t *testing.T) {
		land := sqlmock.NewRows([]string{"id", "user_id", "land_area", "certificate", "created_at", "updated_at", "deleted_at"})
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.LandID, 1).WillReturnRows(land)

		result, err := repo.FindDeletedByID(context.TODO(), ids.LandID)
		assert.Nil(t, result)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}
