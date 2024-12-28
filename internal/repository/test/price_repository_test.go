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
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	repository_implementation "github.com/ryvasa/go-super-farmer/internal/repository/implementation"
	repository_interface "github.com/ryvasa/go-super-farmer/internal/repository/interface"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type PriceRepositoryIDs struct {
	PriceID     uuid.UUID
	CommodityID uuid.UUID
	CityID      int64
}

type PriceRepositoryMockRows struct {
	Price     *sqlmock.Rows
	Notfound  *sqlmock.Rows
	Commodity *sqlmock.Rows
	City      *sqlmock.Rows
}

type PriceRepositoryMocDomain struct {
	Price *domain.Price
}

type PrivcePaginationDTO struct {
	Pagination *dto.PaginationDTO
	Filter     *dto.ParamFilterDTO
}

func PriceRepositorySetup(t *testing.T) (*database.MockDB, repository_interface.PriceRepository, PriceRepositoryIDs, PriceRepositoryMockRows, PriceRepositoryMocDomain, PrivcePaginationDTO) {

	mockDB := database.NewMockDB(t)

	repo := repository_implementation.NewPriceRepository(mockDB.BaseRepo)

	priceID := uuid.New()
	commodityID := uuid.New()
	cityID := int64(1)

	ids := PriceRepositoryIDs{
		PriceID:     priceID,
		CommodityID: commodityID,
		CityID:      cityID,
	}

	rows := PriceRepositoryMockRows{
		Price: sqlmock.NewRows([]string{"id", "price", "commodity_id", "city_id", "created_at", "updated_at"}).
			AddRow(priceID, float64(100), commodityID, cityID, time.Now(), time.Now()),

		Notfound: sqlmock.NewRows([]string{"id", "price", "created_at", "updated_at"}),

		Commodity: sqlmock.NewRows([]string{"id", "name", "code", "duration"}).
			AddRow(commodityID, "commodity name", "commodity code", 3000),

		City: sqlmock.NewRows([]string{"id", "name"}).
			AddRow(cityID, "city name"),
	}

	domains := PriceRepositoryMocDomain{
		Price: &domain.Price{
			ID:          priceID,
			CommodityID: commodityID,
			CityID:      cityID,
			Price:       float64(100),
		},
	}

	dtos := PrivcePaginationDTO{
		Pagination: &dto.PaginationDTO{
			Page:  1,
			Limit: 10,
			Sort:  "created_at desc",
		},
		Filter: &dto.ParamFilterDTO{},
	}

	return mockDB, repo, ids, rows, domains, dtos
}

func TestPriceRepository_Create(t *testing.T) {
	mockDB, repo, ids, _, domains, _ := PriceRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL := `INSERT INTO "prices" ("id","commodity_id","city_id","price","unit","created_at","updated_at","deleted_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`

	t.Run("should not return error when create successfully", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.PriceID, ids.CommodityID, ids.CityID, float64(100), "idr", sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mockDB.Mock.ExpectCommit()

		err := repo.Create(context.TODO(), domains.Price)
		assert.Nil(t, err)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when create failed", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.PriceID, ids.CommodityID, ids.CityID, float64(100), "idr", sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
			WillReturnError(errors.New("database error"))
		mockDB.Mock.ExpectRollback()

		err := repo.Create(context.TODO(), domains.Price)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}

func TestPriceRepository_FindAll(t *testing.T) {
	mockDB, repo, ids, rows, _, dtos := PriceRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL1 := `SELECT * FROM "prices"`

	expectedSQL2 := `SELECT * FROM "commodities" WHERE "commodities"."id" = $1 AND "commodities"."deleted_at" IS NULL`

	expectedSQL3 := `SELECT * FROM "cities" WHERE "cities"."id" = $1`

	t.Run("should return price histories when find all successfully", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL1)).WillReturnRows(rows.Price)

		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL3)).WithArgs(ids.CityID).WillReturnRows(rows.City)

		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL2)).WithArgs(ids.CommodityID).WillReturnRows(rows.Commodity)

		result, err := repo.FindAll(context.TODO(), dtos.Pagination)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 1)
		assert.Equal(t, ids.PriceID, (result)[0].ID)
		assert.Equal(t, float64(100), (result)[0].Price)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
	t.Run("should return error when find all failed", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL1)).WillReturnError(errors.New("database error"))

		result, err := repo.FindAll(context.TODO(), dtos.Pagination)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}

