package repository_test

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/pkg/database"
	"github.com/ryvasa/go-super-farmer/service_api/model/domain"
	"github.com/ryvasa/go-super-farmer/service_api/model/dto"
	repository_implementation "github.com/ryvasa/go-super-farmer/service_api/repository/implementation"
	repository_interface "github.com/ryvasa/go-super-farmer/service_api/repository/interface"
	"github.com/ryvasa/go-super-farmer/utils"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type SaleRepoIDs struct {
	SaleID      uuid.UUID
	CommodityID uuid.UUID
	CityID      int64
}

type SaleRepoRows struct {
	Sale      *sqlmock.Rows
	Notfound  *sqlmock.Rows
	Commodity *sqlmock.Rows
	City      *sqlmock.Rows
}

type SaleRepoDomain struct {
	Sale *domain.Sale
}

type SaleDTO struct {
	Pagination *dto.PaginationDTO
	Filter     *dto.ParamFilterDTO
	Sale       *domain.Sale
}

func SaleRepoSetup(t *testing.T) (*database.MockDB, repository_interface.SaleRepository, *SaleRepoIDs, *SaleRepoRows, *SaleRepoDomain, *SaleDTO) {
	mockDB := database.NewMockDB(t)

	repo := repository_implementation.NewSaleRepository(mockDB.BaseRepo)

	ids := &SaleRepoIDs{
		SaleID:      uuid.New(),
		CommodityID: uuid.New(),
		CityID:      1,
	}

	rows := &SaleRepoRows{
		Sale: sqlmock.NewRows([]string{"id", "commodity_id", "city_id", "quantity", "unit", "price", "sale_date", "created_at", "updated_at"}).
			AddRow(ids.SaleID, ids.CommodityID, ids.CityID,
				float64(1), "kg", float64(100), time.Now(), time.Now(), time.Now()),
		Notfound: sqlmock.NewRows([]string{"id", "commodity_id", "city_id", "quantity", "unit", "price", "sale_date", "created_at", "updated_at"}),
		Commodity: sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
			AddRow(ids.CommodityID, "test", time.Now(), time.Now()),
		City: sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
			AddRow(1, "test", time.Now(), time.Now()),
	}

	domain := &SaleRepoDomain{
		Sale: &domain.Sale{
			ID:          ids.SaleID,
			CommodityID: ids.CommodityID,
			CityID:      ids.CityID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Unit:        "kg",
			Quantity:    float64(1),
			Price:       float64(100),
			SaleDate:    time.Now(),
		},
	}

	dto := &SaleDTO{
		Pagination: &dto.PaginationDTO{
			Page:  1,
			Limit: 10,
			Sort:  "created_at desc",
		},
		Filter: &dto.ParamFilterDTO{},
	}

	return mockDB, repo, ids, rows, domain, dto
}

func TestSaleRepo_Create(t *testing.T) {
	mockDB, repo, ids, _, domains, _ := SaleRepoSetup(t)
	defer mockDB.SqlDB.Close()

	expectedSQL := `INSERT INTO "sales" ("id","city_id","commodity_id","quantity","unit","price","sale_date","created_at","updated_at","deleted_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`

	t.Run("should create a new sale successfully", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.SaleID, ids.CityID, ids.CommodityID, float64(1), "kg", float64(100), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mockDB.Mock.ExpectCommit()

		err := repo.Create(context.TODO(), domains.Sale)
		assert.Nil(t, err)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return an error if create fails", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.SaleID, ids.CityID, ids.CommodityID, float64(1), "kg", float64(100), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
			WillReturnError(utils.NewInternalError("internal error"))
		mockDB.Mock.ExpectRollback()

		err := repo.Create(context.TODO(), domains.Sale)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "internal error")
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}

func TestSaleRepo_FindAll(t *testing.T) {
	mockDB, repo, ids, rows, _, dtos := SaleRepoSetup(t)
	defer mockDB.SqlDB.Close()

	expectedSQL := `SELECT * FROM "sales" WHERE "sales"."deleted_at" IS NULL`

	t.Run("should return sales successfully", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
			WillReturnRows(rows.Sale)

		result, err := repo.FindAll(context.TODO(), dtos.Pagination)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 1, len(result))
		assert.Equal(t, ids.SaleID, result[0].ID)
		assert.Equal(t, ids.CommodityID, result[0].CommodityID)
		assert.Equal(t, ids.CityID, result[0].CityID)
		assert.Equal(t, float64(1), result[0].Quantity)
		assert.Equal(t, "kg", result[0].Unit)
		assert.Equal(t, float64(100), result[0].Price)
	})

	t.Run("should return an error if find all fails", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
			WillReturnError(utils.NewInternalError("internal error"))

		result, err := repo.FindAll(context.TODO(), dtos.Pagination)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "internal error")
		assert.Nil(t, result)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}

