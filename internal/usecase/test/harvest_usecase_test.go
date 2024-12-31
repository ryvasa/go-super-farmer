package usecase_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	mock_repo "github.com/ryvasa/go-super-farmer/internal/repository/mock"
	usecase_implementation "github.com/ryvasa/go-super-farmer/internal/usecase/implementation"
	usecase_interface "github.com/ryvasa/go-super-farmer/internal/usecase/interface"
	"github.com/ryvasa/go-super-farmer/pkg/env"
	mock_pkg "github.com/ryvasa/go-super-farmer/pkg/mock"
	"github.com/ryvasa/go-super-farmer/utils"
	mock_utils "github.com/ryvasa/go-super-farmer/utils/mock"
	"github.com/stretchr/testify/assert"
)

type HarvestRepoMock struct {
	Harvest       *mock_repo.MockHarvestRepository
	City          *mock_repo.MockCityRepository
	LandCommodity *mock_repo.MockLandCommodityRepository
	RabbitMQ      *mock_pkg.MockRabbitMQ
	Cache         *mock_pkg.MockCache
	Glob          *mock_utils.MockGlobFunc
	TxManager     *mock_pkg.MockTransactionManager
}

type HarvestIDs struct {
	HarvestID       uuid.UUID
	LandCommodityID uuid.UUID
	CityID          int64
	LandID          uuid.UUID
	CommodityID     uuid.UUID
}

type HarvestDomainMock struct {
	Harvest         *domain.Harvest
	Harvests        []*domain.Harvest
	UpdatedHarvest  *domain.Harvest
	City            *domain.City
	LandCommodity   *domain.LandCommodity
	LandCommodities []*domain.LandCommodity
	Message         usecase_implementation.HarvestMessage
}

type HarvestDTOMock struct {
	Create *dto.HarvestCreateDTO
	Update *dto.HarvestUpdateDTO
	Params *dto.HarvestParamsDTO
}

