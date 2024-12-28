package usecase_implementation

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	repository_interface "github.com/ryvasa/go-super-farmer/internal/repository/interface"
	usecase_interface "github.com/ryvasa/go-super-farmer/internal/usecase/interface"
	"github.com/ryvasa/go-super-farmer/utils"
)

type ForecastsUsecaseImpl struct {
	landCommodityRepo repository_interface.LandCommodityRepository
	cityRepo          repository_interface.CityRepository
	priceRepo         repository_interface.PriceRepository
	priceHistoryRepo  repository_interface.PriceHistoryRepository
	demandRepo        repository_interface.DemandRepository
	demandHistoryRepo repository_interface.DemandHistoryRepository
	supplyRepo        repository_interface.SupplyRepository
	supplyHistoryRepo repository_interface.SupplyHistoryRepository
	saleRepo          repository_interface.SaleRepository
	harvestRepo       repository_interface.HarvestRepository
	commodityRepo     repository_interface.CommodityRepository
}

func NewForecastsUsecase(
	landCommodityRepo repository_interface.LandCommodityRepository,
	cityRepo repository_interface.CityRepository,
	priceRepo repository_interface.PriceRepository,
	priceHistoryRepo repository_interface.PriceHistoryRepository,
	demandRepo repository_interface.DemandRepository,
	demandHistoryRepo repository_interface.DemandHistoryRepository,
	supplyRepo repository_interface.SupplyRepository,
	supplyHistoryRepo repository_interface.SupplyHistoryRepository,
	saleRepo repository_interface.SaleRepository,
	harvestRepo repository_interface.HarvestRepository,
	commodityRepo repository_interface.CommodityRepository,
) usecase_interface.ForecastsUsecase {
	return &ForecastsUsecaseImpl{
		landCommodityRepo: landCommodityRepo,
		cityRepo:          cityRepo,
		priceRepo:         priceRepo,
		priceHistoryRepo:  priceHistoryRepo,
		demandRepo:        demandRepo,
		demandHistoryRepo: demandHistoryRepo,
		supplyRepo:        supplyRepo,
		supplyHistoryRepo: supplyHistoryRepo,
		saleRepo:          saleRepo,
		harvestRepo:       harvestRepo,
		commodityRepo:     commodityRepo,
	}
}

func (f *ForecastsUsecaseImpl) GetForecastsByCommodityIDAndCityID(ctx context.Context, commodityID uuid.UUID, cityID int64) (*dto.ForecastsResponseDTO, error) {
	_, err := f.commodityRepo.FindByID(ctx, commodityID)
	if err != nil {
		return nil, utils.NewNotFoundError("commodity not found")
	}

	_, err = f.cityRepo.FindByID(ctx, cityID)
	if err != nil {
		return nil, utils.NewNotFoundError("city not found")
	}

	return nil, nil
}
