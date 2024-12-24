package usecase_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/pkg/env"
	mock_pkg "github.com/ryvasa/go-super-farmer/pkg/mock"
	"github.com/ryvasa/go-super-farmer/service_api/model/domain"
	"github.com/ryvasa/go-super-farmer/service_api/model/dto"
	mock_repo "github.com/ryvasa/go-super-farmer/service_api/repository/mock"
	usecase_implementation "github.com/ryvasa/go-super-farmer/service_api/usecase/implementation"
	usecase_interface "github.com/ryvasa/go-super-farmer/service_api/usecase/interface"
	"github.com/ryvasa/go-super-farmer/utils"
	mock_utils "github.com/ryvasa/go-super-farmer/utils/mock"
	"github.com/stretchr/testify/assert"
)

type PriceRepoMock struct {
	Price        *mock_repo.MockPriceRepository
	City         *mock_repo.MockCityRepository
	Commodity    *mock_repo.MockCommodityRepository
	PriceHistory *mock_repo.MockPriceHistoryRepository
	TxManager    *mock_pkg.MockTransactionManager
	RabbitMQ     *mock_pkg.MockRabbitMQ
	Cache        *mock_pkg.MockCache
	Glob         *mock_utils.MockGlobFunc
}

type PriceIDs struct {
	PriceID        uuid.UUID
	PriceHistoryID uuid.UUID
	CommodityID    uuid.UUID
	CityID         int64
}

type Pricemocks struct {
	Prices        []*domain.Price
	Price         *domain.Price
	UpdatedPrice  *domain.Price
	HistoryPrices []*domain.PriceHistory
	HistoryPrice  *domain.PriceHistory
	Commodity     *domain.Commodity
	City          *domain.City
	Message       usecase_implementation.PriceMessage
}

type PriceDTOmocks struct {
	Create *dto.PriceCreateDTO
	Update *dto.PriceUpdateDTO
	Params *dto.PriceParamsDTO
}