func TestSaleRepo_FindByID(t *testing.T) {
	mockDB, repo, ids, rows, _, _ := SaleRepoSetup(t)
	defer mockDB.SqlDB.Close()

	expectedSQL := `SELECT * FROM "sales" WHERE id = $1 AND "sales"."deleted_at" IS NULL ORDER BY "sales"."id" LIMIT $2`

	t.Run("should return a sale successfully", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.SaleID, 1).
			WillReturnRows(rows.Sale)

		result, err := repo.FindByID(context.TODO(), ids.SaleID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.SaleID, result.ID)
		assert.Equal(t, ids.CommodityID, result.CommodityID)
		assert.Equal(t, ids.CityID, result.CityID)
		assert.Equal(t, float64(1), result.Quantity)
		assert.Equal(t, "kg", result.Unit)
		assert.Equal(t, float64(100), result.Price)
	})

	t.Run("should return an error if find by id fails", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.SaleID, 1).
			WillReturnError(utils.NewInternalError("internal error"))

		result, err := repo.FindByID(context.TODO(), ids.SaleID)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "internal error")
		assert.Nil(t, result)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return not found error if not found", func(t *testing.T) {

		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.SaleID, 1).
			WillReturnRows(rows.Notfound)

		result, err := repo.FindByID(context.TODO(), ids.SaleID)
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, result)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}

func TestSaleRepo_FindByCommodityID(t *testing.T) {
	mockDB, repo, ids, rows, _, dtos := SaleRepoSetup(t)
	defer mockDB.SqlDB.Close()

	expectedSQL := `SELECT * FROM "sales" WHERE commodity_id = $1 AND "sales"."deleted_at" IS NULL ORDER BY created_at desc LIMIT $2`

	t.Run("should return sales successfully", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.CommodityID, dtos.Pagination.Limit).
			WillReturnRows(rows.Sale)

		result, err := repo.FindByCommodityID(context.TODO(), dtos.Pagination, ids.CommodityID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 1, len(result))
		assert.Equal(t, ids.SaleID, result[0].ID)
		assert.Equal(t, ids.CommodityID, result[0].CommodityID)
		assert.Equal(t, ids.CityID, result[0].CityID)
		assert.Equal(t, float64(1), result[0].Quantity)
		assert.Equal(t, "kg", result[0].Unit)
		assert.Equal(t, float64(100), result[0].Price)
	})

	t.Run("should return an error if find by commodity id fails", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.CommodityID, dtos.Pagination.Limit).
			WillReturnError(utils.NewInternalError("internal error"))

		result, err := repo.FindByCommodityID(context.TODO(), dtos.Pagination, ids.CommodityID)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "internal error")
		assert.Nil(t, result)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}

func TestSaleRepo_FindByCityID(t *testing.T) {
	mockDB, repo, ids, rows, _, dtos := SaleRepoSetup(t)
	defer mockDB.SqlDB.Close()

	expectedSQL := `SELECT * FROM "sales" WHERE city_id = $1 AND "sales"."deleted_at" IS NULL ORDER BY created_at desc LIMIT $2`

	t.Run("should return sales successfully", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.CityID, dtos.Pagination.Limit).
			WillReturnRows(rows.Sale)

		result, err := repo.FindByCityID(context.TODO(), dtos.Pagination, ids.CityID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 1, len(result))
		assert.Equal(t, ids.SaleID, result[0].ID)
		assert.Equal(t, ids.CommodityID, result[0].CommodityID)
		assert.Equal(t, ids.CityID, result[0].CityID)
		assert.Equal(t, float64(1), result[0].Quantity)
		assert.Equal(t, "kg", result[0].Unit)
		assert.Equal(t, float64(100), result[0].Price)
	})

	t.Run("should return an error if find by city id fails", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.CityID, dtos.Pagination.Limit).
			WillReturnError(utils.NewInternalError("internal error"))

		result, err := repo.FindByCityID(context.TODO(), dtos.Pagination, ids.CityID)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "internal error")
		assert.Nil(t, result)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}

func TestSaleRepo_Update(t *testing.T) {
	mockDB, repo, ids, _, domains, _ := SaleRepoSetup(t)
	defer mockDB.SqlDB.Close()

	expectedSQL := `UPDATE "sales" SET "id"=$1,"city_id"=$2,"commodity_id"=$3,"quantity"=$4,"unit"=$5,"price"=$6,"sale_date"=$7,"created_at"=$8,"updated_at"=$9 WHERE id = $10 AND "sales"."deleted_at" IS NULL`

	t.Run("should update a sale successfully", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.SaleID, ids.CityID, ids.CommodityID, float64(1), "kg", float64(100), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), ids.SaleID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mockDB.Mock.ExpectCommit()

		err := repo.Update(context.TODO(), ids.SaleID, domains.Sale)
		assert.Nil(t, err)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return an error if update fails", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.SaleID, ids.CityID, ids.CommodityID, float64(1), "kg", float64(100), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), ids.SaleID).
			WillReturnError(utils.NewInternalError("internal error"))
		mockDB.Mock.ExpectRollback()

		err := repo.Update(context.TODO(), ids.SaleID, domains.Sale)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "internal error")
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return an error if not found", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.SaleID, ids.CityID, ids.CommodityID, float64(1), "kg", float64(100), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), ids.SaleID).
			WillReturnError(gorm.ErrRecordNotFound)
		mockDB.Mock.ExpectRollback()

		err := repo.Update(context.TODO(), ids.SaleID, domains.Sale)
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}

