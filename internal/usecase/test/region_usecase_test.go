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

type RegionRepoMock struct {
	Region   *mock.MockRegionRepository
	City     *mock.MockCityRepository
	Province *mock.MockProvinceRepository
}

type RegionMocks struct {
	Region        *domain.Region
	Regions       *[]domain.Region
	UpdatedRegion *domain.Region
	City          *domain.City
	Province      *domain.Province
}

type RegionIDs struct {
	RegionID   uuid.UUID
	CityID     int64
	ProvinceID int64
}

type RegionDTOMocks struct {
	Create *dto.RegionCreateDto
}

func RegionUsecaseUtils(t *testing.T) (*RegionIDs, *RegionMocks, *RegionDTOMocks, *RegionRepoMock, usecase.RegionUseCase, context.Context) {
	regionID := uuid.New()
	cityID := int64(1)
	provinceID := int64(1)

	ids := &RegionIDs{
		RegionID:   regionID,
		CityID:     cityID,
		ProvinceID: provinceID,
	}

	mocks := &RegionMocks{
		Region: &domain.Region{
			ID:         regionID,
			CityID:     cityID,
			ProvinceID: provinceID,
		},
		Regions: &[]domain.Region{
			{
				ID:         regionID,
				CityID:     cityID,
				ProvinceID: provinceID,
			},
		},
		City: &domain.City{
			ID:         cityID,
			Name:       "city",
			ProvinceID: provinceID,
		},
		Province: &domain.Province{
			ID:   provinceID,
			Name: "province",
		},
	}

	dtoMocks := &RegionDTOMocks{
		Create: &dto.RegionCreateDto{
			CityID:     cityID,
			ProvinceID: provinceID,
		},
	}

	regionRepoMock := &RegionRepoMock{
		Region:   mock.NewMockRegionRepository(gomock.NewController(t)),
		City:     mock.NewMockCityRepository(gomock.NewController(t)),
		Province: mock.NewMockProvinceRepository(gomock.NewController(t)),
	}

	uc := usecase.NewRegionUsecase(regionRepoMock.Region, regionRepoMock.City, regionRepoMock.Province)
	ctx := context.Background()

	return ids, mocks, dtoMocks, regionRepoMock, uc, ctx
}

func TestCreateRegion(t *testing.T) {
	ids, mocks, dtos, repo, uc, ctx := RegionUsecaseUtils(t)
	t.Run("should create region successfully", func(t *testing.T) {

		repo.City.EXPECT().FindByID(ctx, ids.CityID).Return(mocks.City, nil).Times(1)

		repo.Province.EXPECT().FindByID(ctx, ids.ProvinceID).Return(mocks.Province, nil).Times(1)

		repo.Region.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, r *domain.Region) error {
			r.ID = ids.RegionID
			return nil
		}).Times(1)

		repo.Region.EXPECT().FindByID(ctx, ids.RegionID).Return(mocks.Region, nil).Times(1)

		resp, err := uc.CreateRegion(ctx, dtos.Create)

		assert.NoError(t, err)
		assert.Equal(t, ids.RegionID, resp.ID)
		assert.Equal(t, ids.CityID, resp.CityID)
		assert.Equal(t, ids.ProvinceID, resp.ProvinceID)
	})

	t.Run("should return error validation", func(t *testing.T) {
		req := &dto.RegionCreateDto{
			CityID:     ids.CityID,
			ProvinceID: 0,
		}
		_, err := uc.CreateRegion(ctx, req)
		assert.Error(t, err)
		assert.EqualError(t, err, "Validation failed")
	})

	t.Run("should return error when city not found", func(t *testing.T) {
		repo.City.EXPECT().FindByID(ctx, ids.CityID).Return(nil, utils.NewNotFoundError("city not found")).Times(1)

		resp, err := uc.CreateRegion(ctx, dtos.Create)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "city not found")
	})

	t.Run("should return error when province not found", func(t *testing.T) {
		repo.City.EXPECT().FindByID(ctx, ids.CityID).Return(mocks.City, nil).Times(1)

		repo.Province.EXPECT().FindByID(ctx, ids.ProvinceID).Return(nil, utils.NewNotFoundError("province not found")).Times(1)

		resp, err := uc.CreateRegion(ctx, dtos.Create)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "province not found")
	})

	t.Run("should return error when create region", func(t *testing.T) {
		repo.City.EXPECT().FindByID(ctx, ids.CityID).Return(mocks.City, nil).Times(1)

		repo.Province.EXPECT().FindByID(ctx, ids.ProvinceID).Return(mocks.Province, nil).Times(1)

		repo.Region.EXPECT().Create(ctx, gomock.Any()).Return(utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.CreateRegion(ctx, dtos.Create)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})

	t.Run("should return error when get created region", func(t *testing.T) {
		repo.City.EXPECT().FindByID(ctx, ids.CityID).Return(mocks.City, nil).Times(1)

		repo.Province.EXPECT().FindByID(ctx, ids.ProvinceID).Return(mocks.Province, nil).Times(1)

		repo.Region.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, r *domain.Region) error {
			r.ID = ids.RegionID
			return nil
		}).Times(1)

		repo.Region.EXPECT().FindByID(ctx, ids.RegionID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.CreateRegion(ctx, dtos.Create)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

func TestGetAllRegions(t *testing.T) {
	_, mocks, _, repo, uc, ctx := RegionUsecaseUtils(t)
	t.Run("should return all regions", func(t *testing.T) {
		repo.Region.EXPECT().FindAll(ctx).Return(mocks.Regions, nil).Times(1)

		res, err := uc.GetAllRegions(ctx)

		assert.NotNil(t, res)
		assert.NoError(t, err)
		assert.Len(t, *res, len(*mocks.Regions))
		assert.Equal(t, (*mocks.Regions)[0].ID, (*res)[0].ID)
		assert.Equal(t, (*mocks.Regions)[0].CityID, (*res)[0].CityID)
		assert.Equal(t, (*mocks.Regions)[0].ProvinceID, (*res)[0].ProvinceID)
	})

	t.Run("should return error when get all regions", func(t *testing.T) {
		repo.Region.EXPECT().FindAll(ctx).Return(nil, utils.NewInternalError("internal error")).Times(1)

		_, err := uc.GetAllRegions(ctx)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

func TestGetRegionByID(t *testing.T) {
	ids, mocks, _, repo, uc, ctx := RegionUsecaseUtils(t)
	t.Run("should return region by id success", func(t *testing.T) {
		repo.Region.EXPECT().FindByID(ctx, ids.RegionID).Return(mocks.Region, nil).Times(1)

		res, err := uc.GetRegionByID(ctx, ids.RegionID)
		assert.NoError(t, err)
		assert.Equal(t, ids.RegionID, res.ID)
		assert.Equal(t, ids.CityID, res.CityID)
		assert.Equal(t, ids.ProvinceID, res.ProvinceID)
	})

	t.Run("should return error when get region by id", func(t *testing.T) {
		repo.Region.EXPECT().FindByID(ctx, ids.RegionID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		_, err := uc.GetRegionByID(ctx, ids.RegionID)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}