func PriceUsecaseUtils(t *testing.T) (*PriceIDs, *Pricemocks, *PriceDTOmocks, *PriceRepoMock, usecase_interface.PriceUsecase, context.Context) {
	cityID := int64(1)
	commodityID := uuid.New()
	priceID := uuid.New()
	priceHistoryID := uuid.New()

	ids := &PriceIDs{
		PriceID:        priceID,
		PriceHistoryID: priceHistoryID,
		CommodityID:    commodityID,
		CityID:         cityID,
	}
	startDate, _ := time.Parse("2006-01-02", "2020-01-01")
	endDate, _ := time.Parse("2006-01-02", "2021-01-01-01")
	dtos := &PriceDTOmocks{
		Create: &dto.PriceCreateDTO{
			CommodityID: commodityID,
			CityID:      cityID,
			Price:       100,
		},
		Update: &dto.PriceUpdateDTO{
			Price: 100,
		},
		Params: &dto.PriceParamsDTO{
			CommodityID: commodityID,
			CityID:      cityID,
			StartDate:   startDate,
			EndDate:     endDate,
		},
	}

	mocks := &Pricemocks{
		Prices: []*domain.Price{
			{
				ID:          priceID,
				CommodityID: commodityID,
				CityID:      cityID,
				Price:       100,
			},
		},
		Price: &domain.Price{
			ID:          priceID,
			CommodityID: commodityID,
			CityID:      cityID,
			Price:       100,
		},
		UpdatedPrice: &domain.Price{
			ID:          priceID,
			CommodityID: commodityID,
			CityID:      cityID,
			Price:       900,
		},
		HistoryPrices: []*domain.PriceHistory{
			{
				ID:          priceID,
				CommodityID: commodityID,
				CityID:      cityID,
				Price:       100,
			},
		},
		HistoryPrice: &domain.PriceHistory{
			ID:          priceID,
			CommodityID: commodityID,
			CityID:      cityID,
			Price:       100,
		},
		Commodity: &domain.Commodity{
			ID:   commodityID,
			Name: "string",
		},
		City: &domain.City{
			ID: cityID,
		},
		Message: usecase_implementation.PriceMessage{
			CommodityID: dtos.Params.CommodityID,
			CityID:      dtos.Params.CityID,
			StartDate:   dtos.Params.StartDate,
			EndDate:     dtos.Params.EndDate,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cityRepo := mock_repo.NewMockCityRepository(ctrl)
	commodityRepo := mock_repo.NewMockCommodityRepository(ctrl)
	priceRepo := mock_repo.NewMockPriceRepository(ctrl)
	priceHostoryRepo := mock_repo.NewMockPriceHistoryRepository(ctrl)
	txRepo := mock_pkg.NewMockTransactionManager(ctrl)
	rabbitMQ := mock_pkg.NewMockRabbitMQ(ctrl)
	cache := mock_pkg.NewMockCache(ctrl)
	glob := mock_utils.NewMockGlobFunc(ctrl)
	env := env.Env{}

	uc := usecase_implementation.NewPriceUsecase(priceRepo, priceHostoryRepo, cityRepo, commodityRepo, rabbitMQ, txRepo, cache, glob, &env)
	ctx := context.Background()

	repo := &PriceRepoMock{
		Price:        priceRepo,
		City:         cityRepo,
		Commodity:    commodityRepo,
		PriceHistory: priceHostoryRepo,
		TxManager:    txRepo,
		RabbitMQ:     rabbitMQ,
		Cache:        cache,
		Glob:         glob,
	}

	return ids, mocks, dtos, repo, uc, ctx
}

func TestPriceUsecase_CreatePrice(t *testing.T) {

	ids, mocks, dto, repo, uc, ctx := PriceUsecaseUtils(t)

	t.Run("should create price successfully", func(t *testing.T) {

		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(mocks.Commodity, nil).Times(1)

		repo.City.EXPECT().FindByID(ctx, ids.CityID).Return(mocks.City, nil).Times(1)

		repo.Price.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, p *domain.Price) error {
			p.ID = ids.PriceID
			return nil
		}).Times(1)

		repo.Price.EXPECT().FindByID(ctx, ids.PriceID).Return(mocks.Price, nil).Times(1)

		repo.Cache.EXPECT().DeleteByPattern(ctx, "price").Return(nil).Times(1)

		resp, err := uc.CreatePrice(ctx, dto.Create)

		assert.NoError(t, err)
		assert.Equal(t, dto.Create.CommodityID, resp.CommodityID)
		assert.Equal(t, dto.Create.CityID, resp.CityID)
		assert.Equal(t, mocks.Price.ID, resp.ID)
	})

	t.Run("should return error when commodity not found", func(t *testing.T) {
		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(nil, utils.NewNotFoundError("commodity not found")).Times(1)

		resp, err := uc.CreatePrice(ctx, dto.Create)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "commodity not found")
	})

	t.Run("should return error when city not found", func(t *testing.T) {
		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(mocks.Commodity, nil).Times(1)

		repo.City.EXPECT().FindByID(ctx, ids.CityID).Return(nil, utils.NewNotFoundError("city not found")).Times(1)

		resp, err := uc.CreatePrice(ctx, dto.Create)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "city not found")
	})

	t.Run("should return error when create price", func(t *testing.T) {
		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(mocks.Commodity, nil).Times(1)

		repo.City.EXPECT().FindByID(ctx, ids.CityID).Return(mocks.City, nil).Times(1)

		repo.Price.EXPECT().Create(ctx, gomock.Any()).Return(utils.NewInternalError("service_api error")).Times(1)

		resp, err := uc.CreatePrice(ctx, dto.Create)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "service_api error")
	})

	t.Run("should return error when get created price", func(t *testing.T) {
		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(mocks.Commodity, nil).Times(1)

		repo.City.EXPECT().FindByID(ctx, ids.CityID).Return(mocks.City, nil).Times(1)

		repo.Price.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, p *domain.Price) error {
			p.ID = ids.PriceID
			return nil
		}).Times(1)

		repo.Price.EXPECT().FindByID(ctx, ids.PriceID).Return(nil, utils.NewInternalError("service_api error")).Times(1)

		resp, err := uc.CreatePrice(ctx, dto.Create)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "service_api error")
	})
}