func TestPriceRepository_FindByID(t *testing.T) {
	mockDB, repo, ids, rows, _, _ := PriceRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL1 := `SELECT * FROM "prices" WHERE "prices"."id" = $1 AND "prices"."deleted_at" IS NULL ORDER BY "prices"."id" LIMIT $2`

	expectedSQL2 := `SELECT "commodities"."id","commodities"."name","commodities"."code","commodities"."duration" FROM "commodities" WHERE "commodities"."id" = $1 AND "commodities"."deleted_at" IS NULL`

	expectedSQL3 := `SELECT * FROM "cities" WHERE "cities"."id" = $1`

	t.Run("should return price when find by id successfully", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL1)).WithArgs(ids.PriceID, 1).WillReturnRows(rows.Price)

		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL3)).WithArgs(ids.CityID).WillReturnRows(rows.City)

		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL2)).WithArgs(ids.CommodityID).WillReturnRows(rows.Commodity)

		result, err := repo.FindByID(context.TODO(), ids.PriceID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.PriceID, result.ID)
		assert.Equal(t, float64(100), result.Price)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
	t.Run("should return error when find by id failed", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL1)).WithArgs(ids.PriceID, 1).WillReturnError(errors.New("database error"))
		result, err := repo.FindByID(context.TODO(), ids.PriceID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
	t.Run("should return error when find by id not found", func(t *testing.T) {
		price := sqlmock.NewRows([]string{"id", "price", "created_at", "updated_at"})
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL1)).WithArgs(ids.PriceID, 1).WillReturnRows(price)
		result, err := repo.FindByID(context.TODO(), ids.PriceID)
		assert.Nil(t, result)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}

func TestPriceRepository_FindByCommodityID(t *testing.T) {
	mockDB, repo, ids, rows, _, _ := PriceRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL1 := `SELECT * FROM "prices" WHERE prices.commodity_id = $1 AND "prices"."deleted_at" IS NULL`

	expectedSQL2 := `SELECT "commodities"."id","commodities"."name","commodities"."code","commodities"."duration" FROM "commodities" WHERE "commodities"."id" = $1 AND "commodities"."deleted_at" IS NULL`

	expectedSQL3 := `SELECT * FROM "cities" WHERE "cities"."id" = $1`

	t.Run("should return prices when find by commodity id successfully", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL1)).WithArgs(ids.CommodityID).WillReturnRows(rows.Price)

		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL3)).WithArgs(ids.CityID).WillReturnRows(rows.City)

		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL2)).WithArgs(ids.CommodityID).WillReturnRows(rows.Commodity)

		result, err := repo.FindByCommodityID(context.TODO(), ids.CommodityID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.PriceID, (result)[0].ID)
		assert.Equal(t, ids.CommodityID, (result)[0].CommodityID)
		assert.Equal(t, ids.CityID, (result)[0].CityID)
		assert.Equal(t, float64(100), (result)[0].Price)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when find by commodity id failed", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL1)).WithArgs(ids.CommodityID).WillReturnError(errors.New("database error"))
		result, err := repo.FindByCommodityID(context.TODO(), ids.CommodityID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}

