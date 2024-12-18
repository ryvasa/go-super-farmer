package usecase_implementation

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	"github.com/ryvasa/go-super-farmer/internal/repository/cache"
	repository_interface "github.com/ryvasa/go-super-farmer/internal/repository/interface"
	usecase_interface "github.com/ryvasa/go-super-farmer/internal/usecase/interface"
	"github.com/ryvasa/go-super-farmer/pkg/database/transaction"
	"github.com/ryvasa/go-super-farmer/pkg/logrus"
	"github.com/ryvasa/go-super-farmer/pkg/messages"
	"github.com/ryvasa/go-super-farmer/utils"
)

type PriceUsecaseImpl struct {
	priceRepo        repository_interface.PriceRepository
	priceHistoryRepo repository_interface.PriceHistoryRepository
	regionRepo       repository_interface.RegionRepository
	commodityRepo    repository_interface.CommodityRepository
	rabbitMQ         messages.RabbitMQ
	txManager        transaction.TransactionManager
	cache            cache.Cache
}

func NewPriceUsecase(priceRepo repository_interface.PriceRepository, priceHistoryRepo repository_interface.PriceHistoryRepository, regionRepo repository_interface.RegionRepository, commodityRepo repository_interface.CommodityRepository, rabbitMQ messages.RabbitMQ, txManager transaction.TransactionManager, cache cache.Cache) usecase_interface.PriceUsecase {
	return &PriceUsecaseImpl{priceRepo, priceHistoryRepo, regionRepo, commodityRepo, rabbitMQ, txManager, cache}
}

func (u *PriceUsecaseImpl) CreatePrice(ctx context.Context, req *dto.PriceCreateDTO) (*domain.Price, error) {
	price := domain.Price{}
	if err := utils.ValidateStruct(price); len(err) > 0 {
		return nil, utils.NewValidationError(err)
	}

	if _, err := u.commodityRepo.FindByID(ctx, req.CommodityID); err != nil {
		return nil, utils.NewNotFoundError("commodity not found")
	}

	if _, err := u.regionRepo.FindByID(ctx, req.RegionID); err != nil {
		return nil, utils.NewNotFoundError("region not found")
	}

	price.CommodityID = req.CommodityID
	price.RegionID = req.RegionID
	price.Price = req.Price
	price.ID = uuid.New()

	err := u.priceRepo.Create(ctx, &price)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	createdPrice, err := u.priceRepo.FindByID(ctx, price.ID)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	return createdPrice, nil
}
func (u *PriceUsecaseImpl) GetAllPrices(ctx context.Context) ([]*domain.Price, error) {
	var prices []*domain.Price
	key := fmt.Sprintf("price_%s", "all")

	cachedPrice, err := u.cache.Get(ctx, key)
	if err == nil && cachedPrice != nil {
		err := json.Unmarshal(cachedPrice, &prices)
		if err != nil {
			return nil, err
		}
		return prices, nil
	}

	prices, err = u.priceRepo.FindAll(ctx)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	pricesJSON, err := json.Marshal(prices)
	if err != nil {
		return nil, err
	}
	u.cache.Set(ctx, key, pricesJSON, 1*time.Minute)

	return prices, nil
}

func (u *PriceUsecaseImpl) GetPriceByID(ctx context.Context, id uuid.UUID) (*domain.Price, error) {
	price, err := u.priceRepo.FindByID(ctx, id)
	if err != nil {
		return nil, utils.NewNotFoundError(err.Error())
	}
	return price, nil
}

func (u *PriceUsecaseImpl) GetPricesByCommodityID(ctx context.Context, commodityID uuid.UUID) ([]*domain.Price, error) {
	prices, err := u.priceRepo.FindByCommodityID(ctx, commodityID)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	return prices, nil
}

func (u *PriceUsecaseImpl) GetPricesByRegionID(ctx context.Context, regionID uuid.UUID) ([]*domain.Price, error) {
	prices, err := u.priceRepo.FindByRegionID(ctx, regionID)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	return prices, nil
}

