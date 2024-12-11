package usecase_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	"github.com/ryvasa/go-super-farmer/internal/repository/mock"
	"github.com/ryvasa/go-super-farmer/internal/usecase"
	"github.com/ryvasa/go-super-farmer/utils"
	"github.com/stretchr/testify/assert"
)

func TestCreateRegion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	regionRepo := mock.NewMockRegionRepository(ctrl)
	cityRepo := mock.NewMockCityRepository(ctrl)
	provinceRepo := mock.NewMockProvinceRepository(ctrl)
	uc := usecase.NewRegionUsecase(regionRepo, cityRepo, provinceRepo)
	ctx := context.Background()

	regionID := uuid.New()
	cityID := int64(1)
	provinceID := int64(1)

	mockRegion := &domain.Region{
		ID:         regionID,
		CityID:     cityID,
		ProvinceID: provinceID,
	}

	mockProvince := &domain.Province{
		ID:   provinceID,
		Name: "province",
	}

	mockCity := &domain.City{
		ID:         cityID,
		Name:       "city",
		ProvinceID: provinceID,
	}
	t.Run("should create region successfully", func(t *testing.T) {

		cityRepo.EXPECT().FindByID(ctx, cityID).Return(mockCity, nil).Times(1)

		provinceRepo.EXPECT().FindByID(ctx, provinceID).Return(mockProvince, nil).Times(1)

		regionRepo.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, r *domain.Region) error {
			r.ID = regionID
			return nil
		}).Times(1)

		regionRepo.EXPECT().FindByID(ctx, regionID).Return(mockRegion, nil).Times(1)

		req := &dto.RegionCreateDto{
			CityID:     cityID,
			ProvinceID: provinceID,
		}
		resp, err := uc.CreateRegion(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, req.CityID, resp.CityID)
		assert.Equal(t, req.ProvinceID, resp.ProvinceID)
		assert.Equal(t, regionID, resp.ID)
	})

	t.Run("should return error validation", func(t *testing.T) {
		req := &dto.RegionCreateDto{
			CityID:     cityID,
			ProvinceID: 0,
		}
		_, err := uc.CreateRegion(ctx, req)
		assert.Error(t, err)
		assert.EqualError(t, err, "Validation failed")
	})

	t.Run("should return error when city not found", func(t *testing.T) {
		cityRepo.EXPECT().FindByID(ctx, cityID).Return(nil, utils.NewNotFoundError("city not found")).Times(1)

		req := &dto.RegionCreateDto{
			CityID:     cityID,
			ProvinceID: provinceID,
		}
		_, err := uc.CreateRegion(ctx, req)
		assert.Error(t, err)
		assert.EqualError(t, err, "city not found")
	})

	t.Run("should return error when province not found", func(t *testing.T) {
		cityRepo.EXPECT().FindByID(ctx, cityID).Return(mockCity, nil).Times(1)

		provinceRepo.EXPECT().FindByID(ctx, provinceID).Return(nil, utils.NewNotFoundError("province not found")).Times(1)

		req := &dto.RegionCreateDto{
			CityID:     cityID,
			ProvinceID: provinceID,
		}
		resp, err := uc.CreateRegion(ctx, req)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "province not found")
	})

	t.Run("should return error when create region", func(t *testing.T) {
		cityRepo.EXPECT().FindByID(ctx, cityID).Return(mockCity, nil).Times(1)

		provinceRepo.EXPECT().FindByID(ctx, provinceID).Return(mockProvince, nil).Times(1)

		regionRepo.EXPECT().Create(ctx, gomock.Any()).Return(utils.NewInternalError("internal error")).Times(1)

		req := &dto.RegionCreateDto{
			CityID:     cityID,
			ProvinceID: provinceID,
		}
		_, err := uc.CreateRegion(ctx, req)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})

	t.Run("should return error when get created region", func(t *testing.T) {
		cityRepo.EXPECT().FindByID(ctx, cityID).Return(mockCity, nil).Times(1)

		provinceRepo.EXPECT().FindByID(ctx, provinceID).Return(mockProvince, nil).Times(1)

		regionRepo.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, r *domain.Region) error {
			r.ID = regionID
			return nil
		}).Times(1)

		regionRepo.EXPECT().FindByID(ctx, regionID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		req := &dto.RegionCreateDto{
			CityID:     cityID,
			ProvinceID: provinceID,
		}
		_, err := uc.CreateRegion(ctx, req)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

func TestGetAllRegions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	regionRepo := mock.NewMockRegionRepository(ctrl)
	cityRepo := mock.NewMockCityRepository(ctrl)
	provinceRepo := mock.NewMockProvinceRepository(ctrl)
	uc := usecase.NewRegionUsecase(regionRepo, cityRepo, provinceRepo)
	ctx := context.Background()

	regionID := uuid.New()
	cityID := int64(1)
	provinceID := int64(1)

	mockRegions := &[]domain.Region{
		{
			ID:         regionID,
			CityID:     cityID,
			ProvinceID: provinceID,
		},
	}

	t.Run("should return all regions", func(t *testing.T) {
		regionRepo.EXPECT().FindAll(ctx).Return(mockRegions, nil).Times(1)

		res, err := uc.GetAllRegions(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(*res))
		assert.Equal(t, (*mockRegions)[0].ID, (*res)[0].ID)
		assert.Equal(t, (*mockRegions)[0].CityID, (*res)[0].CityID)
		assert.Equal(t, (*mockRegions)[0].ProvinceID, (*res)[0].ProvinceID)
	})

	t.Run("should return error when get all regions", func(t *testing.T) {
		regionRepo.EXPECT().FindAll(ctx).Return(nil, utils.NewInternalError("internal error")).Times(1)

		_, err := uc.GetAllRegions(ctx)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

func TestGetRegionByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	regionRepo := mock.NewMockRegionRepository(ctrl)
	cityRepo := mock.NewMockCityRepository(ctrl)
	provinceRepo := mock.NewMockProvinceRepository(ctrl)
	uc := usecase.NewRegionUsecase(regionRepo, cityRepo, provinceRepo)
	ctx := context.Background()

	regionID := uuid.New()
	cityID := int64(1)
	provinceID := int64(1)

	mockRegion := &domain.Region{
		ID:         regionID,
		CityID:     cityID,
		ProvinceID: provinceID,
	}

	t.Run("should return region by id success", func(t *testing.T) {
		regionRepo.EXPECT().FindByID(ctx, regionID).Return(mockRegion, nil).Times(1)

		res, err := uc.GetRegionByID(ctx, regionID)
		assert.NoError(t, err)
		assert.Equal(t, regionID, res.ID)
		assert.Equal(t, cityID, res.CityID)
		assert.Equal(t, provinceID, res.ProvinceID)
	})

	t.Run("should return error when get region by id", func(t *testing.T) {
		regionRepo.EXPECT().FindByID(ctx, regionID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		_, err := uc.GetRegionByID(ctx, regionID)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}
