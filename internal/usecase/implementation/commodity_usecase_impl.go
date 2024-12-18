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
	"github.com/ryvasa/go-super-farmer/pkg/database/pagination"
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

	return createdCommodity, nil
}

func (uc *CommodityUsecaseImpl) GetAllCommodities(ctx context.Context, queryParams *dto.PaginationDTO) ([]domain.Commodity, error) {
	var commodities []domain.Commodity
	key := fmt.Sprintf("commodity_%s_start_%s_end_%s", queryParams.CommodityName, queryParams.StartDate, queryParams.EndDate)

	cached, err := uc.cache.Get(ctx, key)
	if err == nil && cached != nil {
		err := json.Unmarshal(cached, &commodities)
		if err != nil {
			return nil, err
		}
		return commodities, nil
	}
	params := &pagination.PaginationParams{
		Limit:         queryParams.Limit,
		Page:          queryParams.Page,
		Sort:          queryParams.Sort,
		CommodityName: queryParams.CommodityName,
		StartDate:     queryParams.StartDate,
		EndDate:       queryParams.EndDate,
	}
	commodities, err = uc.commodityRepository.FindAll(ctx, params)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	commoditiesJSON, err := json.Marshal(commodities)
	if err != nil {
		return nil, err
	}

	uc.cache.Set(ctx, key, commoditiesJSON, 4*time.Minute)

	return commodities, nil
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
