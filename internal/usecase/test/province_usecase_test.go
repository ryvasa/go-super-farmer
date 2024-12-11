package usecase_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	"github.com/ryvasa/go-super-farmer/internal/repository/mock"
	"github.com/ryvasa/go-super-farmer/internal/usecase"
	"github.com/ryvasa/go-super-farmer/utils"
	"github.com/stretchr/testify/assert"
)

func TestCreateProvince(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	provinceRepo := mock.NewMockProvinceRepository(ctrl)
	uc := usecase.NewProvinceUsecase(provinceRepo)
	ctx := context.TODO()

	t.Run("Test CreateProvince, successfully", func(t *testing.T) {
		mockProvince := &domain.Province{
			ID:   int64(1),
			Name: "test province",
		}
		provinceRepo.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, p *domain.Province) error {
			p.ID = 1
			return nil
		}).Times(1)

		provinceRepo.EXPECT().FindByID(ctx, mockProvince.ID).Return(mockProvince, nil).Times(1)

		req := &dto.ProvinceCreateDTO{Name: "test province"}
		resp, err := uc.CreateProvince(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, req.Name, resp.Name)
	})

	t.Run("Test CreateProvince, validation error", func(t *testing.T) {
		req := &dto.ProvinceCreateDTO{}
		resp, err := uc.CreateProvince(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("Test CreateProvince, error create province", func(t *testing.T) {
		provinceRepo.EXPECT().Create(ctx, gomock.Any()).Return(utils.NewInternalError("internal error")).Times(1)

		req := &dto.ProvinceCreateDTO{Name: "test province"}
		resp, err := uc.CreateProvince(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "internal error")
	})

	t.Run("Test CreateProvince, error find province by id", func(t *testing.T) {
		provinceRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil).Times(1)
		provinceRepo.EXPECT().FindByID(ctx, gomock.Any()).Return(nil, utils.NewInternalError("internal error")).Times(1)

		req := &dto.ProvinceCreateDTO{Name: "test province"}
		resp, err := uc.CreateProvince(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "internal error")
	})
}

func TestGetAllProvinces(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	provinceRepo := mock.NewMockProvinceRepository(ctrl)
	uc := usecase.NewProvinceUsecase(provinceRepo)
	ctx := context.TODO()
	t.Run("Test GetAllProvinces, successfully", func(t *testing.T) {
		mockProvinces := &[]domain.Province{
			{ID: int64(1), Name: "test province"},
		}
		provinceRepo.EXPECT().FindAll(ctx).Return(mockProvinces, nil).Times(1)

		resp, err := uc.GetAllProvinces(ctx)

		assert.NoError(t, err)
		assert.Equal(t, 1, len(*resp))
		assert.Equal(t, (*mockProvinces)[0].Name, (*resp)[0].Name)
	})

	t.Run("Test GetAllProvinces, error find all provinces", func(t *testing.T) {
		provinceRepo.EXPECT().FindAll(ctx).Return(nil, utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.GetAllProvinces(ctx)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "internal error")
	})
}

func TestGetProvinceById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	provinceRepo := mock.NewMockProvinceRepository(ctrl)
	uc := usecase.NewProvinceUsecase(provinceRepo)
	ctx := context.TODO()

	t.Run("Test GetProvinceById, successfully", func(t *testing.T) {
		mockProvince := &domain.Province{
			ID:   int64(1),
			Name: "test province",
		}
		provinceRepo.EXPECT().FindByID(ctx, mockProvince.ID).Return(mockProvince, nil).Times(1)

		resp, err := uc.GetProvinceById(ctx, mockProvince.ID)

		assert.NoError(t, err)
		assert.Equal(t, mockProvince.Name, resp.Name)
	})

	t.Run("Test GetProvinceById, error find province by id", func(t *testing.T) {
		provinceRepo.EXPECT().FindByID(ctx, gomock.Any()).Return(nil, utils.NewNotFoundError("province not found")).Times(1)

		resp, err := uc.GetProvinceById(ctx, int64(2))

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "province not found")
	})
}