func TestPriceUsecase_GetAllPrices(t *testing.T) {
	_, mocks, _, repo, uc, ctx := PriceUsecaseUtils(t)

	queryParams := &dto.PaginationDTO{
		Limit:  10,
		Page:   1,
		Filter: dto.PaginationFilterDTO{},
	}

	cacheKey := fmt.Sprintf("price_list_page_%d_limit_%d",
		queryParams.Page,
		queryParams.Limit,
	)

	t.Run("should return prices from cache when cache hit", func(t *testing.T) {
		// Setup
		mockedResponse := &dto.PaginationResponseDTO{
			TotalRows:  10,
			TotalPages: 1,
			Page:       1,
			Limit:      10,
			Data:       mocks.Prices,
		}
		cachedJSON, err := json.Marshal(mockedResponse)
		assert.NoError(t, err)

		// Mock expectations
		repo.Cache.EXPECT().Get(ctx, cacheKey).Return(cachedJSON, nil)

		// Execute
		resp, err := uc.GetAllPrices(ctx, queryParams)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, mockedResponse.TotalRows, resp.TotalRows)
		assert.Equal(t, mockedResponse.TotalPages, resp.TotalPages)
		assert.Equal(t, mockedResponse.Page, resp.Page)
		assert.Equal(t, mockedResponse.Limit, resp.Limit)
		assert.Equal(t, mockedResponse.Data, resp.Data)
	})

	t.Run("should return prices from repository when cache miss", func(t *testing.T) {
		// Mock expectations in order
		repo.Cache.EXPECT().Get(ctx, cacheKey).Return(nil, nil)
		repo.Price.EXPECT().FindAll(ctx, queryParams).Return(mocks.Prices, nil)
		repo.Price.EXPECT().Count(ctx, &queryParams.Filter).Return(int64(10), nil).Times(1)
		repo.Cache.EXPECT().Set(ctx, cacheKey, gomock.Any(), gomock.Any()).Return(nil).Times(1)
		// Execute
		resp, err := uc.GetAllPrices(ctx, queryParams)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, len(mocks.Prices), len(resp.Data.([]*domain.Price)))
		assert.Equal(t, int64(10), resp.TotalRows)
	})

	t.Run("should return validation error", func(t *testing.T) {
		invalidQueryParams := &dto.PaginationDTO{
			Limit: -1, // Invalid limit
			Page:  1,
		}

		resp, err := uc.GetAllPrices(ctx, invalidQueryParams)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "limit must be greater than 0")
	})

	t.Run("should return error when cache is invalid", func(t *testing.T) {
		repo.Cache.EXPECT().Get(ctx, cacheKey).Return([]byte("invalid data"), nil).Times(1)

		resp, err := uc.GetAllPrices(ctx, queryParams)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("should return error when database call fails", func(t *testing.T) {
		repo.Cache.EXPECT().Get(ctx, cacheKey).Return(nil, nil).Times(1)
		repo.Price.EXPECT().FindAll(ctx, queryParams).Return(nil, utils.NewInternalError("db error")).Times(1)

		resp, err := uc.GetAllPrices(ctx, queryParams)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "db error")
	})

	t.Run("should return error when cache set fails", func(t *testing.T) {
		// Mock expectations
		repo.Cache.EXPECT().Get(ctx, cacheKey).Return(nil, nil)
		repo.Price.EXPECT().FindAll(ctx, queryParams).Return(mocks.Prices, nil)
		repo.Price.EXPECT().Count(ctx, &queryParams.Filter).Return(int64(10), nil).Times(1)
		repo.Cache.EXPECT().Set(ctx, cacheKey, gomock.Any(), 4*time.Minute).Return(fmt.Errorf("cache set error"))

		// Execute
		resp, err := uc.GetAllPrices(ctx, queryParams)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "cache set error")
	})
}

func TestPriceUsecase_GetPriceByID(t *testing.T) {
	ids, mocks, _, repo, uc, ctx := PriceUsecaseUtils(t)

	t.Run("should return price by id success", func(t *testing.T) {
		repo.Price.EXPECT().FindByID(ctx, ids.PriceID).Return(mocks.Price, nil).Times(1)

		res, err := uc.GetPriceByID(ctx, ids.PriceID)
		assert.NoError(t, err)
		assert.Equal(t, ids.PriceID, res.ID)
		assert.Equal(t, ids.CommodityID, res.CommodityID)
		assert.Equal(t, ids.CityID, res.CityID)
		assert.Equal(t, mocks.Price.Price, res.Price)
	})

	t.Run("should return error when get price by id", func(t *testing.T) {
		repo.Price.EXPECT().FindByID(ctx, ids.PriceID).Return(nil, utils.NewInternalError("service_api error")).Times(1)

		_, err := uc.GetPriceByID(ctx, ids.PriceID)
		assert.Error(t, err)
		assert.EqualError(t, err, "service_api error")
	})
}