func HarvestUsecaseSetup(t *testing.T) (*HarvestIDs, *HarvestDomainMock, *HarvestDTOMock, *HarvestRepoMock, usecase_interface.HarvestUsecase, context.Context) {
	cityID := int64(1)
	landCommodityID := uuid.New()
	harvestID := uuid.New()
	commodityID := uuid.New()
	landID := uuid.New()

	ids := &HarvestIDs{
		HarvestID:       harvestID,
		LandCommodityID: landCommodityID,
		CityID:          cityID,
		LandID:          landID,
		CommodityID:     commodityID,
	}
	date, _ := time.Parse("2006-01-02", "2022-01-01")
	startDate, _ := time.Parse("2006-01-02", "2020-01-01")
	endDate, _ := time.Parse("2006-01-02", "2021-01-01-01")

	domains := &HarvestDomainMock{
		Harvest: &domain.Harvest{
			ID:              harvestID,
			LandCommodityID: landCommodityID,
			Quantity:        float64(100),
			Unit:            "kg",
			HarvestDate:     date,
		},
		Harvests: []*domain.Harvest{
			{
				ID:              harvestID,
				LandCommodityID: landCommodityID,
				Quantity:        float64(100),
				Unit:            "kg",
				HarvestDate:     date,
			},
		},
		UpdatedHarvest: &domain.Harvest{
			ID:              harvestID,
			LandCommodityID: landCommodityID,
			Quantity:        float64(99),
			Unit:            "kg",
			HarvestDate:     date,
		},
		City: &domain.City{
			ID: cityID,
		},
		LandCommodity: &domain.LandCommodity{
			ID: landCommodityID,
		},
		LandCommodities: []*domain.LandCommodity{
			{
				ID: landCommodityID,
			},
		},
		Message: usecase_implementation.HarvestMessage{
			LandCommodityID: landCommodityID,
			StartDate:       startDate,
			EndDate:         endDate,
		},
	}

	dto := &HarvestDTOMock{
		Create: &dto.HarvestCreateDTO{
			LandCommodityID: landCommodityID,
			Quantity:        float64(100),
			Unit:            "kg",
			HarvestDate:     "2022-01-01",
		},
		Update: &dto.HarvestUpdateDTO{
			Quantity:    99,
			HarvestDate: "2022-01-01",
		},
		Params: &dto.HarvestParamsDTO{
			LandCommodityID: landCommodityID,
			StartDate:       startDate,
			EndDate:         endDate,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cityRepo := mock_repo.NewMockCityRepository(ctrl)
	landCommodityRepo := mock_repo.NewMockLandCommodityRepository(ctrl)
	harvestRepo := mock_repo.NewMockHarvestRepository(ctrl)
	rabbitMQ := mock_pkg.NewMockRabbitMQ(ctrl)
	cache := mock_pkg.NewMockCache(ctrl)
	glob := mock_utils.NewMockGlobFunc(ctrl)
	env := env.Env{}
	txRepo := mock_pkg.NewMockTransactionManager(ctrl)

	uc := usecase_implementation.NewHarvestUsecase(harvestRepo, cityRepo, landCommodityRepo, rabbitMQ, cache, glob, &env, txRepo)
	ctx := context.TODO()

	repo := &HarvestRepoMock{Harvest: harvestRepo, City: cityRepo, LandCommodity: landCommodityRepo, RabbitMQ: rabbitMQ, Cache: cache, Glob: glob}

	return ids, domains, dto, repo, uc, ctx
}

func TestHarvestRepository_CreateHarvest(t *testing.T) {
	ids, domains, dtos, repo, uc, ctx := HarvestUsecaseSetup(t)

	t.Run("should create harvest successfully", func(t *testing.T) {
		repo.City.EXPECT().FindByID(ctx, ids.CityID).Return(domains.City, nil).Times(1)
		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(domains.LandCommodity, nil).Times(1)

		repo.Harvest.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, p *domain.Harvest) error {
			p.ID = ids.HarvestID
			return nil
		}).Times(1)

		repo.Harvest.EXPECT().FindByID(ctx, ids.HarvestID).Return(domains.Harvest, nil).Times(1)

		repo.Cache.EXPECT().DeleteByPattern(ctx, "harvest").Return(nil).Times(1)

		resp, err := uc.CreateHarvest(ctx, dtos.Create)

		assert.NoError(t, err)
		assert.Equal(t, ids.HarvestID, resp.ID)
		assert.Equal(t, ids.LandCommodityID, resp.LandCommodityID)
		assert.Equal(t, float64(100), resp.Quantity)
		assert.Equal(t, "kg", resp.Unit)
		assert.Equal(t, domains.Harvest.HarvestDate, resp.HarvestDate)
		assert.Equal(t, domains.Harvest.Quantity, resp.Quantity)
		assert.Equal(t, domains.Harvest.Unit, resp.Unit)
	})

	t.Run("should return error validation error", func(t *testing.T) {
		resp, err := uc.CreateHarvest(ctx, &dto.HarvestCreateDTO{
			LandCommodityID: ids.LandCommodityID,
			Quantity:        -10,
		})

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "Validation failed")
	})

	t.Run("should return error when land commodity not found", func(t *testing.T) {

		repo.City.EXPECT().FindByID(ctx, ids.CityID).Return(domains.City, nil).Times(1)

		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(nil, utils.NewNotFoundError("land commodity not found")).Times(1)

		resp, err := uc.CreateHarvest(ctx, dtos.Create)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "land commodity not found")
	})

	t.Run("should return error when city not found", func(t *testing.T) {
		repo.City.EXPECT().FindByID(ctx, ids.CityID).Return(nil, utils.NewNotFoundError("city not found")).Times(1)

		resp, err := uc.CreateHarvest(ctx, dtos.Create)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "city not found")
	})

	t.Run("should return error when create harvest fails", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(domains.LandCommodity, nil).Times(1)

		repo.City.EXPECT().FindByID(ctx, ids.CityID).Return(domains.City, nil).Times(1)

		repo.Harvest.EXPECT().Create(ctx, gomock.Any()).Return(utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.CreateHarvest(ctx, dtos.Create)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})

	t.Run("should return error when get created harvest ", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(domains.LandCommodity, nil).Times(1)

		repo.City.EXPECT().FindByID(ctx, ids.CityID).Return(domains.City, nil).Times(1)

		repo.Harvest.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, p *domain.Harvest) error {
			p.ID = ids.HarvestID
			return nil
		}).Times(1)

		repo.Harvest.EXPECT().FindByID(ctx, ids.HarvestID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.CreateHarvest(ctx, dtos.Create)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})

	t.Run("should return error parsing harvest date", func(t *testing.T) {
		repo.City.EXPECT().FindByID(ctx, ids.CityID).Return(domains.City, nil).Times(1)
		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(domains.LandCommodity, nil).Times(1)

		resp, err := uc.CreateHarvest(ctx, &dto.HarvestCreateDTO{
			LandCommodityID: ids.LandCommodityID,
			HarvestDate:     "string",
			Quantity:        float64(10),
			Unit:            "kg",
		})

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "harvest date format is invalid")
	})

	t.Run("should return error when delete cache fails", func(t *testing.T) {
		repo.City.EXPECT().FindByID(ctx, ids.CityID).Return(domains.City, nil).Times(1)
		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(domains.LandCommodity, nil).Times(1)

		repo.Harvest.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, p *domain.Harvest) error {
			p.ID = ids.HarvestID
			return nil
		}).Times(1)

		repo.Harvest.EXPECT().FindByID(ctx, ids.HarvestID).Return(domains.Harvest, nil).Times(1)

		repo.Cache.EXPECT().DeleteByPattern(ctx, "harvest").Return(utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.CreateHarvest(ctx, dtos.Create)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

func TestHarvestUsecase_GetAllHarvest(t *testing.T) {
	_, domains, _, repo, uc, ctx := HarvestUsecaseSetup(t)

	key := fmt.Sprintf("harvest_%s", "all")
	t.Run("should get all harvests successfully from repo", func(t *testing.T) {
		// Setup expectations
		repo.Cache.EXPECT().Get(ctx, key).Return(nil, nil)
		repo.Harvest.EXPECT().FindAll(ctx).Return(domains.Harvests, nil)

		// Expect cache set to be called with any byte array and return nil
		repo.Cache.EXPECT().
			Set(ctx, key, gomock.Any(), 4*time.Minute).
			Return(nil)

		// Execute
		resp, err := uc.GetAllHarvest(ctx)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, len(domains.Harvests), len(resp))
		assert.Equal(t, (domains.Harvests)[0].ID, (resp)[0].ID)
	})

	t.Run("should return error when get all harvests fails", func(t *testing.T) {
		repo.Cache.EXPECT().Get(ctx, key).Return(nil, nil)
		repo.Harvest.EXPECT().FindAll(ctx).Return(nil, utils.NewInternalError("internal error"))

		resp, err := uc.GetAllHarvest(ctx)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})

	t.Run("should return harvests from cache when cache hit", func(t *testing.T) {
		// Setup cached data
		cachedHarvests, err := json.Marshal(domains.Harvests)
		assert.NoError(t, err)

		// Expect cache get to return the cached data
		repo.Cache.EXPECT().Get(ctx, key).Return(cachedHarvests, nil)

		resp, err := uc.GetAllHarvest(ctx)

		assert.NoError(t, err)
		assert.Equal(t, len(domains.Harvests), len(resp))
		assert.Equal(t, (domains.Harvests)[0].ID, (resp)[0].ID)
	})

	t.Run("should return error when cache set fails", func(t *testing.T) {
		repo.Cache.EXPECT().Get(ctx, key).Return(nil, nil)
		repo.Harvest.EXPECT().FindAll(ctx).Return(domains.Harvests, nil)
		repo.Cache.EXPECT().
			Set(ctx, key, gomock.Any(), 4*time.Minute).
			Return(fmt.Errorf("cache error"))

		resp, err := uc.GetAllHarvest(ctx)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cache error")
	})
}

