package usecase_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	mock_pkg "github.com/ryvasa/go-super-farmer/pkg/mock"
	"github.com/ryvasa/go-super-farmer/service_api/model/domain"
	"github.com/ryvasa/go-super-farmer/service_api/model/dto"
	mock_repo "github.com/ryvasa/go-super-farmer/service_api/repository/mock"
	usecase_implementation "github.com/ryvasa/go-super-farmer/service_api/usecase/implementation"
	usecase_interface "github.com/ryvasa/go-super-farmer/service_api/usecase/interface"
	"github.com/ryvasa/go-super-farmer/utils"
	"github.com/stretchr/testify/assert"
)

type LandCommodityRepoMock struct {
	LandCommodity *mock_repo.MockLandCommodityRepository
	Land          *mock_repo.MockLandRepository
	Commodity     *mock_repo.MockCommodityRepository
	Cache         *mock_pkg.MockCache
}

type LandCommodityIDs struct {
	CommodityID     uuid.UUID
	LandID          uuid.UUID
	LandCommodityID uuid.UUID
}

type LandCommodityMocks struct {
	LandCommodity        *domain.LandCommodity
	LandCommodities      []*domain.LandCommodity
	UpdatedLandCommodity *domain.LandCommodity
	Land                 *domain.Land
	Commodity            *domain.Commodity
}

type LandCommodityDTOMocks struct {
	Create *dto.LandCommodityCreateDTO
	Update *dto.LandCommodityUpdateDTO
}