func TestPriceUsecase_GetPricesByCommodityID(t *testing.T) {
	ids, mocks, _, repo, uc, ctx := PriceUsecaseUtils(t)

	t.Run("should return price by commodity id success", func(t *testing.T) {
		repo.Price.EXPECT().FindByCommodityID(ctx, ids.CommodityID).Return(mocks.Prices, nil).Times(1)

		res, err := uc.GetPricesByCommodityID(ctx, ids.CommodityID)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(res))
		assert.Equal(t, (mocks.Prices)[0].ID, (res)[0].ID)
		assert.Equal(t, (mocks.Prices)[0].CommodityID, (res)[0].CommodityID)
		assert.Equal(t, (mocks.Prices)[0].CityID, (res)[0].CityID)
		assert.Equal(t, (mocks.Prices)[0].Price, (res)[0].Price)
	})

	t.Run("should return error when get price by commodity id", func(t *testing.T) {
		repo.Price.EXPECT().FindByCommodityID(ctx, ids.CommodityID).Return(nil, utils.NewInternalError("service_api error")).Times(1)

		_, err := uc.GetPricesByCommodityID(ctx, ids.CommodityID)
		assert.Error(t, err)
		assert.EqualError(t, err, "service_api error")
	})
}

func TestPriceUsecase_GetPricesByCityID(t *testing.T) {
	ids, mocks, _, repo, uc, ctx := PriceUsecaseUtils(t)

	t.Run("should return price by city id success", func(t *testing.T) {
		repo.Price.EXPECT().FindByCityID(ctx, ids.CityID).Return(mocks.Prices, nil).Times(1)

		res, err := uc.GetPricesByCityID(ctx, ids.CityID)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(res))
		assert.Equal(t, (mocks.Prices)[0].ID, (res)[0].ID)
		assert.Equal(t, (mocks.Prices)[0].CommodityID, (res)[0].CommodityID)
		assert.Equal(t, (mocks.Prices)[0].CityID, (res)[0].CityID)
		assert.Equal(t, (mocks.Prices)[0].Price, (res)[0].Price)
	})

	t.Run("should return error when get price by city id", func(t *testing.T) {
		repo.Price.EXPECT().FindByCityID(ctx, ids.CityID).Return(nil, utils.NewInternalError("service_api error")).Times(1)

		_, err := uc.GetPricesByCityID(ctx, ids.CityID)
		assert.Error(t, err)
		assert.EqualError(t, err, "service_api error")
	})
}

