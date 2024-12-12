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

type CityRepoMock struct {
	City *mock.MockCityRepository
}

type CityIDs struct {
	CityID     int64
	ProvinceID int64
}

type CityMocks struct {
	City        *domain.City
	Cities      *[]domain.City
	UpdatedCity *domain.City
}

type CityDTOMock struct {
	Create *dto.CityCreateDTO
	Update *dto.CityUpdateDTO
}

func CityUsecaseUtils(t *testing.T) (*CityIDs, *CityMocks, *CityDTOMock, *CityRepoMock, usecase.CityUsecase, context.Context) {

	cityID := int64(1)
	provinceID := int64(2)

	ids := &CityIDs{
		CityID:     cityID,
		ProvinceID: provinceID,
	}

	mocks := &CityMocks{
		City: &domain.City{
			ID:         cityID,
			Name:       "test",
			ProvinceID: provinceID,
		},
		Cities: &[]domain.City{
			{
				ID:         cityID,
				Name:       "test",
				ProvinceID: provinceID,
			},
		},
		UpdatedCity: &domain.City{
			ID:         cityID,
			Name:       "updated",
			ProvinceID: provinceID,
		},
	}

	dto := &CityDTOMock{
		Create: &dto.CityCreateDTO{
			Name:       "test",
			ProvinceID: provinceID,
		},
		Update: &dto.CityUpdateDTO{
			Name: "updated",
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cityRepo := mock.NewMockCityRepository(ctrl)
	uc := usecase.NewCityUsecase(cityRepo)
	ctx := context.TODO()

	repo := &CityRepoMock{City: cityRepo}

	return ids, mocks, dto, repo, uc, ctx
}

func TestCreateCity(t *testing.T) {
	ids, mocks, dtos, repo, uc, ctx := CityUsecaseUtils(t)
	t.Run("should create city successfully", func(t *testing.T) {
		repo.City.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, p *domain.City) error {
			p.ID = ids.CityID
			return nil
		}).Times(1)

		repo.City.EXPECT().FindByID(ctx, ids.CityID).Return(mocks.City, nil).Times(1)

		resp, err := uc.CreateCity(ctx, dtos.Create)

		assert.NotNil(t, resp)
		assert.NoError(t, err)
		assert.Equal(t, dtos.Create.Name, resp.Name)
	})

	t.Run("should return error validation error", func(t *testing.T) {
		req := &dto.CityCreateDTO{}
		resp, err := uc.CreateCity(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "Validation failed")
	})

	t.Run("should return error when create city", func(t *testing.T) {
		repo.City.EXPECT().Create(ctx, gomock.Any()).Return(utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.CreateCity(ctx, dtos.Create)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "internal error")
	})

	t.Run("should return error when find created city by id", func(t *testing.T) {
		repo.City.EXPECT().Create(ctx, gomock.Any()).Return(nil).Times(1)
		repo.City.EXPECT().FindByID(ctx, gomock.Any()).Return(nil, utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.CreateCity(ctx, dtos.Create)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "internal error")
	})
}

func TestGetAllcities(t *testing.T) {
	_, mocks, _, repo, uc, ctx := CityUsecaseUtils(t)

	t.Run("should return all cities", func(t *testing.T) {
		repo.City.EXPECT().FindAll(ctx).Return(mocks.Cities, nil).Times(1)

		resp, err := uc.GetAllCities(ctx)

		assert.NoError(t, err)
		assert.Equal(t, 1, len(*resp))
		assert.Equal(t, (*mocks.Cities)[0].Name, (*resp)[0].Name)
	})

	t.Run("should return error internal error", func(t *testing.T) {
		repo.City.EXPECT().FindAll(ctx).Return(nil, utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.GetAllCities(ctx)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "internal error")
	})
}

