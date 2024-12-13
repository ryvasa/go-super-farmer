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

type RegionRepositoryIDs struct {
	RegionID   uuid.UUID
	ProvinceID int64
	CityID     int64
}

type RegionRepositoryMockRows struct {
	Region    *sqlmock.Rows
	Notfound  *sqlmock.Rows
	Commodity *sqlmock.Rows
	Province  *sqlmock.Rows
}

type RegionRepositoryMocDomain struct {
	Region *domain.Region
}

func RegionRepositorySetup(t *testing.T) (*sql.DB, sqlmock.Sqlmock, repository.RegionRepository, RegionRepositoryIDs, RegionRepositoryMockRows, RegionRepositoryMocDomain) {

	sqlDB, db, mock := database.DbMock(t)

	repo := repository.NewRegionRepository(db)

	regionID := uuid.New()
	provinceID := int64(1)
	cityID := int64(1)

	ids := RegionRepositoryIDs{
		RegionID:   regionID,
		ProvinceID: provinceID,
		CityID:     cityID,
	}

	rows := RegionRepositoryMockRows{
		Region: sqlmock.NewRows([]string{"id", "province_id", "city_id", "created_at", "updated_at", "deleted_at"}).
			AddRow(regionID, provinceID, cityID, time.Now(), time.Now(), nil),

		Notfound: sqlmock.NewRows([]string{"id", "province_id", "city_id", "created_at", "updated_at", "deleted_at"}),

		Commodity: sqlmock.NewRows([]string{"id", "name"}).
			AddRow(regionID, "commodity name"),

		Province: sqlmock.NewRows([]string{"id", "name"}).
			AddRow(regionID, "province name"),
	}

	domains := RegionRepositoryMocDomain{
		Region: &domain.Region{
			ID:         regionID,
			ProvinceID: provinceID,
			CityID:     cityID,
		},
	}

	return sqlDB, mock, repo, ids, rows, domains
}

