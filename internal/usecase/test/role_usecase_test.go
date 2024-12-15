package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	"github.com/ryvasa/go-super-farmer/internal/repository/mock"
	usecase_implementation "github.com/ryvasa/go-super-farmer/internal/usecase/implementation"
	usecase_interface "github.com/ryvasa/go-super-farmer/internal/usecase/interface"
	"github.com/stretchr/testify/assert"
)

type RoleRepoMock struct {
	Role *mock.MockRoleRepository
}

type RoleIDs struct {
	RoleID int64
}

type RoleMocks struct {
	Role  *domain.Role
	Roles *[]domain.Role
}

type RoleDTOMock struct {
	Create *dto.RoleCreateDTO
}

func RoleUsecaseUtils(t *testing.T) (*RoleIDs, *RoleMocks, *RoleDTOMock, *RoleRepoMock, usecase_interface.RoleUsecase, context.Context) {
	roleID := int64(1)

	ids := &RoleIDs{
		RoleID: roleID,
	}

	mocks := &RoleMocks{
		Role: &domain.Role{
			ID:   roleID,
			Name: "admin",
		},
		Roles: &[]domain.Role{
			{
				ID:   roleID,
				Name: "admin",
			},
			{
				ID:   roleID + 1,
				Name: "farmer",
			},
		},
	}

	dto := &RoleDTOMock{
		Create: &dto.RoleCreateDTO{
			Name: "admin",
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	roleRepo := mock.NewMockRoleRepository(ctrl)
	uc := usecase_implementation.NewRoleUsecase(roleRepo)
	ctx := context.Background()

	repo := &RoleRepoMock{Role: roleRepo}

	return ids, mocks, dto, repo, uc, ctx
}

func TestCreateRole(t *testing.T) {
	ids, mocks, dtos, repo, uc, ctx := RoleUsecaseUtils(t)
	t.Run("should create role successfully", func(t *testing.T) {
		repo.Role.EXPECT().FindAll(ctx).Return(&[]domain.Role{}, nil).Times(1)
		repo.Role.EXPECT().Create(ctx, gomock.Any()).Return(nil).Times(1)
		repo.Role.EXPECT().FindByID(ctx, ids.RoleID).Return(mocks.Role, nil).Times(1)

		resp, err := uc.CreateRole(ctx, dtos.Create)

		assert.NotNil(t, resp)
		assert.NoError(t, err)
		assert.Equal(t, dtos.Create.Name, resp.Name)
	})

	t.Run("Test CreateRole, validation error", func(t *testing.T) {
		resp, err := uc.CreateRole(ctx, &dto.RoleCreateDTO{Name: ""})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "Validation failed")
	})

	t.Run("Test CreateRole, error get created role", func(t *testing.T) {
		repo.Role.EXPECT().FindAll(ctx).Return(&[]domain.Role{}, nil).Times(1)
		repo.Role.EXPECT().Create(ctx, gomock.Any()).Return(nil).Times(1)
		repo.Role.EXPECT().FindByID(ctx, ids.RoleID).Return(nil, errors.New("internal error")).Times(1)

		resp, err := uc.CreateRole(ctx, dtos.Create)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})

	t.Run("Test CreateRole, error get all roles", func(t *testing.T) {
		repo.Role.EXPECT().FindAll(ctx).Return(&[]domain.Role{}, errors.New("internal error")).Times(1)

		resp, err := uc.CreateRole(ctx, dtos.Create)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("Test CreateRole, error create role", func(t *testing.T) {
		repo.Role.EXPECT().FindAll(ctx).Return(&[]domain.Role{}, nil).Times(1)
		repo.Role.EXPECT().Create(ctx, gomock.Any()).Return(errors.New("internal error")).Times(1)

		resp, err := uc.CreateRole(ctx, dtos.Create)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestGetAllRoles(t *testing.T) {

	_, mocks, _, repo, uc, ctx := RoleUsecaseUtils(t)

	t.Run("Test GetAllRoles, successfully", func(t *testing.T) {

		repo.Role.EXPECT().FindAll(ctx).Return(mocks.Roles, nil).Times(1)

		resp, err := uc.GetAllRoles(ctx)

		assert.NoError(t, err)
		assert.Len(t, *resp, len(*mocks.Roles))
	})

	t.Run("Test GetAllRoles, database error", func(t *testing.T) {
		repo.Role.EXPECT().FindAll(ctx).Return(nil, errors.New("internal error")).Times(1)

		resp, err := uc.GetAllRoles(ctx)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}
