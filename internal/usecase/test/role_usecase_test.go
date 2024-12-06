package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	"github.com/ryvasa/go-super-farmer/internal/repository/mock"
	"github.com/ryvasa/go-super-farmer/internal/usecase"
	"github.com/stretchr/testify/assert"
)

func TestCreateRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock.NewMockRoleRepository(ctrl)
	uc := usecase.NewRoleUsecase(repo)
	ctx := context.Background()
	t.Run("Test CreateRole, successfully", func(t *testing.T) {
		mockRole := &domain.Role{ID: 1, Name: "admin"}

		repo.EXPECT().FindAll(ctx).Return(&[]domain.Role{}, nil).Times(1)
		repo.EXPECT().Create(ctx, gomock.Any()).Return(nil).Times(1)
		repo.EXPECT().FindByID(ctx, uint64(1)).Return(mockRole, nil).Times(1)

		req := &dto.RoleCreateDTO{Name: "admin"}
		resp, err := uc.CreateRole(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, req.Name, resp.Name)
	})

	t.Run("Test CreateRole, validation error", func(t *testing.T) {
		req := &dto.RoleCreateDTO{Name: ""}
		resp, err := uc.CreateRole(ctx, req)

		assert.Error(t, err)
		// jika return nya domain role gunakan empty, jika nil gunakan nil
		assert.Nil(t, resp)
	})

	t.Run("Test CreateRole, error get created role", func(t *testing.T) {
		repo.EXPECT().FindAll(ctx).Return(&[]domain.Role{}, nil).Times(1)
		repo.EXPECT().Create(ctx, gomock.Any()).Return(nil).Times(1)
		repo.EXPECT().FindByID(ctx, uint64(1)).Return(nil, errors.New("internal error")).Times(1)

		req := &dto.RoleCreateDTO{Name: "admin"}
		resp, err := uc.CreateRole(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("Test CreateRole, error get all roles", func(t *testing.T) {
		repo.EXPECT().FindAll(ctx).Return(&[]domain.Role{}, errors.New("internal error")).Times(1)

		req := &dto.RoleCreateDTO{Name: "admin"}
		resp, err := uc.CreateRole(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("Test CreateRole, error create role", func(t *testing.T) {
		repo.EXPECT().FindAll(ctx).Return(&[]domain.Role{}, nil).Times(1)
		repo.EXPECT().Create(ctx, gomock.Any()).Return(errors.New("internal error")).Times(1)

		req := &dto.RoleCreateDTO{Name: "admin"}
		resp, err := uc.CreateRole(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestGetAllRoles(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock.NewMockRoleRepository(ctrl)
	uc := usecase.NewRoleUsecase(repo)
	ctx := context.Background()
	t.Run("Test GetAllRoles, successfully", func(t *testing.T) {
		mockRoles := &[]domain.Role{
			{ID: 1, Name: "admin"},
			{ID: 2, Name: "farmer"},
		}

		repo.EXPECT().FindAll(ctx).Return(mockRoles, nil).Times(1)

		resp, err := uc.GetAllRoles(ctx)

		assert.NoError(t, err)
		assert.Len(t, *resp, len(*mockRoles))
	})

	t.Run("Test GetAllRoles, database error", func(t *testing.T) {
		repo.EXPECT().FindAll(ctx).Return(nil, errors.New("internal error")).Times(1)

		resp, err := uc.GetAllRoles(ctx)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}
