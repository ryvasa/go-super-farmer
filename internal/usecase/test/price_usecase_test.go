package usecase_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	"github.com/ryvasa/go-super-farmer/internal/repository/mock"
	usecase_implementation "github.com/ryvasa/go-super-farmer/internal/usecase/implementation"
	usecase_interface "github.com/ryvasa/go-super-farmer/internal/usecase/interface"
	"github.com/ryvasa/go-super-farmer/utils"
	"github.com/stretchr/testify/assert"
)

type PriceRepoMock struct {
	Price        *mock.MockPriceRepository
	Region       *mock.MockRegionRepository
	Commodity    *mock.MockCommodityRepository
	PriceHistory *mock.MockPriceHistoryRepository
}

type PriceIDs struct {
	PriceID        uuid.UUID
	PriceHistoryID uuid.UUID
	CommodityID    uuid.UUID
	RegionID       uuid.UUID
}

type Pricemocks struct {
	Prices        *[]domain.Price
	Price         *domain.Price
	UpdatedPrice  *domain.Price
	HistoryPrices *[]domain.PriceHistory
	HistoryPrice  *domain.PriceHistory
	Commodity     *domain.Commodity
	Region        *domain.Region
}

type PriceDTOmocks struct {
	Create *dto.PriceCreateDTO
	Update *dto.PriceUpdateDTO
}

