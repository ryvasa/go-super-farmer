package repository_test

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ryvasa/go-super-farmer/service_api/model/domain"
	repository_implementation "github.com/ryvasa/go-super-farmer/service_api/repository/implementation"
	repository_interface "github.com/ryvasa/go-super-farmer/service_api/repository/interface"
	"github.com/ryvasa/go-super-farmer/pkg/database"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type RoleIDs struct {
	RoleID int64
}

type RoleMockRows struct {
	Role     *sqlmock.Rows
	Notfound *sqlmock.Rows
}

type RoleMockDomains struct {
	Role *domain.Role
}

func RoleRepositorySetup(t *testing.T) (*database.MockDB, repository_interface.RoleRepository, RoleIDs, RoleMockRows, RoleMockDomains) {

	mockDB := database.NewMockDB(t)

	repo := repository_implementation.NewRoleRepository(mockDB.DB)

	roleID := int64(1)

	ids := RoleIDs{RoleID: roleID}

	rows := RoleMockRows{
		Role: sqlmock.NewRows([]string{
			"id", "name",
		}).AddRow(ids.RoleID, "role1"),
		Notfound: sqlmock.NewRows([]string{
			"id", "name",
		}),
	}

	domains := RoleMockDomains{
		Role: &domain.Role{
			ID:   ids.RoleID,
			Name: "role1",
		},
	}

	return mockDB, repo, ids, rows, domains

}

func TestRoleRepository_Create(t *testing.T) {
	mockDB, repo, ids, _, domains := RoleRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL := `INSERT INTO "roles" ("name","id") VALUES ($1,$2) RETURNING "id"`

	t.Run("should return no error when create successfully", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
			WithArgs(domains.Role.Name, ids.RoleID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))

		mockDB.Mock.ExpectCommit()

		err := repo.Create(context.TODO(), domains.Role)
		assert.Nil(t, err)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when failed create role", func(t *testing.T) {
		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
			WithArgs(domains.Role.Name, ids.RoleID).
			WillReturnError(errors.New("database error"))

		mockDB.Mock.ExpectRollback()

		err := repo.Create(context.TODO(), domains.Role)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}

func TestRoleRepository_FindAll(t *testing.T) {
	mockDB, repo, ids, rows, domains := RoleRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL := `SELECT * FROM "roles"`

	t.Run("should return roles commodities when find all successfully", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WillReturnRows(rows.Role)

		result, err := repo.FindAll(context.TODO())
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 1, len(result))
		assert.Equal(t, ids.RoleID, (result)[0].ID)
		assert.Equal(t, domains.Role.Name, (result)[0].Name)
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

func TestRoleRepository_FindByID(t *testing.T) {
	mockDB, repo, ids, rows, domains := RoleRepositorySetup(t)

	defer mockDB.SqlDB.Close()

	expectedSQL := `SELECT * FROM "roles" WHERE "roles"."id" = $1 ORDER BY "roles"."id" LIMIT $2`

	t.Run("should return role when find by id successfully", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.RoleID, 1).WillReturnRows(rows.Role)

		result, err := repo.FindByID(context.TODO(), ids.RoleID)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ids.RoleID, result.ID)
		assert.Equal(t, domains.Role.Name, result.Name)
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when find by id failed", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.RoleID, 1).WillReturnError(errors.New("database error"))

		result, err := repo.FindByID(context.TODO(), ids.RoleID)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "database error")
		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})

	t.Run("should return error when find by id not found", func(t *testing.T) {
		mockDB.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(ids.RoleID, 1).WillReturnRows(rows.Notfound)

		result, err := repo.FindByID(context.TODO(), ids.RoleID)
		assert.Nil(t, result)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))

		assert.Nil(t, mockDB.Mock.ExpectationsWereMet())
	})
}