func LandCommodityUtils(t *testing.T) (*LandCommodityIDs, *LandCommodityMocks, *LandCommodityDTOMocks, *LandCommodityRepoMock, usecase_interface.LandCommodityUsecase, context.Context) {
	landID := uuid.New()
	commodityID := uuid.New()
	landCommodityID := uuid.New()

	ids := &LandCommodityIDs{
		CommodityID:     commodityID,
		LandID:          landID,
		LandCommodityID: landCommodityID,
	}

	mocks := &LandCommodityMocks{
		LandCommodity: &domain.LandCommodity{
			ID:          landCommodityID,
			CommodityID: commodityID,
			LandID:      landID,
			LandArea:    float64(100),
		},
		LandCommodities: []*domain.LandCommodity{
			{
				ID:          landCommodityID,
				CommodityID: commodityID,
				LandID:      landID,
				LandArea:    float64(100),
			},
			{
				ID:          landCommodityID,
				CommodityID: commodityID,
				// LandID:      landID,
				//LandArea:    float64(100),
			},
		},
		UpdatedLandCommodity: &domain.LandCommodity{
			ID:          landCommodityID,
			CommodityID: commodityID,
			LandID:      landID,
			LandArea:    float64(200),
		},
		Land: &domain.Land{
			ID:       landID,
			LandArea: float64(1000),
		},
		Commodity: &domain.Commodity{
			ID: commodityID,
		},
	}

	dtoMocks := &LandCommodityDTOMocks{
		Create: &dto.LandCommodityCreateDTO{
			LandID:      landID,
			CommodityID: commodityID,
			LandArea:    float64(100),
		},
		Update: &dto.LandCommodityUpdateDTO{
			LandID:      landID,
			CommodityID: commodityID,
			LandArea:    float64(200),
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	landCommodity := mock_repo.NewMockLandCommodityRepository(ctrl)
	land := mock_repo.NewMockLandRepository(ctrl)
	commodity := mock_repo.NewMockCommodityRepository(ctrl)
	city := mock_repo.NewMockCityRepository(ctrl)
	cache := mock_pkg.NewMockCache(ctrl)

	repoMock := &LandCommodityRepoMock{
		LandCommodity: landCommodity,
		Land:          land,
		Commodity:     commodity,
		Cache:         cache,
	}

	uc := usecase_implementation.NewLandCommodityUsecase(landCommodity, land, city, commodity, cache)
	ctx := context.Background()

	return ids, mocks, dtoMocks, repoMock, uc, ctx
}
func TestLandCommodityUsecase_CreateLandCommodity(t *testing.T) {
	ids, mocks, dtos, repo, uc, ctx := LandCommodityUtils(t)
	t.Run("should create land commodity successfully", func(t *testing.T) {
		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(mocks.Commodity, nil).Times(1)

		repo.Land.EXPECT().FindByID(ctx, ids.LandID).Return(mocks.Land, nil).Times(1)

		repo.LandCommodity.EXPECT().SumLandAreaByLandID(ctx, ids.LandID).Return(float64(10), nil).Times(1)

		repo.LandCommodity.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, l *domain.LandCommodity) error {
			l.ID = ids.LandCommodityID
			return nil
		}).Times(1)

		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(mocks.LandCommodity, nil).Times(1)

		resp, err := uc.CreateLandCommodity(ctx, dtos.Create)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, ids.LandCommodityID, resp.ID)
	})

	t.Run("should return error when land area is greater than max area", func(t *testing.T) {

		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(mocks.Commodity, nil).Times(1)

		repo.Land.EXPECT().FindByID(ctx, ids.LandID).Return(mocks.Land, nil).Times(1)

		repo.LandCommodity.EXPECT().SumLandAreaByLandID(ctx, ids.LandID).Return(float64(1000), nil).Times(1)

		req := &dto.LandCommodityCreateDTO{LandID: ids.LandID, CommodityID: ids.CommodityID, LandArea: float64(100)}
		resp, err := uc.CreateLandCommodity(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "land area not enough")
	})

	t.Run("should return error validation error", func(t *testing.T) {
		resp, err := uc.CreateLandCommodity(ctx, &dto.LandCommodityCreateDTO{LandID: ids.LandID, CommodityID: ids.CommodityID, LandArea: 0})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "Validation failed")
	})

	t.Run("should return error when commodity not found", func(t *testing.T) {
		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(nil, utils.NewNotFoundError("commodity not found")).Times(1)

		resp, err := uc.CreateLandCommodity(ctx, dtos.Create)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "commodity not found")
	})

	t.Run("should return error when land not found", func(t *testing.T) {
		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(mocks.Commodity, nil).Times(1)

		repo.Land.EXPECT().FindByID(ctx, ids.LandID).Return(nil, utils.NewNotFoundError("land not found")).Times(1)

		resp, err := uc.CreateLandCommodity(ctx, dtos.Create)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "land not found")
	})
}

func TestLandCommodityUsecase_GetLandCommodityByID(t *testing.T) {
	ids, mocks, _, repo, uc, ctx := LandCommodityUtils(t)

	t.Run("should return error when land commodity not found", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(nil, utils.NewNotFoundError("land commodity not found")).Times(1)

		resp, err := uc.GetLandCommodityByID(ctx, ids.LandCommodityID)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "land commodity not found")
	})
	t.Run("should return land commodity successfully", func(t *testing.T) {

		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(mocks.LandCommodity, nil).Times(1)

		resp, err := uc.GetLandCommodityByID(ctx, ids.LandCommodityID)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, ids.LandCommodityID, resp.ID)
		assert.Equal(t, ids.CommodityID, resp.CommodityID)
		assert.Equal(t, ids.LandID, resp.LandID)
		assert.Equal(t, float64(100), resp.LandArea)
	})
}

func TestLandCommodityUsecase_GetLandCommodityByLandID(t *testing.T) {
	ids, mocks, _, repo, uc, ctx := LandCommodityUtils(t)

	t.Run("should return land commodity successfully", func(t *testing.T) {

		repo.LandCommodity.EXPECT().FindByLandID(ctx, ids.LandID).Return(mocks.LandCommodities, nil).Times(1)

		resp, err := uc.GetLandCommodityByLandID(ctx, ids.LandID)

		assert.NoError(t, err)
		assert.Len(t, resp, 2)
		assert.Equal(t, (mocks.LandCommodities)[0].LandArea, (resp)[0].LandArea)
	})

	t.Run("should return error when service_api error", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByLandID(ctx, ids.LandID).Return(nil, utils.NewInternalError("service_api error")).Times(1)

		resp, err := uc.GetLandCommodityByLandID(ctx, ids.LandID)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "service_api error")
	})
}