func TestPriceRepository_FindByCityID(t *testing.T) {
	mockDB, repo, ids, rows, _, _ := PriceRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL1 := `SELECT * FROM "prices" WHERE prices.city_id = $1 AND "prices"."deleted_at" IS NULL`

	expectedSQL3 := `SELECT * FROM "cities" WHERE "cities"."id" = $1`

	expectedSQL2 := `SELECT "commodities"."id","commodities"."name","commodities"."code","commodities"."duration" FROM "commodities" WHERE "commodities"."id" = $1 AND "commodities"."deleted_at" IS NULL`

	t.Run("should return prices when find by city id successfully", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL1)).WithArgs(ids.CityID).WillReturnRows(rows.Price)

		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL3)).WithArgs(ids.CityID).WillReturnRows(rows.City)

		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL2)).WithArgs(ids.CommodityID).WillReturnRows(rows.Commodity)

		result, err := repo.FindByCityID(context.TODO(), ids.CityID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.PriceID, (result)[0].ID)
		assert.Equal(t, ids.CommodityID, (result)[0].CommodityID)
		assert.Equal(t, ids.CityID, (result)[0].CityID)
		assert.Equal(t, float64(100), (result)[0].Price)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when find by city id failed", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL1)).WithArgs(ids.CityID).WillReturnError(errors.New("database error"))
		result, err := repo.FindByCityID(context.TODO(), ids.CityID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}

func TestPriceRepository_Update(t *testing.T) {
	mockDB, repo, ids, _, domains, _ := PriceRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL := `UPDATE "prices" SET "id"=$1,"commodity_id"=$2,"city_id"=$3,"price"=$4,"updated_at"=$5 WHERE id = $6 AND "prices"."deleted_at" IS NULL`

	t.Run("should not return error when update successfully", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.PriceID, ids.CommodityID, ids.CityID, float64(100), sqlmock.AnyArg(), ids.PriceID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mockDB.Mock.ExpectCommit()

		err := repo.Update(context.TODO(), ids.PriceID, domains.Price)
		assert.Nil(t, err)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when update failed", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.PriceID, ids.CommodityID, ids.CityID, float64(100), sqlmock.AnyArg(), ids.PriceID).
			WillReturnError(errors.New("database error"))
		mockDB.Mock.ExpectRollback()

		err := repo.Update(context.TODO(), ids.PriceID, domains.Price)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when update not found", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.PriceID, ids.CommodityID, ids.CityID, float64(100), sqlmock.AnyArg(), ids.PriceID).
			WillReturnError(gorm.ErrRecordNotFound)
		mockDB.Mock.ExpectRollback()

		err := repo.Update(context.TODO(), ids.PriceID, domains.Price)
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}

func TestPriceRepository_Delete(t *testing.T) {
	mockDB, repo, ids, _, _, _ := PriceRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL := `UPDATE "prices" SET "deleted_at"=$1 WHERE id = $2 AND "prices"."deleted_at" IS NULL`

	t.Run("should not return error when delete successfully", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), ids.PriceID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mockDB.Mock.ExpectCommit()

		err := repo.Delete(context.TODO(), ids.PriceID)
		assert.Nil(t, err)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when delete failed", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), ids.PriceID).
			WillReturnError(errors.New("database error"))
		mockDB.Mock.ExpectRollback()

		err := repo.Delete(context.TODO(), ids.PriceID)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when delete not found", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), ids.PriceID).
			WillReturnError(gorm.ErrRecordNotFound)
		mockDB.Mock.ExpectRollback()

		err := repo.Delete(context.TODO(), ids.PriceID)
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}

func TestPriceRepository_Restore(t *testing.T) {
	mockDB, repo, ids, _, _, _ := PriceRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL := `UPDATE "prices" SET "deleted_at"=$1,"updated_at"=$2 WHERE id = $3`

	t.Run("should not return error when restore successfully", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), ids.PriceID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mockDB.Mock.ExpectCommit()

		err := repo.Restore(context.TODO(), ids.PriceID)
		assert.Nil(t, err)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when restore failed", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), ids.PriceID).
			WillReturnError(errors.New("database error"))
		mockDB.Mock.ExpectRollback()

		err := repo.Restore(context.TODO(), ids.PriceID)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when restore not found", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), ids.PriceID).
			WillReturnError(gorm.ErrRecordNotFound)
		mockDB.Mock.ExpectRollback()

		err := repo.Restore(context.TODO(), ids.PriceID)
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}

func TestPriceRepository_FindDeletedByID(t *testing.T) {
	mockDB, repo, ids, rows, _, _ := PriceRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL1 := `SELECT * FROM "prices" WHERE prices.id = $1 AND prices.deleted_at IS NOT NULL ORDER BY "prices"."id" LIMIT $2`

	expectedSQL2 := `SELECT "commodities"."id","commodities"."name","commodities"."code","commodities"."duration" FROM "commodities" WHERE "commodities"."id" = $1`

	expectedSQL3 := `SELECT * FROM "cities" WHERE "cities"."id" = $1`

	t.Run("should return price when find deleted by id successfully", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL1)).WithArgs(ids.PriceID, 1).WillReturnRows(rows.Price)

		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL3)).WithArgs(ids.CityID).WillReturnRows(rows.City)

		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL2)).WithArgs(ids.CommodityID).WillReturnRows(rows.Commodity)

		result, err := repo.FindDeletedByID(context.TODO(), ids.PriceID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.PriceID, result.ID)
		assert.Equal(t, float64(100), result.Price)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
	t.Run("should return error when find deleted by id failed", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL1)).WithArgs(ids.PriceID, 1).WillReturnError(errors.New("database error"))
		result, err := repo.FindDeletedByID(context.TODO(), ids.PriceID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
	t.Run("should return error when find deleted by id not found", func(t *testing.T) {
		price := sqlmock.NewRows([]string{"id", "price", "created_at", "updated_at"})
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL1)).WithArgs(ids.PriceID, 1).WillReturnRows(price)
		result, err := repo.FindDeletedByID(context.TODO(), ids.PriceID)
		assert.Nil(t, result)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}

