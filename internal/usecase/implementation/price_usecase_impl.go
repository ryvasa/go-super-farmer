package usecase_implementation

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	repository_interface "github.com/ryvasa/go-super-farmer/internal/repository/interface"
	usecase_interface "github.com/ryvasa/go-super-farmer/internal/usecase/interface"
	"github.com/ryvasa/go-super-farmer/pkg/messages"
	"github.com/ryvasa/go-super-farmer/utils"
)

type PriceUsecaseImpl struct {
	priceRepo        repository_interface.PriceRepository
	priceHistoryRepo repository_interface.PriceHistoryRepository
	regionRepo       repository_interface.RegionRepository
	commodityRepo    repository_interface.CommodityRepository
	rabbitMQ         messages.RabbitMQ
}

func NewPriceUsecase(priceRepo repository_interface.PriceRepository, priceHistoryRepo repository_interface.PriceHistoryRepository, regionRepo repository_interface.RegionRepository, commodityRepo repository_interface.CommodityRepository, rabbitMQ messages.RabbitMQ) usecase_interface.PriceUsecase {
	return &PriceUsecaseImpl{priceRepo, priceHistoryRepo, regionRepo, commodityRepo, rabbitMQ}
}

func (uc *PriceUsecaseImpl) CreatePrice(ctx context.Context, req *dto.PriceCreateDTO) (*domain.Price, error) {
	price := domain.Price{}
	if err := utils.ValidateStruct(price); len(err) > 0 {
		return nil, utils.NewValidationError(err)
	}

	if _, err := uc.commodityRepo.FindByID(ctx, req.CommodityID); err != nil {
		return nil, utils.NewNotFoundError("commodity not found")
	}

	if _, err := uc.regionRepo.FindByID(ctx, req.RegionID); err != nil {
		return nil, utils.NewNotFoundError("region not found")
	}

	price.CommodityID = req.CommodityID
	price.RegionID = req.RegionID
	price.Price = req.Price
	price.ID = uuid.New()

	err := uc.priceRepo.Create(ctx, &price)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	createdPrice, err := uc.priceRepo.FindByID(ctx, price.ID)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	return createdPrice, nil
}

func (uc *PriceUsecaseImpl) GetAllPrices(ctx context.Context) (*[]domain.Price, error) {
	prices, err := uc.priceRepo.FindAll(ctx)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	return prices, nil
}

func (uc *PriceUsecaseImpl) GetPriceByID(ctx context.Context, id uuid.UUID) (*domain.Price, error) {
	price, err := uc.priceRepo.FindByID(ctx, id)
	if err != nil {
		return nil, utils.NewNotFoundError(err.Error())
	}
	return price, nil
}

func (uc *PriceUsecaseImpl) GetPricesByCommodityID(ctx context.Context, commodityID uuid.UUID) (*[]domain.Price, error) {
	prices, err := uc.priceRepo.FindByCommodityID(ctx, commodityID)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	return prices, nil
}

func (uc *PriceUsecaseImpl) GetPricesByRegionID(ctx context.Context, regionID uuid.UUID) (*[]domain.Price, error) {
	prices, err := uc.priceRepo.FindByRegionID(ctx, regionID)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	return prices, nil
}

func (uc *PriceUsecaseImpl) UpdatePrice(ctx context.Context, id uuid.UUID, req *dto.PriceUpdateDTO) (*domain.Price, error) {
	price := domain.Price{}
	if err := utils.ValidateStruct(req); len(err) > 0 {
		return nil, utils.NewValidationError(err)
	}

	existingPrice, err := uc.priceRepo.FindByID(ctx, id)
	if err != nil {
		return nil, utils.NewNotFoundError(err.Error())
	}

	historyPrice := domain.PriceHistory{
		ID:          uuid.New(),
		CommodityID: existingPrice.CommodityID,
		RegionID:    existingPrice.RegionID,
		Price:       existingPrice.Price,
		CreatedAt:   existingPrice.CreatedAt,
		UpdatedAt:   existingPrice.UpdatedAt,
	}
	err = uc.priceHistoryRepo.Create(ctx, &historyPrice)
	if err != nil {
		return nil, utils.NewInternalError("failed to create price history")
	}

	price.Price = req.Price
	price.ID = id

	err = uc.priceRepo.Update(ctx, id, &price)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	updatedPrice, err := uc.priceRepo.FindByID(ctx, id)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	return updatedPrice, nil
}

func (uc *PriceUsecaseImpl) DeletePrice(ctx context.Context, id uuid.UUID) error {
	_, err := uc.priceRepo.FindByID(ctx, id)
	if err != nil {
		return utils.NewNotFoundError(err.Error())
	}
	err = uc.priceRepo.Delete(ctx, id)
	if err != nil {
		return utils.NewInternalError(err.Error())
	}
	return nil
}

func (uc *PriceUsecaseImpl) RestorePrice(ctx context.Context, id uuid.UUID) (*domain.Price, error) {
	_, err := uc.priceRepo.FindDeletedByID(ctx, id)
	if err != nil {
		return nil, utils.NewNotFoundError(err.Error())
	}
	err = uc.priceRepo.Restore(ctx, id)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	restoredPrice, err := uc.priceRepo.FindByID(ctx, id)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	return restoredPrice, nil
}

func (uc *PriceUsecaseImpl) GetPriceByCommodityIDAndRegionID(ctx context.Context, commodityID, regionID uuid.UUID) (*domain.Price, error) {
	price, err := uc.priceRepo.FindByCommodityIDAndRegionID(ctx, commodityID, regionID)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	return price, nil
}
func (uc *PriceUsecaseImpl) GetPriceHistoryByCommodityIDAndRegionID(ctx context.Context, commodityID, regionID uuid.UUID) (*[]domain.PriceHistory, error) {
	historyPrices, err := uc.priceHistoryRepo.FindByCommodityIDAndRegionID(ctx, commodityID, regionID)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	currentPrice, err := uc.priceRepo.FindByCommodityIDAndRegionID(ctx, commodityID, regionID)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	currentPriceHistory := domain.PriceHistory{
		ID:          currentPrice.ID,
		CommodityID: currentPrice.CommodityID,
		RegionID:    currentPrice.RegionID,
		Commodity:   currentPrice.Commodity,
		Region:      currentPrice.Region,
		Price:       currentPrice.Price,
		CreatedAt:   currentPrice.CreatedAt,
		UpdatedAt:   currentPrice.UpdatedAt,
		DeletedAt:   currentPrice.DeletedAt,
	}
	newHistoryPrices := append(*historyPrices, currentPriceHistory)
	return &newHistoryPrices, nil
}

func (uc *PriceUsecaseImpl) DownloadPriceHistoryByCommodityIDAndRegionID(ctx context.Context, params *dto.PriceParamsDTO) error {
	type Message struct {
		CommodityID uuid.UUID `json:"CommodityID"`
		RegionID    uuid.UUID `json:"RegionID"`
		StartDate   time.Time `json:"StartDate"`
		EndDate     time.Time `json:"EndDate"`
	}
	msg := Message{
		CommodityID: params.CommodityID,
		RegionID:    params.RegionID,
		StartDate:   params.StartDate,
		EndDate:     params.EndDate,
	}
	err := uc.rabbitMQ.PublishJSON(ctx, "report-exchange", "price-history", msg)
	if err != nil {
		return utils.NewInternalError(err.Error())
	}
	return nil
}