func TestPriceUsecase_UpdatePrice(t *testing.T) {
	ids, mocks, dtos, repo, uc, ctx := PriceUsecaseUtils(t)

	t.Run("should update price successfully", func(t *testing.T) {
		repo.TxManager.EXPECT().
			WithTransaction(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
				return fn(ctx)
			})

		repo.Price.EXPECT().FindByID(ctx, ids.PriceID).Return(mocks.Price, nil).Times(1)

		repo.PriceHistory.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, ph *domain.PriceHistory) error {
			ph.ID = ids.PriceHistoryID
			return nil
		}).Times(1)

		repo.Price.EXPECT().Update(ctx, ids.PriceID, gomock.Any()).DoAndReturn(func(ctx context.Context, id uuid.UUID, p *domain.Price) error {
			p.ID = ids.PriceID
			return nil
		}).Times(1)

		repo.Price.EXPECT().FindByID(ctx, ids.PriceID).Return(mocks.UpdatedPrice, nil).Times(1)

		repo.Cache.EXPECT().DeleteByPattern(ctx, "price").Return(nil).Times(1)

		resp, err := uc.UpdatePrice(ctx, ids.PriceID, dtos.Update)

		assert.NoError(t, err)
		assert.Equal(t, mocks.Price.ID, resp.ID)
		assert.Equal(t, resp.Price, mocks.UpdatedPrice.Price)
	})

	t.Run("should return error when price not found", func(t *testing.T) {
		repo.TxManager.EXPECT().
			WithTransaction(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
				return fn(ctx)
			})

		repo.Price.EXPECT().FindByID(ctx, ids.PriceID).Return(nil, utils.NewNotFoundError("price not found")).Times(1)

		resp, err := uc.UpdatePrice(ctx, ids.PriceID, dtos.Update)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "price not found")
	})

	t.Run("should return error when create price history", func(t *testing.T) {
		repo.TxManager.EXPECT().
			WithTransaction(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
				return fn(ctx)
			})

		repo.Price.EXPECT().FindByID(ctx, ids.PriceID).Return(mocks.Price, nil).Times(1)

		repo.PriceHistory.EXPECT().Create(ctx, gomock.Any()).Return(utils.NewInternalError("failed to create price history")).Times(1)

		resp, err := uc.UpdatePrice(ctx, ids.PriceID, dtos.Update)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "failed to create price history")
	})
	t.Run("should return error when update price", func(t *testing.T) {
		repo.TxManager.EXPECT().
			WithTransaction(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
				return fn(ctx)
			})

		repo.Price.EXPECT().FindByID(ctx, ids.PriceID).Return(mocks.Price, nil).Times(1)

		repo.PriceHistory.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, ph *domain.PriceHistory) error {
			ph.ID = ids.PriceHistoryID
			return nil
		}).Times(1)

		repo.Price.EXPECT().Update(ctx, ids.PriceID, gomock.Any()).Return(utils.NewInternalError("service_api error")).Times(1)

		resp, err := uc.UpdatePrice(ctx, ids.PriceID, dtos.Update)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "service_api error")
	})

	t.Run("should return error when get updated price", func(t *testing.T) {
		repo.TxManager.EXPECT().
			WithTransaction(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
				return fn(ctx)
			})

		repo.Price.EXPECT().FindByID(ctx, ids.PriceID).Return(mocks.Price, nil).Times(1)

		repo.PriceHistory.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, ph *domain.PriceHistory) error {
			ph.ID = ids.PriceHistoryID
			return nil
		}).Times(1)

		repo.Price.EXPECT().Update(ctx, ids.PriceID, gomock.Any()).DoAndReturn(func(ctx context.Context, id uuid.UUID, p *domain.Price) error {
			p.ID = ids.PriceID
			return nil
		}).Times(1)

		repo.Price.EXPECT().FindByID(ctx, ids.PriceID).Return(nil, utils.NewInternalError("service_api error")).Times(1)

		resp, err := uc.UpdatePrice(ctx, ids.PriceID, dtos.Update)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "service_api error")
	})
	t.Run("should return error when validation fails", func(t *testing.T) {

		resp, err := uc.UpdatePrice(ctx, ids.PriceID, &dto.PriceUpdateDTO{
			Price: -10,
		})

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "Validation failed")
	})
}

func TestPriceUsecase_DeletePrice(t *testing.T) {
	ids, mocks, _, repo, uc, ctx := PriceUsecaseUtils(t)

	t.Run("should delete price successfully", func(t *testing.T) {
		repo.Price.EXPECT().FindByID(ctx, ids.PriceID).Return(mocks.Price, nil).Times(1)

		repo.Price.EXPECT().Delete(ctx, ids.PriceID).Return(nil).Times(1)

		repo.Cache.EXPECT().DeleteByPattern(ctx, "price").Return(nil).Times(1)

		err := uc.DeletePrice(ctx, ids.PriceID)

		assert.Nil(t, err)
		assert.NoError(t, err)
	})
	t.Run("should return error when price not found", func(t *testing.T) {

		repo.Price.EXPECT().FindByID(ctx, ids.PriceID).Return(nil, utils.NewNotFoundError("price not found")).Times(1)

		err := uc.DeletePrice(ctx, ids.PriceID)

		assert.Error(t, err)
		assert.EqualError(t, err, "price not found")
	})

	t.Run("should return error when delete price", func(t *testing.T) {
		repo.Price.EXPECT().FindByID(ctx, ids.PriceID).Return(mocks.Price, nil).Times(1)

		repo.Price.EXPECT().Delete(ctx, ids.PriceID).Return(utils.NewInternalError("service_api error")).Times(1)

		err := uc.DeletePrice(ctx, ids.PriceID)

		assert.Error(t, err)
		assert.EqualError(t, err, "service_api error")
	})

}

