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

func TestCreateCity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cityRepo := mock.NewMockCityRepository(ctrl)
	uc := usecase.NewCityUsecase(cityRepo)
	ctx := context.TODO()

	cityID := int64(1)
	provinceID := int64(2)
	t.Run("Test CreateCity, successfully", func(t *testing.T) {
		mockCity := &domain.City{
			ID:         cityID,
			ProvinceID: provinceID,
			Name:       "test city",
		}
		cityRepo.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, p *domain.City) error {
			p.ID = 1
			return nil
		}).Times(1)

		cityRepo.EXPECT().FindByID(ctx, mockCity.ID).Return(mockCity, nil).Times(1)

		req := &dto.CityCreateDTO{Name: "test city", ProvinceID: provinceID}
		resp, err := uc.CreateCity(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, req.Name, resp.Name)
	})

	t.Run("Test CreateCity, validation error", func(t *testing.T) {
		req := &dto.CityCreateDTO{}
		resp, err := uc.CreateCity(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("Test CreateCity, error create city", func(t *testing.T) {
		cityRepo.EXPECT().Create(ctx, gomock.Any()).Return(utils.NewInternalError("internal error")).Times(1)

		req := &dto.CityCreateDTO{Name: "test city", ProvinceID: provinceID}
		resp, err := uc.CreateCity(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "internal error")
	})

	t.Run("Test CreateCity, error find city by id", func(t *testing.T) {
		cityRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil).Times(1)
		cityRepo.EXPECT().FindByID(ctx, gomock.Any()).Return(nil, utils.NewInternalError("internal error")).Times(1)

		req := &dto.CityCreateDTO{Name: "test city", ProvinceID: provinceID}
		resp, err := uc.CreateCity(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "internal error")
	})
}

func TestGetAllcities(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cityRepo := mock.NewMockCityRepository(ctrl)
	uc := usecase.NewCityUsecase(cityRepo)
	ctx := context.TODO()

	cityID := int64(1)
	provinceID := int64(2)
	t.Run("Test GetAllCities, successfully", func(t *testing.T) {
		mockCities := &[]domain.City{
			{ID: cityID, Name: "test city", ProvinceID: provinceID},
		}
		cityRepo.EXPECT().FindAll(ctx).Return(mockCities, nil).Times(1)

		resp, err := uc.GetAllCities(ctx)

		assert.NoError(t, err)
		assert.Equal(t, 1, len(*resp))
		assert.Equal(t, (*mockCities)[0].Name, (*resp)[0].Name)
	})

	t.Run("Test GetAllCities, error find all citys", func(t *testing.T) {
		cityRepo.EXPECT().FindAll(ctx).Return(nil, utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.GetAllCities(ctx)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "internal error")
	})
}

func TestGetCityById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cityRepo := mock.NewMockCityRepository(ctrl)
	uc := usecase.NewCityUsecase(cityRepo)
	ctx := context.TODO()

	cityID := int64(1)
	provinceID := int64(2)
	t.Run("Test GetCityById, successfully", func(t *testing.T) {
		mockCity := &domain.City{
			ID:         cityID,
			ProvinceID: provinceID,
			Name:       "test city",
		}
		cityRepo.EXPECT().FindByID(ctx, mockCity.ID).Return(mockCity, nil).Times(1)

		resp, err := uc.GetCityById(ctx, mockCity.ID)

		assert.NoError(t, err)
		assert.Equal(t, mockCity.Name, resp.Name)
	})

	t.Run("Test GetCityById, error find city by id", func(t *testing.T) {
		cityRepo.EXPECT().FindByID(ctx, gomock.Any()).Return(nil, utils.NewNotFoundError("city not found")).Times(1)

		resp, err := uc.GetCityById(ctx, int64(2))

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "city not found")
	})
}

