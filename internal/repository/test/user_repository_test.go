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

func TestUserRepository_Create(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewUserRepository(db)

	expectedSQL := `INSERT INTO "users" ("id","name","email","password","role_id","phone","created_at","updated_at","deleted_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`

	userID := uuid.New()
	phone := "122233"

	t.Run("Test Create, successfully", func(t *testing.T) {

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(userID, "admin", "admin@example.com", "password", 1, &phone, sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repoImpl.Create(context.TODO(), &domain.User{ID: userID, Name: "admin", Email: "admin@example.com", Password: "password", RoleID: 1, Phone: &phone})
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test Create, error database", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(userID, "admin", "admin@example.com", "password", 1, &phone, sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
			WillReturnError(errors.New("database error"))
		mock.ExpectRollback()

		err := repoImpl.Create(context.TODO(), &domain.User{ID: userID, Name: "admin", Email: "admin@example.com", Password: "password", RoleID: 1, Phone: &phone})
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestUserRepository_FindAll(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewUserRepository(db)

	expectedSQL := `SELECT users.id,users.name,users.email,users.phone,users.created_at,users.updated_at FROM "users" WHERE "users"."deleted_at" IS NULL`

	t.Run("Test FindAll, error database", func(t *testing.T) {

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WillReturnError(errors.New("database error"))

		result, err := repoImpl.FindAll(context.TODO())
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")

		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test FindAll, successfully", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "email", "phone", "created_at", "updated_at"}).
			AddRow(uuid.New(), "admin", "admin@example.com", "122233", time.Now(), time.Now()).
			AddRow(uuid.New(), "user", "user@example.com", "12233", time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WillReturnRows(rows)

		result, err := repoImpl.FindAll(context.TODO())
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Len(t, *result, 2)
		assert.Equal(t, "admin@example.com", (*result)[0].Email)
		assert.Equal(t, "user@example.com", (*result)[1].Email)

		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestUserRepository_FindByID(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewUserRepository(db)

	expectedSQL := `SELECT users.id,users.name,users.email,users.phone,users.created_at,users.updated_at FROM "users" WHERE "users"."id" = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`

	userID := uuid.New()

	t.Run("Test FindByID, error notfound", func(t *testing.T) {
		user := sqlmock.NewRows([]string{"id", "name", "email", "phone", "created_at", "updated_at"})

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(userID, 1).WillReturnRows(user)

		result, err := repoImpl.FindByID(context.TODO(), userID)
		assert.Nil(t, result)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))

		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test FindByID, error database", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(userID, 1).WillReturnError(errors.New("database error"))

		result, err := repoImpl.FindByID(context.TODO(), userID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")

		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test FindByID, successfully", func(t *testing.T) {
		user := sqlmock.NewRows([]string{"id", "name", "email", "phone", "created_at", "updated_at"}).
			AddRow(userID, "admin", "admin@example.com", "122233", time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(userID, 1).WillReturnRows(user)

		result, err := repoImpl.FindByID(context.TODO(), userID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "admin", result.Name)

		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestUserRepository_Update(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewUserRepository(db)

	expectedSQL := `UPDATE "users" SET "name"=$1,"email"=$2,"password"=$3,"phone"=$4,"updated_at"=$5 WHERE id = $6 AND "users"."deleted_at" IS NULL`

	userID := uuid.New()
	phone := "122233"

	t.Run("Test Update, successfully", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs("admin", "admin@example.com", "password", &phone, sqlmock.AnyArg(), userID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repoImpl.Update(context.TODO(), userID, &domain.User{Name: "admin", Email: "admin@example.com", Password: "password", Phone: &phone})
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test Update, error database", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs("admin", "admin@example.com", "password", &phone, sqlmock.AnyArg(), userID).
			WillReturnError(errors.New("database error"))
		mock.ExpectRollback()

		err := repoImpl.Update(context.TODO(), userID, &domain.User{Name: "admin", Email: "admin@example.com", Password: "password", Phone: &phone})
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test Update, not found", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs("admin", "admin@example.com", "password", &phone, sqlmock.AnyArg(), userID).
			WillReturnError(gorm.ErrRecordNotFound)
		mock.ExpectRollback()

		err := repoImpl.Update(context.TODO(), userID, &domain.User{Name: "admin", Email: "admin@example.com", Password: "password", Phone: &phone})
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestUserRepository_Delete(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewUserRepository(db)

	expectedSQL := `UPDATE "users" SET "deleted_at"=$1 WHERE id = $2 AND "users"."deleted_at" IS NULL`

	userID := uuid.New()

	t.Run("Test Delete, successfully", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), userID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repoImpl.Delete(context.TODO(), userID)
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test Delete, error database", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), userID).
			WillReturnError(errors.New("database error"))
		mock.ExpectRollback()

		err := repoImpl.Delete(context.TODO(), userID)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test Delete, not found", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), userID).
			WillReturnError(gorm.ErrRecordNotFound)
		mock.ExpectRollback()

		err := repoImpl.Delete(context.TODO(), userID)
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestUserRepository_Restore(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewUserRepository(db)

	expectedSQL := `UPDATE "users" SET "deleted_at"=$1,"updated_at"=$2 WHERE id = $3`

	userID := uuid.New()

	t.Run("Test Restore, successfully", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), userID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repoImpl.Restore(context.TODO(), userID)
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test Restore, error database", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), userID).
			WillReturnError(errors.New("database error"))
		mock.ExpectRollback()

		err := repoImpl.Restore(context.TODO(), userID)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test Restore, not found", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), userID).
			WillReturnError(gorm.ErrRecordNotFound)
		mock.ExpectRollback()

		err := repoImpl.Restore(context.TODO(), userID)
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestUserRepository_FindDeletedByID(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewUserRepository(db)

	expectedSQL := `SELECT * FROM "users" WHERE id = $1 AND deleted_at IS NOT NULL ORDER BY "users"."id" LIMIT $2`

	userID := uuid.New()

	parsedTime, err := time.Parse(time.RFC3339, "2022-01-01T00:00:00Z")
	assert.Nil(t, err)

	t.Run("Test FindDeletedByID, error notfound", func(t *testing.T) {
		user := sqlmock.NewRows([]string{"id", "name", "email", "phone", "created_at", "updated_at"})

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(userID, 1).WillReturnRows(user)

		result, err := repoImpl.FindDeletedByID(context.TODO(), userID)
		assert.Nil(t, result)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))

		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test FindDeletedByID, error database", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(userID, 1).WillReturnError(errors.New("database error"))

		result, err := repoImpl.FindDeletedByID(context.TODO(), userID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")

		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test FindDeletedByID, successfully", func(t *testing.T) {
		user := sqlmock.NewRows([]string{"id", "name", "email", "phone", "created_at", "updated_at"}).
			AddRow(userID, "admin", "admin@example.com", "122233", parsedTime, parsedTime)

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(userID, 1).WillReturnRows(user)

		result, err := repoImpl.FindDeletedByID(context.TODO(), userID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "admin", result.Name)

		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestUserRepository_FindByEmail(t *testing.T) {
	sqlDB, db, mock := database.DbMock(t)
	defer sqlDB.Close()

	repoImpl := repository.NewUserRepository(db)

	expectedSQL := `SELECT * FROM "users" WHERE email = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`

	email := "admin@example.com"

	parsedTime, err := time.Parse(time.RFC3339, "2022-01-01T00:00:00Z")
	assert.Nil(t, err)

	t.Run("Test FindByEmail, error notfound", func(t *testing.T) {
		user := sqlmock.NewRows([]string{"id", "name", "email", "phone", "created_at", "updated_at"})

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(email, 1).WillReturnRows(user)

		result, err := repoImpl.FindByEmail(context.TODO(), email)
		assert.Nil(t, result)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))

		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test FindByEmail, error database", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(email, 1).WillReturnError(errors.New("database error"))

		result, err := repoImpl.FindByEmail(context.TODO(), email)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")

		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Test FindByEmail, successfully", func(t *testing.T) {
		user := sqlmock.NewRows([]string{"id", "name", "email", "phone", "created_at", "updated_at"}).
			AddRow(uuid.New(), "admin", "admin@example.com", "122233", parsedTime, parsedTime)

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(email, 1).WillReturnRows(user)

		result, err := repoImpl.FindByEmail(context.TODO(), email)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "admin", result.Name)

		assert.Nil(t, mock.ExpectationsWereMet())
	})
}
