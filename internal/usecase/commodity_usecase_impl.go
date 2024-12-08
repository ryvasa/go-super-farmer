package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	"github.com/ryvasa/go-super-farmer/internal/repository"
	"github.com/ryvasa/go-super-farmer/utils"
)

type CommodityUsecaseImpl struct {
	commodityRepository repository.CommodityRepository
}

func NewCommodityUsecase(commodityRepository repository.CommodityRepository) CommodityUsecase {
	return &CommodityUsecaseImpl{commodityRepository}
}

func (c *CommodityUsecaseImpl) CreateCommodity(ctx context.Context, req *dto.CommodityCreateDTO) (*domain.Commodity, error) {
	commodity := domain.Commodity{}
	if err := utils.ValidateStruct(req); len(err) > 0 {
		return nil, utils.NewValidationError(err)
	}

	commodity.Name = req.Name
	commodity.Description = req.Description
	commodity.ID = uuid.New()

	err := c.commodityRepository.Create(ctx, &commodity)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	createdCommodity, err := c.commodityRepository.FindByID(ctx, commodity.ID)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	return createdCommodity, nil
}

func (c *CommodityUsecaseImpl) GetAllCommodities(ctx context.Context) (*[]domain.Commodity, error) {
	commodities, err := c.commodityRepository.FindAll(ctx)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	return commodities, nil
}

func (c *CommodityUsecaseImpl) GetCommodityById(ctx context.Context, id uuid.UUID) (*domain.Commodity, error) {
	commodity, err := c.commodityRepository.FindByID(ctx, id)
	if err != nil {
		return nil, utils.NewNotFoundError("commodity not found")
	}
	return commodity, nil
}

func (c *CommodityUsecaseImpl) UpdateCommodity(ctx context.Context, id uuid.UUID, req *dto.CommodityUpdateDTO) (*domain.Commodity, error) {
	if err := utils.ValidateStruct(req); len(err) > 0 {
		return nil, utils.NewValidationError(err)
	}

	commodity, err := c.commodityRepository.FindByID(ctx, id)
	if err != nil {
		return nil, utils.NewNotFoundError("commodity not found")
	}

	commodity.Name = req.Name
	commodity.Description = req.Description

	err = c.commodityRepository.Update(ctx, id, commodity)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	updatedCommodity, err := c.commodityRepository.FindByID(ctx, id)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	return updatedCommodity, nil
}

func (c *CommodityUsecaseImpl) DeleteCommodity(ctx context.Context, id uuid.UUID) error {
	_, err := c.commodityRepository.FindByID(ctx, id)
	if err != nil {
		return utils.NewNotFoundError("commodity not found")
	}

	err = c.commodityRepository.Delete(ctx, id)
	if err != nil {
		return utils.NewInternalError(err.Error())
	}

	return nil
}

func (c *CommodityUsecaseImpl) RestoreCommodity(ctx context.Context, id uuid.UUID) (*domain.Commodity, error) {
	_, err := c.commodityRepository.FindDeletedByID(ctx, id)
	if err != nil {
		return nil, utils.NewNotFoundError("deleted commodity not found")
	}

	err = c.commodityRepository.Restore(ctx, id)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	restoredCommodity, err := c.commodityRepository.FindByID(ctx, id)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	return restoredCommodity, nil
}
