package usecase_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	"github.com/ryvasa/go-super-farmer/internal/repository/mock"
	usecase_implementation "github.com/ryvasa/go-super-farmer/internal/usecase/implementation"
	usecase_interface "github.com/ryvasa/go-super-farmer/internal/usecase/interface"
	"github.com/ryvasa/go-super-farmer/utils"
	"github.com/stretchr/testify/assert"
)

type ProvinceRepoMock struct {
	Province *mock.MockProvinceRepository
}

type ProvinceIDs struct {
	ProvinceID int64
}

type ProvinceMocks struct {
	Province        *domain.Province
	Provinces       []*domain.Province
	UpdatedProvince *domain.Province
}

type ProvinceDTOMock struct {
	Create *dto.ProvinceCreateDTO
	Update *dto.ProvinceUpdateDTO
}

func ProvinceUsecaseUtils(t *testing.T) (*ProvinceIDs, *ProvinceMocks, *ProvinceDTOMock, *ProvinceRepoMock, usecase_interface.ProvinceUsecase, context.Context) {
	provinceID := int64(1)

	ids := &ProvinceIDs{
		ProvinceID: provinceID,
	}

	mocks := &ProvinceMocks{
		Province: &domain.Province{
			ID:   provinceID,
			Name: "test province",
		},
		Provinces: []*domain.Province{
			{
				ID:   provinceID,
				Name: "test province",
			},
		},
		UpdatedProvince: &domain.Province{
			ID:   provinceID,
			Name: "updated province",
		},
	}

	dto := &ProvinceDTOMock{
		Create: &dto.ProvinceCreateDTO{
			Name: "test province",
		},
		Update: &dto.ProvinceUpdateDTO{
			Name: "updated province",
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	provinceRepo := mock.NewMockProvinceRepository(ctrl)
	uc := usecase_implementation.NewProvinceUsecase(provinceRepo)
	ctx := context.TODO()

	repo := &ProvinceRepoMock{Province: provinceRepo}

	return ids, mocks, dto, repo, uc, ctx
}

func TestProvinceUsecase_CreateProvince(t *testing.T) {
	ids, mocks, dtos, repo, uc, ctx := ProvinceUsecaseUtils(t)

	t.Run("should create province successfully", func(t *testing.T) {
		repo.Province.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, p *domain.Province) error {
			p.ID = 1
			return nil
		}).Times(1)

		repo.Province.EXPECT().FindByID(ctx, ids.ProvinceID).Return(mocks.Province, nil).Times(1)

		resp, err := uc.CreateProvince(ctx, dtos.Create)

		assert.NotNil(t, resp)
		assert.NoError(t, err)
		assert.Equal(t, dtos.Create.Name, resp.Name)
	})

	t.Run("should return error validation error", func(t *testing.T) {
		req := &dto.ProvinceCreateDTO{}
		resp, err := uc.CreateProvince(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("should return error when create province", func(t *testing.T) {
		repo.Province.EXPECT().Create(ctx, gomock.Any()).Return(utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.CreateProvince(ctx, dtos.Create)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "internal error")
	})

	t.Run("should return error when find created province by id", func(t *testing.T) {
		repo.Province.EXPECT().Create(ctx, gomock.Any()).Return(nil).Times(1)
		repo.Province.EXPECT().FindByID(ctx, gomock.Any()).Return(nil, utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.CreateProvince(ctx, dtos.Create)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "internal error")
	})
}

func TestProvinceUsecase_GetAllProvinces(t *testing.T) {
	_, mocks, _, repo, uc, ctx := ProvinceUsecaseUtils(t)

	t.Run("should return all provinces", func(t *testing.T) {
		repo.Province.EXPECT().FindAll(ctx).Return(mocks.Provinces, nil).Times(1)

		resp, err := uc.GetAllProvinces(ctx)

		assert.NoError(t, err)
		assert.Equal(t, 1, len(resp))
		assert.Equal(t, (mocks.Provinces)[0].Name, (resp)[0].Name)
	})

	t.Run("should return error internal error", func(t *testing.T) {
		repo.Province.EXPECT().FindAll(ctx).Return(nil, utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.GetAllProvinces(ctx)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "internal error")
	})
}

func TestProvinceUsecase_GetProvinceById(t *testing.T) {
	ids, mocks, _, repo, uc, ctx := ProvinceUsecaseUtils(t)

	t.Run("shpuld return province successfully", func(t *testing.T) {
		repo.Province.EXPECT().FindByID(ctx, ids.ProvinceID).Return(mocks.Province, nil).Times(1)

		resp, err := uc.GetProvinceByID(ctx, ids.ProvinceID)

		assert.NoError(t, err)
		assert.Equal(t, ids.ProvinceID, resp.ID)
		assert.Equal(t, mocks.Province.Name, resp.Name)
	})

	t.Run("should return error when get province by id", func(t *testing.T) {
		repo.Province.EXPECT().FindByID(ctx, gomock.Any()).Return(nil, utils.NewNotFoundError("province not found")).Times(1)

		resp, err := uc.GetProvinceByID(ctx, int64(2))

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "province not found")
	})
}

func TestProvinceUsecase_UpdateProvince(t *testing.T) {
	ids, mocks, dtos, repo, uc, ctx := ProvinceUsecaseUtils(t)
	t.Run("should update province successfully", func(t *testing.T) {
		repo.Province.EXPECT().FindByID(ctx, ids.ProvinceID).Return(mocks.Province, nil).Times(1)

		repo.Province.EXPECT().Update(ctx, ids.ProvinceID, gomock.Any()).DoAndReturn(func(ctx context.Context, id int64, p *domain.Province) error {
			p.ID = ids.ProvinceID
			return nil
		}).Times(1)

		repo.Province.EXPECT().FindByID(ctx, ids.ProvinceID).Return(mocks.UpdatedProvince, nil).Times(1)

		resp, err := uc.UpdateProvince(ctx, ids.ProvinceID, dtos.Update)

		assert.NoError(t, err)
		assert.Equal(t, mocks.UpdatedProvince.ID, resp.ID)
		assert.Equal(t, dtos.Update.Name, resp.Name)
	})

	t.Run("should return error validation error", func(t *testing.T) {
		req := &dto.ProvinceUpdateDTO{}
		resp, err := uc.UpdateProvince(ctx, ids.ProvinceID, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "Validation failed")
	})

	t.Run("should return error when province not found", func(t *testing.T) {
		repo.Province.EXPECT().FindByID(ctx, gomock.Any()).Return(nil, utils.NewNotFoundError("province not found")).Times(1)

		resp, err := uc.UpdateProvince(ctx, ids.ProvinceID, dtos.Update)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "province not found")
	})

	t.Run("should return error when update province", func(t *testing.T) {
		repo.Province.EXPECT().FindByID(ctx, ids.ProvinceID).Return(mocks.Province, nil).Times(1)

		repo.Province.EXPECT().Update(ctx, ids.ProvinceID, gomock.Any()).Return(utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.UpdateProvince(ctx, ids.ProvinceID, dtos.Update)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})

	t.Run("should return error when province not found after update", func(t *testing.T) {
		repo.Province.EXPECT().FindByID(ctx, ids.ProvinceID).Return(mocks.Province, nil).Times(1)

		repo.Province.EXPECT().Update(ctx, ids.ProvinceID, gomock.Any()).DoAndReturn(func(ctx context.Context, id int64, p *domain.Province) error {
			p.ID = ids.ProvinceID
			return nil
		}).Times(1)

		repo.Province.EXPECT().FindByID(ctx, ids.ProvinceID).Return(nil, utils.NewNotFoundError("province not found")).Times(1)

		req := &dto.ProvinceUpdateDTO{Name: "updated province"}
		resp, err := uc.UpdateProvince(ctx, ids.ProvinceID, req)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "province not found")
	})
}

func TestProvinceUsecase_DeleteProvince(t *testing.T) {
	ids, mocks, _, repo, uc, ctx := ProvinceUsecaseUtils(t)

	t.Run("should delete province successfully", func(t *testing.T) {
		repo.Province.EXPECT().FindByID(ctx, ids.ProvinceID).Return(mocks.Province, nil).Times(1)

		repo.Province.EXPECT().Delete(ctx, ids.ProvinceID).Return(nil).Times(1)

		err := uc.DeleteProvince(ctx, ids.ProvinceID)

		assert.NoError(t, err)
	})

	t.Run("should return error when province not found", func(t *testing.T) {
		repo.Province.EXPECT().FindByID(ctx, gomock.Any()).Return(nil, utils.NewNotFoundError("province not found")).Times(1)

		err := uc.DeleteProvince(ctx, int64(2))

		assert.Error(t, err)
		assert.EqualError(t, err, "province not found")
	})

	t.Run("should return error when delete province", func(t *testing.T) {
		repo.Province.EXPECT().FindByID(ctx, ids.ProvinceID).Return(mocks.Province, nil).Times(1)

		repo.Province.EXPECT().Delete(ctx, ids.ProvinceID).Return(utils.NewInternalError("internal error")).Times(1)

		err := uc.DeleteProvince(ctx, ids.ProvinceID)

		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}
