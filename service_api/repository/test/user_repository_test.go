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
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type UserRepositoryIDs struct {
	UserID uuid.UUID
}

type UserRepositoryMockRows struct {
	User     *sqlmock.Rows
	NotFound *sqlmock.Rows
}

type UserRepositoryMocDomain struct {
	User *domain.User
}

type UserDTO struct {
	Pagination *dto.PaginationDTO
	Filter     *dto.ParamFilterDTO
}

func UserRepositorySetup(t *testing.T) (*database.MockDB, repository_interface.UserRepository, UserRepositoryIDs, UserRepositoryMockRows, UserRepositoryMocDomain, UserDTO) {

	mockDB := database.NewMockDB(t)

	repo := repository_implementation.NewUserRepository(mockDB.DB)

	userID := uuid.New()

	ids := UserRepositoryIDs{
		UserID: userID,
	}

	rows := UserRepositoryMockRows{
		User: sqlmock.NewRows([]string{"id", "name", "email", "phone", "password", "created_at", "updated_at", "deleted_at"}).
			AddRow(userID, "user name", "user@email.com", "123456789", "password", time.Now(), time.Now(), nil),

		NotFound: sqlmock.NewRows([]string{"id", "name", "email", "phone", "password", "created_at", "updated_at", "deleted_at"}),
	}

	phone := "123456789"
	domains := UserRepositoryMocDomain{
		User: &domain.User{
			ID:        userID,
			Name:      "user name",
			Email:     "user@email.com",
			Phone:     &phone,
			Password:  "password",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	dtos := UserDTO{
		Pagination: &dto.PaginationDTO{
			Page:  1,
			Limit: 10,
			Sort:  "created_at desc",
		},
		Filter: &dto.ParamFilterDTO{},
	}

	return mockDB, repo, ids, rows, domains, dtos
}

func TestUserRepository_Create(t *testing.T) {
	mockDB, repo, ids, _, domains, _ := UserRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL := `INSERT INTO "users" ("id","name","email","password","role_id","phone","created_at","updated_at","deleted_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`

	t.Run("should not return error when create successfully", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.UserID, "user name", "user@email.com", "password", 1, "123456789", sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mockDB.Mock.ExpectCommit()

		err := repo.Create(context.TODO(), domains.User)
		assert.Nil(t, err)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when create failed", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.UserID, "user name", "user@email.com", "password", 1, "123456789", sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
			WillReturnError(errors.New("database error"))
		mockDB.Mock.ExpectRollback()

		err := repo.Create(context.TODO(), domains.User)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}

func TestUserRepository_FindByID(t *testing.T) {
	mockDB, repo, ids, rows, _, _ := UserRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL := `SELECT users.id,users.name,users.email,users.phone,users.created_at,users.updated_at FROM "users" WHERE "users"."id" = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`

	t.Run("should return user when find by id successfully", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.UserID, 1).WillReturnRows(rows.User)
		result, err := repo.FindByID(context.TODO(), ids.UserID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.UserID, result.ID)
		assert.Equal(t, "user name", result.Name)
		assert.Equal(t, "user@email.com", result.Email)
		assert.Equal(t, "123456789", *result.Phone)
		assert.Equal(t, "password", result.Password)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
	t.Run("should return error when find by id failed", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.UserID, 1).WillReturnError(errors.New("database error"))
		result, err := repo.FindByID(context.TODO(), ids.UserID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
	t.Run("should return error when find by id not found", func(t *testing.T) {
		user := sqlmock.NewRows([]string{"id", "name", "email", "phone", "password", "created_at", "updated_at", "deleted_at"})
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.UserID, 1).WillReturnRows(user)
		result, err := repo.FindByID(context.TODO(), ids.UserID)
		assert.Nil(t, result)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}

func TestUserRepository_FindAll(t *testing.T) {
	mockDB, repo, ids, rows, _, dtos := UserRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL := `SELECT users.id,users.name,users.email,users.phone,users.created_at,users.updated_at FROM "users" WHERE "users"."deleted_at" IS NULL`

	t.Run("should return users when find all successfully", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WillReturnRows(rows.User)

		result, err := repo.FindAll(context.TODO(), dtos.Pagination)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.UserID, (result)[0].ID)
		assert.Equal(t, "user name", (result)[0].Name)
		assert.Equal(t, "user@email.com", (result)[0].Email)
		assert.Equal(t, "123456789", *(result)[0].Phone)
		assert.Equal(t, "password", (result)[0].Password)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when find all failed", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WillReturnError(errors.New("database error"))

		result, err := repo.FindAll(context.TODO(), dtos.Pagination)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}

func TestUserRepository_Update(t *testing.T) {
	mockDB, repo, ids, _, domains, _ := UserRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL := `UPDATE "users" SET "id"=$1,"name"=$2,"email"=$3,"password"=$4,"phone"=$5,"created_at"=$6,"updated_at"=$7 WHERE id = $8 AND "users"."deleted_at" IS NULL`

	t.Run("should not return error when update successfully", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.UserID, "user name", "user@email.com", "password", "123456789", sqlmock.AnyArg(), sqlmock.AnyArg(), ids.UserID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mockDB.Mock.ExpectCommit()

		err := repo.Update(context.TODO(), ids.UserID, domains.User)
		assert.Nil(t, err)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when update failed", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.UserID, "user name", "user@email.com", "password", "123456789", sqlmock.AnyArg(), sqlmock.AnyArg(), ids.UserID).
			WillReturnError(errors.New("database error"))
		mockDB.Mock.ExpectRollback()

		err := repo.Update(context.TODO(), ids.UserID, domains.User)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when update not found", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(ids.UserID, "user name", "user@email.com", "password", "123456789", sqlmock.AnyArg(), sqlmock.AnyArg(), ids.UserID).
			WillReturnError(gorm.ErrRecordNotFound)
		mockDB.Mock.ExpectRollback()

		err := repo.Update(context.TODO(), ids.UserID, domains.User)
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}

func TestUserRepository_Delete(t *testing.T) {
	mockDB, repo, ids, _, _, _ := UserRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL := `UPDATE "users" SET "deleted_at"=$1 WHERE id = $2 AND "users"."deleted_at" IS NULL`

	t.Run("should not return error when delete successfully", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), ids.UserID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mockDB.Mock.ExpectCommit()

		err := repo.Delete(context.TODO(), ids.UserID)
		assert.Nil(t, err)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when delete failed", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), ids.UserID).
			WillReturnError(errors.New("error"))

		mockDB.Mock.ExpectRollback()

		err := repo.Delete(context.TODO(), ids.UserID)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "error")
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when delete not found", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), ids.UserID).
			WillReturnError(gorm.ErrRecordNotFound)

		mockDB.Mock.ExpectRollback()

		err := repo.Delete(context.TODO(), ids.UserID)
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}

func TestUserRepository_Restore(t *testing.T) {
	mockDB, repo, ids, _, _, _ := UserRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL := `UPDATE "users" SET "deleted_at"=$1,"updated_at"=$2 WHERE id = $3`

	t.Run("should not return error when restore successfully", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), ids.UserID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mockDB.Mock.ExpectCommit()

		err := repo.Restore(context.TODO(), ids.UserID)
		assert.Nil(t, err)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when restore failed", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), ids.UserID).
			WillReturnError(errors.New("error"))

		mockDB.Mock.ExpectRollback()

		err := repo.Restore(context.TODO(), ids.UserID)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "error")
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when restore not found", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), ids.UserID).
			WillReturnError(gorm.ErrRecordNotFound)

		mockDB.Mock.ExpectRollback()

		err := repo.Restore(context.TODO(), ids.UserID)
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}

func TestUserRepository_FindDeleted(t *testing.T) {
	mockDB, repo, ids, rows, _, _ := UserRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL := `SELECT * FROM "users" WHERE id = $1 AND deleted_at IS NOT NULL ORDER BY "users"."id" LIMIT $2`

	t.Run("should return user when find deleted by id successfully", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.UserID, 1).WillReturnRows(rows.User)

		result, err := repo.FindDeletedByID(context.TODO(), ids.UserID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.UserID, result.ID)
		assert.Equal(t, "user name", result.Name)
		assert.Equal(t, "user@email.com", result.Email)
		assert.Equal(t, "123456789", *result.Phone)
		assert.Equal(t, "password", result.Password)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
	t.Run("should return error when find deleted by id failed", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.UserID, 1).WillReturnError(errors.New("database error"))
		result, err := repo.FindDeletedByID(context.TODO(), ids.UserID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
	t.Run("should return error when find deleted by id not found", func(t *testing.T) {
		user := sqlmock.NewRows([]string{"id", "name", "email", "phone", "password", "created_at", "updated_at", "deleted_at"})
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.UserID, 1).WillReturnRows(user)
		result, err := repo.FindDeletedByID(context.TODO(), ids.UserID)
		assert.Nil(t, result)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}

func TestUserRepository_FindByEmail(t *testing.T) {
	mockDB, repo, ids, rows, _, _ := UserRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL := `SELECT * FROM "users" WHERE email = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`

	t.Run("should return user when find by email successfully", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs("user@email.com", 1).WillReturnRows(rows.User)
		result, err := repo.FindByEmail(context.TODO(), "user@email.com")
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.UserID, result.ID)
		assert.Equal(t, "user name", result.Name)
		assert.Equal(t, "user@email.com", result.Email)
		assert.Equal(t, "123456789", *result.Phone)
		assert.Equal(t, "password", result.Password)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
	// t.Run("should return error when find by email failed", func(t *testing.T) {
	// 	mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs("user@email.com", 1).WillReturnError(errors.New("database error"))
	// 	result, err := repo.FindByEmail(context.TODO(), "user@email.com")
	// 	assert.Nil(t, result)
	// 	assert.NotNil(t, err)
	// 	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
	// 	assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	// })
	t.Run("should return error when find by email not found", func(t *testing.T) {
		user := sqlmock.NewRows([]string{"id", "name", "email", "phone", "password", "created_at", "updated_at", "deleted_at"})
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs("user@email.com", 1).WillReturnRows(user)
		result, err := repo.FindByEmail(context.TODO(), "user@email.com")
		assert.Nil(t, result.Phone)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}
func TestUserRepository_Count(t *testing.T) {
	mockDB, repo, _, _, _, dtos := UserRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL := `SELECT count(*) FROM "users" WHERE "users"."deleted_at" IS NULL`

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