func TestLandCommodityUsecase_GetLandCommodityByCommodityID(t *testing.T) {
	ids, mocks, _, repo, uc, ctx := LandCommodityUtils(t)

	t.Run("should return land commodity successfully", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByCommodityID(ctx, ids.CommodityID).Return(mocks.LandCommodities, nil).Times(1)

		resp, err := uc.GetLandCommodityByCommodityID(ctx, ids.CommodityID)

		assert.NoError(t, err)
		assert.Len(t, resp, 2)
		assert.Equal(t, (mocks.LandCommodities)[0].LandArea, (resp)[0].LandArea)
	})

	t.Run("should return error when service_api error", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByCommodityID(ctx, ids.CommodityID).Return(nil, utils.NewInternalError("service_api error")).Times(1)

		resp, err := uc.GetLandCommodityByCommodityID(ctx, ids.CommodityID)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "service_api error")
	})
}

func TestLandCommodityUsecase_GetAllLandCommodity(t *testing.T) {
	_, mocks, _, repo, uc, ctx := LandCommodityUtils(t)

	key := fmt.Sprintf("land_commodity_%s", "all")
	t.Run("should get all land commodity successfully from repo", func(t *testing.T) {
		// Setup expectations
		repo.Cache.EXPECT().Get(ctx, key).Return(nil, nil)
		repo.LandCommodity.EXPECT().FindAll(ctx).Return(mocks.LandCommodities, nil)

		// Expect cache set to be called with any byte array and return nil
		repo.Cache.EXPECT().
			Set(ctx, key, gomock.Any(), 4*time.Minute).
			Return(nil)

		// Execute
		resp, err := uc.GetAllLandCommodity(ctx)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, len(mocks.LandCommodities), len(resp))
		assert.Equal(t, (mocks.LandCommodities)[0].ID, (resp)[0].ID)
	})

	t.Run("should return error when get all land commodity fails", func(t *testing.T) {
		repo.Cache.EXPECT().Get(ctx, key).Return(nil, nil)
		repo.LandCommodity.EXPECT().FindAll(ctx).Return(nil, utils.NewInternalError("service_api error"))

		resp, err := uc.GetAllLandCommodity(ctx)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "service_api error")
	})

	t.Run("should return land commodity from cache when cache hit", func(t *testing.T) {
		// Setup cached data
		cachedLandCommodities, err := json.Marshal(mocks.LandCommodities)
		assert.NoError(t, err)

		// Expect cache get to return the cached data
		repo.Cache.EXPECT().Get(ctx, key).Return(cachedLandCommodities, nil)

		resp, err := uc.GetAllLandCommodity(ctx)

		assert.NoError(t, err)
		assert.Equal(t, len(mocks.LandCommodities), len(resp))
		assert.Equal(t, (mocks.LandCommodities)[0].ID, (resp)[0].ID)
	})

	t.Run("should return error when cache set fails", func(t *testing.T) {
		repo.Cache.EXPECT().Get(ctx, key).Return(nil, nil)
		repo.LandCommodity.EXPECT().FindAll(ctx).Return(mocks.LandCommodities, nil)
		repo.Cache.EXPECT().
			Set(ctx, key, gomock.Any(), 4*time.Minute).
			Return(fmt.Errorf("cache error"))

		resp, err := uc.GetAllLandCommodity(ctx)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cache error")
	})
}

