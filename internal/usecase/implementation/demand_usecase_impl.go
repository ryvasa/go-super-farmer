package usecase_implementation

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	repository_interface "github.com/ryvasa/go-super-farmer/internal/repository/interface"
	usecase_interface "github.com/ryvasa/go-super-farmer/internal/usecase/interface"
	"github.com/ryvasa/go-super-farmer/utils"
)

type DemandUsecaseImpl struct {
	demandRepo        repository_interface.DemandRepository
	demandHistoryRepo repository_interface.DemandHistoryRepository
	commodityRepo     repository_interface.CommodityRepository
	regionRepo        repository_interface.RegionRepository
}

func NewDemandUsecase(demandRepo repository_interface.DemandRepository, demandHistoryRepo repository_interface.DemandHistoryRepository, commodityRepo repository_interface.CommodityRepository, regionRepo repository_interface.RegionRepository) usecase_interface.DemandUsecase {
	return &DemandUsecaseImpl{
		demandRepo:        demandRepo,
		demandHistoryRepo: demandHistoryRepo,
		commodityRepo:     commodityRepo,
		regionRepo:        regionRepo,
	}
}

func (u *DemandUsecaseImpl) CreateDemand(ctx context.Context, req *dto.DemandCreateDTO) (*domain.Demand, error) {
	demand := domain.Demand{}
	if err := utils.ValidateStruct(req); err != nil {
		return nil, utils.NewValidationError(err)
	}
	_, err := u.commodityRepo.FindByID(ctx, req.CommodityID)
	if err != nil {
		return nil, utils.NewNotFoundError("commodity not found")
	}
	_, err = u.regionRepo.FindByID(ctx, req.RegionID)
	if err != nil {
		return nil, utils.NewNotFoundError("region not found")
	}
	demand.CommodityID = req.CommodityID
	demand.RegionID = req.RegionID
	demand.Quantity = req.Quantity
	demand.ID = uuid.New()

	err = u.demandRepo.Create(ctx, &demand)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	createdDemand, err := u.demandRepo.FindByID(ctx, demand.ID)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	return createdDemand, nil
}

func (u *DemandUsecaseImpl) GetAllDemands(ctx context.Context) (*[]domain.Demand, error) {
	demands, err := u.demandRepo.FindAll(ctx)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	return demands, nil
}
func (u *DemandUsecaseImpl) GetDemandByID(ctx context.Context, id uuid.UUID) (*domain.Demand, error) {
	demand, err := u.demandRepo.FindByID(ctx, id)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	return demand, nil
}

func (u *DemandUsecaseImpl) GetDemandsByCommodityID(ctx context.Context, commodityID uuid.UUID) (*[]domain.Demand, error) {
	_, err := u.commodityRepo.FindByID(ctx, commodityID)
	if err != nil {
		return nil, utils.NewNotFoundError("commodity not found")
	}

	demands, err := u.demandRepo.FindByCommodityID(ctx, commodityID)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	return demands, nil
}

func (u *DemandUsecaseImpl) GetDemandsByRegionID(ctx context.Context, regionID uuid.UUID) (*[]domain.Demand, error) {
	_, err := u.regionRepo.FindByID(ctx, regionID)
	if err != nil {
		return nil, utils.NewNotFoundError("region not found")
	}

	demands, err := u.demandRepo.FindByRegionID(ctx, regionID)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	return demands, nil
}

func (u *DemandUsecaseImpl) UpdateDemand(ctx context.Context, id uuid.UUID, req *dto.DemandUpdateDTO) (*domain.Demand, error) {
	if err := utils.ValidateStruct(req); err != nil {
		return nil, utils.NewValidationError(err)
	}
	demand, err := u.demandRepo.FindByID(ctx, id)
	if err != nil {
		return nil, utils.NewNotFoundError("demand not found")
	}

	demandHistory := &domain.DemandHistory{
		ID:          uuid.New(),
		CommodityID: demand.CommodityID,
		Commodity:   demand.Commodity,
		RegionID:    demand.RegionID,
		Region:      demand.Region,
		Quantity:    demand.Quantity,
	}

	err = u.demandHistoryRepo.Create(ctx, demandHistory)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	demand.Quantity = req.Quantity
	err = u.demandRepo.Update(ctx, id, demand)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	updatedDemand, err := u.demandRepo.FindByID(ctx, id)

	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	return updatedDemand, nil
}

func (u *DemandUsecaseImpl) DeleteDemand(ctx context.Context, id uuid.UUID) error {
	_, err := u.demandRepo.FindByID(ctx, id)
	if err != nil {
		return utils.NewNotFoundError("demand not found")
	}
	err = u.demandRepo.Delete(ctx, id)
	if err != nil {
		return utils.NewInternalError(err.Error())
	}
	return nil
}

func (u *DemandUsecaseImpl) GetDemandHistoryByCommodityIDAndRegionID(ctx context.Context, commodityID uuid.UUID, regionID uuid.UUID) (*[]domain.DemandHistory, error) {
	demands, err := u.demandHistoryRepo.FindByCommodityIDAndRegionID(ctx, commodityID, regionID)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	demand, err := u.demandRepo.FindByCommodityIDAndRegionID(ctx, commodityID, regionID)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	currentDemand := domain.DemandHistory{
		ID:          demand.ID,
		CommodityID: demand.CommodityID,
		Commodity:   demand.Commodity,
		RegionID:    demand.RegionID,
		Region:      demand.Region,
		Quantity:    demand.Quantity,
		CreatedAt:   demand.CreatedAt,
		UpdatedAt:   demand.UpdatedAt,
		Unit:        demand.Unit,
		DeletedAt:   demand.DeletedAt,
	}

	allDemandHistory := append(*demands, currentDemand)
	return &allDemandHistory, nil
}
