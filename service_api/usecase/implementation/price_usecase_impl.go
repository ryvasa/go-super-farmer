package usecase_implementation

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/pkg/database/cache"
	"github.com/ryvasa/go-super-farmer/pkg/database/transaction"
	"github.com/ryvasa/go-super-farmer/pkg/env"
	"github.com/ryvasa/go-super-farmer/pkg/logrus"
	"github.com/ryvasa/go-super-farmer/pkg/messages"
	"github.com/ryvasa/go-super-farmer/service_api/model/domain"
	"github.com/ryvasa/go-super-farmer/service_api/model/dto"
	repository_interface "github.com/ryvasa/go-super-farmer/service_api/repository/interface"
	usecase_interface "github.com/ryvasa/go-super-farmer/service_api/usecase/interface"
	"github.com/ryvasa/go-super-farmer/utils"
)

type PriceMessage struct {
	CommodityID uuid.UUID `json:"CommodityID"`
	CityID      int64     `json:"CityID"`
	StartDate   time.Time `json:"StartDate"`
	EndDate     time.Time `json:"EndDate"`
}

type PriceUsecaseImpl struct {
	priceRepo        repository_interface.PriceRepository
	priceHistoryRepo repository_interface.PriceHistoryRepository
	cityRepo         repository_interface.CityRepository
	commodityRepo    repository_interface.CommodityRepository
	rabbitMQ         messages.RabbitMQ
	txManager        transaction.TransactionManager
	cache            cache.Cache
	globFunc         utils.GlobFunc
	env              *env.Env
}

func NewPriceUsecase(priceRepo repository_interface.PriceRepository, priceHistoryRepo repository_interface.PriceHistoryRepository, cityRepo repository_interface.CityRepository, commodityRepo repository_interface.CommodityRepository, rabbitMQ messages.RabbitMQ, txManager transaction.TransactionManager, cache cache.Cache, globFunc utils.GlobFunc, env *env.Env) usecase_interface.PriceUsecase {
	return &PriceUsecaseImpl{priceRepo, priceHistoryRepo, cityRepo, commodityRepo, rabbitMQ, txManager, cache, globFunc, env}
}

func (u *PriceUsecaseImpl) CreatePrice(ctx context.Context, req *dto.PriceCreateDTO) (*domain.Price, error) {
	price := domain.Price{}
	if err := utils.ValidateStruct(price); len(err) > 0 {
		return nil, utils.NewValidationError(err)
	}

	if _, err := u.commodityRepo.FindByID(ctx, req.CommodityID); err != nil {
		return nil, utils.NewNotFoundError("commodity not found")
	}

	if _, err := u.cityRepo.FindByID(ctx, req.CityID); err != nil {
		return nil, utils.NewNotFoundError("city not found")
	}

	price.CommodityID = req.CommodityID
	price.CityID = req.CityID
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

	err = u.cache.DeleteByPattern(ctx, "price")
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	return createdPrice, nil
}
func (u *PriceUsecaseImpl) GetAllPrices(ctx context.Context, queryParams *dto.PaginationDTO) (*dto.PaginationResponseDTO, error) {
	if err := queryParams.Validate(); err != nil {
		return nil, utils.NewBadRequestError(err.Error())
	}

	var response *dto.PaginationResponseDTO

	cacheKey := fmt.Sprintf("price_list_page_%d_limit_%d",
		queryParams.Page,
		queryParams.Limit,
	)

	cached, err := u.cache.Get(ctx, cacheKey)
	if err == nil && cached != nil {
		err := json.Unmarshal(cached, &response)
		if err != nil {
			logrus.Log.Errorf("Error: %v", err)
			return nil, utils.NewInternalError("invalid data")
		}

		if data, ok := response.Data.([]interface{}); ok {
			prices := make([]*domain.Price, len(data))
			for i, item := range data {
				if priceMap, ok := item.(map[string]interface{}); ok {
					priceJSON, _ := json.Marshal(priceMap)
					var price domain.Price
					json.Unmarshal(priceJSON, &price)
					prices[i] = &price
				}
			}
			response.Data = prices
		}

		logrus.Log.Info("Cache hit")
		return response, nil
	}

	prices, err := u.priceRepo.FindAll(ctx, queryParams)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	count, err := u.priceRepo.Count(ctx, &queryParams.Filter)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	response = &dto.PaginationResponseDTO{
		TotalRows:  int64(count),
		TotalPages: int(math.Ceil(float64(count) / float64(queryParams.Limit))),
		Page:       queryParams.Page,
		Limit:      queryParams.Limit,
		Data:       prices,
	}

	pricesJSON, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}
	err = u.cache.Set(ctx, cacheKey, pricesJSON, 4*time.Minute)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	return response, nil
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