func TestPriceUsecase_RestorePrice(t *testing.T) {
	ids, mocks, _, repo, uc, ctx := PriceUsecaseUtils(t)

	t.Run("should restore price successfully", func(t *testing.T) {
		repo.Price.EXPECT().FindDeletedByID(ctx, ids.PriceID).Return(mocks.Price, nil).Times(1)

		repo.Price.EXPECT().Restore(ctx, ids.PriceID).Return(nil).Times(1)

		repo.Price.EXPECT().FindByID(ctx, ids.PriceID).Return(mocks.Price, nil).Times(1)

		repo.Cache.EXPECT().DeleteByPattern(ctx, "price")

		resp, err := uc.RestorePrice(ctx, ids.PriceID)

		assert.NotNil(t, resp)
		assert.Nil(t, err)
		assert.NoError(t, err)
	})

	t.Run("should return error when restore price", func(t *testing.T) {
		repo.Price.EXPECT().FindDeletedByID(ctx, ids.PriceID).Return(mocks.Price, nil).Times(1)

		repo.Price.EXPECT().Restore(ctx, ids.PriceID).Return(utils.NewInternalError("service_api error")).Times(1)

		resp, err := uc.RestorePrice(ctx, ids.PriceID)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "service_api error")
	})

	t.Run("should return error when get price by id", func(t *testing.T) {
		repo.Price.EXPECT().FindDeletedByID(ctx, ids.PriceID).Return(mocks.Price, nil).Times(1)

		repo.Price.EXPECT().Restore(ctx, ids.PriceID).Return(nil).Times(1)

		repo.Price.EXPECT().FindByID(ctx, ids.PriceID).Return(nil, utils.NewInternalError("service_api error")).Times(1)

		resp, err := uc.RestorePrice(ctx, ids.PriceID)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "service_api error")
	})

	t.Run("should return error when get deleted price by id", func(t *testing.T) {
		repo.Price.EXPECT().FindDeletedByID(ctx, ids.PriceID).Return(nil, utils.NewNotFoundError("price not found")).Times(1)

		resp, err := uc.RestorePrice(ctx, ids.PriceID)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "price not found")
	})
}

func TestPriceUsecase_GetPriceByCommodityIDAndCityID(t *testing.T) {
	ids, mocks, _, repo, uc, ctx := PriceUsecaseUtils(t)

	t.Run("should return price by commodity id and city id success", func(t *testing.T) {
		repo.Price.EXPECT().FindByCommodityIDAndCityID(ctx, ids.CommodityID, ids.CityID).Return(mocks.Price, nil).Times(1)

		res, err := uc.GetPriceByCommodityIDAndCityID(ctx, ids.CommodityID, ids.CityID)
		assert.NoError(t, err)
		assert.Equal(t, ids.PriceID, res.ID)
		assert.Equal(t, ids.CommodityID, res.CommodityID)
		assert.Equal(t, ids.CityID, res.CityID)
		assert.Equal(t, mocks.Price.Price, res.Price)
	})

	t.Run("should return error when get price by commodity id and city id", func(t *testing.T) {
		repo.Price.EXPECT().FindByCommodityIDAndCityID(ctx, ids.CommodityID, ids.CityID).Return(nil, utils.NewInternalError("service_api error")).Times(1)

		_, err := uc.GetPriceByCommodityIDAndCityID(ctx, ids.CommodityID, ids.CityID)
		assert.Error(t, err)
		assert.EqualError(t, err, "service_api error")
	})
}

