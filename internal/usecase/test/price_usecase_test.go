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

type IDs struct {
	PriceID        uuid.UUID
	PriceHistoryID uuid.UUID
	CommodityID    uuid.UUID
	RegionID       uuid.UUID
}

type Mocks struct {
	Prices        *[]domain.Price
	Price         *domain.Price
	UpdatedPrice  *domain.Price
	HistoryPrices *[]domain.PriceHistory
	HistoryPrice  *domain.PriceHistory
	Commodity     *domain.Commodity
	Region        *domain.Region
}

func PriceUsecaseUtils() (*IDs, *Mocks) {
	regionID := uuid.New()
	commodityID := uuid.New()
	priceID := uuid.New()
	priceHistoryID := uuid.New()

	ids := &IDs{
		PriceID:        priceID,
		PriceHistoryID: priceHistoryID,
		CommodityID:    commodityID,
		RegionID:       regionID,
	}

	mocks := &Mocks{
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

	return ids, mocks
}

func TestCreatePrice(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	regionRepo := mock.NewMockRegionRepository(ctrl)
	commodityRepo := mock.NewMockCommodityRepository(ctrl)
	priceRepo := mock.NewMockPriceRepository(ctrl)
	priceHostoryRepo := mock.NewMockPriceHistoryRepository(ctrl)

	uc := usecase.NewPriceUsecase(priceRepo, priceHostoryRepo, regionRepo, commodityRepo)
	ctx := context.Background()

	IDs, Mocks := PriceUsecaseUtils()

	req := &dto.PriceCreateDTO{
		CommodityID: IDs.CommodityID,
		RegionID:    IDs.RegionID,
		Price:       Mocks.Price.Price,
	}

	t.Run("should create price successfully", func(t *testing.T) {

		commodityRepo.EXPECT().FindByID(ctx, IDs.CommodityID).Return(Mocks.Commodity, nil).Times(1)

		regionRepo.EXPECT().FindByID(ctx, IDs.RegionID).Return(Mocks.Region, nil).Times(1)

		priceRepo.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, p *domain.Price) error {
			p.ID = IDs.PriceID
			return nil
		}).Times(1)

		priceRepo.EXPECT().FindByID(ctx, IDs.PriceID).Return(Mocks.Price, nil).Times(1)

		resp, err := uc.CreatePrice(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, req.CommodityID, resp.CommodityID)
		assert.Equal(t, req.RegionID, resp.RegionID)
		assert.Equal(t, Mocks.Price.ID, resp.ID)
	})

	t.Run("should return error when commodity not found", func(t *testing.T) {
		commodityRepo.EXPECT().FindByID(ctx, IDs.CommodityID).Return(nil, utils.NewNotFoundError("commodity not found")).Times(1)

		resp, err := uc.CreatePrice(ctx, req)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "commodity not found")
	})

	t.Run("should return error when region not found", func(t *testing.T) {
		commodityRepo.EXPECT().FindByID(ctx, IDs.CommodityID).Return(Mocks.Commodity, nil).Times(1)

		regionRepo.EXPECT().FindByID(ctx, IDs.RegionID).Return(nil, utils.NewNotFoundError("region not found")).Times(1)

		resp, err := uc.CreatePrice(ctx, req)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "region not found")
	})

	t.Run("should return error when create price", func(t *testing.T) {
		commodityRepo.EXPECT().FindByID(ctx, IDs.CommodityID).Return(Mocks.Commodity, nil).Times(1)

		regionRepo.EXPECT().FindByID(ctx, IDs.RegionID).Return(Mocks.Region, nil).Times(1)

		priceRepo.EXPECT().Create(ctx, gomock.Any()).Return(utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.CreatePrice(ctx, req)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})

	t.Run("should return error when get created price", func(t *testing.T) {
		commodityRepo.EXPECT().FindByID(ctx, IDs.CommodityID).Return(Mocks.Commodity, nil).Times(1)

		regionRepo.EXPECT().FindByID(ctx, IDs.RegionID).Return(Mocks.Region, nil).Times(1)

		priceRepo.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, p *domain.Price) error {
			p.ID = IDs.PriceID
			return nil
		}).Times(1)

		priceRepo.EXPECT().FindByID(ctx, IDs.PriceID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.CreatePrice(ctx, req)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

func TestGetAllPrices(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	regionRepo := mock.NewMockRegionRepository(ctrl)
	commodityRepo := mock.NewMockCommodityRepository(ctrl)
	priceRepo := mock.NewMockPriceRepository(ctrl)
	priceHostoryRepo := mock.NewMockPriceHistoryRepository(ctrl)

	uc := usecase.NewPriceUsecase(priceRepo, priceHostoryRepo, regionRepo, commodityRepo)
	ctx := context.Background()

	_, Mocks := PriceUsecaseUtils()

	t.Run("should return all prices", func(t *testing.T) {
		priceRepo.EXPECT().FindAll(ctx).Return(Mocks.Prices, nil).Times(1)

		res, err := uc.GetAllPrices(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(*res))
		assert.Equal(t, (*Mocks.Prices)[0].ID, (*res)[0].ID)
		assert.Equal(t, (*Mocks.Prices)[0].CommodityID, (*res)[0].CommodityID)
		assert.Equal(t, (*Mocks.Prices)[0].RegionID, (*res)[0].RegionID)
		assert.Equal(t, (*Mocks.Prices)[0].Price, (*res)[0].Price)
	})

	t.Run("should return error when get all prices", func(t *testing.T) {
		priceRepo.EXPECT().FindAll(ctx).Return(nil, utils.NewInternalError("internal error")).Times(1)

		_, err := uc.GetAllPrices(ctx)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

func TestGetPriceByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	regionRepo := mock.NewMockRegionRepository(ctrl)
	commodityRepo := mock.NewMockCommodityRepository(ctrl)
	priceRepo := mock.NewMockPriceRepository(ctrl)
	priceHostoryRepo := mock.NewMockPriceHistoryRepository(ctrl)

	uc := usecase.NewPriceUsecase(priceRepo, priceHostoryRepo, regionRepo, commodityRepo)
	ctx := context.Background()

	IDs, Mocks := PriceUsecaseUtils()

	t.Run("should return price by id success", func(t *testing.T) {
		priceRepo.EXPECT().FindByID(ctx, IDs.PriceID).Return(Mocks.Price, nil).Times(1)

		res, err := uc.GetPriceByID(ctx, IDs.PriceID)
		assert.NoError(t, err)
		assert.Equal(t, IDs.PriceID, res.ID)
		assert.Equal(t, IDs.CommodityID, res.CommodityID)
		assert.Equal(t, IDs.RegionID, res.RegionID)
		assert.Equal(t, Mocks.Price.Price, res.Price)
	})

	t.Run("should return error when get price by id", func(t *testing.T) {
		priceRepo.EXPECT().FindByID(ctx, IDs.PriceID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		_, err := uc.GetPriceByID(ctx, IDs.PriceID)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

func TestGetPricesByCommodityID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	regionRepo := mock.NewMockRegionRepository(ctrl)
	commodityRepo := mock.NewMockCommodityRepository(ctrl)
	priceRepo := mock.NewMockPriceRepository(ctrl)
	priceHostoryRepo := mock.NewMockPriceHistoryRepository(ctrl)

	uc := usecase.NewPriceUsecase(priceRepo, priceHostoryRepo, regionRepo, commodityRepo)
	ctx := context.Background()

	IDs, Mocks := PriceUsecaseUtils()

	t.Run("should return price by commodity id success", func(t *testing.T) {
		priceRepo.EXPECT().FindByCommodityID(ctx, IDs.CommodityID).Return(Mocks.Prices, nil).Times(1)

		res, err := uc.GetPricesByCommodityID(ctx, IDs.CommodityID)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(*res))
		assert.Equal(t, (*Mocks.Prices)[0].ID, (*res)[0].ID)
		assert.Equal(t, (*Mocks.Prices)[0].CommodityID, (*res)[0].CommodityID)
		assert.Equal(t, (*Mocks.Prices)[0].RegionID, (*res)[0].RegionID)
		assert.Equal(t, (*Mocks.Prices)[0].Price, (*res)[0].Price)
	})

	t.Run("should return error when get price by commodity id", func(t *testing.T) {
		priceRepo.EXPECT().FindByCommodityID(ctx, IDs.CommodityID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		_, err := uc.GetPricesByCommodityID(ctx, IDs.CommodityID)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

func TestGetPricesByRegionID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	regionRepo := mock.NewMockRegionRepository(ctrl)
	commodityRepo := mock.NewMockCommodityRepository(ctrl)
	priceRepo := mock.NewMockPriceRepository(ctrl)
	priceHostoryRepo := mock.NewMockPriceHistoryRepository(ctrl)

	uc := usecase.NewPriceUsecase(priceRepo, priceHostoryRepo, regionRepo, commodityRepo)
	ctx := context.Background()

	IDs, Mocks := PriceUsecaseUtils()

	t.Run("should return price by region id success", func(t *testing.T) {
		priceRepo.EXPECT().FindByRegionID(ctx, IDs.RegionID).Return(Mocks.Prices, nil).Times(1)

		res, err := uc.GetPricesByRegionID(ctx, IDs.RegionID)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(*res))
		assert.Equal(t, (*Mocks.Prices)[0].ID, (*res)[0].ID)
		assert.Equal(t, (*Mocks.Prices)[0].CommodityID, (*res)[0].CommodityID)
		assert.Equal(t, (*Mocks.Prices)[0].RegionID, (*res)[0].RegionID)
		assert.Equal(t, (*Mocks.Prices)[0].Price, (*res)[0].Price)
	})

	t.Run("should return error when get price by region id", func(t *testing.T) {
		priceRepo.EXPECT().FindByRegionID(ctx, IDs.RegionID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		_, err := uc.GetPricesByRegionID(ctx, IDs.RegionID)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

// TODO : update besok
func TestUpdatePrice(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	regionRepo := mock.NewMockRegionRepository(ctrl)
	commodityRepo := mock.NewMockCommodityRepository(ctrl)
	priceRepo := mock.NewMockPriceRepository(ctrl)
	priceHostoryRepo := mock.NewMockPriceHistoryRepository(ctrl)

	uc := usecase.NewPriceUsecase(priceRepo, priceHostoryRepo, regionRepo, commodityRepo)
	ctx := context.Background()

	IDs, Mocks := PriceUsecaseUtils()

	req := &dto.PriceUpdateDTO{
		Price: Mocks.UpdatedPrice.Price,
	}

	t.Run("should update price successfully", func(t *testing.T) {
		priceRepo.EXPECT().FindByID(ctx, IDs.PriceID).Return(Mocks.Price, nil).Times(1)

		priceHostoryRepo.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, ph *domain.PriceHistory) error {
			ph.ID = IDs.PriceHistoryID
			return nil
		}).Times(1)

		priceRepo.EXPECT().Update(ctx, IDs.PriceID, gomock.Any()).DoAndReturn(func(ctx context.Context, p *domain.Price) error {
			p.ID = IDs.PriceID
			return nil
		}).Times(1)

		priceRepo.EXPECT().FindByID(ctx, IDs.PriceID).Return(Mocks.UpdatedPrice, nil).Times(1)

		resp, err := uc.UpdatePrice(ctx, IDs.PriceID, req)

		assert.NoError(t, err)
		assert.Equal(t, Mocks.Price.ID, resp.ID)
		assert.Equal(t, resp.Price, Mocks.UpdatedPrice.Price)
	})

	t.Run("should return error when commodity not found", func(t *testing.T) {
		commodityRepo.EXPECT().FindByID(ctx, IDs.CommodityID).Return(nil, utils.NewNotFoundError("commodity not found")).Times(1)

		resp, err := uc.UpdatePrice(ctx, IDs.PriceID, req)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "commodity not found")
	})

	t.Run("should return error when region not found", func(t *testing.T) {
		commodityRepo.EXPECT().FindByID(ctx, IDs.CommodityID).Return(Mocks.Commodity, nil).Times(1)

		regionRepo.EXPECT().FindByID(ctx, IDs.RegionID).Return(nil, utils.NewNotFoundError("region not found")).Times(1)

		resp, err := uc.UpdatePrice(ctx, IDs.PriceID, req)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "region not found")
	})

	t.Run("should return error when update price", func(t *testing.T) {
		commodityRepo.EXPECT().FindByID(ctx, IDs.CommodityID).Return(Mocks.Commodity, nil).Times(1)

		regionRepo.EXPECT().FindByID(ctx, IDs.RegionID).Return(Mocks.Region, nil).Times(1)

		priceRepo.EXPECT().Update(ctx, IDs.PriceID, gomock.Any()).Return(utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.UpdatePrice(ctx, IDs.PriceID, req)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})

	t.Run("should return error when get updated price", func(t *testing.T) {
		commodityRepo.EXPECT().FindByID(ctx, IDs.CommodityID).Return(Mocks.Commodity, nil).Times(1)

		regionRepo.EXPECT().FindByID(ctx, IDs.RegionID).Return(Mocks.Region, nil).Times(1)

		priceRepo.EXPECT().Update(ctx, IDs.PriceID, gomock.Any()).DoAndReturn(func(ctx context.Context, p *domain.Price) error {
			p.ID = IDs.PriceID
			return nil
		}).Times(1)

		priceRepo.EXPECT().FindByID(ctx, IDs.PriceID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.UpdatePrice(ctx, IDs.PriceID, req)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}