func TestRegionRepository_Create(t *testing.T) {
	db, mock, repo, ids, _, domains := RegionRepositorySetup(t)

	defer db.Close()

	expectedSQL := `INSERT INTO "regions" ("id","province_id","city_id","created_at","updated_at","deleted_at") VALUES ($1,$2,$3,$4,$5,$6)`

	t.Run("should not return error when create successfully", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.RegionID, ids.ProvinceID, ids.CityID, sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Create(context.TODO(), domains.Region)
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when create failed", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.RegionID, ids.ProvinceID, ids.CityID, sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
			WillReturnError(errors.New("database error"))
		mock.ExpectRollback()

		err := repo.Create(context.TODO(), domains.Region)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestRegionRepository_FindByID(t *testing.T) {
	db, mock, repo, ids, rows, _ := RegionRepositorySetup(t)

	defer db.Close()

	expectedSQL := `SELECT * FROM "regions" WHERE "regions"."id" = $1 AND "regions"."deleted_at" IS NULL ORDER BY "regions"."id" LIMIT $2`

	t.Run("should return region when find by id successfully", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.RegionID, 1).WillReturnRows(rows.Region)
		result, err := repo.FindByID(context.TODO(), ids.RegionID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.RegionID, result.ID)
		assert.Equal(t, int64(1), result.ProvinceID)
		assert.Equal(t, int64(1), result.CityID)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
	t.Run("should return error when find by id failed", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.RegionID, 1).WillReturnError(errors.New("database error"))
		result, err := repo.FindByID(context.TODO(), ids.RegionID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})
	t.Run("should return error when find by id not found", func(t *testing.T) {
		region := sqlmock.NewRows([]string{"id", "province_id", "city_id", "created_at", "updated_at", "deleted_at"})

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.RegionID, 1).WillReturnRows(region)

		result, err := repo.FindByID(context.TODO(), ids.RegionID)

		assert.NotNil(t, err)
		assert.EqualError(t, err, "record not found")
		assert.Nil(t, result)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestRegionRepository_FindAll(t *testing.T) {
	db, mock, repo, ids, rows, _ := RegionRepositorySetup(t)

	defer db.Close()

	expectedSQL := `SELECT * FROM "regions" WHERE "regions"."deleted_at" IS NULL`

	t.Run("should return regions when find all successfully", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WillReturnRows(rows.Region)

		result, err := repo.FindAll(context.TODO())

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.RegionID, (*result)[0].ID)
		assert.Equal(t, ids.ProvinceID, (*result)[0].ProvinceID)
		assert.Equal(t, int64(1), (*result)[0].CityID)
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

func TestRegionRepository_FindByProvinceID(t *testing.T) {
	db, mock, repo, ids, rows, _ := RegionRepositorySetup(t)

	defer db.Close()

	expectedSQL := `SELECT * FROM "regions" WHERE province_id = $1 AND "regions"."deleted_at" IS NULL`

	t.Run("should return regions when find by province id successfully", func(t *testing.T) {

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.ProvinceID).WillReturnRows(rows.Region)

		result, err := repo.FindByProvinceID(context.TODO(), ids.ProvinceID)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.RegionID, (*result)[0].ID)
		assert.Equal(t, ids.ProvinceID, (*result)[0].ProvinceID)
		assert.Equal(t, int64(1), (*result)[0].CityID)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when find by province id failed", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.ProvinceID).WillReturnError(errors.New("database error"))
		result, err := repo.FindByProvinceID(context.TODO(), ids.ProvinceID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")

		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestRegionRepository_Update(t *testing.T) {
	db, mock, repo, ids, _, domains := RegionRepositorySetup(t)

	defer db.Close()

	expectedSQL := `UPDATE "regions" SET "id"=$1,"province_id"=$2,"city_id"=$3,"updated_at"=$4 WHERE id = $5 AND "regions"."deleted_at" IS NULL`

	t.Run("should not return error when update successfully", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.RegionID, ids.ProvinceID, ids.CityID, sqlmock.AnyArg(), ids.RegionID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Update(context.TODO(),
			ids.RegionID, domains.Region)
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when update failed", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.RegionID, ids.ProvinceID, ids.CityID, sqlmock.AnyArg(), ids.RegionID).
			WillReturnError(errors.New("database error"))
		mock.ExpectRollback()

		err := repo.Update(context.TODO(),
			ids.RegionID, domains.Region)
		assert.NotNil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when update not found", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.RegionID, ids.ProvinceID, ids.CityID, sqlmock.AnyArg(), ids.RegionID).
			WillReturnError(gorm.ErrRecordNotFound)
		mock.ExpectRollback()

		err := repo.Update(context.TODO(),
			ids.RegionID, domains.Region)
		assert.NotNil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestRegionRepository_Delete(t *testing.T) {
	db, mock, repo, ids, _, _ := RegionRepositorySetup(t)

	defer db.Close()

	expectedSQL := `UPDATE "regions" SET "deleted_at"=$1 WHERE "regions"."id" = $2 AND "regions"."deleted_at" IS NULL`

	t.Run("should not return error when delete successfully", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), ids.RegionID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		err := repo.Delete(context.TODO(), ids.RegionID)
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when delete failed", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), ids.RegionID).
			WillReturnError(errors.New("error"))

		mock.ExpectRollback()

		err := repo.Delete(context.TODO(), ids.RegionID)
		assert.NotNil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when delete not found", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), ids.RegionID).
			WillReturnError(gorm.ErrRecordNotFound)

		mock.ExpectRollback()

		err := repo.Delete(context.TODO(), ids.RegionID)
		assert.NotNil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestRegionRepository_Restore(t *testing.T) {
	db, mock, repo, ids, _, _ := RegionRepositorySetup(t)

	defer db.Close()

	expectedSQL := `UPDATE "regions" SET "deleted_at"=$1,"updated_at"=$2 WHERE id = $3`

	t.Run("should not return error when restore successfully", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), ids.RegionID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		err := repo.Restore(context.TODO(), ids.RegionID)
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when restore failed", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), ids.RegionID).
			WillReturnError(errors.New("error"))

		mock.ExpectRollback()

		err := repo.Restore(context.TODO(), ids.RegionID)
		assert.NotNil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when restore not found", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), ids.RegionID).
			WillReturnError(gorm.ErrRecordNotFound)

		mock.ExpectRollback()

		err := repo.Restore(context.TODO(), ids.RegionID)
		assert.NotNil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestRegionRepository_FindDeleted(t *testing.T) {
	db, mock, repo, ids, rows, _ := RegionRepositorySetup(t)

	defer db.Close()

	expectedSQL := `SELECT * FROM "regions" WHERE "regions"."id" = $1 ORDER BY "regions"."id" LIMIT $2`

	t.Run("should return region when find deleted by id successfully", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.RegionID, 1).WillReturnRows(rows.Region)

		result, err := repo.FindDeleted(context.TODO(), ids.RegionID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.RegionID, result.ID)
		assert.Equal(t, int64(1), result.ProvinceID)
		assert.Equal(t, int64(1), result.CityID)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
	t.Run("should return error when find deleted by id failed", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.RegionID, 1).WillReturnError(gorm.ErrRecordNotFound)
		result, err := repo.FindDeleted(context.TODO(), ids.RegionID)
		assert.Nil(t, result)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))

		assert.Nil(t, mock.ExpectationsWereMet())
	})
}