func PriceUsecaseUtils(t *testing.T) (*PriceIDs, *Pricemocks, *PriceDTOmocks, *PriceRepoMock, usecase_interface.PriceUsecase, context.Context) {
	regionID := uuid.New()
	commodityID := uuid.New()
	priceID := uuid.New()
	priceHistoryID := uuid.New()

	ids := &PriceIDs{
		PriceID:        priceID,
		PriceHistoryID: priceHistoryID,
		CommodityID:    commodityID,
		RegionID:       regionID,
	}

	mocks := &Pricemocks{
		Prices: &[]domain.Price{
			{
				ID:          priceID,
				CommodityID: commodityID,
				RegionID:    regionID,
				Price:       100,
			},
		},
		Price: &domain.Price{
			ID:          priceID,
			CommodityID: commodityID,
			RegionID:    regionID,
			Price:       100,
		},
		UpdatedPrice: &domain.Price{
			ID:          priceID,
			CommodityID: commodityID,
			RegionID:    regionID,
			Price:       900,
		},
		HistoryPrices: &[]domain.PriceHistory{
			{
				ID:          priceID,
				CommodityID: commodityID,
				RegionID:    regionID,
				Price:       100,
			},
		},
		HistoryPrice: &domain.PriceHistory{
			ID:          priceID,
			CommodityID: commodityID,
			RegionID:    regionID,
			Price:       100,
		},
		Commodity: &domain.Commodity{
			ID:   commodityID,
			Name: "string",
		},
		Region: &domain.Region{
			ID: regionID,
		},
	}

	dto := &PriceDTOmocks{
		Create: &dto.PriceCreateDTO{
			CommodityID: commodityID,
			RegionID:    regionID,
			Price:       100,
		},
		Update: &dto.PriceUpdateDTO{
			Price: 100,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	regionRepo := mock.NewMockRegionRepository(ctrl)
	commodityRepo := mock.NewMockCommodityRepository(ctrl)
	priceRepo := mock.NewMockPriceRepository(ctrl)
	priceHostoryRepo := mock.NewMockPriceHistoryRepository(ctrl)

	uc := usecase_implementation.NewPriceUsecase(priceRepo, priceHostoryRepo, regionRepo, commodityRepo)
	ctx := context.Background()

	repo := &PriceRepoMock{
		Price:        priceRepo,
		Region:       regionRepo,
		Commodity:    commodityRepo,
		PriceHistory: priceHostoryRepo,
	}

	return ids, mocks, dto, repo, uc, ctx
}

func TestCreatePrice(t *testing.T) {

	ids, mocks, dto, repo, uc, ctx := PriceUsecaseUtils(t)

	t.Run("should create price successfully", func(t *testing.T) {

		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(mocks.Commodity, nil).Times(1)

		repo.Region.EXPECT().FindByID(ctx, ids.RegionID).Return(mocks.Region, nil).Times(1)

		repo.Price.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, p *domain.Price) error {
			p.ID = ids.PriceID
			return nil
		}).Times(1)

		repo.Price.EXPECT().FindByID(ctx, ids.PriceID).Return(mocks.Price, nil).Times(1)

		resp, err := uc.CreatePrice(ctx, dto.Create)

		assert.NoError(t, err)
		assert.Equal(t, dto.Create.CommodityID, resp.CommodityID)
		assert.Equal(t, dto.Create.RegionID, resp.RegionID)
		assert.Equal(t, mocks.Price.ID, resp.ID)
	})

	t.Run("should return error when commodity not found", func(t *testing.T) {
		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(nil, utils.NewNotFoundError("commodity not found")).Times(1)

		resp, err := uc.CreatePrice(ctx, dto.Create)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "commodity not found")
	})

	t.Run("should return error when region not found", func(t *testing.T) {
		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(mocks.Commodity, nil).Times(1)

		repo.Region.EXPECT().FindByID(ctx, ids.RegionID).Return(nil, utils.NewNotFoundError("region not found")).Times(1)

		resp, err := uc.CreatePrice(ctx, dto.Create)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "region not found")
	})

	t.Run("should return error when create price", func(t *testing.T) {
		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(mocks.Commodity, nil).Times(1)

		repo.Region.EXPECT().FindByID(ctx, ids.RegionID).Return(mocks.Region, nil).Times(1)

		repo.Price.EXPECT().Create(ctx, gomock.Any()).Return(utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.CreatePrice(ctx, dto.Create)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})

	t.Run("should return error when get created price", func(t *testing.T) {
		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(mocks.Commodity, nil).Times(1)

		repo.Region.EXPECT().FindByID(ctx, ids.RegionID).Return(mocks.Region, nil).Times(1)

		repo.Price.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, p *domain.Price) error {
			p.ID = ids.PriceID
			return nil
		}).Times(1)

		repo.Price.EXPECT().FindByID(ctx, ids.PriceID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.CreatePrice(ctx, dto.Create)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

func TestGetAllPrices(t *testing.T) {

	_, mocks, _, repo, uc, ctx := PriceUsecaseUtils(t)

	t.Run("should return all prices", func(t *testing.T) {
		repo.Price.EXPECT().FindAll(ctx).Return(mocks.Prices, nil).Times(1)

		res, err := uc.GetAllPrices(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(*res))
		assert.Equal(t, (*mocks.Prices)[0].ID, (*res)[0].ID)
		assert.Equal(t, (*mocks.Prices)[0].CommodityID, (*res)[0].CommodityID)
		assert.Equal(t, (*mocks.Prices)[0].RegionID, (*res)[0].RegionID)
		assert.Equal(t, (*mocks.Prices)[0].Price, (*res)[0].Price)
	})

	t.Run("should return error when get all prices", func(t *testing.T) {
		repo.Price.EXPECT().FindAll(ctx).Return(nil, utils.NewInternalError("internal error")).Times(1)

		_, err := uc.GetAllPrices(ctx)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

func TestGetPriceByID(t *testing.T) {
	ids, mocks, _, repo, uc, ctx := PriceUsecaseUtils(t)

	t.Run("should return price by id success", func(t *testing.T) {
		repo.Price.EXPECT().FindByID(ctx, ids.PriceID).Return(mocks.Price, nil).Times(1)

		res, err := uc.GetPriceByID(ctx, ids.PriceID)
		assert.NoError(t, err)
		assert.Equal(t, ids.PriceID, res.ID)
		assert.Equal(t, ids.CommodityID, res.CommodityID)
		assert.Equal(t, ids.RegionID, res.RegionID)
		assert.Equal(t, mocks.Price.Price, res.Price)
	})

	t.Run("should return error when get price by id", func(t *testing.T) {
		repo.Price.EXPECT().FindByID(ctx, ids.PriceID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		_, err := uc.GetPriceByID(ctx, ids.PriceID)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

func TestGetPricesByCommodityID(t *testing.T) {
	ids, mocks, _, repo, uc, ctx := PriceUsecaseUtils(t)

	t.Run("should return price by commodity id success", func(t *testing.T) {
		repo.Price.EXPECT().FindByCommodityID(ctx, ids.CommodityID).Return(mocks.Prices, nil).Times(1)

		res, err := uc.GetPricesByCommodityID(ctx, ids.CommodityID)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(*res))
		assert.Equal(t, (*mocks.Prices)[0].ID, (*res)[0].ID)
		assert.Equal(t, (*mocks.Prices)[0].CommodityID, (*res)[0].CommodityID)
		assert.Equal(t, (*mocks.Prices)[0].RegionID, (*res)[0].RegionID)
		assert.Equal(t, (*mocks.Prices)[0].Price, (*res)[0].Price)
	})

	t.Run("should return error when get price by commodity id", func(t *testing.T) {
		repo.Price.EXPECT().FindByCommodityID(ctx, ids.CommodityID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		_, err := uc.GetPricesByCommodityID(ctx, ids.CommodityID)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

func TestGetPricesByRegionID(t *testing.T) {
	ids, mocks, _, repo, uc, ctx := PriceUsecaseUtils(t)

	t.Run("should return price by region id success", func(t *testing.T) {
		repo.Price.EXPECT().FindByRegionID(ctx, ids.RegionID).Return(mocks.Prices, nil).Times(1)

		res, err := uc.GetPricesByRegionID(ctx, ids.RegionID)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(*res))
		assert.Equal(t, (*mocks.Prices)[0].ID, (*res)[0].ID)
		assert.Equal(t, (*mocks.Prices)[0].CommodityID, (*res)[0].CommodityID)
		assert.Equal(t, (*mocks.Prices)[0].RegionID, (*res)[0].RegionID)
		assert.Equal(t, (*mocks.Prices)[0].Price, (*res)[0].Price)
	})

	t.Run("should return error when get price by region id", func(t *testing.T) {
		repo.Price.EXPECT().FindByRegionID(ctx, ids.RegionID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		_, err := uc.GetPricesByRegionID(ctx, ids.RegionID)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

func TestUpdatePrice(t *testing.T) {
	ids, mocks, dto, repo, uc, ctx := PriceUsecaseUtils(t)

	t.Run("should update price successfully", func(t *testing.T) {
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

		resp, err := uc.UpdatePrice(ctx, ids.PriceID, dto.Update)

		assert.NoError(t, err)
		assert.Equal(t, mocks.Price.ID, resp.ID)
		assert.Equal(t, resp.Price, mocks.UpdatedPrice.Price)
	})

	t.Run("should return error when price not found", func(t *testing.T) {
		repo.Price.EXPECT().FindByID(ctx, ids.PriceID).Return(nil, utils.NewNotFoundError("price not found")).Times(1)

		resp, err := uc.UpdatePrice(ctx, ids.PriceID, dto.Update)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "price not found")
	})

	t.Run("should return error when create price history", func(t *testing.T) {
		repo.Price.EXPECT().FindByID(ctx, ids.PriceID).Return(mocks.Price, nil).Times(1)

		repo.PriceHistory.EXPECT().Create(ctx, gomock.Any()).Return(utils.NewInternalError("failed to create price history")).Times(1)

		resp, err := uc.UpdatePrice(ctx, ids.PriceID, dto.Update)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "failed to create price history")
	})
	t.Run("should return error when update price", func(t *testing.T) {
		repo.Price.EXPECT().FindByID(ctx, ids.PriceID).Return(mocks.Price, nil).Times(1)

		repo.PriceHistory.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, ph *domain.PriceHistory) error {
			ph.ID = ids.PriceHistoryID
			return nil
		}).Times(1)

		repo.Price.EXPECT().Update(ctx, ids.PriceID, gomock.Any()).Return(utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.UpdatePrice(ctx, ids.PriceID, dto.Update)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})

	t.Run("should return error when get updated price", func(t *testing.T) {
		repo.Price.EXPECT().FindByID(ctx, ids.PriceID).Return(mocks.Price, nil).Times(1)

		repo.PriceHistory.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, ph *domain.PriceHistory) error {
			ph.ID = ids.PriceHistoryID
			return nil
		}).Times(1)

		repo.Price.EXPECT().Update(ctx, ids.PriceID, gomock.Any()).DoAndReturn(func(ctx context.Context, id uuid.UUID, p *domain.Price) error {
			p.ID = ids.PriceID
			return nil
		}).Times(1)

		repo.Price.EXPECT().FindByID(ctx, ids.PriceID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.UpdatePrice(ctx, ids.PriceID, dto.Update)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

func TestDeletePrice(t *testing.T) {
	ids, mocks, _, repo, uc, ctx := PriceUsecaseUtils(t)

	t.Run("should delete price successfully", func(t *testing.T) {
		repo.Price.EXPECT().FindByID(ctx, ids.PriceID).Return(mocks.Price, nil).Times(1)

		repo.Price.EXPECT().Delete(ctx, ids.PriceID).Return(nil).Times(1)

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

		repo.Price.EXPECT().Delete(ctx, ids.PriceID).Return(utils.NewInternalError("internal error")).Times(1)

		err := uc.DeletePrice(ctx, ids.PriceID)

		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})

}

func TestRestorePrice(t *testing.T) {
	ids, mocks, _, repo, uc, ctx := PriceUsecaseUtils(t)

	t.Run("should restore price successfully", func(t *testing.T) {
		repo.Price.EXPECT().FindDeletedByID(ctx, ids.PriceID).Return(mocks.Price, nil).Times(1)

		repo.Price.EXPECT().Restore(ctx, ids.PriceID).Return(nil).Times(1)

		repo.Price.EXPECT().FindByID(ctx, ids.PriceID).Return(mocks.Price, nil).Times(1)

		resp, err := uc.RestorePrice(ctx, ids.PriceID)

		assert.NotNil(t, resp)
		assert.Nil(t, err)
		assert.NoError(t, err)
	})

	t.Run("should return error when restore price", func(t *testing.T) {
		repo.Price.EXPECT().FindDeletedByID(ctx, ids.PriceID).Return(mocks.Price, nil).Times(1)

		repo.Price.EXPECT().Restore(ctx, ids.PriceID).Return(utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.RestorePrice(ctx, ids.PriceID)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})

	t.Run("should return error when get price by id", func(t *testing.T) {
		repo.Price.EXPECT().FindDeletedByID(ctx, ids.PriceID).Return(mocks.Price, nil).Times(1)

		repo.Price.EXPECT().Restore(ctx, ids.PriceID).Return(nil).Times(1)

		repo.Price.EXPECT().FindByID(ctx, ids.PriceID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.RestorePrice(ctx, ids.PriceID)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})

	t.Run("should return error when get deleted price by id", func(t *testing.T) {
		repo.Price.EXPECT().FindDeletedByID(ctx, ids.PriceID).Return(nil, utils.NewNotFoundError("price not found")).Times(1)

		resp, err := uc.RestorePrice(ctx, ids.PriceID)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "price not found")
	})
}

func TestGetPriceByCommodityIDAndRegionID(t *testing.T) {
	ids, mocks, _, repo, uc, ctx := PriceUsecaseUtils(t)

	t.Run("should return price by commodity id and region id success", func(t *testing.T) {
		repo.Price.EXPECT().FindByCommodityIDAndRegionID(ctx, ids.CommodityID, ids.RegionID).Return(mocks.Price, nil).Times(1)

		res, err := uc.GetPriceByCommodityIDAndRegionID(ctx, ids.CommodityID, ids.RegionID)
		assert.NoError(t, err)
		assert.Equal(t, ids.PriceID, res.ID)
		assert.Equal(t, ids.CommodityID, res.CommodityID)
		assert.Equal(t, ids.RegionID, res.RegionID)
		assert.Equal(t, mocks.Price.Price, res.Price)
	})

	t.Run("should return error when get price by commodity id and region id", func(t *testing.T) {
		repo.Price.EXPECT().FindByCommodityIDAndRegionID(ctx, ids.CommodityID, ids.RegionID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		_, err := uc.GetPriceByCommodityIDAndRegionID(ctx, ids.CommodityID, ids.RegionID)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

func TestGetPriceHistoryByCommodityIDAndRegionID(t *testing.T) {
	ids, mocks, _, repo, uc, ctx := PriceUsecaseUtils(t)

	t.Run("should return price history by commodity id and region id success", func(t *testing.T) {
		repo.PriceHistory.EXPECT().FindByCommodityIDAndRegionID(ctx, ids.CommodityID, ids.RegionID).Return(mocks.HistoryPrices, nil).Times(1)

		repo.Price.EXPECT().FindByCommodityIDAndRegionID(ctx, ids.CommodityID, ids.RegionID).Return(mocks.Price, nil).Times(1)

		res, err := uc.GetPriceHistoryByCommodityIDAndRegionID(ctx, ids.CommodityID, ids.RegionID)

		assert.NoError(t, err)
		assert.Equal(t, 2, len(*res))
		assert.Equal(t, ids.PriceID, (*res)[0].ID)
		assert.Equal(t, ids.CommodityID, (*res)[0].CommodityID)
		assert.Equal(t, ids.RegionID, (*res)[0].RegionID)
		assert.Equal(t, mocks.Price.Price, (*res)[0].Price)
	})

	t.Run("should return error when get price history by commodity id and region id", func(t *testing.T) {
		repo.PriceHistory.EXPECT().FindByCommodityIDAndRegionID(ctx, ids.CommodityID, ids.RegionID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		_, err := uc.GetPriceHistoryByCommodityIDAndRegionID(ctx, ids.CommodityID, ids.RegionID)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})

	t.Run("should return error when get price by commodity id and region id", func(t *testing.T) {
		repo.PriceHistory.EXPECT().FindByCommodityIDAndRegionID(ctx, ids.CommodityID, ids.RegionID).Return(mocks.HistoryPrices, nil).Times(1)

		repo.Price.EXPECT().FindByCommodityIDAndRegionID(ctx, ids.CommodityID, ids.RegionID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		_, err := uc.GetPriceHistoryByCommodityIDAndRegionID(ctx, ids.CommodityID, ids.RegionID)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}