func TestPriceRepository_FindByCommodityIDAndCityID(t *testing.T) {
	mockDB, repo, ids, rows, _, _ := PriceRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL := `SELECT * FROM "prices" WHERE (prices.commodity_id = $1 AND prices.city_id = $2) AND "prices"."deleted_at" IS NULL ORDER BY "prices"."id" LIMIT $3`

	expectedSQL2 := `SELECT "commodities"."id","commodities"."name","commodities"."code","commodities"."duration" FROM "commodities" WHERE "commodities"."id" = $1`

	expectedSQL3 := `SELECT * FROM "cities" WHERE "cities"."id" = $1`

	t.Run("should return price when find by commodity id and city id successfully", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.CommodityID, ids.CityID, 1).WillReturnRows(rows.Price)

		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL3)).WithArgs(ids.CityID).WillReturnRows(rows.City)

		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL2)).WithArgs(ids.CommodityID).WillReturnRows(rows.Commodity)

		result, err := repo.FindByCommodityIDAndCityID(context.TODO(), ids.CommodityID, ids.CityID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.PriceID, result.ID)
		assert.Equal(t, ids.CommodityID, result.CommodityID)
		assert.Equal(t, ids.CityID, result.CityID)
		assert.Equal(t, float64(100), result.Price)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when find by commodity id and city id failed", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.CommodityID, ids.CityID, 1).WillReturnError(errors.New("database error"))
		result, err := repo.FindByCommodityIDAndCityID(context.TODO(), ids.CommodityID, ids.CityID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when find by commodity id and city id not found", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.CommodityID, ids.CityID, 1).WillReturnRows(rows.Notfound)
		result, err := repo.FindByCommodityIDAndCityID(context.TODO(), ids.CommodityID, ids.CityID)
		assert.Nil(t, result)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}

func TestPriceRepository_Count(t *testing.T) {
	mockDB, repo, _, _, _, dtos := PriceRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL := `SELECT count(*) FROM "prices" WHERE "prices"."deleted_at" IS NULL`

	t.Run("should return count when count successfully", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
		result, err := repo.Count(context.TODO(), dtos.Filter)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(1), result)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
	t.Run("should return error when count failed", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WillReturnError(errors.New("database error"))
		result, err := repo.Count(context.TODO(), dtos.Filter)
		assert.Equal(t, int64(0), result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}