func TestLandCommodityUsecase_UpdateLandCommodity(t *testing.T) {
	ids, mocks, dtos, repo, uc, ctx := LandCommodityUtils(t)

	t.Run("should return error when validation error", func(t *testing.T) {
		resp, err := uc.UpdateLandCommodity(ctx, ids.LandCommodityID, &dto.LandCommodityUpdateDTO{LandArea: 0})

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("should return error when land commodity not found", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(nil, utils.NewNotFoundError("land commodity not found")).Times(1)

		resp, err := uc.UpdateLandCommodity(ctx, ids.LandCommodityID, dtos.Update)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "land commodity not found")
	})

	t.Run("should return error when commodity not found", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(mocks.LandCommodity, nil).Times(1)

		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(nil, utils.NewNotFoundError("commodity not found")).Times(1)

		resp, err := uc.UpdateLandCommodity(ctx, ids.LandCommodityID, dtos.Update)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "commodity not found")
	})

	t.Run("should return error when land not found", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(mocks.LandCommodity, nil).Times(1)

		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(mocks.Commodity, nil).Times(1)

		repo.Land.EXPECT().FindByID(ctx, ids.LandID).Return(nil, utils.NewNotFoundError("land not found")).Times(1)

		resp, err := uc.UpdateLandCommodity(ctx, ids.LandCommodityID, dtos.Update)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "land not found")
	})

	t.Run("should return error when land area not enough", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(mocks.LandCommodity, nil).Times(1)

		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(mocks.Commodity, nil).Times(1)

		repo.Land.EXPECT().FindByID(ctx, ids.LandID).Return(mocks.Land, nil).Times(1)

		repo.LandCommodity.EXPECT().SumLandAreaByLandID(ctx, ids.LandID).Return(float64(1000), nil).Times(1)

		resp, err := uc.UpdateLandCommodity(ctx, ids.LandCommodityID, dtos.Update)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "land area not enough")
	})

	t.Run("should return error when service_api error", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(mocks.LandCommodity, nil).Times(1)

		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(mocks.Commodity, nil).Times(1)

		repo.Land.EXPECT().FindByID(ctx, ids.LandID).Return(mocks.Land, nil).Times(1)

		repo.LandCommodity.EXPECT().SumLandAreaByLandID(ctx, ids.LandID).Return(float64(100), nil).Times(1)

		repo.LandCommodity.EXPECT().Update(ctx, ids.LandCommodityID, mocks.LandCommodity).Return(utils.NewInternalError("service_api error")).Times(1)

		resp, err := uc.UpdateLandCommodity(ctx, ids.LandCommodityID, dtos.Update)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "service_api error")
	})

	t.Run("should update land commodity successfully", func(t *testing.T) {

		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(mocks.LandCommodity, nil).Times(1)

		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(mocks.Commodity, nil).Times(1)

		repo.Land.EXPECT().FindByID(ctx, ids.LandID).Return(mocks.Land, nil).Times(1)

		repo.LandCommodity.EXPECT().SumLandAreaByLandID(ctx, ids.LandID).Return(float64(100), nil).Times(1)

		repo.LandCommodity.EXPECT().Update(ctx, ids.LandCommodityID, mocks.LandCommodity).Return(nil).Times(1)

		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(mocks.UpdatedLandCommodity, nil).Times(1)

		resp, err := uc.UpdateLandCommodity(ctx, ids.LandCommodityID, dtos.Update)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, ids.LandCommodityID, resp.ID)
		assert.Equal(t, float64(200), resp.LandArea)
	})
}