func TestUpdateCity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cityRepo := mock.NewMockCityRepository(ctrl)
	uc := usecase.NewCityUsecase(cityRepo)
	ctx := context.TODO()

	cityID := int64(1)
	provinceID := int64(2)
	t.Run("Test UpdateCity, successfully", func(t *testing.T) {
		mockCity := &domain.City{
			ID:         cityID,
			Name:       "test city",
			ProvinceID: provinceID,
		}
		mockUpdate := &domain.City{
			ID:         cityID,
			Name:       "updated city",
			ProvinceID: provinceID,
		}
		cityRepo.EXPECT().FindByID(ctx, mockCity.ID).Return(mockCity, nil).Times(1)

		cityRepo.EXPECT().Update(ctx, cityID, mockUpdate).Return(nil).Times(1)

		cityRepo.EXPECT().FindByID(ctx, mockCity.ID).Return(mockUpdate, nil).Times(1)

		req := &dto.CityUpdateDTO{Name: "updated city", ProvinceID: provinceID}
		resp, err := uc.UpdateCity(ctx, mockCity.ID, req)

		assert.NoError(t, err)
		assert.Equal(t, req.Name, resp.Name)
	})

	t.Run("Test UpdateCity, validation error", func(t *testing.T) {
		req := &dto.CityUpdateDTO{}
		resp, err := uc.UpdateCity(ctx, provinceID, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("Test UpdateCity, error find city by id", func(t *testing.T) {
		cityRepo.EXPECT().FindByID(ctx, gomock.Any()).Return(nil, utils.NewNotFoundError("city not found")).Times(1)

		req := &dto.CityUpdateDTO{Name: "updated city", ProvinceID: provinceID}
		resp, err := uc.UpdateCity(ctx, int64(2), req)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "city not found")
	})

	t.Run("Test UpdateCity, error update city", func(t *testing.T) {
		mockCity := &domain.City{
			ID:         cityID,
			Name:       "test city",
			ProvinceID: provinceID,
		}
		cityRepo.EXPECT().FindByID(ctx, mockCity.ID).Return(mockCity, nil).Times(1)

		cityRepo.EXPECT().Update(ctx, mockCity.ID, gomock.Any()).Return(utils.NewInternalError("internal error")).Times(1)

		req := &dto.CityUpdateDTO{Name: "updated city", ProvinceID: provinceID}
		resp, err := uc.UpdateCity(ctx, mockCity.ID, req)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})

	t.Run("Test UpdateCity, error find city by id after update", func(t *testing.T) {
		mockCity := &domain.City{
			ID:         cityID,
			Name:       "test city",
			ProvinceID: provinceID,
		}
		mockUpdate := &domain.City{
			ID:         cityID,
			Name:       "updated city",
			ProvinceID: provinceID,
		}
		cityRepo.EXPECT().FindByID(ctx, mockCity.ID).Return(mockCity, nil).Times(1)

		cityRepo.EXPECT().Update(ctx, mockCity.ID, mockUpdate).Return(nil).Times(1)

		cityRepo.EXPECT().FindByID(ctx, mockCity.ID).Return(nil, utils.NewNotFoundError("city not found")).Times(1)

		req := &dto.CityUpdateDTO{Name: "updated city", ProvinceID: provinceID}
		resp, err := uc.UpdateCity(ctx, mockCity.ID, req)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "city not found")
	})
}

func TestDeleteCity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cityRepo := mock.NewMockCityRepository(ctrl)
	uc := usecase.NewCityUsecase(cityRepo)
	ctx := context.TODO()

	cityID := int64(1)

	t.Run("Test DeleteCity, successfully", func(t *testing.T) {
		mockCity := &domain.City{
			ID:   cityID,
			Name: "test city",
		}
		cityRepo.EXPECT().FindByID(ctx, mockCity.ID).Return(mockCity, nil).Times(1)

		cityRepo.EXPECT().Delete(ctx, mockCity.ID).Return(nil).Times(1)

		err := uc.DeleteCity(ctx, mockCity.ID)

		assert.NoError(t, err)
	})

	t.Run("Test DeleteCity, error find city by id", func(t *testing.T) {
		cityRepo.EXPECT().FindByID(ctx, gomock.Any()).Return(nil, utils.NewNotFoundError("city not found")).Times(1)

		err := uc.DeleteCity(ctx, int64(2))

		assert.Error(t, err)
		assert.EqualError(t, err, "city not found")
	})

	t.Run("Test DeleteCity, error delete city", func(t *testing.T) {
		mockCity := &domain.City{
			ID:   cityID,
			Name: "test city",
		}
		cityRepo.EXPECT().FindByID(ctx, mockCity.ID).Return(mockCity, nil).Times(1)

		cityRepo.EXPECT().Delete(ctx, mockCity.ID).Return(utils.NewInternalError("internal error")).Times(1)

		err := uc.DeleteCity(ctx, mockCity.ID)

		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}