func (u *PriceUsecaseImpl) GetPricesByCityID(ctx context.Context, cityID int64) ([]*domain.Price, error) {
	prices, err := u.priceRepo.FindByCityID(ctx, cityID)
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
			CityID:      existingPrice.CityID,
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

	err = u.cache.DeleteByPattern(ctx, "price")
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

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

	err = u.cache.DeleteByPattern(ctx, "price")
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

	err = u.cache.DeleteByPattern(ctx, "price")
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	return restoredPrice, nil
}

func (u *PriceUsecaseImpl) GetPriceByCommodityIDAndCityID(ctx context.Context, commodityID uuid.UUID, cityID int64) (*domain.Price, error) {
	price, err := u.priceRepo.FindByCommodityIDAndCityID(ctx, commodityID, cityID)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	return price, nil
}
func (u *PriceUsecaseImpl) GetPriceHistoryByCommodityIDAndCityID(ctx context.Context, commodityID uuid.UUID, cityID int64) ([]*domain.PriceHistory, error) {
	cacheKey := fmt.Sprintf("price_history_%s_%d", commodityID, cityID)
	cachedPriceHistory, err := u.cache.Get(ctx, cacheKey)
	if err == nil && cachedPriceHistory != nil {
		var priceHistories []*domain.PriceHistory
		err := json.Unmarshal(cachedPriceHistory, &priceHistories)
		if err != nil {
			return nil, err
		}
		return priceHistories, nil
	}
	historyPrices, err := u.priceHistoryRepo.FindByCommodityIDAndCityID(ctx, commodityID, cityID)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	currentPrice, err := u.priceRepo.FindByCommodityIDAndCityID(ctx, commodityID, cityID)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	currentPriceHistory := &domain.PriceHistory{
		ID:          currentPrice.ID,
		CommodityID: currentPrice.CommodityID,
		CityID:      currentPrice.CityID,
		Commodity:   currentPrice.Commodity,
		City:        currentPrice.City,
		Price:       currentPrice.Price,
		CreatedAt:   currentPrice.CreatedAt,
		UpdatedAt:   currentPrice.UpdatedAt,
		DeletedAt:   currentPrice.DeletedAt,
	}
	newHistoryPrices := append(historyPrices, currentPriceHistory)

	priceJSON, err := json.Marshal(newHistoryPrices)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	err = u.cache.Set(ctx, cacheKey, priceJSON, 4*time.Minute)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	return newHistoryPrices, nil
}

func (u *PriceUsecaseImpl) DownloadPriceHistoryByCommodityIDAndCityID(ctx context.Context, params *dto.PriceParamsDTO) (*dto.DownloadResponseDTO, error) {
	_, err := u.priceRepo.FindByCommodityIDAndCityID(ctx, params.CommodityID, params.CityID)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	msg := PriceMessage{
		CommodityID: params.CommodityID,
		CityID:      params.CityID,
		StartDate:   params.StartDate,
		EndDate:     params.EndDate,
	}
	err = u.rabbitMQ.PublishJSON(ctx, "report-exchange", "price-history", msg)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	response := dto.DownloadResponseDTO{
		Message: "Price history report generation in progress. Please check back in a few moments.",
		DownloadURL: fmt.Sprintf("http://localhost%s/prices/history/commodity/%s/city/%d/download/file?start_date=%s&end_date=%s",
			u.env.Report.Port, params.CommodityID, params.CityID, params.StartDate.Format("2006-01-02"), params.EndDate.Format("2006-01-02")),
	}
	return &response, nil
}
