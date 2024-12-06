package repository_test

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/repository"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUserRepository_Create(t *testing.T) {
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewUserRepository(db)

	t.Run("Test Create, error database", func(t *testing.T) {
		expectedSQL := `INSERT INTO "users" ("created_at","updated_at","deleted_at","name","email","password","role_id","phone") VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING "id","id"`
		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "UserTest", "", "", 1, sqlmock.AnyArg()).
			WillReturnError(errors.New("database error"))
		mock.ExpectRollback()

		err := repoImpl.Create(context.TODO(), &domain.User{Name: "UserTest"})
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test Create, successfully", func(t *testing.T) {
		expectedSQL := `INSERT INTO "users" ("created_at","updated_at","deleted_at","name","email","password","role_id","phone") VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING "id","id"`
		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "UserTest", "test@example.com", "string", 1, sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectCommit()

		err := repoImpl.Create(context.TODO(), &domain.User{Name: "UserTest", Email: "test@example.com", Password: "string"})
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestUserRepository_FindById(t *testing.T) {
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewUserRepository(db)

	t.Run("Test FindByID, error notfound", func(t *testing.T) {
		user := sqlmock.NewRows([]string{"id", "name", "email", "phone", "created_at", "updated_at"}) // Empty rows

		expectedSQL := `SELECT users.id,users.name,users.email,users.phone,users.created_at,users.updated_at FROM "users" WHERE "users"."id" = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(2, 1).WillReturnRows(user)

		result, err := repoImpl.FindByID(context.TODO(), 2)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)

		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test FindByID, error database", func(t *testing.T) {
		expectedSQL := `SELECT users.id,users.name,users.email,users.phone,users.created_at,users.updated_at FROM "users" WHERE "users"."id" = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(1, 1).WillReturnError(errors.New("database error"))

		result, err := repoImpl.FindByID(context.TODO(), 1)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")

		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test FindByID, successfully", func(t *testing.T) {
		parsedTime, err := time.Parse("2006-01-02 15:04:05", "2024-12-05 10:00:00")
		if err != nil {
			fmt.Println("Error parsing time:", err)
			return
		}

		user := sqlmock.NewRows([]string{"id", "name", "email", "phone", "created_at", "updated_at"}).
			AddRow(1, "admin", "admin@example.com", "123456789", parsedTime, parsedTime)

		expectedSQL := `SELECT users.id,users.name,users.email,users.phone,users.created_at,users.updated_at FROM "users" WHERE "users"."id" = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(1, 1).WillReturnRows(user)

		result, err := repoImpl.FindByID(context.TODO(), 1)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "admin", result.Name)
		assert.Equal(t, "admin@example.com", result.Email)

		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestUserRepository_FindAll(t *testing.T) {
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewUserRepository(db)
	t.Run("Test FindAll, error database", func(t *testing.T) {
		expectedSQL := `SELECT users.id,users.name,users.email,users.phone,users.created_at,users.updated_at FROM "users" WHERE "users"."deleted_at" IS NULL`
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WillReturnError(errors.New("database error"))

		result, err := repoImpl.FindAll(context.TODO())
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")

		assert.Nil(t, mock.ExpectationsWereMet())
	})
	t.Run("Test FindAll, successfully", func(t *testing.T) {
		parsedTime, err := time.Parse("2006-01-02 15:04:05", "2024-12-05 10:00:00")
		if err != nil {
			fmt.Println("Error parsing time:", err)
			return
		}

		rows := sqlmock.NewRows([]string{"id", "name", "email", "phone", "created_at", "updated_at"}).
			AddRow(1, "ryan", "ryan@example.com", "00000000", parsedTime, parsedTime).
			AddRow(2, "oktavian", "oktavian@example.com", "11111111", parsedTime, parsedTime)

		expectedSQL := `SELECT users.id,users.name,users.email,users.phone,users.created_at,users.updated_at FROM "users" WHERE "users"."deleted_at" IS NULL`
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WillReturnRows(rows)

		result, err := repoImpl.FindAll(context.TODO())
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Len(t, *result, 2)
		assert.Equal(t, "ryan", (*result)[0].Name)
		assert.Equal(t, "oktavian", (*result)[1].Name)
		assert.Equal(t, "ryan@example.com", (*result)[0].Email)
		assert.Equal(t, "oktavian@example.com", (*result)[1].Email)
		assert.Equal(t, "00000000", *(*result)[0].Phone)
		assert.Equal(t, "11111111", *(*result)[1].Phone)
		assert.Equal(t, parsedTime, (*result)[0].CreatedAt)
		assert.Equal(t, parsedTime, (*result)[1].CreatedAt)
		assert.Equal(t, parsedTime, (*result)[0].UpdatedAt)
		assert.Equal(t, parsedTime, (*result)[1].UpdatedAt)

		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestUserRepository_Delete(t *testing.T) {
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewUserRepository(db)

	t.Run("Test Delete, error database", func(t *testing.T) {
		expectedSQL := `UPDATE "users" SET "deleted_at"=$1 WHERE id = $2`
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), 1).
			WillReturnError(errors.New("database error"))
		mock.ExpectRollback()

		err := repoImpl.Delete(context.TODO(), 1)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test Delete, successfully", func(t *testing.T) {
		expectedSQL := `UPDATE "users" SET "deleted_at"=$1 WHERE id = $2`
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), 1).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := repoImpl.Delete(context.TODO(), 1)
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestUserRepository_Restore(t *testing.T) {
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewUserRepository(db)

	t.Run("Test Restore, error database", func(t *testing.T) {
		expectedSQL := `UPDATE "users" SET "deleted_at"=$1,"updated_at"=$2 WHERE id = $3`
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), 1).
			WillReturnError(errors.New("database error"))
		mock.ExpectRollback()

		err := repoImpl.Restore(context.TODO(), 1)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test Restore, successfully", func(t *testing.T) {

		expectedSQL := `UPDATE "users" SET "deleted_at"=$1,"updated_at"=$2 WHERE id = $3`
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(nil, sqlmock.AnyArg(), 1).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := repoImpl.Restore(context.TODO(), 1)
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestUserRepository_Update(t *testing.T) {
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewUserRepository(db)

	t.Run("Test Update, error database", func(t *testing.T) {
		expectedSQL := `UPDATE "users" SET "id"=$1,"updated_at"=$2,"name"=$3,"email"=$4,"phone"=$5 WHERE id = $6 AND "users"."deleted_at" IS NULL`
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(1, sqlmock.AnyArg(), "new_name", "new_email@example.com", "987654321", 1).
			WillReturnError(errors.New("database error"))
		mock.ExpectRollback()

		phone := "987654321"
		user := &domain.User{
			ID:    1,
			Name:  "new_name",
			Email: "new_email@example.com",
			Phone: &phone,
		}
		err := repoImpl.Update(context.TODO(), user.ID, user)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test Update, successfully", func(t *testing.T) {
		expectedSQL := `UPDATE "users" SET "id"=$1,"updated_at"=$2,"name"=$3,"email"=$4,"phone"=$5 WHERE id = $6 AND "users"."deleted_at" IS NULL`

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(1, sqlmock.AnyArg(), "new_name", "new_email@example.com", "987654321", 1).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()
		phone := "987654321"
		user := &domain.User{
			ID:    1,
			Name:  "new_name",
			Email: "new_email@example.com",
			Phone: &phone,
		}
		err := repoImpl.Update(context.TODO(), user.ID, user)
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestUserRepository_FindDeletedByID(t *testing.T) {
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewUserRepository(db)
	t.Run("Test FindDeletedByID, error database", func(t *testing.T) {
		user := sqlmock.NewRows([]string{"id", "name", "email", "phone", "created_at", "updated_at"})

		expectedSQL := `SELECT * FROM "users" WHERE id = $1 AND deleted_at IS NOT NULL ORDER BY "users"."id" LIMIT $2`
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(2, 1).WillReturnRows(user)

		result, err := repoImpl.FindDeletedByID(context.TODO(), 2)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)

		assert.Nil(t, mock.ExpectationsWereMet())
	})
	t.Run("Test FindDeletedByID, error not found", func(t *testing.T) {
		user := sqlmock.NewRows([]string{"id", "name", "email", "phone", "created_at", "updated_at", "deleted_at"})

		expectedSQL := `SELECT * FROM "users" WHERE id = $1 AND deleted_at IS NOT NULL ORDER BY "users"."id" LIMIT $2`
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(1, 1).WillReturnRows(user)

		result, err := repoImpl.FindDeletedByID(context.TODO(), 1)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test FindDeletedByID, successfully", func(t *testing.T) {
		parsedTime, err := time.Parse("2006-01-02 15:04:05", "2024-12-05 10:00:00")
		if err != nil {
			fmt.Println("Error parsing time:", err)
			return
		}
		user := sqlmock.NewRows([]string{"id", "deleted_at"}).
			AddRow(1, parsedTime)

		expectedSQL := `SELECT * FROM "users" WHERE id = $1 AND deleted_at IS NOT NULL ORDER BY "users"."id" LIMIT $2`
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(1, 1).WillReturnRows(user)
		result, err := repoImpl.FindDeletedByID(context.TODO(), 1)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, uint64(1), result.ID)
		assert.Equal(t, parsedTime, result.DeletedAt.Time)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}