func TestHarvestUsecase_GetHarvestByID(t *testing.T) {
	ids, domains, _, repo, uc, ctx := HarvestUsecaseSetup(t)

	t.Run("should get harvest by id successfully", func(t *testing.T) {
		repo.Harvest.EXPECT().FindByID(ctx, ids.HarvestID).Return(domains.Harvest, nil).Times(1)

		resp, err := uc.GetHarvestByID(ctx, ids.HarvestID)

		assert.NoError(t, err)
		assert.Equal(t, ids.HarvestID, resp.ID)
	})

	t.Run("should return error when get harvest by id not found", func(t *testing.T) {
		repo.Harvest.EXPECT().FindByID(ctx, ids.HarvestID).Return(nil, utils.NewNotFoundError("harvest not found")).Times(1)

		resp, err := uc.GetHarvestByID(ctx, ids.HarvestID)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "harvest not found")
	})
}

func TestHarvestUsecase_GetHarvestByCommodityID(t *testing.T) {
	ids, domains, _, repo, uc, ctx := HarvestUsecaseSetup(t)

	t.Run("should get harvests by commodity id successfully", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(domains.LandCommodity, nil).Times(1)

		repo.Harvest.EXPECT().FindByCommodityID(ctx, ids.LandCommodityID).Return(domains.Harvests, nil).Times(1)

		resp, err := uc.GetHarvestByCommodityID(ctx, ids.LandCommodityID)

		assert.NoError(t, err)
		assert.Equal(t, len(domains.Harvests), len(resp))
		assert.Equal(t, (domains.Harvests)[0].ID, (resp)[0].ID)
	})

	t.Run("should return error when get harvests by commodity id fails", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(domains.LandCommodity, nil).Times(1)

		repo.Harvest.EXPECT().FindByCommodityID(ctx, ids.LandCommodityID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.GetHarvestByCommodityID(ctx, ids.LandCommodityID)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

func TestHarvestUsecase_GetHarvestByLandID(t *testing.T) {
	ids, domains, _, repo, uc, ctx := HarvestUsecaseSetup(t)

	t.Run("should get harvests by land id successfully", func(t *testing.T) {
		repo.Harvest.EXPECT().FindByLandID(ctx, ids.LandID).Return(domains.Harvests, nil).Times(1)

		resp, err := uc.GetHarvestByLandID(ctx, ids.LandID)

		assert.NoError(t, err)
		assert.Equal(t, len(domains.Harvests), len(resp))
		assert.Equal(t, (domains.Harvests)[0].ID, (resp)[0].ID)
	})

	t.Run("should return error when get harvests by land id fails", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandID).Return(domains.LandCommodity, nil).Times(1)

		repo.Harvest.EXPECT().FindByLandID(ctx, ids.LandID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.GetHarvestByLandID(ctx, ids.LandID)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

func TestHarvestUsecase_GetHarvestByLandCommodityID(t *testing.T) {
	ids, domains, _, repo, uc, ctx := HarvestUsecaseSetup(t)

	t.Run("should get harvests by land commodity id successfully", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(domains.LandCommodity, nil).Times(1)

		repo.Harvest.EXPECT().FindByLandCommodityID(ctx, ids.LandCommodityID).Return(domains.Harvests, nil).Times(1)

		resp, err := uc.GetHarvestByLandCommodityID(ctx, ids.LandCommodityID)

		assert.NoError(t, err)
		assert.Equal(t, len(domains.Harvests), len(resp))
		assert.Equal(t, (domains.Harvests)[0].ID, (resp)[0].ID)
	})

	t.Run("should return error when get harvests by land commodity id fails", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(domains.LandCommodity, nil).Times(1)

		repo.Harvest.EXPECT().FindByLandCommodityID(ctx, ids.LandCommodityID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.GetHarvestByLandCommodityID(ctx, ids.LandCommodityID)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

func TestHarvestUsecase_GetHarvestByCityID(t *testing.T) {
	ids, domains, _, repo, uc, ctx := HarvestUsecaseSetup(t)

	t.Run("should get harvests by city id successfully", func(t *testing.T) {
		repo.City.EXPECT().FindByID(ctx, ids.CityID).Return(domains.City, nil).Times(1)

		repo.Harvest.EXPECT().FindByCityID(ctx, ids.CityID).Return(domains.Harvests, nil).Times(1)

		resp, err := uc.GetHarvestByCityID(ctx, ids.CityID)

		assert.NoError(t, err)
		assert.Equal(t, len(domains.Harvests), len(resp))
		assert.Equal(t, (domains.Harvests)[0].ID, (resp)[0].ID)
	})

	t.Run("should return error when get harvests by city id fails", func(t *testing.T) {
		repo.City.EXPECT().FindByID(ctx, ids.CityID).Return(domains.City, nil).Times(1)

		repo.Harvest.EXPECT().FindByCityID(ctx, ids.CityID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.GetHarvestByCityID(ctx, ids.CityID)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

func TestHarvestUsecase_UpdateHarvest(t *testing.T) {
	ids, domains, dtos, repo, uc, ctx := HarvestUsecaseSetup(t)

	t.Run("should update harvest successfully", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(domains.LandCommodity, nil).Times(1)

		repo.Harvest.EXPECT().FindByID(ctx, ids.HarvestID).Return(domains.Harvest, nil).Times(1)

		repo.Harvest.EXPECT().Update(ctx, ids.HarvestID, gomock.Any()).DoAndReturn(func(ctx context.Context, id uuid.UUID, p *domain.Harvest) error {
			p.ID = ids.HarvestID
			return nil
		}).Times(1)

		repo.Harvest.EXPECT().FindByID(ctx, ids.HarvestID).Return(domains.UpdatedHarvest, nil).Times(1)

		repo.Cache.EXPECT().DeleteByPattern(ctx, "harvest").Return(nil).Times(1)

		resp, err := uc.UpdateHarvest(ctx, ids.HarvestID, dtos.Update)

		assert.NoError(t, err)
		assert.Equal(t, ids.HarvestID, resp.ID)
		assert.Equal(t, float64(99), resp.Quantity)
		assert.Equal(t, "kg", resp.Unit)
	})

	t.Run("should return error when harvest not found", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(domains.LandCommodity, nil).Times(1)

		repo.Harvest.EXPECT().FindByID(ctx, ids.HarvestID).Return(nil, utils.NewNotFoundError("harvest not found")).Times(1)

		resp, err := uc.UpdateHarvest(ctx, ids.HarvestID, dtos.Update)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "harvest not found")
	})

	t.Run("should return error validation error", func(t *testing.T) {
		resp, err := uc.UpdateHarvest(ctx, ids.HarvestID, &dto.HarvestUpdateDTO{
			Quantity: -10,
		})

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "Validation failed")
	})

	t.Run("should return error when update harvest fails", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(domains.LandCommodity, nil).Times(1)

		repo.Harvest.EXPECT().FindByID(ctx, ids.HarvestID).Return(domains.Harvest, nil).Times(1)

		repo.Harvest.EXPECT().Update(ctx, ids.HarvestID, gomock.Any()).Return(utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.UpdateHarvest(ctx, ids.HarvestID, dtos.Update)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})

	t.Run("should return error when get updated harvest error", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(domains.LandCommodity, nil).Times(1)

		repo.Harvest.EXPECT().FindByID(ctx, ids.HarvestID).Return(domains.Harvest, nil).Times(1)

		repo.Harvest.EXPECT().Update(ctx, ids.HarvestID, gomock.Any()).DoAndReturn(func(ctx context.Context, id uuid.UUID, p *domain.Harvest) error {
			p.ID = ids.HarvestID
			return nil
		}).Times(1)

		repo.Harvest.EXPECT().FindByID(ctx, ids.HarvestID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.UpdateHarvest(ctx, ids.HarvestID, dtos.Update)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})

	t.Run("should return error when date format is invalid", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(domains.LandCommodity, nil).Times(1)

		repo.Harvest.EXPECT().FindByID(ctx, ids.HarvestID).Return(domains.Harvest, nil).Times(1)

		repo.Harvest.EXPECT().Update(ctx, ids.HarvestID, gomock.Any()).DoAndReturn(func(ctx context.Context, id uuid.UUID, p *domain.Harvest) error {
			p.ID = ids.HarvestID
			return nil
		}).Times(1)

		repo.Harvest.EXPECT().FindByID(ctx, ids.HarvestID).Return(domains.UpdatedHarvest, nil).Times(1)

		resp, err := uc.UpdateHarvest(ctx, ids.HarvestID, &dto.HarvestUpdateDTO{
			Quantity:    99,
			HarvestDate: "string",
		})

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "Validation failed")
	})

	t.Run("should return error when delete cache fails", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(domains.LandCommodity, nil).Times(1)

		repo.Harvest.EXPECT().FindByID(ctx, ids.HarvestID).Return(domains.Harvest, nil).Times(1)

		repo.Harvest.EXPECT().Update(ctx, ids.HarvestID, gomock.Any()).DoAndReturn(func(ctx context.Context, id uuid.UUID, p *domain.Harvest) error {
			p.ID = ids.HarvestID
			return nil
		}).Times(1)

		repo.Harvest.EXPECT().FindByID(ctx, ids.HarvestID).Return(domains.UpdatedHarvest, nil).Times(1)

		repo.Cache.EXPECT().DeleteByPattern(ctx, "harvest").Return(utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.UpdateHarvest(ctx, ids.HarvestID, dtos.Update)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

func TestHarvestUsecase_DeleteHarvest(t *testing.T) {
	ids, domains, _, repo, uc, ctx := HarvestUsecaseSetup(t)

	t.Run("should delete harvest successfully", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(domains.LandCommodity, nil).Times(1)

		repo.Harvest.EXPECT().FindByID(ctx, ids.HarvestID).Return(domains.Harvest, nil).Times(1)

		repo.Harvest.EXPECT().Delete(ctx, ids.HarvestID).Return(nil).Times(1)

		repo.Cache.EXPECT().DeleteByPattern(ctx, "harvest").Return(nil).Times(1)

		err := uc.DeleteHarvest(ctx, ids.HarvestID)

		assert.NoError(t, err)
	})

	t.Run("should return error when harvest not found", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(domains.LandCommodity, nil).Times(1)

		repo.Harvest.EXPECT().FindByID(ctx, ids.HarvestID).Return(nil, utils.NewNotFoundError("harvest not found")).Times(1)

		err := uc.DeleteHarvest(ctx, ids.HarvestID)

		assert.Error(t, err)
		assert.EqualError(t, err, "harvest not found")
	})

	t.Run("should return error when delete harvest fails", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(domains.LandCommodity, nil).Times(1)

		repo.Harvest.EXPECT().FindByID(ctx, ids.HarvestID).Return(domains.Harvest, nil).Times(1)

		repo.Harvest.EXPECT().Delete(ctx, ids.HarvestID).Return(utils.NewInternalError("internal error")).Times(1)

		err := uc.DeleteHarvest(ctx, ids.HarvestID)

		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})

	t.Run("should return error when delete cache fails", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(domains.LandCommodity, nil).Times(1)

		repo.Harvest.EXPECT().FindByID(ctx, ids.HarvestID).Return(domains.Harvest, nil).Times(1)

		repo.Harvest.EXPECT().Delete(ctx, ids.HarvestID).Return(nil).Times(1)

		repo.Cache.EXPECT().DeleteByPattern(ctx, "harvest").Return(utils.NewInternalError("internal error")).Times(1)

		err := uc.DeleteHarvest(ctx, ids.HarvestID)

		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

func TestHarvestUsecase_RestoreHarvest(t *testing.T) {
	ids, domains, _, repo, uc, ctx := HarvestUsecaseSetup(t)

	t.Run("should restore harvest successfully", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(domains.LandCommodity, nil).Times(1)

		repo.Harvest.EXPECT().FindDeletedByID(ctx, ids.HarvestID).Return(domains.Harvest, nil).Times(1)

		repo.Harvest.EXPECT().Restore(ctx, ids.HarvestID).Return(nil).Times(1)

		repo.Harvest.EXPECT().FindByID(ctx, ids.HarvestID).Return(domains.Harvest, nil).Times(1)

		repo.Cache.EXPECT().DeleteByPattern(ctx, "harvest").Return(nil).Times(1)

		resp, err := uc.RestoreHarvest(ctx, ids.HarvestID)

		assert.NoError(t, err)
		assert.Equal(t, ids.HarvestID, resp.ID)
	})

	t.Run("should return error when harvest not found", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(domains.LandCommodity, nil).Times(1)

		repo.Harvest.EXPECT().FindDeletedByID(ctx, ids.HarvestID).Return(nil, utils.NewNotFoundError("deleted harvest not found")).Times(1)

		resp, err := uc.RestoreHarvest(ctx, ids.HarvestID)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "deleted harvest not found")
	})

	t.Run("should return error when restore harvest fails", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(domains.LandCommodity, nil).Times(1)

		repo.Harvest.EXPECT().FindDeletedByID(ctx, ids.HarvestID).Return(domains.Harvest, nil).Times(1)

		repo.Harvest.EXPECT().Restore(ctx, ids.HarvestID).Return(utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.RestoreHarvest(ctx, ids.HarvestID)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})

	t.Run("should return error when get restored harvest", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(domains.LandCommodity, nil).Times(1)

		repo.Harvest.EXPECT().FindDeletedByID(ctx, ids.HarvestID).Return(domains.Harvest, nil).Times(1)

		repo.Harvest.EXPECT().Restore(ctx, ids.HarvestID).Return(utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.RestoreHarvest(ctx, ids.HarvestID)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})

	t.Run("should return error when delete cache fails", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(domains.LandCommodity, nil).Times(1)

		repo.Harvest.EXPECT().FindDeletedByID(ctx, ids.HarvestID).Return(domains.Harvest, nil).Times(1)

		repo.Harvest.EXPECT().Restore(ctx, ids.HarvestID).Return(nil).Times(1)

		repo.Harvest.EXPECT().FindByID(ctx, ids.HarvestID).Return(domains.Harvest, nil).Times(1)

		repo.Cache.EXPECT().DeleteByPattern(ctx, "harvest").Return(utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.RestoreHarvest(ctx, ids.HarvestID)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})

	t.Run("should return error when get restored harvest fails", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(domains.LandCommodity, nil).Times(1)

		repo.Harvest.EXPECT().FindDeletedByID(ctx, ids.HarvestID).Return(domains.Harvest, nil).Times(1)

		repo.Harvest.EXPECT().Restore(ctx, ids.HarvestID).Return(nil).Times(1)

		repo.Harvest.EXPECT().FindByID(ctx, ids.HarvestID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.RestoreHarvest(ctx, ids.HarvestID)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

func TestHarvestUsecase_GetAllDeletedHarvest(t *testing.T) {
	_, domains, _, repo, uc, ctx := HarvestUsecaseSetup(t)

	t.Run("should get all deleted harvests successfully", func(t *testing.T) {
		repo.Harvest.EXPECT().FindAllDeleted(ctx).Return(domains.Harvests, nil).Times(1)

		resp, err := uc.GetAllDeletedHarvest(ctx)

		assert.NoError(t, err)
		assert.Equal(t, len(domains.Harvests), len(resp))
		assert.Equal(t, (domains.Harvests)[0].ID, (resp)[0].ID)
	})

	t.Run("should return error when get all deleted harvests fails", func(t *testing.T) {
		repo.Harvest.EXPECT().FindAllDeleted(ctx).Return(nil, utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.GetAllDeletedHarvest(ctx)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

func TestHarvestUsecase_GetHarvestDeletedByID(t *testing.T) {
	ids, domains, _, repo, uc, ctx := HarvestUsecaseSetup(t)

	t.Run("should get deleted harvest by id successfully", func(t *testing.T) {
		repo.Harvest.EXPECT().FindDeletedByID(ctx, ids.HarvestID).Return(domains.Harvest, nil).Times(1)

		resp, err := uc.GetHarvestDeletedByID(ctx, ids.HarvestID)

		assert.NoError(t, err)
		assert.Equal(t, ids.HarvestID, resp.ID)
	})

	t.Run("should return error when get deleted harvest by id fails", func(t *testing.T) {
		repo.Harvest.EXPECT().FindDeletedByID(ctx, ids.HarvestID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.GetHarvestDeletedByID(ctx, ids.HarvestID)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "deleted harvest not found")
	})
}

func TestHarvestUsecase_DownloadHarvestByLandCommodityID(t *testing.T) {
	_, domains, dtos, repo, uc, ctx := HarvestUsecaseSetup(t)

	t.Run("should publish harvests to rabbitmq successfully", func(t *testing.T) {
		repo.Harvest.EXPECT().FindByLandCommodityID(context.TODO(), dtos.Params.LandCommodityID).Return(domains.Harvests, nil).Times(1)
		repo.RabbitMQ.EXPECT().
			PublishJSON(ctx, "report-exchange", "harvest", domains.Message).
			Return(nil)

		res, err := uc.DownloadHarvestByLandCommodityID(ctx, dtos.Params)

		url := fmt.Sprintf("http://localhost:8081/harvests/land_commodity/%s/download/file?start_date=%s&end_date=%s",
			dtos.Params.LandCommodityID, dtos.Params.StartDate.Format("2006-01-02"), dtos.Params.EndDate.Format("2006-01-02"))

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, res.Message, "Report generation in progress. Please check back in a few moments.")
		assert.Equal(t, url, res.DownloadURL)
	})

	t.Run("should return error when publish harvests to rabbitmq fails", func(t *testing.T) {
		repo.Harvest.EXPECT().FindByLandCommodityID(context.TODO(), dtos.Params.LandCommodityID).Return(domains.Harvests, nil).Times(1)

		repo.RabbitMQ.EXPECT().
			PublishJSON(ctx, "report-exchange", "harvest", domains.Message).
			Return(utils.NewInternalError("internal error"))

		resp, err := uc.DownloadHarvestByLandCommodityID(ctx, dtos.Params)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})

	t.Run("should return error when get harvests by land commodity id fails", func(t *testing.T) {
		repo.Harvest.EXPECT().FindByLandCommodityID(context.TODO(), dtos.Params.LandCommodityID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.DownloadHarvestByLandCommodityID(ctx, dtos.Params)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

// func TestHarvestUsecase_GetHarvestExcelFile(t *testing.T) {
// 	ids, domains, _, repo, uc, ctx := HarvestUsecaseSetup(t)

// 	// Buat direktori ./public/reports jika belum ada
// 	reportsDir := "./public/reports"
// 	err := os.MkdirAll(reportsDir, 0755) // 0755 adalah permission untuk direktori
// 	if err != nil {
// 		t.Fatalf("Failed to create reports directory: %v", err)
// 	}
// 	defer os.RemoveAll(reportsDir) // Hapus direktori setelah tes selesai

// 	// Buat file dummy Excel di ./public/reports
// 	file := fmt.Sprintf("harvests_%s_%s_%s_*.xlsx",
// 		ids.LandCommodityID,
// 		domains.Harvests[0].HarvestDate.Format("2006-01-02"),
// 		domains.Harvests[0].HarvestDate.Format("2006-01-02"))

// 	dummyFilePath := fmt.Sprintf("%s/%s", reportsDir, file)
// 	err = os.WriteFile(dummyFilePath, []byte("Dummy Excel content"), 0644)
// 	if err != nil {
// 		t.Fatalf("Failed to create dummy Excel file: %v", err)
// 	}

// 	t.Run("should get harvest excel file successfully", func(t *testing.T) {

// 		repo.Glob.EXPECT().Glob(dummyFilePath).Return([]string{dummyFilePath}, nil)

// 		resp, err := uc.GetHarvestExcelFile(ctx, &dto.HarvestParamsDTO{
// 			LandCommodityID: ids.LandCommodityID,
// 			StartDate:       domains.Harvests[0].HarvestDate,
// 			EndDate:         domains.Harvests[0].HarvestDate,
// 		})

// 		assert.NoError(t, err)
// 		assert.NotNil(t, resp)
// 		assert.Equal(t, dummyFilePath, *resp)
// 	})

// 	t.Run("should return error when glob fails", func(t *testing.T) {
// 		repo.Glob.EXPECT().Glob(dummyFilePath).Return(nil, utils.NewInternalError("Error finding report file"))

// 		resp, err := uc.GetHarvestExcelFile(ctx, &dto.HarvestParamsDTO{
// 			LandCommodityID: ids.LandCommodityID,
// 			StartDate:       domains.Harvests[0].HarvestDate,
// 			EndDate:         domains.Harvests[0].HarvestDate,
// 		})

// 		assert.Nil(t, resp)
// 		assert.Error(t, err)
// 		assert.EqualError(t, err, "Error finding report file")
// 	})

// 	t.Run("should return error when no matching files", func(t *testing.T) {

// 		repo.Glob.EXPECT().Glob(dummyFilePath).Return([]string{}, nil)

// 		resp, err := uc.GetHarvestExcelFile(ctx, &dto.HarvestParamsDTO{
// 			LandCommodityID: ids.LandCommodityID,
// 			StartDate:       domains.Harvests[0].HarvestDate,
// 			EndDate:         domains.Harvests[0].HarvestDate,
// 		})

// 		logrus.Info(resp)
// 		logrus.Info(err)
// 		assert.Nil(t, resp)
// 		assert.Error(t, err)
// 		assert.EqualError(t, err, "Report file not found")
// 	})
// }