func (u *PriceUsecaseImpl) UpdatePrice(ctx context.Context, id uuid.UUID, req *dto.PriceUpdateDTO) (*domain.Price, error) {
	price := domain.Price{}

	if err := utils.ValidateStruct(req); len(err) > 0 {
		return nil, utils.NewValidationError(err)
	}

	err := u.txManager.WithTransaction(ctx, func(txCtx context.Context) error {

		logrus.Log.Info("starting price update transaction")

		existingPrice, err := u.priceRepo.FindByID(txCtx, id)
		if err != nil {
			logrus.Log.Error(err, "failed to find price")
			return utils.NewNotFoundError(err.Error())
		}

		historyPrice := domain.PriceHistory{
			ID:          uuid.New(),
			CommodityID: existingPrice.CommodityID,
			RegionID:    existingPrice.RegionID,
			Price:       existingPrice.Price,
			CreatedAt:   existingPrice.CreatedAt,
			UpdatedAt:   existingPrice.UpdatedAt,
		}

		err = u.priceHistoryRepo.Create(txCtx, &historyPrice)
		if err != nil {
			logrus.Log.Error(err, "failed to create price history")
			return err
		}
		logrus.Log.Info("price history created")

		price.Price = req.Price
		price.ID = id

		err = u.priceRepo.Update(txCtx, id, &price)
		if err != nil {
			logrus.Log.Error(err, "failed to update price")
			return err
		}
		updatedPrice, err := u.priceRepo.FindByID(txCtx, id)
		if err != nil {
			logrus.Log.Error(err, "failed to find price")
			return err
		}

		price = *updatedPrice

		return nil
	})

	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	logrus.Log.Info("price update transaction completed")

	return &price, nil
}

func (u *PriceUsecaseImpl) DeletePrice(ctx context.Context, id uuid.UUID) error {
	_, err := u.priceRepo.FindByID(ctx, id)
	if err != nil {
		return utils.NewNotFoundError(err.Error())
	}
	err = u.priceRepo.Delete(ctx, id)
	if err != nil {
		return utils.NewInternalError(err.Error())
	}
	return nil
}

func (u *PriceUsecaseImpl) RestorePrice(ctx context.Context, id uuid.UUID) (*domain.Price, error) {
	_, err := u.priceRepo.FindDeletedByID(ctx, id)
	if err != nil {
		return nil, utils.NewNotFoundError(err.Error())
	}
	err = u.priceRepo.Restore(ctx, id)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	restoredPrice, err := u.priceRepo.FindByID(ctx, id)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	return restoredPrice, nil
}

func (u *PriceUsecaseImpl) GetPriceByCommodityIDAndRegionID(ctx context.Context, commodityID, regionID uuid.UUID) (*domain.Price, error) {
	price, err := u.priceRepo.FindByCommodityIDAndRegionID(ctx, commodityID, regionID)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	return price, nil
}
func (u *PriceUsecaseImpl) GetPriceHistoryByCommodityIDAndRegionID(ctx context.Context, commodityID, regionID uuid.UUID) ([]*domain.PriceHistory, error) {
	cacheKey := fmt.Sprintf("price_history_%s_%s", commodityID, regionID)
	cachedPriceHistory, err := u.cache.Get(ctx, cacheKey)
	if err == nil && cachedPriceHistory != nil {
		var priceHistories []*domain.PriceHistory
		err := json.Unmarshal(cachedPriceHistory, &priceHistories)
		if err != nil {
			return nil, err
		}
		return priceHistories, nil
	}
	historyPrices, err := u.priceHistoryRepo.FindByCommodityIDAndRegionID(ctx, commodityID, regionID)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	currentPrice, err := u.priceRepo.FindByCommodityIDAndRegionID(ctx, commodityID, regionID)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	currentPriceHistory := &domain.PriceHistory{
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
	newHistoryPrices := append(historyPrices, currentPriceHistory)

	userJSON, _ := json.Marshal(newHistoryPrices)
	u.cache.Set(ctx, cacheKey, userJSON, 1*time.Minute)
	return newHistoryPrices, nil
}

func (u *PriceUsecaseImpl) DownloadPriceHistoryByCommodityIDAndRegionID(ctx context.Context, params *dto.PriceParamsDTO) error {
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
	err := u.rabbitMQ.PublishJSON(ctx, "report-exchange", "price-history", msg)
	if err != nil {
		return utils.NewInternalError(err.Error())
	}
	return nil
}