func TestSaleRepo_Delete(t *testing.T) {
	mockDB, repo, ids, _, _, _ := SaleRepoSetup(t)
	defer mockDB.SqlDB.Close()

	expectedSQL := `UPDATE "sales" SET "deleted_at"=$1 WHERE "sales"."id" = $2 AND "sales"."deleted_at" IS NULL`

	t.Run("should delete a sale successfully", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), ids.SaleID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mockDB.Mock.ExpectCommit()

		err := repo.Delete(context.TODO(), ids.SaleID)
		assert.Nil(t, err)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return an error if delete fails", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), ids.SaleID).
			WillReturnError(utils.NewInternalError("internal error"))
		mockDB.Mock.ExpectRollback()

		err := repo.Delete(context.TODO(), ids.SaleID)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "internal error")
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return an error if not found", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), ids.SaleID).
			WillReturnError(gorm.ErrRecordNotFound)
		mockDB.Mock.ExpectRollback()

		err := repo.Delete(context.TODO(), ids.SaleID)
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}

func TestSaleRepo_Restore(t *testing.T) {
	mockDB, repo, ids, _, _, _ := SaleRepoSetup(t)
	defer mockDB.SqlDB.Close()

	expectedSQL := `UPDATE "sales" SET "deleted_at"=$1,"updated_at"=$2 WHERE id = $3`

	t.Run("should restore a sale successfully", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), ids.SaleID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mockDB.Mock.ExpectCommit()

		err := repo.Restore(context.TODO(), ids.SaleID)
		assert.Nil(t, err)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return an error if restore fails", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), ids.SaleID).
			WillReturnError(utils.NewInternalError("internal error"))
		mockDB.Mock.ExpectRollback()

		err := repo.Restore(context.TODO(), ids.SaleID)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "internal error")
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return an error if not found", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), ids.SaleID).
			WillReturnError(gorm.ErrRecordNotFound)
		mockDB.Mock.ExpectRollback()

		err := repo.Restore(context.TODO(), ids.SaleID)
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}

func TestSaleRepo_FindAllDeleted(t *testing.T) {
	mockDB, repo, ids, rows, _, dtos := SaleRepoSetup(t)
	defer mockDB.SqlDB.Close()

	expectedSQL := `SELECT * FROM "sales" WHERE deleted_at IS NOT NULL`

	t.Run("should return sales successfully", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
			WillReturnRows(rows.Sale)

		result, err := repo.FindAllDeleted(context.TODO(), dtos.Pagination)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 1, len(result))
		assert.Equal(t, ids.SaleID, result[0].ID)
		assert.Equal(t, ids.CommodityID, result[0].CommodityID)
		assert.Equal(t, ids.CityID, result[0].CityID)
		assert.Equal(t, float64(1), result[0].Quantity)
		assert.Equal(t, "kg", result[0].Unit)
		assert.Equal(t, float64(100), result[0].Price)
	})

	t.Run("should return an error if find all deleted fails", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
			WillReturnError(utils.NewInternalError("internal error"))

		result, err := repo.FindAllDeleted(context.TODO(), dtos.Pagination)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "internal error")
		assert.Nil(t, result)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}

func TestSaleRepo_FindDeletedByID(t *testing.T) {
	mockDB, repo, ids, rows, _, _ := SaleRepoSetup(t)
	defer mockDB.SqlDB.Close()

	expectedSQL := `SELECT * FROM "sales" WHERE id = $1 AND deleted_at IS NOT NULL ORDER BY "sales"."id" LIMIT $2`

	t.Run("should return a sale successfully", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.SaleID, 1).
			WillReturnRows(rows.Sale)

		result, err := repo.FindDeletedByID(context.TODO(), ids.SaleID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.SaleID, result.ID)
		assert.Equal(t, ids.CommodityID, result.CommodityID)
		assert.Equal(t, ids.CityID, result.CityID)
		assert.Equal(t, float64(1), result.Quantity)
		assert.Equal(t, "kg", result.Unit)
		assert.Equal(t, float64(100), result.Price)
	})

	t.Run("should return an error if find deleted by id fails", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.SaleID, 1).
			WillReturnError(utils.NewInternalError("internal error"))

		result, err := repo.FindDeletedByID(context.TODO(), ids.SaleID)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "internal error")
		assert.Nil(t, result)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when find deleted by id not found", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.SaleID, 1).
			WillReturnRows(rows.Notfound)

		result, err := repo.FindDeletedByID(context.TODO(), ids.SaleID)

		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, result)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}
