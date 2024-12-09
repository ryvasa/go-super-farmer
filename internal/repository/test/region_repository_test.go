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

func TestRegionRepository_Create(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewRegionRepository(db)

	expectedSQL := `INSERT INTO "regions" ("id","province_id","city_id","created_at","updated_at","deleted_at") VALUES ($1,$2,$3,$4,$5,$6)`

	regionID := uuid.New()
	provinceID := int64(1)
	cityID := int64(1)
	t.Run("Test Create, successfully", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(regionID, provinceID, cityID, sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repoImpl.Create(context.TODO(), &domain.Region{ID: regionID, ProvinceID: provinceID, CityID: cityID})

		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test Create, error database", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(regionID, provinceID, cityID, sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
			WillReturnError(errors.New("database error"))
		mock.ExpectRollback()

		err := repoImpl.Create(context.TODO(), &domain.Region{ID: regionID, ProvinceID: provinceID, CityID: cityID})

		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestRegionRepository_FindByID(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewRegionRepository(db)

	expectedSQL := `SELECT * FROM "regions" WHERE "regions"."id" = $1 AND "regions"."deleted_at" IS NULL ORDER BY "regions"."id" LIMIT $2`

	regionID := uuid.New()
	cityID := int64(1)
	provinceID := int64(1)
	t.Run("Test FindByID, successfully", func(t *testing.T) {
		region := sqlmock.NewRows([]string{"id", "province_id", "city_id", "created_at", "updated_at", "deleted_at"}).
			AddRow(regionID, cityID, provinceID, time.Now(), time.Now(), nil)

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
			WithArgs(regionID, 1).
			WillReturnRows(region)

		result, err := repoImpl.FindByID(context.TODO(), regionID)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(1), result.ProvinceID)
		assert.Equal(t, int64(1), result.CityID)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test FindByID, error database", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
			WithArgs(regionID, 1).
			WillReturnError(errors.New("database error"))

		result, err := repoImpl.FindByID(context.TODO(), regionID)

		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, result)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
	t.Run("Test FindByID, error not found", func(t *testing.T) {
		region := sqlmock.NewRows([]string{"id", "province_id", "city_id", "created_at", "updated_at", "deleted_at"})

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
			WithArgs(regionID, 1).
			WillReturnRows(region)

		result, err := repoImpl.FindByID(context.TODO(), regionID)

		assert.NotNil(t, err)
		assert.EqualError(t, err, "record not found")
		assert.Nil(t, result)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestRegionRepository_FindByProvinceID(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewRegionRepository(db)

	expectedSQL := `SELECT * FROM "regions" WHERE province_id = $1 AND "regions"."deleted_at" IS NULL`

	regionID := uuid.New()
	provinceID := int64(1)

	t.Run("Test FindByProvinceID, successfully", func(t *testing.T) {
		regions := sqlmock.NewRows([]string{"id", "province_id", "city_id", "created_at", "updated_at", "deleted_at"}).AddRow(regionID, provinceID, 1, time.Now(), time.Now(), nil)

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
			WithArgs(provinceID).
			WillReturnRows(regions)

		result, err := repoImpl.FindByProvinceID(context.TODO(), provinceID)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(1), (*result)[0].ProvinceID)
		assert.Equal(t, regionID, (*result)[0].ID)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test FindByProvinceID, database error", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
			WithArgs(provinceID).
			WillReturnError(errors.New("database error"))

		result, err := repoImpl.FindByProvinceID(context.TODO(), provinceID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")

		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestRegionRepository_FindAll(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewRegionRepository(db)

	expectedSQL := `SELECT * FROM "regions" WHERE "regions"."deleted_at" IS NULL`

	regionID := uuid.New()
	cityID := int64(1)
	provinceID := int64(1)
	t.Run("Test FindAll, successfully", func(t *testing.T) {
		regions := sqlmock.NewRows([]string{"id", "province_id", "city_id", "created_at", "updated_at", "deleted_at"}).
			AddRow(regionID, provinceID, cityID, time.Now(), time.Now(), nil)

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
			WillReturnRows(regions)

		result, err := repoImpl.FindAll(context.TODO())

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, provinceID, (*result)[0].ProvinceID)
		assert.Equal(t, regionID, (*result)[0].ID)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test FindAll, database error", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
			WillReturnError(errors.New("database error"))

		result, err := repoImpl.FindAll(context.TODO())
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")

		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestRegionRepository_Update(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewRegionRepository(db)

	expectedSQL := `UPDATE "regions" SET "province_id"=$1,"city_id"=$2,"updated_at"=$3 WHERE id = $4 AND "regions"."deleted_at" IS NULL`

	regionID := uuid.New()
	cityID := int64(1)
	provinceID := int64(1)
	t.Run("Test Update, successfully", func(t *testing.T) {

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(provinceID, cityID, sqlmock.AnyArg(), regionID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repoImpl.Update(context.TODO(),
			regionID, &domain.Region{
				ProvinceID: provinceID,
				CityID:     cityID,
			})
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test Update, database error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(provinceID, cityID, sqlmock.AnyArg(), regionID).
			WillReturnError(errors.New("database error"))
		mock.ExpectRollback()

		err := repoImpl.Update(context.TODO(),
			regionID, &domain.Region{
				ProvinceID: provinceID,
				CityID:     cityID,
			})
		assert.NotNil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test Update, record not found", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(provinceID, cityID, sqlmock.AnyArg(), regionID).
			WillReturnError(gorm.ErrRecordNotFound)
		mock.ExpectRollback()

		err := repoImpl.Update(context.TODO(),
			regionID, &domain.Region{
				ProvinceID: provinceID,
				CityID:     cityID,
			})
		assert.NotNil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestRegionRepository_Delete(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewRegionRepository(db)

	expectedSQL := `UPDATE "regions" SET "deleted_at"=$1 WHERE "regions"."id" = $2 AND "regions"."deleted_at" IS NULL`

	regionID := uuid.New()

	t.Run("Test Delete, successfully", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), regionID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		err := repoImpl.Delete(context.TODO(), regionID)
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test Delete, database error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), regionID).
			WillReturnError(errors.New("error"))

		mock.ExpectRollback()

		err := repoImpl.Delete(context.TODO(), regionID)
		assert.NotNil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test Delete, not found", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), regionID).
			WillReturnError(gorm.ErrRecordNotFound)

		mock.ExpectRollback()

		err := repoImpl.Delete(context.TODO(), regionID)
		assert.NotNil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestRegionRepository_Restore(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewRegionRepository(db)

	expectedSQL := `UPDATE "regions" SET "deleted_at"=$1,"updated_at"=$2 WHERE id = $3`

	regionID := uuid.New()

	t.Run("Test Restore, successfully", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), regionID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		err := repoImpl.Restore(context.TODO(), regionID)
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test Restore, database error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), regionID).
			WillReturnError(errors.New("error"))

		mock.ExpectRollback()

		err := repoImpl.Restore(context.TODO(), regionID)
		assert.NotNil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test Restore, not found", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), regionID).
			WillReturnError(gorm.ErrRecordNotFound)

		mock.ExpectRollback()

		err := repoImpl.Restore(context.TODO(), regionID)
		assert.NotNil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestRegionRepository_FindDeleted(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewRegionRepository(db)

	expectedSQL := `SELECT * FROM "regions" WHERE "regions"."id" = $1 ORDER BY "regions"."id" LIMIT $2`

	regionID := uuid.New()
	provinceID := int64(1)
	cityID := int64(1)
	t.Run("Test FindDeleted, error notfound", func(t *testing.T) {
		region := sqlmock.NewRows([]string{"id", "province_id", "city_id", "created_at", "updated_at", "deleted_at"})

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(regionID, 1).WillReturnRows(region)

		result, err := repoImpl.FindDeleted(context.TODO(), regionID)
		assert.Nil(t, result)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))

		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test FindDeleted, error database", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(regionID, 1).WillReturnError(errors.New("database error"))

		result, err := repoImpl.FindDeleted(context.TODO(), regionID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")

		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test FindDeleted, successfully", func(t *testing.T) {
		region := sqlmock.NewRows([]string{"id", "province_id", "city_id", "created_at", "updated_at", "deleted_at"}).AddRow(regionID, provinceID, cityID, time.Now(), time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(regionID, 1).WillReturnRows(region)

		result, err := repoImpl.FindDeleted(context.TODO(), regionID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, regionID, result.ID)

		assert.Nil(t, mock.ExpectationsWereMet())
	})
}