func TestUpdateProvince(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	provinceRepo := mock.NewMockProvinceRepository(ctrl)
	uc := usecase.NewProvinceUsecase(provinceRepo)
	ctx := context.TODO()

	provinceID := int64(1)
	t.Run("Test UpdateProvince, successfully", func(t *testing.T) {
		mockProvince := &domain.Province{
			ID:   provinceID,
			Name: "test province",
		}
		mockUpdate := &domain.Province{
			ID:   provinceID,
			Name: "updated province",
		}
		provinceRepo.EXPECT().FindByID(ctx, mockProvince.ID).Return(mockProvince, nil).Times(1)

		provinceRepo.EXPECT().Update(ctx, provinceID, mockUpdate).Return(nil).Times(1)

		provinceRepo.EXPECT().FindByID(ctx, mockProvince.ID).Return(mockUpdate, nil).Times(1)

		req := &dto.ProvinceUpdateDTO{Name: "updated province"}
		resp, err := uc.UpdateProvince(ctx, mockProvince.ID, req)

		assert.NoError(t, err)
		assert.Equal(t, req.Name, resp.Name)
	})

	t.Run("Test UpdateProvince, validation error", func(t *testing.T) {
		req := &dto.ProvinceUpdateDTO{}
		resp, err := uc.UpdateProvince(ctx, int64(2), req)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("Test UpdateProvince, error find province by id", func(t *testing.T) {
		provinceRepo.EXPECT().FindByID(ctx, gomock.Any()).Return(nil, utils.NewNotFoundError("province not found")).Times(1)

		req := &dto.ProvinceUpdateDTO{Name: "updated province"}
		resp, err := uc.UpdateProvince(ctx, int64(2), req)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "province not found")
	})

	t.Run("Test UpdateProvince, error update province", func(t *testing.T) {
		mockProvince := &domain.Province{
			ID:   provinceID,
			Name: "test province",
		}
		provinceRepo.EXPECT().FindByID(ctx, mockProvince.ID).Return(mockProvince, nil).Times(1)

		provinceRepo.EXPECT().Update(ctx, mockProvince.ID, gomock.Any()).Return(utils.NewInternalError("internal error")).Times(1)

		req := &dto.ProvinceUpdateDTO{Name: "updated province"}
		resp, err := uc.UpdateProvince(ctx, mockProvince.ID, req)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})

	t.Run("Test UpdateProvince, error find province by id after update", func(t *testing.T) {
		mockProvince := &domain.Province{
			ID:   provinceID,
			Name: "test province",
		}
		mockUpdate := &domain.Province{
			ID:   provinceID,
			Name: "updated province",
		}
		provinceRepo.EXPECT().FindByID(ctx, mockProvince.ID).Return(mockProvince, nil).Times(1)

		provinceRepo.EXPECT().Update(ctx, mockProvince.ID, mockUpdate).Return(nil).Times(1)

		provinceRepo.EXPECT().FindByID(ctx, mockProvince.ID).Return(nil, utils.NewNotFoundError("province not found")).Times(1)

		req := &dto.ProvinceUpdateDTO{Name: "updated province"}
		resp, err := uc.UpdateProvince(ctx, mockProvince.ID, req)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "province not found")
	})
}

func TestDeleteProvince(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	provinceRepo := mock.NewMockProvinceRepository(ctrl)
	uc := usecase.NewProvinceUsecase(provinceRepo)
	ctx := context.TODO()

	provinceID := int64(1)

	t.Run("Test DeleteProvince, successfully", func(t *testing.T) {
		mockProvince := &domain.Province{
			ID:   provinceID,
			Name: "test province",
		}
		provinceRepo.EXPECT().FindByID(ctx, mockProvince.ID).Return(mockProvince, nil).Times(1)

		provinceRepo.EXPECT().Delete(ctx, mockProvince.ID).Return(nil).Times(1)

		err := uc.DeleteProvince(ctx, mockProvince.ID)

		assert.NoError(t, err)
	})

	t.Run("Test DeleteProvince, error find province by id", func(t *testing.T) {
		provinceRepo.EXPECT().FindByID(ctx, gomock.Any()).Return(nil, utils.NewNotFoundError("province not found")).Times(1)

		err := uc.DeleteProvince(ctx, int64(2))

		assert.Error(t, err)
		assert.EqualError(t, err, "province not found")
	})

	t.Run("Test DeleteProvince, error delete province", func(t *testing.T) {
		mockProvince := &domain.Province{
			ID:   provinceID,
			Name: "test province",
		}
		provinceRepo.EXPECT().FindByID(ctx, mockProvince.ID).Return(mockProvince, nil).Times(1)

		provinceRepo.EXPECT().Delete(ctx, mockProvince.ID).Return(utils.NewInternalError("internal error")).Times(1)

		err := uc.DeleteProvince(ctx, mockProvince.ID)

		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}