func TestGetCityById(t *testing.T) {
	ids, mocks, _, repo, uc, ctx := CityUsecaseUtils(t)

	t.Run("shpuld return city successfully", func(t *testing.T) {
		repo.City.EXPECT().FindByID(ctx, ids.CityID).Return(mocks.City, nil).Times(1)

		resp, err := uc.GetCityById(ctx, ids.CityID)

		assert.NotNil(t, resp)
		assert.NoError(t, err)
		assert.Equal(t, mocks.City.Name, resp.Name)
	})

	t.Run("should return error when find city by id", func(t *testing.T) {
		repo.City.EXPECT().FindByID(ctx, gomock.Any()).Return(nil, utils.NewNotFoundError("city not found")).Times(1)

		resp, err := uc.GetCityById(ctx, int64(2))

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "city not found")
	})
}

func TestUpdateCity(t *testing.T) {
	ids, mocks, dtos, repo, uc, ctx := CityUsecaseUtils(t)

	t.Run("should update city successfully", func(t *testing.T) {
		repo.City.EXPECT().FindByID(ctx, ids.CityID).Return(mocks.City, nil).Times(1)

		repo.City.EXPECT().Update(ctx, ids.CityID, gomock.Any()).DoAndReturn(func(ctx context.Context, id int64, p *domain.City) error {
			p.ID = ids.CityID
			return nil
		}).Times(1)

		repo.City.EXPECT().FindByID(ctx, ids.CityID).Return(mocks.UpdatedCity, nil).Times(1)

		resp, err := uc.UpdateCity(ctx, ids.CityID, dtos.Update)

		assert.NoError(t, err)
		assert.Equal(t, dtos.Update.Name, resp.Name)
	})

	t.Run("should return error validation error", func(t *testing.T) {
		req := &dto.CityUpdateDTO{}
		resp, err := uc.UpdateCity(ctx, ids.CityID, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("should return error when find city by id", func(t *testing.T) {
		repo.City.EXPECT().FindByID(ctx, gomock.Any()).Return(nil, utils.NewNotFoundError("city not found")).Times(1)

		resp, err := uc.UpdateCity(ctx, ids.CityID, dtos.Update)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "city not found")
	})

	t.Run("should return error when update city", func(t *testing.T) {
		repo.City.EXPECT().FindByID(ctx, ids.CityID).Return(mocks.City, nil).Times(1)

		repo.City.EXPECT().Update(ctx, ids.CityID, gomock.Any()).Return(utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.UpdateCity(ctx, ids.CityID, dtos.Update)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})

	t.Run("should return error when find city by id after update", func(t *testing.T) {
		repo.City.EXPECT().FindByID(ctx, ids.CityID).Return(mocks.City, nil).Times(1)

		repo.City.EXPECT().Update(ctx, ids.CityID, gomock.Any()).DoAndReturn(func(ctx context.Context, id int64, p *domain.City) error {
			p.ID = ids.CityID
			return nil
		}).Times(1)

		repo.City.EXPECT().FindByID(ctx, ids.CityID).Return(nil, utils.NewNotFoundError("city not found")).Times(1)

		resp, err := uc.UpdateCity(ctx, ids.CityID, dtos.Update)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "city not found")
	})
}

func TestDeleteCity(t *testing.T) {
	ids, mocks, _, repo, uc, ctx := CityUsecaseUtils(t)

	t.Run("should delete city successfully", func(t *testing.T) {
		repo.City.EXPECT().FindByID(ctx, ids.CityID).Return(mocks.City, nil).Times(1)

		repo.City.EXPECT().Delete(ctx, ids.CityID).Return(nil).Times(1)

		err := uc.DeleteCity(ctx, ids.CityID)

		assert.NoError(t, err)
	})

	t.Run("should return error when find city by id", func(t *testing.T) {
		repo.City.EXPECT().FindByID(ctx, gomock.Any()).Return(nil, utils.NewNotFoundError("city not found")).Times(1)

		err := uc.DeleteCity(ctx, int64(2))

		assert.Error(t, err)
		assert.EqualError(t, err, "city not found")
	})

	t.Run("should return error when delete city", func(t *testing.T) {
		repo.City.EXPECT().FindByID(ctx, ids.CityID).Return(mocks.City, nil).Times(1)

		repo.City.EXPECT().Delete(ctx, ids.CityID).Return(utils.NewInternalError("internal error")).Times(1)

		err := uc.DeleteCity(ctx, ids.CityID)

		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}