func TestPriceUsecase_GetPriceHistoryByCommodityIDAndCityID(t *testing.T) {
	ids, mocks, _, repo, uc, ctx := PriceUsecaseUtils(t)
	key := fmt.Sprintf("price_history_%s_%d", ids.CommodityID, ids.CityID)
	t.Run("should return price history by commodity id and city id success by repo", func(t *testing.T) {
		repo.Cache.EXPECT().Get(ctx, key).Return(nil, nil)
		repo.PriceHistory.EXPECT().FindByCommodityIDAndCityID(ctx, ids.CommodityID, ids.CityID).Return(mocks.HistoryPrices, nil).Times(1)

		repo.Price.EXPECT().FindByCommodityIDAndCityID(ctx, ids.CommodityID, ids.CityID).Return(mocks.Price, nil).Times(1)

		repo.Cache.EXPECT().Set(ctx, key, gomock.Any(), 4*time.Minute).Return(nil)

		res, err := uc.GetPriceHistoryByCommodityIDAndCityID(ctx, ids.CommodityID, ids.CityID)

		assert.NoError(t, err)
		assert.Equal(t, 2, len(res))
		assert.Equal(t, ids.PriceID, (res)[0].ID)
		assert.Equal(t, ids.CommodityID, (res)[0].CommodityID)
		assert.Equal(t, ids.CityID, (res)[0].CityID)
		assert.Equal(t, mocks.Price.Price, (res)[0].Price)
	})

	t.Run("should return price from cache when cache hit", func(t *testing.T) {
		// Setup
		expectedResponse := mocks.HistoryPrices
		cachedJSON, err := json.Marshal(expectedResponse)
		assert.NoError(t, err)

		// Mock expectations
		repo.Cache.EXPECT().Get(ctx, key).Return(cachedJSON, nil)

		// Execute
		resp, err := uc.GetPriceHistoryByCommodityIDAndCityID(ctx, ids.CommodityID, ids.CityID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, len(expectedResponse), len(resp))
		for i, price := range resp {
			assert.Equal(t, expectedResponse[i].ID, price.ID)
			assert.Equal(t, expectedResponse[i].CommodityID, price.CommodityID)
			assert.Equal(t, expectedResponse[i].CityID, price.CityID)
			assert.Equal(t, expectedResponse[i].Price, price.Price)
			assert.Equal(t, expectedResponse[i].CreatedAt, price.CreatedAt)
			assert.Equal(t, expectedResponse[i].UpdatedAt, price.UpdatedAt)
			assert.Equal(t, expectedResponse[i].DeletedAt, price.DeletedAt)
		}
	})

	t.Run("should return error when get price history by commodity id and city id", func(t *testing.T) {
		repo.Cache.EXPECT().Get(ctx, key).Return(nil, nil)

		repo.PriceHistory.EXPECT().FindByCommodityIDAndCityID(ctx, ids.CommodityID, ids.CityID).Return(nil, utils.NewInternalError("service_api error")).Times(1)

		_, err := uc.GetPriceHistoryByCommodityIDAndCityID(ctx, ids.CommodityID, ids.CityID)
		assert.Error(t, err)
		assert.EqualError(t, err, "service_api error")
	})

	t.Run("should return error when get price by commodity id and city id", func(t *testing.T) {
		repo.Cache.EXPECT().Get(ctx, key).Return(nil, nil)

		repo.PriceHistory.EXPECT().FindByCommodityIDAndCityID(ctx, ids.CommodityID, ids.CityID).Return(mocks.HistoryPrices, nil).Times(1)

		repo.Price.EXPECT().FindByCommodityIDAndCityID(ctx, ids.CommodityID, ids.CityID).Return(nil, utils.NewInternalError("service_api error")).Times(1)

		_, err := uc.GetPriceHistoryByCommodityIDAndCityID(ctx, ids.CommodityID, ids.CityID)
		assert.Error(t, err)
		assert.EqualError(t, err, "service_api error")
	})

	t.Run("should return error when cache set fails", func(t *testing.T) {
		// Mock expectations
		repo.Cache.EXPECT().Get(ctx, key).Return(nil, nil)
		repo.PriceHistory.EXPECT().FindByCommodityIDAndCityID(ctx, ids.CommodityID, ids.CityID).Return(mocks.HistoryPrices, nil).Times(1)
		repo.Price.EXPECT().FindByCommodityIDAndCityID(ctx, ids.CommodityID, ids.CityID).Return(mocks.Price, nil).Times(1)
		repo.Cache.EXPECT().Set(ctx, key, gomock.Any(), 4*time.Minute).Return(fmt.Errorf("cache set error"))

		// Execute
		resp, err := uc.GetPriceHistoryByCommodityIDAndCityID(ctx, ids.CommodityID, ids.CityID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "cache set error")
	})
}