func TestLandCommodityUsecase_DeleteLandCommodity(t *testing.T) {
	ids, mocks, _, repo, uc, ctx := LandCommodityUtils(t)

	t.Run("should delete land commodity successfully", func(t *testing.T) {

		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(mocks.LandCommodity, nil).Times(1)

		repo.LandCommodity.EXPECT().Delete(ctx, ids.LandCommodityID).Return(nil).Times(1)

		err := uc.DeleteLandCommodity(ctx, ids.LandCommodityID)

		assert.NoError(t, err)
	})

	t.Run("should return error when land commodity not found", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(nil, utils.NewNotFoundError("land commodity not found")).Times(1)

		err := uc.DeleteLandCommodity(ctx, ids.LandCommodityID)

		assert.Error(t, err)
		assert.EqualError(t, err, "land commodity not found")
	})

	t.Run("should return error when service_api error", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(mocks.LandCommodity, nil).Times(1)

		repo.LandCommodity.EXPECT().Delete(ctx, ids.LandCommodityID).Return(utils.NewInternalError("database error")).Times(1)

		err := uc.DeleteLandCommodity(ctx, ids.LandCommodityID)

		assert.Error(t, err)
		assert.EqualError(t, err, "database error")
	})
}

func TestLandCommodityUsecase_RestoreLandCommodity(t *testing.T) {
	ids, mocks, _, repo, uc, ctx := LandCommodityUtils(t)

	t.Run("should restore land commodity successfully", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindDeletedByID(ctx, ids.LandCommodityID).Return(mocks.LandCommodity, nil).Times(1)
		repo.Land.EXPECT().FindByID(ctx, ids.LandID).Return(mocks.Land, nil).Times(1)
		repo.LandCommodity.EXPECT().SumLandAreaByLandID(ctx, ids.LandID).Return(float64(100), nil).Times(1)

		repo.LandCommodity.EXPECT().Restore(ctx, ids.LandCommodityID).Return(nil).Times(1)

		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(mocks.LandCommodity, nil).Times(1)

		resp, err := uc.RestoreLandCommodity(ctx, ids.LandCommodityID)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, ids.LandCommodityID, resp.ID)
	})

	t.Run("should return error when deleted land commodity not found", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindDeletedByID(ctx, ids.LandCommodityID).Return(nil, utils.NewNotFoundError("deleted land commodity not found")).Times(1)

		resp, err := uc.RestoreLandCommodity(ctx, ids.LandCommodityID)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "deleted land commodity not found")
	})

	t.Run("should return error when service_api error", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindDeletedByID(ctx, ids.LandCommodityID).Return(mocks.LandCommodity, nil).Times(1)
		repo.Land.EXPECT().FindByID(ctx, ids.LandID).Return(mocks.Land, nil).Times(1)

		repo.LandCommodity.EXPECT().SumLandAreaByLandID(ctx, ids.LandID).Return(float64(100), nil).Times(1)

		repo.LandCommodity.EXPECT().Restore(ctx, ids.LandCommodityID).Return(utils.NewInternalError("service_api error")).Times(1)

		resp, err := uc.RestoreLandCommodity(ctx, ids.LandCommodityID)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "service_api error")
	})

	t.Run("should return error when land area not enough", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindDeletedByID(ctx, ids.LandCommodityID).Return(mocks.LandCommodity, nil).Times(1)

		repo.Land.EXPECT().FindByID(ctx, ids.LandID).Return(mocks.Land, nil).Times(1)

		repo.LandCommodity.EXPECT().SumLandAreaByLandID(ctx, ids.LandID).Return(float64(1000), nil).Times(1)

		resp, err := uc.RestoreLandCommodity(ctx, ids.LandCommodityID)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "land area not enough")
	})

	t.Run("should return error when restored land not found", func(t *testing.T) {

		repo.LandCommodity.EXPECT().FindDeletedByID(ctx, ids.LandCommodityID).Return(mocks.LandCommodity, nil).Times(1)
		repo.Land.EXPECT().FindByID(ctx, ids.LandID).Return(mocks.Land, nil).Times(1)

		repo.LandCommodity.EXPECT().SumLandAreaByLandID(ctx, ids.LandID).Return(float64(100), nil).Times(1)

		repo.LandCommodity.EXPECT().Restore(ctx, ids.LandCommodityID).Return(nil).Times(1)

		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(nil, utils.NewNotFoundError("restored land not found")).Times(1)

		resp, err := uc.RestoreLandCommodity(ctx, ids.LandCommodityID)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "restored land not found")
	})
}
