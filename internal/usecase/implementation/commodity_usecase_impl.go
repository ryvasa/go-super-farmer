package usecase_implementation

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	repository_interface "github.com/ryvasa/go-super-farmer/internal/repository/interface"
	usecase_interface "github.com/ryvasa/go-super-farmer/internal/usecase/interface"
	"github.com/ryvasa/go-super-farmer/pkg/database/cache"
	"github.com/ryvasa/go-super-farmer/pkg/logrus"
	"github.com/ryvasa/go-super-farmer/utils"
)

type CommodityUsecaseImpl struct {
	commodityRepository repository_interface.CommodityRepository
	cache               cache.Cache
}

func NewCommodityUsecase(commodityRepository repository_interface.CommodityRepository, cache cache.Cache) usecase_interface.CommodityUsecase {
	return &CommodityUsecaseImpl{commodityRepository, cache}
}

func (uc *CommodityUsecaseImpl) CreateCommodity(ctx context.Context, req *dto.CommodityCreateDTO) (*domain.Commodity, error) {
	commodity := domain.Commodity{}
	if err := utils.ValidateStruct(req); len(err) > 0 {
		return nil, utils.NewValidationError(err)
	}

	commodity.Name = req.Name
	commodity.Description = req.Description
	commodity.Code = req.Code
	commodity.ID = uuid.New()

	err := uc.commodityRepository.Create(ctx, &commodity)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	createdCommodity, err := uc.commodityRepository.FindByID(ctx, commodity.ID)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	err = uc.cache.DeleteByPattern(ctx, "commodity")
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	return createdCommodity, nil
}
func (uc *CommodityUsecaseImpl) GetAllCommodities(ctx context.Context, queryParams *dto.PaginationDTO) (*dto.PaginationResponseDTO, error) {
	// Validasi pagination params
	if err := queryParams.Validate(); err != nil {
		return nil, utils.NewBadRequestError(err.Error())
	}

	// Cek cache
	cacheKey := fmt.Sprintf("commodity_list_page_%d_limit_%d_%s",
		queryParams.Page,
		queryParams.Limit,
		queryParams.Filter.CommodityName,
	)

	var response *dto.PaginationResponseDTO
	cached, err := uc.cache.Get(ctx, cacheKey)
	if err == nil && cached != nil {
		err := json.Unmarshal(cached, &response)
		if err != nil {
			logrus.Log.Error("Error get from cache")
			return nil, err
		}
		logrus.Log.Info("Cache hit")
		return response, nil
	}

	// Get data dari repository
	commodities, err := uc.commodityRepository.FindAll(ctx, queryParams)
	if err != nil {
		logrus.Log.Error("Error get from repository")
		return nil, utils.NewInternalError(err.Error())
	}

	count, err := uc.commodityRepository.Count(ctx, &queryParams.Filter)
	if err != nil {
		logrus.Log.Error("Error get count from repository")
		return nil, utils.NewInternalError(err.Error())
	}

	// Create response
	response = &dto.PaginationResponseDTO{
		TotalRows:  int64(count),
		TotalPages: int(math.Ceil(float64(count) / float64(queryParams.Limit))),
		Page:       queryParams.Page,
		Limit:      queryParams.Limit,
		Data:       commodities,
	}

	// Set cache
	responseJSON, err := json.Marshal(response)
	if err != nil {
		logrus.Log.Error("Error marshal response")
		return nil, utils.NewInternalError(err.Error())
	}
	err = uc.cache.Set(ctx, cacheKey, responseJSON, 4*time.Minute)
	if err != nil {
		logrus.Log.Error("Error set cache")
		return nil, utils.NewInternalError(err.Error())
	}

	return response, nil
}

func (uc *CommodityUsecaseImpl) GetCommodityById(ctx context.Context, id uuid.UUID) (*domain.Commodity, error) {
	commodity, err := uc.commodityRepository.FindByID(ctx, id)
	if err != nil {
		return nil, utils.NewNotFoundError("commodity not found")
	}
	return commodity, nil
}

func (uc *CommodityUsecaseImpl) UpdateCommodity(ctx context.Context, id uuid.UUID, req *dto.CommodityUpdateDTO) (*domain.Commodity, error) {
	if err := utils.ValidateStruct(req); len(err) > 0 {
		return nil, utils.NewValidationError(err)
	}

	commodity, err := uc.commodityRepository.FindByID(ctx, id)
	if err != nil {
		return nil, utils.NewNotFoundError("commodity not found")
	}
	if req.Name != nil {
		commodity.Name = *req.Name
	}
	if req.Code != nil {
		commodity.Code = *req.Code
	}
	if req.Description != nil {
		commodity.Description = *req.Description
	}

	err = uc.cache.DeleteByPattern(ctx, "commodity")
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	err = uc.commodityRepository.Update(ctx, id, commodity)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	updatedCommodity, err := uc.commodityRepository.FindByID(ctx, id)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	return updatedCommodity, nil
}

func (uc *CommodityUsecaseImpl) DeleteCommodity(ctx context.Context, id uuid.UUID) error {
	_, err := uc.commodityRepository.FindByID(ctx, id)
	if err != nil {
		return utils.NewNotFoundError("commodity not found")
	}

	err = uc.commodityRepository.Delete(ctx, id)
	if err != nil {
		return utils.NewInternalError(err.Error())
	}

	err = uc.cache.DeleteByPattern(ctx, "commodity")
	if err != nil {
		return utils.NewInternalError(err.Error())
	}

	return nil
}

func (uc *CommodityUsecaseImpl) RestoreCommodity(ctx context.Context, id uuid.UUID) (*domain.Commodity, error) {
	_, err := uc.commodityRepository.FindDeletedByID(ctx, id)
	if err != nil {
		return nil, utils.NewNotFoundError("deleted commodity not found")
	}

	err = uc.commodityRepository.Restore(ctx, id)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	restoredCommodity, err := uc.commodityRepository.FindByID(ctx, id)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	return restoredCommodity, nil
}