func TestPriceUsecase_DownloadPriceHistoryByCommodityIDAndCityID(t *testing.T) {
	_, mocks, dtos, repo, uc, ctx := PriceUsecaseUtils(t)

	t.Run("should publish message successfully", func(t *testing.T) {

		// Mock RabbitMQ publish
		repo.RabbitMQ.EXPECT().
			PublishJSON(ctx, "report-exchange", "price-history", mocks.Message).
			Return(nil)

		// Execute
		res, err := uc.DownloadPriceHistoryByCommodityIDAndCityID(ctx, dtos.Params)

		url := fmt.Sprintf("http://localhost:8080/api/prices/history/commodity/%s/city/%d/download/file?start_date=%s&end_date=%s",
			dtos.Params.CommodityID, dtos.Params.CityID, dtos.Params.StartDate.Format("2006-01-02"), dtos.Params.EndDate.Format("2006-01-02"))

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, res.Message, "Price history report generation in progress. Please check back in a few moments.")
		assert.Equal(t, res.DownloadURL, url)
	})

	t.Run("should return error when publish fails", func(t *testing.T) {

		repo.RabbitMQ.EXPECT().
			PublishJSON(ctx, "report-exchange", "price-history", mocks.Message).
			Return(fmt.Errorf("publish error"))

		// Execute
		resp, err := uc.DownloadPriceHistoryByCommodityIDAndCityID(ctx, dtos.Params)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "publish error")
	})

}

// func TestPriceUsecase_GetPriceExcelFile(t *testing.T) {
// 	ids, domains, _, repo, uc, ctx := PriceUsecaseUtils(t)

// 	// Buat direktori ./public/reports jika belum ada
// 	reportsDir := "./public/reports"
// 	err := os.MkdirAll(reportsDir, 0755) // 0755 adalah permission untuk direktori
// 	if err != nil {
// 		t.Fatalf("Failed to create reports directory: %v", err)
// 	}
// 	defer os.RemoveAll(reportsDir) // Hapus direktori setelah tes selesai

// 	// Buat file dummy Excel di ./public/reports
// 	file := fmt.Sprintf("price_history_%s_%d_%s_%s_*.xlsx",
// 		ids.CommodityID,
// 		ids.CityID,
// 		domains.Prices[0].CreatedAt.Format("2006-01-02"), domains.Prices[0].CreatedAt.Format("2006-01-02"),
// 	)

// 	dummyFilePath := fmt.Sprintf("%s/%s", reportsDir, file)
// 	err = os.WriteFile(dummyFilePath, []byte("Dummy Excel content"), 0644)
// 	if err != nil {
// 		t.Fatalf("Failed to create dummy Excel file: %v", err)
// 	}

// 	t.Run("should get price excel file successfully", func(t *testing.T) {

// 		repo.Glob.EXPECT().Glob(dummyFilePath).Return([]string{dummyFilePath}, nil)

// 		resp, err := uc.GetPriceExcelFile(ctx, &dto.PriceParamsDTO{
// 			CommodityID: ids.CommodityID,
// 			CityID:      ids.CityID,
// 			StartDate:   domains.Prices[0].CreatedAt,
// 			EndDate:     domains.Prices[0].CreatedAt,
// 		})

// 		assert.NoError(t, err)
// 		assert.NotNil(t, resp)
// 		assert.Equal(t, dummyFilePath, *resp)
// 	})

// 	t.Run("should return error when glob fails", func(t *testing.T) {
// 		repo.Glob.EXPECT().Glob(dummyFilePath).Return(nil, utils.NewInternalError("Error finding report file"))

// 		resp, err := uc.GetPriceExcelFile(ctx, &dto.PriceParamsDTO{
// 			CommodityID: ids.CommodityID,
// 			CityID:      ids.CityID,
// 			StartDate:   domains.Prices[0].CreatedAt,
// 			EndDate:     domains.Prices[0].CreatedAt,
// 		})

// 		assert.Nil(t, resp)
// 		assert.Error(t, err)
// 		assert.EqualError(t, err, "Error finding report file")
// 	})

// 	t.Run("should return error when no matching files", func(t *testing.T) {

// 		repo.Glob.EXPECT().Glob(dummyFilePath).Return([]string{}, nil)

// 		resp, err := uc.GetPriceExcelFile(ctx, &dto.PriceParamsDTO{
// 			CommodityID: ids.CommodityID,
// 			CityID:      ids.CityID,
// 			StartDate:   domains.Prices[0].CreatedAt,
// 			EndDate:     domains.Prices[0].CreatedAt,
// 		})

// 		logrus.Info(resp)
// 		logrus.Info(err)
// 		assert.Nil(t, resp)
// 		assert.Error(t, err)
// 		assert.EqualError(t, err, "Report file not found")
// 	})
// }
