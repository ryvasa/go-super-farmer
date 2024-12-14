package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	"github.com/ryvasa/go-super-farmer/internal/repository"
	"github.com/ryvasa/go-super-farmer/utils"
)

type SupplyUsecaseImpl struct {
	supplyRepo        repository.SupplyRepository
	supplyHistoryRepo repository.SupplyHistoryRepository
	commodityRepo     repository.CommodityRepository
	regionRepo        repository.RegionRepository
}

func NewSupplyUsecase(supplyRepo repository.SupplyRepository, supplyHistoryRepo repository.SupplyHistoryRepository, commodityRepo repository.CommodityRepository, regionRepo repository.RegionRepository) SupplyUsecase {
	return &SupplyUsecaseImpl{
		supplyRepo:        supplyRepo,
		supplyHistoryRepo: supplyHistoryRepo,
		commodityRepo:     commodityRepo,
		regionRepo:        regionRepo,
	}
}
func (u *SupplyUsecaseImpl) CreateSupply(ctx context.Context, req *dto.SupplyCreateDTO) (*domain.Supply, error) {
	supply := domain.Supply{}
	if err := utils.ValidateStruct(req); err != nil {
		return nil, utils.NewValidationError(err)
	}
	commodity, err := u.commodityRepo.FindByID(ctx, req.CommodityID)
	if err != nil {
		return nil, utils.NewNotFoundError("commodity not found")
	}
	region, err := u.regionRepo.FindByID(ctx, req.RegionID)
	if err != nil {
		return nil, utils.NewNotFoundError("region not found")
	}
	supply.CommodityID = req.CommodityID
	supply.RegionID = req.RegionID
	supply.Quantity = req.Quantity
	supply.Commodity = commodity
	supply.Region = region
	supply.ID = uuid.New()

	err = u.supplyRepo.Create(ctx, &supply)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	createdSupply, err := u.supplyRepo.FindByID(ctx, supply.ID)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	return createdSupply, nil
}

func (u *SupplyUsecaseImpl) GetAllSupply(ctx context.Context) (*[]domain.Supply, error) {
	supplies, err := u.supplyRepo.FindAll(ctx)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	return supplies, nil
}

func (u *SupplyUsecaseImpl) GetSupplyByID(ctx context.Context, id uuid.UUID) (*domain.Supply, error) {
	supply, err := u.supplyRepo.FindByID(ctx, id)
	if err != nil {
		return nil, utils.NewInternalError("internal error")
	}
	return supply, nil
}

func (u *SupplyUsecaseImpl) GetSupplyByCommodityID(ctx context.Context, commodityID uuid.UUID) (*[]domain.Supply, error) {
	supplies, err := u.supplyRepo.FindByCommodityID(ctx, commodityID)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	return supplies, nil
}

func (u *SupplyUsecaseImpl) GetSupplyByRegionID(ctx context.Context, regionID uuid.UUID) (*[]domain.Supply, error) {
	supplies, err := u.supplyRepo.FindByRegionID(ctx, regionID)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	return supplies, nil
}

func (u *SupplyUsecaseImpl) UpdateSupply(ctx context.Context, id uuid.UUID, req *dto.SupplyUpdateDTO) (*domain.Supply, error) {
	if err := utils.ValidateStruct(req); err != nil {
		return nil, utils.NewValidationError(err)
	}
	supply, err := u.supplyRepo.FindByID(ctx, id)
	if err != nil {
		return nil, utils.NewNotFoundError("supply not found")
	}

	supplyHistory := &domain.SupplyHistory{
		ID:          uuid.New(),
		CommodityID: supply.CommodityID,
		Commodity:   supply.Commodity,
		RegionID:    supply.RegionID,
		Region:      supply.Region,
		Quantity:    supply.Quantity,
	}

	err = u.supplyHistoryRepo.Create(ctx, supplyHistory)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	supply.Quantity = req.Quantity
	err = u.supplyRepo.Update(ctx, id, supply)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	updatedSupply, err := u.supplyRepo.FindByID(ctx, id)

	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	return updatedSupply, nil
}

func (u *SupplyUsecaseImpl) DeleteSupply(ctx context.Context, id uuid.UUID) error {
	_, err := u.supplyRepo.FindByID(ctx, id)
	if err != nil {
		return utils.NewNotFoundError("supply not found")
	}
	err = u.supplyRepo.Delete(ctx, id)
	if err != nil {
		return utils.NewInternalError(err.Error())
	}
	return nil
}

func (u *SupplyUsecaseImpl) GetSupplyHistoryByCommodityIDAndRegionID(ctx context.Context, commodityID uuid.UUID, regionID uuid.UUID) (*[]domain.SupplyHistory, error) {
	supplys, err := u.supplyHistoryRepo.FindByCommodityIDAndRegionID(ctx, commodityID, regionID)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	supply, err := u.supplyRepo.FindByCommodityIDAndRegionID(ctx, commodityID, regionID)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	currentSupply := domain.SupplyHistory{
		ID:          supply.ID,
		CommodityID: supply.CommodityID,
		Commodity:   supply.Commodity,
		RegionID:    supply.RegionID,
		Region:      supply.Region,
		Quantity:    supply.Quantity,
		CreatedAt:   supply.CreatedAt,
		UpdatedAt:   supply.UpdatedAt,
		DeletedAt:   supply.DeletedAt,
	}

	allSupplyHistory := append(*supplys, currentSupply)
	return &allSupplyHistory, nil
}
